package model

import (
	"database/sql"
	"fmt"

	"gitlab.com/alienspaces/go-boilerplate/server/service/template/internal/record"
)

// GetTemplateRecs -
func (m *Model) GetTemplateRecs(params map[string]interface{}, operators map[string]string, forUpdate bool) ([]*record.Template, error) {

	m.Log.Info("Getting template records params >%s<", params)

	r := m.TemplateRepository()

	return r.GetMany(params, operators, forUpdate)
}

// GetTemplateRec -
func (m *Model) GetTemplateRec(recID string, forUpdate bool) (*record.Template, error) {

	m.Log.Info("Getting template rec ID >%s<", recID)

	r := m.TemplateRepository()

	// validate UUID
	if !m.IsUUID(recID) {
		return nil, fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	rec, err := r.GetOne(recID, forUpdate)
	if err == sql.ErrNoRows {
		m.Log.Warn("No record found ID >%s<", recID)
		return nil, nil
	}

	return rec, err
}

// CreateTemplateRec -
func (m *Model) CreateTemplateRec(rec *record.Template) error {

	m.Log.Info("Creating template rec >%#v<", rec)

	r := m.TemplateRepository()

	err := m.ValidateTemplateRec(rec)
	if err != nil {
		m.Log.Info("Failed model validation >%v<", err)
		return err
	}

	return r.CreateOne(rec)
}

// UpdateTemplateRec -
func (m *Model) UpdateTemplateRec(rec *record.Template) error {

	m.Log.Info("Updating template rec >%#v<", rec)

	r := m.TemplateRepository()

	err := m.ValidateTemplateRec(rec)
	if err != nil {
		m.Log.Info("Failed model validation >%v<", err)
		return err
	}

	return r.UpdateOne(rec)
}

// DeleteTemplateRec -
func (m *Model) DeleteTemplateRec(recID string) error {

	m.Log.Info("Deleting template rec ID >%s<", recID)

	r := m.TemplateRepository()

	// validate UUID
	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.ValidateDeleteTemplateRec(recID)
	if err != nil {
		m.Log.Info("Failed model validation >%v<", err)
		return err
	}

	return r.DeleteOne(recID)
}

// RemoveTemplateRec -
func (m *Model) RemoveTemplateRec(recID string) error {

	m.Log.Info("Removing template rec ID >%s<", recID)

	r := m.TemplateRepository()

	// validate UUID
	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.ValidateDeleteTemplateRec(recID)
	if err != nil {
		m.Log.Info("Failed model validation >%v<", err)
		return err
	}

	return r.RemoveOne(recID)
}
