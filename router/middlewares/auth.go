package middleware

import (
	"errors"
	"wooble/models"
	"wooble/router/helpers"

	jwt_lib "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// Authenticate is a handler middleware that authorizes a token (header field Authorization)
func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := helper.ParseToken(c)

		if token == nil {
			abort(c, errors.New("Token invalid or missing"))
			return
		}

		// If token has expired, refresh it and returns it in the header
		if ve, ok := err.(*jwt_lib.ValidationError); ok {
			// check signature
			if model.IsTokenInvalid(ve) {
				abort(c, err)
				return
			}

			if model.IsTokenExpired(ve) {
				if newToken, refTokenErr := model.RefreshToken(token); refTokenErr == nil {
					var tokenRaw string
					tokenRaw, _ = newToken.SignedString(model.TokenKey())
					token, err = jwt_lib.Parse(tokenRaw, func(token *jwt_lib.Token) (interface{}, error) {
						return model.TokenKey(), nil
					})
					c.Header("Authorization", tokenRaw)
				}
			}
		}

		// check other errors
		if err != nil {
			abort(c, err)
			return
		}

		user, err := model.UserByToken(token)
		if err != nil {
			// Invalid usertoken
			abort(c, err)
			return
		}

		c.Set("user", user)
		c.Next()
	}
}

func abort(c *gin.Context, err error) {
	c.Header("Location", "/signin")
	c.AbortWithError(401, err)
}
