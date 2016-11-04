package handler

import (
  "wooble/model"
  "wooble/lib"

  "gopkg.in/gin-gonic/gin.v1"
)

func GETCreations(c *gin.Context) {
  var data interface{}
  var err error

  res := NewRes()

  opts := lib.ParseOptions(c)
  title := c.Param("title")

  if title != "" {
    data, err = model.CreationByTitle(title, opts)
    if err != nil {
      res.Error(ResNotFound, "Creation", title)
    }
  } else {
    data, err = model.AllCreations(opts)
    if err != nil {
      res.Error(NotFound, "creations")
    }
  }

  res.Response(data)

  c.JSON(res.HttpStatus(), res)
}
