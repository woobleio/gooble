package models

import (
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
)

type Style struct {
  Id bson.ObjectId `bson:"_id"`
  Style string
}

func (this *Style) Create(s *mgo.Session) {
  this.Id = bson.NewObjectId()
}

func (this *Style) Save(s *mgo.Session) {
  err := s.DB("").C(STYLE_C).Insert(this)
  if err != nil {
    panic("Stylesheet failed to be save")
  }
}

func (this *Style) FindOne(s *mgo.Session, o bson.ObjectId) {
  s.DB("").C(STYLE_C).FindId(o).One(&this)
}
