package repository

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewRecordID(t *testing.T) {
	ids := map[string]struct{}{}

	for range make([]int, 10000) {
		id := NewRecordID()

		_, ok := ids[id]
		ids[id] = struct{}{}

		require.False(t, ok, "Record ID should be unique")

		uuidLength := 36
		require.NotEmpty(t, id)
		require.Equal(t, uuidLength, len(id))
	}
}
