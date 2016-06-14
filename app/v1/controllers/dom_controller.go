package controllers

import (
  m "wobblapp/app/v1/models"
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
)

const DOM_C = "documents"

type DomCtrl struct {
  Id bson.ObjectId
  Model *m.DOM
}

func (this *DomCtrl) Create(s *mgo.Session) {
  this.Id = bson.NewObjectId()
}

func (this *DomCtrl) Save(s *mgo.Session) {
  err := s.DB("").C(DOM_C).Insert(this.Model)
  if err != nil {
    panic("HTML document failed to be saved")
  }
}

func (this *DomCtrl) FindOne(s *mgo.Session, o bson.ObjectId) {
  this.Model = &m.DOM{}
  err := s.DB("").C(DOM_C).FindId(o).One(&this.Model)
  if err != nil {
    panic("HTML document not found")
  }
}
