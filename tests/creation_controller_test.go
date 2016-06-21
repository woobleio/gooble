package tests

import (
  "fmt"
  "testing"
  "gopkg.in/mgo.v2"
)

func TestCreationPOST(t *testing.T) {
  var s *mgo.Session = GetSession()
  defer s.Close()

  fmt.Print(GetSession())
}
