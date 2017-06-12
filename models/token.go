package model

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"time"
	"wooble/lib"

	"golang.org/x/crypto/scrypt"

	"github.com/dgrijalva/jwt-go"
)

// CustomClaims is wooble token claims
type CustomClaims struct {
	Name         string `json:"name"`
	Plan         Plan   `json:"plan"`
	RefreshToken string `json:"refresh_token"`
	jwt.StandardClaims
}

// NewToken generates a new token
func NewToken(user *User, refreshToken string) *jwt.Token {
	if refreshToken == "" {
		refreshToken = genRefreshToken(user)
	}

	// Set default plan if the user has no plan (this isn't an exceptionnal case)
	if user.Plan.Label == nil || !user.Plan.Label.Valid {
		user.Plan, _ = DefaultPlan(user.ID)
	}

	claims := &CustomClaims{
		user.Name,
		*user.Plan,
		refreshToken,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * time.Duration(lib.GetTokenLifetime())).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    lib.GetOrigins()[0],
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

	user, err = UserPrivateByID(user.ID)
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
	// Lifelong refresh token 2 weeks
	if tokenOld.Hours() >= (24 * 14) {
		return nil, fmt.Errorf("Token has expired, re-auth required : %s", token.Raw)
	}

	return NewToken(user, validToken), nil
}

// TokenKey returns token private key
func TokenKey() []byte {
	return []byte(lib.GetTokenKey())
}

// UserByToken return the user whom the token belong to
func UserByToken(token interface{}) (*User, error) {
	if token.(*jwt.Token) == nil {
		return nil, errors.New("No tokens")
	}
	claims := token.(*jwt.Token).Claims.(jwt.MapClaims)

	idStr := claims["sub"]
	userID, err := strconv.ParseUint(idStr.(string), 10, 64)

	if err != nil {
		return nil, err
	}

	name := claims["name"]
	planInf := claims["plan"].(map[string]interface{})

	labelSrc := planInf["label"]
	nbPkgSrc := planInf["nbPkg"]
	nbCreaSrc := planInf["nbCrea"]

	layout := "2006-01-02T15:04:05Z"
	var startDate time.Time
	var endDate time.Time
	var dateErr error

	switch planInf["startDate"].(type) {
	case interface{}:
		startDate, _ = time.Parse(layout, time.Now().String())
	case string:
		startDate, dateErr = time.Parse(layout, planInf["startDate"].(string))
		if dateErr != nil {
			return nil, dateErr
		}
	}

	switch planInf["endDate"].(type) {
	case interface{}:
		endDate, _ = time.Parse(layout, time.Now().String())
	case string:
		endDate, dateErr = time.Parse(layout, planInf["endDate"].(string))
		if dateErr != nil {
			return nil, dateErr
		}
	}

	nbPkg, okPkg := nbPkgSrc.(float64)
	nbCrea, okCrea := nbCreaSrc.(float64)

	if !okPkg || !okCrea {
		return nil, errors.New("Parsing error on nbPkg or nbCrea")
	}

	plan := &Plan{
		Label:     lib.InitNullString(labelSrc.(string)),
		NbPkg:     lib.InitNullInt64(int64(nbPkg)),
		NbCrea:    lib.InitNullInt64(int64(nbCrea)),
		StartDate: lib.InitNullTime(startDate),
		EndDate:   lib.InitNullTime(endDate),
	}

	// DefaultPlan if the current plan has exprired
	if plan.HasExpired() {
		plan, _ = DefaultPlan(userID)
	}

	return &User{
		ID:   userID,
		Name: name.(string),
		Plan: plan,
	}, nil
}

func genRefreshToken(user *User) string {
	id := fmt.Sprintf("%v", user.ID)
	cp, err := scrypt.Key([]byte(user.Salt+id), []byte(TokenKey()), 16384, 8, 1, 32)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(cp)
}
