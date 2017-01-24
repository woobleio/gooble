package router

import (
	"time"

	"wooble/lib"
	"wooble/router/handler"
	"wooble/router/middleware"

	cors "gopkg.in/gin-contrib/cors.v1"
	"gopkg.in/gin-gonic/gin.v1"
)

// Load initializes the router and loads all handlers
func Load() {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

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

		v1.POST("/token/generate", handler.GenerateToken)
		v1.POST("/token/refresh", handler.RefreshToken)

		v1.GET("/creations", handler.GETCreations)
		v1.GET("/creations/:id", handler.GETCreations)

		v1.Use(middleware.Authenticate())
		{
			v1.POST("/creations", handler.POSTCreations)

			// packages is private, so those requests are about the authenticated user only
			v1.GET("/packages", handler.GETPackages)
			v1.GET("/packages/:id", handler.GETPackages)
			v1.POST("/packages", handler.POSTPackages)
			v1.POST("/packages/:id/push", handler.PushCreations)
			v1.PUT("/packages/:id/build", handler.BuildPackage)
		}
	}

	r.Run()
}
