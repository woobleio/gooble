package models

import (
  "gopkg.in/mgo.v2"
)

const CREA_C = "creations"
const DOM_C = "documents"
const SCRIPT_C = "scripts"
const STYLE_C = "styles"

type Model interface {
  Create(*mgo.Session)
  Save(*mgo.Session)
  Populate(*mgo.Session)
}
