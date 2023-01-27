package test

// NOTE: query tests are run is the public space so we are
// able to use common setup and teardown tooling for all queries

import (
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/dependencies"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/harness"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/model"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/query/dungeonentityinstanceturn"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/record"
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
		expectRecs func(d harness.Data) []*record.DungeonEntityInstanceTurn
		expectErr  bool
	}{
		{
			name: "With no params",
			params: func(d harness.Data) map[string]interface{} {
				return nil
			},
			// TODO: 8-implement-turns: Complete query test
			expectRecs: func(d harness.Data) []*record.DungeonEntityInstanceTurn {
				dir, _ := d.GetDungeonInstanceRecByName(harness.DungeonNameCave)
				cir, _ := d.GetCharacterInstanceRecByName(harness.CharacterNameBarricade)
				mir, _ := d.GetMonsterInstanceRecByName(harness.MonsterNameGrumpyDwarf)
				return []*record.DungeonEntityInstanceTurn{
					&record.DungeonEntityInstanceTurn{
						DungeonInstanceID:         dir.ID,
						DungeonName:               harness.DungeonNameCave,
						DungeonInstanceTurnNumber: 2,
						EntityType:                string(model.EntityTypeCharacter),
						EntityInstanceID:          cir.ID,
						EntityName:                harness.CharacterNameBarricade,
						EntityInstanceTurnNumber:  1,
					},
					&record.DungeonEntityInstanceTurn{
						DungeonInstanceID:         dir.ID,
						DungeonName:               harness.DungeonNameCave,
						DungeonInstanceTurnNumber: 2,
						EntityType:                string(model.EntityTypeMonster),
						EntityInstanceID:          mir.ID,
						EntityName:                harness.MonsterNameGrumpyDwarf,
						EntityInstanceTurnNumber:  1,
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
			q, err := dungeonentityinstanceturn.NewQuery(l, h.QueryPreparer, h.Tx())
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
				require.Equal(t, len(expectRecs), len(recs), "GetMany returns the expected number of records")
				// Expected records and returned records may be in any order
				for idx := range recs {
					found := false
					rec := recs[idx]
					for eidx := range expectRecs {
						expectRec := expectRecs[eidx]
						if rec.DungeonInstanceID == expectRec.DungeonInstanceID &&
							rec.EntityInstanceID == expectRec.EntityInstanceID {
							found = true
							require.Equal(t, expectRec.DungeonName, harness.NormalName(rec.DungeonName), "DungeonName equals expected")
							require.Equal(t, expectRec.DungeonInstanceTurnNumber, rec.DungeonInstanceTurnNumber, "DungeonInstanceTurnNumber equals expected")
							require.Equal(t, expectRec.EntityName, harness.NormalName(rec.EntityName), "EntityName equals expected")
							require.Equal(t, expectRec.EntityInstanceTurnNumber, rec.EntityInstanceTurnNumber, "EntityInstanceTurnNumber equals expected")
						}
					}
					require.True(t, found, "Returned record with DungeonInstance ID >%s< and EntityInstance ID >%s< found", recs[idx].DungeonInstanceID, recs[idx].EntityInstanceID)
				}
			}

			h.RollbackTx()
		})
	}
}
