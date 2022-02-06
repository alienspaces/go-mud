package model

import (
	"database/sql"
	"fmt"

	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// GetDungeonActionMonsterObjectRecs -
func (m *Model) GetDungeonActionMonsterObjectRecs(params map[string]interface{}, operators map[string]string, forUpdate bool) ([]*record.DungeonActionMonsterObject, error) {

	m.Log.Info("Getting dungeon action monster object records params >%#v<", params)

	r := m.DungeonActionMonsterObjectRepository()

	return r.GetMany(params, operators, forUpdate)
}

// GetDungeonActionMonsterObjectRec -
func (m *Model) GetDungeonActionMonsterObjectRec(recID string, forUpdate bool) (*record.DungeonActionMonsterObject, error) {

	m.Log.Info("Getting dungeon action monster object rec ID >%s<", recID)

	r := m.DungeonActionMonsterObjectRepository()

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

// CreateDungeonActionMonsterObjectRec -
func (m *Model) CreateDungeonActionMonsterObjectRec(rec *record.DungeonActionMonsterObject) error {

	m.Log.Info("Creating dungeon action monster object rec >%#v<", rec)

	r := m.DungeonActionMonsterObjectRepository()

	err := m.ValidateDungeonActionMonsterObjectRec(rec)
	if err != nil {
		m.Log.Info("Failed model validation >%v<", err)
		return err
	}

	return r.CreateOne(rec)
}

// UpdateDungeonActionMonsterObjectRec -
func (m *Model) UpdateDungeonActionMonsterObjectRec(rec *record.DungeonActionMonsterObject) error {

	m.Log.Info("Updating dungeon action monster object rec >%#v<", rec)

	r := m.DungeonActionMonsterObjectRepository()

	err := m.ValidateDungeonActionMonsterObjectRec(rec)
	if err != nil {
		m.Log.Info("Failed model validation >%v<", err)
		return err
	}

	return r.UpdateOne(rec)
}

// DeleteDungeonActionMonsterObjectRec -
func (m *Model) DeleteDungeonActionMonsterObjectRec(recID string) error {

	m.Log.Info("Deleting dungeon action monster object rec ID >%s<", recID)

	r := m.DungeonActionMonsterObjectRepository()

	// validate UUID
	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.ValidateDeleteDungeonActionMonsterObjectRec(recID)
	if err != nil {
		m.Log.Info("Failed model validation >%v<", err)
		return err
	}

	return r.DeleteOne(recID)
}

// RemoveDungeonActionMonsterObjectRec -
func (m *Model) RemoveDungeonActionMonsterObjectRec(recID string) error {

	m.Log.Info("Removing dungeon action monster object rec ID >%s<", recID)

	r := m.DungeonActionMonsterObjectRepository()

	// validate UUID
	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.ValidateDeleteDungeonActionMonsterObjectRec(recID)
	if err != nil {
		m.Log.Info("Failed model validation >%v<", err)
		return err
	}

	return r.RemoveOne(recID)
}
