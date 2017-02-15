package handler

import (
	"database/sql"
	"fmt"

	"wooble/lib"
	"wooble/models"

	"gopkg.in/gin-gonic/gin.v1"
)

// GETCreations is a handler that returns one or more creations
func GETCreations(c *gin.Context) {
	var data interface{}
	var err error

	opts := lib.ParseOptions(c)

	creaID := c.Param("id")

	if creaID != "" {
		data, err = model.CreationByID(creaID)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Error(err).SetMeta(ErrResNotFound)
			} else {
				c.Error(err).SetMeta(ErrDBSelect)
			}
			return
		}
	} else {
		data, err = model.AllCreations(opts)
		if err != nil {
			c.Error(err).SetMeta(ErrDBSelect)
			return
		}
	}

	c.JSON(OK, NewRes(data))
}

// POSTCreations is a handler that retrieve a form and create a creation
func POSTCreations(c *gin.Context) {
	var data model.CreationForm

	if err := c.BindJSON(&data); err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind).SetMeta(ErrBadForm)
		return
	}

	if data.Version == "" {
		data.Version = model.BaseVersion
	}

	user, _ := c.Get("user")

	data.CreatorID = user.(*model.User).ID

	creaID, err := model.NewCreation(&data)
	if err != nil {
		c.Error(err).SetMeta(ErrDBSave)
		return
	}

	eng, err := model.EngineByName(data.Engine)
	if err != nil {
		c.Error(err).SetMeta(ErrServ)
		return
	}

	storage := lib.NewStorage(lib.SrcCreations, data.Version)

	userIDStr := fmt.Sprintf("%d", user.(*model.User).ID)

	if data.Document != "" {
		storage.StoreFile(data.Document, "text/html", userIDStr, creaID, "doc.html")
	}
	if data.Script != "" {
		storage.StoreFile(data.Script, eng.ContentType, userIDStr, creaID, "script"+eng.Extension)
	}
	if data.Style != "" {
		storage.StoreFile(data.Style, "text/css", userIDStr, creaID, "style.css")
	}

	if storage.Error != nil {
		// Delete the crea since files failed to be save in the cloud
		model.DeleteCreation(creaID)
		c.Error(storage.Error).SetMeta(ErrServ)
	}

	c.Header("Location", fmt.Sprintf("/%s/%s", "creations", creaID))

	c.JSON(Created, nil)
}

// BuyCreation is a handler that purchases a creation
func BuyCreation(c *gin.Context) {
	var cardForm struct {
		CardToken string `json:"cardToken"`
	}

	if err := c.BindJSON(&cardForm); err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind).SetMeta(ErrBadForm)
		return
	}

	creaID := c.Param("id")

	crea, err := model.CreationByID(creaID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.Error(err).SetMeta(ErrResNotFound)
		} else {
			c.Error(err).SetMeta(ErrDBSelect)
		}
		return
	}

	user, _ := c.Get("user")
	userID := user.(*model.User).ID

	if crea.CreatorID == userID {
		c.Error(nil).SetMeta(ErrCantBuy)
		return
	}

	customerID, err := model.UserCustomerID(userID)
	if err != nil || customerID == "" {
		c.Error(err).SetMeta(ErrDBSelect)
		return
	}

	var chargeID string
	if cardForm.CardToken == "" {
		charge, chargeErr := lib.ChargeCustomerForCreations(customerID, crea.Price, []string{crea.ID.ValueEncoded})
		if chargeErr != nil {
			c.Error(chargeErr).SetMeta(ErrCharge)
			return
		}
		chargeID = charge.ID
	} else {
		charge, chargeErr := lib.ChargeOneTimeForCreations(crea.Price, []string{crea.ID.ValueEncoded}, cardForm.CardToken)
		if chargeErr != nil {
			c.Error(chargeErr).SetMeta(ErrCharge)
			return
		}
		chargeID = charge.ID
	}

	creaPurchase := model.CreationPurchase{
		UserID:   userID,
		CreaID:   crea.ID.ValueDecoded,
		Total:    crea.Price,
		ChargeID: chargeID,
	}

	if err := model.NewCreationPurchase(&creaPurchase); err != nil {
		c.Error(err).SetMeta(ErrDBSave)
		return
	}

	if err := model.UpdateUserTotalDue(userID, crea.Price); err != nil {
		c.Error(err).SetMeta(ErrDBSave)
		return
	}

	c.Header("Location", fmt.Sprintf("/%s/%s", "creations", creaID))

	c.JSON(OK, nil)
}
