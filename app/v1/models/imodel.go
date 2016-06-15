package models

import (
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
)

const CREA_C = "creations"
const DOM_C = "documents"
const SCRIPT_C = "scripts"
const STYLE_C = "styles"

type IModel interface {
  Create(*mgo.Session)
  Save(*mgo.Session)
  FindOne(*mgo.Session, bson.ObjectId)
  FindOneWithKey(*mgo.Session, string)
}
