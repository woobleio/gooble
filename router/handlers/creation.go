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

	creaID := c.Param("encid")

	if creaID != "" {
		data, err = model.CreationByID(creaID)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Error(err).SetMeta(ErrResNotFound.SetParams("source", "Creation", "id", creaID))
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

	user, _ := c.Get("user")

	data.CreatorID = user.(*model.User).ID

	creaID, err := model.NewCreation(&data)
	if err != nil {
		c.Error(err).SetMeta(ErrDBSave)
		return
	}

	eng, err := model.EngineByName(data.Engine)
	if err != nil {
		c.Error(err).SetMeta(ErrServ.SetParams("source", "script"))
		return
	}

	storage := lib.NewStorage(lib.SrcCreations)

	userIDStr := fmt.Sprintf("%d", user.(*model.User).ID)

	if data.Document != "" {
		storage.StoreFile(data.Document, "text/html", userIDStr, creaID, data.Version, "doc.html")
	}
	if data.Script != "" {
		storage.StoreFile(data.Script, eng.ContentType, userIDStr, creaID, data.Version, "script"+eng.Extension)
	}
	if data.Style != "" {
		storage.StoreFile(data.Style, "text/css", userIDStr, creaID, data.Version, "style.css")
	}

	if storage.Error != nil {
		// Delete the crea since files failed to be save in the cloud
		model.DeleteCreation(creaID)
		c.Error(storage.Error).SetMeta(ErrServ.SetParams("source", "files"))
		return
	}

	c.Header("Location", fmt.Sprintf("/%s/%s", "creations", creaID))

	c.JSON(Created, nil)
}

// BuyCreations is a handler that purchases creations
func BuyCreations(c *gin.Context) {
	var buyForm struct {
		Creations []string `json:"creations,omitempty" binding:"required"`
		CardToken string   `json:"cardToken"`
	}

	if err := c.BindJSON(&buyForm); err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind).SetMeta(ErrBadForm)
		return
	}

	user, _ := c.Get("user")
	userID := user.(*model.User).ID
	customerID, err := model.UserCustomerID(userID)
	if err != nil || customerID == "" {
		c.Error(err).SetMeta(ErrDBSelect)
		return
	}

	totalAmount := uint64(0)
	creas := make([]model.Creation, 0)
	for _, creaID := range buyForm.Creations {
		crea, err := model.CreationByID(creaID)
		if err != nil {
			c.Error(err).SetMeta(ErrResNotFound.SetParams("source", "Creation", "id", creaID))
			return
		}
		if crea.CreatorID == userID {
			c.Error(nil).SetMeta(ErrCantBuy.SetParams("id", crea.ID.ValueEncoded))
			return
		}
		totalAmount = totalAmount + crea.Price
		creas = append(creas, *crea)
	}

	var chargeID string
	if buyForm.CardToken == "" {
		charge, chargeErr := lib.ChargeCustomerForCreations(customerID, totalAmount, buyForm.Creations)
		if chargeErr != nil {
			c.Error(chargeErr).SetMeta(ErrCharge)
			return
		}
		chargeID = charge.ID
	} else {
		charge, chargeErr := lib.ChargeOneTimeForCreations(totalAmount, buyForm.Creations, buyForm.CardToken)
		if chargeErr != nil {
			c.Error(chargeErr).SetMeta(ErrCharge)
			return
		}
		chargeID = charge.ID
	}

	if err := model.NewCreationPurchases(userID, chargeID, &creas); err != nil {
		c.Error(err).SetMeta(ErrDBSave)
		return
	}

	lib.CaptureCharge(chargeID)

	// TODO location to mycreations
	c.Header("Location", fmt.Sprintf("/%s/%s", "creations", buyForm.Creations[0]))

	c.JSON(OK, nil)
}

// EditCreation return private creation view
func EditCreation(c *gin.Context) {
	var data model.CreationForm

	creaID := c.Param("encid")

	user, _ := c.Get("user")

	userID := user.(*model.User).ID

	crea, err := model.CreationEditByID(creaID, userID)
	if err != nil {
		c.Error(err).SetMeta(ErrResNotFound.SetParams("source", "creation", "id", creaID))
		return
	}

	storage := lib.NewStorage(lib.SrcCreations)

	latestVersion := crea.Versions[len(crea.Versions)-1]

	if crea.HasDoc {
		data.Document = storage.GetFileContent(fmt.Sprintf("%d", crea.CreatorID), crea.ID.ValueEncoded, latestVersion, "doc.html")
	}
	if crea.HasScript {
		data.Script = storage.GetFileContent(fmt.Sprintf("%d", crea.CreatorID), crea.ID.ValueEncoded, latestVersion, "script.js")
	}
	if crea.HasStyle {
		data.Style = storage.GetFileContent(fmt.Sprintf("%d", crea.CreatorID), crea.ID.ValueEncoded, latestVersion, "style.css")
	}

	if storage.Error != nil {
		c.Error(storage.Error).SetMeta(ErrServ.SetParams("source", "files"))
		return
	}

	data.CreatorID = crea.CreatorID
	data.Description = crea.Description.String
	data.Engine = crea.Engine.Name
	data.Version = latestVersion
	data.Price = crea.Price

	c.JSON(OK, data)
}
