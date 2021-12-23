package test

// NOTE: model tests are run is the public space to avoid cyclic dependencies

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/go-mud/server/service/game/internal/dependencies"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/harness"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/model"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

func TestProcessDungeonCharacterAction(t *testing.T) {

	// harness
	config := harness.DefaultDataConfig

	c, l, s, m, err := dependencies.Default()
	require.NoError(t, err, "NewTesting returns without error")

	th, err := harness.NewTesting(c, l, s, m, config)
	require.NoError(t, err, "NewTesting returns without error")

	// harness commit data
	th.CommitData = true

	tests := []struct {
		name                         string
		dungeonID                    func(data harness.Data) string
		dungeonCharacterID           func(data harness.Data) string
		sentence                     string
		expectDungeonActionRecordSet func(data harness.Data) *model.DungeonActionRecordSet
		expectError                  bool
	}{
		{
			name: "Command",
			dungeonID: func(data harness.Data) string {
				return data.DungeonRecs[0].ID
			},
			dungeonCharacterID: func(data harness.Data) string {
				return data.DungeonCharacterRecs[0].ID
			},
			sentence: "look",
			expectDungeonActionRecordSet: func(data harness.Data) *model.DungeonActionRecordSet {
				return &model.DungeonActionRecordSet{
					ActionRec: &record.DungeonAction{
						DungeonID:          data.DungeonRecs[0].ID,
						DungeonLocationID:  data.DungeonLocationRecs[0].ID,
						DungeonCharacterID: sql.NullString{String: data.DungeonCharacterRecs[0].ID, Valid: true},
					},
					ActionCharacterRec: &record.DungeonActionCharacter{
						Name: data.DungeonCharacterRecs[0].Name,
					},
					TargetLocation: &model.DungeonActionLocationRecordSet{
						LocationRec: &record.DungeonLocation{
							Name: data.DungeonLocationRecs[0].Name,
						},
					},
				}
			},
			expectError: false,
		},
		{
			name: "Command",
			dungeonID: func(data harness.Data) string {
				return data.DungeonRecs[0].ID
			},
			dungeonCharacterID: func(data harness.Data) string {
				return data.DungeonCharacterRecs[0].ID
			},
			sentence: "look north",
			expectDungeonActionRecordSet: func(data harness.Data) *model.DungeonActionRecordSet {
				return &model.DungeonActionRecordSet{
					ActionRec: &record.DungeonAction{
						DungeonID:          data.DungeonRecs[0].ID,
						DungeonLocationID:  data.DungeonLocationRecs[0].ID,
						DungeonCharacterID: sql.NullString{String: data.DungeonCharacterRecs[0].ID, Valid: true},
					},
					ActionCharacterRec: &record.DungeonActionCharacter{
						Name: data.DungeonCharacterRecs[0].Name,
					},
					TargetLocation: &model.DungeonActionLocationRecordSet{
						LocationRec: &record.DungeonLocation{
							Name: data.DungeonLocationRecs[1].Name,
						},
					},
				}
			},
			expectError: false,
		},
		{
			name: "Command",
			dungeonID: func(data harness.Data) string {
				return data.DungeonRecs[0].ID
			},
			dungeonCharacterID: func(data harness.Data) string {
				return data.DungeonCharacterRecs[0].ID
			},
			sentence: "look rusted sword",
			expectDungeonActionRecordSet: func(data harness.Data) *model.DungeonActionRecordSet {
				return &model.DungeonActionRecordSet{
					ActionRec: &record.DungeonAction{
						DungeonID:          data.DungeonRecs[0].ID,
						DungeonLocationID:  data.DungeonLocationRecs[0].ID,
						DungeonCharacterID: sql.NullString{String: data.DungeonCharacterRecs[0].ID, Valid: true},
					},
					ActionCharacterRec: &record.DungeonActionCharacter{
						Name: data.DungeonCharacterRecs[0].Name,
					},
					TargetActionObjectRec: &record.DungeonActionObject{
						Name: data.DungeonObjectRecs[0].Name,
					},
				}
			},
			expectError: false,
		},
		{
			name: "Command",
			dungeonID: func(data harness.Data) string {
				return data.DungeonRecs[0].ID
			},
			dungeonCharacterID: func(data harness.Data) string {
				return data.DungeonCharacterRecs[0].ID
			},
			sentence: "look white cat",
			expectDungeonActionRecordSet: func(data harness.Data) *model.DungeonActionRecordSet {
				return &model.DungeonActionRecordSet{
					ActionRec: &record.DungeonAction{
						DungeonID:          data.DungeonRecs[0].ID,
						DungeonLocationID:  data.DungeonLocationRecs[0].ID,
						DungeonCharacterID: sql.NullString{String: data.DungeonCharacterRecs[0].ID, Valid: true},
					},
					ActionCharacterRec: &record.DungeonActionCharacter{
						Name: data.DungeonCharacterRecs[0].Name,
					},
					TargetActionMonsterRec: &record.DungeonActionMonster{
						Name: data.DungeonMonsterRecs[0].Name,
					},
				}
			},
			expectError: false,
		},
		{
			name: "Command",
			dungeonID: func(data harness.Data) string {
				return data.DungeonRecs[0].ID
			},
			dungeonCharacterID: func(data harness.Data) string {
				return data.DungeonCharacterRecs[0].ID
			},
			sentence: "look barricade",
			expectDungeonActionRecordSet: func(data harness.Data) *model.DungeonActionRecordSet {
				return &model.DungeonActionRecordSet{
					ActionRec: &record.DungeonAction{
						DungeonID:          data.DungeonRecs[0].ID,
						DungeonLocationID:  data.DungeonLocationRecs[0].ID,
						DungeonCharacterID: sql.NullString{String: data.DungeonCharacterRecs[0].ID, Valid: true},
					},
					ActionCharacterRec: &record.DungeonActionCharacter{
						Name: data.DungeonCharacterRecs[0].Name,
					},
					TargetActionCharacterRec: &record.DungeonActionCharacter{
						Name: data.DungeonCharacterRecs[0].Name,
					},
				}
			},
			expectError: false,
		},
	}

	for _, tc := range tests {

		t.Logf("Run test >%s< - >%s<", tc.name, tc.sentence)

		func() {

			// Test harness
			err = th.Setup()
			require.NoError(t, err, "Setup returns without error")
			defer func() {
				err = th.RollbackTx()
				require.NoError(t, err, "RollbackTx returns without error")
				err = th.Teardown()
				require.NoError(t, err, "Teardown returns without error")
			}()

			// init tx
			err = th.InitTx(nil)
			require.NoError(t, err, "InitTx returns without error")

			dungeonID := tc.dungeonID(th.Data)
			dungeonCharacterID := tc.dungeonCharacterID(th.Data)

			rslt, err := th.Model.(*model.Model).ProcessDungeonCharacterAction(dungeonID, dungeonCharacterID, tc.sentence)
			if tc.expectError == true {
				require.Error(t, err, "CreateDungeonObjectRec returns error")
				return
			}
			require.NoError(t, err, "ProcessDungeonCharacterAction returns without error")
			require.NotNil(t, rslt.ActionRec, "ProcessDungeonCharacterAction returns DungeonActionRecordSet with ActionRec")

			xrslt := tc.expectDungeonActionRecordSet(th.Data)
			if xrslt == nil {
				return
			}

			if xrslt.ActionRec != nil {
				require.NotNil(t, rslt.ActionRec, "ActionRec is not nil")
				require.Equal(t, xrslt.ActionRec.DungeonID, rslt.ActionRec.DungeonID, "ActionRec.DungeonID equals expected")
				require.Equal(t, xrslt.ActionRec.DungeonCharacterID, rslt.ActionRec.DungeonCharacterID, "ActionRec.DungeonCharacterID equals expected")
				require.Equal(t, xrslt.ActionRec.DungeonMonsterID, rslt.ActionRec.DungeonMonsterID, "ActionRec.DungeonMonsterID equals expected")
				require.Equal(t, xrslt.ActionRec.DungeonLocationID, rslt.ActionRec.DungeonLocationID, "ActionRec.DungeonLocationID equals expected")
			}
			if xrslt.ActionCharacterRec != nil {
				require.NotNil(t, rslt.ActionCharacterRec, "ActionCharacterRec is not nil")
				require.Equal(t, xrslt.ActionCharacterRec.Name, rslt.ActionCharacterRec.Name, "ActionCharacterRec.Name equals expected")
			}
			if xrslt.ActionMonsterRec != nil {
				require.NotNil(t, rslt.ActionMonsterRec, "ActionMonsterRec is not nil")
				require.Equal(t, xrslt.ActionMonsterRec.Name, rslt.ActionMonsterRec.Name, "ActionMonsterRec.Name equals expected")
			}
			if xrslt.EquippedActionObjectRec != nil {
				require.NotNil(t, rslt.EquippedActionObjectRec, "EquippedActionObjectRec is not nil")
				require.Equal(t, xrslt.EquippedActionObjectRec.Name, rslt.EquippedActionObjectRec.Name, "EquippedActionObjectRec.Name equals expected")
			}
			if xrslt.StashedActionObjectRec != nil {
				require.NotNil(t, rslt.StashedActionObjectRec, "StashedActionObjectRec is not nil")
				require.Equal(t, xrslt.StashedActionObjectRec.Name, rslt.StashedActionObjectRec.Name, "StashedActionObjectRec.Name equals expected")
			}
			if xrslt.TargetActionObjectRec != nil {
				require.NotNil(t, rslt.TargetActionObjectRec, "TargetActionObjectRec is not nil")
				require.Equal(t, xrslt.TargetActionObjectRec.Name, rslt.TargetActionObjectRec.Name, "TargetActionObjectRec.Name equals expected")
			}
			if xrslt.TargetActionCharacterRec != nil {
				require.NotNil(t, rslt.TargetActionCharacterRec, "TargetActionCharacterRec is not nil")
				require.Equal(t, xrslt.TargetActionCharacterRec.Name, rslt.TargetActionCharacterRec.Name, "TargetActionCharacterRec.Name equals expected")
			}
			if xrslt.TargetActionMonsterRec != nil {
				require.NotNil(t, rslt.TargetActionMonsterRec, "TargetActionMonsterRec is not nil")
				require.Equal(t, xrslt.TargetActionMonsterRec.Name, rslt.TargetActionMonsterRec.Name, "TargetActionMonsterRec.Name equals expected")
			}
			if xrslt.TargetLocation != nil {
				require.NotNil(t, rslt.TargetLocation, "TargetLocation is not nil")
				if xrslt.TargetLocation.LocationRec != nil {
					require.NotNil(t, rslt.TargetLocation.LocationRec, "TargetLocation.LocationRec is not nil")
					require.Equal(t, xrslt.TargetLocation.LocationRec.Name, rslt.TargetLocation.LocationRec.Name, "TargetLocation.LocationRec equals expected")
				}
			}
		}()
	}
}
