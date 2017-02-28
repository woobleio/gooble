package handler

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	version "github.com/mcuadros/go-version"

	"wooble/forms"
	"wooble/lib"
	"wooble/models"
	enum "wooble/models/enums"

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
				c.Error(err).SetMeta(ErrDB)
			}
			return
		}
	} else {
		data, err = model.AllCreations(opts)
		if err != nil {
			c.Error(err).SetMeta(ErrDB)
			return
		}
	}

	c.JSON(OK, NewRes(data))
}

// POSTCreation creates a new creation
func POSTCreation(c *gin.Context) {
	var data form.CreationForm
	data.State = enum.Draft

	if err := c.BindJSON(&data); err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind).SetMeta(ErrBadForm)
		return
	}

	user, _ := c.Get("user")

	crea := new(model.Creation)
	crea.CreatorID = user.(*model.User).ID
	crea.State = data.State
	crea.Title = data.Title
	crea.Description = lib.InitNullString(data.Description)
	crea.Price = data.Price
	crea.Engine = model.Engine{Name: data.Engine}

	var errCrea error
	crea, errCrea = model.NewCreation(crea)
	if errCrea != nil {
		c.Error(errCrea).SetMeta(ErrDB)
		return
	}

	c.Header("Location", fmt.Sprintf("/%s/%s", "creations", crea.ID.ValueEncoded))

	c.JSON(Created, NewRes(crea))
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
		c.Error(err).SetMeta(ErrDB)
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
		charge, chargeErr := model.ChargeCustomerForCreations(customerID, totalAmount, buyForm.Creations)
		if chargeErr != nil {
			c.Error(chargeErr).SetMeta(ErrCharge)
			return
		}
		chargeID = charge.ID
	} else {
		charge, chargeErr := model.ChargeOneTimeForCreations(totalAmount, buyForm.Creations, buyForm.CardToken)
		if chargeErr != nil {
			c.Error(chargeErr).SetMeta(ErrCharge)
			return
		}
		chargeID = charge.ID
	}

	if err := model.NewCreationPurchases(userID, chargeID, &creas); err != nil {
		c.Error(err).SetMeta(ErrDB)
		return
	}

	model.CaptureCharge(chargeID)

	// TODO location to mycreations
	// c.Header("Location", fmt.Sprintf("/creations/%s/code", buyForm.Creations[0]))

	c.JSON(NoContent, nil)
}

// GETCodeCreation return private creation view
func GETCodeCreation(c *gin.Context) {
	var data form.CreationCodeForm

	creaID := c.Param("encid")

	user, _ := c.Get("user")

	var errCrea error
	crea := new(model.Creation)
	crea.ID = lib.InitID(creaID)
	crea.CreatorID = user.(*model.User).ID

	crea, errCrea = model.CreationEditByID(crea)
	if errCrea != nil {
		c.Error(errCrea).SetMeta(ErrResNotFound.SetParams("source", "creation", "id", creaID))
		return
	}

	storage := lib.NewStorage(lib.SrcCreations)

	latestVersion := crea.Versions[len(crea.Versions)-1]
	data.Script = storage.GetFileContent(fmt.Sprintf("%d", crea.CreatorID), crea.ID.ValueEncoded, latestVersion, "script.js")

	if crea.HasDoc {
		data.Document = storage.GetFileContent(fmt.Sprintf("%d", crea.CreatorID), crea.ID.ValueEncoded, latestVersion, "doc.html")
	}
	if crea.HasStyle {
		data.Style = storage.GetFileContent(fmt.Sprintf("%d", crea.CreatorID), crea.ID.ValueEncoded, latestVersion, "style.css")
	}

	if storage.Error() != nil {
		c.Error(storage.Error()).SetMeta(ErrServ.SetParams("source", "files"))
		return
	}

	c.JSON(OK, NewRes(data))
}

// PUTCreation edits creation information
func PUTCreation(c *gin.Context) {
	var creaForm form.CreationForm

	if err := c.BindJSON(&creaForm); err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind).SetMeta(ErrBadForm)
		return
	}

	creaID := c.Param("encid")

	user, _ := c.Get("user")

	var crea model.Creation
	crea.ID = lib.InitID(creaID)
	crea.CreatorID = user.(*model.User).ID
	crea.Title = creaForm.Title
	crea.Description = lib.InitNullString(creaForm.Description)
	crea.Price = creaForm.Price

	if err := model.UpdateCreation(&crea); err != nil {
		c.Error(err).SetMeta(ErrDB)
		return
	}

	c.Header("Location", fmt.Sprintf("/%s/%s", "creations", creaID))

	c.JSON(NoContent, nil)
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

	userIDStr := fmt.Sprintf("%d", user.(*model.User).ID)

	version = strings.Replace(version, "_", ".", -1)

	var crea model.Creation
	crea.ID = lib.InitID(creaID)
	crea.HasDoc = codeForm.Document != ""
	crea.HasStyle = codeForm.Style != ""
	crea.Version = version

	if nbRowAffected, err := model.UpdateCreationCode(&crea); err != nil || nbRowAffected == 0 {
		c.Error(err).SetMeta(ErrDB)
		return
	}

	storage := lib.NewStorage(lib.SrcCreations)
	if codeForm.Document != "" {
		storage.StoreFile(codeForm.Document, "text/html", userIDStr, creaID, version, "doc.html")
	}
	if codeForm.Script != "" {
		storage.StoreFile(codeForm.Script, "application/javascript", userIDStr, creaID, version, "script.js")
	}
	if codeForm.Style != "" {
		storage.StoreFile(codeForm.Style, "text/css", userIDStr, creaID, version, "style.css")
	}

	if storage.Error() != nil {
		c.Error(storage.Error()).SetMeta(ErrServ.SetParams("source", "files"))
		return
	}

	c.Header("Location", fmt.Sprintf("/%s/%s", "creations", creaID))

	c.JSON(NoContent, nil)
}

// PublishCreation make a creation public
func PublishCreation(c *gin.Context) {
	user, _ := c.Get("user")

	var crea model.Creation
	crea.ID = lib.InitID(c.Param("encid"))
	crea.CreatorID = user.(*model.User).ID

	if err := model.PublishCreation(&crea); err != nil {
		c.Error(err).SetMeta(ErrDB)
		return
	}

	c.JSON(NoContent, nil)
}

// POSTCreationVersion creates a new version
func POSTCreationVersion(c *gin.Context) {
	var versionForm struct {
		Version string `json:"version" validate:"required"`
	}

	if err := c.BindJSON(&versionForm); err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind).SetMeta(ErrBadForm)
		return
	}

	user, _ := c.Get("user")

	var errCrea error
	crea := new(model.Creation)
	crea.ID = lib.InitID(c.Param("encid"))
	crea.CreatorID = user.(*model.User).ID
	crea, errCrea = model.CreationEditByID(crea)
	if errCrea != nil {
		c.Error(errCrea).SetMeta(ErrResNotFound.SetParams("source", "creation", "id", crea.ID.ValueEncoded))
		return
	}

	curVersion := crea.Versions[len(crea.Versions)-1]
	if version.Compare(curVersion, versionForm.Version, ">=") {
		c.Error(errors.New("Version POST issue, either malformed or posted version lesser than the latest")).SetMeta(ErrCreaVersion.SetParams("version", versionForm.Version))
		return
	}

	crea.Versions = append(crea.Versions, versionForm.Version)

	if nbRowsAff, err := model.NewCreationVersion(crea); err != nil || nbRowsAff == 0 {
		c.Error(err).SetMeta(ErrDB)
		return
	}

	storage := lib.NewStorage(lib.SrcCreations)
	storage.CopyAndStoreFile(fmt.Sprintf("%d", crea.CreatorID), crea.ID.ValueEncoded, curVersion, versionForm.Version, "script.js")

	if crea.HasDoc {
		storage.CopyAndStoreFile(fmt.Sprintf("%d", crea.CreatorID), crea.ID.ValueEncoded, curVersion, versionForm.Version, "doc.html")
	}
	if crea.HasStyle {
		storage.CopyAndStoreFile(fmt.Sprintf("%d", crea.CreatorID), crea.ID.ValueEncoded, curVersion, versionForm.Version, "style.css")
	}

	if storage.Error() != nil {
		c.Error(storage.Error()).SetMeta(ErrServ.SetParams("source", "copy"))
		return
	}

	c.Header("Location", fmt.Sprintf("/creations/%s/code", c.Param("encid")))

	c.JSON(OK, nil)
}
