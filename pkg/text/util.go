package text

import "strings"

func IsBlankString(value string) bool {
	return len(strings.TrimSpace(value)) == 0
}
