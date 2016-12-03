package handler

import (
	"database/sql"

	"wooble/model"

	"gopkg.in/gin-gonic/gin.v1"
)

func SignIn(c *gin.Context) {
	type SigninForm struct {
		Login  string `json:"login" binding:"required"`
		Passwd string `json:"passwd"`
	}

	var form SigninForm

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
		token, err := model.NewToken(user.ID, user.Name)
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

func SignUp(c *gin.Context) {
	var data model.User

	res := NewRes()

	// FIXME workaroun gin issue with Bind (https://github.com/gin-gonic/gin/issues/633)
	c.Header("Content-Type", gin.MIMEJSON)
	if c.BindJSON(&data) != nil {
		res.Error(ErrBadForm, "username (string) and email (string) are required")
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
