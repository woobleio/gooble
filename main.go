package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

import (
	"wobblapp/config"
	ctrl "wobblapp/app/v1/controllers"
)

func main() {
	StartGin()
}

func StartGin() {
	router := gin.Default()

	// http://stackoverflow.com/questions/36420863/file-path-in-golang
	fmt.Printf("%s", config.DBCredentials())

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
