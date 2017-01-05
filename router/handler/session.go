package handler

import (
	"database/sql"

	jwt "github.com/dgrijalva/jwt-go"

	"wooble/model"
	"wooble/router/helper"

	"gopkg.in/gin-gonic/gin.v1"
)

// GenerateToken generates a new token
func GenerateToken(c *gin.Context) {
	type CredsForm struct {
		Login  string `json:"login" binding:"required"`
		Passwd string `json:"secret"`
	}

	var form CredsForm

	res := NewRes()

	c.Header("Content-Type", gin.MIMEJSON)
	if c.BindJSON(&form) != nil {
		res.Error(ErrBadForm, "login (string) is required")
		c.JSON(res.HTTPStatus(), res)
		return
	}

	user, err := model.UserByLogin(form.Login)
	if err != nil {
		if err == sql.ErrNoRows {
			res.Error(ErrBadCreds, "Username or email do not exist")
		} else {
			res.Error(ErrDBSelect)
		}
		c.JSON(res.HTTPStatus(), res)
		return
	}

	if user.IsPasswordValid(form.Passwd) {
		token := model.NewToken(user, "")
		tokenS, err := token.SignedString(model.TokenKey())

		if err != nil {
			res.Error(ErrServ, "token generation")
			c.JSON(res.HTTPStatus(), res)
			return
		}

		res.Response(tokenS)
	} else {
		res.Error(ErrBadCreds, "Password invalid")
	}

	res.Status = Created

	c.JSON(res.HTTPStatus(), res)
}

// RefreshToken refreshes a token (lifetime is given in the server conf file $CONFPATH)
func RefreshToken(c *gin.Context) {
	res := NewRes()

	token, err := helper.ParseToken(c)

	if ve, ok := err.(*jwt.ValidationError); ok && model.IsTokenExpired(ve) {
		newToken, err := model.RefreshToken(token)
		if err != nil {
			res.Error(ErrServ, "token refresh")
		}

		res.Response(newToken)

		res.Status = Created
	} else if ve != nil {
		res.Error(ErrServ, "token refresh")
	}

	c.JSON(res.HTTPStatus(), res)
}
