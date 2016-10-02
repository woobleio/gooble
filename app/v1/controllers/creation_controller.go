package controllers

import (
  "net/http"
  "wooblapp/lib"
  m "wooblapp/app/v1/models"
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

  var form CreationForm
  if err := c.BindJSON(&form); err != nil {
    panic("Can't validate data, either dom or title is missing")
  }

  var dom m.Model = &m.Dom{Dom: form.Dom}
  var script m.Model = &m.Script{Script: form.Script}
  var style m.Model = &m.Style{Style: form.Style}

  dom.Create(s)
  dom.Save(s)

  if form.Script != "" {
    script.Create(s)
    script.Save(s)
  }

  if form.Style != "" {
    style.Create(s)
    style.Save(s)
  }

  var creation m.Model = &m.Creation{
    Title: form.Title,
    Dom: dom.(*m.Dom).Id,
    Style: style.(*m.Style).Id,
    Script: script.(*m.Script).Id,
  }

  creation.Save(s)

  RequestSuccessHandler(c, "Creation '" + form.Title + "' successfuly created")
}

func CreationGET(c *gin.Context) {
  defer RequestErrorHandler(c)

  s := lib.GetSession()
  defer s.Close()

  title := c.Param("title")

  var creation m.Model = &m.Creation{Title: title}
  creation.Populate(s)

  var dom m.Model = &m.Dom{Id: creation.(*m.Creation).Dom}
  var script m.Model = &m.Script{Id: creation.(*m.Creation).Script}
  var style m.Model = &m.Style{Id: creation.(*m.Creation).Style}

  dom.Populate(s)
  script.Populate(s)
  style.Populate(s)

  json := &CreationForm{
    creation.(*m.Creation).Title,
    dom.(*m.Dom).Dom,
    script.(*m.Script).Script,
    style.(*m.Style).Style,
  }

  if json.Title == "" || json.Dom == "" {
    panic("Creation " + creation.(*m.Creation).Title + " not found or dom is empty")
  }

  c.JSON(http.StatusOK, json)
}
