package session

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnknownCaps(t *testing.T) {
	capsStr := `{"browserVersion":"102.0","browserkube:options":{"reportportal":{}},"options":{"option1":"option1Val"}}`

	var caps Capabilities
	err := json.Unmarshal([]byte(capsStr), &caps)
	require.NoError(t, err)

	res, err := json.Marshal(caps)
	require.NoError(t, err)
	require.Equal(t, capsStr, string(res))
}
