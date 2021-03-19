package pkg

import "strings"

var (
	//nolint
	sls = strings.NewReplacer("\n", " ", "  ", " ")
)

func SingleLineString(s string) string {
	return sls.Replace(s)
}
