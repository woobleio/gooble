package handler

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/js"
	"github.com/woobleio/wooblizer/wbzr"
	"github.com/woobleio/wooblizer/wbzr/engine"

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

	s := lib.NewStorage(lib.SrcPreview)
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

		creaLastVersion := fmt.Sprintf("%d", data.(*model.Creation).Version)
		creatorID := fmt.Sprintf("%d", data.(*model.Creation).Creator.ID)
		creaID := fmt.Sprintf("%d", data.(*model.Creation).ID.ValueDecoded)
		previewURL := s.GetPathFor(creatorID, creaID, creaLastVersion, "index.html")
		spltPath := strings.Split(previewURL, "/")
		data.(*model.Creation).PreviewURL = strings.Join(spltPath[1:], "/")
	} else {
		switch c.DefaultQuery("list", "") {
		case "popular":
			data, err = model.AllPopularCreations(opts, authUserID)
		case "used":
			if authUserID > 0 {
				data, err = model.AllUsedCreations(opts, authUserID)
			}
		case "draft":
			if authUserID > 0 {
				data, err = model.AllDraftCreations(opts, authUserID)
			}
		default:
			data, err = model.AllCreations(opts, authUserID)
		}

		if err != nil {
			c.Error(err).SetMeta(ErrDB)
			return
		}

		for i, crea := range data.([]model.Creation) {
			nbVersions := len(crea.Versions)
			if nbVersions > 0 && crea.IsThumbPreview { // Means the creation is in draft and has never been published, it can't have a preview
				url := s.GetPathFor(fmt.Sprintf("%d", crea.Creator.ID), fmt.Sprintf("%d", crea.ID.ValueDecoded), fmt.Sprintf("%d", crea.Versions[nbVersions-1]), "index.html")
				spltPath := strings.Split(url, "/")
				data.([]model.Creation)[i].PreviewURL = strings.Join(spltPath[1:], "/")
			}
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
	crea.Alias = data.Alias
	if data.Engine == "" {
		data.Engine = "JS"
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
		crea, err := model.CreationPrivateByID(uID, creaID, "")
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

// GETCreationCode return private creation view
func GETCreationCode(c *gin.Context) {
	user, _ := c.Get("user")

	creaID := lib.InitID(c.Param("encid"))
	crea, err := model.CreationPrivateByID(user.(*model.User).ID, creaID, c.Query("v"))
	if err != nil {
		c.Error(err).SetMeta(ErrResNotFound.SetParams("source", "creation", "id", creaID.ValueEncoded))
		return
	}

	filter := c.Query("filter")
	creaVersion := fmt.Sprintf("%d", crea.Version)
	if filter != "" {
		fFields := strings.Split(filter, ",")
		for _, field := range fFields {
			switch field {
			case "script":
				crea.RetrieveSourceCode(creaVersion, enum.Script)
			case "document":
				crea.RetrieveSourceCode(creaVersion, enum.Document)
			case "style":
				crea.RetrieveSourceCode(creaVersion, enum.Style)
			}
		}
	} else {
		if storErr := crea.RetrieveSourceCode(creaVersion, enum.Script, enum.Document, enum.Style); storErr != nil {
			c.Error(storErr)
			crea.Script = ""
		}
	}

	if crea.Script == "" {
		// TODO put this in wooblizer lib
		crea.Script = wbzr.WooblyJS
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
	crea.ThumbPath = lib.InitNullString(creaForm.ThumbPath)
	crea.IsThumbPreview = creaForm.IsThumbPreview
	crea.State = creaForm.State
	crea.Alias = creaForm.Alias
	crea.Params = creaForm.Params
	crea.Version = uint64(creaForm.Version)

	if err := model.UpdateCreation(&crea); err != nil {
		c.Error(err).SetMeta(ErrDB)
		return
	}

	version := fmt.Sprintf("%d", model.CreationLastVersion(crea.ID))
	if err := crea.RetrieveSourceCode(version, enum.Script, enum.Document, enum.Style); err != nil {
		c.Error(err)
	}
	buildPreview(&crea, fmt.Sprintf("%d", crea.CreatorID), version)

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
	crea, err := model.CreationPrivateByID(user.(*model.User).ID, creaID, "latest")
	if err != nil {
		c.Error(err).SetMeta(ErrResNotFound.SetParams("source", "creation", "id", creaID))
		return
	}

	if crea.State != enum.Draft {
		// TODO state error
		return
	}

	userIDStr := fmt.Sprintf("%d", user.(*model.User).ID)
	version := fmt.Sprintf("%d", crea.Version)
	creaIDStr := fmt.Sprintf("%d", crea.ID.ValueDecoded)
	storage := lib.NewStorage(lib.SrcCreations)

	storage.StoreFile(codeForm.Document, "text/html", userIDStr, creaIDStr, version, enum.Document)
	storage.StoreFile(codeForm.Script, "application/javascript", userIDStr, creaIDStr, version, enum.Script)
	storage.StoreFile(codeForm.ParsedScript, "application/javascript", userIDStr, creaIDStr, version, enum.ParsedScript)
	storage.StoreFile(codeForm.Style, "text/css", userIDStr, creaIDStr, version, enum.Style)

	minifier := minify.New()
	minifier.AddFunc("text/javascript", js.Minify)

	var minErr error
	codeForm.Script, minErr = minifier.String("text/javascript", codeForm.Script)

	if minErr != nil {
		c.Error(minErr).SetMeta(ErrServ.SetParams("source", "minifier"))
		return
	}

	crea.Script = codeForm.Script
	crea.Document = codeForm.Document
	crea.Style = codeForm.Style
	crea.Params = codeForm.Params

	model.UpdateCreationParams(crea)

	buildPreview(crea, userIDStr, version)

	if storage.Error() != nil {
		c.Error(storage.Error()).SetMeta(ErrServ.SetParams("source", "files"))
		return
	}

	if codeForm.ParsedScript != "" {
		if _, errs := engine.NewJS("control", codeForm.ParsedScript, make([]engine.JSParam, 0)); len(errs) > 0 {
			switch errs[0] {
			case engine.ErrNoClassFound:
				c.Error(errs[0]).SetMeta(ErrBadScriptClass)
			case engine.ErrNoConstructor:
				c.Error(errs[0]).SetMeta(ErrBadScriptConst.SetParams("example", "constructor() { }"))
			case engine.ErrNoDocInit:
				c.Error(errs[0]).SetMeta(ErrBadScriptDoc.SetParams("example", "this.document = document"))
			}
			return
		}
	}

	c.Header("Location", fmt.Sprintf("/creations/%s", creaID.ValueEncoded))

	c.AbortWithStatus(NoContent)
}

// POSTCreationVersion creates a new version
func POSTCreationVersion(c *gin.Context) {
	user, _ := c.Get("user")
	uID := user.(*model.User).ID

	creaID := lib.InitID(c.Param("encid"))
	crea, err := model.CreationPrivateByID(uID, creaID, "latest")
	if err != nil {
		c.Error(err).SetMeta(ErrResNotFound.SetParams("source", "creation", "id", creaID.ValueEncoded))
		return
	}

	newVersion := crea.Version + 1

	uIDStr := fmt.Sprintf("%d", crea.CreatorID)
	creaIDStr := fmt.Sprintf("%d", crea.ID.ValueDecoded)
	storage := lib.NewStorage(lib.SrcCreations)
	curVersionStr := fmt.Sprintf("%d", crea.Version)
	newVersionStr := fmt.Sprintf("%d", newVersion)

	storage.CopyAndStoreFile(uIDStr, creaIDStr, curVersionStr, newVersionStr, enum.Script)
	storage.CopyAndStoreFile(uIDStr, creaIDStr, curVersionStr, newVersionStr, enum.Document)
	storage.CopyAndStoreFile(uIDStr, creaIDStr, curVersionStr, newVersionStr, enum.Style)

	if storage.Error() != nil {
		c.Error(storage.Error()).SetMeta(ErrServ.SetParams("source", "copy"))
		return
	}

	if err := model.NewCreationVersion(crea, newVersion); err != nil {
		c.Error(err).SetMeta(ErrDB)
		return
	}

	c.Header("Location", fmt.Sprintf("/creations/%s/code", c.Param("encid")))

	c.AbortWithStatus(NoContent)
}

func buildPreview(crea *model.Creation, userID string, version string) {
	params := ""

	for _, p := range crea.Params {
		params += fmt.Sprintf(`"%s":%s,`, p.Field, p.Value)
	}
	params = strings.TrimRight(params, ",")

	preview := `<html>
		<head>
			<script type="text/javascript">` + crea.Script + `</script>
			<script type="text/javascript" src="https://cdnjs.cloudflare.com/ajax/libs/webcomponentsjs/1.0.0-rc.11/webcomponents-lite.js"></script>
			<script type="text/javascript">window.onload = function(){new Woobly({` + params + `});}</script>
			<style>html, body { height: 100%; width: 100%; margin: 0; }` + crea.Style + `</style>
		</head>
		<body>
			` + crea.Document + `
		</body>
	</html>`

	storage := lib.NewStorage(lib.SrcPreview)

	creaIDStr := fmt.Sprintf("%d", crea.ID.ValueDecoded)

	storage.StoreFile(preview, "text/html", userID, creaIDStr, version, enum.Preview)
}
