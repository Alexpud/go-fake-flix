package apierrors

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err
		var appErr *AppError
		if errors.As(err, &appErr) {
			fmt.Println("AppError:", *appErr)
			c.JSON(http.StatusBadRequest, ProblemDetails{
				Type:   "https://datatracker.ietf.org/doc/html/rfc7807",
				Title:  "An unexpected error occurred  SDA",
				Status: http.StatusBadRequest,
				Detail: err.Error(),
				Extras: map[string]any{"code": appErr.Code},
			})
			return
		}

		c.JSON(http.StatusInternalServerError, ProblemDetails{
			Type:   "https://datatracker.ietf.org/doc/html/rfc7807",
			Title:  "An unexpected error occurred ADSASD",
			Status: http.StatusInternalServerError,
			Detail: err.Error(),
			Extras: map[string]any{"code": "INTERNAL"},
		})
	}
}
