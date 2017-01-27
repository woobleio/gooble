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
		Email  string `json:"email" binding:"required"`
		Secret string `json:"secret" binding:"required"`
	}

	var form CredsForm

	res := NewRes()

	c.Header("Content-Type", gin.MIMEJSON)
	if c.BindJSON(&form) != nil {
		res.Error(ErrBadForm, "email (string) and secret (string) are required")
		c.JSON(res.HTTPStatus(), res)
		return
	}

	user, err := model.UserByEmail(form.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			res.Error(ErrBadCreds, "Email does not exist")
		} else {
			res.Error(ErrDBSelect)
		}
		c.JSON(res.HTTPStatus(), res)
		return
	}

	if user.IsPasswordValid(form.Secret) {
		token := model.NewToken(user, "")
		tokenS, err := token.SignedString(model.TokenKey())

		if err != nil {
			res.Error(ErrServ, "token generation")
			c.JSON(res.HTTPStatus(), res)
			return
		}

		res.Response(struct {
			Token string `json:"token"`
		}{Token: tokenS})

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

		tokenRaw, _ := newToken.SignedString(model.TokenKey())
		res.Response(struct {
			Token string `json:"token"`
		}{Token: tokenRaw})

		res.Status = Created
	} else if ve != nil {
		res.Error(ErrServ, "token refresh")
	}

	c.JSON(res.HTTPStatus(), res)
}

// SignUp saves a new user in the database
func SignUp(c *gin.Context) {
	var data model.UserForm

	res := NewRes()

	// FIXME workaroun gin issue with Bind (https://github.com/gin-gonic/gin/issues/633)
	c.Header("Content-Type", gin.MIMEJSON)
	if c.BindJSON(&data) != nil {
		res.Error(ErrBadForm, "name (string), email (string) and secret (string) are required")
		c.JSON(res.HTTPStatus(), res)
		return
	}

	_, err := model.NewUser(&data)
	if err != nil {
		res.Error(ErrDBSave, "- Name should be unique\n - Email should be unique")
	} else {
		c.Header("Location", "/token/generate")
	}

	res.Status = Created

	c.JSON(res.HTTPStatus(), res)
}
