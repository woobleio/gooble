package handler

import (
	"database/sql"

	jwt "github.com/dgrijalva/jwt-go"

	"wooble/models"
	"wooble/router/helpers"

	"gopkg.in/gin-gonic/gin.v1"
)

// GenerateToken generates a new token
func GenerateToken(c *gin.Context) {
	type CredsForm struct {
		Email  string `json:"email" validate:"required"`
		Secret string `json:"secret" validate:"required"`
	}

	var form CredsForm

	if err := c.BindJSON(&form); err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind).SetMeta(ErrBadForm)
		return
	}

	user, err := model.UserByEmail(form.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			c.Error(err).SetMeta(ErrBadCreds)
		} else {
			c.Error(err).SetMeta(ErrDB)
		}
		return
	}

	res := NewRes(nil)

	if user.IsPasswordValid(form.Secret) {
		token := model.NewToken(user, "")
		tokenS, err := token.SignedString(model.TokenKey())

		if err != nil {
			c.Error(err).SetMeta(ErrIntServ)
			return
		}

		res.Response(struct {
			Token string `json:"token"`
		}{Token: tokenS})

	} else {
		c.Error(nil).SetMeta(ErrBadCreds)
		return
	}

	c.JSON(Created, res)
}

// RefreshToken refreshes a token (lifetime is given in the server conf file $CONFPATH)
func RefreshToken(c *gin.Context) {
	res := NewRes(nil)

	token, err := helper.ParseToken(c)

	if ve, ok := err.(*jwt.ValidationError); ok && model.IsTokenExpired(ve) {
		newToken, err := model.RefreshToken(token)
		if err != nil {
			c.Error(err).SetMeta(ErrServ.SetParams("source", "token"))
			return
		}

		tokenRaw, _ := newToken.SignedString(model.TokenKey())
		res.Response(struct {
			Token string `json:"token"`
		}{Token: tokenRaw})
	} else if ve != nil {
		c.Error(ve).SetMeta(ErrServ.SetParams("source", "token"))
		return
	}

	c.JSON(Created, res)
}
