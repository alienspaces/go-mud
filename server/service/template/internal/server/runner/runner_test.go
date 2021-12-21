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

	c, l, s, m, err := dependencies.Default()
	if err != nil {
		return nil, err
	}

	h, err := harness.NewTesting(c, l, s, m, config)
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
	err = r.Init(th.Config, th.Log, th.Store, th.Model)
	require.NoError(t, err, "Init returns without error")
}
