package controllers

import "gopkg.in/mgo.v2"

func GetStyleC(s *mgo.Session) *mgo.Collection {
  return s.DB("").C("Style")
}
