package revuuid

import (
	"math/rand/v2"
	"sort"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestNewV7Reverse(t *testing.T) {
	u1 := uuid.Must(NewV7Reverse()).String()
	time.Sleep(10 * time.Millisecond)
	u2 := uuid.Must(NewV7Reverse()).String()
	time.Sleep(10 * time.Millisecond)
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
