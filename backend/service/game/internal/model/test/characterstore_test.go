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

func TestCreateCharacterRec(t *testing.T) {

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
		rec  func(data harness.Data) *record.Character
		err  bool
	}{
		{
			name: "Without ID",
			rec: func(data harness.Data) *record.Character {
				return &record.Character{
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
			rec: func(data harness.Data) *record.Character {
				rec := &record.Character{
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

			err = th.Model.(*model.Model).CreateCharacterRec(rec)
			if tc.err == true {
				require.Error(t, err, "CreateCharacterRec returns error")
				return
			}
			require.NoError(t, err, "CreateCharacterRec returns without error")
			require.NotEmpty(t, rec.CreatedAt, "CreateCharacterRec returns record with CreatedAt")
		}()
	}
}

func TestGetCharacterRec(t *testing.T) {

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
				return h.Data.CharacterRecs[0].ID
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

			rec, err := h.Model.(*model.Model).GetCharacterRec(tc.id(), false)
			if tc.err == true {
				require.Error(t, err, "GetCharacterRec returns error")
				return
			}
			require.NoError(t, err, "GetCharacterRec returns without error")
			require.NotNil(t, rec, "GetCharacterRec returns record")
			require.NotEmpty(t, rec.ID, "Record ID is not empty")
		}()
	}
}

func TestUpdateCharacterRec(t *testing.T) {

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
		rec  func() *record.Character
		err  bool
	}{
		{
			name: "With ID",
			rec: func() *record.Character {
				return h.Data.CharacterRecs[0]
			},
			err: false,
		},
		{
			name: "Without ID",
			rec: func() *record.Character {
				rec := h.Data.CharacterRecs[0]
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

			err := h.Model.(*model.Model).UpdateCharacterRec(rec)
			if tc.err == true {
				require.Error(t, err, "UpdateCharacterRec returns error")
				return
			}
			require.NoError(t, err, "UpdateCharacterRec returns without error")
			require.NotEmpty(t, rec.UpdatedAt, "UpdateCharacterRec returns record with UpdatedAt")
		}()
	}
}

func TestDeleteCharacterRec(t *testing.T) {

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
				return h.Data.CharacterRecs[0].ID
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

			err := h.Model.(*model.Model).DeleteCharacterRec(tc.id())
			if tc.err == true {
				require.Error(t, err, "DeleteCharacterRec returns error")
				return
			}
			require.NoError(t, err, "DeleteCharacterRec returns without error")

			rec, err := h.Model.(*model.Model).GetCharacterRec(tc.id(), false)
			require.NoError(t, err, "GetCharacterRec returns without error")
			require.Nil(t, rec, "GetCharacterRec does not return record")
		}()
	}
}
