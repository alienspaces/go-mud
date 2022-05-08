package runner

import (
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/go-mud/server/service/game/internal/dependencies"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/harness"
)

func NewTestHarness(config *harness.DataConfig) (*harness.Testing, error) {

	// harness
	if config == nil {
		config = &harness.DefaultDataConfig
	}

	c, l, s, m, err := dependencies.Default()
	if err != nil {
		return nil, err
	}

	h, err := harness.NewTesting(c, l, s, m, *config)
	if err != nil {
		return nil, err
	}

	// harness commit data
	h.CommitData = true

	return h, nil
}

func TestNewRunner(t *testing.T) {

	// Test harness
	th, err := NewTestHarness(nil)
	require.NoError(t, err, "New test data returns without error")

	rnr := NewRunner()
	err = rnr.Init(th.Store)
	require.NoError(t, err, "Init returns without error")
}
