package handler

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	form "wooble/forms"
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
		data, err = model.PackageByID(user.(*model.User).ID, lib.InitID(pkgID))
		if err != nil {
			if err == sql.ErrNoRows {
				c.Error(err).SetMeta(ErrResNotFound.SetParams("source", "Package", "id", pkgID))
			} else {
				c.Error(err).SetMeta(ErrDB)
			}
		}
	} else {
		opts := lib.ParseOptions(c)
		data, err = model.AllPackages(&opts, user.(*model.User).ID)
		if err != nil {
			c.Error(err).SetMeta(ErrDB)
			return
		}
	}

	c.JSON(OK, NewRes(data))
}

// POSTPackage is a handler that create am empty Wooble package
func POSTPackage(c *gin.Context) {
	var data form.PackageForm

	if err := c.BindJSON(&data); err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind).SetMeta(ErrBadForm)
		return
	}

	user, _ := c.Get("user")

	pkg := new(model.Package)
	pkg.UserID = user.(*model.User).ID

	plan := user.(*model.User).Plan

	limitNbPkg := plan.NbPkg.Int64
	limitNbDomains := plan.NbDomains.Int64
	userNbPkg := model.UserNbPackages(pkg.UserID)

	// 0 means unlimited
	if limitNbPkg != 0 && userNbPkg >= limitNbPkg {
		c.Error(nil).SetMeta(ErrPlanLimit.SetParams("source", "package", "plan", plan.Label.String))
		return
	}

	if limitNbDomains != 0 && int64(len(data.Domains)) > limitNbDomains {
		c.Error(nil).SetMeta(ErrPlanLimit.SetParams("source", "domains", "plan", plan.Label.String))
		return
	}

	pkg.Title = data.Title
	pkg.Domains = data.Domains

	var errPkg error
	pkg, errPkg = model.NewPackage(pkg)
	if errPkg != nil {
		c.Error(errPkg).SetMeta(ErrDB)
		return
	}

	c.Header("Location", fmt.Sprintf("/packages/%v", pkg.ID.ValueEncoded))

	c.JSON(Created, NewRes(pkg))
}

// DELETEPackage delete a package
func DELETEPackage(c *gin.Context) {
	user, _ := c.Get("user")

	uID := user.(*model.User).ID
	pkgID := lib.InitID(c.Param("encid"))

	if err := model.DeletePackage(uID, pkgID); err != nil {
		c.Error(err).SetMeta(ErrDB)
		return
	}

	storage := lib.NewStorage(lib.SrcPackages)
	storage.DeleteFile(fmt.Sprintf("%d", uID), fmt.Sprintf("%d", pkgID.ValueDecoded), "", enum.Wooble)

	if storage.Error() != nil {
		c.Error(storage.Error())
	}

	c.Header("Location", "/packages")

	c.AbortWithStatus(NoContent)
}

// PUTPackage is an handler that updates a package
func PUTPackage(c *gin.Context) {
	var data form.PackageForm

	if err := c.BindJSON(&data); err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind).SetMeta(ErrBadForm)
		return
	}

	user, _ := c.Get("user")

	pkg := new(model.Package)
	pkg.ID = lib.InitID(c.Param("encid"))
	pkg.UserID = user.(*model.User).ID
	pkg.Title = data.Title

	plan := user.(*model.User).Plan
	limitNbDomains := plan.NbDomains.Int64

	if limitNbDomains != 0 && int64(len(data.Domains)) > limitNbDomains {
		c.Error(nil).SetMeta(ErrPlanLimit.SetParams("source", "domains", "plan", plan.Label.String))
		return
	}

	pkg.Domains = data.Domains

	if err := model.UpdatePackage(pkg); err != nil {
		c.Error(err).SetMeta(ErrDB)
		return
	}

	c.Header("Location", fmt.Sprintf("/packages/%s", pkg.ID.ValueEncoded))

	c.AbortWithStatus(NoContent)
}

// PUTPackageCreation updates a package creation
func PUTPackageCreation(c *gin.Context) {
	var data form.PackageCreationForm

	if err := c.BindJSON(&data); err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind).SetMeta(ErrBadForm)
		return
	}

	user, _ := c.Get("user")

	crea := new(model.Creation)
	pkgCrea := new(model.Package)
	pkgCrea.ID = lib.InitID(c.Param("encid"))
	pkgCrea.UserID = user.(*model.User).ID
	crea.Alias = lib.InitNullString(data.Alias)
	crea.Version = data.Version
	pkgCrea.Creations = []model.Creation{*crea}

	if err := model.UpdatePackageCreation(pkgCrea); err != nil {
		c.Error(err).SetMeta(ErrDB)
		return
	}

	c.Header("Location", fmt.Sprintf("/packages/%s", pkgCrea.ID.ValueEncoded))

	c.AbortWithStatus(NoContent)
}

// PushCreation is an handler that pushes one or more creations in a package
func PushCreation(c *gin.Context) {
	var pkgCreaForm form.PackageCreationForm

	if err := c.BindJSON(&pkgCreaForm); err != nil {
		c.Error(err).SetMeta(ErrBadForm)
		return
	}

	crea, err := model.CreationByIDAndVersion(lib.InitID(pkgCreaForm.CreationID), pkgCreaForm.Version)
	if err != nil {
		c.Error(err).SetMeta(ErrResNotFound.SetParams("source", "creation", "id", pkgCreaForm.CreationID+"/"+pkgCreaForm.Version))
		return
	}

	if crea.State == enum.Draft && crea.Versions[len(crea.Versions)-1] == pkgCreaForm.Version {
		c.Error(errors.New("Creation not available")).SetMeta(ErrCreaNotAvail.SetParams("id", crea.ID.ValueEncoded+"/"+pkgCreaForm.Version))
		return
	}

	user, _ := c.Get("user")

	pkgID := lib.InitID(c.Param("encid"))
	pkg, err := model.PackageByID(user.(*model.User).ID, pkgID)
	if err != nil {
		c.Error(err).SetMeta(ErrResNotFound.SetParams("source", "Package", "id", pkgID.ValueEncoded))
		return
	}

	plan := user.(*model.User).Plan
	limitNbCrea := uint64(plan.NbCrea.Int64)
	pkgNbCrea := model.PackageNbCrea(pkg.ID)

	if limitNbCrea != 0 && pkgNbCrea >= limitNbCrea {
		c.Error(nil).SetMeta(ErrPlanLimit.SetParams("source", "creation", "plan", plan.Label.String))
		return
	}

	if pkgCreaForm.Version == "" {
		pkgCreaForm.Version = crea.Versions[len(crea.Versions)-1]
	}

	if err := model.NewPackageCreation(pkg.ID, crea.ID, pkgCreaForm.Version); err != nil {
		c.Error(err).SetMeta(ErrDB)
		return
	}

	c.Header("Location", fmt.Sprintf("/packages/%s", pkg.ID.ValueEncoded))

	c.AbortWithStatus(NoContent)
}

// BuildPackage is a handler action that builds the Wooble lib of a package
// (a Wooble lib is a file that bundles everything contained in a package,
// the file is stored in the cloud)
func BuildPackage(c *gin.Context) {
	user, _ := c.Get("user")

	pkgID := lib.InitID(c.Param("encid"))
	pkg, err := model.PackageByID(user.(*model.User).ID, pkgID)
	if err != nil {
		c.Error(err).SetMeta(ErrResNotFound.SetParams("source", "Package", "id", pkg.ID.ValueEncoded))
		return
	}

	// Check if at least one creation should be bought to be build in the package
	for _, creation := range pkg.Creations {
		if creation.IsToBuy {
			c.Error(nil).SetMeta(ErrMustBuy)
			return
		}
	}

	storage := lib.NewStorage(lib.SrcCreations)

	wb := wbzr.New(wbzr.JSES5)
	for _, creation := range pkg.Creations {
		var script engine.Script

		creatorIDStr := fmt.Sprintf("%d", creation.CreatorID)

		objName := creation.Title
		if creation.Alias != nil {
			objName = creation.Alias.String
		}

		creaIDStr := fmt.Sprintf("%d", creation.ID.ValueDecoded)
		src := storage.GetFileContent(creatorIDStr, creaIDStr, creation.Version, enum.Script)
		script, err = wb.Inject(src, objName)

		if err != nil {
			if err == wbzr.ErrUniqueName {
				c.Error(err).SetMeta(ErrAliasRequired.SetParams("name", creation.Title))
				return
			}
			panic(err)
		}

		if creation.HasDoc {
			src := storage.GetFileContent(creatorIDStr, creaIDStr, creation.Version, enum.Document)
			err = script.IncludeHtml(src)
		}
		if creation.HasStyle {
			src := storage.GetFileContent(creatorIDStr, creaIDStr, creation.Version, enum.Style)
			err = script.IncludeCss(src)
		}

		if err != nil {
			panic(err)
		}
	}

	storage.Source = lib.SrcPackages

	bf, err := wb.SecureAndWrap(pkg.Domains...)

	if err != nil || storage.Error() != nil {
		c.Error(storage.Error()).SetMeta(ErrServ.SetParams("source", "package"))
		return
	}

	// TODO if multitype allowed, package should have an engine too
	path := storage.StoreFile(bf.String(), "application/javascript", fmt.Sprintf("%d", user.(*model.User).ID), fmt.Sprintf("%d", pkg.ID.ValueDecoded), "", enum.Wooble)

	spltPath := strings.Split(path, "/")
	spltPath[0] = ""

	source := "https://pkg.wooble.io" + strings.Join(spltPath, "/")
	if err := model.UpdatePackageSource(pkg.UserID, pkg.ID, source); err != nil {
		c.Error(err).SetMeta(ErrDB)
		return
	}

	pkg.Source = lib.InitNullString(source)

	c.JSON(OK, NewRes(pkg))
}

// RemovePackageCreation remove a creation from a package
func RemovePackageCreation(c *gin.Context) {
	var data form.PackageCreationForm

	user, _ := c.Get("user")

	if err := c.BindJSON(&data); err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind).SetMeta(ErrBadForm)
		return
	}

	if err := model.DeletePackageCreation(user.(*model.User).ID, lib.InitID(c.Param("encid")), lib.InitID(data.CreationID)); err != nil {
		c.Error(err).SetMeta(ErrDB)
		return
	}

	c.AbortWithStatus(NoContent)
}
