package lib

import "gopkg.in/mgo.v2"

var session *mgo.Session

func GetSession() *mgo.Session {
  return session.Copy()
}

/**
 * Singleton
 */
func SetSession(s *mgo.Session) {
  if session == nil {
    session = s
  }
}
