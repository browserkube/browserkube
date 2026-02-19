//nolint:nonamedreturns
package revuuid

import (
	"encoding/binary"

	"github.com/google/uuid"
)

// This is slightly changed version of original
// Google's UUID-V7 generator
// This version replaces time.Now with (MaxInt64 - time.Now)
// to ensure descending order of generated UIDs

// Copyright 2023 Google Inc.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// UUID version 7 features a time-ordered value field derived from the widely
// implemented and well known Unix Epoch timestamp source,
// the number of milliseconds seconds since midnight 1 Jan 1970 UTC, leap seconds excluded.
// As well as improved entropy characteristics over versions 1 or 6.
//
// see https://datatracker.ietf.org/doc/html/draft-peabody-dispatch-new-uuid-format-03#name-uuid-version-7
//
// Implementations SHOULD utilize UUID version 7 over UUID version 1 and 6 if possible.
//
// NewV7 returns a Version 7 UUID based on the current time(Unix Epoch).
// Uses the randomness pool if it was enabled with EnableRandPool.
// On error, NewV7 returns Nil and an error
func NewV7Reverse() (uuid.UUID, error) {
	u, err := uuid.NewV7()
	if err != nil {
		return uuid.Nil, err
	}

	// Extract first 6 bytes (timestamp)
	ts := binary.BigEndian.Uint64(append([]byte{0, 0}, u[0:6]...))

	const max48 = (1 << 48) - 1
	inverted := max48 - ts

	// Put inverted timestamp back
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, inverted)

	copy(u[0:6], buf[2:8])

	return u, nil
}
