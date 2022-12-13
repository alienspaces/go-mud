package tag

import (
	"reflect"
)

type ExcludeFilterPredicate func(f reflect.StructField) bool

func ExcludeStructTypes(types map[interface{}]struct{}) ExcludeFilterPredicate {
	return func(f reflect.StructField) bool {
		name := f.Type.Name()

		return f.Type.Kind() == reflect.Struct && isExcludedTypeName(name, types)
	}
}

func isExcludedTypeName(name string, excludedTypes map[interface{}]struct{}) bool {
	_, ok := excludedTypes[name]
	return ok
}

func ExcludeTagValues(values map[string]struct{}, tag string) ExcludeFilterPredicate {
	return func(f reflect.StructField) bool {
		value, ok := f.Tag.Lookup(tag)
		if !ok {
			return false
		}

		return isExcludedTagValue(value, values)
	}
}

func isExcludedTagValue(value string, excludedValues map[string]struct{}) bool {
	_, ok := excludedValues[value]
	return ok
}