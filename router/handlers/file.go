package handler

import (
	"errors"
	"fmt"
	"strings"
	"wooble/lib"
	model "wooble/models"

	"github.com/gin-gonic/gin"
)

// POSTFile is a handler to upload a file
func POSTFile(c *gin.Context) {
	var res struct {
		Path string `json:"path"`
	}
	file, header, _ := c.Request.FormFile("file")
	mimeType := header.Header.Get("Content-Type")

	if mimeType != "image/jpeg" && mimeType != "image/gif" && mimeType != "image/png" {
		c.Error(errors.New("File upload is not an image")).SetMeta(ErrBadFileFormat.SetParams("formats", "image/jpeg, image/gif, image/png"))
		return
	}

	user, _ := c.Get("user")
	storage := lib.NewStorage(lib.SrcProfile)
	res.Path = storage.StoreFile(file, mimeType, fmt.Sprintf("%d", user.(*model.User).ID), lib.SrcProfile, "", header.Filename)

	splPath := strings.Split(res.Path, "/")
	res.Path = strings.Join(splPath[1:len(splPath)], "/")

	c.JSON(OK, NewRes(res))
}
