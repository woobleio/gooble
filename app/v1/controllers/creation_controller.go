package controllers

import (
  "net/http"
  "wobblapp/lib"
  m "wobblapp/app/v1/models"
  "github.com/gin-gonic/gin"
  "gopkg.in/mgo.v2"
  // "gopkg.in/mgo.v2/bson"
)

const CREA_C = "creations"

type CreationForm struct {
  Title string `json:"title" binding:"required"`
  Dom string `json:"dom" binding:"required"`
  Style string `json:"style"`
  Script string `json:"script"`
}

type CreationCtrl struct {
  Form *CreationForm
  Model *m.Creation
}

func (this *CreationCtrl) Save(s *mgo.Session) {
  if err := s.DB("").C(CREA_C).Insert(&this.Model); err != nil {
    // TODO be more specific for the error
    panic("Creation '" + this.Model.Title + "' failed to be saved")
  }
}

func (this *CreationCtrl) FindOneWithKey(s *mgo.Session, k string) {
  this.Model = &m.Creation{}
  err := s.DB("").C(CREA_C).FindId(k).One(&this.Model)
  if err != nil {
    panic("Creation '" + k + "' not found")
  }
}

func (this *CreationCtrl) ValidateAndSet(c *gin.Context) {
  var form CreationForm
  if err := c.BindJSON(&form); err != nil {
    panic("Can't validate data, either dom or title is missing")
  }
  this.Form = &form
}

func CreationPOST(c *gin.Context) {
  defer RequestErrorHandler(c)

  s := lib.GetSession()
  defer s.Close()

  this := new(CreationCtrl)
  domCtrl := new(DomCtrl)
  scriptCtrl := new(ScriptCtrl)
  styleCtrl := new(StyleCtrl)

  this.ValidateAndSet(c)
  form := this.Form

  // Create obj ; Populate ; Push TODO better populate
  domCtrl.Create(s)
  domCtrl.Model = &m.DOM{ID: domCtrl.Id, Dom: form.Dom}
  domCtrl.Save(s)

  if form.Script != "" {
    scriptCtrl.Create(s)
    scriptCtrl.Model = &m.Script{ID: scriptCtrl.Id, Script: form.Script}
    scriptCtrl.Save(s)
  }

  if form.Style != "" {
    styleCtrl.Create(s)
    styleCtrl.Model = &m.Style{ID: styleCtrl.Id, Style: form.Style}
    styleCtrl.Save(s)
  }

  this.Model = &m.Creation{
    Title: form.Title,
    Dom: domCtrl.Id,
    Style: styleCtrl.Id,
    Script: scriptCtrl.Id,
  }

  this.Save(s)

  RequestSuccessHandler(c, "Creation '" + this.Model.Title + "' successfuly created")
}

func CreationGET(c *gin.Context) {
  defer RequestErrorHandler(c)

  s := lib.GetSession()
  defer s.Close()

  this := new(CreationCtrl)
  domCtrl := new(DomCtrl)
  scriptCtrl := new(ScriptCtrl)
  styleCtrl := new(StyleCtrl)

  title := c.Query("title")

  // TODO better error handling (think to validate instead of panicing error when not found)
  this.FindOneWithKey(s, title)
  domCtrl.FindOne(s, this.Model.Dom)
  scriptCtrl.FindOne(s, this.Model.Script)
  styleCtrl.FindOne(s, this.Model.Style)

  json := &CreationForm{
    this.Model.Title,
    domCtrl.Model.Dom,
    scriptCtrl.Model.Script,
    styleCtrl.Model.Style,
  }

  c.JSON(http.StatusOK, json)
}
