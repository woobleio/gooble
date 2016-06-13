package controllers

import (
  m "wobblapp/app/v1/models"
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
)

const DOM_C = "DOM"

type DomCtrl struct {
  Id bson.ObjectId
  Form *m.DOM
}

func (ctrl *DomCtrl) Create(s *mgo.Session) {
  ctrl.Id = bson.NewObjectId()
}

func (ctrl *DomCtrl) Save(s *mgo.Session) {
  err := s.DB("").C(DOM_C).Insert(ctrl.Form)
  if err != nil {
    panic(err)
  }
}

func (ctrl *DomCtrl) FindOne(s *mgo.Session, o bson.ObjectId) {
  ctrl.Form = &m.DOM{}
  s.DB("").C(DOM_C).FindId(o).One(&ctrl.Form)
}
