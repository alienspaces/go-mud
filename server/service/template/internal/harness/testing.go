package harness

import (
	"gitlab.com/alienspaces/go-mud/server/core/harness"
	"gitlab.com/alienspaces/go-mud/server/core/type/configurer"
	"gitlab.com/alienspaces/go-mud/server/core/type/logger"
	"gitlab.com/alienspaces/go-mud/server/core/type/modeller"
	"gitlab.com/alienspaces/go-mud/server/core/type/storer"
	"gitlab.com/alienspaces/go-mud/server/service/template/internal/model"
	"gitlab.com/alienspaces/go-mud/server/service/template/internal/record"
)

// Testing -
type Testing struct {
	harness.Testing
	Data         Data
	DataConfig   DataConfig
	teardownData teardownData
}

// DataConfig -
type DataConfig struct {
	TemplateConfig []TemplateConfig
}

// TemplateConfig -
type TemplateConfig struct {
	Record record.Template
}

// Data -
type Data struct {
	TemplateRecs []record.Template
}

// teardownData -
type teardownData struct {
	TemplateRecs []record.Template
}

// NewTesting -
func NewTesting(c configurer.Configurer, l logger.Logger, s storer.Storer, m modeller.Modeller, config DataConfig) (t *Testing, err error) {

	t = &Testing{
		Testing: harness.Testing{
			Config: c,
			Log:    l,
			Store:  s,
			Model:  m,
		},
	}

	// modeller
	t.ModellerFunc = t.Modeller

	// data
	t.CreateDataFunc = t.CreateData
	t.RemoveDataFunc = t.RemoveData

	t.DataConfig = config
	t.Data = Data{}
	t.teardownData = teardownData{}

	err = t.Init()
	if err != nil {
		return nil, err
	}

	return t, nil
}

// Modeller -
func (t *Testing) Modeller() (modeller.Modeller, error) {
	return t.Model, nil
}

// CreateData - Custom data
func (t *Testing) CreateData() error {

	for _, templateConfig := range t.DataConfig.TemplateConfig {

		templateRec, err := t.createTemplateRec(templateConfig)
		if err != nil {
			t.Log.Warn("Failed creating template record >%v<", err)
			return err
		}
		t.Data.TemplateRecs = append(t.Data.TemplateRecs, templateRec)
		t.teardownData.TemplateRecs = append(t.teardownData.TemplateRecs, templateRec)
	}

	return nil
}

// RemoveData -
func (t *Testing) RemoveData() error {

TEMPLATE_RECS:
	for {
		if len(t.teardownData.TemplateRecs) == 0 {
			break TEMPLATE_RECS
		}
		rec := record.Template{}
		rec, t.teardownData.TemplateRecs = t.teardownData.TemplateRecs[0], t.teardownData.TemplateRecs[1:]

		err := t.Model.(*model.Model).RemoveTemplateRec(rec.ID)
		if err != nil {
			t.Log.Warn("Failed removing template record >%v<", err)
			return err
		}
	}

	return nil
}
