package model

import (
	"database/sql"
	"fmt"

	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// GetDungeonLocationRecs -
func (m *Model) GetDungeonLocationRecs(params map[string]interface{}, operators map[string]string, forUpdate bool) ([]*record.DungeonLocation, error) {

	m.Log.Debug("Getting dungeon location records params >%s<", params)

	r := m.DungeonLocationRepository()

	return r.GetMany(params, operators, forUpdate)
}

// GetDungeonLocationRec -
func (m *Model) GetDungeonLocationRec(recID string, forUpdate bool) (*record.DungeonLocation, error) {

	m.Log.Debug("Getting dungeon location rec ID >%s<", recID)

	r := m.DungeonLocationRepository()

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

// CreateDungeonLocationRec -
func (m *Model) CreateDungeonLocationRec(rec *record.DungeonLocation) error {

	m.Log.Debug("Creating dungeon location rec >%#v<", rec)

	r := m.DungeonLocationRepository()

	err := m.ValidateDungeonLocationRec(rec)
	if err != nil {
		m.Log.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.CreateOne(rec)
}

// UpdateDungeonLocationRec -
func (m *Model) UpdateDungeonLocationRec(rec *record.DungeonLocation) error {

	m.Log.Debug("Updating dungeon location rec >%#v<", rec)

	r := m.DungeonLocationRepository()

	err := m.ValidateDungeonLocationRec(rec)
	if err != nil {
		m.Log.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.UpdateOne(rec)
}

// DeleteDungeonLocationRec -
func (m *Model) DeleteDungeonLocationRec(recID string) error {

	m.Log.Debug("Deleting dungeon location rec ID >%s<", recID)

	r := m.DungeonLocationRepository()

	// validate UUID
	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.ValidateDeleteDungeonLocationRec(recID)
	if err != nil {
		m.Log.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.DeleteOne(recID)
}

// RemoveDungeonLocationRec -
func (m *Model) RemoveDungeonLocationRec(recID string) error {

	m.Log.Debug("Removing dungeon location rec ID >%s<", recID)

	r := m.DungeonLocationRepository()

	// validate UUID
	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.ValidateDeleteDungeonLocationRec(recID)
	if err != nil {
		m.Log.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.RemoveOne(recID)
}
