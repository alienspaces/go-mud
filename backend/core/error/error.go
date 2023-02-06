package error

import (
	"fmt"
	"net/http"
	"strings"
)

const (
	// TODO: (core) Change to data
	ErrorCategoryValidation ErrorCategory = "validation"
	// TODO: (core) Additional categories data, query_parameter and path parameter
	ErrorCategoryResource ErrorCategory = "resource"
	ErrorCategoryClient   ErrorCategory = "client"
	ErrorCategoryServer   ErrorCategory = "server"
)

const (
	// TODO: (core) Change the following to data.invalid
	ErrorCodeValidationSchema ErrorCode = "validation.schema"
	// TODO: (core) Change the following to data.invalid_json
	ErrorCodeValidationJSON ErrorCode = "validation.invalid_json"
	// TODO: (core) Change the following to query_parameter.invalid
	ErrorCodeValidationQueryParam ErrorCode = "validation.invalid_query_parameter"
	// TODO: (core) Change the following to path_parameter.invalid
	ErrorCodeValidationPathParam   ErrorCode = "validation.invalid_path_parameter"
	ErrorCodeResourceNotFound      ErrorCode = "resource.not_found"
	ErrorCodeClientUnauthorized    ErrorCode = "client.unauthorized"
	ErrorCodeClientUnauthenticated ErrorCode = "client.unauthenticated"
	ErrorCodeServerUnavailable     ErrorCode = "server.unavailable"
	ErrorCodeServerInternal        ErrorCode = "server.internal_error"
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

type ErrorCategory string
type ErrorCode string

func NewErrorCode(ec ErrorCategory, et ErrorType, s string) ErrorCode {
	return ErrorCode(fmt.Sprintf("%s.%s_%s", ec, et, s))
}

func HasErrorCode(err error, ec ErrorCode) bool {
	e, err := ToError(err)
	if err != nil {
		return false
	}

	return e.ErrorCode == ec
}

func NewResourceNotFoundError(resourceName string, resourceID string) error {
	e := GetRegistryError(ErrorCodeResourceNotFound)
	e.Message = fmt.Sprintf("%s with ID >%s< not found", resourceName, resourceID)

	return e
}

func NewServerInternalError() error {
	return GetRegistryError(ErrorCodeServerInternal)
}

func NewServerUnavailableError() error {
	return GetRegistryError(ErrorCodeServerUnavailable)
}

func NewClientUnauthorizedError() error {
	return GetRegistryError(ErrorCodeClientUnauthorized)
}

func NewClientUnauthenticatedError(message string) error {
	e := GetRegistryError(ErrorCodeClientUnauthenticated)
	e.Message = message
	return e
}

func NewValidationQueryParamError(message string) error {
	e := GetRegistryError(ErrorCodeValidationQueryParam)
	e.Message = message
	return e
}

func ProcessValidationQueryParamError(err error) error {
	e, conversionErr := ToError(err)
	if conversionErr != nil {
		return err
	}

	if len(e.SchemaValidationErrors) == 0 {
		return NewValidationQueryParamError(e.Error())
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

	return NewValidationQueryParamError(formattedErrString)
}

func NewValidationPathParamTypeError(param string, id string) error {
	e := GetRegistryError(ErrorCodeValidationPathParam)
	e.Message = fmt.Sprintf("Path parameter %s value >%s< is not a valid UUID", param, id)
	return e
}

func NewValidationPathParamError(param, id, message string) error {
	e := GetRegistryError(ErrorCodeValidationPathParam)
	if message == "" {
		e.Message = fmt.Sprintf("Path parameter %s value >%s< is not valid", param, id)
	} else {
		e.Message = fmt.Sprintf("Path parameter %s value >%s< is not valid, %s", param, id, message)
	}
	return e
}

func NewValidationJSONError(message string) error {
	e := GetRegistryError(ErrorCodeValidationJSON)
	if message != "" {
		e.Message = message
	}
	return e
}

type ErrorType string

const (
	ErrorTypeUnsupported ErrorType = "unsupported"
	ErrorTypeInvalid     ErrorType = "invalid"
)

func NewValidationInvalidError(errorCodeSuffix string, message string) error {
	return Error{
		HttpStatusCode: http.StatusBadRequest,
		ErrorCode:      NewErrorCode(ErrorCategoryValidation, ErrorTypeInvalid, errorCodeSuffix),
		Message:        message,
	}
}

func NewValidationUnsupportedError(errorCodeSuffix string, message string) error {
	return Error{
		HttpStatusCode: http.StatusBadRequest,
		ErrorCode:      NewErrorCode(ErrorCategoryValidation, ErrorTypeUnsupported, errorCodeSuffix),
		Message:        message,
	}
}
