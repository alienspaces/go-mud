package model

import (
	"database/sql"
	"fmt"

	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// GetObjectInstanceRecs -
func (m *Model) GetObjectInstanceRecs(params map[string]interface{}, operators map[string]string, forUpdate bool) ([]*record.ObjectInstance, error) {

	m.Log.Debug("Getting object instance records params >%s<", params)

	r := m.ObjectInstanceRepository()

	return r.GetMany(params, operators, forUpdate)
}

// GetObjectInstanceRec -
func (m *Model) GetObjectInstanceRec(recID string, forUpdate bool) (*record.ObjectInstance, error) {

	m.Log.Debug("Getting object instance record ID >%s<", recID)

	r := m.ObjectInstanceRepository()

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

// GetObjectInstanceViewRecs -
func (m *Model) GetObjectInstanceViewRecs(params map[string]interface{}, operators map[string]string) ([]*record.ObjectInstanceView, error) {

	m.Log.Debug("Getting object instance view records params >%s<", params)

	r := m.ObjectInstanceViewRepository()

	return r.GetMany(params, operators, false)
}

// GetObjectInstanceViewRec -
func (m *Model) GetObjectInstanceViewRec(recID string) (*record.ObjectInstanceView, error) {

	m.Log.Debug("Getting object instance view record ID >%s<", recID)

	r := m.ObjectInstanceViewRepository()

	// validate UUID
	if !m.IsUUID(recID) {
		return nil, fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	rec, err := r.GetOne(recID, false)
	if err == sql.ErrNoRows {
		m.Log.Warn("No record found ID >%s<", recID)
		return nil, nil
	}

	return rec, err
}

// CreateObjectInstanceRec -
func (m *Model) CreateObjectInstanceRec(rec *record.ObjectInstance) error {

	m.Log.Debug("Creating object rec >%#v<", rec)

	r := m.ObjectInstanceRepository()

	err := m.ValidateObjectInstanceRec(rec)
	if err != nil {
		m.Log.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.CreateOne(rec)
}

// UpdateObjectInstanceRec -
func (m *Model) UpdateObjectInstanceRec(rec *record.ObjectInstance) error {

	m.Log.Debug("Updating object rec >%#v<", rec)

	r := m.ObjectInstanceRepository()

	err := m.ValidateObjectInstanceRec(rec)
	if err != nil {
		m.Log.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.UpdateOne(rec)
}

// DeleteObjectInstanceRec -
func (m *Model) DeleteObjectInstanceRec(recID string) error {

	m.Log.Debug("Deleting object rec ID >%s<", recID)

	r := m.ObjectInstanceRepository()

	// validate UUID
	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.ValidateDeleteObjectInstanceRec(recID)
	if err != nil {
		m.Log.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.DeleteOne(recID)
}

// RemoveObjectInstanceRec -
func (m *Model) RemoveObjectInstanceRec(recID string) error {

	m.Log.Debug("Removing object rec ID >%s<", recID)

	r := m.ObjectInstanceRepository()

	// validate UUID
	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.ValidateDeleteObjectInstanceRec(recID)
	if err != nil {
		m.Log.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.RemoveOne(recID)
}
