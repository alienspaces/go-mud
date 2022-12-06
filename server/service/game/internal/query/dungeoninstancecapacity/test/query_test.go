package test

// NOTE: query tests are run is the public space so we are
// able to use common setup and teardown tooling for all queries

import (
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/go-mud/server/service/game/internal/dependencies"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/harness"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/query/dungeoninstancecapacity"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

func TestGetOne(t *testing.T) {

	// harness
	config := harness.DefaultDataConfig

	c, l, s, err := dependencies.Default()
	require.NoError(t, err, "Default dependencies returns without error")

	h, err := harness.NewTesting(c, l, s, config)
	require.NoError(t, err, "NewTesting returns without error")

	// harness commit data
	h.CommitData = true

	tests := []struct {
		name       string
		params     func(d harness.Data) map[string]interface{}
		expectRecs func(d harness.Data) []*record.DungeonInstanceCapacity
		expectErr  bool
	}{
		{
			name: "With no params",
			params: func(d harness.Data) map[string]interface{} {
				return nil
			},
			expectRecs: func(d harness.Data) []*record.DungeonInstanceCapacity {
				return []*record.DungeonInstanceCapacity{
					{
						DungeonInstanceID:             d.DungeonInstanceRecs[0].ID,
						DungeonID:                     d.DungeonInstanceRecs[0].DungeonID,
						DungeonInstanceCharacterCount: 2,
						DungeonLocationCount:          6,
					},
				}
			},
			expectErr: false,
		},
	}

	// harness setup
	err = h.Setup()
	require.NoError(t, err, "Setup returns without error")
	defer func() {
		err = h.Teardown()
		require.NoError(t, err, "Teardown returns without error")
	}()

	for _, tc := range tests {

		t.Logf("Run test >%s<", tc.name)

		t.Run(tc.name, func(t *testing.T) {

			// init tx
			err = h.InitTx(nil)
			require.NoError(t, err, "InitTx returns without error")

			// query
			q, err := dungeoninstancecapacity.NewQuery(l, h.QueryPreparer, h.Tx())
			require.NoError(t, err, "NewQuery returns without error")
			require.NotNil(t, q, "Query is not nil")

			recs, err := q.GetMany(tc.params(h.Data), nil)
			if tc.expectErr == true {
				require.Error(t, err, "GetMany returns error")
				return
			}
			require.NoError(t, err, "GetMany returns without error")
			if tc.expectRecs != nil {
				expectRecs := tc.expectRecs(h.Data)
				require.Equal(t, expectRecs, recs, "GetMany returns expected records")
			}

			h.RollbackTx()
		})
	}
}
