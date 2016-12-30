package handler

import (
	"database/sql"

	"wooble/model"

	"gopkg.in/gin-gonic/gin.v1"
)

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
		c.JSON(res.HttpStatus(), res)
		return
	}

	user, err := model.UserByLogin(form.Login)
	if err != nil {
		if err == sql.ErrNoRows {
			res.Error(ErrBadCreds, "Username or email do not exist")
		} else {
			res.Error(ErrDBSelect)
		}
		c.JSON(res.HttpStatus(), res)
		return
	}

	if user.IsPasswordValid(form.Passwd) {
		token, err := model.NewToken(user, "")
		if err != nil {
			res.Error(ErrServ, "token generation")
			c.JSON(res.HttpStatus(), res)
			return
		}
		res.Response(&token)
	} else {
		res.Error(ErrBadCreds, "Password invalid")
	}

	res.Status = Created

	c.JSON(res.HttpStatus(), res)
}

func RefreshToken(c *gin.Context) {
	type TokenForm struct {
		Token string `json:"accessToken" binding:"required"`
	}

	var form TokenForm

	res := NewRes()

	c.Header("Content-Type", gin.MIMEJSON)
	if c.BindJSON(&form) != nil {
		res.Error(ErrBadForm, "accessToken (string) is required")
		c.JSON(res.HttpStatus(), res)
		return
	}

	token, err := model.RefreshToken(form.Token)
	if err != nil {
		res.Error(ErrServ, "token refresh")
	}

	res.Response(&token)

	res.Status = Created

	c.JSON(res.HttpStatus(), res)
}

func SignUp(c *gin.Context) {
	var data model.User

	res := NewRes()

	// FIXME workaroun gin issue with Bind (https://github.com/gin-gonic/gin/issues/633)
	c.Header("Content-Type", gin.MIMEJSON)
	if c.BindJSON(&data) != nil {
		res.Error(ErrBadForm, "name (string), email (string) and secret (string) are required")
		c.JSON(res.HttpStatus(), res)
		return
	}

	_, err := model.NewUser(&data)
	if err != nil {
		res.Error(ErrDBSave, "- Name should be unique")
	} else {
		c.Header("Location", "/signin")
	}

	res.Status = Created

	c.JSON(res.HttpStatus(), res)
}
