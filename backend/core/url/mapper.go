package url

import (
	"regexp"
	"strings"
)

var reInitialNonWord = regexp.MustCompile(`^\W+`)
var reTrailingNonWord = regexp.MustCompile(`\W+$`)
var reQuotes = regexp.MustCompile(`[\'\"\\']`)
var reNonWord = regexp.MustCompile(`\W+`)

func Slugify(s string) string {
	// Replace any initial non-word character with nothing
	s = reInitialNonWord.ReplaceAllString(s, `$1`)

	// Replace any trailing non-word character with nothing
	s = reTrailingNonWord.ReplaceAllString(s, `$1`)

	// Replace any quotes with nothing
	s = reQuotes.ReplaceAllString(s, `$1`)

	// Replace any remaining non-word character with -
	s = reNonWord.ReplaceAllString(s, `$1-`)

	s = strings.ToLower(s)

	return s
}
