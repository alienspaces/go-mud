package runner

import (
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/dependencies"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/harness"
)

func NewTestHarness(config *harness.DataConfig) (*harness.Testing, error) {

	// harness
	if config == nil {
		config = &harness.DefaultDataConfig
	}

	c, l, s, err := dependencies.Default()
	if err != nil {
		return nil, err
	}

	h, err := harness.NewTesting(c, l, s, *config)
	if err != nil {
		return nil, err
	}

	// harness commit data
	h.ShouldCommitData = true

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
