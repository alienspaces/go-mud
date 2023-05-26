package test

// NOTE: query tests are run is the public space so we are
// able to use common setup and teardown tooling for all queries

import (
	"testing"

	"github.com/stretchr/testify/require"

	coresql "gitlab.com/alienspaces/go-mud/backend/core/sql"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/dependencies"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/harness"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/model"
)

func TestGetMany(t *testing.T) {

	// harness
	config := harness.DefaultDataConfig

	c, l, s, err := dependencies.Default()
	require.NoError(t, err, "Default dependencies returns without error")

	h, err := harness.NewTesting(c, l, s, config)
	require.NoError(t, err, "NewTesting returns without error")

	// harness commit data
	h.ShouldCommitData = true

	tests := []struct {
		name      string
		options   func(d harness.Data) *coresql.Options
		expectErr bool
	}{
		{
			name: "With no params",
			options: func(d harness.Data) *coresql.Options {
				return nil
			},
			expectErr: false,
		},
	}

	// harness setup
	_, err = h.Setup()
	require.NoError(t, err, "Setup returns without error")
	defer func() {
		err = h.Teardown()
		require.NoError(t, err, "Teardown returns without error")
	}()

	for _, tc := range tests {

		t.Logf("Run test >%s<", tc.name)

		t.Run(tc.name, func(t *testing.T) {

			// init tx
			_, err = h.InitTx()
			require.NoError(t, err, "InitTx returns without error")

			// query
			q := h.Model.(*model.Model).DungeonEntityInstanceTurnQuery()
			require.NotNil(t, q, "Query is not nil")

			recs, err := q.GetMany(tc.options(h.Data))
			if tc.expectErr == true {
				require.Error(t, err, "GetMany returns error")
				return
			}

			require.NoError(t, err, "GetMany returns without error")

			for idx := range recs {
				require.NotEmpty(t, recs[idx].DungeonInstanceID, "DungeonInstanceID is not empty")
				require.NotEmpty(t, recs[idx].DungeonName, "DungeonName is not empty")
				require.NotEmpty(t, recs[idx].DungeonInstanceTurnNumber, "DungeonInstanceTurnNumber is not empty")
				require.NotEmpty(t, recs[idx].EntityType, "EntityType is not empty")
				require.NotEmpty(t, recs[idx].EntityInstanceID, "EntityInstanceID is not empty")
				require.NotEmpty(t, recs[idx].EntityName, "EntityName is not empty")
				require.NotEmpty(t, recs[idx].EntityInstanceTurnNumber, "EntityInstanceTurnNumber is not empty")
			}

			h.RollbackTx()
		})
	}
}
