package error

import (
	"errors"
	"fmt"
)

type Code string

const (
	SchemaValidation  Code = "validation.body_not_matching_json_schema"
	InvalidJSON       Code = "validation.invalid_json"
	InvalidQueryParam Code = "validation.invalid_query_parameter"
	InvalidPathParam  Code = "validation.invalid_path_parameter"
	NotFound          Code = "resource_not_found"
	Unauthorized      Code = "unauthorized"
	Unauthenticated   Code = "unauthenticated"
	Unavailable       Code = "unavailable"
	Internal          Code = "internal_error"
)

type Error struct {
	HttpStatusCode         int                     `json:"-"`
	ErrorCode              Code                    `json:"code"`
	Message                string                  `json:"message"`
	SchemaValidationErrors []SchemaValidationError `json:"validationErrors,omitempty"`
}

func (e Error) Error() string {
	return e.Message
}

func IsError(err error) bool {
	var errorPtr Error
	return errors.As(err, &errorPtr)
}

func HasErrorCode(err error, c Code) bool {
	e, err := ToError(err)
	if err != nil {
		return false
	}

	return e.ErrorCode == c
}

func ToError(err error) (Error, error) {
	if err == nil {
		return Error{}, fmt.Errorf("err is nil when converting to coreerror.Error type")
	}

	var errorPtr Error
	if !errors.As(err, &errorPtr) {
		return Error{}, fmt.Errorf("failed to convert to coreerror.Error type >%v<", err)
	}

	if len(errorPtr.SchemaValidationErrors) == 0 {
		errorPtr.SchemaValidationErrors = nil
	}

	return errorPtr, nil
}
