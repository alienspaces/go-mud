package jsonschema

import (
	"fmt"
	"strings"
	"sync"

	"github.com/xeipuuv/gojsonschema"
)

type SchemaWithReferences struct {
	Main       Schema
	References []Schema
}

// Schema defines the file path of a JSON schema
type Schema struct {
	// LocationRoot is the path from the file systems root director
	LocationRoot string
	// Location is the relative path from the root location
	Location string
	// Name is the file name
	Name string
}

type cacheKey string

var schemaCache = map[cacheKey]*gojsonschema.Schema{}
var mu = sync.Mutex{}

func (s SchemaWithReferences) IsEmpty() bool {
	return s.Main.Name == "" || s.Main.Location == ""
}

func (s Schema) GetFilePath() string {
	parts := []string{}
	if s.LocationRoot != "" {
		parts = append(parts, s.LocationRoot)
	}
	if s.Location != "" {
		parts = append(parts, s.Location)
	}
	return strings.Join(parts, "/")
}

func (s Schema) GetFileName() string {
	loc := s.GetFilePath()

	parts := []string{}
	if loc != "" {
		parts = append(parts, loc)
	}
	if s.Name != "" {
		parts = append(parts, s.Name)
	}
	return strings.Join(parts, "/")
}

func ValidateJSON(schema SchemaWithReferences) func([]byte) error {
	return func(document []byte) error {
		result, err := Validate(schema, document)
		if err != nil {
			return err
		}

		return MapError(result)
	}
}

func Validate(schema SchemaWithReferences, data interface{}) (*gojsonschema.Result, error) {

	s, err := Compile(schema)
	if err != nil {
		return nil, err
	}

	var dataLoader gojsonschema.JSONLoader
	switch d := data.(type) {
	case nil:
		return nil, fmt.Errorf("data is nil")
	case []byte:
		dataLoader = gojsonschema.NewStringLoader(string(d[:]))
	case string:
		dataLoader = gojsonschema.NewStringLoader(d)
	default:
		dataLoader = gojsonschema.NewGoLoader(d)
	}

	result, err := s.Validate(dataLoader)
	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, fmt.Errorf("failed validation, result is nil")
	}

	if !result.Valid() {
		return result, fmt.Errorf("failed validation >%+v<", result.Errors())
	}

	return result, nil
}

func MapError(result *gojsonschema.Result) error {
	if result.Valid() {
		return nil
	}

	var errStr string

	for _, e := range result.Errors() {
		if errStr == "" {
			errStr = e.String()
			continue
		}
		errStr = fmt.Sprintf("%s, %s", errStr, e.String())
	}

	return fmt.Errorf(errStr)
}

// Compile caches JSON schema compilation
func Compile(sr SchemaWithReferences) (*gojsonschema.Schema, error) {

	key := generateCacheKey(sr)
	cached, ok := schemaCache[key]
	if !ok {
		mu.Lock()
		defer mu.Unlock()
		if cached, ok = schemaCache[key]; ok {
			return cached, nil
		}
	} else {
		return cached, nil
	}

	s, err := compile(sr)
	if err != nil {
		return nil, err
	}

	schemaCache[key] = s

	return s, nil
}

func generateCacheKey(s SchemaWithReferences) cacheKey {
	var refs []string
	for _, r := range s.References {
		refs = append(refs, r.GetFileName())
	}

	key := s.Main.GetFileName() + strings.Join(refs, "-")
	return cacheKey(key)
}

// Internal non-caching JSON schema compilation
func compile(sr SchemaWithReferences) (*gojsonschema.Schema, error) {

	sl := gojsonschema.NewSchemaLoader()
	sl.Validate = true

	for _, ref := range sr.References {
		refPath := fmt.Sprintf("file://%s", ref.GetFileName())
		loader := gojsonschema.NewReferenceLoader(refPath)
		err := sl.AddSchemas(loader)
		if err != nil {
			return nil, fmt.Errorf("failed adding reference schema >%s< err >%w<", refPath, err)
		}
	}

	mainPath := fmt.Sprintf("file://%s", sr.Main.GetFileName())
	loader := gojsonschema.NewReferenceLoader(mainPath)
	s, err := sl.Compile(loader)
	if err != nil {
		return nil, fmt.Errorf("failed adding main schema >%s< err >%w<", mainPath, err)
	}

	return s, nil
}

func ResolveSchemaLocationRoot(root string, cfg SchemaWithReferences) SchemaWithReferences {
	cfg.Main.LocationRoot = resolveString(cfg.Main.LocationRoot, root)
	for i := range cfg.References {
		cfg.References[i].LocationRoot = resolveString(cfg.References[i].LocationRoot, root)
	}
	return cfg
}

func ResolveSchemaLocation(location string, cfg SchemaWithReferences) SchemaWithReferences {
	cfg.Main.Location = resolveString(cfg.Main.Location, location)
	for i := range cfg.References {
		cfg.References[i].Location = resolveString(cfg.References[i].Location, location)
	}
	return cfg
}

func resolveString(str string, defaultStr string) string {
	if str != "" {
		return str
	}
	return defaultStr
}
