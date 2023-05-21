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

type Schema struct {
	LocationRoot string
	Location     string
	Name         string
}

type cacheKey string

var schemaCache = map[cacheKey]*gojsonschema.Schema{}
var mu = sync.Mutex{}

func (s SchemaWithReferences) IsEmpty() bool {
	return s.Main.Name == "" || s.Main.Location == ""
}

func (s SchemaWithReferences) GetReferencesFullPaths() []string {
	var paths []string

	for _, r := range s.References {
		paths = append(paths, r.GetFullPath())
	}

	return paths
}

func (s Schema) GetLocation() string {
	if s.LocationRoot == "" || s.Location == "" {
		return ""
	}

	return fmt.Sprintf("%s/%s", s.LocationRoot, s.Location)
}

func (s Schema) GetFullPath() string {
	loc := s.GetLocation()
	if loc == "" || s.Name == "" {
		return ""
	}

	return fmt.Sprintf("%s/%s", loc, s.Name)
}

func ValidateJSON(schema SchemaWithReferences, document []byte) error {
	result, err := Validate(schema, document)
	if err != nil {
		return err
	}
	return MapError(result)
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

	return s.Validate(dataLoader)
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
		refs = append(refs, r.GetFullPath())
	}

	key := s.Main.GetFullPath() + strings.Join(refs, "-")
	return cacheKey(key)
}

// Internal non-caching JSON schema compilation
func compile(sr SchemaWithReferences) (*gojsonschema.Schema, error) {

	sl := gojsonschema.NewSchemaLoader()
	sl.Validate = true

	for _, ref := range sr.References {
		refPath := fmt.Sprintf("file://%s", ref.GetFullPath())
		loader := gojsonschema.NewReferenceLoader(refPath)
		err := sl.AddSchemas(loader)
		if err != nil {
			return nil, fmt.Errorf("failed adding reference schema >%s< err >%w<", refPath, err)
		}
	}

	mainPath := fmt.Sprintf("file://%s", sr.Main.GetFullPath())
	loader := gojsonschema.NewReferenceLoader(mainPath)
	s, err := sl.Compile(loader)
	if err != nil {
		return nil, fmt.Errorf("failed adding main schema >%s< err >%w<, are you sure you've loaded all required reference schemas?", mainPath, err)
	}

	return s, nil
}

func ResolveSchemaLocationRoot(cfg SchemaWithReferences, root string) SchemaWithReferences {
	cfg.Main.LocationRoot = resolveString(cfg.Main.LocationRoot, root)

	for i := range cfg.References {
		cfg.References[i].LocationRoot = resolveString(cfg.References[i].LocationRoot, root)
	}

	return cfg
}

func ResolveSchemaLocation(cfg SchemaWithReferences, location string) SchemaWithReferences {
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
