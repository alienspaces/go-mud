package model

import (
	"database/sql"
	"fmt"

	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// GetDungeonRecs -
func (m *Model) GetDungeonRecs(params map[string]interface{}, operators map[string]string, forUpdate bool) ([]*record.Dungeon, error) {

	m.Log.Info("Getting game records params >%s<", params)

	r := m.DungeonRepository()

	return r.GetMany(params, operators, forUpdate)
}

// GetDungeonRec -
func (m *Model) GetDungeonRec(recID string, forUpdate bool) (*record.Dungeon, error) {

	m.Log.Info("Getting game rec ID >%s<", recID)

	r := m.DungeonRepository()

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

// CreateDungeonRec -
func (m *Model) CreateDungeonRec(rec *record.Dungeon) error {

	m.Log.Info("Creating game rec >%#v<", rec)

	r := m.DungeonRepository()

	err := m.ValidateDungeonRec(rec)
	if err != nil {
		m.Log.Info("Failed model validation >%v<", err)
		return err
	}

	return r.CreateOne(rec)
}

// UpdateDungeonRec -
func (m *Model) UpdateDungeonRec(rec *record.Dungeon) error {

	m.Log.Info("Updating game rec >%#v<", rec)

	r := m.DungeonRepository()

	err := m.ValidateDungeonRec(rec)
	if err != nil {
		m.Log.Info("Failed model validation >%v<", err)
		return err
	}

	return r.UpdateOne(rec)
}

// DeleteDungeonRec -
func (m *Model) DeleteDungeonRec(recID string) error {

	m.Log.Info("Deleting game rec ID >%s<", recID)

	r := m.DungeonRepository()

	// validate UUID
	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.ValidateDeleteDungeonRec(recID)
	if err != nil {
		m.Log.Info("Failed model validation >%v<", err)
		return err
	}

	return r.DeleteOne(recID)
}

// RemoveDungeonRec -
func (m *Model) RemoveDungeonRec(recID string) error {

	m.Log.Info("Removing game rec ID >%s<", recID)

	r := m.DungeonRepository()

	// validate UUID
	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.ValidateDeleteDungeonRec(recID)
	if err != nil {
		m.Log.Info("Failed model validation >%v<", err)
		return err
	}

	return r.RemoveOne(recID)
}
