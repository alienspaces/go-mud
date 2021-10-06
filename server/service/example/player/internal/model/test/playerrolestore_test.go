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

func TestCreatePlayerRoleRec(t *testing.T) {

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
		rec  func(data *harness.Data) *record.PlayerRole
		err  bool
	}{
		{
			name: "Without ID",
			rec: func(data *harness.Data) *record.PlayerRole {
				return &record.PlayerRole{
					PlayerID: data.PlayerRecs[0].ID,
					Role:      record.PlayerRoleDefault,
				}
			},
			err: false,
		},
		{
			name: "With ID",
			rec: func(data *harness.Data) *record.PlayerRole {
				rec := &record.PlayerRole{
					PlayerID: data.PlayerRecs[0].ID,
					Role:      record.PlayerRoleDefault,
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

			rec := tc.rec(h.Data)

			err = h.Model.(*model.Model).CreatePlayerRoleRec(rec)
			if tc.err == true {
				require.Error(t, err, "CreatePlayerRoleRec returns error")
				return
			}
			require.NoError(t, err, "CreatePlayerRoleRec returns without error")
			require.NotEmpty(t, rec.CreatedAt, "CreatePlayerRoleRec returns record with CreatedAt")
		}()
	}
}

func TestGetPlayerRoleRec(t *testing.T) {

	// harness
	config := harness.DataConfig{
		PlayerConfig: []harness.PlayerConfig{
			{
				Record: record.Player{},
				PlayerRoleConfig: []harness.PlayerRoleConfig{
					{
						Record: record.PlayerRole{},
					},
				},
			},
		},
	}

	h, err := harness.NewTesting(config)
	require.NoError(t, err, "NewTesting returns without error")

	// harness commit data
	h.CommitData = true

	tests := []struct {
		name string
		id   func(data *harness.Data) string
		err  bool
	}{
		{
			name: "With ID",
			id: func(data *harness.Data) string {
				return data.PlayerRoleRecs[0].ID
			},
			err: false,
		},
		{
			name: "Without ID",
			id: func(data *harness.Data) string {
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

			rec, err := h.Model.(*model.Model).GetPlayerRoleRec(tc.id(h.Data), false)
			if tc.err == true {
				require.Error(t, err, "GetPlayerRoleRec returns error")
				return
			}
			require.NoError(t, err, "GetPlayerRoleRec returns without error")
			require.NotNil(t, rec, "GetPlayerRoleRec returns record")
			require.NotEmpty(t, rec.ID, "Record ID is not empty")
		}()
	}
}

func TestUpdatePlayerRoleRec(t *testing.T) {

	// harness
	config := harness.DataConfig{
		PlayerConfig: []harness.PlayerConfig{
			{
				Record: record.Player{},
				PlayerRoleConfig: []harness.PlayerRoleConfig{
					{
						Record: record.PlayerRole{},
					},
				},
			},
		},
	}

	h, err := harness.NewTesting(config)

	// harness commit data
	h.CommitData = true

	require.NoError(t, err, "NewTesting returns without error")

	tests := []struct {
		name string
		rec  func(data *harness.Data) record.PlayerRole
		err  bool
	}{
		{
			name: "With ID",
			rec: func(data *harness.Data) record.PlayerRole {
				return data.PlayerRoleRecs[0]
			},
			err: false,
		},
		{
			name: "Without ID",
			rec: func(data *harness.Data) record.PlayerRole {
				rec := data.PlayerRoleRecs[0]
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

			rec := tc.rec(h.Data)

			err := h.Model.(*model.Model).UpdatePlayerRoleRec(&rec)
			if tc.err == true {
				require.Error(t, err, "UpdatePlayerRoleRec returns error")
				return
			}
			require.NoError(t, err, "UpdatePlayerRoleRec returns without error")
			require.NotEmpty(t, rec.UpdatedAt, "UpdatePlayerRoleRec returns record with UpdatedAt")
		}()
	}
}

func TestDeletePlayerRoleRec(t *testing.T) {

	// harness
	config := harness.DataConfig{
		PlayerConfig: []harness.PlayerConfig{
			{
				Record: record.Player{},
				PlayerRoleConfig: []harness.PlayerRoleConfig{
					{
						Record: record.PlayerRole{},
					},
				},
			},
		},
	}

	h, err := harness.NewTesting(config)
	require.NoError(t, err, "NewTesting returns without error")

	// harness commit data
	h.CommitData = true

	tests := []struct {
		name string
		id   func(data *harness.Data) string
		err  bool
	}{
		{
			name: "With ID",
			id: func(data *harness.Data) string {
				return data.PlayerRoleRecs[0].ID
			},
			err: false,
		},
		{
			name: "Without ID",
			id: func(data *harness.Data) string {
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

			err := h.Model.(*model.Model).DeletePlayerRoleRec(tc.id(h.Data))
			if tc.err == true {
				require.Error(t, err, "DeletePlayerRoleRec returns error")
				return
			}
			require.NoError(t, err, "DeletePlayerRoleRec returns without error")

			rec, err := h.Model.(*model.Model).GetPlayerRoleRec(tc.id(h.Data), false)
			require.NoError(t, err, "GetPlayerRoleRec returns without error")
			require.Nil(t, rec, "GetPlayerRoleRec does not return record")
		}()
	}
}
