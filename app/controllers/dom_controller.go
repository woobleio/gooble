package controllers

import (
  "github.com/gin-gonic/gin"
  "wobblapp/app/models"
)

func DomGET(c *gin.Context) {
  var json models.DOM
  if c.BindJSON(&json) == nil {
    c.JSON(200, gin.H{ "message": "DOM" })
  } else {
    c.JSON(200, gin.H{ "message": "error" })
  }
}
