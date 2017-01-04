package helper

import (
	"wooble/model"

	jwt_lib "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"gopkg.in/gin-gonic/gin.v1"
)

func ParseToken(c *gin.Context) (*jwt_lib.Token, error) {
	return request.ParseFromRequest(c.Request, request.OAuth2Extractor, func(token *jwt_lib.Token) (interface{}, error) {
		return model.TokenKey(), nil
	})
}