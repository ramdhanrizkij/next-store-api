package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ramdhanrizkij/next-store-api/internal/modules/auth/application"
	"github.com/ramdhanrizkij/next-store-api/internal/shared/response"
	"github.com/ramdhanrizkij/next-store-api/internal/shared/validator"
)

type AuthHandler struct {
	authService application.AuthService
}

func NewAuthHandler(authService application.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req application.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Validation error", validator.FormatValidationError(err))
		return
	}

	res, err := h.authService.Register(c.Request.Context(), &req)
	if err != nil {
		response.FromAppError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "User registered successfully", res)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req application.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Validation error", validator.FormatValidationError(err))
		return
	}

	res, err := h.authService.Login(c.Request.Context(), &req)
	if err != nil {
		response.FromAppError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Login successful", res)
}
