package test

// NOTE: model tests are run is the public space so we are
// able to use common setup and teardown tooling for all models

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/go-boilerplate/server/service/player/internal/harness"
	"gitlab.com/alienspaces/go-boilerplate/server/service/player/internal/model"
	"gitlab.com/alienspaces/go-boilerplate/server/service/player/internal/record"
)

func TestCreatePlayerRec(t *testing.T) {

	// harness
	config := harness.DataConfig{}

	h, err := harness.NewTesting(config)
	require.NoError(t, err, "NewTesting returns without error")

	// harness commit data
	h.CommitData = true

	tests := []struct {
		name string
		rec  func() *record.Player
		err  bool
	}{
		{
			name: "Without ID",
			rec: func() *record.Player {
				return &record.Player{
					Name:              "Scary Susan",
					Email:             "scarysusan@example.com",
					Provider:          record.AccountProviderGoogle,
					ProviderAccountID: "abcdefg",
				}
			},
			err: false,
		},
		{
			name: "With ID",
			rec: func() *record.Player {
				rec := &record.Player{
					Name:              "Horrific Harry",
					Email:             "horrificharry@example.com",
					Provider:          record.AccountProviderGoogle,
					ProviderAccountID: "abcdefg",
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

			err = h.Model.(*model.Model).CreatePlayerRec(rec)
			if tc.err == true {
				require.Error(t, err, "CreatePlayerRec returns error")
				return
			}
			require.NoError(t, err, "CreatePlayerRec returns without error")
			require.NotEmpty(t, rec.CreatedAt, "CreatePlayerRec returns record with CreatedAt")
		}()
	}
}

func TestGetPlayerRec(t *testing.T) {

	// harness
	config := harness.DataConfig{
		PlayerConfig: []harness.PlayerConfig{
			{
				Record: record.Player{},
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
				return h.Data.PlayerRecs[0].ID
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

			rec, err := h.Model.(*model.Model).GetPlayerRec(tc.id(), false)
			if tc.err == true {
				require.Error(t, err, "GetPlayerRec returns error")
				return
			}
			require.NoError(t, err, "GetPlayerRec returns without error")
			require.NotNil(t, rec, "GetPlayerRec returns record")
			require.NotEmpty(t, rec.ID, "Record ID is not empty")
		}()
	}
}

func TestUpdatePlayerRec(t *testing.T) {

	// harness
	config := harness.DataConfig{
		PlayerConfig: []harness.PlayerConfig{
			{
				Record: record.Player{},
			},
		},
	}

	h, err := harness.NewTesting(config)

	// harness commit data
	h.CommitData = true

	require.NoError(t, err, "NewTesting returns without error")

	tests := []struct {
		name string
		rec  func() record.Player
		err  bool
	}{
		{
			name: "With ID",
			rec: func() record.Player {
				return h.Data.PlayerRecs[0]
			},
			err: false,
		},
		{
			name: "Without ID",
			rec: func() record.Player {
				rec := h.Data.PlayerRecs[0]
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

			err := h.Model.(*model.Model).UpdatePlayerRec(&rec)
			if tc.err == true {
				require.Error(t, err, "UpdatePlayerRec returns error")
				return
			}
			require.NoError(t, err, "UpdatePlayerRec returns without error")
			require.NotEmpty(t, rec.UpdatedAt, "UpdatePlayerRec returns record with UpdatedAt")
		}()
	}
}

func TestDeletePlayerRec(t *testing.T) {

	// harness
	config := harness.DataConfig{
		PlayerConfig: []harness.PlayerConfig{
			{
				Record: record.Player{},
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
				return h.Data.PlayerRecs[0].ID
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

			err := h.Model.(*model.Model).DeletePlayerRec(tc.id())
			if tc.err == true {
				require.Error(t, err, "DeletePlayerRec returns error")
				return
			}
			require.NoError(t, err, "DeletePlayerRec returns without error")

			rec, err := h.Model.(*model.Model).GetPlayerRec(tc.id(), false)
			require.NoError(t, err, "GetPlayerRec returns without error")
			require.Nil(t, rec, "GetPlayerRec does not return record")
		}()
	}
}
