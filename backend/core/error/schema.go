package error

import (
	"strings"

	"github.com/xeipuuv/gojsonschema"
)

type SchemaValidationError struct {
	DataPath string `json:"dataPath"`
	Message  string `json:"message"`
}

func (sve SchemaValidationError) GetField() string {
	field := strings.Split(sve.DataPath, ".")
	lastField := field[len(field)-1]
	return lastField
}

func NewSchemaValidationError(resultErrors []gojsonschema.ResultError) error {
	e := GetRegistryError(ErrorCodeValidationSchema)

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
