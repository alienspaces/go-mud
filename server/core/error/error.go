package error

import (
	"errors"
	"fmt"
	"strings"
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

// Registry is a map of error codes to errors
type Registry map[Code]Error

// Merge merges another error collection with this error collection returning a new error collection
func (c Registry) Merge(a Registry) Registry {
	for k, v := range c {
		a[k] = v
	}
	return a
}

type SchemaValidationError struct {
	DataPath string `json:"dataPath"`
	Message  string `json:"message"`
}

func (sve SchemaValidationError) GetField() string {
	field := strings.Split(sve.DataPath, ".")
	lastField := field[len(field)-1]
	return lastField
}

func IsError(e error) bool {
	var errorPtr Error
	return errors.As(e, &errorPtr)
}

func HasErrorCode(err error, c Code) bool {
	e, err := ToError(err)
	if err != nil {
		return false
	}

	return e.ErrorCode == c
}

func ToError(e error) (Error, error) {
	if e == nil {
		return Error{}, fmt.Errorf("err is nil when converting to coreerror.Error type")
	}

	var err Error
	if !errors.As(e, &err) {
		return Error{}, fmt.Errorf("failed to convert to coreerror.Error type >%v<", e)
	}

	if len(err.SchemaValidationErrors) == 0 {
		err.SchemaValidationErrors = nil
	}

	return err, nil
}
