package model

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"time"

	"golang.org/x/crypto/scrypt"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
)

type Token struct {
	Token string `json:"token"`
}

type CustomClaims struct {
	Name         string `json:"name"`
	RefreshToken string `json:"refresh_token"`
	jwt.StandardClaims
}

func NewToken(user *User, refreshToken string) (*Token, error) {
	if refreshToken == "" {
		refreshToken = genRefreshToken(user)
	}
	claims := &CustomClaims{
		user.Name,
		refreshToken,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 1).Unix(), // TODO * 30
			IssuedAt:  time.Now().Unix(),
			Issuer:    "wooble.io", // TODO
			Subject:   fmt.Sprintf("%v", user.ID),
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

func IsTokenExpired(err *jwt.ValidationError) bool {
	return err.Errors&jwt.ValidationErrorExpired != 0
}

func IsTokenMalformed(err *jwt.ValidationError) bool {
	return err.Errors&jwt.ValidationErrorMalformed != 0
}

func RefreshToken(tokenStr string) (*Token, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return TokenKey(), nil
	})

	if ve, ok := err.(*jwt.ValidationError); ok {
		if IsTokenMalformed(ve) {
			return nil, err
		} else if IsTokenExpired(ve) { // Refresh only if token has expired
			user, err := UserByToken(token)
			if err != nil {
				return nil, err
			}

			user, err = UserByID(user.ID)
			if err != nil {
				return nil, err
			}

			validToken := genRefreshToken(user)
			if validToken != token.Claims.(jwt.MapClaims)["refresh_token"] {
				return nil, errors.New(fmt.Sprintf("Invalid refresh token %s", validToken))
			}

			iat := token.Claims.(jwt.MapClaims)["iat"].(float64)
			iatUnix := int64(iat)

			tokenOld := time.Now().Sub(time.Unix(iatUnix, 0))
			// Token expired since 7 days ago
			if tokenOld.Hours() >= (24 * 7) {
				return nil, errors.New(fmt.Sprintf("Token has expired, re-auth required : %s", token.Raw))
			}

			return NewToken(user, validToken)
		}
	}

	if err != nil {
		return nil, err
	}

	return nil, errors.New("Failed to refresh the token, it is still valid")
}

func TokenKey() []byte {
	return []byte(viper.GetString("token_key"))
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

func genRefreshToken(user *User) string {
	id := fmt.Sprintf("%v", user.ID)
	cp, err := scrypt.Key([]byte(user.Salt+id), []byte(TokenKey()), 16384, 8, 1, 32)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(cp)
}
