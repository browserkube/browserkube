package utils

import "fmt"

const readinessConfigmapKeyTpl = "%s.%s"

// CreateReadinessConfigmapKey creates key by given prefix and value
func CreateReadinessConfigmapKey(prefix, key string) string {
	return fmt.Sprintf(readinessConfigmapKeyTpl, prefix, key)
}
