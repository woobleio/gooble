package controllers

import "gopkg.in/mgo.v2"

func GetScriptC(s *mgo.Session) *mgo.Collection {
  return s.DB("").C("Script")
}
