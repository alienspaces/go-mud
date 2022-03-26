package jsonschema

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCompile(t *testing.T) {
	cwd, err := os.Getwd()
	require.NoError(t, err, "Getwd returns without error")

	schema := SchemaWithReferences{
		Main: Schema{
			LocationRoot: cwd,
			Location:     "testdata",
			Name:         "test.main.schema.json",
		},
		References: []Schema{
			{
				LocationRoot: cwd,
				Location:     "testdata",
				Name:         "test.data.schema.json",
			},
		},
	}

	s, err := Compile(schema)
	require.NoError(t, err, "Compile returns without error")
	require.NotNil(t, s, "Compile returns a compiled schema")

	schema.References = append(schema.References, Schema{
		LocationRoot: cwd,
		Location:     "testdata",
		Name:         "test.missing.schema.json",
	})

	s, err = Compile(schema)
	require.Error(t, err, "Compile returns with error")
	require.Nil(t, s, "Compile does not return a compiled schema with error")
}

func BenchmarkCompile(b *testing.B) {
	cwd, err := os.Getwd()
	require.NoError(b, err, "Getpwd returns without error")

	schema := SchemaWithReferences{
		Main: Schema{
			LocationRoot: cwd,
			Location:     "testdata",
			Name:         "test.main.schema.json",
		},
		References: []Schema{
			{
				LocationRoot: cwd,
				Location:     "testdata",
				Name:         "test.data.schema.json",
			},
		},
	}

	for n := 0; n < b.N; n++ {
		Compile(schema)
	}
}

func Benchmark_compile(b *testing.B) {
	cwd, err := os.Getwd()
	require.NoError(b, err, "Getpwd returns without error")

	schema := SchemaWithReferences{
		Main: Schema{
			LocationRoot: cwd,
			Location:     "testdata",
			Name:         "test.main.schema.json",
		},
		References: []Schema{
			{
				LocationRoot: cwd,
				Location:     "testdata",
				Name:         "test.data.schema.json",
			},
		},
	}

	for n := 0; n < b.N; n++ {
		compile(schema)
	}
}

func TestValidate(t *testing.T) {
	//
}
