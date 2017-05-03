package handler

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"strings"
	"wooble/lib"
	model "wooble/models"

	"github.com/gin-gonic/gin"
	"github.com/nfnt/resize"
)

// POSTFile is a handler to upload a file
func POSTFile(c *gin.Context) {
	var res struct {
		Path string `json:"path"`
	}
	file, header, _ := c.Request.FormFile("file")
	mimeType := header.Header.Get("Content-Type")

	var source string
	var sizeW uint
	switch c.Query("source") {
	case "profile":
		source = lib.SrcProfile
		sizeW = 128
	case "crea_thumb":
		source = lib.SrcCreaThumb
		sizeW = 350
	default:
		c.AbortWithStatus(NoContent)
		return
	}

	image, _, _ := image.Decode(file)
	buff := new(bytes.Buffer)

	switch mimeType {
	case "image/jpeg":
		newImage := resize.Resize(sizeW, 0, image, resize.Lanczos3)
		jpeg.Encode(buff, newImage, nil)
	case "image/png":
		newImage := resize.Resize(sizeW, 0, image, resize.Lanczos3)
		png.Encode(buff, newImage)
	case "image/gif":
		newImage := resize.Resize(sizeW, 0, image, resize.Lanczos3)
		gif.Encode(buff, newImage, nil)
	default:
		c.Error(errors.New("File upload is not an image")).SetMeta(ErrBadFileFormat.SetParams("formats", "image/jpeg, image/gif, image/png"))
		return
	}

	id := c.Query("id")

	storage := lib.NewStorage(source)
	user, _ := c.Get("user")
	res.Path = storage.StoreFile(buff, mimeType, fmt.Sprintf("%d", user.(*model.User).ID), source+id, "", header.Filename)

	fmt.Print(storage.Error())

	splPath := strings.Split(res.Path, "/")
	res.Path = strings.Join(splPath[1:len(splPath)], "/")

	c.JSON(OK, NewRes(res))
}
