package test

// NOTE: repository tests are run is the public space so we are
// able to use common setup and teardown tooling for all repositories

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/go-boilerplate/server/service/player/internal/harness"
	"gitlab.com/alienspaces/go-boilerplate/server/service/player/internal/model"
	"gitlab.com/alienspaces/go-boilerplate/server/service/player/internal/record"
)

func TestCreateOne(t *testing.T) {

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
				err = h.Teardown()
				require.NoError(t, err, "Teardown returns without error")
			}()

			// init tx
			err = h.InitTx(nil)
			require.NoError(t, err, "InitTx returns without error")

			// repository
			r := h.Model.(*model.Model).PlayerRepository()
			require.NotNil(t, r, "Repository is not nil")

			rec := tc.rec()

			err = r.CreateOne(rec)
			if tc.err == true {
				require.Error(t, err, "CreateOne returns error")
				return
			}
			require.NoError(t, err, "CreateOne returns without error")
			require.NotEmpty(t, rec.CreatedAt, "CreateOne returns record with CreatedAt")

			h.RollbackTx()
		}()
	}
}

func TestGetOne(t *testing.T) {

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
				err = h.Teardown()
				require.NoError(t, err, "Teardown returns without error")
			}()

			// init tx
			err = h.InitTx(nil)
			require.NoError(t, err, "InitTx returns without error")

			// repository
			r := h.Model.(*model.Model).PlayerRepository()
			require.NotNil(t, r, "Repository is not nil")

			rec, err := r.GetOne(tc.id(), false)
			if tc.err == true {
				require.Error(t, err, "GetOne returns error")
				return
			}
			require.NoError(t, err, "GetOne returns without error")
			require.NotNil(t, rec, "GetOne returns record")
			require.NotEmpty(t, rec.ID, "Record ID is not empty")

			h.RollbackTx()
		}()
	}
}

func TestUpdateOne(t *testing.T) {

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
				err = h.Teardown()
				require.NoError(t, err, "Teardown returns without error")
			}()

			// init tx
			err = h.InitTx(nil)
			require.NoError(t, err, "InitTx returns without error")

			// repository
			r := h.Model.(*model.Model).PlayerRepository()
			require.NotNil(t, r, "Repository is not nil")

			rec := tc.rec()

			err := r.UpdateOne(&rec)
			if tc.err == true {
				require.Error(t, err, "UpdateOne returns error")
				return
			}
			require.NoError(t, err, "UpdateOne returns without error")
			require.NotEmpty(t, rec.UpdatedAt, "UpdateOne returns record with UpdatedAt")

			h.RollbackTx()
		}()
	}
}

func TestDeleteOne(t *testing.T) {

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
				err = h.Teardown()
				require.NoError(t, err, "Teardown returns without error")
			}()

			// init tx
			err = h.InitTx(nil)
			require.NoError(t, err, "InitTx returns without error")

			// repository
			r := h.Model.(*model.Model).PlayerRepository()
			require.NotNil(t, r, "Repository is not nil")

			err := r.DeleteOne(tc.id())
			if tc.err == true {
				require.Error(t, err, "DeleteOne returns error")
				return
			}
			require.NoError(t, err, "DeleteOne returns without error")

			rec, err := r.GetOne(tc.id(), false)
			require.Error(t, err, "GetOne returns error")
			require.Nil(t, rec, "GetOne does not return record")

			h.RollbackTx()
		}()
	}
}
