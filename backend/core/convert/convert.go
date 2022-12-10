package convert

import "github.com/lib/pq"

func Boolp(b bool) *bool {
	return &b
}

func Stringp(s string) *string {
	return &s
}

func PqStringArrayToStrSlice(a pq.StringArray) []string {
	return a
}
