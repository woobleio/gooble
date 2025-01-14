package router

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	cors "gopkg.in/gin-contrib/cors.v1"

	"wooble/lib"
	"wooble/router/handlers"
	middleware "wooble/router/middlewares"

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
	config.AllowMethods = []string{"POST", "GET", "PUT", "PATCH", "DELETE"}
	config.AllowHeaders = []string{"Content-Type", "Authorization"}
	config.ExposeHeaders = []string{"Authorization", "Location"}

	r.Use(cors.New(config))

	v1 := r.Group("/v1")
	{
		v1.GET("/heartbeat", func(c *gin.Context) {
			c.String(200, "%s", time.Now())
		})

		v1.POST("/webhooks/:event", handler.POSTWebhooks)

		v1.POST("/users", handler.POSTUser)
		v1.GET("/users/:name", handler.GETUser)

		v1.POST("/tokens", handler.GenerateToken)
		v1.PUT("/tokens", handler.RefreshToken)

		v1.GET("/plans", handler.GETPlans)

		v1.GET("/tags", handler.GETTags)

		v1.GET("/creations", handler.GETCreations)
		v1.GET("/creations/:encid", handler.GETCreations)
		v1.GET("/creations/:encid/code", handler.GETCreationCode)

		v1.Use(middleware.Authenticate())
		{
			users := v1.Group("/users")
			{
				users.PATCH("", handler.PATCHUser)
				users.DELETE("", handler.DELETEUser)
			}

			v1.POST("/files", handler.POSTFile)

			v1.Use(middleware.IsActive())
			{
				creations := v1.Group("/creations")
				{
					creations.POST("", handler.POSTCreation)
					creations.PUT("/:encid", handler.PUTCreation)
					creations.PATCH("/:encid", handler.PATCHCreation)
					creations.DELETE("/:encid", handler.DELETECreation)

					creations.POST("/:encid/versions", handler.POSTCreationVersion)
					creations.PUT("/:encid/versions", handler.SaveVersion)
				}

				packages := v1.Group("/packages")
				{
					packages.GET("", handler.GETPackages)
					packages.POST("", handler.POSTPackage)
					packages.GET("/:encid", handler.GETPackages)
					packages.PUT("/:encid", handler.PUTPackage)
					packages.PATCH("/:encid", handler.PATCHPackage)
					packages.DELETE("/:encid", handler.DELETEPackage)

					packages.POST("/:encid/creations", handler.PushCreation)
					packages.DELETE("/:encid/creations/:creaid", handler.RemovePackageCreation)
					packages.PUT("/:encid/creations/:creaid", handler.PUTPackageCreation)
				}
			}
		}
	}

	r.Run()
}
