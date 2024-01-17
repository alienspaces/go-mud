package error

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/xeipuuv/gojsonschema"
	"gitlab.com/alienspaces/go-mud/backend/core/type/logger"
)

var (
	reArray = regexp.MustCompile(`(?m)\.(\d+)(\.)?`)
)

var internal = Error{
	HttpStatusCode: http.StatusInternalServerError,
	ErrorCode:      Internal,
	Message:        "An internal error has occurred.",
	LogLevel:       logger.ErrorLevel,
	LogMessage:     "",
}

func NewInternalError(message string, args ...any) error {
	e := internal
	if message != "" {
		e.LogLevel = logger.ErrorLevel
		e.LogMessage = fmt.Sprintf(message, args...)
	}
	return e
}

var notFound = Error{
	HttpStatusCode: http.StatusNotFound,
	ErrorCode:      NotFound,
	Message:        "Resource not found",
}

func NewNotFoundError(resource string, id string) error {
	e := notFound
	if resource != "" && id != "" {
		e.Message = fmt.Sprintf("%s with ID >%s< not found", resource, id)
	}
	return e
}

var unavailable = Error{
	HttpStatusCode: http.StatusServiceUnavailable,
	ErrorCode:      Unavailable,
	Message:        "Server overloaded: unable to process request",
}

func NewUnavailableError() error {
	return unavailable
}

var malformed = Error{
	HttpStatusCode: http.StatusBadRequest,
	ErrorCode:      Malformed,
	Message:        "Malformed data: unable to process the request",
}

func NewMalformedError() error {
	return malformed
}

var unauthorized = Error{
	HttpStatusCode: http.StatusForbidden,
	ErrorCode:      Unauthorized,
	Message:        "Permission to the requested resource is denied.",
}

func NewUnauthorizedError() error {
	e := unauthorized
	return e
}

var unauthenticated = Error{
	HttpStatusCode: http.StatusUnauthorized,
	ErrorCode:      Unauthenticated,
	Message:        "Request is unauthenticated",
}

func NewUnauthenticatedError(message string, args ...any) error {
	e := unauthenticated
	if message != "" {
		e.Message = fmt.Sprintf(message, args...)
	}
	return e
}

var invalidParam = Error{
	HttpStatusCode: http.StatusBadRequest,
	ErrorCode:      InvalidParam,
	Message:        "Request contains invalid parameters",
}

func NewParamError(message string, args ...any) error {
	e := invalidParam
	if message != "" {
		e.Message = fmt.Sprintf(message, args...)
	}
	return e
}

var invalidHeader = Error{
	HttpStatusCode: http.StatusBadRequest,
	ErrorCode:      InvalidHeader,
	Message:        "Request contains invalid headers",
}

func NewHeaderError(message string, args ...any) error {
	e := invalidHeader
	if message != "" {
		e.Message = fmt.Sprintf(message, args...)
	}
	return e
}

var invalidData = Error{
	HttpStatusCode: http.StatusBadRequest,
	ErrorCode:      InvalidData,
	Message:        "Request body contains invalid data.",
}

func NewInvalidDataError(message string, args ...any) error {
	e := invalidData
	if message != "" {
		e.Message = fmt.Sprintf(message, args...)
	}
	return e
}

var invalidAction = Error{
	HttpStatusCode: http.StatusConflict,
	ErrorCode:      InvalidAction,
	Message:        "The request conflicts with the current state of the target resource.",
}

func NewInvalidActionError(message string) error {
	e := invalidAction
	if message != "" {
		e.Message = message
	}
	return e
}

var invalidJSON = Error{
	HttpStatusCode: http.StatusBadRequest,
	ErrorCode:      InvalidJSON,
	Message:        "Request body failed JSON schema validation.",
}

func NewInvalidJSONError(resultErrors []gojsonschema.ResultError) error {
	e := invalidJSON

	resultErrors = filterNonUserFriendlyErrors(resultErrors)

	for _, re := range resultErrors {
		sve := setDataPath(SchemaValidationError{}, re)
		sve = setMessage(sve, re)
		e.SchemaValidationErrors = append(e.SchemaValidationErrors, sve)
	}

	return e
}

func filterNonUserFriendlyErrors(re []gojsonschema.ResultError) []gojsonschema.ResultError {
	var friendly []gojsonschema.ResultError
	var unfriendly []gojsonschema.ResultError

	// These errors refer to conditionals in the schema that may not be understood by end-users.
	for _, err := range re {
		errType := err.Type()
		switch errType {
		case "number_any_of", "number_one_of", "number_all_of", "number_not", "condition_then", "condition_else":
			unfriendly = append(unfriendly, err)
		default:
			friendly = append(friendly, err)
		}
	}

	// The non-user friendly errors are _usually_ accompanied by a more specific user-friendly error.
	if len(friendly) == 0 {
		return unfriendly
	}

	return friendly
}

func setDataPath(sve SchemaValidationError, re gojsonschema.ResultError) SchemaValidationError {
	var field string
	if re.Type() == "required" {
		field = re.Details()["property"].(string)
	} else {
		field = re.Field()
	}

	sve.DataPath = "$"

	// not sure if it is possible for the field to be empty, but to be safe the path is set to "$"
	switch field {
	case "", "(root)":
		return sve
	}

	// reformat fields with array index and prefix with "$." (e.g contacts.0.type -> $.contacts[0].type, contacts.0 -> $.contacts[0])
	sve.DataPath = sve.DataPath + "." + reArray.ReplaceAllString(field, "[$1]$2")

	return sve
}

// setMessage sets the detail of the validation error with the reformatted errors returned from the validation.
func setMessage(sve SchemaValidationError, re gojsonschema.ResultError) SchemaValidationError {
	switch re.Type() {
	case "number_gte", "number_gt", "number_lte", "number_lt", "format", "pattern", "array_min_items", "array_max_items":
		sve.Message = re.String()
		if strings.Contains(sve.Message, " 1 items") {
			sve.Message = strings.ReplaceAll(sve.Message, " 1 items", " 1 item")
		}
	default:
		sve.Message = re.Description()
	}

	// clean up message to avoid repeating the property
	if strings.Contains(sve.Message, re.Field()+": ") {
		sve.Message = strings.ReplaceAll(sve.Message, re.Field()+": ", "")
	}
	if strings.Contains(sve.Message, re.Field()+" must") {
		sve.Message = strings.ReplaceAll(sve.Message, re.Field()+" must", "Must")
	}
	if strings.Contains(sve.Message, re.Field()+" does") {
		sve.Message = strings.ReplaceAll(sve.Message, re.Field()+" does", "Does")
	}

	return sve
}

// Dynamic error generation functions
func CreateInvalidError(errorCodeSuffix string, message string, args ...any) error {
	return Error{
		HttpStatusCode: http.StatusBadRequest,
		ErrorCode:      CreateErrorCode("validation.invalid", errorCodeSuffix),
		Message:        fmt.Sprintf(message, args...),
	}
}

func CreateErrorCode(errorType string, field string) ErrorCode {
	return ErrorCode(fmt.Sprintf("%s_%s", errorType, field))
}
