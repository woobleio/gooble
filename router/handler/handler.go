package handler

import (
	"fmt"
	"net/http"
)

// errors code
const (
	resNotFound int = iota + 100
	dbSelect        // 101
	badForm         // 102
	badParam        // 103
	dbSave          // 104
	servErr         // 105
	badCreds        // 106
	notOwner        // 107
)

// API errors
var (
	ErrBadCreds    = ReqError{badCreds, "Wrong credentials", "%s", http.StatusUnauthorized}
	ErrBadForm     = ReqError{badForm, "Form not valid", "%s", http.StatusBadRequest}
	ErrBadParam    = ReqError{badParam, "Bad param", "Param should be of type %s", http.StatusBadRequest}
	ErrDBSave      = ReqError{dbSave, "Database error", "One or many issues encountered while saving the data :\n %s", http.StatusConflict}
	ErrDBSelect    = ReqError{dbSelect, "Database error", "Failed to select the resources requested", http.StatusInternalServerError}
	ErrNotOwner    = ReqError{notOwner, "Unauthorized", "Authenticated user is not the owner of the resource", http.StatusUnauthorized}
	ErrResNotFound = ReqError{resNotFound, "Resource not found", "%s %s not found", http.StatusNotFound}
	ErrServ        = ReqError{servErr, "Internal server error", "Something wrong happened while processing %s", http.StatusInternalServerError}
)

const (
	Created int = http.StatusCreated
	OK      int = http.StatusOK
)

type JSONRes struct {
	Data   interface{} `json:"data"`
	Errors []ReqError  `json:"errors,omitempty"`
	Status int         `json:"-"`
}

func NewRes() JSONRes {
	return JSONRes{
		nil,
		make([]ReqError, 0),
		http.StatusOK,
	}
}

// Lookup all errors found and return the prioritized HTTP status (the greatest value)
func (j *JSONRes) HttpStatus() int {
	cStatus := 0
	for _, err := range j.Errors {
		cStatus = err.Status
	}
	if cStatus > j.Status {
		j.Status = cStatus
	}
	return j.Status
}

func (j *JSONRes) Response(data interface{}) {
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
