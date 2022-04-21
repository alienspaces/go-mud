package model

import (
	"database/sql"
	"fmt"

	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// GetActionObjectRecs -
func (m *Model) GetActionObjectRecs(params map[string]interface{}, operators map[string]string, forUpdate bool) ([]*record.ActionObject, error) {

	m.Log.Debug("Getting dungeon action object records params >%s<", params)

	r := m.ActionObjectRepository()

	return r.GetMany(params, operators, forUpdate)
}

// GetActionObjectRec -
func (m *Model) GetActionObjectRec(recID string, forUpdate bool) (*record.ActionObject, error) {

	m.Log.Debug("Getting dungeon action object rec ID >%s<", recID)

	r := m.ActionObjectRepository()

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

// CreateActionObjectRec -
func (m *Model) CreateActionObjectRec(rec *record.ActionObject) error {

	m.Log.Debug("Creating dungeon action object rec >%#v<", rec)

	r := m.ActionObjectRepository()

	err := m.ValidateActionObjectRec(rec)
	if err != nil {
		m.Log.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.CreateOne(rec)
}

// UpdateActionObjectRec -
func (m *Model) UpdateActionObjectRec(rec *record.ActionObject) error {

	m.Log.Debug("Updating dungeon action object rec >%#v<", rec)

	r := m.ActionObjectRepository()

	err := m.ValidateActionObjectRec(rec)
	if err != nil {
		m.Log.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.UpdateOne(rec)
}

// DeleteActionObjectRec -
func (m *Model) DeleteActionObjectRec(recID string) error {

	m.Log.Debug("Deleting dungeon action object rec ID >%s<", recID)

	r := m.ActionObjectRepository()

	// validate UUID
	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.ValidateDeleteActionObjectRec(recID)
	if err != nil {
		m.Log.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.DeleteOne(recID)
}

// RemoveActionObjectRec -
func (m *Model) RemoveActionObjectRec(recID string) error {

	m.Log.Debug("Removing dungeon action object rec ID >%s<", recID)

	r := m.ActionObjectRepository()

	// validate UUID
	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.ValidateDeleteActionObjectRec(recID)
	if err != nil {
		m.Log.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.RemoveOne(recID)
}
