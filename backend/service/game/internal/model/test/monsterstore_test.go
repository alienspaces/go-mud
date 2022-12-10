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

func TestCreateMonsterRec(t *testing.T) {

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
		rec  func(data harness.Data) *record.Monster
		err  bool
	}{
		{
			name: "Without ID",
			rec: func(data harness.Data) *record.Monster {
				return &record.Monster{
					Name:         gofakeit.StreetName() + gofakeit.Name(),
					Description:  gofakeit.Sentence(7),
					Strength:     10,
					Dexterity:    10,
					Intelligence: 10,
				}
			},
			err: false,
		},
		{
			name: "With ID",
			rec: func(data harness.Data) *record.Monster {
				rec := &record.Monster{
					Name:         gofakeit.StreetName() + gofakeit.Name(),
					Description:  gofakeit.Sentence(7),
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

			err = th.Model.(*model.Model).CreateMonsterRec(rec)
			if tc.err == true {
				require.Error(t, err, "CreateMonsterRec returns error")
				return
			}
			require.NoError(t, err, "CreateMonsterRec returns without error")
			require.NotEmpty(t, rec.CreatedAt, "CreateMonsterRec returns record with CreatedAt")
		}()
	}
}

func TestGetMonsterRec(t *testing.T) {

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
				return h.Data.MonsterRecs[0].ID
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

			rec, err := h.Model.(*model.Model).GetMonsterRec(tc.id(), false)
			if tc.err == true {
				require.Error(t, err, "GetMonsterRec returns error")
				return
			}
			require.NoError(t, err, "GetMonsterRec returns without error")
			require.NotNil(t, rec, "GetMonsterRec returns record")
			require.NotEmpty(t, rec.ID, "Record ID is not empty")
		}()
	}
}

func TestUpdateMonsterRec(t *testing.T) {

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
		rec  func() *record.Monster
		err  bool
	}{
		{
			name: "With ID",
			rec: func() *record.Monster {
				return h.Data.MonsterRecs[0]
			},
			err: false,
		},
		{
			name: "Without ID",
			rec: func() *record.Monster {
				rec := h.Data.MonsterRecs[0]
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

			err := h.Model.(*model.Model).UpdateMonsterRec(rec)
			if tc.err == true {
				require.Error(t, err, "UpdateMonsterRec returns error")
				return
			}
			require.NoError(t, err, "UpdateMonsterRec returns without error")
			require.NotEmpty(t, rec.UpdatedAt, "UpdateMonsterRec returns record with UpdatedAt")
		}()
	}
}

func TestDeleteMonsterRec(t *testing.T) {

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
				return h.Data.MonsterRecs[0].ID
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

			err := h.Model.(*model.Model).DeleteMonsterRec(tc.id())
			if tc.err == true {
				require.Error(t, err, "DeleteMonsterRec returns error")
				return
			}
			require.NoError(t, err, "DeleteMonsterRec returns without error")

			rec, err := h.Model.(*model.Model).GetMonsterRec(tc.id(), false)
			require.NoError(t, err, "GetMonsterRec returns without error")
			require.Nil(t, rec, "GetMonsterRec does not return record")
		}()
	}
}
