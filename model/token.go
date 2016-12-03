package model

import (
	"fmt"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
)

type Token struct {
	Token string `json:"token"`
}

type CustomClaims struct {
	Name string `json:"name"`
	jwt.StandardClaims
}

func NewToken(id uint64, name string) (*Token, error) {
	claims := &CustomClaims{
		name,
		jwt.StandardClaims{
			Issuer:  "wooble.io", // TODO
			Subject: fmt.Sprintf("%v", id),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tString, err := token.SignedString(TokenKey())
	if err != nil {
		return nil, err
	}
	return &Token{
		tString,
	}, nil
}

func UserByToken(token interface{}) (*User, error) {
	claims := token.(*jwt.Token).Claims.(jwt.MapClaims)

	idStr := claims["sub"]
	id, err := strconv.ParseUint(idStr.(string), 10, 640)

	name := claims["name"]

	return &User{
		ID:   id,
		Name: name.(string),
	}, err
}

func TokenKey() []byte {
	return []byte(viper.GetString("token_key"))
}
