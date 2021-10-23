package model

import (
	"database/sql"
	"fmt"

	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// GetDungeonObjectRecs -
func (m *Model) GetDungeonObjectRecs(params map[string]interface{}, operators map[string]string, forUpdate bool) ([]*record.DungeonObject, error) {

	m.Log.Info("Getting dungeon object records params >%s<", params)

	r := m.DungeonObjectRepository()

	return r.GetMany(params, operators, forUpdate)
}

// GetDungeonObjectRec -
func (m *Model) GetDungeonObjectRec(recID string, forUpdate bool) (*record.DungeonObject, error) {

	m.Log.Info("Getting dungeon object rec ID >%s<", recID)

	r := m.DungeonObjectRepository()

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

// CreateDungeonObjectRec -
func (m *Model) CreateDungeonObjectRec(rec *record.DungeonObject) error {

	m.Log.Info("Creating dungeon object rec >%#v<", rec)

	r := m.DungeonObjectRepository()

	err := m.ValidateDungeonObjectRec(rec)
	if err != nil {
		m.Log.Info("Failed model validation >%v<", err)
		return err
	}

	return r.CreateOne(rec)
}

// UpdateDungeonObjectRec -
func (m *Model) UpdateDungeonObjectRec(rec *record.DungeonObject) error {

	m.Log.Info("Updating dungeon object rec >%#v<", rec)

	r := m.DungeonObjectRepository()

	err := m.ValidateDungeonObjectRec(rec)
	if err != nil {
		m.Log.Info("Failed model validation >%v<", err)
		return err
	}

	return r.UpdateOne(rec)
}

// DeleteDungeonObjectRec -
func (m *Model) DeleteDungeonObjectRec(recID string) error {

	m.Log.Info("Deleting dungeon object rec ID >%s<", recID)

	r := m.DungeonObjectRepository()

	// validate UUID
	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.ValidateDeleteDungeonObjectRec(recID)
	if err != nil {
		m.Log.Info("Failed model validation >%v<", err)
		return err
	}

	return r.DeleteOne(recID)
}

// RemoveDungeonObjectRec -
func (m *Model) RemoveDungeonObjectRec(recID string) error {

	m.Log.Info("Removing dungeon object rec ID >%s<", recID)

	r := m.DungeonObjectRepository()

	// validate UUID
	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.ValidateDeleteDungeonObjectRec(recID)
	if err != nil {
		m.Log.Info("Failed model validation >%v<", err)
		return err
	}

	return r.RemoveOne(recID)
}
