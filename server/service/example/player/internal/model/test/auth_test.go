package test

// NOTE: model tests are run is the public space so we are
// able to use common setup and teardown tooling for all models

import (
	"testing"

	"github.com/brianvoe/gofakeit"
	"gitlab.com/alienspaces/go-boilerplate/server/service/player/internal/record"

	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/go-boilerplate/server/service/player/internal/harness"
	"gitlab.com/alienspaces/go-boilerplate/server/service/player/internal/model"
)

func TestVerifyAuth(t *testing.T) {

	// harness
	config := harness.DataConfig{}

	h, err := harness.NewTesting(config)
	require.NoError(t, err, "NewTesting returns without error")

	// harness commit data
	h.CommitData = true

	tests := []struct {
		name string
		data func() model.AuthData
		err  bool
	}{
		{
			name: "Anonymous",
			data: func() model.AuthData {
				return model.AuthData{
					Provider:          record.AccountProviderAnonymous,
					ProviderAccountID: gofakeit.UUID(),
				}
			},
			err: false,
		},
		{
			name: "Google",
			data: func() model.AuthData {
				return model.AuthData{
					Provider:          record.AccountProviderGoogle,
					ProviderAccountID: gofakeit.UUID(),
					PlayerName:        gofakeit.Name(),
					PlayerEmail:       gofakeit.Email(),
				}
			},
			err: true,
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

			data := tc.data()

			rec, err := h.Model.(*model.Model).AuthVerify(data)
			if tc.err == true {
				require.Error(t, err, "AuthVerify returns error")
				return
			}
			require.NoError(t, err, "AuthVerify returns without error")
			require.NotEmpty(t, rec.Provider, "AuthVerify returns record with Provider")
			require.NotEmpty(t, rec.ProviderAccountID, "AuthVerify returns record with ProviderAccountID")
			require.NotEmpty(t, rec.CreatedAt, "AuthVerify returns record with CreatedAt")
		}()
	}
}
