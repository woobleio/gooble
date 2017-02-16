package handler

import (
	"database/sql"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"

	"wooble/lib"
	"wooble/models"
	"wooble/router/helpers"

	"gopkg.in/gin-gonic/gin.v1"
)

// GenerateToken generates a new token
func GenerateToken(c *gin.Context) {
	type CredsForm struct {
		Email  string `json:"email" binding:"required"`
		Secret string `json:"secret" binding:"required"`
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
			c.Error(err).SetMeta(ErrDBSelect)
		}
		return
	}

	res := NewRes(nil)

	if user.IsPasswordValid(form.Secret) {
		token := model.NewToken(user, "")
		tokenS, err := token.SignedString(model.TokenKey())

		if err != nil {
			c.Error(err).SetMeta(ErrServ)
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
			c.Error(err).SetMeta(ErrServ)
			return
		}

		tokenRaw, _ := newToken.SignedString(model.TokenKey())
		res.Response(struct {
			Token string `json:"token"`
		}{Token: tokenRaw})
	} else if ve != nil {
		c.Error(ve).SetMeta(ErrServ)
		return
	}

	c.JSON(Created, res)
}

// SignUp saves a new user in the database
func SignUp(c *gin.Context) {
	var data model.UserForm

	if err := c.BindJSON(&data); err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind).SetMeta(ErrBadForm)
		return
	}

	uID, err := model.NewUser(&data)
	if err != nil {
		c.Error(err).SetMeta(ErrDBSave)
		return
	}

	// Saves the customer in third party (Stripe for now)
	customer, errCust := lib.NewCustomer(data.Email, data.Plan, data.CardToken)
	if errCust != nil {
		model.DeleteUser(uID)
		c.Error(errCust).SetMeta(ErrDBSave)
		return
	}

	data.CustomerID = customer.ID

	// Sets customer id to User
	model.UpdateCustomerID(uID, customer.ID)

	// Logs customer subscription in the DB
	if _, err := model.NewPlanUser(uID, strings.Split(data.Plan, "_")[0], customer.Subs.Values[0].PeriodEnd); err != nil {
		// TODO logs subscription error somewhere to keep track
		c.Error(err)
	}

	c.Header("Location", "/token/generate")

	c.JSON(Created, nil)
}
