package models

import "gopkg.in/mgo.v2/bson"

// IN CASE NODE GEN IS NECESSARY
/*type DOM struct {
  Node bson.ObjectId
  Children []DOM
}*/
type DOM struct {
  ID bson.ObjectId `bson:"_id"`
  Dom string
}
