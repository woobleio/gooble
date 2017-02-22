package handler

import (
	"fmt"
	"strings"
	"wooble/lib"
	model "wooble/models"
	helper "wooble/router/helpers"

	"gopkg.in/gin-gonic/gin.v1"
)

// GETUsers returns one users with private infos if authenticated, with public infos if not
func GETUsers(c *gin.Context) {
	var err error
	data := new(model.User)

	token, _ := helper.ParseToken(c)

	tokenUser, _ := model.UserByToken(token)

	username := c.Param("username")

	if tokenUser != nil && username == tokenUser.Name {
		data, err = model.UserPrivateByID(tokenUser.ID)
	} else {
		data, err = model.UserPublicByName(c.Param("username"))
	}

	if err != nil {
		c.Error(err).SetMeta(ErrResNotFound.SetParams("source", "user", "name", c.Param("username")))
		return
	}

	c.JSON(OK, data)
}

// POSTUsers saves a new user in the database
func POSTUsers(c *gin.Context) {
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
		// TODO when fail stripe should't charge
		c.Error(err).SetMeta(ErrIntServ)
		return
	}

	c.Header("Location", "/token/generate")

	c.JSON(Created, nil)
}

// UpdatePassword update authenticated user's password
func UpdatePassword(c *gin.Context) {
	var passwordForm struct {
		OldSecret string `json:"oldSecret" binding:"required"`
		NewSecret string `json:"newSecret" binding:"required"`
	}

	if err := c.BindJSON(&passwordForm); err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind).SetMeta(ErrBadForm)
		return
	}

	user, _ := c.Get("user")

	if !user.(*model.User).IsPasswordValid(passwordForm.OldSecret) {
		fmt.Println("invalid")
		return
	}

	if err := model.UpdateUserPassword(user.(*model.User).ID, passwordForm.NewSecret); err != nil {
		c.Error(err).SetMeta(ErrUpdate.SetParams("source", "user", "name", user.(*model.User).Name))
		return
	}

	c.Header("Location", "/users/me")

	c.JSON(OK, nil)
}
