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
}

/**
 * Handle all controllers errors
 */
func RequestErrorHandler(c *gin.Context) {
  if err := recover(); err != nil {
    // TODO should get the message
    c.JSON(http.StatusBadRequest, gin.H{
      "error": err,
    })
  }
}
