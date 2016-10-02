package main

import (
	"os"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gopkg.in/mgo.v2"
	"time"
)

import (
	ctrl "wooblapp/app/v1/controllers"
	"wooblapp/lib"
)

func main() {
	InitConf()
	InitSession()
	StartGin()
}

func InitSession() {
	var host, dbName, port, username, passwd =
		viper.GetString("db_host"),
		viper.GetString("db_name"),
		viper.GetString("db_port"),
		viper.GetString("db_username"),
		viper.GetString("db_password")

	var dbUrl = "mongodb://"

	switch {
		case username != "":
			dbUrl += username
			dbUrl += ":" + passwd + "@"
		case host != "":
			dbUrl += host
		case port != "":
			dbUrl += ":" + port
		case dbName != "":
			dbUrl += "/" + dbName
			break
	}

	session, err := mgo.Dial(dbUrl)
	if err != nil {
    panic(err)
  }

	lib.SetSession(session)
}

func InitConf() {
	viper.SetConfigName(os.Getenv("GOENV"))
	viper.SetConfigType("yaml")
	viper.AddConfigPath(os.Getenv("CONFPATH"))

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
		v1.GET("/creation/:title", ctrl.CreationGET)
		v1.POST("/creation", ctrl.CreationPOST)
	}

	router.Run();
}
