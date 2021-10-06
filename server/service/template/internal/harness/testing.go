package harness

import (
	"gitlab.com/alienspaces/go-boilerplate/server/core/harness"
	"gitlab.com/alienspaces/go-boilerplate/server/core/type/modeller"
	"gitlab.com/alienspaces/go-boilerplate/server/service/template/internal/model"
	"gitlab.com/alienspaces/go-boilerplate/server/service/template/internal/record"
)

// Testing -
type Testing struct {
	harness.Testing
	Data       *Data
	DataConfig DataConfig
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

// NewTesting -
func NewTesting(config DataConfig) (t *Testing, err error) {

	// harness
	t = &Testing{}

	// modeller
	t.ModellerFunc = t.Modeller

	// data
	t.CreateDataFunc = t.CreateData
	t.RemoveDataFunc = t.RemoveData

	t.DataConfig = config
	t.Data = &Data{}

	err = t.Init()
	if err != nil {
		return nil, err
	}

	return t, nil
}

// Modeller -
func (t *Testing) Modeller() (modeller.Modeller, error) {

	m, err := model.NewModel(t.Config, t.Log, t.Store)
	if err != nil {
		t.Log.Warn("Failed new model >%v<", err)
		return nil, err
	}

	return m, nil
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
	}

	return nil
}

// RemoveData -
func (t *Testing) RemoveData() error {

TEMPLATE_RECS:
	for {
		if len(t.Data.TemplateRecs) == 0 {
			break TEMPLATE_RECS
		}
		rec := record.Template{}
		rec, t.Data.TemplateRecs = t.Data.TemplateRecs[0], t.Data.TemplateRecs[1:]

		err := t.Model.(*model.Model).RemoveTemplateRec(rec.ID)
		if err != nil {
			t.Log.Warn("Failed removing template record >%v<", err)
			return err
		}
	}

	return nil
}
