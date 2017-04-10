package handler

import (
	"fmt"
	"wooble/lib"
	model "wooble/models"

	"github.com/gin-gonic/gin"
)

// POSTFile is a handler to upload a file
func POSTFile(c *gin.Context) {
	file, header, _ := c.Request.FormFile("file")
	mimeType := header.Header.Get("Content-Type")

	if mimeType != "image/jpeg" && mimeType != "image/gif" && mimeType != "image/png" {
		// error not good format
		return
	}

	user, _ := c.Get("user")
	storage := lib.NewStorage(lib.SrcProfile)
	storage.StoreFile(file, mimeType, fmt.Sprintf("%d", user.(*model.User).ID), lib.SrcProfile, "", header.Filename)

	// TODO return profile path (build for user on the fly)

	c.AbortWithStatus(NoContent)
}
