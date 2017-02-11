package handler

import (
	"database/sql"
	"fmt"
	"strings"

	"wooble/lib"
	"wooble/models"

	"github.com/woobleio/wooblizer/wbzr"
	"github.com/woobleio/wooblizer/wbzr/engine"
	"gopkg.in/gin-gonic/gin.v1"
)

// GETPackages is a handler that returns one or more packages
func GETPackages(c *gin.Context) {
	var data interface{}
	var err error

	res := NewRes()

	opts := lib.ParseOptions(c)

	user, _ := c.Get("user")

	pkgID := c.Param("id")

	if pkgID != "" {
		data, err = model.PackageByID(pkgID, user.(*model.User).ID)
		if err != nil {
			if err == sql.ErrNoRows {
				res.Error(ErrResNotFound, "Package", pkgID)
			} else {
				res.Error(ErrDBSelect)
			}
		}
	} else {
		data, err = model.AllPackages(opts, user.(*model.User).ID)
		if err != nil {
			res.Error(ErrDBSelect)
		}
	}

	res.Response(data)

	c.JSON(res.HTTPStatus(), res)

}

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

	plan := user.(*model.User).Plan

	limitNbPkg := plan.NbPkg.Int64
	limitNbDomains := plan.NbDomains.Int64
	userNbPkg := model.UserNbPackages(data.UserID)

	// 0 means unlimited
	if limitNbPkg != 0 && userNbPkg >= limitNbPkg {
		res.Error(ErrPlanLimit, "Packages", plan.Label.String)
		c.JSON(res.HTTPStatus(), res)
		return
	}

	if limitNbDomains != 0 && int64(len(data.Domains)) > limitNbDomains {
		res.Error(ErrPlanLimit, "Domains per package", plan.Label.String)
		c.JSON(res.HTTPStatus(), res)
		return
	}

	pkgID, err := model.NewPackage(&data)
	if err != nil {
		res.Error(ErrDBSave, "- Title should be unique for the creator")
		c.JSON(res.HTTPStatus(), res)
		return
	}

	c.Header("Location", fmt.Sprintf("/%s/%v", "packages", pkgID))

	res.Status = Created

	c.JSON(res.HTTPStatus(), res)
}

// PushCreation is an handler that pushes one or more creations in a package
func PushCreation(c *gin.Context) {
	type PackageCreationForm struct {
		PackageID  uint64
		CreationID string `json:"creation" binding:"required"`
	}

	var data PackageCreationForm

	res := NewRes()

	// FIXME workaround gin issue with Bind (https://github.com/gin-gonic/gin/issues/633)
	c.Header("Content-Type", gin.MIMEJSON)
	if c.BindJSON(&data) != nil {
		res.Error(ErrBadForm, "creations (string) is required")
		c.JSON(res.HTTPStatus(), res)
		return
	}

	pkgID := c.Param("id")

	user, _ := c.Get("user")

	pkg, err := model.PackageByID(pkgID, user.(*model.User).ID)
	if err != nil {
		res.Error(ErrResNotFound, "package", "")
		c.JSON(res.HTTPStatus(), res)
		return
	}

	plan := user.(*model.User).Plan
	limitNbCrea := plan.NbCrea.Int64
	pkgNbCrea := model.PackageNbCrea(pkg.ID.ValueEncoded)

	if limitNbCrea != 0 && pkgNbCrea >= limitNbCrea {
		res.Error(ErrPlanLimit, "Creations per package", plan.Label.String)
		c.JSON(res.HTTPStatus(), res)
		return
	}

	if err := model.PushCreation(pkg.ID.ValueDecoded, data.CreationID); err != nil {
		res.Error(ErrDBSave, fmt.Sprintf("failed to push creation %v in the package", data.CreationID))
	}

	c.Header("Location", fmt.Sprintf("/%s/%s", "packages", pkgID))

	res.Status = Created

	c.JSON(res.HTTPStatus(), res)
}

// BuildPackage is a handler action that builds the Wooble lib of a package
// (a Wooble lib is a file that bundles everything contained in a package,
// the file is stored in the cloud)
func BuildPackage(c *gin.Context) {
	res := NewRes()

	pkgID := c.Param("id")

	user, _ := c.Get("user")

	pkg, err := model.PackageByID(pkgID, user.(*model.User).ID)
	if err != nil {
		res.Error(ErrResNotFound, "package", pkgID)
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
			src := storage.GetFileContent(creation.Creator.Name, creation.Title, "script.js", "")

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

	// FIXME this change could crash
	// pkg.Source = lib.InitNullString("https://pkg.wooble.io" + strings.Join(spltPath, "/"))

	if err := model.UpdatePackageSource("https://pkg.wooble.io"+strings.Join(spltPath, "/"), pkg.ID); err != nil {
		res.Error(ErrUpdate, "package", pkg.ID)
	}

	res.Response(pkg)

	res.Status = OK

	c.JSON(res.HTTPStatus(), res)
}
