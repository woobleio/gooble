package helper

import (
	"wooble/models"

	jwt_lib "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/gin-gonic/gin"
)

// ParseToken parses a token from a request
func ParseToken(c *gin.Context) (*jwt_lib.Token, error) {
	return request.ParseFromRequest(c.Request, request.OAuth2Extractor, func(token *jwt_lib.Token) (interface{}, error) {
		return model.TokenKey(), nil
	})
}
