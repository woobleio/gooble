package handler

import (
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
		data, err = model.CreationByTitle(title, opts)
		if err != nil {
			res.Error(ErrResNotFound, "Creation", title)
		}
	} else {
		data, err = model.AllCreations(opts)
		if err != nil {
			res.Error(ErrNotFound, "creations")
		}
	}

	res.Response(&data)

	c.JSON(res.HttpStatus(), res)
}
