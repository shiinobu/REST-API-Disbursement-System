package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"rest-api-disbursement-system/internal/services"
)

type AuthHandler struct {
	auth services.AuthService
}

type loginRequest struct {
	Username string `json:"username" binding:"required,min=3,max=100"`
	Password string `json:"password" binding:"required,min=8,max=72"`
}

func NewAuthHandler(auth services.AuthService) *AuthHandler {
	return &AuthHandler{auth: auth}
}

func (h *AuthHandler) MountRoutes(router *gin.RouterGroup) {
	router.POST("/login", h.Login)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var request loginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		ValidationError(c, err, request)
		return
	}

	response, err := h.auth.Login(services.LoginInput{
		Username: request.Username,
		Password: request.Password,
	})
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			Error(c, http.StatusUnauthorized, "Username atau password salah", gin.H{"credentials": "username atau password salah"})
			return
		}
		Error(c, http.StatusInternalServerError, "Login gagal", nil)
		return
	}

	Success(c, http.StatusOK, "Login berhasil", response)
}
