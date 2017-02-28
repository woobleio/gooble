package handler

import (
	"fmt"
	"strings"
	form "wooble/forms"
	"wooble/lib"
	model "wooble/models"
	helper "wooble/router/helpers"

	"gopkg.in/gin-gonic/gin.v1"
)

// GETUser returns one users with private infos if authenticated, with public infos if not
func GETUser(c *gin.Context) {
	var data *model.User
	var err error

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

	c.JSON(OK, NewRes(data))
}

// POSTUser saves a new user in the database
func POSTUser(c *gin.Context) {
	var data form.UserForm

	if err := c.BindJSON(&data); err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind).SetMeta(ErrBadForm)
		return
	}

	var user model.User
	user.Name = data.Name
	user.Email = data.Email
	user.IsCreator = data.IsCreator
	user.Secret = data.Secret

	uID, err := model.NewUser(&user)
	if err != nil {
		c.Error(err).SetMeta(ErrDB)
		return
	}

	// Saves the customer in third party (Stripe for now)
	customer, errCust := model.NewCustomer(data.Email, data.Plan, data.CardToken)
	if errCust != nil {
		model.DeleteUser(uID)
		c.Error(errCust).SetMeta(ErrDB)
		return
	}

	user.CustomerID = customer.ID

	// Sets customer id to User
	model.UpdateCustomerID(uID, customer.ID)

	// Logs customer subscription in the DB
	if _, err := model.NewPlanUser(uID, strings.Split(data.Plan, "_")[0], customer.Subs.Values[0].PeriodEnd); err != nil {
		// TODO when fail stripe should't charge
		c.Error(err).SetMeta(ErrIntServ)
		return
	}

	c.Header("Location", "/token/generate")

	c.JSON(Created, NewRes(user))
}

// DELETEUser delete the authenticated user
func DELETEUser(c *gin.Context) {
	user, _ := c.Get("user")
	uID := user.(*model.User).ID

	pkgs, _ := model.AllPackages(nil, uID)

	storage := lib.NewStorage(lib.SrcPackages)

	var pkgToUpdt lib.StringSlice
	for _, pkg := range *pkgs {
		if pkg.Source != nil {
			storage.DeleteFile(fmt.Sprintf("%d", uID), fmt.Sprintf("%d", pkg.ID.ValueDecoded), "wooble.js")
			pkgToUpdt = append(pkgToUpdt, fmt.Sprintf("%d", pkg.ID.ValueDecoded))
		}
	}

	if storage.Error() != nil {
		c.Error(storage.Error())
	}

	model.BulkUpdatePackageSource(pkgToUpdt, "")

	// TODO Unsub plan

	if err := model.SafeDeleteUser(uID); err != nil {
		c.Error(err).SetMeta(ErrDB)
		return
	}

	c.JSON(NoContent, nil)
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
		c.Error(nil).SetMeta(ErrBadCreds)
		return
	}

	if err := model.UpdateUserPassword(user.(*model.User).ID, passwordForm.NewSecret); err != nil {
		c.Error(err).SetMeta(ErrDB)
		return
	}

	c.Header("Location", fmt.Sprintf("/users/%s", user.(*model.User).Name))

	c.JSON(NoContent, nil)
}
