package error

import (
	"net/http"
	"regexp"
)

// Registry is a map of predefined error codes to errors
type Registry map[ErrorCode]Error

// Merge merges another error collection with this error collection returning a new error collection
func (c Registry) Merge(a Registry) Registry {
	for k, v := range c {
		a[k] = v
	}
	return a
}

var reArray = regexp.MustCompile(`(?m)\.(\d+)(\.)?`)

var registry = Registry{
	ErrorCodeValidationSchema: Error{
		HttpStatusCode: http.StatusBadRequest,
		ErrorCode:      ErrorCodeValidationSchema,
		Message:        "Request body failed JSON schema validation.",
	},
	ErrorCodeValidationJSON: Error{
		HttpStatusCode: http.StatusBadRequest,
		ErrorCode:      ErrorCodeValidationJSON,
		Message:        "Request body contains invalid JSON.",
	},
	ErrorCodeValidationQueryParam: Error{
		HttpStatusCode: http.StatusBadRequest,
		ErrorCode:      ErrorCodeValidationQueryParam,
		Message:        "The value for the query parameter is invalid.",
	},
	ErrorCodeValidationPathParam: Error{
		HttpStatusCode: http.StatusBadRequest,
		ErrorCode:      ErrorCodeValidationPathParam,
		Message:        "The value for the path parameter is invalid.",
	},
	ErrorCodeResourceNotFound: Error{
		HttpStatusCode: http.StatusNotFound,
		ErrorCode:      ErrorCodeResourceNotFound,
		Message:        "Resource not found.",
	},
	ErrorCodeClientUnauthorized: Error{
		HttpStatusCode: http.StatusForbidden,
		ErrorCode:      ErrorCodeClientUnauthorized,
		Message:        "Permission to the requested resource is denied.",
	},
	ErrorCodeClientUnauthenticated: Error{
		HttpStatusCode: http.StatusUnauthorized,
		ErrorCode:      ErrorCodeClientUnauthenticated,
		Message:        "Authentication information is missing or invalid.",
	},
	ErrorCodeServerUnavailable: Error{
		HttpStatusCode: http.StatusServiceUnavailable,
		ErrorCode:      ErrorCodeServerUnavailable,
		Message:        "Server overloaded: unable to process request",
	},
	ErrorCodeServerInternal: Error{
		HttpStatusCode: http.StatusInternalServerError,
		ErrorCode:      ErrorCodeServerInternal,
		Message:        "An internal error has occurred.",
	},
}

func GetRegistryError(ec ErrorCode) Error {
	return deepcopy(registry[ec])
}

func deepcopy(e Error) Error {
	detail := e.SchemaValidationErrors

	if len(detail) > 0 {
		e.SchemaValidationErrors = make([]SchemaValidationError, len(detail))
		copy(e.SchemaValidationErrors, detail)
	}

	return e
}
