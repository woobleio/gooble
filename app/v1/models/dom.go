package models

import (
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
)

// IN CASE NODE GEN IS NECESSARY
/*type DOM struct {
  Node bson.ObjectId
  Children []DOM
}*/
type Dom struct {
  Id bson.ObjectId `bson:"_id"`
  Dom string
}

func (this *Dom) Create(s *mgo.Session) {
  this.Id = bson.NewObjectId()
}

func (this *Dom) Save(s *mgo.Session) {
  err := s.DB("").C(DOM_C).Insert(this)
  if err != nil {
    panic("HTML document failed to be saved")
  }
}

func (this *Dom) FindOne(s *mgo.Session, o bson.ObjectId) {
  err := s.DB("").C(DOM_C).FindId(o).One(&this)
  if err != nil {
    panic("HTML document not found")
  }
}
