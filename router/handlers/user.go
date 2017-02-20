package handler

import (
	"strings"
	"wooble/lib"
	model "wooble/models"

	"gopkg.in/gin-gonic/gin.v1"
)

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
		// TODO logs subscription error somewhere to keep track
		c.Error(err)
	}

	c.Header("Location", "/token/generate")

	c.JSON(Created, nil)
}
