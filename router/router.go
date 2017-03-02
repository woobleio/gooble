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

		v1.POST("/tokens", handler.GenerateToken)
		v1.PUT("/tokens", handler.RefreshToken)

		v1.GET("/creations", handler.GETCreations)
		v1.GET("/creations/:encid", handler.GETCreations)

		v1.Use(middleware.Authenticate())
		{
			creations := v1.Group("/creations")
			{
				creations.POST("", handler.POSTCreation)
				creations.PUT("/:encid", handler.PUTCreation)
				creations.DELETE("/:encid", handler.DELETECreation)
				creations.PATCH("/:encid/publish", handler.PublishCreation)
				creations.GET("/:encid/code", handler.GETCreationCode)

				creations.POST("/:encid/versions", handler.POSTCreationVersion)
				creations.PUT("/:encid/versions", handler.SaveVersion)
			}

			v1.POST("/buy", handler.BuyCreations)

			users := v1.Group("/users")
			{
				users.PATCH("/password", handler.UpdatePassword)
				users.DELETE("", handler.DELETEUser)
				users.POST("/funds/bank", handler.POSTUserBank)      // FIXME Stripe version, managed account don't work for now
				users.POST("/funds/withdraw", handler.WithdrawFunds) // FIXME Stripe version, managed account don't work for now
			}

			packages := v1.Group("/packages")
			{
				packages.GET("", handler.GETPackages)
				packages.POST("", handler.POSTPackage)
				packages.GET("/:encid", handler.GETPackages)
				packages.PUT("/:encid", handler.PUTPackage)
				packages.DELETE("/:encid", handler.DELETEPackage)

				packages.POST("/:encid/creations", handler.PushCreation)
				packages.DELETE("/:encid/creations", handler.RemovePackageCreation)
				packages.PUT("/:encid/creations/:creaid", handler.PUTPackageCreation)
				packages.POST("/:encid/build", handler.BuildPackage)
			}
		}
	}

	r.Run()
}
