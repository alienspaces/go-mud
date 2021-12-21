package repository

import (
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/dependencies"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/harness"
)

func NewTestHarness(config harness.DataConfig) (*harness.Testing, error) {
	c, l, s, m, err := dependencies.Default()
	if err != nil {
		return nil, err
	}

	h, err := harness.NewTesting(c, l, s, m, config)
	if err != nil {
		return nil, err
	}

	return h, nil
}
