package router

import (
	"time"

	"github.com/gin-gonic/gin/binding"

	"wooble/lib"
	"wooble/router/handlers"
	middleware "wooble/router/middlewares"

	cors "gopkg.in/gin-contrib/cors.v1"
	"gopkg.in/gin-gonic/gin.v1"
	validator "gopkg.in/go-playground/validator.v9"
)

type Validator struct {
	*validator.Validate
}

func (v *Validator) ValidateStruct(i interface{}) error {
	return v.Struct(i)
}

// Load initializes the router and loads all handlers
func Load() {
	r := gin.New()

	binding.Validator = &Validator{validator.New()}

	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(handler.HandleErrors)

	config := cors.DefaultConfig()
	config.AllowOrigins = lib.GetOrigins()
	config.AllowHeaders = []string{"Content-Type", "Authorization"}
	config.ExposeHeaders = []string{"Authorization", "Location"}

	r.Use(cors.New(config))

	v1 := r.Group("/v1")
	{
		v1.GET("/heartbeat", func(c *gin.Context) {
			c.String(200, "%s", time.Now())
		})

		v1.POST("/users", handler.POSTUser)
		v1.GET("/users/:username", handler.GETUser)

		v1.POST("/token/generate", handler.GenerateToken)
		v1.POST("/token/refresh", handler.RefreshToken)

		v1.GET("/creations", handler.GETCreations)
		v1.GET("/creations/:encid", handler.GETCreations)

		v1.Use(middleware.Authenticate())
		{
			v1.POST("/creations", handler.POSTCreation)
			v1.PUT("/creations/:encid", handler.PUTCreation)
			// v1.DELETE("/creation/:encid", handler.DELETECreation)
			v1.GET("/creations/:encid/code", handler.GETCodeCreation)
			v1.PATCH("/creations/:encid/publish", handler.PublishCreation)
			v1.POST("/creations/:encid/versions", handler.POSTCreationVersion)
			v1.PUT("/creations/:encid/versions/:version", handler.SaveVersion)

			v1.POST("/buy", handler.BuyCreations)

			v1.POST("/users/password", handler.UpdatePassword)
			v1.DELETE("/users", handler.DELETEUser)
			v1.POST("/users/funds/bank", handler.POSTUserBank)      // FIXME Stripe version, don't work for now
			v1.POST("/users/funds/withdraw", handler.WithdrawFunds) // FIXME Stripe version, don't work for now

			// packages is private, so those requests are about the authenticated user only
			v1.GET("/packages", handler.GETPackages)
			v1.GET("/packages/:encid", handler.GETPackages)
			v1.PUT("/packages/:encid", handler.PUTPackage)
			v1.DELETE("/packages/:encid", handler.DELETEPackage)
			v1.POST("/packages", handler.POSTPackage)
			v1.POST("/packages/:encid/creations", handler.PushCreation)
			v1.DELETE("/packages/:encid/creations", handler.RemovePackageCreation)
			v1.PUT("/packages/:encid/creations/:creaid", handler.PUTPackageCreation)
			v1.PATCH("/packages/:encid/build", handler.BuildPackage)
		}
	}

	r.Run()
}
