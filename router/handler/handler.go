package handler

import (
	"fmt"
	"net/http"
)

// errors code
const (
	cResNotFound int = iota + 100
	cNotFound
)

// API errors
var (
	ErrResNotFound = ReqError{cResNotFound, "Resource not found", "%s %s not found", http.StatusNotFound}
	ErrNotFound    = ReqError{cNotFound, "Resource not found", "No %s found", http.StatusNotFound}
)

type JSONRes struct {
	Data   *interface{} `json:"data"`
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
		cStatus = err.Status
	}
	if cStatus > status {
		status = cStatus
	}
	return status
}

func (j *JSONRes) Response(data *interface{}) {
	j.Data = data
}

type ReqError struct {
	Code    interface{} `json:"code"`
	Title   string      `json:"title"`
	Details string      `json:"details"`
	Status  int         `json:"-"`
}

func (j *JSONRes) Error(err ReqError, args ...interface{}) {
	err.Details = fmt.Sprintf(err.Details, args...)
	j.Errors = append(j.Errors, err)
}
