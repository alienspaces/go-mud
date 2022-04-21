package model

import (
	"database/sql"
	"fmt"

	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// GetObjectRecs -
func (m *Model) GetObjectRecs(params map[string]interface{}, operators map[string]string, forUpdate bool) ([]*record.Object, error) {

	m.Log.Debug("Getting dungeon object records params >%s<", params)

	r := m.ObjectRepository()

	return r.GetMany(params, operators, forUpdate)
}

// GetObjectRec -
func (m *Model) GetObjectRec(recID string, forUpdate bool) (*record.Object, error) {

	m.Log.Debug("Getting dungeon object rec ID >%s<", recID)

	r := m.ObjectRepository()

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

// CreateObjectRec -
func (m *Model) CreateObjectRec(rec *record.Object) error {

	m.Log.Debug("Creating dungeon object rec >%#v<", rec)

	r := m.ObjectRepository()

	err := m.ValidateObjectRec(rec)
	if err != nil {
		m.Log.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.CreateOne(rec)
}

// UpdateObjectRec -
func (m *Model) UpdateObjectRec(rec *record.Object) error {

	m.Log.Debug("Updating dungeon object rec >%#v<", rec)

	r := m.ObjectRepository()

	err := m.ValidateObjectRec(rec)
	if err != nil {
		m.Log.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.UpdateOne(rec)
}

// DeleteObjectRec -
func (m *Model) DeleteObjectRec(recID string) error {

	m.Log.Debug("Deleting dungeon object rec ID >%s<", recID)

	r := m.ObjectRepository()

	// validate UUID
	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.ValidateDeleteObjectRec(recID)
	if err != nil {
		m.Log.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.DeleteOne(recID)
}

// RemoveObjectRec -
func (m *Model) RemoveObjectRec(recID string) error {

	m.Log.Debug("Removing dungeon object rec ID >%s<", recID)

	r := m.ObjectRepository()

	// validate UUID
	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.ValidateDeleteObjectRec(recID)
	if err != nil {
		m.Log.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.RemoveOne(recID)
}
