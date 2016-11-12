package handler

import (
	"fmt"
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
			res.Error(ErrResNotFound, "Creation", title)
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
		res.Error(ErrBadForm, "creatorId (int) and title (string) are required")
		c.JSON(res.HttpStatus(), res)
		return
	}

	if data.Version == "" {
		data.Version = model.BASE_VERSION
	}

	// TODO Authenticated user and put in CreatorID
	user, _ := model.UserByID(1)

	// TODO delete row if push to bucket fails
	_, err := model.NewCreation(&data)
	if err != nil {
		res.Error(ErrDBSave, "- Title should be unique for the creator\n - Creator should exist")
		fmt.Print(err)
		c.JSON(res.HttpStatus(), res)
		return
	}

	eng, _ := model.EngineByName(data.Engine)

	storage := lib.NewStorage(lib.SrcCreations, user.Name, data.Version)

	if data.Document != "" {
		storage.StoreFile(data.Document, eng.ContentType, data.Title, "doc.html")
	}
	if data.Script != "" {
		storage.StoreFile(data.Script, eng.ContentType, data.Title, "script"+eng.Extension)
	}
	if data.Style != "" {
		storage.StoreFile(data.Style, eng.ContentType, data.Title, "style.css")
	}

	if storage.Error != nil {
		res.Error(ErrServ, "doc, script and style files")
	}

	res.Status = Created

	c.JSON(res.HttpStatus(), res)
}
