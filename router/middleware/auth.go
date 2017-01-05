package middleware

import (
	"wooble/model"
	"wooble/router/helper"

	jwt_lib "github.com/dgrijalva/jwt-go"
	"gopkg.in/gin-gonic/gin.v1"
)

// Authenticate is a handler middleware that authorizes a token (header field Authorization)
func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := helper.ParseToken(c)

		if ve, ok := err.(*jwt_lib.ValidationError); ok && model.IsTokenExpired(ve) {
			if newToken, refTokenErr := model.RefreshToken(token); refTokenErr == nil {
				var tokenRaw string
				tokenRaw, _ = newToken.SignedString(model.TokenKey())
				token, err = jwt_lib.Parse(tokenRaw, func(token *jwt_lib.Token) (interface{}, error) {
					return model.TokenKey(), nil
				})
				c.Header("Authorization", tokenRaw)
			}
		}

		if !token.Valid || err != nil {
			c.Header("Location", "/signin")
			c.AbortWithError(401, err)
			return
		}

		user, err := model.UserByToken(token)
		if err != nil {
			// Invalid token
			c.Header("Location", "/signin")
			c.AbortWithError(401, err)
		}

		c.Set("user", user)
		c.Next()
	}
}
