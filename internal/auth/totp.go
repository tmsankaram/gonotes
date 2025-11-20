package auth

import (
	"encoding/base64"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"github.com/skip2/go-qrcode"

	"github.com/tmsankram/gonotes/internal/response"
	"github.com/tmsankram/gonotes/internal/users"
)

type TOTPHandler struct {
	users *users.Service
}
type VerifyReq struct {
	Token string `json:"token" binding:"required"`
}

func NewTOTPHandler(user *users.Service) *TOTPHandler {
	return &TOTPHandler{users: user}
}

// RegisterRoutes registers the protected TOTP routes.
func (h *TOTPHandler) RegisterRoutes(r *gin.Engine) {
	authProtected := r.Group("/auth")
	authProtected.Use(AuthRequired())
	authProtected.POST("/totp/enable", h.EnableTOTP)
	authProtected.POST("/totp/verify", h.VerifyTOTP)
}

func (h *TOTPHandler) EnableTOTP(c *gin.Context) {
	userID := c.GetUint("userID")

	// Generate secret
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "GoNotes",
		AccountName: c.GetString("userEmail"),
		Algorithm:   otp.AlgorithmSHA1,
	})
	if err != nil {
		response.Internal(c, err)
		return
	}

	// store secret in DB
	err = h.users.SetTOTP(userID, key.Secret())
	if err != nil {
		response.Internal(c, err)
	}

	img, err := qrcode.Encode(key.URL(), qrcode.Medium, 256)
	if err != nil {
		response.Internal(c, err)
		return
	}
	response.Success(c, "TOTP enabled", gin.H{
		"secret": key.Secret(),
		"uri":    key.URL(),
		"qr":     "data:image/png;base64," + base64.StdEncoding.EncodeToString(img),
	})
}

func (h *TOTPHandler) VerifyTOTP(c *gin.Context) {
	userID := c.GetUint("userID")

	u, err := h.users.GetByID(userID)
	if err != nil {
		response.NotFound(c, err)
		return
	}

	if u.TOTPSecret == "" {
		response.BadRequest(c, errors.New("TOTP not enabled"))
		return
	}

	var req VerifyReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	ok := totp.Validate(req.Token, u.TOTPSecret)
	if !ok {
		response.Unauthorized(c, errors.New("invalid TOTP code"))
		return
	}
	response.Success(c, "TOTP verified", nil)
}
