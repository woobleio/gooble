package controllers

import "gopkg.in/mgo.v2"

func GetDomC(s *mgo.Session) *mgo.Collection {
  return s.DB("").C("DOM")
}
