package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	appErrors "github.com/ramdhanrizkij/next-store-api/internal/shared/errors"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

func Success(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func SuccessWithMeta(c *gin.Context, statusCode int, message string, data interface{}, meta interface{}) {
	c.JSON(statusCode, Response{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

func Error(c *gin.Context, statusCode int, message string, errs interface{}) {
	c.JSON(statusCode, Response{
		Success: false,
		Message: message,
		Errors:  errs,
	})
}

func FromAppError(c *gin.Context, err error) {
	if appErr, ok := appErrors.IsAppError(err); ok {
		c.JSON(appErr.StatusCode, Response{
			Success: false,
			Message: appErr.Message,
			Errors:  appErr.Code,
		})
		return
	}

	c.JSON(http.StatusInternalServerError, Response{
		Success: false,
		Message: "Internal server error",
		Errors:  err.Error(),
	})
}
