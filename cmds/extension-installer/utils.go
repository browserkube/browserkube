package extensioninstaller

import (
	"strings"
)

func checkSlice(slice []string, item string) bool {
	for _, v := range slice {
		v = strings.TrimSuffix(v, "\n")
		if v == item {
			return true
		}
	}
	return false
}
