package formatter

import (
	"regexp"
	"strings"
)

type Unfold struct {
}

var newlineRegexp = regexp.MustCompile(`(\\r)?\\n`)

func (_ Unfold) Handle(raw string) string {
	return strings.Replace(
		newlineRegexp.ReplaceAllString(raw, "\r\n"),
		`\t`, "\t", -1,
	)
}
