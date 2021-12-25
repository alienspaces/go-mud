package schema

import (
	"fmt"

	"github.com/xeipuuv/gojsonschema"
	"gitlab.com/alienspaces/go-mud/server/core/type/configurer"
	"gitlab.com/alienspaces/go-mud/server/core/type/logger"
)

// schemaCache - key, schema
var schemaCache map[string]*gojsonschema.Schema

type Config struct {
	Key        string
	Location   string
	Main       string
	References []string
}

type Validator struct {
	Config configurer.Configurer
	Log    logger.Logger
}

func NewValidator(c configurer.Configurer, l logger.Logger) (*Validator, error) {
	return &Validator{
		Config: c,
		Log:    l,
	}, nil
}

// ValidateBytes -
func (v *Validator) ValidateBytes(schemaConfig Config, data []byte) error {
	return v.Validate(schemaConfig, string(data))
}

// Validate -
func (v *Validator) Validate(schemaConfig Config, data string) error {

	v.Log.Warn("Validating key >%s< data >%s<", schemaConfig.Key, string(data))

	var err error

	s := schemaCache[schemaConfig.Key]
	if s == nil {
		s, err = v.LoadSchema(schemaConfig)
		if err != nil {
			v.Log.Warn("Failed validate >%v<", err)
			return err
		}
	}

	dataLoader := gojsonschema.NewStringLoader(string(data))

	result, err := s.Validate(dataLoader)
	if err != nil {
		v.Log.Warn("Failed validate >%v<", err)
		if err.Error() == "EOF" {
			v.Log.Warn("Data is empty")
			return err
		}
		return err
	}

	if !result.Valid() {
		errStr := ""
		for _, e := range result.Errors() {
			v.Log.Warn("Invalid data >%s<", e)
			if errStr == "" {
				errStr = e.String()
				continue
			}
			errStr = fmt.Sprintf("%s, %s", errStr, e.String())
		}
		return fmt.Errorf(errStr)
	}

	return nil
}

func (v *Validator) SchemaCached(key string) bool {

	if _, ok := schemaCache[key]; ok {
		return true
	}

	return false
}

func (v *Validator) LoadSchema(schemaConfig Config) (*gojsonschema.Schema, error) {

	v.Log.Warn("Loading schema key >%s< location >%s<", schemaConfig.Key, schemaConfig.Location)

	if schemaConfig.Location == "" {
		return nil, fmt.Errorf("missing Location, invalid config")
	}
	if schemaConfig.Main == "" {
		return nil, fmt.Errorf("missing Main, invalid config")
	}

	schemaLoc := schemaConfig.Location
	schema := schemaConfig.Main
	schemaReferences := schemaConfig.References

	schemaPath := v.Config.Get("APP_SERVER_SCHEMA_PATH")
	schemaLoc = fmt.Sprintf("file://%s/%s", schemaPath, schemaLoc)

	v.Log.Info("Loading schema %s/%s", schemaLoc, schema)

	sl := gojsonschema.NewSchemaLoader()
	sl.Validate = true

	for _, schemaName := range schemaReferences {
		v.Log.Info("Adding schema reference %s/%s", schemaLoc, schemaName)
		loader := gojsonschema.NewReferenceLoader(fmt.Sprintf("%s/%s", schemaLoc, schemaName))
		err := sl.AddSchemas(loader)
		if err != nil {
			v.Log.Warn("Failed adding schema reference %v", err)
			return nil, err
		}
	}

	loader := gojsonschema.NewReferenceLoader(fmt.Sprintf("%s/%s", schemaLoc, schema))
	s, err := sl.Compile(loader)
	if err != nil {
		v.Log.Warn("Failed compiling schema's >%v<", err)
		return nil, err
	}

	if schemaCache == nil {
		schemaCache = map[string]*gojsonschema.Schema{}
	}
	schemaCache[schemaConfig.Key] = s

	return s, nil
}
