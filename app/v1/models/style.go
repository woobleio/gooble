package models

import "gopkg.in/mgo.v2/bson"

type Style struct {
  ID bson.ObjectId `bson:"_id"`
  Style string
}
