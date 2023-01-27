package model

import (
	"net/http"

	coreerror "gitlab.com/alienspaces/go-mud/backend/core/error"
)

const (
	ErrorCodeActionInvalid  coreerror.ErrorCode = "action.invalid"
	ErrorCodeActionTooEarly coreerror.ErrorCode = "action.too_early"
)

func NewActionInvalidError(message string) error {
	return coreerror.Error{
		HttpStatusCode: http.StatusBadRequest,
		ErrorCode:      ErrorCodeActionInvalid,
		Message:        message,
	}
}

func NewActionTooEarlyError(message string) error {
	return coreerror.Error{
		HttpStatusCode: http.StatusBadRequest,
		ErrorCode:      ErrorCodeActionTooEarly,
		Message:        message,
	}
}
