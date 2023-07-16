package runner

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunDaemon(t *testing.T) {

	c, l, s, err := newDefaultDependencies()
	require.NoError(t, err, "newDefaultDependencies returns without error")

	r, err := newTestRunner(c, l, s)
	require.NoError(t, err, "newTestRunner returns without error")
	require.NotNil(t, r, "newTestRunner returns a Runner")

	err = r.RunDaemon(map[string]interface{}{ArgCycleLimit: 1})
	require.NoError(t, err, "Runner runs once without error")
}
