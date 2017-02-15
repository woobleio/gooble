package router

import (
	"time"

	"wooble/lib"
	"wooble/router/handlers"

	cors "gopkg.in/gin-contrib/cors.v1"
	"gopkg.in/gin-gonic/gin.v1"
)

// Load initializes the router and loads all handlers
func Load() {
	r := gin.New()
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

		v1.POST("/users", handler.SignUp)

		// v1.POST("/token/generate", handler.Handle(handler.GenerateToken))
		// v1.POST("/token/refresh", handler.Handle(handler.RefreshToken))
		//
		// v1.GET("/creations", handler.Handle(handler.GETCreations))
		// v1.GET("/creations/:id", handler.Handle(handler.GETCreations))
		//
		// v1.Use(middleware.Authenticate())
		// {
		// 	v1.POST("/creations", handler.Handle(handler.POSTCreations))
		// 	v1.PUT("/creations/:id/buy", handler.Handle(handler.BuyCreation))
		//
		// 	// packages is private, so those requests are about the authenticated user only
		// 	v1.GET("/packages", handler.Handle(handler.GETPackages))
		// 	v1.GET("/packages/:id", handler.Handle(handler.GETPackages))
		// 	v1.POST("/packages", handler.Handle(handler.POSTPackages))
		// 	v1.POST("/packages/:id/push", handler.Handle(handler.PushCreation))
		// 	v1.PATCH("/packages/:id/build", handler.Handle(handler.BuildPackage))
		// }
	}

	r.Run()
}
