package error

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/xeipuuv/gojsonschema"
)

var (
	reArray = regexp.MustCompile(`(?m)\.(\d+)(\.)?`)
)

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

func NewPathParamError(param string, id string) error {
	e := GetRegistryError(InvalidPathParam)
	e.Message = fmt.Sprintf("%s >%s< is not a valid UUID", param, id)

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
		ErrorCode:      CreateErrorCode(ValidationErrorInvalid, errorCodeSuffix),
		Message:        message,
	}
}

func NewUnsupportedError(errorCodeSuffix string, message string) error {
	return Error{
		HttpStatusCode: http.StatusBadRequest,
		ErrorCode:      CreateErrorCode(ValidationErrorUnsupported, errorCodeSuffix),
		Message:        message,
	}
}

func CreateErrorCode(errorType ValidationErrorType, field string) Code {
	return Code(fmt.Sprintf("%s_%s", errorType, field))
}

func NewSchemaValidationError(resultErrors []gojsonschema.ResultError) error {
	e := GetRegistryError(SchemaValidation)

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
