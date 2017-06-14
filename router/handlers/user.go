package handler

import (
	"errors"
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

	users, _ := model.AllActiveUsers()

	var user model.User
	user.Name = data.Name
	user.Fullname = lib.InitNullString(data.Fullname)
	user.Email = data.Email
	user.IsCreator = data.IsCreator
	user.Secret = data.Secret
	user.IsVIP = len(users) < 310

	uID, err := model.NewUser(&user)
	if err != nil {
		c.Error(err).SetMeta(ErrDB)
		return
	}

	// Saves the customer in third party (Stripe for now)
	customer, errCust := model.NewCustomer(data.Email, data.Plan.Label, data.CardToken)
	if errCust != nil {
		model.DeleteUser(uID)
		c.Error(errCust).SetMeta(ErrDB)
		return
	}

	user.CustomerID = customer.ID

	// Sets customer id to User
	model.UpdateCustomerID(uID, customer.ID)

	// Logs customer subscription in the DB
	model.NewPlanUser(uID, strings.Split(data.Plan.Label, "_")[0], 0)

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

	if userPatchForm.NewSecret != nil {
		if !user.(*model.User).IsPasswordValid(*userPatchForm.OldSecret) {
			c.Error(nil).SetMeta(ErrBadCreds)
			return
		}
		userPatchForm.Salt = new(string)
		*userPatchForm.Salt = lib.GenKey()
		*userPatchForm.NewSecret, _ = lib.Encrypt(*userPatchForm.NewSecret, []byte(*userPatchForm.Salt))
	}

	if userPatchForm.Plan != nil {
		privUser, err := model.UserPrivateByID(user.(*model.User).ID)
		if err != nil {
			c.Error(err).SetMeta(ErrServ)
			return
		}

		// Stipe uses "label_period" id style, so let's split to get the label
		planLabel := strings.Split(userPatchForm.Plan.Label, "_")[0]
		plan, err := model.PlanByLabel(planLabel)
		if err != nil {
			c.Error(err).SetMeta(ErrResNotFound.SetParams("source", "plan", "label", planLabel))
			// return error, plan not known
			return
		}

		if planLabel != model.Free {
			// Means that the user has currently a better plan
			if privUser.Plan.Level.Int64 > plan.Level.Int64 {
				c.Error(errors.New("Already has a better plan")).SetMeta(ErrBetterPlan.SetParams("currentPlan", privUser.Plan.Label.String, "reqPlan", plan.Label))
				return
			}
			// Unsub the precedent plan if it wasn't a free one and not already unsubbed
			if privUser.Plan.UnsubDate == nil && privUser.Plan.Level.Int64 > 0 {
				if err = model.UnsubCustomer(privUser.CustomerID); err != nil {
					c.Error(err).SetMeta(ErrIntServ)
					return
				}
			}
			sub, err := model.SubCustomer(privUser.CustomerID, userPatchForm.Plan.Label, *userPatchForm.CardToken)
			if err != nil {
				c.Error(err).SetMeta(ErrIntServ)
				return
			}
			model.NewPlanUser(privUser.ID, planLabel, sub.PeriodEnd)

			privUser.Plan = plan

			// Refresh token
			token := model.NewToken(privUser, "")
			tokenS, _ := token.SignedString(model.TokenKey())
			c.Header("Authorization", tokenS)

		} else if planLabel == model.Free && privUser.Plan.UnsubDate == nil {
			if err := model.UnsubCustomer(privUser.CustomerID); err != nil {
				c.Error(err).SetMeta(ErrIntServ)
				return
			}
			model.UnsubUserPlan(uint64(privUser.PlanID.Int64))
		}
	}

	if err := model.UpdateUserPatch(user.(*model.User).ID, lib.SQLPatches(userPatchForm)); err != nil {
		c.Error(err).SetMeta(ErrDB)
		return
	}

	if userPatchForm.Operation != nil {
		switch *userPatchForm.Operation {
		case "resendActivationEmail":
			fullUser, _ := model.UserPrivateByID(user.(*model.User).ID)
			if fullUser.IsActive {
				c.Error(errors.New("User already active"))
				return
			}
			name := fullUser.Name
			if fullUser.Fullname != nil {
				name = strings.Split(fullUser.Fullname.String, " ")[0]
			}
			if err := model.SendActivationEmail(name, fullUser.Email); err != nil {
				c.Error(err)
				return
			}
		case "activateEmail":
			fullUser, _ := model.UserPrivateByID(user.(*model.User).ID)
			if fullUser.IsActive {
				c.Error(errors.New("User already active")).SetMeta(ErrActivation)
				return
			}
			if err := model.ActivateUser(fullUser.Email, *userPatchForm.ValidationToken); err != nil {
				c.Error(err).SetMeta(ErrActivation)
				return
			}

			fullUser.IsActive = true

			// Refresh token
			token := model.NewToken(fullUser, "")
			tokenS, _ := token.SignedString(model.TokenKey())
			c.Header("Authorization", tokenS)
		}
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
