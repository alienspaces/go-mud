package model

import (
	"database/sql"
	"fmt"

	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// GetDungeonActionCharacterObjectRecs -
func (m *Model) GetDungeonActionCharacterObjectRecs(params map[string]interface{}, operators map[string]string, forUpdate bool) ([]*record.DungeonActionCharacterObject, error) {

	m.Log.Info("Getting dungeon action character object records params >%s<", params)

	r := m.DungeonActionCharacterObjectRepository()

	return r.GetMany(params, operators, forUpdate)
}

// GetDungeonActionCharacterObjectRec -
func (m *Model) GetDungeonActionCharacterObjectRec(recID string, forUpdate bool) (*record.DungeonActionCharacterObject, error) {

	m.Log.Info("Getting dungeon action character object rec ID >%s<", recID)

	r := m.DungeonActionCharacterObjectRepository()

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

// CreateDungeonActionCharacterObjectRec -
func (m *Model) CreateDungeonActionCharacterObjectRec(rec *record.DungeonActionCharacterObject) error {

	m.Log.Info("Creating dungeon action character object rec >%#v<", rec)

	r := m.DungeonActionCharacterObjectRepository()

	err := m.ValidateDungeonActionCharacterObjectRec(rec)
	if err != nil {
		m.Log.Info("Failed model validation >%v<", err)
		return err
	}

	return r.CreateOne(rec)
}

// UpdateDungeonActionCharacterObjectRec -
func (m *Model) UpdateDungeonActionCharacterObjectRec(rec *record.DungeonActionCharacterObject) error {

	m.Log.Info("Updating dungeon action character object rec >%#v<", rec)

	r := m.DungeonActionCharacterObjectRepository()

	err := m.ValidateDungeonActionCharacterObjectRec(rec)
	if err != nil {
		m.Log.Info("Failed model validation >%v<", err)
		return err
	}

	return r.UpdateOne(rec)
}

// DeleteDungeonActionCharacterObjectRec -
func (m *Model) DeleteDungeonActionCharacterObjectRec(recID string) error {

	m.Log.Info("Deleting dungeon action character object rec ID >%s<", recID)

	r := m.DungeonActionCharacterObjectRepository()

	// validate UUID
	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.ValidateDeleteDungeonActionCharacterObjectRec(recID)
	if err != nil {
		m.Log.Info("Failed model validation >%v<", err)
		return err
	}

	return r.DeleteOne(recID)
}

// RemoveDungeonActionCharacterObjectRec -
func (m *Model) RemoveDungeonActionCharacterObjectRec(recID string) error {

	m.Log.Info("Removing dungeon action character object rec ID >%s<", recID)

	r := m.DungeonActionCharacterObjectRepository()

	// validate UUID
	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.ValidateDeleteDungeonActionCharacterObjectRec(recID)
	if err != nil {
		m.Log.Info("Failed model validation >%v<", err)
		return err
	}

	return r.RemoveOne(recID)
}
