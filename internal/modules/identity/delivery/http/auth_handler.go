package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ramdhanrizkij/next-store-api/internal/modules/identity/application"
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

	response.Success(c, http.StatusCreated, "User registered successfully. Please verify your email.", res)
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

func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	var req application.VerifyEmailRequest

	// Support token from either query parameter (?token=...) or JSON body ({"token": "..."})
	tokenQuery := c.Query("token")
	if tokenQuery != "" {
		req.Token = tokenQuery
	} else {
		if err := c.ShouldBindJSON(&req); err != nil {
			response.Error(c, http.StatusBadRequest, "Validation error", validator.FormatValidationError(err))
			return
		}
	}

	if err := h.authService.VerifyEmail(c.Request.Context(), &req); err != nil {
		response.FromAppError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Email successfully verified. Your account is now active.", nil)
}
