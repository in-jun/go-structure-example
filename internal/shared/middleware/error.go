package middleware

import (
	stderrors "errors"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/in-jun/go-structure-example/internal/shared/errors"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err
		var customErr errors.CustomError
		if stderrors.As(err, &customErr) {
			if customErr.Status >= 500 {
				slog.Error("internal error", "message", customErr.Message)
				customErr.Message = "Internal Server Error"
			}
			c.JSON(customErr.Status, customErr)
			return
		}

		slog.Error("unhandled error", "error", err)
		c.JSON(500, errors.Internal("Internal Server Error"))
	}
}
