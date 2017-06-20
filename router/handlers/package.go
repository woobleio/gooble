package handler

import (
	"database/sql"
	"fmt"
	"time"

	form "wooble/forms"
	"wooble/lib"
	"wooble/models"
	enum "wooble/models/enums"

	"github.com/gin-gonic/gin"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/js"
	"github.com/woobleio/wooblizer/wbzr"
	"github.com/woobleio/wooblizer/wbzr/engine"
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
			return
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
	userNbPkg := model.UserNbPackages(pkg.UserID)

	// 0 means unlimited
	if limitNbPkg != 0 && userNbPkg >= limitNbPkg {
		c.Error(nil).SetMeta(ErrPlanLimit.SetParams("source", "package", "plan", plan.Label.String))
		return
	}

	pkg.Title = data.Title
	pkg.Referer = lib.InitNullString(data.Referer)

	var errPkg error
	pkg, errPkg = model.NewPackage(pkg)
	if errPkg != nil {
		c.Error(errPkg).SetMeta(ErrDB)
		return
	}

	pkg.CreatedAt = lib.InitNullTime(time.Now())

	c.Header("Location", fmt.Sprintf("/packages/%v", pkg.ID.ValueEncoded))

	c.JSON(Created, NewRes(pkg))
}

// PATCHPackage patches a packages
func PATCHPackage(c *gin.Context) {
	var pkgPatchForm form.PackagePatchForm
	var pkg model.Package

	if err := c.BindJSON(&pkgPatchForm); err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind).SetMeta(ErrBadForm)
		return
	}

	user, _ := c.Get("user")
	pkg.ID = lib.InitID(c.Param("encid"))

	if pkgPatchForm.Operation != nil && *pkgPatchForm.Operation == "build" {
		fullPkg, err := model.PackageByID(user.(*model.User).ID, pkg.ID)
		if err != nil {
			c.Error(err).SetMeta(ErrResNotFound.SetParams("source", "Package", "id", fullPkg.ID.ValueEncoded))
			return
		}

		storage := lib.NewStorage(lib.SrcCreations)

		wb := wbzr.New(wbzr.JS)

		minifier := minify.New()
		minifier.AddFunc("text/javascript", js.Minify)

		for _, creation := range fullPkg.Creations {
			creatorIDStr := fmt.Sprintf("%d", creation.CreatorID)

			objName := creation.Alias

			creaIDStr := fmt.Sprintf("%d", creation.ID.ValueDecoded)
			creaVersionStr := fmt.Sprintf("%d", creation.Version)

			src := storage.GetFileContent(creatorIDStr, creaIDStr, creaVersionStr, enum.ParsedScript)

			// Minify to remove comments and white spaces
			src, _ = minifier.String("text/javascript", src)

			var jsParams = make([]interface{}, len(creation.Params))
			for i, p := range creation.Params {
				jsParams[i] = engine.JSParam(engine.JSParam{Field: p.Field, Value: p.Value})
			}

			script, errsScript := wb.Inject(src, objName, jsParams)

			if len(errsScript) > 0 {
				if errsScript[0] == wbzr.ErrUniqueName {
					c.Error(errsScript[0]).SetMeta(ErrAliasRequired.SetParams("name", creation.Title))
					return
				}
				panic(errsScript)
			}

			docSrc := storage.GetFileContent(creatorIDStr, creaIDStr, creaVersionStr, enum.Document)
			styleSrc := storage.GetFileContent(creatorIDStr, creaIDStr, creaVersionStr, enum.Style)
			err = script.IncludeHTMLCSS(docSrc, styleSrc)

			if err != nil {
				panic(err)
			}
		}

		storage.Source = lib.SrcPackages

		// localhost and empty domain (file:// for instance) are for user's tests
		bf, err := wb.SecureAndWrap([]string{fullPkg.Referer.String, "localhost", ""}...)

		if err != nil || storage.Error() != nil {
			c.Error(storage.Error()).SetMeta(ErrServ.SetParams("source", "package"))
			return
		}

		wbSrc, minErr := minifier.String("text/javascript", bf.String())

		if minErr != nil {
			c.Error(minErr).SetMeta(ErrServ.SetParams("source", "minifier"))
			return
		}

		fullPkg.NbBuild++

		// TODO if multitype allowed, package should have an engine too
		path := storage.StoreFile(
			wbSrc,
			"application/javascript",
			fmt.Sprintf("%d", user.(*model.User).ID),
			fmt.Sprintf("%d", fullPkg.ID.ValueDecoded),
			fmt.Sprintf("%d", fullPkg.NbBuild),
			enum.Wooble,
		)

		pkgPatchForm.Source = new(string)
		*pkgPatchForm.Source = lib.GetPkgURL() + "/" + path

		pkgPatchForm.NbBuild = new(uint64)
		*pkgPatchForm.NbBuild = fullPkg.NbBuild

		pkgPatchForm.BuiltAt = new(string)
		*pkgPatchForm.BuiltAt = time.Now().Format("2006-01-02 15:04:05 -07:00")
	}

	if err := model.UpdatePackagePatch(user.(*model.User).ID, pkg.ID, lib.SQLPatches(pkgPatchForm)); err != nil {
		c.Error(err).SetMeta(ErrDB)
		return
	}

	pkg.Source = lib.InitNullString(*pkgPatchForm.Source)
	pkg.BuiltAt = lib.InitNullTime(time.Now())

	c.JSON(OK, NewRes(pkg))
}

// DELETEPackage deletes a package
func DELETEPackage(c *gin.Context) {
	user, _ := c.Get("user")

	uID := user.(*model.User).ID
	pkgID := lib.InitID(c.Param("encid"))

	pkg, err := model.PackageByID(uID, pkgID)
	if err != nil {
		c.Error(err).SetMeta(ErrResNotFound.SetParams("source", "package", "id", pkgID.ValueEncoded))
		return
	}

	if err := model.DeletePackage(uID, pkgID); err != nil {
		c.Error(err).SetMeta(ErrDB)
		return
	}

	storage := lib.NewStorage(lib.SrcPackages)

	if pkg.NbBuild > 0 {
		for v := pkg.NbBuild; v > 0; v-- {
			storage.PushBulkFile(fmt.Sprintf("%d", uID), fmt.Sprintf("%d", pkg.ID.ValueDecoded), fmt.Sprintf("%d", v), enum.Wooble)
		}
		storage.BulkDeleteFiles()
	}

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
	pkg.Referer = lib.InitNullString(data.Referer)

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
	crea.ID = lib.InitID(c.Param("creaid"))
	crea.Alias = data.Alias
	crea.Version = data.Version
	pkgCrea.Creations = []model.Creation{*crea} // It'll update for only this creation

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
	creaVersionStr := fmt.Sprintf("%d", pkgCreaForm.Version)
	if err != nil {
		c.Error(err).SetMeta(ErrResNotFound.SetParams("source", "creation", "id", pkgCreaForm.CreationID+"/v"+creaVersionStr))
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

	if pkgCreaForm.Version == 0 || pkgCreaForm.Version > crea.Versions[len(crea.Versions)-1] {
		pkgCreaForm.Version = crea.Versions[len(crea.Versions)-1]
	}

	if err := model.NewPackageCreation(pkg.ID, crea.ID, pkgCreaForm.Version, crea.Alias); err != nil {
		c.Error(err).SetMeta(ErrDB)
		return
	}

	c.Header("Location", fmt.Sprintf("/packages/%s", pkg.ID.ValueEncoded))

	c.JSON(Created, NewRes(crea))
}

// RemovePackageCreation remove a creation from a package
func RemovePackageCreation(c *gin.Context) {
	user, _ := c.Get("user")

	if err := model.DeletePackageCreation(user.(*model.User).ID, lib.InitID(c.Param("encid")), lib.InitID(c.Param("creaid"))); err != nil {
		c.Error(err).SetMeta(ErrDB)
		return
	}

	c.AbortWithStatus(NoContent)
}
