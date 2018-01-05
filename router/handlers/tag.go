package handler

import (
	"wooble/lib"
	model "wooble/models"

	"github.com/gin-gonic/gin"
)

// GETTags is a handler that returns tags
func GETTags(c *gin.Context) {
	opts := lib.ParseOptions(c)

	tags, err := model.AllTags(opts)
	if err != nil {
		c.Error(err).SetMeta(ErrDB)
		return
	}

	c.JSON(OK, NewRes(tags))
}
