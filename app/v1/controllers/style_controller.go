package controllers

import (
  m "wobblapp/app/v1/models"
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
)

const STYLE_C = "styles"

type StyleCtrl struct {
  Id bson.ObjectId
  Model *m.Style
}

func (this *StyleCtrl) Create(s *mgo.Session) {
  this.Id = bson.NewObjectId()
}

func (this *StyleCtrl) Save(s *mgo.Session) {
  err := s.DB("").C(STYLE_C).Insert(this.Model)
  if err != nil {
    panic("Stylesheet failed to be save")
  }
}

func (this *StyleCtrl) FindOne(s *mgo.Session, o bson.ObjectId) {
  this.Model = &m.Style{}
  s.DB("").C(STYLE_C).FindId(o).One(&this.Model)
}
