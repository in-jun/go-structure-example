package middleware

import (
	stderrors "errors"

	"github.com/gin-gonic/gin"
	"github.com/in-jun/go-structure-example/internal/shared/errors"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			var customErr errors.CustomError
			if stderrors.As(err, &customErr) {
				c.JSON(customErr.Status, customErr)
				return
			}
			c.JSON(500, errors.Internal("Internal Server Error"))
		}
	}
}
