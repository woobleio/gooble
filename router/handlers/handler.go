package handler

import "net/http"

// Http status
const (
	Created   int = http.StatusCreated
	OK        int = http.StatusOK
	NoContent int = http.StatusNoContent
)

// JSONRes is a standardized JSON response
type JSONRes struct {
	Data interface{} `json:"data,omitempty"`
}

// NewRes initializes a reponse
func NewRes(data interface{}) *JSONRes {
	return &JSONRes{
		data,
	}
}

// Response sets the response
func (j *JSONRes) Response(data interface{}) {
	j.Data = data
}
