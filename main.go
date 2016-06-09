package main

import (
	"github.com/gin-gonic/gin"
	"time"
)

import ctrl "wobblapp/app/controllers"

func main() {
	router := gin.Default()
	// MIDDLEWARE

	v1 := router.Group("/v1")
	{
		v1.GET("/heartbeat", func(c *gin.Context) {
			c.String(200, "%s", time.Now())
		})
		v1.GET("/dom", ctrl.DomGET)
	}

	router.Run();
}
