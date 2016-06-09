package controllers

import (
  "github.com/gin-gonic/gin"
)

func DomGET(c *gin.Context) {
  c.JSON(200, gin.H{ "message": "DOM" })
}
