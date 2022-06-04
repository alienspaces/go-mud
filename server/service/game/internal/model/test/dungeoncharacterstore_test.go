package test

// NOTE: model tests are run is the public space to avoid cyclic dependencies

import (
	"testing"

	"github.com/brianvoe/gofakeit"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/go-mud/server/service/game/internal/dependencies"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/harness"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/model"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

func TestCreateDungeonCharacterRec(t *testing.T) {

	// harness
	config := harness.DefaultDataConfig

	c, l, s, err := dependencies.Default()
	require.NoError(t, err, "NewTesting returns without error")

	th, err := harness.NewTesting(c, l, s, config)
	require.NoError(t, err, "NewTesting returns without error")

	// harness commit data
	th.CommitData = true

	tests := []struct {
		name string
		rec  func(data harness.Data) *record.DungeonCharacter
		err  bool
	}{
		{
			name: "Without ID",
			rec: func(data harness.Data) *record.DungeonCharacter {
				return &record.DungeonCharacter{
					DungeonID:    data.DungeonRecs[0].ID,
					Name:         gofakeit.StreetName() + gofakeit.Name(),
					Strength:     10,
					Dexterity:    10,
					Intelligence: 10,
				}
			},
			err: false,
		},
		{
			name: "With ID",
			rec: func(data harness.Data) *record.DungeonCharacter {
				rec := &record.DungeonCharacter{
					DungeonID:    data.DungeonRecs[0].ID,
					Name:         gofakeit.StreetName() + gofakeit.Name(),
					Strength:     10,
					Dexterity:    10,
					Intelligence: 10,
				}
				id, _ := uuid.NewRandom()
				rec.ID = id.String()
				return rec
			},
			err: false,
		},
	}

	for _, tc := range tests {

		t.Logf("Run test >%s<", tc.name)

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

			rec := tc.rec(th.Data)

			err = th.Model.(*model.Model).CreateDungeonCharacterRec(rec)
			if tc.err == true {
				require.Error(t, err, "CreateDungeonCharacterRec returns error")
				return
			}
			require.NoError(t, err, "CreateDungeonCharacterRec returns without error")
			require.NotEmpty(t, rec.CreatedAt, "CreateDungeonCharacterRec returns record with CreatedAt")
		}()
	}
}

func TestGetDungeonCharacterRec(t *testing.T) {

	// harness
	config := harness.DefaultDataConfig

	c, l, s, err := dependencies.Default()
	require.NoError(t, err, "NewTesting returns without error")

	h, err := harness.NewTesting(c, l, s, config)
	require.NoError(t, err, "NewTesting returns without error")

	// harness commit data
	h.CommitData = true

	tests := []struct {
		name string
		id   func() string
		err  bool
	}{
		{
			name: "With ID",
			id: func() string {
				return h.Data.DungeonCharacterRecs[0].ID
			},
			err: false,
		},
		{
			name: "Without ID",
			id: func() string {
				return ""
			},
			err: true,
		},
	}

	for _, tc := range tests {

		t.Logf("Run test >%s<", tc.name)

		func() {

			// harness setup
			err = h.Setup()
			require.NoError(t, err, "Setup returns without error")
			defer func() {
				err = h.RollbackTx()
				require.NoError(t, err, "RollbackTx returns without error")
				err = h.Teardown()
				require.NoError(t, err, "Teardown returns without error")
			}()

			// init tx
			err = h.InitTx(nil)
			require.NoError(t, err, "InitTx returns without error")

			rec, err := h.Model.(*model.Model).GetDungeonCharacterRec(tc.id(), false)
			if tc.err == true {
				require.Error(t, err, "GetDungeonCharacterRec returns error")
				return
			}
			require.NoError(t, err, "GetDungeonCharacterRec returns without error")
			require.NotNil(t, rec, "GetDungeonCharacterRec returns record")
			require.NotEmpty(t, rec.ID, "Record ID is not empty")
		}()
	}
}

func TestUpdateDungeonCharacterRec(t *testing.T) {

	// harness
	config := harness.DefaultDataConfig

	c, l, s, err := dependencies.Default()
	require.NoError(t, err, "NewTesting returns without error")

	h, err := harness.NewTesting(c, l, s, config)
	require.NoError(t, err, "NewTesting returns without error")

	// harness commit data
	h.CommitData = true

	require.NoError(t, err, "NewTesting returns without error")

	tests := []struct {
		name string
		rec  func() *record.DungeonCharacter
		err  bool
	}{
		{
			name: "With ID",
			rec: func() *record.DungeonCharacter {
				return h.Data.DungeonCharacterRecs[0]
			},
			err: false,
		},
		{
			name: "Without ID",
			rec: func() *record.DungeonCharacter {
				rec := h.Data.DungeonCharacterRecs[0]
				rec.ID = ""
				return rec
			},
			err: true,
		},
	}

	for _, tc := range tests {

		t.Logf("Run test >%s<", tc.name)

		func() {

			// harness setup
			err = h.Setup()
			require.NoError(t, err, "Setup returns without error")
			defer func() {
				err = h.RollbackTx()
				require.NoError(t, err, "RollbackTx returns without error")
				err = h.Teardown()
				require.NoError(t, err, "Teardown returns without error")
			}()

			// init tx
			err = h.InitTx(nil)
			require.NoError(t, err, "InitTx returns without error")

			rec := tc.rec()

			err := h.Model.(*model.Model).UpdateDungeonCharacterRec(rec)
			if tc.err == true {
				require.Error(t, err, "UpdateDungeonCharacterRec returns error")
				return
			}
			require.NoError(t, err, "UpdateDungeonCharacterRec returns without error")
			require.NotEmpty(t, rec.UpdatedAt, "UpdateDungeonCharacterRec returns record with UpdatedAt")
		}()
	}
}

func TestDeleteDungeonCharacterRec(t *testing.T) {

	// harness
	config := harness.DefaultDataConfig

	c, l, s, err := dependencies.Default()
	require.NoError(t, err, "NewTesting returns without error")

	h, err := harness.NewTesting(c, l, s, config)
	require.NoError(t, err, "NewTesting returns without error")

	// harness commit data
	h.CommitData = true

	tests := []struct {
		name string
		id   func() string
		err  bool
	}{
		{
			name: "With ID",
			id: func() string {
				return h.Data.DungeonCharacterRecs[0].ID
			},
			err: false,
		},
		{
			name: "Without ID",
			id: func() string {
				return ""
			},
			err: true,
		},
	}

	for _, tc := range tests {

		t.Logf("Run test >%s<", tc.name)

		func() {

			// harness setup
			err = h.Setup()
			require.NoError(t, err, "Setup returns without error")
			defer func() {
				err = h.RollbackTx()
				require.NoError(t, err, "RollbackTx returns without error")
				err = h.Teardown()
				require.NoError(t, err, "Teardown returns without error")
			}()

			// init tx
			err = h.InitTx(nil)
			require.NoError(t, err, "InitTx returns without error")

			err := h.Model.(*model.Model).DeleteDungeonCharacterRec(tc.id())
			if tc.err == true {
				require.Error(t, err, "DeleteDungeonCharacterRec returns error")
				return
			}
			require.NoError(t, err, "DeleteDungeonCharacterRec returns without error")

			rec, err := h.Model.(*model.Model).GetDungeonCharacterRec(tc.id(), false)
			require.NoError(t, err, "GetDungeonCharacterRec returns without error")
			require.Nil(t, rec, "GetDungeonCharacterRec does not return record")
		}()
	}
}
