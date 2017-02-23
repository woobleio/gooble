package handler

import (
	"database/sql"
	"fmt"
	"strings"

	"wooble/forms"
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

// POSTCreations creates a new creation
func POSTCreations(c *gin.Context) {
	var data form.CreationForm

	if err := c.BindJSON(&data); err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind).SetMeta(ErrBadForm)
		return
	}

	user, _ := c.Get("user")

	var crea model.Creation
	crea.CreatorID = user.(*model.User).ID
	crea.Title = data.Title
	crea.Description = lib.InitNullString(data.Description)
	crea.Price = data.Price
	crea.Engine = model.Engine{Name: data.Engine}

	creaID, err := model.NewCreation(&crea)
	if err != nil {
		c.Error(err).SetMeta(ErrDBSave)
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

// GETCodeCreation return private creation view
func GETCodeCreation(c *gin.Context) {
	var data form.CreationCodeForm

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

	c.JSON(OK, data)
}

// PUTCreations edits creation information
func PUTCreations(c *gin.Context) {
	var creaForm form.CreationForm

	if err := c.Bind(&creaForm); err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind).SetMeta(ErrBadForm)
		return
	}

	creaID := c.Param("encid")

	user, _ := c.Get("user")

	var crea model.Creation
	crea.CreatorID = user.(*model.User).ID
	crea.Title = creaForm.Title
	crea.Description = lib.InitNullString(creaForm.Description)
	crea.Price = creaForm.Price

	if err := model.UpdateCreation(creaID, &crea); err != nil {
		c.Error(err).SetMeta(ErrDBSave.SetParams("source", "creation", "id", creaID))
		return
	}

	c.Header("Location", fmt.Sprintf("/%s/%s", "creations", creaID))

	c.JSON(OK, nil)
}

// SaveVersion save the current code for a version (must be in draft state)
func SaveVersion(c *gin.Context) {
	var codeForm form.CreationCodeForm

	version := c.Param("version")
	creaID := c.Param("encid")

	if err := c.BindJSON(&codeForm); err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind).SetMeta(ErrBadForm)
		return
	}

	user, _ := c.Get("user")

	storage := lib.NewStorage(lib.SrcCreations)

	userIDStr := fmt.Sprintf("%d", user.(*model.User).ID)

	version = strings.Replace(version, "_", ".", -1)

	var crea model.Creation
	crea.ID = lib.InitID(creaID)
	crea.HasDoc = codeForm.Document != ""
	crea.HasStyle = codeForm.Style != ""
	crea.HasScript = codeForm.Script != ""
	crea.Version = version

	if err := model.UpdateCreationCode(&crea); err != nil {
		c.Error(err).SetMeta(ErrDBSave)
		return
	}

	if codeForm.Document != "" {
		storage.StoreFile(codeForm.Document, "text/html", userIDStr, creaID, version, "doc.html")
	}
	if codeForm.Script != "" {
		storage.StoreFile(codeForm.Script, "application/javascript", userIDStr, creaID, version, "script.js") // TODO Engine extension instead of .js
	}
	if codeForm.Style != "" {
		storage.StoreFile(codeForm.Style, "text/css", userIDStr, creaID, version, "style.css")
	}

	if storage.Error != nil {
		c.Error(storage.Error).SetMeta(ErrServ.SetParams("source", "files"))
		return
	}

	c.Header("Location", fmt.Sprintf("/%s/%s", "creations", creaID))

	c.JSON(Created, nil)
}
