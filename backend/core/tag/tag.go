package tag

import (
	"database/sql"
	"reflect"
	"time"
)

var nonRecurseTypeNames = map[interface{}]struct{}{
	getTypeName(sql.NullInt16{}):  {},
	getTypeName(sql.NullString{}): {},
	getTypeName(sql.NullTime{}):   {},
	getTypeName(sql.NullBool{}):   {},
	getTypeName(time.Time{}):      {},
}

func getTypeName(i interface{}) string {
	return reflect.TypeOf(i).Name()
}

func GetValues(entity interface{}, tag string, fp ...ExcludeFilterPredicate) []string {
	entityType := reflect.TypeOf(entity)

	return getValuesRecur(entityType, tag, fp...)
}

func getValuesRecur(t reflect.Type, tag string, fp ...ExcludeFilterPredicate) []string {
	var tagValues []string

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if shouldExclude(field, fp...) {
			continue
		}

		if shouldRecurse(field) {
			tagValues = append(tagValues, getValuesRecur(field.Type, tag)...)
		} else if attr, ok := shouldGetTagValue(field, tag); ok {
			tagValues = append(tagValues, attr)
		}
	}

	return tagValues
}

func shouldRecurse(field reflect.StructField) bool {
	return field.Type.Kind() == reflect.Struct && isRecurseFieldType(field.Type.Name())
}

func isRecurseFieldType(name string) bool {
	_, ok := nonRecurseTypeNames[name]
	return !ok
}

func shouldGetTagValue(field reflect.StructField, tag string) (string, bool) {
	return field.Tag.Lookup(tag)
}

func shouldExclude(field reflect.StructField, fp ...ExcludeFilterPredicate) bool {
	for _, exclude := range fp {
		if exclude(field) {
			return true
		}
	}

	return false
}
