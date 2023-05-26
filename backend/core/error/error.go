package error

import (
	"errors"
	"fmt"
	"strings"
)

type ErrorCode string

const (
	SchemaValidation ErrorCode = "validation.body_not_matching_json_schema"
	InvalidAction    ErrorCode = "invalid_action"
	InvalidJSON      ErrorCode = "validation.invalid_json"
	InvalidHeader    ErrorCode = "validation.invalid_header"
	InvalidParam     ErrorCode = "validation.invalid_parameter"
	NotFound         ErrorCode = "resource_not_found"
	Unauthorized     ErrorCode = "unauthorized"
	Unauthenticated  ErrorCode = "unauthenticated"
	Unavailable      ErrorCode = "unavailable"
	Malformed        ErrorCode = "malformed"
	Internal         ErrorCode = "internal_error"
)

type Error struct {
	HttpStatusCode         int                     `json:"-"`
	ErrorCode              ErrorCode               `json:"code"`
	Message                string                  `json:"message"`
	SchemaValidationErrors []SchemaValidationError `json:"validationErrors,omitempty"`
}

func (e Error) Error() string {
	return fmt.Sprintf("%s: %s", e.ErrorCode, e.Message)
}

// Registry is a map of error codes to errors
type Registry map[ErrorCode]Error

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

func HasErrorCode(err error, c ErrorCode) bool {
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

func ToErrors(errs ...error) ([]Error, error) {
	var results []Error

	for _, e := range errs {
		result, err := ToError(e)
		if err != nil {
			return nil, fmt.Errorf("failed to convert err to error result >%#v<", err)
		}

		results = append(results, result)
	}

	return results, nil
}

func ProcessParamError(err error) error {
	e, conversionErr := ToError(err)
	if conversionErr != nil {
		return err
	}

	if len(e.SchemaValidationErrors) == 0 {
		return NewParamError(e.Error())
	}

	errStr := strings.Builder{}
	errStr.WriteString("Invalid parameter(s): ")
	for i, sve := range e.SchemaValidationErrors {
		if sve.GetField() == "$" {
			errStr.WriteString(fmt.Sprintf("(%d) %s; ", i+1, sve.Message))
		} else {
			errStr.WriteString(fmt.Sprintf("(%d) %s: %s; ", i+1, sve.GetField(), sve.Message))
		}
	}

	formattedErrString := errStr.String()
	formattedErrString = formattedErrString[0 : len(formattedErrString)-2] // remove extra space and semicolon
	formattedErrString += "."
	return NewParamError(formattedErrString)
}
