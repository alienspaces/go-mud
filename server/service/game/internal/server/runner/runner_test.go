package runner

import (
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/go-mud/server/service/game/internal/harness"
)

func NewTestHarness() (*harness.Testing, error) {

	// harness
	config := harness.DefaultDataConfig

	h, err := harness.NewTesting(config)
	if err != nil {
		return nil, err
	}

	// harness commit data
	h.CommitData = true

	return h, nil
}

func TestNewRunner(t *testing.T) {

	// Test harness
	th, err := NewTestHarness()
	require.NoError(t, err, "New test data returns without error")

	r := NewRunner()
	err = r.Init(th.Config, th.Log, th.Store)
	require.NoError(t, err, "Init returns without error")
}
