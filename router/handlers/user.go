package handler

import (
	"fmt"
	"strings"
	form "wooble/forms"
	"wooble/lib"
	model "wooble/models"
	enum "wooble/models/enums"
	helper "wooble/router/helpers"

	"github.com/gin-gonic/gin"
)

// GETUser returns one users with private infos if authenticated, with public infos if not
func GETUser(c *gin.Context) {
	var data *model.User
	var err error

	token, _ := helper.ParseToken(c)

	tokenUser, _ := model.UserByToken(token)

	username := c.Param("name")

	if tokenUser != nil && username == tokenUser.Name {
		data, err = model.UserPrivateByID(tokenUser.ID)
	} else {
		data, err = model.UserPublicByName(c.Param("name"))
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

	c.Header("Location", "/tokens")

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
			storage.PushBulkFile(fmt.Sprintf("%d", uID), fmt.Sprintf("%d", pkg.ID.ValueDecoded), "", enum.Wooble)
			pkgToUpdt = append(pkgToUpdt, fmt.Sprintf("%d", pkg.ID.ValueDecoded))
		}
	}
	storage.BulkDeleteFiles()

	if storage.Error() != nil {
		c.Error(storage.Error())
	}

	model.BulkUpdatePackageSource(pkgToUpdt, "")

	privUser, _ := model.UserPrivateByID(uID)
	if err := model.UnsubCustomer(privUser.CustomerID); err != nil {
		c.Error(err).SetMeta(ErrIntServ)
		return
	}

	if err := model.SafeDeleteUser(uID); err != nil {
		c.Error(err).SetMeta(ErrDB)
		return
	}

	c.AbortWithStatus(NoContent)
}

// PATCHUser update authenticated user's password
func PATCHUser(c *gin.Context) {
	var userPatchForm form.UserPatchForm

	if err := c.BindJSON(&userPatchForm); err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind).SetMeta(ErrBadForm)
		return
	}

	user, _ := c.Get("user")

	// Check the password if is in the patch
	if userPatchForm.NewSecret != nil {
		if !user.(*model.User).IsPasswordValid(*userPatchForm.OldSecret) {
			c.Error(nil).SetMeta(ErrBadCreds)
			return
		}
		userPatchForm.Salt = new(string)
		*userPatchForm.Salt = lib.GenKey()
		*userPatchForm.NewSecret, _ = lib.Encrypt(*userPatchForm.NewSecret, []byte(*userPatchForm.Salt))
	}

	if err := model.UpdateUserPatch(user.(*model.User).ID, lib.SQLPatches(userPatchForm)); err != nil {
		c.Error(err).SetMeta(ErrDB)
		return
	}

	c.Header("Location", fmt.Sprintf("/users/%s", user.(*model.User).Name))

	c.AbortWithStatus(NoContent)
}

// POSTUserBank creates bank info for customer (only for funds)
func POSTUserBank(c *gin.Context) {
	var bankForm struct {
		BankToken string `json:"bankToken" validate:"required"`
	}

	if err := c.BindJSON(&bankForm); err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind).SetMeta(ErrBadForm)
		return
	}

	var errUser error
	privateUser := new(model.User)
	user, _ := c.Get("user")
	privateUser, errUser = model.UserPrivateByID(user.(*model.User).ID)
	if errUser != nil {
		c.Error(errUser).SetMeta(ErrDB)
		return
	}

	acc, err := model.RegisterBank(privateUser.Email, bankForm.BankToken)
	if err != nil {
		c.Error(err).SetMeta(ErrServ.SetParams("source", "bank register"))
		return
	}

	if err := model.UpdateUserAccountID(privateUser.ID, acc.ID); err != nil {
		c.Error(err).SetMeta(ErrDB)
		return
	}

	c.AbortWithStatus(NoContent)
}

// WithdrawFunds withdraws users fund to a registered bank account
func WithdrawFunds(c *gin.Context) {
	user, _ := c.Get("user")

	var errUser error
	privateUser := new(model.User)
	privateUser, errUser = model.UserPrivateByID(user.(*model.User).ID)
	if errUser != nil {
		c.Error(errUser).SetMeta(ErrDB)
		return
	}

	if privateUser.Fund <= 0 {
		// Error nothing to withdraw TODO
		return
	}

	if _, err := model.PayUser(privateUser.AccountID.String, privateUser.Fund); err != nil {
		c.Error(err).SetMeta(ErrIntServ)
		return
	}

	if err := model.UserSubFund(privateUser.ID, privateUser.Fund); err != nil {
		c.Error(err).SetMeta(ErrDB)
		return
	}

	c.AbortWithStatus(NoContent)
}
