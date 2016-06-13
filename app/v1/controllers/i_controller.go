package controllers

import (
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
)

type Controller interface {
  Create(*mgo.Session)
  Save(*mgo.Session)
  FindOne(*mgo.Session, bson.ObjectId)
  FindOneWithKey(*mgo.Session, string)
}
