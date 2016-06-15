package controllers

import (
  "net/http"
  "wobblapp/lib"
  m "wobblapp/app/v1/models"
  "github.com/gin-gonic/gin"
)

type CreationForm struct {
  Title string `json:"title" binding:"required"`
  Dom string `json:"dom" binding:"required"`
  Style string `json:"style"`
  Script string `json:"script"`
}

func CreationPOST(c *gin.Context) {
  defer RequestErrorHandler(c)

  s := lib.GetSession()
  defer s.Close()

  dom := new(m.Dom)
  script := new(m.Script)
  style := new(m.Style)

  var form CreationForm
  if err := c.BindJSON(&form); err != nil {
    panic("Can't validate data, either dom or title is missing")
  }

  dom.Create(s)
  dom = &m.Dom{Id: dom.Id, Dom: form.Dom}
  dom.Save(s)

  if form.Script != "" {
    script.Create(s)
    script = &m.Script{Id: script.Id, Script: form.Script}
    script.Save(s)
  }

  if form.Style != "" {
    style.Create(s)
    style = &m.Style{Id: style.Id, Style: form.Style}
    style.Save(s)
  }

  creation := &m.Creation{
    Title: form.Title,
    Dom: dom.Id,
    Style: style.Id,
    Script: script.Id,
  }

  creation.Save(s)

  RequestSuccessHandler(c, "Creation '" + creation.Title + "' successfuly created")
}

func CreationGET(c *gin.Context) {
  defer RequestErrorHandler(c)

  s := lib.GetSession()
  defer s.Close()

  creation := new(m.Creation)
  dom := new(m.Dom)
  script := new(m.Script)
  style := new(m.Style)

  title := c.Param("title")

  // TODO better error handling (think to validate instead of panicing error when not found)
  creation.FindOneWithKey(s, title)
  dom.FindOne(s, creation.Dom)
  script.FindOne(s, creation.Script)
  style.FindOne(s, creation.Style)

  json := &CreationForm{
    creation.Title,
    dom.Dom,
    script.Script,
    style.Style,
  }

  c.JSON(http.StatusOK, json)
}
