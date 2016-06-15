package models

import "gopkg.in/mgo.v2/bson"

type Creation struct {
  Title string `bson:"_id"`
  Dom bson.ObjectId `json:"dom"`
  Style bson.ObjectId `bson:",omitempty" json:"style"`
  Script bson.ObjectId `bson:",omitempty" json:"script"`
}
