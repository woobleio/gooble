package handler

import (
  "net/http"
  "fmt"

  "wooble/model"
  "wooble/lib"

  "github.com/gin-gonic/gin"
)

func GETAllCreations(c *gin.Context) {
  opts := lib.ParseOptions(c)

  creas, err := model.AllCreations(opts)
  if err != nil {
    fmt.Println(err)
  }

  c.JSON(http.StatusOK, creas)
}

func GETCreation(c *gin.Context) {
  opts := lib.ParseOptions(c)
  title := c.Param("title")

  crea, err := model.CreationByTitle(title, opts)
  if err != nil {
    fmt.Println(err)
  }

  c.JSON(http.StatusOK, crea)
}
