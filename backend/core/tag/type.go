package tag

import (
	"reflect"

	"gitlab.com/alienspaces/go-mud/backend/core/collection/set"
)

func GetArrayFieldTagValues(entity any, tag string) set.Set[string] {
	t := reflect.TypeOf(entity)
	if t.Kind() != reflect.Struct {
		return nil
	}

	values := GetFieldTagValues(entity, tag, func(f reflect.StructField) bool {
		kind := f.Type.Kind()

		// []byte
		if (kind == reflect.Slice || kind == reflect.Array) && f.Type.Elem().Kind() == reflect.Uint8 {
			return true
		}

		return kind != reflect.Slice && kind != reflect.Array
	})

	return set.New(values...)
}
