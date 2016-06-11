package main

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"time"
)

import (
	ctrl "wobblapp/app/v1/controllers"
)

func main() {
	InitConf()
	StartGin()
}

func InitConf() {
	// TODO one conf per env
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/go/configs")

	errViper := viper.ReadInConfig()

	if errViper != nil {
		panic(errViper)
	}
}

func StartGin() {
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
