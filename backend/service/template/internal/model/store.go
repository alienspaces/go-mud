package model

import (
	"database/sql"
	"errors"
	"fmt"

	coreerror "gitlab.com/alienspaces/go-mud/backend/core/error"
	coresql "gitlab.com/alienspaces/go-mud/backend/core/sql"
	templateerror "gitlab.com/alienspaces/go-mud/backend/service/template/internal/error"
	"gitlab.com/alienspaces/go-mud/backend/service/template/internal/record"
)

// GetTemplateRecs -
func (m *Model) GetTemplateRecs(opts *coresql.Options) ([]*record.Template, error) {

	m.Log.Debug("Getting template records options >%#v<", opts)

	r := m.TemplateRepository()

	recs, err := r.GetMany(opts)
	if err != nil {
		return recs, templateerror.NewDatabaseError(m.Log, err)
	}

	return recs, nil
}

// GetTemplateRec -
func (m *Model) GetTemplateRec(recID string, lock *coresql.Lock) (*record.Template, error) {

	m.Log.Debug("Getting template record ID >%s<", recID)

	if !m.IsUUID(recID) {
		return nil, fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	r := m.TemplateRepository()

	rec, err := r.GetOne(recID, lock)
	if errors.Is(err, sql.ErrNoRows) {
		return rec, coreerror.NewNotFoundError("template", recID)
	} else if err != nil {
		return rec, templateerror.NewDatabaseError(m.Log, err)
	}

	return rec, err
}

// CreateTemplateRec -
func (m *Model) CreateTemplateRec(rec *record.Template) (*record.Template, error) {

	m.Log.Debug("Creating template rec >%v<", rec)

	r := m.TemplateRepository()

	if err := m.ValidateTemplateRec(rec); err != nil {
		m.Log.Warn("Failed model validation >%v<", err)
		return rec, err
	}

	if err := r.CreateOne(rec); err != nil {
		return rec, templateerror.NewDatabaseError(m.Log, err)
	}

	return rec, nil
}

// UpdateTemplateRec -
func (m *Model) UpdateTemplateRec(rec *record.Template) (*record.Template, error) {

	m.Log.Debug("Updating template rec >%v<", rec)

	r := m.TemplateRepository()

	if err := m.ValidateTemplateRec(rec); err != nil {
		m.Log.Warn("Failed model validation >%v<", err)
		return rec, err
	}

	if err := r.UpdateOne(rec); err != nil {
		return rec, templateerror.NewDatabaseError(m.Log, err)
	}

	return rec, nil
}

// DeleteTemplateRec -
func (m *Model) DeleteTemplateRec(recID string) error {

	m.Log.Debug("Deleting template record ID >%s<", recID)

	r := m.TemplateRepository()

	if err := m.ValidateDeleteTemplateRec(recID); err != nil {
		m.Log.Warn("Failed model validation >%v<", err)
		return err
	}

	if err := r.DeleteOne(recID); err != nil {
		return templateerror.NewDatabaseError(m.Log, err)
	}

	return nil
}

// RemoveTemplateRec -
func (m *Model) RemoveTemplateRec(recID string) error {

	m.Log.Debug("Removing template record ID >%s<", recID)

	r := m.TemplateRepository()

	err := m.ValidateDeleteTemplateRec(recID)
	if err != nil {
		m.Log.Warn("Failed model validation >%v<", err)
		return err
	}

	if err := r.RemoveOne(recID); err != nil {
		return templateerror.NewDatabaseError(m.Log, err)
	}

	return nil
}
