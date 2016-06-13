package main

import (
	"os"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gopkg.in/mgo.v2"
	"time"
)

import (
	ctrl "wobblapp/app/v1/controllers"
	"wobblapp/lib"
)

func main() {
	InitConf()
	InitSession()
	StartGin()
}

func InitSession() {
	var host, dbName, port, username, passwd = viper.GetString("db_host"), viper.GetString("db_name"), viper.GetString("db_port"), viper.GetString("db_username"), viper.GetString("db_password")
	session, err := mgo.Dial("mongodb://" + username + ":" + passwd + "@" + host + ":" + port + "/" + dbName)
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
		v1.GET("/dom", ctrl.DomGET)
	}

	router.Run();
}
