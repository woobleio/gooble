package models

import (
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
)

type Creation struct {
  Id bson.ObjectId `bson:",omitempty"`
  Title string `bson:"_id"`
  Dom bson.ObjectId `json:"dom"`
  Style bson.ObjectId `bson:",omitempty"`
  Script bson.ObjectId `bson:",omitempty"`
}

func (this *Creation) Save(s *mgo.Session) {
  if err := s.DB("").C(CREA_C).Insert(&this); err != nil {
    // TODO be more specific for the error
    panic("Creation '" + this.Title + "' failed to be saved")
  }
}

func (this *Creation) Populate(s *mgo.Session) {
  if err := s.DB("").C(CREA_C).FindId(this.Title).One(&this); err != nil {
    panic("Creation '" + this.Title + "' not found")
  }
}

func (this *Creation) Create(s *mgo.Session) {
  this.Id = bson.NewObjectId()
}
