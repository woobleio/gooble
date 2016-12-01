package model

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"

	"github.com/spf13/viper"
)

type Token struct {
	Token string `json:"token"`
}

func NewToken(id uint64) (*Token, error) {
	claims := &jwt.StandardClaims{
		Issuer:  "TODO",
		Subject: fmt.Sprintf("%v", id),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tString, err := token.SignedString([]byte(viper.GetString("token_key")))
	if err != nil {
		return nil, err
	}
	return &Token{
		tString,
	}, nil
}
