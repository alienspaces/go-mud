package runner

import (
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/dependencies"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/harness"
)

func NewTestHarness() (*harness.Testing, error) {

	// Default test harness data configuration
	config := harness.DefaultDataConfig

	// Default dependencies
	c, l, s, err := dependencies.Default()
	if err != nil {
		return nil, err
	}

	// Test harness
	h, err := harness.NewTesting(c, l, s, config)
	if err != nil {
		return nil, err
	}

	// For handler testst the test harness needs to commit data as the handler
	// creates a new database transaction. Data created as a result of a test
	// case bust be specifically added to the test harness to be torn down.
	h.CommitData = true

	return h, nil
}

func TestNewRunner(t *testing.T) {

	c, l, s, err := dependencies.Default()
	require.NoError(t, err, "NewDefaultDependencies returns without error")

	r, err := NewRunner(c, l)
	require.NoError(t, err, "NewRunner returns without error")

	err = r.Init(s)
	require.NoError(t, err, "Init returns without error")
}
