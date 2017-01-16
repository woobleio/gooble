package model

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"time"
	"wooble/lib"

	"golang.org/x/crypto/scrypt"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
)

// CustomClaims is wooble token claims
type CustomClaims struct {
	Name         string `json:"name"`
	RefreshToken string `json:"refresh_token"`
	jwt.StandardClaims
}

// NewToken generates a new token
func NewToken(user *User, refreshToken string) *jwt.Token {
	if refreshToken == "" {
		refreshToken = genRefreshToken(user)
	}
	claims := &CustomClaims{
		user.Name,
		refreshToken,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * time.Duration(lib.GetTokenLifetime())).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "wooble.io", // TODO
			Subject:   fmt.Sprintf("%v", user.ID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token
}

// IsTokenInvalid tells if the token signature is wrong
func IsTokenInvalid(err *jwt.ValidationError) bool {
	return err.Errors&jwt.ValidationErrorSignatureInvalid != 0
}

// IsTokenExpired tells if a token has expired (claim exp)
func IsTokenExpired(err *jwt.ValidationError) bool {
	return err.Errors&jwt.ValidationErrorExpired != 0
}

// IsTokenMalformed tells if a token is not JWT standard
func IsTokenMalformed(err *jwt.ValidationError) bool {
	return err.Errors&jwt.ValidationErrorMalformed != 0
}

// RefreshToken refreshes old token (depending on exp claim)
func RefreshToken(token *jwt.Token) (*jwt.Token, error) {
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
		return nil, fmt.Errorf("Invalid refresh token %s", validToken)
	}

	iat := token.Claims.(jwt.MapClaims)["iat"].(float64)
	iatUnix := int64(iat)

	tokenOld := time.Now().Sub(time.Unix(iatUnix, 0))
	// Lifelong refresh token 1 month
	if tokenOld.Hours() >= (24 * 30) {
		return nil, fmt.Errorf("Token has expired, re-auth required : %s", token.Raw)
	}

	return NewToken(user, validToken), nil
}

// TokenKey returns token private key
func TokenKey() []byte {
	return []byte(viper.GetString("token_key"))
}

// UserByToken return the user whom the token belong to
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
