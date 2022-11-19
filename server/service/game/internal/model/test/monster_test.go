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

func TestGetMonsterEquippedObjectRecs(t *testing.T) {

	// harness
	config := harness.DefaultDataConfig

	c, l, s, err := dependencies.Default()
	require.NoError(t, err, "NewTesting returns without error")

	th, err := harness.NewTesting(c, l, s, config)
	require.NoError(t, err, "NewTesting returns without error")

	// harness commit data
	th.CommitData = true

	tests := []struct {
		name             string
		dungeonMonsterID func(data harness.Data) string
		expectObjectRecs func(data harness.Data) []*record.Object
		expectError      bool
	}{
		{
			name: "Returns expected equipped objects",
			dungeonMonsterID: func(data harness.Data) string {
				cRec, _ := data.GetMonsterRecByName("grumpy dwarf")
				return cRec.ID
			},
			expectObjectRecs: func(data harness.Data) []*record.Object {
				ceoRec, _ := data.GetObjectRecByName("bone dagger")
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

			dungeonMonsterID := tc.dungeonMonsterID(th.Data)

			recs, err := th.Model.(*model.Model).GetMonsterEquippedObjectRecs(dungeonMonsterID)
			if tc.expectError == true {
				require.Error(t, err, "GetMonsterEquippedObjectRecs returns error")
				return
			}
			require.NoError(t, err, "GetMonsterEquippedObjectRecs returns without error")

			expectedRecs := tc.expectObjectRecs(th.Data)
			if expectedRecs == nil {
				return
			}

			for idx, expectedRec := range expectedRecs {
				require.Equal(t, expectedRec.ID, recs[idx].ID, "Returned monster object ID equals expected")
				require.Equal(t, expectedRec.Name, recs[idx].Name, "Returned monster object Name equals expected")
			}
		})
	}
}

func TestGetMonsterStashedObjectRecs(t *testing.T) {

	// harness
	config := harness.DefaultDataConfig

	c, l, s, err := dependencies.Default()
	require.NoError(t, err, "NewTesting returns without error")

	th, err := harness.NewTesting(c, l, s, config)
	require.NoError(t, err, "NewTesting returns without error")

	// harness commit data
	th.CommitData = true

	tests := []struct {
		name             string
		dungeonMonsterID func(data harness.Data) string
		expectObjectRecs func(data harness.Data) []*record.Object
		expectError      bool
	}{
		{
			name: "Returns expected stashed objects",
			dungeonMonsterID: func(data harness.Data) string {
				cRec, _ := data.GetMonsterRecByName("grumpy dwarf")
				return cRec.ID
			},
			expectObjectRecs: func(data harness.Data) []*record.Object {
				csoRec, _ := data.GetObjectRecByName("vial of ogre blood")
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

			dungeonMonsterID := tc.dungeonMonsterID(th.Data)

			recs, err := th.Model.(*model.Model).GetMonsterStashedObjectRecs(dungeonMonsterID)
			if tc.expectError == true {
				require.Error(t, err, "GetMonsterStashedObjectRecs returns error")
				return
			}
			require.NoError(t, err, "GetMonsterStashedObjectRecs returns without error")

			expectedRecs := tc.expectObjectRecs(th.Data)
			if expectedRecs == nil {
				return
			}

			for idx, expectedRec := range expectedRecs {
				require.Equal(t, expectedRec.ID, recs[idx].ID, "Returned monster object ID equals expected")
				require.Equal(t, expectedRec.Name, recs[idx].Name, "Returned monster object Name equals expected")
			}
		})
	}
}
