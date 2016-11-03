package handler

import (
  "fmt"
)

const (
  NotFound = "%s %s not found"
)

type ReqError struct {
  Message string `json:"message"`
  Code    int    `json:"code"`
}

func NewError(mess string, args ...interface{}) ReqError {
  return ReqError{
    fmt.Sprintf(mess, args...),
    1,
  }
}
