package controllers

import (
  m "wobblapp/app/v1/models"
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
)

const STYLE_C = "Style"

type StyleCtrl struct {
  Id bson.ObjectId
  Form *m.Style
}

func (ctrl *StyleCtrl) Create(s *mgo.Session) {
  ctrl.Id = bson.NewObjectId()
}

func (ctrl *StyleCtrl) Save(s *mgo.Session) {
  err := s.DB("").C(STYLE_C).Insert(ctrl.Form)
  if err != nil {
    panic(err)
  }
}

func (ctrl *StyleCtrl) FindOne(s *mgo.Session, o bson.ObjectId) {
  ctrl.Form = &m.Style{}
  s.DB("").C(STYLE_C).FindId(o).One(&ctrl.Form)
}
