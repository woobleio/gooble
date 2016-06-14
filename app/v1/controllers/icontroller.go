package controllers

import (
  "net/http"
  "github.com/gin-gonic/gin"
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
)

type Controller interface {
  Create(*mgo.Session)
  Save(*mgo.Session)
  FindOne(*mgo.Session, bson.ObjectId)
  FindOneWithKey(*mgo.Session, string)
  ValidateAndSet(*gin.Context)
}

/**
 * Handle all requests errors
 */
func RequestErrorHandler(c *gin.Context) {
  if err := recover(); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{
      "message": err,
    })
  }
}

func RequestSuccessHandler(c *gin.Context, message string) {
  c.JSON(http.StatusOK, gin.H{
    "message": message,
  })
}
