package tag

import (
	"reflect"
	"strings"

	"gitlab.com/alienspaces/go-mud/backend/core/collection/set"
)

func GetArrayFieldTagValues(entity any, tag string) set.Set[string] {
	t := reflect.TypeOf(entity)
	if t.Kind() != reflect.Struct {
		return nil
	}

	values := GetFieldTagValues(entity, tag, func(f reflect.StructField) bool {
		kind := f.Type.Kind()
		if (kind == reflect.Slice || kind == reflect.Array) && f.Type.Elem().Kind() == reflect.Uint8 {
			return true
		}
		return kind != reflect.Slice && kind != reflect.Array
	})

	return set.New(values...)
}

func GetReadonlyFieldTagValues(entity any, tag string) []string {
	values := GetFieldTagValues(entity, tag, func(f reflect.StructField) bool {
		if attr, ok := f.Tag.Lookup(tag); ok {
			if strings.Contains(attr, "readonly") {
				return false
			}
		}
		return true
	})
	return values
}
