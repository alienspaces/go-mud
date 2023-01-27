package error

import (
	"fmt"
	"net/http"
	"strings"
)

type ErrorCode string

const (
	SchemaValidation  ErrorCode = "validation.body_not_matching_json_schema"
	InvalidJSON       ErrorCode = "validation.invalid_json"
	InvalidQueryParam ErrorCode = "validation.invalid_query_parameter"
	InvalidPathParam  ErrorCode = "validation.invalid_path_parameter"
	NotFound          ErrorCode = "resource_not_found"
	Unauthorized      ErrorCode = "unauthorized"
	Unauthenticated   ErrorCode = "unauthenticated"
	Unavailable       ErrorCode = "unavailable"
	Internal          ErrorCode = "internal_error"
)

type ErrorType string

const (
	ErrorTypeUnsupported ErrorType = "unsupported"
	ErrorTypeInvalid     ErrorType = "invalid"
)

type Error struct {
	HttpStatusCode         int                     `json:"-"`
	ErrorCode              ErrorCode               `json:"code"`
	Message                string                  `json:"message"`
	SchemaValidationErrors []SchemaValidationError `json:"validationErrors,omitempty"`
}

func (e Error) Error() string {
	return e.Message
}

func NewInternalError() error {
	return GetRegistryError(Internal)
}

func NewNotFoundError(resourceName string, resourceID string) error {
	e := GetRegistryError(NotFound)
	e.Message = fmt.Sprintf("%s with ID >%s< not found", resourceName, resourceID)

	return e
}

func NewUnavailableError() error {
	return GetRegistryError(Unavailable)
}

func NewUnauthorizedError() error {
	return GetRegistryError(Unauthorized)
}

func NewUnauthenticatedError(message string) error {
	e := GetRegistryError(Unauthenticated)
	e.Message = message
	return e
}

func NewQueryParamError(message string) error {
	e := GetRegistryError(InvalidQueryParam)
	e.Message = message
	return e
}

func ProcessQueryParamError(err error) error {
	e, conversionErr := ToError(err)
	if conversionErr != nil {
		return err
	}

	if len(e.SchemaValidationErrors) == 0 {
		return NewQueryParamError(e.Error())
	}

	errStr := strings.Builder{}
	errStr.WriteString("Invalid query parameter(s): ")
	for i, sve := range e.SchemaValidationErrors {
		if sve.GetField() == "$" {
			errStr.WriteString(fmt.Sprintf("(%d) %s; ", i+1, sve.Message))
		} else {
			errStr.WriteString(fmt.Sprintf("(%d) %s: %s; ", i+1, sve.GetField(), sve.Message))
		}
	}

	formattedErrString := errStr.String()

	// remove extra space and semicolon
	formattedErrString = formattedErrString[0 : len(formattedErrString)-2]
	formattedErrString += "."

	return NewQueryParamError(formattedErrString)
}

func NewPathParamInvalidTypeError(param string, id string) error {
	e := GetRegistryError(InvalidPathParam)
	e.Message = fmt.Sprintf("Path parameter %s value >%s< is not a valid UUID", param, id)
	return e
}

func NewPathParamInvalidError(param, id, message string) error {
	e := GetRegistryError(InvalidPathParam)
	if message == "" {
		e.Message = fmt.Sprintf("Path parameter %s value >%s< is not valid", param, id)
	} else {
		e.Message = fmt.Sprintf("Path parameter %s value >%s< is not valid, %s", param, id, message)
	}
	return e
}

func NewInvalidBodyError(message string) error {
	e := GetRegistryError(InvalidJSON)
	if message != "" {
		e.Message = message
	}
	return e
}

func NewInvalidError(errorCodeSuffix string, message string) error {
	return Error{
		HttpStatusCode: http.StatusBadRequest,
		ErrorCode:      NewErrorCode(ErrorTypeInvalid, errorCodeSuffix),
		Message:        message,
	}
}

func NewUnsupportedError(errorCodeSuffix string, message string) error {
	return Error{
		HttpStatusCode: http.StatusBadRequest,
		ErrorCode:      NewErrorCode(ErrorTypeUnsupported, errorCodeSuffix),
		Message:        message,
	}
}
