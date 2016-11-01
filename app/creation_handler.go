package app

import (
  "net/http"
  "wooblapp/app/model"
  "wooblapp/app/util"
  "fmt"

  "github.com/gin-gonic/gin"
)

func GETCreation(c *gin.Context) {
  opt := util.ParseOptions(c)
  title := c.Param("title")

  crea, err := model.CreationByTitle(title, opt)
  if err != nil {
    fmt.Println(err)
  }

  c.JSON(http.StatusOK, crea)
}
