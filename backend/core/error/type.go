package error

import (
	"fmt"
	"strings"
)

type ErrorCode string

const (
	InvalidAction    ErrorCode = "invalid_action"
	SchemaValidation ErrorCode = "validation.body_not_matching_json_schema"
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

type SchemaValidationError struct {
	DataPath string `json:"dataPath"`
	Message  string `json:"message"`
}

func (sve SchemaValidationError) GetField() string {
	field := strings.Split(sve.DataPath, ".")
	lastField := field[len(field)-1]
	return lastField
}
