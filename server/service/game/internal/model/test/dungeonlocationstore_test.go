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

func TestCreateDungeonLocationRec(t *testing.T) {

	// harness
	config := harness.DefaultDataConfig

	c, l, s, m, err := dependencies.Default()
	require.NoError(t, err, "NewTesting returns without error")

	th, err := harness.NewTesting(c, l, s, m, config)
	require.NoError(t, err, "NewTesting returns without error")

	// harness commit data
	th.CommitData = true

	tests := []struct {
		name string
		rec  func(data harness.Data) *record.DungeonLocation
		err  bool
	}{
		{
			name: "Without ID",
			rec: func(data harness.Data) *record.DungeonLocation {
				return &record.DungeonLocation{
					DungeonID: data.DungeonRecs[0].ID,
					Name:      gofakeit.StreetName() + gofakeit.Name(),
				}
			},
			err: false,
		},
		{
			name: "With ID",
			rec: func(data harness.Data) *record.DungeonLocation {
				rec := &record.DungeonLocation{
					DungeonID: data.DungeonRecs[0].ID,
					Name:      gofakeit.StreetName() + gofakeit.Name(),
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

			err = th.Model.(*model.Model).CreateDungeonLocationRec(rec)
			if tc.err == true {
				require.Error(t, err, "CreateDungeonLocationRec returns error")
				return
			}
			require.NoError(t, err, "CreateDungeonLocationRec returns without error")
			require.NotEmpty(t, rec.CreatedAt, "CreateDungeonLocationRec returns record with CreatedAt")
		}()
	}
}

func TestGetDungeonLocationRec(t *testing.T) {

	// harness
	config := harness.DefaultDataConfig

	c, l, s, m, err := dependencies.Default()
	require.NoError(t, err, "NewTesting returns without error")

	h, err := harness.NewTesting(c, l, s, m, config)
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
				return h.Data.DungeonLocationRecs[0].ID
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

			rec, err := h.Model.(*model.Model).GetDungeonLocationRec(tc.id(), false)
			if tc.err == true {
				require.Error(t, err, "GetDungeonLocationRec returns error")
				return
			}
			require.NoError(t, err, "GetDungeonLocationRec returns without error")
			require.NotNil(t, rec, "GetDungeonLocationRec returns record")
			require.NotEmpty(t, rec.ID, "Record ID is not empty")
		}()
	}
}

func TestUpdateDungeonLocationRec(t *testing.T) {

	// harness
	config := harness.DefaultDataConfig

	c, l, s, m, err := dependencies.Default()
	require.NoError(t, err, "NewTesting returns without error")

	h, err := harness.NewTesting(c, l, s, m, config)
	require.NoError(t, err, "NewTesting returns without error")

	// harness commit data
	h.CommitData = true

	require.NoError(t, err, "NewTesting returns without error")

	tests := []struct {
		name string
		rec  func() *record.DungeonLocation
		err  bool
	}{
		{
			name: "With ID",
			rec: func() *record.DungeonLocation {
				return h.Data.DungeonLocationRecs[0]
			},
			err: false,
		},
		{
			name: "Without ID",
			rec: func() *record.DungeonLocation {
				rec := h.Data.DungeonLocationRecs[0]
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

			err := h.Model.(*model.Model).UpdateDungeonLocationRec(rec)
			if tc.err == true {
				require.Error(t, err, "UpdateDungeonLocationRec returns error")
				return
			}
			require.NoError(t, err, "UpdateDungeonLocationRec returns without error")
			require.NotEmpty(t, rec.UpdatedAt, "UpdateDungeonLocationRec returns record with UpdatedAt")
		}()
	}
}

func TestDeleteDungeonLocationRec(t *testing.T) {

	// harness
	config := harness.DefaultDataConfig

	c, l, s, m, err := dependencies.Default()
	require.NoError(t, err, "NewTesting returns without error")

	h, err := harness.NewTesting(c, l, s, m, config)
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
				return h.Data.DungeonLocationRecs[0].ID
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

			err := h.Model.(*model.Model).DeleteDungeonLocationRec(tc.id())
			if tc.err == true {
				require.Error(t, err, "DeleteDungeonLocationRec returns error")
				return
			}
			require.NoError(t, err, "DeleteDungeonLocationRec returns without error")

			rec, err := h.Model.(*model.Model).GetDungeonLocationRec(tc.id(), false)
			require.NoError(t, err, "GetDungeonLocationRec returns without error")
			require.Nil(t, rec, "GetDungeonLocationRec does not return record")
		}()
	}
}
