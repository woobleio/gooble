package handler

import (
	"database/sql"
	"fmt"
	"strings"

	"wooble/lib"
	"wooble/models"
	enum "wooble/models/enums"

	"github.com/woobleio/wooblizer/wbzr"
	"github.com/woobleio/wooblizer/wbzr/engine"
	"gopkg.in/gin-gonic/gin.v1"
)

// GETPackages is a handler that returns one or more packages
func GETPackages(c *gin.Context) {
	var data interface{}
	var err error

	user, _ := c.Get("user")

	pkgID := c.Param("encid")

	if pkgID != "" {
		data, err = model.PackageByID(pkgID, user.(*model.User).ID)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Error(err).SetMeta(ErrResNotFound.SetParams("source", "Package", "id", pkgID))
			} else {
				c.Error(err).SetMeta(ErrDBSelect)
			}
		}
	} else {
		opts := lib.ParseOptions(c)
		data, err = model.AllPackages(opts, user.(*model.User).ID)
		if err != nil {
			c.Error(err).SetMeta(ErrDBSelect)
			return
		}
	}

	c.JSON(OK, NewRes(data))
}

// POSTPackages is a handler that create am empty Wooble package
func POSTPackages(c *gin.Context) {
	var data model.PackageForm

	if err := c.BindJSON(&data); err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind).SetMeta(ErrBadForm)
		return
	}

	user, _ := c.Get("user")
	data.UserID = user.(*model.User).ID

	plan := user.(*model.User).Plan

	limitNbPkg := plan.NbPkg.Int64
	limitNbDomains := plan.NbDomains.Int64
	userNbPkg := model.UserNbPackages(data.UserID)

	// 0 means unlimited
	if limitNbPkg != 0 && userNbPkg >= limitNbPkg {
		c.Error(nil).SetMeta(ErrPlanLimit.SetParams("source", "package", "plan", plan.Label.String))
		return
	}

	if limitNbDomains != 0 && int64(len(data.Domains)) > limitNbDomains {
		c.Error(nil).SetMeta(ErrPlanLimit.SetParams("source", "domains", "plan", plan.Label.String))
		return
	}

	pkgID, err := model.NewPackage(&data)
	if err != nil {
		c.Error(err).SetMeta(ErrDBSave)
		return
	}

	c.Header("Location", fmt.Sprintf("/%s/%v", "packages", pkgID))

	c.JSON(Created, nil)
}

// PushCreation is an handler that pushes one or more creations in a package
func PushCreation(c *gin.Context) {
	var packageCreationForm model.PackageCreationForm

	if err := c.BindJSON(&packageCreationForm); err != nil {
		c.Error(err).SetMeta(ErrBadForm)
		return
	}

	if packageCreationForm.Version == "" {
		packageCreationForm.Version = model.BaseVersion
	}

	crea, err := model.CreationByIDAndVersion(packageCreationForm.CreationID, packageCreationForm.Version)
	if err != nil {
		c.Error(err).SetMeta(ErrResNotFound.SetParams("source", "creation", "id", packageCreationForm.CreationID+"/"+packageCreationForm.Version))
		return
	}

	if crea.State == enum.Draft && crea.Versions[len(crea.Versions)-1] == packageCreationForm.Version {
		c.Error(err).SetMeta(ErrCreationNotAvail.SetParams("id", crea.ID.ValueEncoded+"/"+packageCreationForm.Version))
		return
	}

	pkgID := c.Param("encid")

	user, _ := c.Get("user")

	pkg, err := model.PackageByID(pkgID, user.(*model.User).ID)
	if err != nil {
		c.Error(err).SetMeta(ErrResNotFound.SetParams("source", "Package", "id", pkgID))
		return
	}

	plan := user.(*model.User).Plan
	limitNbCrea := plan.NbCrea.Int64
	pkgNbCrea := model.PackageNbCrea(pkg.ID.ValueEncoded)

	if limitNbCrea != 0 && pkgNbCrea >= limitNbCrea {
		c.Error(nil).SetMeta(ErrPlanLimit.SetParams("source", "creation", "plan", plan.Label.String))
		return
	}

	if err := model.PushCreation(pkg.ID.ValueDecoded, &packageCreationForm); err != nil {
		c.Error(err).SetMeta(ErrDBSave)
		return
	}

	c.Header("Location", fmt.Sprintf("/%s/%s", "packages", pkgID))

	c.JSON(Created, nil)
}

// BuildPackage is a handler action that builds the Wooble lib of a package
// (a Wooble lib is a file that bundles everything contained in a package,
// the file is stored in the cloud)
func BuildPackage(c *gin.Context) {
	pkgID := c.Param("encid")

	user, _ := c.Get("user")

	pkg, err := model.PackageByID(pkgID, user.(*model.User).ID)
	if err != nil {
		c.Error(err).SetMeta(ErrResNotFound.SetParams("source", "Package", "id", pkgID))
		return
	}

	// Check if at least one creation should be bought to be build in the package
	for _, creation := range pkg.Creations {
		if creation.IsToBuy {
			c.Error(nil).SetMeta(ErrMustBuy)
			return
		}
	}

	storage := lib.NewStorage(lib.SrcPackages)

	storage.Source = lib.SrcCreations
	wb := wbzr.New(wbzr.JSES5)
	for _, creation := range pkg.Creations {
		var script engine.Script

		creatorIDStr := fmt.Sprintf("%d", creation.CreatorID)

		objName := creation.Title
		if creation.Alias != nil {
			objName = creation.Alias.String
		}

		if creation.HasScript {
			src := storage.GetFileContent(creatorIDStr, creation.ID.ValueEncoded, creation.Version, "script.js")

			script, err = wb.Inject(src, objName)
		} else {
			script, err = wb.Inject("", objName)
		}

		if err != nil {
			if err == wbzr.ErrUniqueName {
				c.Error(err).SetMeta(ErrAliasRequired.SetParams("name", creation.Title))
				return
			}
			panic(err)
		}

		if creation.HasDoc {
			src := storage.GetFileContent(creatorIDStr, creation.ID.ValueEncoded, creation.Version, "doc.html")
			err = script.IncludeHtml(src)
		}
		if creation.HasStyle {
			src := storage.GetFileContent(creatorIDStr, creation.ID.ValueEncoded, creation.Version, "style.css")
			err = script.IncludeCss(src)
		}

		if err != nil {
			panic(err)
		}
	}

	storage.Source = lib.SrcPackages

	bf, err := wb.SecureAndWrap(pkg.Domains...)

	if err != nil || storage.Error != nil {
		c.Error(storage.Error).SetMeta(ErrServ.SetParams("source", "package"))
		return
	}

	// TODO if multitype allowed, package should have an engine too
	path := storage.StoreFile(bf.String(), "application/javascript", fmt.Sprintf("%d", user.(*model.User).ID), pkg.ID.ValueEncoded, "", "wooble.js")

	spltPath := strings.Split(path, "/")
	spltPath[0] = ""

	source := "https://pkg.wooble.io" + strings.Join(spltPath, "/")
	if err := model.UpdatePackageSource(source, pkg.ID); err != nil {
		c.Error(err).SetMeta(ErrUpdate.SetParams("source", "package", "id", pkg.ID.ValueEncoded))
		return
	}

	pkg.Source = lib.InitNullString(source)

	c.JSON(OK, NewRes(pkg))
}
