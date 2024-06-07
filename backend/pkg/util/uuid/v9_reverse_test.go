package revuuid

import (
	"sort"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"math/rand/v2"
)

func TestNewV7Reverse(t *testing.T) {
	u1 := uuid.Must(NewV7Reverse()).String()
	u2 := uuid.Must(NewV7Reverse()).String()
	u3 := uuid.Must(NewV7Reverse()).String()
	u := []string{u1, u2, u3}

	// shuffle
	rand.Shuffle(len(u), func(i, j int) {
		u[i], u[j] = u[j], u[i]
	})
	// sort
	sort.Strings(u)
	require.Equal(t, []string{u3, u2, u1}, u)
}
