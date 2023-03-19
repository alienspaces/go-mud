package convert

import "github.com/lib/pq"

func Boolp(b bool) *bool {
	return &b
}

func Stringp(s string) *string {
	return &s
}

func Int16p(i int16) *int16 {
	return &i
}

func PqStringArrayToStrSlice(a pq.StringArray) []string {
	return a
}
