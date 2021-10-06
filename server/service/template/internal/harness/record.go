package harness

import (
	"gitlab.com/alienspaces/go-boilerplate/server/service/template/internal/model"
	"gitlab.com/alienspaces/go-boilerplate/server/service/template/internal/record"
)

func (t *Testing) createTemplateRec(templateConfig TemplateConfig) (record.Template, error) {

	rec := templateConfig.Record

	// NOTE: Add default values for required properties here

	t.Log.Info("Creating testing record >%#v<", rec)

	err := t.Model.(*model.Model).CreateTemplateRec(&rec)
	if err != nil {
		t.Log.Warn("Failed creating testing template record >%v<", err)
		return rec, err
	}
	return rec, nil
}
