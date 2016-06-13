package controllers

import (
  "wobblapp/lib"
  m "wobblapp/app/v1/models"
  "github.com/gin-gonic/gin"
  "gopkg.in/mgo.v2/bson"
)

func CreationPOST(c *gin.Context) {
  title := c.PostForm("title") // TODO should not be empty
  dom := c.PostForm("dom") // TODO should not be empty
  script := c.DefaultPostForm("script", "")
  style := c.DefaultPostForm("style", "")

  s := lib.GetSession()
  creationC := s.DB("").C("Creation")
  domC := s.DB("").C("DOM")
  scriptC := s.DB("").C("Script")
  styleC := s.DB("").C("Style")
  defer s.Close()

  domID := bson.NewObjectId()
  domC.Insert(&m.DOM{ID: domID, Dom: dom})

  scriptID := bson.NewObjectId()
  scriptC.Insert(&m.Script{ID: scriptID, Script: script})

  styleID := bson.NewObjectId()
  styleC.Insert(&m.Style{ID: styleID, Style: style})

  insertErr := creationC.Insert(&m.Creation{Title: title, DOM: domID, Style: styleID, Script: scriptID})
  if insertErr != nil {
    panic(insertErr)
  }
}

func CreationGET(c *gin.Context) {

}
