package pkg

import "strings"

var (
	//nolint
	sls = strings.NewReplacer("\n", " ", "  ", " ")
)

func SingleLineString(s string) string {
	return strings.TrimSpace(sls.Replace(s))
}
