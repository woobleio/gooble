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

	res := NewRes()

	opts := lib.ParseOptions(c)

	creaID := c.Param("id")

	if creaID != "" {
		data, err = model.CreationByID(creaID)
		if err != nil {
			if err == sql.ErrNoRows {
				res.Error(ErrResNotFound, "Creation", creaID)
			} else {
				res.Error(ErrDBSelect)
			}
		}
	} else {
		data, err = model.AllCreations(opts)
		if err != nil {
			res.Error(ErrDBSelect)
		}
	}

	res.Response(data)

	c.JSON(res.HTTPStatus(), res)
}

// POSTCreations is a handler that retrieve a form and create a creation
func POSTCreations(c *gin.Context) {
	var data model.CreationForm

	res := NewRes()

	// FIXME workaround gin issue with Bind (https://github.com/gin-gonic/gin/issues/633)
	c.Header("Content-Type", gin.MIMEJSON)
	if c.BindJSON(&data) != nil {
		res.Error(ErrBadForm, "title (string), engine (string) are required")
		c.JSON(res.HTTPStatus(), res)
		return
	}

	if data.Version == "" {
		data.Version = model.BaseVersion
	}

	user, _ := c.Get("user")

	data.CreatorID = user.(*model.User).ID

	creaID, err := model.NewCreation(&data)
	if err != nil {
		fmt.Print(err)
		res.Error(ErrDBSave, "")
		c.JSON(res.HTTPStatus(), res)
		return
	}

	eng, err := model.EngineByName(data.Engine)
	if err != nil {
		res.Error(ErrServ, "engine : "+data.Engine+" does not exist")
		c.JSON(res.HTTPStatus(), res)
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
		res.Error(ErrServ, "doc, script and style files")
	}

	c.Header("Location", fmt.Sprintf("/%s/%s", "creations", creaID))

	res.Status = Created

	c.JSON(res.HTTPStatus(), res)
}

// BuyCreation is a handler that purchases a creation
func BuyCreation(c *gin.Context) {
	var cardForm struct {
		CardToken string `json:"cardToken"`
	}

	res := NewRes()

	// FIXME workaround gin issue with Bind (https://github.com/gin-gonic/gin/issues/633)
	c.Header("Content-Type", gin.MIMEJSON)
	if c.BindJSON(&cardForm) != nil {
		res.Error(ErrBadForm, "")
		c.JSON(res.HTTPStatus(), res)
		return
	}

	creaID := c.Param("id")

	crea, err := model.CreationByID(creaID)
	if err != nil {
		if err == sql.ErrNoRows {
			res.Error(ErrResNotFound, "Creation", creaID)
		} else {
			res.Error(ErrDBSelect)
		}
		c.JSON(res.HTTPStatus(), res)
		return
	}

	user, _ := c.Get("user")
	userID := user.(*model.User).ID

	if crea.CreatorID == userID {
		res.Error(ErrCantBuy, "User is the owner of the creation")
		c.JSON(res.HTTPStatus(), res)
		return
	}

	customerID, err := model.UserCustomerID(userID)
	if err != nil || customerID == "" {
		res.Error(ErrDBSelect)
		c.JSON(res.HTTPStatus(), res)
		return
	}

	var chargeID string
	if cardForm.CardToken == "" {
		charge, chargeErr := lib.ChargeCustomerForCreations(customerID, crea.Price, []string{crea.ID.ValueEncoded})
		if chargeErr != nil {
			res.Error(ErrCharge, "creations", "customer to charge not found")
			c.JSON(res.HTTPStatus(), res)
			return
		}
		chargeID = charge.ID
	} else {
		charge, chargeErr := lib.ChargeOneTimeForCreations(crea.Price, []string{crea.ID.ValueEncoded}, cardForm.CardToken)
		if chargeErr != nil {
			res.Error(ErrCharge, "creations", "wrong billing info")
			c.JSON(res.HTTPStatus(), res)
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
		res.Error(ErrDBSave, "- Creation already purchased")
		c.JSON(res.HTTPStatus(), res)
		return
	}

	if err := model.UpdateUserTotalDue(userID, crea.Price); err != nil {
		res.Error(ErrDBSave, "- Failed to credit the creator")
	}

	c.Header("Location", fmt.Sprintf("/%s/%s", "creations", creaID))

	res.Status = OK

	c.JSON(res.HTTPStatus(), res)
}
