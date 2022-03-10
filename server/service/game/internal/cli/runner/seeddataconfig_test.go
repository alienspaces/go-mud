package runner

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_seedDataConfig(t *testing.T) {

	// Test harness
	config := SeedDataConfig()

	th, err := NewTestHarness(&config)
	require.NoError(t, err, "New test harness returns without error")

	defer func() {
		th.Teardown()
	}()

	err = th.Setup()
	require.NoError(t, err, "Seed data setup returns without error")
}
