package models

import (
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
)

type Script struct {
  Id bson.ObjectId `bson:"_id,omitempty"`
  Script string `bson:",omitempty" json:"script"`
}

func (this *Script) Create(s *mgo.Session) {
  this.Id = bson.NewObjectId()
}

func (this *Script) Save(s *mgo.Session) {
  if err := s.DB("").C(SCRIPT_C).Insert(this); err != nil {
    panic("Script failed to be saved")
  }
}

func (this *Script) Populate(s *mgo.Session) {
  s.DB("").C(SCRIPT_C).FindId(this.Id).One(&this)
}
