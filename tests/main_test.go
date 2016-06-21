package tests

import (
  "os"
  "github.com/spf13/viper"
  "testing"
  "github.com/gin-gonic/gin"
  . "github.com/smartystreets/goconvey/convey"
  "gopkg.in/mgo.v2"
  "net/http"
  "net/http/httptest"
)

func init() {
  viper.SetConfigName("test")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(os.Getenv("CONFPATH"))

	errViper := viper.ReadInConfig()
	if errViper != nil {
		panic(errViper)
	}
}

func GetSession() *mgo.Session {
  var host, dbName, port, username, passwd = viper.GetString("db_host"), viper.GetString("db_name"), viper.GetString("db_port"), viper.GetString("db_username"), viper.GetString("db_password")
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

  return session
}

func TestMain(t *testing.T) {
  router := gin.New()
  router.GET("/v1/heartbeat", func(c *gin.Context) {})
  Convey("Ask the server", t, func() {
    req, _ := http.NewRequest("GET", "/v1/heartbeat", nil)
    res := httptest.NewRecorder()
    router.ServeHTTP(res, req)
    So(res.Code, ShouldEqual, http.StatusOK)
  })
}
