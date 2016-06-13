package controllers

import (
  "net/http"
  "wobblapp/lib"
  m "wobblapp/app/v1/models"
  "github.com/gin-gonic/gin"
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
)

func GetCreationC(s *mgo.Session) *mgo.Collection {
  return s.DB("").C("Creation")
}

func CreationPOST(c *gin.Context) {
  title := c.PostForm("title") // TODO should not be empty
  dom := c.PostForm("dom") // TODO should not be empty
  script := c.DefaultPostForm("script", "")
  style := c.DefaultPostForm("style", "")

  s := lib.GetSession()
  creationC := GetCreationC(s)
  domC := GetDomC(s)
  scriptC := GetScriptC(s)
  styleC := GetStyleC(s)
  defer s.Close()

  domID := bson.NewObjectId()
  domC.Insert(&m.DOM{ID: domID, Dom: dom})

  scriptID := bson.NewObjectId()
  scriptC.Insert(&m.Script{ID: scriptID, Script: script})

  styleID := bson.NewObjectId()
  styleC.Insert(&m.Style{ID: styleID, Style: style})

  insertErr := creationC.Insert(&m.Creation{Title: title, Dom: domID, Style: styleID, Script: scriptID})
  if insertErr != nil {
    panic(insertErr)
  }
}

func CreationGET(c *gin.Context) {
  // TODO populate function in lib & refactor Collections
  s := lib.GetSession()
  creationC := GetCreationC(s)
  domC := GetDomC(s)
  scriptC := GetScriptC(s)
  styleC := GetStyleC(s)
  defer s.Close()

  title := c.Query("title")

  creaResult := m.Creation{}
  err := creationC.FindId(title).One(&creaResult)
  if err != nil {
    panic(err)
  }
  domRes := m.DOM{}
  scriptRes := m.Script{}
  styleRes := m.Style{}
  domC.FindId(creaResult.Dom).One(&domRes)
  scriptC.FindId(creaResult.Script).One(&scriptRes)
  styleC.FindId(creaResult.Style).One(&styleRes)
  c.JSON(http.StatusOK, gin.H{
    "title": creaResult.Title,
    "dom": domRes.Dom,
    "script": scriptRes.Script,
    "style": styleRes.Style,
  })
}
