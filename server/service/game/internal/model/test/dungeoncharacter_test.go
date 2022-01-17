package test

// NOTE: model tests are run is the public space to avoid cyclic dependencies

import (
	"testing"

	"gitlab.com/alienspaces/go-mud/server/core/repository"

	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/go-mud/server/service/game/internal/dependencies"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/harness"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/model"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

func TestGetCharacterEquippedDungeonObjectRecs(t *testing.T) {

	// harness
	config := harness.DefaultDataConfig

	c, l, s, m, err := dependencies.Default()
	require.NoError(t, err, "NewTesting returns without error")

	th, err := harness.NewTesting(c, l, s, m, config)
	require.NoError(t, err, "NewTesting returns without error")

	// harness commit data
	th.CommitData = true

	tests := []struct {
		name                    string
		dungeonCharacterID      func(data harness.Data) string
		expectDungeonObjectRecs func(data harness.Data) []*record.DungeonObject
		expectError             bool
	}{
		{
			name: "Returns objects for character with objects",
			dungeonCharacterID: func(data harness.Data) string {
				return data.DungeonCharacterRecs[0].ID
			},
			expectDungeonObjectRecs: func(data harness.Data) []*record.DungeonObject {
				return []*record.DungeonObject{
					&record.DungeonObject{
						Record: repository.Record{
							ID: data.DungeonObjectRecs[2].ID,
						},
						Name: data.DungeonObjectRecs[2].Name,
					},
				}
			},
			expectError: false,
		},
	}

	for _, tc := range tests {

		t.Run(tc.name, func(t *testing.T) {
			t.Logf("Run test >%s<", tc.name)

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

			dungeonCharacterID := tc.dungeonCharacterID(th.Data)

			recs, err := th.Model.(*model.Model).GetCharacterEquippedDungeonObjectRecs(dungeonCharacterID)
			if tc.expectError == true {
				require.Error(t, err, "GetCharacterEquippedDungeonObjectRecs returns error")
				return
			}
			require.NoError(t, err, "GetCharacterEquippedDungeonObjectRecs returns without error")

			expectedRecs := tc.expectDungeonObjectRecs(th.Data)
			if expectedRecs == nil {
				return
			}

			for idx, expectedRec := range expectedRecs {
				require.Equal(t, expectedRec.ID, recs[idx].ID, "Returned character object ID equals expected")
				require.Equal(t, expectedRec.Name, recs[idx].Name, "Returned character object Name equals expected")
			}
		})
	}
}

func TestGetCharacterStashedDungeonObjectRecs(t *testing.T) {

	// harness
	config := harness.DefaultDataConfig

	c, l, s, m, err := dependencies.Default()
	require.NoError(t, err, "NewTesting returns without error")

	th, err := harness.NewTesting(c, l, s, m, config)
	require.NoError(t, err, "NewTesting returns without error")

	// harness commit data
	th.CommitData = true

	tests := []struct {
		name                    string
		dungeonCharacterID      func(data harness.Data) string
		expectDungeonObjectRecs func(data harness.Data) []*record.DungeonObject
		expectError             bool
	}{
		{
			name: "Returns objects for character with objects",
			dungeonCharacterID: func(data harness.Data) string {
				return data.DungeonCharacterRecs[0].ID
			},
			expectDungeonObjectRecs: func(data harness.Data) []*record.DungeonObject {
				return []*record.DungeonObject{
					&record.DungeonObject{
						Record: repository.Record{
							ID: data.DungeonObjectRecs[3].ID,
						},
						Name: data.DungeonObjectRecs[3].Name,
					},
				}
			},
			expectError: false,
		},
	}

	for _, tc := range tests {

		t.Run(tc.name, func(t *testing.T) {
			t.Logf("Run test >%s<", tc.name)

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

			dungeonCharacterID := tc.dungeonCharacterID(th.Data)

			recs, err := th.Model.(*model.Model).GetCharacterStashedDungeonObjectRecs(dungeonCharacterID)
			if tc.expectError == true {
				require.Error(t, err, "GetCharacterStashedDungeonObjectRecs returns error")
				return
			}
			require.NoError(t, err, "GetCharacterStashedDungeonObjectRecs returns without error")

			expectedRecs := tc.expectDungeonObjectRecs(th.Data)
			if expectedRecs == nil {
				return
			}

			for idx, expectedRec := range expectedRecs {
				require.Equal(t, expectedRec.ID, recs[idx].ID, "Returned character object ID equals expected")
				require.Equal(t, expectedRec.Name, recs[idx].Name, "Returned character object Name equals expected")
			}
		})
	}
}
