package server

import (
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	coreerror "gitlab.com/alienspaces/go-mud/backend/core/error"
	"gitlab.com/alienspaces/go-mud/backend/core/jsonschema"
)

func Test_extractPathParams(t *testing.T) {

	type testcase struct {
		name         string
		path         string
		expectParams []string
	}

	tests := []testcase{
		{
			name:         "two parameters with trailing parameter",
			path:         "/animals/:animal_id/humans/:human_id",
			expectParams: []string{"animal_id", "human_id"},
		},
		{
			name:         "two parameters without trailing parameter",
			path:         "/animals/:animal_id/humans/:human_id/overlords",
			expectParams: []string{"animal_id", "human_id"},
		},
		{
			name:         "one parameter without trailing parameter",
			path:         "/animals/:animal_id/humans",
			expectParams: []string{"animal_id"},
		},
		{
			name:         "no parameters",
			path:         "/animals",
			expectParams: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := extractPathParams(tt.path)
			require.Equal(t, tt.expectParams, params, "Extracted params equals expected")
		})
	}
}

func Test_validateQueryParameters(t *testing.T) {

	_, l, _, err := NewDefaultDependencies()
	require.NoError(t, err, "NewDefaultDependencies should not return err")

	type args struct {
		q           url.Values
		paramSchema *jsonschema.SchemaWithReferences
	}

	type testcase struct {
		name    string
		args    args
		errCode coreerror.ErrorCode
	}

	cwd, err := os.Getwd()
	require.NoError(t, err, "Getwd returns without error")

	paramSchema := &jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			LocationRoot: cwd,
			Location:     "testdata",
			Name:         "test.main.schema.json",
		},
		References: []jsonschema.Schema{
			{
				LocationRoot: cwd,
				Location:     "testdata",
				Name:         "test.data.schema.json",
			},
		},
	}

	tests := []testcase{
		{
			name: "nil",
			args: args{
				q:           nil,
				paramSchema: paramSchema,
			},
		},
		{
			name: "id",
			args: args{
				q: url.Values{
					"id": []string{"a87feca8-d6f0-4794-98c7-037b30219520"},
				},
				paramSchema: paramSchema,
			},
		},
		{
			name: "string",
			args: args{
				q: url.Values{
					"string": []string{"asdf"},
				},
				paramSchema: paramSchema,
			},
		},
		{
			name: "number",
			args: args{
				q: url.Values{
					"number": []string{"123"},
				},
				paramSchema: paramSchema,
			},
		},
		{
			name: "multiple",
			args: args{
				q: url.Values{
					"id":     []string{"a87feca8-d6f0-4794-98c7-037b30219520"},
					"string": []string{"asdf"},
					"number": []string{"123"},
				},
				paramSchema: paramSchema,
			},
		},
		{
			name: "string invalid empty",
			args: args{
				q: url.Values{
					"string": []string{""},
				},
				paramSchema: paramSchema,
			},
			errCode: coreerror.GetRegistryError(coreerror.SchemaValidation).ErrorCode,
		},
		{
			name: "number invalid below min",
			args: args{
				q: url.Values{
					"number": []string{"0"},
				},
				paramSchema: paramSchema,
			},
			errCode: coreerror.GetRegistryError(coreerror.SchemaValidation).ErrorCode,
		},
		{
			name: "number array",
			args: args{
				q: url.Values{
					"number": []string{"0", "1"},
				},
				paramSchema: paramSchema,
			},
			errCode: coreerror.GetRegistryError(coreerror.SchemaValidation).ErrorCode,
		},
		{
			name: "multiple",
			args: args{
				q: url.Values{
					"string": []string{""},
					"number": []string{"0", "1"},
				},
				paramSchema: paramSchema,
			},
			errCode: coreerror.GetRegistryError(coreerror.SchemaValidation).ErrorCode,
		},
		{
			name: "additional property",
			args: args{
				q: url.Values{
					"asdf": []string{"0"},
				},
				paramSchema: paramSchema,
			},
			errCode: coreerror.GetRegistryError(coreerror.SchemaValidation).ErrorCode,
		},
	}

	noSchemaTests := []testcase{
		{
			name: "query param with no schema",
			args: args{
				q: url.Values{
					"asdf": []string{"0"},
				},
			},
		},
	}

	tests = append(tests, noSchemaTests...)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err = validateParams(l, tt.args.q, tt.args.paramSchema)
			if tt.errCode != "" {
				require.Error(t, err, "validateQueryParameters should return err")
				coreerrorErr, conversionErr := coreerror.ToError(err)
				require.Nil(t, conversionErr, "should not have an err that is not wrapped")

				require.Equal(t, tt.errCode, coreerrorErr.ErrorCode)

				e := coreerror.ProcessParamError(err)
				coreerrorErr, conversionErr = coreerror.ToError(e)
				require.Nil(t, conversionErr, "should not have an err that is not wrapped")

				require.Equal(t, coreerror.GetRegistryError(coreerror.InvalidParam).ErrorCode, coreerrorErr.ErrorCode)
			} else {
				require.NoError(t, err, "validateParams should not return err")
			}
		})
	}
}

func Test_paramsToJSON(t *testing.T) {

	type args struct {
		q url.Values
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "nil",
			want: "",
		},
		{
			name: "empty",
			args: args{
				q: url.Values{},
			},
			want: "",
		},
		{
			name: "single empty value",
			args: args{
				q: url.Values{
					"1": []string{},
				},
			},
			want: `{"1":""}`,
		},
		{
			name: "single string value",
			args: args{
				q: url.Values{
					"a": []string{"a"},
				},
			},
			want: `{"a":"a"}`,
		},
		{
			name: "single int value",
			args: args{
				q: url.Values{
					"i": []string{"123"},
				},
			},
			want: `{"i":123}`,
		},
		{
			name: "single int value with operator",
			args: args{
				q: url.Values{
					"i:lte": []string{"123"},
				},
			},
			want: `{"i":123}`,
		},
		{
			name: "array - empty",
			args: args{
				q: url.Values{
					"str[]": []string{},
				},
			},
			want: `{"str[]":[]}`,
		},
		{
			name: "array - single value",
			args: args{
				q: url.Values{
					"str[]": []string{"a"},
				},
			},
			want: `{"str[]":["a"]}`,
		},
		{
			name: "array - multiple value",
			args: args{
				q: url.Values{
					"str[]": []string{"a", "az"},
				},
			},
			want: `{"str[]":["a","az"]}`,
		},
		{
			name: "array with mixed types",
			args: args{
				q: url.Values{
					"mixed[]": []string{"123", "az"},
				},
			},
			want: `{"mixed[]":[123,"az"]}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := paramsToJSON(tt.args.q)
			require.Equal(t, tt.want, got)
		})
	}
}
