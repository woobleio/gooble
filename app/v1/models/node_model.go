package models

// IN CASE NODE GEN IS NECESSARY

import "gopkg.in/mgo.v2/bson"

type Node struct {
  ID bson.ObjectId `bson:"_id"`
  ElId string
  Classes string
  El string
}
