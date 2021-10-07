package test

// NOTE: model tests are run is the public space so we are
// able to use common setup and teardown tooling for all models

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/go-mud/server/service/game/internal/harness"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/model"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

func TestCreateGameRec(t *testing.T) {

	// harness
	config := harness.DataConfig{}

	h, err := harness.NewTesting(config)
	require.NoError(t, err, "NewTesting returns without error")

	// harness commit data
	h.CommitData = true

	tests := []struct {
		name string
		rec  func() *record.Game
		err  bool
	}{
		{
			name: "Without ID",
			rec: func() *record.Game {
				return &record.Game{}
			},
			err: false,
		},
		{
			name: "With ID",
			rec: func() *record.Game {
				rec := &record.Game{}
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

			err = h.Model.(*model.Model).CreateGameRec(rec)
			if tc.err == true {
				require.Error(t, err, "CreateGameRec returns error")
				return
			}
			require.NoError(t, err, "CreateGameRec returns without error")
			require.NotEmpty(t, rec.CreatedAt, "CreateGameRec returns record with CreatedAt")
		}()
	}
}

func TestGetGameRec(t *testing.T) {

	// harness
	config := harness.DataConfig{
		GameConfig: []harness.GameConfig{
			{
				Record: record.Game{},
			},
		},
	}

	h, err := harness.NewTesting(config)
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
				return h.Data.GameRecs[0].ID
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

			rec, err := h.Model.(*model.Model).GetGameRec(tc.id(), false)
			if tc.err == true {
				require.Error(t, err, "GetGameRec returns error")
				return
			}
			require.NoError(t, err, "GetGameRec returns without error")
			require.NotNil(t, rec, "GetGameRec returns record")
			require.NotEmpty(t, rec.ID, "Record ID is not empty")
		}()
	}
}

func TestUpdateGameRec(t *testing.T) {

	// harness
	config := harness.DataConfig{
		GameConfig: []harness.GameConfig{
			{
				Record: record.Game{},
			},
		},
	}

	h, err := harness.NewTesting(config)

	// harness commit data
	h.CommitData = true

	require.NoError(t, err, "NewTesting returns without error")

	tests := []struct {
		name string
		rec  func() record.Game
		err  bool
	}{
		{
			name: "With ID",
			rec: func() record.Game {
				return h.Data.GameRecs[0]
			},
			err: false,
		},
		{
			name: "Without ID",
			rec: func() record.Game {
				rec := h.Data.GameRecs[0]
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

			err := h.Model.(*model.Model).UpdateGameRec(&rec)
			if tc.err == true {
				require.Error(t, err, "UpdateGameRec returns error")
				return
			}
			require.NoError(t, err, "UpdateGameRec returns without error")
			require.NotEmpty(t, rec.UpdatedAt, "UpdateGameRec returns record with UpdatedAt")
		}()
	}
}

func TestDeleteGameRec(t *testing.T) {

	// harness
	config := harness.DataConfig{
		GameConfig: []harness.GameConfig{
			{
				Record: record.Game{},
			},
		},
	}

	h, err := harness.NewTesting(config)
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
				return h.Data.GameRecs[0].ID
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

			err := h.Model.(*model.Model).DeleteGameRec(tc.id())
			if tc.err == true {
				require.Error(t, err, "DeleteGameRec returns error")
				return
			}
			require.NoError(t, err, "DeleteGameRec returns without error")

			rec, err := h.Model.(*model.Model).GetGameRec(tc.id(), false)
			require.NoError(t, err, "GetGameRec returns without error")
			require.Nil(t, rec, "GetGameRec does not return record")
		}()
	}
}
