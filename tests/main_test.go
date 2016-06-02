package tests

import (
  "testing"
  "github.com/gin-gonic/gin"
  . "github.com/smartystreets/goconvey/convey"
  "net/http"
  "net/http/httptest"
)

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
