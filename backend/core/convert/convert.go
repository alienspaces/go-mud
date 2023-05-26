package convert

import (
	"fmt"
	"math"
	"strconv"

	"github.com/lib/pq"
)

var ordinalDictionary = map[int]string{
	0: "th",
	1: "st",
	2: "nd",
	3: "rd",
	4: "th",
	5: "th",
	6: "th",
	7: "th",
	8: "th",
	9: "th",
}

var shortToLongDays = map[string]string{
	"mon": "Monday",
	"tue": "Tuesday",
	"wed": "Wednesday",
	"thu": "Thursday",
	"fri": "Friday",
	"sat": "Saturday",
	"sun": "Sunday",
}

func Intp(i int) *int {
	return &i
}

func Int16p(i int16) *int16 {
	return &i
}

func Int32p(i int32) *int32 {
	return &i
}

func Int64p(i int64) *int64 {
	return &i
}

func Boolp(b bool) *bool {
	return &b
}

func Bool(b *bool) bool {
	return b != nil && *b
}

func Stringp(s string) *string {
	return &s
}

func PqStringArrayToStrSlice(a pq.StringArray) []string {
	return a
}

func Ordinalize(n int) string {
	n = int(math.Abs(float64(n)))

	if ((n % 100) >= 11) && ((n % 100) <= 13) {
		return strconv.Itoa(n) + "th"
	}

	return strconv.Itoa(n) + ordinalDictionary[n%10]
}

func ShortToLongWeekdays(shortDay string) (string, error) {
	if val, ok := shortToLongDays[shortDay]; ok {
		return val, nil
	}
	return "", fmt.Errorf("invalid shortday string >%s<", shortDay)
}
