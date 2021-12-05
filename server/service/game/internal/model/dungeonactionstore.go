package model

import (
	"database/sql"
	"fmt"

	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// GetDungeonActionRecs -
func (m *Model) GetDungeonActionRecs(params map[string]interface{}, operators map[string]string, forUpdate bool) ([]*record.DungeonAction, error) {

	m.Log.Debug("Getting dungeon action records params >%s<", params)

	r := m.DungeonActionRepository()

	return r.GetMany(params, operators, forUpdate)
}

// GetDungeonActionRec -
func (m *Model) GetDungeonActionRec(recID string, forUpdate bool) (*record.DungeonAction, error) {

	m.Log.Debug("Getting dungeon action rec ID >%s<", recID)

	r := m.DungeonActionRepository()

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

// CreateDungeonActionRec -
func (m *Model) CreateDungeonActionRec(rec *record.DungeonAction) error {

	m.Log.Debug("Creating dungeon action rec >%#v<", rec)

	r := m.DungeonActionRepository()

	err := m.ValidateDungeonActionRec(rec)
	if err != nil {
		m.Log.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.CreateOne(rec)
}

// UpdateDungeonActionRec -
func (m *Model) UpdateDungeonActionRec(rec *record.DungeonAction) error {

	m.Log.Debug("Updating dungeon action rec >%#v<", rec)

	r := m.DungeonActionRepository()

	err := m.ValidateDungeonActionRec(rec)
	if err != nil {
		m.Log.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.UpdateOne(rec)
}

// DeleteDungeonActionRec -
func (m *Model) DeleteDungeonActionRec(recID string) error {

	m.Log.Debug("Deleting dungeon action rec ID >%s<", recID)

	r := m.DungeonActionRepository()

	// validate UUID
	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.ValidateDeleteDungeonActionRec(recID)
	if err != nil {
		m.Log.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.DeleteOne(recID)
}

// RemoveDungeonActionRec -
func (m *Model) RemoveDungeonActionRec(recID string) error {

	m.Log.Debug("Removing dungeon action rec ID >%s<", recID)

	r := m.DungeonActionRepository()

	// validate UUID
	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.ValidateDeleteDungeonActionRec(recID)
	if err != nil {
		m.Log.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.RemoveOne(recID)
}
