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

func TestCreateDungeonMonsterRec(t *testing.T) {

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
		rec  func(data harness.Data) *record.DungeonMonster
		err  bool
	}{
		{
			name: "Without ID",
			rec: func(data harness.Data) *record.DungeonMonster {
				return &record.DungeonMonster{
					DungeonID:         data.DungeonRecs[0].ID,
					LocationID: data.LocationRecs[0].ID,
					Name:              gofakeit.StreetName() + gofakeit.Name(),
					Strength:          10,
					Dexterity:         10,
					Intelligence:      10,
				}
			},
			err: false,
		},
		{
			name: "With ID",
			rec: func(data harness.Data) *record.DungeonMonster {
				rec := &record.DungeonMonster{
					DungeonID:         data.DungeonRecs[0].ID,
					LocationID: data.LocationRecs[0].ID,
					Name:              gofakeit.StreetName() + gofakeit.Name(),
					Strength:          10,
					Dexterity:         10,
					Intelligence:      10,
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

			err = th.Model.(*model.Model).CreateDungeonMonsterRec(rec)
			if tc.err == true {
				require.Error(t, err, "CreateDungeonMonsterRec returns error")
				return
			}
			require.NoError(t, err, "CreateDungeonMonsterRec returns without error")
			require.NotEmpty(t, rec.CreatedAt, "CreateDungeonMonsterRec returns record with CreatedAt")
		}()
	}
}

func TestGetDungeonMonsterRec(t *testing.T) {

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
				return h.Data.DungeonMonsterRecs[0].ID
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

			rec, err := h.Model.(*model.Model).GetDungeonMonsterRec(tc.id(), false)
			if tc.err == true {
				require.Error(t, err, "GetDungeonMonsterRec returns error")
				return
			}
			require.NoError(t, err, "GetDungeonMonsterRec returns without error")
			require.NotNil(t, rec, "GetDungeonMonsterRec returns record")
			require.NotEmpty(t, rec.ID, "Record ID is not empty")
		}()
	}
}

func TestUpdateDungeonMonsterRec(t *testing.T) {

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
		rec  func() *record.DungeonMonster
		err  bool
	}{
		{
			name: "With ID",
			rec: func() *record.DungeonMonster {
				return h.Data.DungeonMonsterRecs[0]
			},
			err: false,
		},
		{
			name: "Without ID",
			rec: func() *record.DungeonMonster {
				rec := h.Data.DungeonMonsterRecs[0]
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

			err := h.Model.(*model.Model).UpdateDungeonMonsterRec(rec)
			if tc.err == true {
				require.Error(t, err, "UpdateDungeonMonsterRec returns error")
				return
			}
			require.NoError(t, err, "UpdateDungeonMonsterRec returns without error")
			require.NotEmpty(t, rec.UpdatedAt, "UpdateDungeonMonsterRec returns record with UpdatedAt")
		}()
	}
}

func TestDeleteDungeonMonsterRec(t *testing.T) {

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
				return h.Data.DungeonMonsterRecs[0].ID
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

			err := h.Model.(*model.Model).DeleteDungeonMonsterRec(tc.id())
			if tc.err == true {
				require.Error(t, err, "DeleteDungeonMonsterRec returns error")
				return
			}
			require.NoError(t, err, "DeleteDungeonMonsterRec returns without error")

			rec, err := h.Model.(*model.Model).GetDungeonMonsterRec(tc.id(), false)
			require.NoError(t, err, "GetDungeonMonsterRec returns without error")
			require.Nil(t, rec, "GetDungeonMonsterRec does not return record")
		}()
	}
}
