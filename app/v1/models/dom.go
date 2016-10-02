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
  if err := s.DB("").C(DOM_C).Insert(this); err != nil {
    panic("HTML document failed to be saved")
  }
}

func (this *Dom) Populate(s *mgo.Session) {
  s.DB("").C(DOM_C).FindId(this.Id).One(&this)
}
