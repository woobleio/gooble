package controllers

import (
  "net/http"
  "github.com/gin-gonic/gin"
)

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
