package harness

import (
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/go-mud/server/service/game/internal/dependencies"
)

func TestSetupTeardown(t *testing.T) {

	// harness
	config := DefaultDataConfig

	c, l, s, m, err := dependencies.Default()
	require.NoError(t, err, "NewTesting returns without error")

	h, err := NewTesting(c, l, s, m, config)
	require.NoError(t, err, "NewTesting returns without error")

	// harness commit data
	h.CommitData = true

	func() {
		err = h.Setup()
		require.NoError(t, err, "Setup returns without error")
		defer func() {
			err = h.Teardown()
			require.NoError(t, err, "Teardown returns without error")
		}()
	}()
}
