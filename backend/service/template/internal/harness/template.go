package harness

import (
	"fmt"

	"gitlab.com/alienspaces/go-mud/backend/service/template/internal/model"
	"gitlab.com/alienspaces/go-mud/backend/service/template/internal/record"
)

func (t *Testing) processTemplateConfig(templateConfig TemplateConfig) (*record.Template, error) {
	l := t.Logger("processTemplateConfig")

	rec := templateConfig.Record
	rec, err := t.assignTemplateRecDefaults(rec)
	if err != nil {
		l.Warn("failed to assign template rec defaults >%v<", err)
		return nil, err
	}

	rec, err = t.createTemplateRec(rec)
	if err != nil {
		l.Warn("failed creating template record >%v<", err)
		return nil, err
	}

	l.Info("** created template record >%#v<", rec)

	t.appendTemplateRec(*rec)

	return rec, nil
}

func (t *Testing) assignTemplateRecDefaults(rec *record.Template) (*record.Template, error) {
	if rec == nil {
		rec = &record.Template{}
	}

	return rec, nil
}

func (t *Testing) createTemplateRec(rec *record.Template) (*record.Template, error) {
	l := t.Logger("createTemplateRec")

	if rec == nil {
		msg := "template record is nil, cannot create template record"
		l.Warn(msg)
		return nil, fmt.Errorf(msg)
	}

	rec, err := t.Model.(*model.Model).CreateTemplateRec(rec)
	if err != nil {
		l.Warn("failed creating template record >%v<", err)
		return rec, err
	}

	return rec, nil
}

func (t *Testing) appendTemplateRec(rec record.Template) {
	for _, r := range t.Data.TemplateRecs {
		if r.ID == rec.ID {
			t.Log.Debug("- template rec already appended >%s<", rec.ID)
			return
		}
	}

	l := t.Logger("appendTemplateRec")

	l.Debug("- appending template record ID >%s<", rec.ID)
	t.Data.TemplateRecs = append(t.Data.TemplateRecs, rec)
	t.AddTemplateTeardownID(rec.ID)
}

func (t *Testing) AddTemplateTeardownID(id string) {
	l := t.Logger("AddTemplateTeardownID")

	for _, rec := range t.teardownData.TemplateRecs {
		if rec.ID == id {
			l.Debug("- template teardown ID already appended >%s<", id)
			return
		}
	}

	l.Debug("- adding template teardown ID >%s<", id)
	rec := record.Template{}
	rec.ID = id
	t.teardownData.TemplateRecs = append(t.teardownData.TemplateRecs, rec)
}

func (t *Testing) removeTemplateRec(rec *record.Template) error {
	l := t.Logger("removeTemplateRec")

	l.Debug("- removing template record ID >%s<", rec.ID)

	err := t.Model.(*model.Model).RemoveTemplateRec(rec.ID)
	if err != nil {
		l.Warn("failed removing template record >%v<", err)
		return err
	}
	return nil
}
