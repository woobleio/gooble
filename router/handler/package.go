package handler

import (
	"fmt"
	"strings"

	"wooble/lib"
	"wooble/model"

	"github.com/woobleio/wooblizer/wbzr"
	"github.com/woobleio/wooblizer/wbzr/engine"
	"gopkg.in/gin-gonic/gin.v1"
)

// POSTPackages is a handler that create am empty Wooble package
func POSTPackages(c *gin.Context) {
	var data model.PackageForm

	res := NewRes()

	// FIXME workaroun gin issue with Bind (https://github.com/gin-gonic/gin/issues/633)
	c.Header("Content-Type", gin.MIMEJSON)
	if c.BindJSON(&data) != nil {
		res.Error(ErrBadForm, "title (string) is required")
		c.JSON(res.HTTPStatus(), res)
		return
	}

	user, _ := c.Get("user")
	data.UserID = user.(*model.User).ID
	data.Key = lib.GenKey()

	pkgID, err := model.NewPackage(&data)
	if err != nil {
		res.Error(ErrDBSave, "- Title should be unique for the creator")
		c.JSON(res.HTTPStatus(), res)
		return
	}

	c.Header("Location", fmt.Sprintf("/%s/%s/%s/%v", "users", user.(*model.User).Name, "packages", pkgID))

	res.Status = Created

	c.JSON(res.HTTPStatus(), res)
}

// PushCreations is an handler that pushes one or more creations in a package
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
		c.JSON(res.HTTPStatus(), res)
		return
	}

	pkgID := c.Param("id")
	username := c.Param("username")

	pkg, err := model.PackageByID(pkgID)
	if err != nil {
		res.Error(ErrResNotFound, "package", "")
		c.JSON(res.HTTPStatus(), res)
		return
	}

	user, _ := c.Get("user")
	// TODO ownership change
	if pkg.UserID != user.(*model.User).ID || username != user.(*model.User).Name || pkg.User.Name != username {
		res.Error(ErrNotOwner)
		c.JSON(res.HTTPStatus(), res)
		return
	}

	for _, creaID := range data.CreationID {
		if err := model.PushCreation(pkg.ID, creaID); err != nil {
			res.Error(ErrDBSave, fmt.Sprintf("failed to push creation %v in the package", creaID))
		}
	}

	c.Header("Location", fmt.Sprintf("/%s/%v", "packages", pkgID))

	res.Status = Created

	c.JSON(res.HTTPStatus(), res)
}

// BuildPackage is a handler action that builds the Wooble lib of a package
// (a Wooble lib is a file that bundles everything contained in a package,
// the file is stored in the cloud)
func BuildPackage(c *gin.Context) {
	type Build struct {
		Source string `json:"source"`
	}
	var data Build
	res := NewRes()

	pkgID := c.Param("id")
	username := c.Param("username")

	pkg, err := model.PackageByID(pkgID)
	if err != nil {
		res.Error(ErrResNotFound, "package", pkgID)
		c.JSON(res.HTTPStatus(), res)
		return
	}

	user, _ := c.Get("user")
	// TODO ownership change
	if pkg.UserID != user.(*model.User).ID || username != user.(*model.User).Name || pkg.User.Name != username {
		res.Error(ErrNotOwner)
		c.JSON(res.HTTPStatus(), res)
		return
	}

	storage := lib.NewStorage(lib.SrcPackages, "1.0")

	storage.Source = lib.SrcCreations
	wb := wbzr.New(wbzr.JSES5)
	for _, creation := range pkg.Creations {
		var script engine.Script

		storage.Version = creation.Version

		if creation.HasScript {
			src := storage.GetFileContent(creation.Creator.Name, creation.Title, "script"+creation.Engine.Extension, "")

			script, err = wb.Inject(src, creation.Title)
		} else {
			script, err = wb.Inject("", creation.Title)
		}

		if err != nil {
			panic(err)
		}

		if creation.HasDoc {
			src := storage.GetFileContent(creation.Creator.Name, creation.Title, "doc.html", "")
			err = script.IncludeHtml(src)
		}
		if creation.HasStyle {
			src := storage.GetFileContent(creation.Creator.Name, creation.Title, "style.css", "")
			err = script.IncludeCss(src)
		}

		if err != nil {
			panic(err)
		}
	}

	storage.Source = lib.SrcPackages
	storage.Version = ""

	bf, err := wb.SecureAndWrap(pkg.Domains...)

	if err != nil || storage.Error != nil {
		res.Error(ErrServ, "creations packaging")
		c.JSON(res.HTTPStatus(), res)
		return
	}

	// TODO if multitype allowed, package should have an engine too
	path := storage.StoreFile(bf.String(), "application/javascript", user.(*model.User).Name, pkg.Title, "wooble.js", pkg.Key)

	spltPath := strings.Split(path, "/")
	spltPath[0] = ""
	data.Source = "https://pkg.wooble.io" + strings.Join(spltPath, "/")

	res.Response(data)

	res.Status = OK

	c.JSON(res.HTTPStatus(), res)
}
