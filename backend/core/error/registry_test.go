package error

import (
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/go-mud/backend/core/jsonschema"
)

func TestNewSchemaError(t *testing.T) {
	tests := []struct {
		name string
		data string
		want Error
	}{
		{
			name: "Missing required field",
			data: `{}`,
			want: Error{
				SchemaValidationErrors: []SchemaValidationError{
					{
						DataPath: "$.first_name",
						Message:  "first_name is required",
					},
					{
						DataPath: "$.last_name",
						Message:  "last_name is required",
					},
					{
						DataPath: "$.age",
						Message:  "age is required",
					},
				},
			},
		},
		{
			name: "Nullable field with invalid type",
			data: `
				{
					"first_name": "John",
					"last_name": "Doe",
					"age": 20,
					"proof_of_age": {
						"id_type": "ssn",
						"id_number": "A8693170"
					},
					"nullable_string_field": 1234
				}
					`,
			want: Error{
				SchemaValidationErrors: []SchemaValidationError{
					{
						DataPath: "$.nullable_string_field",
						Message:  "Invalid type. Expected: string, given: integer",
					},
				},
			},
		},
		{
			name: "Field with invalid value - array_min_items",
			data: `
			{
				"first_name": "John",
				"last_name": "Doe",
				"age": 20,
				"proof_of_age": {
					"id_type": "ssn",
					"id_number": "A8693170"
				},
				"parents": []
			}
					`,
			want: Error{
				SchemaValidationErrors: []SchemaValidationError{
					{
						DataPath: "$.parents",
						Message:  "Array must have at least 1 item",
					},
				},
			},
		},
		{
			name: "Field with invalid type",
			data: `
			{
				"first_name": "John",
				"last_name": "Doe",
				"age": "20",
				"proof_of_age": {
					"id_type": "ssn",
					"id_number": "A8693170"
				}
			}
					`,
			want: Error{
				SchemaValidationErrors: []SchemaValidationError{
					{
						DataPath: "$.age",
						Message:  "Invalid type. Expected: integer, given: string",
					},
				},
			},
		},
		{
			name: "Field with invalid enum value",
			data: `
			{
				"first_name": "John",
				"last_name": "Doe",
				"age": 20,
				"proof_of_age": {
					"id_type": "ssn",
					"id_number": "A8693170"
				},
				"title": "Engr"
			}
					`,
			want: Error{
				SchemaValidationErrors: []SchemaValidationError{
					{
						DataPath: "$.title",
						Message:  `Must be one of the following: "Dr", "Prof", "Mr", "Mrs", "Ms"`,
					},
				},
			},
		},
		{
			name: "Field with validation error in array",
			data: `
			{
				"first_name": "John",
				"last_name": "Doe",
				"age": 20,
				"proof_of_age": {
					"id_type": "ssn",
					"id_number": "A8693170"
				},
				"parents": [
					"24a78917-da7d-42f1-b1c8-8f9895f48e54",
					"invalid-uuid"
				]
			}
					`,
			want: Error{
				SchemaValidationErrors: []SchemaValidationError{
					{
						DataPath: "$.parents[1]",
						Message:  "Does not match format 'uuid'",
					},
				},
			},
		},
		{
			name: "Field with validation error in sub property",
			data: `
			{
				"first_name": "John",
				"last_name": "Doe",
				"age": 20,
				"proof_of_age": {
					"id_type": "ssn",
					"id_number": "A8693170"
				},
				"primary_contact": {
					"phone": "1"
				}
			}
					`,
			want: Error{
				SchemaValidationErrors: []SchemaValidationError{
					{
						DataPath: "$.primary_contact.phone",
						Message:  "String length must be greater than or equal to 2",
					},
				},
			},
		},
		{
			name: "Field with validation error in sub property array",
			data: `
			{
				"first_name": "John",
				"last_name": "Doe",
				"age": 20,
				"proof_of_age": {
					"id_type": "ssn",
					"id_number": "A8693170"
				},
				"primary_contact": {
					"email": []
				}
			}
					`,
			want: Error{
				SchemaValidationErrors: []SchemaValidationError{
					{
						DataPath: "$.primary_contact.email",
						Message:  "Array must have at least 1 item",
					},
				},
			},
		},
		{
			name: "Field with validation error in sub property array with invalid value",
			data: `
			{
				"first_name": "John",
				"last_name": "Doe",
				"age": 20,
				"proof_of_age": {
					"id_type": "ssn",
					"id_number": "A8693170"
				},
				"primary_contact": {
					"email": [
						"test@example.com",
						"notavalidemail"
					]
				}
			}
					`,
			want: Error{
				SchemaValidationErrors: []SchemaValidationError{
					{
						DataPath: "$.primary_contact.email[1]",
						Message:  "Does not match format 'email'",
					},
				},
			},
		},
	}

	for i := range tests {
		tests[i].want.HttpStatusCode = http.StatusBadRequest
		tests[i].want.ErrorCode = ErrorCodeValidationSchema
		tests[i].want.Message = registry[ErrorCodeValidationSchema].Message
	}

	cwd, err := os.Getwd()
	require.NoError(t, err, "Getwd returns without error")

	schema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			LocationRoot: cwd,
			Location:     "testdata",
			Name:         "test.main.schema.json",
		},
		References: []jsonschema.Schema{
			{
				LocationRoot: cwd,
				Location:     "testdata",
				Name:         "test.reference.schema.json",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := jsonschema.Validate(schema, []byte(tt.data))
			require.NoError(t, err, "should not return an error")

			schemaResultErrors := result.Errors()
			require.NotEmpty(t, schemaResultErrors, "schema validation should return errors")

			sve := NewSchemaValidationError(schemaResultErrors)
			require.Equal(t, tt.want, sve, "schema validation error should be as expected")
		})
	}
}
