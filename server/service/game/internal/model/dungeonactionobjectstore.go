package model

import (
	"database/sql"
	"fmt"

	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// GetDungeonActionObjectRecs -
func (m *Model) GetDungeonActionObjectRecs(params map[string]interface{}, operators map[string]string, forUpdate bool) ([]*record.DungeonActionObject, error) {

	m.Log.Debug("Getting dungeon action object records params >%s<", params)

	r := m.DungeonActionObjectRepository()

	return r.GetMany(params, operators, forUpdate)
}

// GetDungeonActionObjectRec -
func (m *Model) GetDungeonActionObjectRec(recID string, forUpdate bool) (*record.DungeonActionObject, error) {

	m.Log.Debug("Getting dungeon action object rec ID >%s<", recID)

	r := m.DungeonActionObjectRepository()

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

// CreateDungeonActionObjectRec -
func (m *Model) CreateDungeonActionObjectRec(rec *record.DungeonActionObject) error {

	m.Log.Debug("Creating dungeon action object rec >%#v<", rec)

	r := m.DungeonActionObjectRepository()

	err := m.ValidateDungeonActionObjectRec(rec)
	if err != nil {
		m.Log.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.CreateOne(rec)
}

// UpdateDungeonActionObjectRec -
func (m *Model) UpdateDungeonActionObjectRec(rec *record.DungeonActionObject) error {

	m.Log.Debug("Updating dungeon action object rec >%#v<", rec)

	r := m.DungeonActionObjectRepository()

	err := m.ValidateDungeonActionObjectRec(rec)
	if err != nil {
		m.Log.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.UpdateOne(rec)
}

// DeleteDungeonActionObjectRec -
func (m *Model) DeleteDungeonActionObjectRec(recID string) error {

	m.Log.Debug("Deleting dungeon action object rec ID >%s<", recID)

	r := m.DungeonActionObjectRepository()

	// validate UUID
	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.ValidateDeleteDungeonActionObjectRec(recID)
	if err != nil {
		m.Log.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.DeleteOne(recID)
}

// RemoveDungeonActionObjectRec -
func (m *Model) RemoveDungeonActionObjectRec(recID string) error {

	m.Log.Debug("Removing dungeon action object rec ID >%s<", recID)

	r := m.DungeonActionObjectRepository()

	// validate UUID
	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.ValidateDeleteDungeonActionObjectRec(recID)
	if err != nil {
		m.Log.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.RemoveOne(recID)
}
