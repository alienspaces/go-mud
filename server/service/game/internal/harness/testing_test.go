package harness

// NOTE: repository tests are run is the public space so we are
// able to use common setup and teardown tooling for all repositories

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateOne(t *testing.T) {

	// harness
	config := DefaultDataConfig

	h, err := NewTesting(config)
	require.NoError(t, err, "NewTesting returns without error")

	// harness commit data
	h.CommitData = true

	func() {

		// Test harness
		err = h.Setup()
		require.NoError(t, err, "Setup returns without error")
		defer func() {
			err = h.Teardown()
			require.NoError(t, err, "Teardown returns without error")
		}()
	}()
}
