package test

// NOTE: model tests are run is the public space to avoid cyclic dependencies

import (
	"testing"

	"github.com/brianvoe/gofakeit"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/dependencies"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/harness"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/model"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/record"
)

func TestCreateObjectRec(t *testing.T) {

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
		rec  func(data harness.Data) *record.Object
		err  bool
	}{
		{
			name: "Without ID",
			rec: func(data harness.Data) *record.Object {
				return &record.Object{
					Name:                gofakeit.StreetName() + gofakeit.Name(),
					Description:         gofakeit.StreetName() + gofakeit.Name(),
					DescriptionDetailed: gofakeit.StreetName() + gofakeit.Name(),
				}
			},
			err: false,
		},
		{
			name: "With ID",
			rec: func(data harness.Data) *record.Object {
				rec := &record.Object{
					Name:                gofakeit.StreetName() + gofakeit.Name(),
					Description:         gofakeit.StreetName() + gofakeit.Name(),
					DescriptionDetailed: gofakeit.StreetName() + gofakeit.Name(),
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

			err = th.Model.(*model.Model).CreateObjectRec(rec)
			if tc.err == true {
				require.Error(t, err, "CreateObjectRec returns error")
				return
			}
			require.NoError(t, err, "CreateObjectRec returns without error")
			require.NotEmpty(t, rec.CreatedAt, "CreateObjectRec returns record with CreatedAt")
		}()
	}
}

func TestGetObjectRec(t *testing.T) {

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
				return h.Data.ObjectRecs[0].ID
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

			rec, err := h.Model.(*model.Model).GetObjectRec(tc.id(), false)
			if tc.err == true {
				require.Error(t, err, "GetObjectRec returns error")
				return
			}
			require.NoError(t, err, "GetObjectRec returns without error")
			require.NotNil(t, rec, "GetObjectRec returns record")
			require.NotEmpty(t, rec.ID, "Record ID is not empty")
		}()
	}
}

func TestUpdateObjectRec(t *testing.T) {

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
		rec  func() *record.Object
		err  bool
	}{
		{
			name: "With ID",
			rec: func() *record.Object {
				return h.Data.ObjectRecs[0]
			},
			err: false,
		},
		{
			name: "Without ID",
			rec: func() *record.Object {
				rec := h.Data.ObjectRecs[0]
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

			err := h.Model.(*model.Model).UpdateObjectRec(rec)
			if tc.err == true {
				require.Error(t, err, "UpdateObjectRec returns error")
				return
			}
			require.NoError(t, err, "UpdateObjectRec returns without error")
			require.NotEmpty(t, rec.UpdatedAt, "UpdateObjectRec returns record with UpdatedAt")
		}()
	}
}

func TestDeleteObjectRec(t *testing.T) {

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
				return h.Data.ObjectRecs[0].ID
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

			err := h.Model.(*model.Model).DeleteObjectRec(tc.id())
			if tc.err == true {
				require.Error(t, err, "DeleteObjectRec returns error")
				return
			}
			require.NoError(t, err, "DeleteObjectRec returns without error")

			rec, err := h.Model.(*model.Model).GetObjectRec(tc.id(), false)
			require.NoError(t, err, "GetObjectRec returns without error")
			require.Nil(t, rec, "GetObjectRec does not return record")
		}()
	}
}
