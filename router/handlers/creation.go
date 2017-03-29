package handler

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/js"

	"wooble/forms"
	"wooble/lib"
	"wooble/models"
	"wooble/models/enums"
	helper "wooble/router/helpers"
)

// GETCreations is a handler that returns one or more creations
func GETCreations(c *gin.Context) {
	var data interface{}
	var err error

	opts := lib.ParseOptions(c)

	creaID := c.Param("encid")

	token, _ := helper.ParseToken(c)

	// Auth not mandatory, only here to know either the auth'd user owns the creation or not
	authUserID := uint64(0)
	if token != nil {
		if user, _ := model.UserByToken(token); user != nil {
			authUserID = user.ID
		}
	}

	if creaID != "" {
		data, err = model.CreationByID(lib.InitID(creaID), authUserID)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Error(err).SetMeta(ErrResNotFound.SetParams("source", "Creation", "id", creaID))
			} else {
				c.Error(err).SetMeta(ErrDB)
			}
			return
		}

		creaLastVersion := fmt.Sprintf("%d", data.(*model.Creation).Versions[len(data.(*model.Creation).Versions)-1])
		s := lib.NewStorage(lib.SrcPreview)
		creatorID := fmt.Sprintf("%d", data.(*model.Creation).Creator.ID)
		creaID := fmt.Sprintf("%d", data.(*model.Creation).ID.ValueDecoded)
		previewURL := s.GetPathFor(creatorID, creaID, creaLastVersion, "index.html")
		spltPath := strings.Split(previewURL, "/")
		data.(*model.Creation).PreviewURL = lib.GetAssetURL() + "/" + strings.Join(spltPath[1:], "/")

	} else {
		data, err = model.AllCreations(opts, authUserID)
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

	if err := c.BindJSON(&data); err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind).SetMeta(ErrBadForm)
		return
	}

	user, _ := c.Get("user")

	crea := new(model.Creation)
	crea.CreatorID = user.(*model.User).ID
	crea.State = enum.Draft
	crea.Title = data.Title
	crea.Description = lib.InitNullString(data.Description)
	crea.Price = data.Price
	crea.Alias = data.Alias
	if data.Engine == "" {
		data.Engine = "JSES5"
	}
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

// DELETECreation delete the creation of the authenticated user
func DELETECreation(c *gin.Context) {
	creaID := lib.InitID(c.Param("encid"))

	user, _ := c.Get("user")
	uID := user.(*model.User).ID

	if model.CreationInUse(creaID) {
		if err := model.SafeDeleteCreation(uID, creaID); err != nil {
			c.Error(err).SetMeta(ErrDB)
			return
		}
	} else {
		crea, err := model.CreationPrivateByID(uID, creaID)
		if err != nil {
			c.Error(err).SetType(gin.ErrorTypeBind).SetMeta(ErrResNotFound.SetParams("source", "creation", "id", creaID.ValueEncoded))
			return
		}
		if err := model.DeleteCreation(uID, creaID); err != nil {
			c.Error(err).SetMeta(ErrDB)
			return
		}

		uIDStr := fmt.Sprintf("%d", uID)
		creaIDStr := fmt.Sprintf("%d", crea.ID.ValueDecoded)

		storage := lib.NewStorage(lib.SrcCreations)

		for _, v := range crea.Versions {
			storage.PushBulkFile(uIDStr, creaIDStr, fmt.Sprintf("%d", v), enum.Script)
			storage.PushBulkFile(uIDStr, creaIDStr, fmt.Sprintf("%d", v), enum.Document)
			storage.PushBulkFile(uIDStr, creaIDStr, fmt.Sprintf("%d", v), enum.Style)
		}

		storage.BulkDeleteFiles()

		if storage.Error() != nil {
			c.Error(storage.Error()) // log error
		}
	}

	c.AbortWithStatus(NoContent)
}

// BuyCreations is a handler that purchases creations
func BuyCreations(c *gin.Context) {
	var buyForm struct {
		Creations []string `json:"creations,omitempty" validate:"required"`
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
		crea, err := model.CreationByID(lib.InitID(creaID), 0)
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

	c.AbortWithStatus(NoContent)
}

// GETCreationCode return private creation view
func GETCreationCode(c *gin.Context) {
	user, _ := c.Get("user")

	creaID := lib.InitID(c.Param("encid"))
	crea, err := model.CreationPrivateByID(user.(*model.User).ID, creaID)
	if err != nil {
		c.Error(err).SetMeta(ErrResNotFound.SetParams("source", "creation", "id", creaID.ValueEncoded))
		return
	}

	storage := lib.NewStorage(lib.SrcCreations)

	latestVersion := fmt.Sprintf("%d", crea.Versions[len(crea.Versions)-1])
	uIDStr := fmt.Sprintf("%d", crea.CreatorID)
	creaIDStr := fmt.Sprintf("%d", crea.ID.ValueDecoded)

	filter := c.Query("filter")
	if filter != "" {
		fFields := strings.Split(filter, ",")
		for _, field := range fFields {
			switch field {
			case "script":
				crea.Script = storage.GetFileContent(uIDStr, creaIDStr, latestVersion, enum.Script)
			case "document":
				crea.Document = storage.GetFileContent(uIDStr, creaIDStr, latestVersion, enum.Document)
			case "style":
				crea.Style = storage.GetFileContent(uIDStr, creaIDStr, latestVersion, enum.Style)
			}
		}
	} else {
		crea.Script = storage.GetFileContent(uIDStr, creaIDStr, latestVersion, enum.Script)
		crea.Document = storage.GetFileContent(uIDStr, creaIDStr, latestVersion, enum.Document)
		crea.Style = storage.GetFileContent(uIDStr, creaIDStr, latestVersion, enum.Style)
	}

	if storage.Error() != nil { // TODO manage storage errors better
		c.Error(storage.Error())
		crea.Script = ""
	}

	if crea.Script == "" {
		// TODO put this in wooblizer lib
		crea.Script = `var woobly = {
  "attribute": "a value (optionnal)",
  "_init": function() {
    // Creation code at runtime
  },
  "method": function(a, b) {
    // a method (optionnal)
  }
}`
	}

	c.JSON(OK, NewRes(crea))
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
	crea.State = creaForm.State
	crea.Alias = creaForm.Alias

	if err := model.UpdateCreation(&crea); err != nil {
		c.Error(err).SetMeta(ErrDB)
		return
	}

	c.Header("Location", fmt.Sprintf("/creations/%s", creaID))

	c.AbortWithStatus(NoContent)
}

// SaveVersion save the current code for a version (must be in draft state)
func SaveVersion(c *gin.Context) {
	var codeForm form.CreationCodeForm

	creaID := lib.InitID(c.Param("encid"))

	if err := c.BindJSON(&codeForm); err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind).SetMeta(ErrBadForm)
		return
	}

	user, _ := c.Get("user")
	crea, err := model.CreationPrivateByID(user.(*model.User).ID, creaID)
	if err != nil {
		c.Error(err).SetMeta(ErrResNotFound.SetParams("source", "creation", "id", creaID))
		return
	}

	if crea.State != enum.Draft {
		// TODO state error
		return
	}

	userIDStr := fmt.Sprintf("%d", user.(*model.User).ID)
	version := fmt.Sprintf("%d", crea.Versions[len(crea.Versions)-1])
	creaIDStr := fmt.Sprintf("%d", crea.ID.ValueDecoded)
	storage := lib.NewStorage(lib.SrcCreations)

	storage.StoreFile(codeForm.Document, "text/html", userIDStr, creaIDStr, version, enum.Document)
	storage.StoreFile(codeForm.Script, "application/javascript", userIDStr, creaIDStr, version, enum.Script)
	storage.StoreFile(codeForm.Style, "text/css", userIDStr, creaIDStr, version, enum.Style)

	minifier := minify.New()
	minifier.AddFunc("text/javascript", js.Minify)

	var minErr error
	codeForm.Script, minErr = minifier.String("text/javascript", codeForm.Script)

	if minErr != nil {
		c.Error(minErr).SetMeta(ErrServ.SetParams("source", "minifier"))
		return
	}

	preview := `<html>
		<head>
			<script type="text/javascript">` + codeForm.Script + `</script>
			<script type="text/javascript">(function(){woobly._init()})()</script>
			<style>` + codeForm.Style + `</style>
		</head>
		<body>
			` + codeForm.Document + `
		</body>
	</html>`

	storage.SetSource(lib.SrcPreview)
	storage.StoreFile(preview, "text/html", userIDStr, creaIDStr, version, enum.Preview)

	if storage.Error() != nil {
		c.Error(storage.Error()).SetMeta(ErrServ.SetParams("source", "files"))
		return
	}

	c.Header("Location", fmt.Sprintf("/creations/%s", creaID.ValueEncoded))

	c.AbortWithStatus(NoContent)
}

// POSTCreationVersion creates a new version
func POSTCreationVersion(c *gin.Context) {
	user, _ := c.Get("user")
	uID := user.(*model.User).ID

	creaID := lib.InitID(c.Param("encid"))
	crea, err := model.CreationPrivateByID(uID, creaID)
	if err != nil {
		c.Error(err).SetMeta(ErrResNotFound.SetParams("source", "creation", "id", creaID.ValueEncoded))
		return
	}

	curVersion := crea.Versions[len(crea.Versions)-1]
	newVersion := curVersion + 1

	uIDStr := fmt.Sprintf("%d", crea.CreatorID)
	creaIDStr := fmt.Sprintf("%d", crea.ID.ValueDecoded)
	storage := lib.NewStorage(lib.SrcCreations)
	curVersionStr := fmt.Sprintf("%d", curVersion)
	newVersionStr := fmt.Sprintf("%d", newVersion)

	storage.CopyAndStoreFile(uIDStr, creaIDStr, curVersionStr, newVersionStr, enum.Script)
	storage.CopyAndStoreFile(uIDStr, creaIDStr, curVersionStr, newVersionStr, enum.Document)
	storage.CopyAndStoreFile(uIDStr, creaIDStr, curVersionStr, newVersionStr, enum.Style)

	if storage.Error() != nil {
		c.Error(storage.Error()).SetMeta(ErrServ.SetParams("source", "copy"))
		return
	}

	if err := model.NewCreationVersion(uID, crea.ID, append(crea.Versions, newVersion)); err != nil {
		c.Error(err).SetMeta(ErrDB)
		return
	}

	c.Header("Location", fmt.Sprintf("/creations/%s/code", c.Param("encid")))

	c.AbortWithStatus(NoContent)
}
