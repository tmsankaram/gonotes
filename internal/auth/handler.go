package auth

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/pquerna/otp/totp"
	"github.com/tmsankram/gonotes/internal/response"
	"github.com/tmsankram/gonotes/internal/users"
)

type Handler struct {
	users *users.Service
}

type RegisterReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	TOTP     string `json:"totp"`
}

func NewHandler(users *users.Service) *Handler {
	return &Handler{users: users}
}

// RegisterRoutes registers the public auth routes.
// These do not require authentication.
func (h *Handler) RegisterPublicRoutes(r *gin.Engine) {
	r.POST("/auth/register", h.Register)
	r.POST("/auth/login", h.Login)
}

// RegisterProtectedRoutes registers authenticated auth routes like /auth/me.
func (h *Handler) RegisterProtectedRoutes(r *gin.Engine) {
	protected := r.Group("/auth")
	protected.Use(AuthRequired())
	protected.GET("/me", func(c *gin.Context) {
		userID := c.GetUint("userID")
		u, err := h.users.GetByID(userID)
		if err != nil {
			response.Internal(c, err)
			return
		}
		c.JSON(200, gin.H{"user": u})
	})
}

func (h *Handler) Register(c *gin.Context) {
	var req RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	hash, err := users.HashPassword(req.Password)
	if err != nil {
		response.Internal(c, err)
		return
	}

	u, err := h.users.Create(users.User{
		Email:    req.Email,
		Password: hash,
	})

	if err != nil {
		response.Internal(c, err)
		return
	}

	response.Created(c, "user registered", gin.H{"id": u.ID, "email": u.Email})
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginReq

	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	u, err := h.users.GetByEmail(req.Email)
	if err != nil {
		response.Unauthorized(c, errors.New("invalid email or password"))
		return
	}

	if !users.CheckPasswordHash(req.Password, u.Password) {
		response.Unauthorized(c, errors.New("invalid email or password"))
		return
	}

	token, err := GenerateToken(u.ID)

	if err != nil {
		response.Internal(c, err)
		return
	}

	if u.TOTPEnabled {
		if req.TOTP == "" {
			response.Unauthorized(c, errors.New("TOTP required"))
			return
		}

		if !totp.Validate(req.TOTP, u.TOTPSecret) {
			response.Unauthorized(c, errors.New("invalid TOTP"))
			return
		}
	}

	response.Success(c, "login successful", gin.H{"token": token})
}
