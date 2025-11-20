package auth

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/tmsankram/gonotes/internal/response"
	"github.com/tmsankram/gonotes/internal/users"
)

type Handler struct {
	users *users.Service
}

func NewHandler(users *users.Service) *Handler {
	return &Handler{users: users}
}

type RegisterReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
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

type LoginReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
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

	response.Success(c, "login successful", gin.H{"token": token})
}
