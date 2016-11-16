package handler

import (
	"fmt"
	"strconv"

	"wooble/model"

	"gopkg.in/gin-gonic/gin.v1"
)

func POSTPackages(c *gin.Context) {
	var data model.PackageForm

	res := NewRes()

	// FIXME workaround gin issue with Bind (https://github.com/gin-gonic/gin/issues/633)
	c.Header("Content-Type", gin.MIMEJSON)
	if c.BindJSON(&data) != nil {
		res.Error(ErrBadForm, "title (string) is required")
		c.JSON(res.HttpStatus(), res)
		return
	}

	// TODO Authenticated user and put in CreatorID
	data.UserID = 5

	_, err := model.NewPackage(&data)
	if err != nil {
		res.Error(ErrDBSave, "- Title should be unique for the creator\n")
		c.JSON(res.HttpStatus(), res)
		return
	}

	res.Status = Created

	c.JSON(res.HttpStatus(), res)
}

func PushCreations(c *gin.Context) {
	type PackageCreationForm struct {
		PackageID  uint64
		CreationID []uint64 `json:"creations" binding:"required"`
	}

	var data PackageCreationForm

	res := NewRes()

	// FIXME workaround gin issue with Bind (https://github.com/gin-gonic/gin/issues/633)
	c.Header("Content-Type", gin.MIMEJSON)
	if c.BindJSON(&data) != nil {
		res.Error(ErrBadForm, "creations (int[]) is required")
		c.JSON(res.HttpStatus(), res)
		return
	}

	param := c.Param("id")
	pkgID, err := strconv.ParseUint(param, 10, 64)
	if err != nil {
		res.Error(ErrBadParam, "int")
		c.JSON(res.HttpStatus(), res)
		return
	}

	for _, creaID := range data.CreationID {
		if err := model.PushCreation(pkgID, creaID); err != nil {
			res.Error(ErrDBSave, fmt.Sprintf("failed to push creation %v in the package", creaID))
		}
	}

	res.Status = Created

	c.JSON(res.HttpStatus(), res)
}
