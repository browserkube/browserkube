package browserkubeutil

import (
	"io"

	"go.uber.org/zap"
)

// FirstNonEmpty first non-empty string
func FirstNonEmpty(strs ...string) string {
	for _, s := range strs {
		if s != "" {
			return s
		}
	}
	return ""
}

func CloseQuietly(c io.Closer) {
	if err := c.Close(); err != nil {
		zap.S().Errorf("unable to close: %+v", err)
	}
}
