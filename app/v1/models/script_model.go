package models

import "gopkg.in/mgo.v2/bson"

type Script struct {
  ID bson.ObjectId `bson:"_id,omitempty"`
  Script string `bson:",omitempty" json:"script"`
}
