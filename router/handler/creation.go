package handler

import (
	"database/sql"
	"fmt"

	"wooble/lib"
	"wooble/model"

	"gopkg.in/gin-gonic/gin.v1"
)

// GETCreations is a handler that returns one or more creations
func GETCreations(c *gin.Context) {
	var data interface{}
	var err error

	res := NewRes()

	opts := lib.ParseOptions(c)
	title := c.Param("title")

	if title != "" {
		data, err = model.CreationByTitle(title)
		if err != nil {
			if err == sql.ErrNoRows {
				res.Error(ErrResNotFound, "Creation", title)
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
		res.Error(ErrDBSave, "- Title should be unique for the creator\n")
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

	if data.Document != "" {
		storage.StoreFile(data.Document, "text/html", user.(*model.User).Name, data.Title, "doc.html", "")
	}
	if data.Script != "" {
		storage.StoreFile(data.Script, eng.ContentType, user.(*model.User).Name, data.Title, "script"+eng.Extension, "")
	}
	if data.Style != "" {
		storage.StoreFile(data.Style, "text/css", user.(*model.User).Name, data.Title, "style.css", "")
	}

	if storage.Error != nil {
		// Delete the crea since files failed to be save in the cloud
		model.DeleteCreation(creaID)
		res.Error(ErrServ, "doc, script and style files")
	}

	c.Header("Location", fmt.Sprintf("/%s/%v", "creations", creaID))

	res.Status = Created

	c.JSON(res.HTTPStatus(), res)
}
