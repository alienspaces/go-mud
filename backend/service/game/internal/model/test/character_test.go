package test

// NOTE: model tests are run is the public space to avoid cyclic dependencies

import (
	"testing"

	"gitlab.com/alienspaces/go-mud/backend/core/repository"

	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/dependencies"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/harness"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/model"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/record"
)

func TestGetCharacterEquippedObjectRecs(t *testing.T) {

	// harness
	config := harness.DefaultDataConfig

	c, l, s, err := dependencies.Default()
	require.NoError(t, err, "NewTesting returns without error")

	th, err := harness.NewTesting(c, l, s, config)
	require.NoError(t, err, "NewTesting returns without error")

	// harness commit data
	th.ShouldCommitData = true

	tests := []struct {
		name               string
		dungeonCharacterID func(data harness.Data) string
		expectObjectRecs   func(data harness.Data) []*record.Object
		expectError        bool
	}{
		{
			name: "Returns expected equipped objects",
			dungeonCharacterID: func(data harness.Data) string {
				cRec, _ := data.GetCharacterRecByName("barricade")
				return cRec.ID
			},
			expectObjectRecs: func(data harness.Data) []*record.Object {
				ceoRec, _ := data.GetObjectRecByName("dull bronze ring")
				return []*record.Object{
					&record.Object{
						Record: repository.Record{
							ID: ceoRec.ID,
						},
						Name: ceoRec.Name,
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
			_, err = th.Setup()
			require.NoError(t, err, "Setup returns without error")
			defer func() {
				err = th.RollbackTx()
				require.NoError(t, err, "RollbackTx returns without error")
				err = th.Teardown()
				require.NoError(t, err, "Teardown returns without error")
			}()

			// init tx
			_, err = th.InitTx()
			require.NoError(t, err, "InitTx returns without error")

			dungeonCharacterID := tc.dungeonCharacterID(th.Data)

			recs, err := th.Model.(*model.Model).GetCharacterEquippedObjectRecs(dungeonCharacterID)
			if tc.expectError == true {
				require.Error(t, err, "GetCharacterEquippedObjectRecs returns error")
				return
			}
			require.NoError(t, err, "GetCharacterEquippedObjectRecs returns without error")

			expectedRecs := tc.expectObjectRecs(th.Data)
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

func TestGetCharacterStashedObjectRecs(t *testing.T) {

	// harness
	config := harness.DefaultDataConfig

	c, l, s, err := dependencies.Default()
	require.NoError(t, err, "NewTesting returns without error")

	th, err := harness.NewTesting(c, l, s, config)
	require.NoError(t, err, "NewTesting returns without error")

	// harness commit data
	th.ShouldCommitData = true

	tests := []struct {
		name               string
		dungeonCharacterID func(data harness.Data) string
		expectObjectRecs   func(data harness.Data) []*record.Object
		expectError        bool
	}{
		{
			name: "Returns expected stashed objects",
			dungeonCharacterID: func(data harness.Data) string {
				cRec, _ := data.GetCharacterRecByName("barricade")
				return cRec.ID
			},
			expectObjectRecs: func(data harness.Data) []*record.Object {
				csoRec, _ := data.GetObjectRecByName("blood stained pouch")
				return []*record.Object{
					&record.Object{
						Record: repository.Record{
							ID: csoRec.ID,
						},
						Name: csoRec.Name,
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
			_, err = th.Setup()
			require.NoError(t, err, "Setup returns without error")
			defer func() {
				err = th.RollbackTx()
				require.NoError(t, err, "RollbackTx returns without error")
				err = th.Teardown()
				require.NoError(t, err, "Teardown returns without error")
			}()

			// init tx
			_, err = th.InitTx()
			require.NoError(t, err, "InitTx returns without error")

			dungeonCharacterID := tc.dungeonCharacterID(th.Data)

			recs, err := th.Model.(*model.Model).GetCharacterStashedObjectRecs(dungeonCharacterID)
			if tc.expectError == true {
				require.Error(t, err, "GetCharacterStashedObjectRecs returns error")
				return
			}
			require.NoError(t, err, "GetCharacterStashedObjectRecs returns without error")

			expectedRecs := tc.expectObjectRecs(th.Data)
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
