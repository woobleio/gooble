package models

import (
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
)

type Creation struct {
  Title string `bson:"_id"`
  Dom bson.ObjectId `json:"dom"`
  Style bson.ObjectId `bson:",omitempty" json:"style"`
  Script bson.ObjectId `bson:",omitempty" json:"script"`
}

func (this *Creation) Save(s *mgo.Session) {
  if err := s.DB("").C(CREA_C).Insert(&this); err != nil {
    // TODO be more specific for the error
    panic("Creation '" + this.Title + "' failed to be saved")
  }
}

func (this *Creation) FindOneWithKey(s *mgo.Session, k string) {
  err := s.DB("").C(CREA_C).FindId(k).One(&this)
  if err != nil {
    panic("Creation '" + k + "' not found")
  }
}
