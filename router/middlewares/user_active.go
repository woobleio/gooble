package middleware

import (
	"errors"
	model "wooble/models"

	"github.com/gin-gonic/gin"
)

// IsActive verifies if a user is activated
func IsActive() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, _ := c.Get("user")

		if !user.(*model.User).IsActive {
			c.Header("Location", "/")
			c.AbortWithError(401, errors.New("User not activated"))
			return
		}

		c.Next()
	}
}
