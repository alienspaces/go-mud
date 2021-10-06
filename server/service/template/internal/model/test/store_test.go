package test

// NOTE: model tests are run is the public space so we are
// able to use common setup and teardown tooling for all models

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/go-boilerplate/server/service/template/internal/harness"
	"gitlab.com/alienspaces/go-boilerplate/server/service/template/internal/model"
	"gitlab.com/alienspaces/go-boilerplate/server/service/template/internal/record"
)

func TestCreateTemplateRec(t *testing.T) {

	// harness
	config := harness.DataConfig{}

	h, err := harness.NewTesting(config)
	require.NoError(t, err, "NewTesting returns without error")

	// harness commit data
	h.CommitData = true

	tests := []struct {
		name string
		rec  func() *record.Template
		err  bool
	}{
		{
			name: "Without ID",
			rec: func() *record.Template {
				return &record.Template{}
			},
			err: false,
		},
		{
			name: "With ID",
			rec: func() *record.Template {
				rec := &record.Template{}
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

			err = h.Model.(*model.Model).CreateTemplateRec(rec)
			if tc.err == true {
				require.Error(t, err, "CreateTemplateRec returns error")
				return
			}
			require.NoError(t, err, "CreateTemplateRec returns without error")
			require.NotEmpty(t, rec.CreatedAt, "CreateTemplateRec returns record with CreatedAt")
		}()
	}
}

func TestGetTemplateRec(t *testing.T) {

	// harness
	config := harness.DataConfig{
		TemplateConfig: []harness.TemplateConfig{
			{
				Record: record.Template{},
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
				return h.Data.TemplateRecs[0].ID
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

			rec, err := h.Model.(*model.Model).GetTemplateRec(tc.id(), false)
			if tc.err == true {
				require.Error(t, err, "GetTemplateRec returns error")
				return
			}
			require.NoError(t, err, "GetTemplateRec returns without error")
			require.NotNil(t, rec, "GetTemplateRec returns record")
			require.NotEmpty(t, rec.ID, "Record ID is not empty")
		}()
	}
}

func TestUpdateTemplateRec(t *testing.T) {

	// harness
	config := harness.DataConfig{
		TemplateConfig: []harness.TemplateConfig{
			{
				Record: record.Template{},
			},
		},
	}

	h, err := harness.NewTesting(config)

	// harness commit data
	h.CommitData = true

	require.NoError(t, err, "NewTesting returns without error")

	tests := []struct {
		name string
		rec  func() record.Template
		err  bool
	}{
		{
			name: "With ID",
			rec: func() record.Template {
				return h.Data.TemplateRecs[0]
			},
			err: false,
		},
		{
			name: "Without ID",
			rec: func() record.Template {
				rec := h.Data.TemplateRecs[0]
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

			err := h.Model.(*model.Model).UpdateTemplateRec(&rec)
			if tc.err == true {
				require.Error(t, err, "UpdateTemplateRec returns error")
				return
			}
			require.NoError(t, err, "UpdateTemplateRec returns without error")
			require.NotEmpty(t, rec.UpdatedAt, "UpdateTemplateRec returns record with UpdatedAt")
		}()
	}
}

func TestDeleteTemplateRec(t *testing.T) {

	// harness
	config := harness.DataConfig{
		TemplateConfig: []harness.TemplateConfig{
			{
				Record: record.Template{},
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
				return h.Data.TemplateRecs[0].ID
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

			err := h.Model.(*model.Model).DeleteTemplateRec(tc.id())
			if tc.err == true {
				require.Error(t, err, "DeleteTemplateRec returns error")
				return
			}
			require.NoError(t, err, "DeleteTemplateRec returns without error")

			rec, err := h.Model.(*model.Model).GetTemplateRec(tc.id(), false)
			require.NoError(t, err, "GetTemplateRec returns without error")
			require.Nil(t, rec, "GetTemplateRec does not return record")
		}()
	}
}
