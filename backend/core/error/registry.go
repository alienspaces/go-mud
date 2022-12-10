package error

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

// Registry is a map of error codes to errors
type Registry map[Code]Error

// Merge merges another error collection with this error collection returning a new error collection
func (c Registry) Merge(a Registry) Registry {
	for k, v := range c {
		a[k] = v
	}
	return a
}

var reArray = regexp.MustCompile(`(?m)\.(\d+)(\.)?`)

var registry = Registry{
	SchemaValidation: Error{
		HttpStatusCode: http.StatusBadRequest,
		ErrorCode:      SchemaValidation,
		Message:        "Request body failed JSON schema validation.",
	},
	InvalidJSON: Error{
		HttpStatusCode: http.StatusBadRequest,
		ErrorCode:      InvalidJSON,
		Message:        "Request body contains invalid JSON.",
	},
	InvalidQueryParam: Error{
		HttpStatusCode: http.StatusBadRequest,
		ErrorCode:      InvalidQueryParam,
		Message:        "The value for the query parameter is invalid.",
	},
	InvalidPathParam: Error{
		HttpStatusCode: http.StatusBadRequest,
		ErrorCode:      InvalidPathParam,
		Message:        "The value for the path parameter is invalid.",
	},
	NotFound: Error{
		HttpStatusCode: http.StatusNotFound,
		ErrorCode:      NotFound,
		Message:        "Resource not found.",
	},
	Unauthorized: Error{
		HttpStatusCode: http.StatusForbidden,
		ErrorCode:      Unauthorized,
		Message:        "Permission to the requested resource is denied.",
	},
	Unauthenticated: Error{
		HttpStatusCode: http.StatusUnauthorized,
		ErrorCode:      Unauthenticated,
		Message:        "Authentication information is missing or invalid.",
	},
	Unavailable: Error{
		HttpStatusCode: http.StatusServiceUnavailable,
		ErrorCode:      Unavailable,
		Message:        "Server overloaded: unable to process request",
	},
	Internal: Error{
		HttpStatusCode: http.StatusInternalServerError,
		ErrorCode:      Internal,
		Message:        "An internal error has occurred.",
	},
}

func GetRegistryError(code Code) Error {
	return deepcopy(registry[code])
}

func deepcopy(e Error) Error {
	detail := e.SchemaValidationErrors

	if len(detail) > 0 {
		e.SchemaValidationErrors = make([]SchemaValidationError, len(detail))
		copy(e.SchemaValidationErrors, detail)
	}

	return e
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
	formattedErrString = formattedErrString[0 : len(formattedErrString)-2] // remove extra space and semicolon
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
		ErrorCode:      createErrorCode(ValidationErrorInvalid, errorCodeSuffix),
		Message:        message,
	}
}

func NewUnsupportedError(errorCodeSuffix string, message string) error {
	return Error{
		HttpStatusCode: http.StatusBadRequest,
		ErrorCode:      createErrorCode(ValidationErrorUnsupported, errorCodeSuffix),
		Message:        message,
	}
}

func createErrorCode(errorType ValidationErrorType, field string) Code {
	return Code(fmt.Sprintf("%s_%s", errorType, field))
}
