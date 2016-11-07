package router

import (
	"time"

	"wooble/router/handler"

	"gopkg.in/gin-gonic/gin.v1"
)

func Load() {
	r := gin.Default()

	// MIDDLEWARE
	v1 := r.Group("/v1")
	{
		v1.GET("/heartbeat", func(c *gin.Context) {
			c.String(200, "%s", time.Now())
		})
		v1.GET("/creations/:title", handler.GETCreations)
		v1.GET("/creations", handler.GETCreations)
	}

	r.Run()
}
