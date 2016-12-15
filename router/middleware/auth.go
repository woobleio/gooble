package middleware

import (
	"fmt"
	"wooble/model"

	jwt_lib "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"gopkg.in/gin-gonic/gin.v1"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := request.ParseFromRequest(c.Request, request.OAuth2Extractor, func(token *jwt_lib.Token) (interface{}, error) {
			return model.TokenKey(), nil
		})

		// TODO Oauth spec
		if ve, ok := err.(*jwt_lib.ValidationError); ok {
			if model.IsTokenMalformed(ve) {
				fmt.Println("That's not even a token")
			} else if model.IsTokenExpired(ve) {
				// Token is either expired
				fmt.Println("Timing is everything")
			}
		}

		if err != nil || !token.Valid {
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
