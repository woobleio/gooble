package controllers

import (
  m "wobblapp/app/v1/models"
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
)

const SCRIPT_C = "scripts"

type ScriptCtrl struct {
  Id bson.ObjectId
  Model *m.Script
}

func (this *ScriptCtrl) Create(s *mgo.Session) {
  this.Id = bson.NewObjectId()
}

func (this *ScriptCtrl) Save(s *mgo.Session) {
  err := s.DB("").C(SCRIPT_C).Insert(this.Model)
  if err != nil {
    panic("Script failed to be saved")
  }
}

func (this *ScriptCtrl) FindOne(s *mgo.Session, o bson.ObjectId) {
  this.Model = &m.Script{}
  s.DB("").C(SCRIPT_C).FindId(o).One(&this.Model)
}
