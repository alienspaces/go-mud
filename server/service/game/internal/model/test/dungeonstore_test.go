package test

// NOTE: model tests are run is the public space to avoid cyclic dependencies

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/go-mud/server/service/game/internal/dependencies"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/harness"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/model"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

func TestCreateDungeonRec(t *testing.T) {

	// harness
	config := harness.DataConfig{}

	c, l, s, err := dependencies.Default()
	require.NoError(t, err, "NewTesting returns without error")

	h, err := harness.NewTesting(c, l, s, config)
	require.NoError(t, err, "NewTesting returns without error")

	// harness commit data
	h.CommitData = true

	tests := []struct {
		name string
		rec  func() *record.Dungeon
		err  bool
	}{
		{
			name: "Without ID",
			rec: func() *record.Dungeon {
				return &record.Dungeon{}
			},
			err: false,
		},
		{
			name: "With ID",
			rec: func() *record.Dungeon {
				rec := &record.Dungeon{}
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

			err = h.Model.(*model.Model).CreateDungeonRec(rec)
			if tc.err == true {
				require.Error(t, err, "CreateDungeonRec returns error")
				return
			}
			require.NoError(t, err, "CreateDungeonRec returns without error")
			require.NotEmpty(t, rec.CreatedAt, "CreateDungeonRec returns record with CreatedAt")
		}()
	}
}

func TestGetDungeonRec(t *testing.T) {

	// harness
	config := harness.DataConfig{
		DungeonConfig: []harness.DungeonConfig{
			{
				Record: record.Dungeon{},
			},
		},
	}

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
				return h.Data.DungeonRecs[0].ID
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

			rec, err := h.Model.(*model.Model).GetDungeonRec(tc.id(), false)
			if tc.err == true {
				require.Error(t, err, "GetDungeonRec returns error")
				return
			}
			require.NoError(t, err, "GetDungeonRec returns without error")
			require.NotNil(t, rec, "GetDungeonRec returns record")
			require.NotEmpty(t, rec.ID, "Record ID is not empty")
		}()
	}
}

func TestUpdateDungeonRec(t *testing.T) {

	// harness
	config := harness.DataConfig{
		DungeonConfig: []harness.DungeonConfig{
			{
				Record: record.Dungeon{},
			},
		},
	}

	c, l, s, err := dependencies.Default()
	require.NoError(t, err, "NewTesting returns without error")

	h, err := harness.NewTesting(c, l, s, config)
	require.NoError(t, err, "NewTesting returns without error")

	// harness commit data
	h.CommitData = true

	require.NoError(t, err, "NewTesting returns without error")

	tests := []struct {
		name string
		rec  func() *record.Dungeon
		err  bool
	}{
		{
			name: "With ID",
			rec: func() *record.Dungeon {
				return h.Data.DungeonRecs[0]
			},
			err: false,
		},
		{
			name: "Without ID",
			rec: func() *record.Dungeon {
				rec := h.Data.DungeonRecs[0]
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

			err := h.Model.(*model.Model).UpdateDungeonRec(rec)
			if tc.err == true {
				require.Error(t, err, "UpdateDungeonRec returns error")
				return
			}
			require.NoError(t, err, "UpdateDungeonRec returns without error")
			require.NotEmpty(t, rec.UpdatedAt, "UpdateDungeonRec returns record with UpdatedAt")
		}()
	}
}

func TestDeleteDungeonRec(t *testing.T) {

	// harness
	config := harness.DataConfig{
		DungeonConfig: []harness.DungeonConfig{
			{
				Record: record.Dungeon{},
			},
		},
	}

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
				return h.Data.DungeonRecs[0].ID
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

			err := h.Model.(*model.Model).DeleteDungeonRec(tc.id())
			if tc.err == true {
				require.Error(t, err, "DeleteDungeonRec returns error")
				return
			}
			require.NoError(t, err, "DeleteDungeonRec returns without error")

			rec, err := h.Model.(*model.Model).GetDungeonRec(tc.id(), false)
			require.NoError(t, err, "GetDungeonRec returns without error")
			require.Nil(t, rec, "GetDungeonRec does not return record")
		}()
	}
}
