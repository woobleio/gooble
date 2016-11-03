package handler

import (
  "net/http"

  "wooble/model"
  "wooble/lib"

  "gopkg.in/gin-gonic/gin.v1"
)

func GETCreations(c *gin.Context) {
  var res interface{}
  var err error

  opts := lib.ParseOptions(c)
  title := c.Param("title")

  if title != "" {
    res, err = model.CreationByTitle(title, opts)
    if err != nil {
      c.JSON(http.StatusNotFound, NewError(NotFound, "Creation", title))
      return
    }
  } else {
    res, _ = model.AllCreations(opts)
  }

  c.JSON(http.StatusOK, res)
}
