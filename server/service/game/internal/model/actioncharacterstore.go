package model

import (
	"database/sql"
	"fmt"

	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// GetActionCharacterRecs -
func (m *Model) GetActionCharacterRecs(params map[string]interface{}, operators map[string]string, forUpdate bool) ([]*record.ActionCharacter, error) {

	m.Log.Info("Getting dungeon action character records params >%s<", params)

	r := m.ActionCharacterRepository()

	return r.GetMany(params, operators, forUpdate)
}

// GetActionCharacterRec -
func (m *Model) GetActionCharacterRec(recID string, forUpdate bool) (*record.ActionCharacter, error) {

	m.Log.Info("Getting dungeon action character rec ID >%s<", recID)

	r := m.ActionCharacterRepository()

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

// CreateActionCharacterRec -
func (m *Model) CreateActionCharacterRec(rec *record.ActionCharacter) error {

	m.Log.Info("Creating dungeon action character rec >%#v<", rec)

	r := m.ActionCharacterRepository()

	err := m.ValidateActionCharacterRec(rec)
	if err != nil {
		m.Log.Info("Failed model validation >%v<", err)
		return err
	}

	return r.CreateOne(rec)
}

// UpdateActionCharacterRec -
func (m *Model) UpdateActionCharacterRec(rec *record.ActionCharacter) error {

	m.Log.Info("Updating dungeon action character rec >%#v<", rec)

	r := m.ActionCharacterRepository()

	err := m.ValidateActionCharacterRec(rec)
	if err != nil {
		m.Log.Info("Failed model validation >%v<", err)
		return err
	}

	return r.UpdateOne(rec)
}

// DeleteActionCharacterRec -
func (m *Model) DeleteActionCharacterRec(recID string) error {

	m.Log.Info("Deleting dungeon action character rec ID >%s<", recID)

	r := m.ActionCharacterRepository()

	// validate UUID
	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.ValidateDeleteActionCharacterRec(recID)
	if err != nil {
		m.Log.Info("Failed model validation >%v<", err)
		return err
	}

	return r.DeleteOne(recID)
}

// RemoveActionCharacterRec -
func (m *Model) RemoveActionCharacterRec(recID string) error {

	m.Log.Info("Removing dungeon action character rec ID >%s<", recID)

	r := m.ActionCharacterRepository()

	// validate UUID
	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.ValidateDeleteActionCharacterRec(recID)
	if err != nil {
		m.Log.Info("Failed model validation >%v<", err)
		return err
	}

	return r.RemoveOne(recID)
}
