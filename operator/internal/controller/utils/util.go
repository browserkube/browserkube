package utils

// FirstNonEmpty first non-empty string
func FirstNonEmpty(strs ...string) string {
	for _, s := range strs {
		if s != "" {
			return s
		}
	}
	return ""
}
