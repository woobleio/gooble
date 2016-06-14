package controllers

import (
  m "wobblapp/app/v1/models"
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
)

const SCRIPT_C = "Script"

type ScriptCtrl struct {
  Id bson.ObjectId
  Form *m.Script
}

func (ctrl *ScriptCtrl) Create(s *mgo.Session) {
  ctrl.Id = bson.NewObjectId()
}

func (ctrl *ScriptCtrl) Save(s *mgo.Session) {
  err := s.DB("").C(SCRIPT_C).Insert(ctrl.Form)
  if err != nil {
    panic("Script failed to be saved")
  }
}

func (ctrl *ScriptCtrl) FindOne(s *mgo.Session, o bson.ObjectId) {
  ctrl.Form = &m.Script{}
  err := s.DB("").C(SCRIPT_C).FindId(o).One(&ctrl.Form)
  if err != nil {
    panic("Script not found")
  }
}
