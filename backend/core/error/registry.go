package error

import (
	"net/http"
	"regexp"
)

// Registry is a map of error codes to errors
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
