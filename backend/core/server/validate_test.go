package server

import (
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	coreerror "gitlab.com/alienspaces/go-mud/backend/core/error"
	"gitlab.com/alienspaces/go-mud/backend/core/jsonschema"
)

func Test_validateQueryParameters(t *testing.T) {
	_, l, _, err := NewDefaultDependencies()
	require.NoError(t, err, "NewDefaultDependencies should not return err")

	type args struct {
		q           url.Values
		paramSchema jsonschema.SchemaWithReferences
	}

	type testcase struct {
		name    string
		args    args
		want    map[string]interface{}
		errCode coreerror.Code
	}
	tests := []testcase{
		{
			name: "nil",
			args: args{
				q: nil,
			},
		},
		{
			name: "id",
			args: args{
				q: url.Values{
					"id": []string{"a87feca8-d6f0-4794-98c7-037b30219520"},
				},
			},
		},
		{
			name: "string",
			args: args{
				q: url.Values{
					"string": []string{"asdf"},
				},
			},
		},
		{
			name: "number",
			args: args{
				q: url.Values{
					"number": []string{"123"},
				},
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
			},
		},
		{
			name: "string invalid empty",
			args: args{
				q: url.Values{
					"string": []string{""},
				},
			},
			errCode: coreerror.GetRegistryError(coreerror.SchemaValidation).ErrorCode,
		},
		{
			name: "number invalid below min",
			args: args{
				q: url.Values{
					"number": []string{"0"},
				},
			},
			errCode: coreerror.GetRegistryError(coreerror.SchemaValidation).ErrorCode,
		},
		{
			name: "number array",
			args: args{
				q: url.Values{
					"number": []string{"0", "1"},
				},
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
			},
			errCode: coreerror.GetRegistryError(coreerror.SchemaValidation).ErrorCode,
		},
		{
			name: "additional property",
			args: args{
				q: url.Values{
					"asdf": []string{"0"},
				},
			},
			errCode: coreerror.GetRegistryError(coreerror.SchemaValidation).ErrorCode,
		},
	}

	cwd, err := os.Getwd()
	require.NoError(t, err, "Getwd returns without error")

	for i := range tests {
		if tests[i].args.q != nil {
			tests[i].want = buildQueryParams(tests[i].args.q)
		}

		tests[i].args.paramSchema = jsonschema.SchemaWithReferences{
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
	}

	noSchemaTests := []testcase{
		{
			name: "query param with no schema",
			args: args{
				q: url.Values{
					"asdf": []string{"0"},
				},
			},
			errCode: coreerror.GetRegistryError(coreerror.InvalidQueryParam).ErrorCode,
		},
		{
			name: "no query param with no schema",
			want: buildQueryParams(nil),
		},
	}

	tests = append(tests, noSchemaTests...)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := validateQueryParameters(l, tt.args.q, tt.args.paramSchema)
			if tt.errCode != "" {
				require.NotNil(t, err, "validateQueryParameters should return err")
				coreerrorErr, conversionErr := coreerror.ToError(err)
				require.Nil(t, conversionErr, "should not have an err that is not wrapped")

				require.Equal(t, tt.errCode, coreerrorErr.ErrorCode)

				e := coreerror.ProcessQueryParamError(err)
				coreerrorErr, conversionErr = coreerror.ToError(e)
				require.Nil(t, conversionErr, "should not have an err that is not wrapped")

				require.Equal(t, coreerror.GetRegistryError(coreerror.InvalidQueryParam).ErrorCode, coreerrorErr.ErrorCode)
				return
			}

			require.Nil(t, err, "validateQueryParameters should not return err")
			require.Equal(t, tt.want, got)
		})
	}
}

func Test_queryParamsToJSON(t *testing.T) {
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
			name: "array value",
			args: args{
				q: url.Values{
					"str": []string{"a", "az"},
				},
			},
			want: `{"str":["a","az"]}`,
		},
		{
			name: "array with mixed types",
			args: args{
				q: url.Values{
					"mixed": []string{"123", "az"},
				},
			},
			want: `{"mixed":[123,"az"]}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := queryParamsToJSON(tt.args.q)
			require.Equal(t, tt.want, got)
		})
	}
}
