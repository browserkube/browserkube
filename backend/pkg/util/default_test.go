package browserkubeutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Checking if the first string is found. Gotta find the string "first_string"
func TestFirstNonEmptyTrue(t *testing.T) {
	str := FirstNonEmpty("", "first_string", "")
	assert.Equal(t, "first_string", str)
}

// Checking if the first string is found. Gotta find nothing
func TestFirstNonEmptyFalse(t *testing.T) {
	str := FirstNonEmpty("", "", "")
	assert.Equal(t, "", str)
}
