package controllers

import (
  "wobblapp/lib"
  m "wobblapp/app/v1/models"
  "github.com/gin-gonic/gin"
)

func DomGET(c *gin.Context) {
  s := lib.GetSession()
  doc := s.DB("").C("test")
  defer s.Close()

  doc.Insert(&m.DOM{ Message: "test" })

  c.JSON(200, gin.H{ "message": "DOM" })
}
