package handler

import (
	"database/sql"
	"wooble/lib"
	"wooble/model"

	"gopkg.in/gin-gonic/gin.v1"
)

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

	res.Response(&data)

	c.JSON(res.HttpStatus(), res)
}

func POSTCreations(c *gin.Context) {
	var data model.CreationForm

	res := NewRes()

	// FIXME workaround gin issue with Bind (https://github.com/gin-gonic/gin/issues/633)
	c.Header("Content-Type", gin.MIMEJSON)
	if c.BindJSON(&data) != nil {
		res.Error(ErrBadForm, "title (string) is required")
		c.JSON(res.HttpStatus(), res)
		return
	}

	if data.Version == "" {
		data.Version = model.BASE_VERSION
	}

	// TODO Authenticated user and put in CreatorID
	user, _ := model.UserByID(4)
	data.CreatorID = 4

	creaId, err := model.NewCreation(&data)
	if err != nil {
		res.Error(ErrDBSave, "- Title should be unique for the creator\n")
		c.JSON(res.HttpStatus(), res)
		return
	}

	eng, _ := model.EngineByName(data.Engine)

	storage := lib.NewStorage(lib.SrcCreations, user.Name, data.Version)

	if data.Document != "" {
		storage.StoreFile(data.Document, "text/html", data.Title, "doc.html")
	}
	if data.Script != "" {
		storage.StoreFile(data.Script, eng.ContentType, data.Title, "script"+eng.Extension)
	}
	if data.Style != "" {
		storage.StoreFile(data.Style, "text/css", data.Title, "style.css")
	}

	if storage.Error != nil {
		// Delete the crea since files failed to be save in the cloud
		model.DeleteCreation(creaId)
		res.Error(ErrServ, "doc, script and style files")
	}

	res.Status = Created

	c.JSON(res.HttpStatus(), res)
}
