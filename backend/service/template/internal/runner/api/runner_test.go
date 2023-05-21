package runner

import (
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/go-mud/backend/core/log"
	"gitlab.com/alienspaces/go-mud/backend/core/store"
	"gitlab.com/alienspaces/go-mud/backend/core/type/configurer"
	"gitlab.com/alienspaces/go-mud/backend/core/type/logger"
	"gitlab.com/alienspaces/go-mud/backend/core/type/storer"

	templateConfig "gitlab.com/alienspaces/go-mud/backend/service/template/internal/config"
	"gitlab.com/alienspaces/go-mud/backend/service/template/internal/harness"
	"gitlab.com/alienspaces/go-mud/backend/service/template/internal/record"
)

// newDefaultDependencies -
func newDefaultDependencies() (configurer.Configurer, logger.Logger, storer.Storer, error) {

	// configurer
	c, err := templateConfig.NewConfig(nil, false)
	if err != nil {
		return nil, nil, nil, err
	}

	// logger
	l, err := log.NewLogger(c)
	if err != nil {
		return nil, nil, nil, err
	}

	// storer
	s, err := store.NewStore(c, l)
	if err != nil {
		return nil, nil, nil, err
	}

	return c, l, s, nil
}

func newTestHarness(c configurer.Configurer, l logger.Logger, s storer.Storer, config func() harness.DataConfig) (*harness.Testing, error) {
	if config == nil {
		config = func() harness.DataConfig {
			return harness.DataConfig{
				TemplateConfig: []harness.TemplateConfig{
					{
						Record: &record.Template{},
					},
				},
			}
		}
	}

	h, err := harness.NewTesting(c, l, s, config(), true)
	if err != nil {
		return nil, err
	}

	return h, nil
}

func newTestRunner(c configurer.Configurer, l logger.Logger, s storer.Storer) (*Runner, error) {
	r, err := NewRunner(c, l)
	if err != nil {
		return nil, err
	}

	if err = r.Init(s); err != nil {
		return nil, err
	}

	return r, nil
}

func TestNewRunner(t *testing.T) {

	c, l, s, err := newDefaultDependencies()
	require.NoError(t, err, "newDefaultDependencies returns without error")

	r, err := NewRunner(c, l)
	require.NoError(t, err, "NewRunner returns without error")

	err = r.Init(s)
	require.NoError(t, err, "Init returns without error")
}
