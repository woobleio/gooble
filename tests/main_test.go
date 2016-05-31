package tests

import (
  "testing"
  "github.com/gin-gonic/gin"
  describe "github.com/smartystreets/goconvey/convey"
  "net/http"
  "net/http/httptest"
)

func TestMain(t *testing.T) {
  router := gin.New()
  router.GET("/v1/heartbeat", func(c *Context) {})
  describe("Ask the server", t, func() {
    req, _ := http.NewRequest("GET", "/v1/heartbeat", nil)
    resp := httptest.NewRecorder()
    router.ServeHTTP(resp, req)
    So(res.Code, ShouldEqual, http.StatusBadRequest)
  })
}
