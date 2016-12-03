package router

import (
	"time"

	"wooble/router/handler"
	"wooble/router/middleware"

	"gopkg.in/gin-gonic/gin.v1"
)

func Load() {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	v1 := r.Group("/v1")
	{
		v1.GET("/heartbeat", func(c *gin.Context) {
			c.String(200, "%s", time.Now())
		})

		v1.POST("/signin", handler.SignIn)
		v1.POST("/signup", handler.SignUp)

		v1.GET("/creations/:title", handler.GETCreations) // TODO /creations/:username/:title
		v1.GET("/creations", handler.GETCreations)

		v1.Use(middleware.Authenticate())
		{
			v1.POST("/creations", handler.POSTCreations)
			v1.GET("/toto", handler.POSTCreations)

			v1.POST("/packages", handler.POSTPackages)

			// Resouce owner commands
			v1.POST("/packages/:id/push", handler.PushCreations)
			v1.GET("/packages/:id/build", handler.BuildPackage)
		}
	}

	r.Run()
}
