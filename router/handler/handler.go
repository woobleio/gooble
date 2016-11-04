package handler

import (
  "net/http"
  "fmt"
)

// Each err type is for common error cases
// One err type has one or many errors code, error code is defined by iota
type errNotFound int
const (
  ResNotFound errNotFound = iota + 100
  NotFound
)

var errsDetail = map[interface{}]string{
  ResNotFound: "%s %s not found",
  NotFound: "No %s found",
}

var errsTitle = map[interface{}]string{
  0: "Unknown error",
  ResNotFound: "Resource not found",
}

type JSONRes struct {
  Data   interface{}  `json:"data"`
  Errors []ReqError   `json:"errors,omitempty"`
}

func NewRes() JSONRes {
  return JSONRes{
    nil,
    make([]ReqError, 0),
  }
}

// Lookup all errors found and return the prioritized HTTP status (the greatest value)
func (j *JSONRes) HttpStatus() int {
  cStatus := 0
  status := http.StatusOK
  for _, err := range j.Errors {
    switch err.Code.(type) {
    case errNotFound :
      cStatus = http.StatusNotFound
    }
  }
  if cStatus > status {
    status = cStatus
  }
  return status
}

func (j *JSONRes) Response(data interface{}) {
  j.Data = data
}

type ReqError struct {
  Code    interface{}    `json:"code"`
  Title   string `json:"title"`
  Message string `json:"detail"`
}

func (j *JSONRes) Error(err interface{}, args ...interface{}) {
  reqErr := ReqError{
    err,
    getErrTitle(err),
    fmt.Sprintf(errsDetail[err], args...),
  }
  j.Errors = append(j.Errors, reqErr)
}

// A title concerns one or many errors
func getErrTitle(err interface{}) string {
  var title string
  switch err {
  case ResNotFound:
    fallthrough
  case NotFound:
    title = errsTitle[ResNotFound]
  default:
    title = errsTitle[0]
  }
  return title
}
