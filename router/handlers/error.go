package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	validator "gopkg.in/go-playground/validator.v9"
)

type errCode string

// errors code
const (
	activationFail   errCode = "activation_fail"
	aliasRequired    errCode = "alias_required"
	alreadyCreaOwner errCode = "already_owner"
	badCreds         errCode = "bad_credentials"
	badFileFormat    errCode = "bad_file_format"
	badForm          errCode = "bad_form"
	badScriptDoc     errCode = "bad_script_document"
	badScriptClass   errCode = "bad_script_class"
	badScriptConst   errCode = "bad_script_constructor"
	betterPlan       errCode = "better_plan"
	chargeErr        errCode = "charge_fail"
	creaNotAvail     errCode = "creation_not_available"
	creaVersion      errCode = "creation_version"
	dbFail           errCode = "db_failure"
	planLimit        errCode = "plan_limit"
	resNotFound      errCode = "res_not_found"
	servErr          errCode = "server_error"
	servIntErr       errCode = "server_internal_error"
)

// API errors
var (
	ErrActivation     = NewAPIError(activationFail, "Activation failed", "The token might be wrong", http.StatusConflict)
	ErrAliasRequired  = NewAPIError(aliasRequired, "Alias required", "Creation name should be unique in a package : creation %s should have an alias", http.StatusBadRequest)
	ErrBadCreds       = NewAPIError(badCreds, "Wrong credentials", "Unknown email or password invalid", http.StatusUnauthorized)
	ErrBadFileFormat  = NewAPIError(badFileFormat, "Server can't handle this file", "Allowed formats are : %s", http.StatusBadRequest)
	ErrBadForm        = NewAPIError(badForm, "Form invalid", "", http.StatusBadRequest)
	ErrBadScriptDoc   = NewAPIError(badScriptDoc, "Missing document initializer", "Document initializer is required and shloud be in the constructor : %s", http.StatusBadRequest)
	ErrBadScriptClass = NewAPIError(badScriptClass, "Missing class", "Class Woobly is mandatory", http.StatusBadRequest)
	ErrBadScriptConst = NewAPIError(badScriptConst, "Missing constructor", "A constructor is mandatory : %s", http.StatusBadRequest)
	ErrBetterPlan     = NewAPIError(betterPlan, "Already has a better plan", "Current plan %s and requested plan is %s", http.StatusConflict)
	ErrCreaNotAvail   = NewAPIError(creaNotAvail, "Creation not available", "The creation %s is not available", http.StatusConflict)
	ErrCreaVersion    = NewAPIError(creaVersion, "Bad version", "Version %s can't be created", http.StatusBadRequest)
	ErrDB             = NewAPIError(dbFail, "Database error", "Database failed to process the request", http.StatusConflict)
	ErrIntServ        = NewAPIError(servIntErr, "Internal server error", "Something wrong happened", http.StatusInternalServerError)
	ErrPlanLimit      = NewAPIError(planLimit, "Plan limit exceeded", "Number of %s limited by actual plan %s", http.StatusUnauthorized)
	ErrResNotFound    = NewAPIError(resNotFound, "Resource not found", "%s %v not found", http.StatusNotFound)
	ErrServ           = NewAPIError(servErr, "Internal server error", "Something wrong happened while processing %s", http.StatusInternalServerError)
)

// APIError is a struct that standardize a Wooble error
type APIError struct {
	Code    errCode                `json:"code"`
	Title   string                 `json:"title"`
	Details string                 `json:"details,omitempty"`
	Status  int                    `json:"status"`
	Params  map[string]interface{} `json:"params,omitempty"`
}

// APIErrors wrap all API errors
type APIErrors struct {
	Errors []APIError `json:"errors,omitempty"`
}

// NewAPIError creates an APIError
func NewAPIError(code errCode, title string, details string, status int) APIError {
	return APIError{
		code,
		title,
		details,
		status,
		make(map[string]interface{}),
	}
}

// SetParams adds params to APIError, parameters must be in as the following : key(string), value(interface) ...
func (e APIError) SetParams(params ...interface{}) APIError {
	lenParams := len(params)
	if lenParams%2 > 0 {
		panic("Params in APIErrors should be even such as key:value")
	}

	detailsParams := make([]interface{}, 0)

	for i := 1; i < lenParams; i = i + 2 {
		index := fmt.Sprintf("%v", params[i-1])
		e.Params[index] = params[i]
		detailsParams = append(detailsParams, params[i])
	}

	e.Details = fmt.Sprintf(e.Details, detailsParams...)

	return e
}

// ValidationError builds and sets validation errors
func (e APIError) ValidationError(ve validator.FieldError) APIError {
	switch ve.Tag() {
	case "required":
		e.Details = "%s is required"
		e = e.SetParams("field", ve.Field())
	case "max":
		e.Details = "%s cannot be longer than %s"
		e = e.SetParams("field", ve.Field(), "param", ve.Param())
	case "min":
		e.Details = "%s must be longer than %s"
		e = e.SetParams("field", ve.Field(), "param", ve.Param())
	case "email":
		e.Details = "Invalid email format"
	case "len":
		e.Details = "%s must be %s characters long"
		e = e.SetParams("field", ve.Field(), "param", ve.Param())
	case "alpha":
		e.Details = "%s must be one word"
		e = e.SetParams("field", ve.Field())
	}

	return e
}

// HTTPStatus returns the HTTP status of errors
func (e *APIErrors) HTTPStatus() int {
	cStatus := 0
	status := 0
	for _, err := range e.Errors {
		cStatus = err.Status
	}
	if cStatus > status {
		status = cStatus
	}
	return status
}

// Error appends a new API error
func (e *APIErrors) Error(err APIError) {
	e.Errors = append(e.Errors, err)
}

// HasErrors tells if the is any error
func (e *APIErrors) HasErrors() bool {
	return len(e.Errors) > 0
}

// HandleErrors handle API errors
func HandleErrors(c *gin.Context) {
	// FIXME workaroun gin issue with Bind (https://github.com/gin-gonic/gin/issues/633)
	c.Header("Content-Type", gin.MIMEJSON)

	c.Next()

	apiErrors := &APIErrors{
		make([]APIError, 0),
	}
	for _, err := range c.Errors {
		if err.Meta != nil {
			publicError := err.Meta.(APIError)
			switch err.Type {
			case gin.ErrorTypeBind:
				fmt.Print(err.Err)
				valErrors := err.Err.(validator.ValidationErrors)
				for _, valErr := range valErrors {
					publicError = publicError.ValidationError(valErr)
					apiErrors.Error(publicError)
				}
			default:
				apiErrors.Error(publicError)
			}
		}
	}
	if apiErrors.HasErrors() {
		c.JSON(apiErrors.HTTPStatus(), apiErrors)
	}
}
