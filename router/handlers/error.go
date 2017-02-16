package handler

import (
	"fmt"
	"net/http"

	"gopkg.in/gin-gonic/gin.v1"
	validator "gopkg.in/go-playground/validator.v8"
)

type errCode string

// errors code
const (
	aliasRequired    errCode = "alias_required"
	alreadyCreaOwner errCode = "already_owner"
	badCreds         errCode = "bad_credentials"
	badForm          errCode = "bad_form"
	chargeErr        errCode = "charge_fail"
	dbFailSave       errCode = "db_failed_save"
	dbFailSelect     errCode = "db_failed_select"
	dbFailUpdate     errCode = "db_failed_update"
	mustBuy          errCode = "must_buy"
	planLimit        errCode = "plan_limit"
	resNotFound      errCode = "res_not_found"
	servErr          errCode = "server_error"
)

// API errors
var (
	ErrAliasRequired = NewAPIError(aliasRequired, "Alias required", "Creation name should be unique in a package : creation %s should have an alias", http.StatusBadRequest)
	ErrBadCreds      = NewAPIError(badCreds, "Wrong credentials", "Unknown email or password invalid", http.StatusUnauthorized)
	ErrBadForm       = NewAPIError(badForm, "Form not valid", "", http.StatusBadRequest)
	ErrCantBuy       = NewAPIError(alreadyCreaOwner, "Purchase failed", "Can't buy the creation %s because you already own it", http.StatusBadRequest)
	ErrCharge        = NewAPIError(chargeErr, "Charge failed", "Couldn't charge", http.StatusBadRequest)
	ErrDBSave        = NewAPIError(dbFailSave, "Database error", "Failed to save the data", http.StatusConflict)
	ErrDBSelect      = NewAPIError(dbFailSelect, "Database error", "Failed to select the resources requested", http.StatusInternalServerError)
	ErrMustBuy       = NewAPIError(mustBuy, "Must purchase before doing this", "One or some creations must be purchased to do this", http.StatusUnauthorized)
	ErrPlanLimit     = NewAPIError(planLimit, "Plan limit exceeded", "Number of %s limited by actual plan %s", http.StatusUnauthorized)
	ErrResNotFound   = NewAPIError(resNotFound, "Resource not found", "%s %v not found", http.StatusNotFound)
	ErrServ          = NewAPIError(servErr, "Internal server error", "Something wrong happened while processing %s", http.StatusInternalServerError)
	ErrUpdate        = NewAPIError(dbFailUpdate, "Database error", "Failed to update %s %v", http.StatusInternalServerError)
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
func NewAPIError(code errCode, title string, details string, status int) *APIError {
	return &APIError{
		code,
		title,
		details,
		status,
		make(map[string]interface{}),
	}
}

// SetParams adds params to APIError, parameters must be in as the following : key(string), value(interface) ...
func (e *APIError) SetParams(params ...interface{}) *APIError {
	lenParams := len(params)
	if lenParams%2 > 0 {
		panic("Params in APIErrors should be even such as key:value")
	}

	detailsParams := make([]interface{}, 0)

	for i := 1; i < lenParams; i++ {
		index := fmt.Sprintf("%v", params[i-1])
		e.Params[index] = params[i]
		detailsParams = append(detailsParams, params[i])
	}

	e.Details = fmt.Sprintf(e.Details, detailsParams...)

	// FIXME
	fmt.Println("-->", detailsParams)

	return e
}

// ValidationError builds and sets validation errors
func (e *APIError) ValidationError(ve *validator.FieldError) {
	switch ve.Tag {
	case "required":
		e.Details = "%s is required"
		e.SetParams("field", ve.Field)
	case "max":
		e.Details = "%s cannot be longer than %s"
		e.SetParams("field", ve.Field, "param", ve.Param)
	case "min":
		e.Details = "%s must be longer than %s"
		e.SetParams("field", ve.Field, "param", ve.Param)
	case "email":
		e.Details = "Invalid email format"
	case "len":
		e.Details = "%s must be %s characters long"
		e.SetParams("field", ve.Field, "param", ve.Param)
	}
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

// HandleErrors handle API errors TODO put warnings
func HandleErrors(c *gin.Context) {
	// FIXME workaroun gin issue with Bind (https://github.com/gin-gonic/gin/issues/633)
	c.Header("Content-Type", gin.MIMEJSON)

	c.Next()

	apiErrors := &APIErrors{
		make([]APIError, 0),
	}
	for _, err := range c.Errors {
		if err.Meta != nil {
			publicError := err.Meta.(*APIError)
			switch err.Type {
			case gin.ErrorTypeBind:
				valErrors := err.Err.(validator.ValidationErrors)
				for _, valErr := range valErrors {
					publicError.ValidationError(valErr)
					break
				}
				apiErrors.Error(*publicError)
			}
			apiErrors.Error(*publicError)
		}
	}
	if apiErrors.HasErrors() {
		c.JSON(apiErrors.HTTPStatus(), apiErrors)
	}
}
