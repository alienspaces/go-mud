package runner

import (
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/go-mud/server/service/template/internal/dependencies"
	"gitlab.com/alienspaces/go-mud/server/service/template/internal/harness"
	"gitlab.com/alienspaces/go-mud/server/service/template/internal/record"
)

func NewTestHarness() (*harness.Testing, error) {

	// harness
	config := harness.DataConfig{
		TemplateConfig: []harness.TemplateConfig{
			{
				Record: record.Template{},
			},
		},
	}

	c, l, s, err := dependencies.Default()
	if err != nil {
		return nil, err
	}

	h, err := harness.NewTesting(c, l, s, config)
	if err != nil {
		return nil, err
	}

	// harness commit data
	h.CommitData = true

	return h, nil
}

func TestNewRunner(t *testing.T) {

	c, l, s, err := dependencies.Default()
	require.NoError(t, err, "NewDefaultDependencies returns without error")

	rnr, err := NewRunner(c, l)
	require.NoError(t, err, "NewRunner returns without error")

	err = rnr.Init(s)
	require.NoError(t, err, "Init returns without error")
}
