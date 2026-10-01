package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ramdhanrizkij/next-store-api/internal/modules/identity/application"
	appErrors "github.com/ramdhanrizkij/next-store-api/internal/shared/errors"
	"github.com/ramdhanrizkij/next-store-api/internal/shared/pagination"
	"github.com/ramdhanrizkij/next-store-api/internal/shared/response"
	"github.com/ramdhanrizkij/next-store-api/internal/shared/validator"
)

type UserHandler struct {
	userService application.UserService
}

func NewUserHandler(userService application.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetMe retrieves the authenticated user's profile
func (h *UserHandler) GetMe(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		response.Error(c, http.StatusUnauthorized, "Unauthorized", "User ID not found in context")
		return
	}

	user, err := h.userService.GetByID(c.Request.Context(), userID)
	if err != nil {
		response.FromAppError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "User profile retrieved successfully", user)
}

// GetByID retrieves a user by ID
func (h *UserHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	user, err := h.userService.GetByID(c.Request.Context(), id)
	if err != nil {
		response.FromAppError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "User retrieved successfully", user)
}

// List retrieves paginated users
func (h *UserHandler) List(c *gin.Context) {
	p := pagination.GetPagination(c)
	users, meta, err := h.userService.List(c.Request.Context(), p)
	if err != nil {
		response.FromAppError(c, err)
		return
	}

	response.SuccessWithMeta(c, http.StatusOK, "Users retrieved successfully", users, meta)
}

// Update updates the current authenticated user's profile
func (h *UserHandler) Update(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		response.Error(c, http.StatusUnauthorized, "Unauthorized", "User ID not found in context")
		return
	}

	var req application.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Validation error", validator.FormatValidationError(err))
		return
	}

	user, err := h.userService.Update(c.Request.Context(), userID, &req)
	if err != nil {
		response.FromAppError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "User updated successfully", user)
}

// Delete removes a user by ID
func (h *UserHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.userService.Delete(c.Request.Context(), id); err != nil {
		if _, ok := appErrors.IsAppError(err); ok {
			response.FromAppError(c, err)
			return
		}
		response.Error(c, http.StatusInternalServerError, "Failed to delete user", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "User deleted successfully", nil)
}
