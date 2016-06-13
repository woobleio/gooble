package models

import "gopkg.in/mgo.v2/bson"

type Creation struct {
  Title string `bson:"_id"`
  Dom bson.ObjectId
  Style bson.ObjectId
  Script bson.ObjectId
}
