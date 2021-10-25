package server

import (
	"fmt"
	"net/http"

	"gitlab.com/alienspaces/go-mud/server/core/type/logger"
)

// ErrorCode -
const (
	ErrorCodeUnauthorized   string = "unauthorized_error"
	ErrorDetailUnauthorized string = "request could not be authorized"
	ErrorCodeSystem         string = "internal_error"
	ErrorDetailSystem       string = "an internal error has occurred"
	ErrorCodeValidation     string = "validation_error"
	ErrorDetailValidation   string = "request contains validation errors"
	ErrorCodeNotFound       string = "not_found"
	ErrorDetailNotFound     string = "requested resource could not be found"
)

// WriteUnauthorizedError -
func (rnr *Runner) WriteUnauthorizedError(l logger.Logger, w http.ResponseWriter, err error) {

	l.Warn("Unauthorized error >%v<", err)

	// Unauthorized error
	res := rnr.UnauthorizedError(err)

	err = rnr.WriteResponse(l, w, res)
	if err != nil {
		l.Warn("Failed writing response >%v<", err)
		return
	}
}

// UnauthorizedError -
func (rnr *Runner) UnauthorizedError(err error) Response {

	rnr.Log.Error("Error >%v<", err)

	return Response{
		Error: &ResponseError{
			Code:   ErrorCodeUnauthorized,
			Detail: err.Error(),
		},
	}
}

// WriteModelError -
func (rnr *Runner) WriteModelError(l logger.Logger, w http.ResponseWriter, err error) {

	l.Warn("Model error >%v<", err)

	// model error
	res := rnr.ModelError(err)

	err = rnr.WriteResponse(l, w, res)
	if err != nil {
		l.Warn("Failed writing response >%v<", err)
		return
	}
}

// ModelError -
func (rnr *Runner) ModelError(err error) Response {

	rnr.Log.Error("Error >%v<", err)

	return Response{
		Error: &ResponseError{
			Code:   ErrorCodeValidation,
			Detail: err.Error(),
		},
	}
}

// WriteSystemError -
func (rnr *Runner) WriteSystemError(l logger.Logger, w http.ResponseWriter, err error) {

	l.Warn("System error >%v<", err)

	// system error
	res := rnr.SystemError(err)

	err = rnr.WriteResponse(l, w, res)
	if err != nil {
		rnr.Log.Warn("Failed writing response >%v<", err)
		return
	}
}

// SystemError -
func (rnr *Runner) SystemError(err error) Response {

	rnr.Log.Error("Error >%v<", err)

	// NOTE: never expose actual system error details
	return Response{
		Error: &ResponseError{
			Code:   ErrorCodeSystem,
			Detail: ErrorDetailSystem,
		},
	}
}

// ValidationError -
func (rnr *Runner) ValidationError(err error) Response {

	rnr.Log.Error("error >%v<", err)

	if err == nil {
		err = fmt.Errorf(ErrorDetailValidation)
	}

	return Response{
		Error: &ResponseError{
			Code:   ErrorCodeValidation,
			Detail: err.Error(),
		},
	}
}

// WriteNotFoundError -
func (rnr *Runner) WriteNotFoundError(l logger.Logger, w http.ResponseWriter, id string) {

	err := fmt.Errorf("resource with ID >%s< not found", id)

	l.Warn("Not found error >%v<", err)

	// not found error
	res := rnr.NotFoundError(err)

	err = rnr.WriteResponse(l, w, res)
	if err != nil {
		l.Warn("Failed writing response >%v<", err)
		return
	}
}

// NotFoundError -
func (rnr *Runner) NotFoundError(err error) Response {

	rnr.Log.Error("Error >%v<", err)

	if err == nil {
		err = fmt.Errorf(ErrorDetailNotFound)
	}

	return Response{
		Error: &ResponseError{
			Code:   ErrorCodeNotFound,
			Detail: err.Error(),
		},
	}
}
