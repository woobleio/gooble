package controllers

import (
  "net/http"
  "wobblapp/lib"
  m "wobblapp/app/v1/models"
  "github.com/gin-gonic/gin"
  "gopkg.in/mgo.v2"
)

const CREA_C = "Creation"

type CreationCtrl struct {
  Form *m.Creation
}

func (ctrl *CreationCtrl) Save(s *mgo.Session) {
  err := s.DB("").C(CREA_C).Insert(ctrl.Form)
  if err != nil {
    panic(err)
  }
}

func (ctrl *CreationCtrl) FindOneWithKey(s *mgo.Session, k string) {
  ctrl.Form = &m.Creation{}
  s.DB("").C(CREA_C).FindId(k).One(&ctrl.Form)
}

func CreationPOST(c *gin.Context) {
  title := c.PostForm("title") // TODO field verification & error
  dom := c.PostForm("dom") // TODO should not be empty
  script := c.DefaultPostForm("script", "")
  style := c.DefaultPostForm("style", "")

  s := lib.GetSession()
  defer s.Close()

  this := new(CreationCtrl)
  domCtrl := new(DomCtrl)
  scriptCtrl := new(ScriptCtrl)
  styleCtrl := new(StyleCtrl)

  // Create obj ; Populate ; Push
  domCtrl.Create(s)
  domCtrl.Form = &m.DOM{ID: domCtrl.Id, Dom: dom}
  domCtrl.Save(s)

  scriptCtrl.Create(s)
  scriptCtrl.Form = &m.Script{ID: scriptCtrl.Id, Script: script}
  scriptCtrl.Save(s)

  styleCtrl.Create(s)
  styleCtrl.Form = &m.Style{ID: styleCtrl.Id, Style: style}
  styleCtrl.Save(s)

  this.Form = &m.Creation{
    Title: title,
    Dom: domCtrl.Id,
    Style: styleCtrl.Id,
    Script: scriptCtrl.Id,
  }

  this.Save(s)
}

func CreationGET(c *gin.Context) {
  // TODO populate function in lib & refactor Collections
  s := lib.GetSession()
  defer s.Close()

  title := c.Query("title")

  this := new(CreationCtrl)
  this.FindOneWithKey(s, title)

  domCtrl := new(DomCtrl)
  domCtrl.FindOne(s, this.Form.Dom)

  scriptCtrl := new(ScriptCtrl)
  scriptCtrl.FindOne(s, this.Form.Script)

  styleCtrl := new(StyleCtrl)
  styleCtrl.FindOne(s, this.Form.Style)

  c.JSON(http.StatusOK, gin.H{
    "title": this.Form.Title,
    "dom": domCtrl.Form.Dom,
    "script": scriptCtrl.Form.Script,
    "style": styleCtrl.Form.Style,
  })
}
