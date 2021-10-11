package model

import (
	"database/sql"
	"fmt"

	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// GetDungeonCharacterRecs -
func (m *Model) GetDungeonCharacterRecs(params map[string]interface{}, operators map[string]string, forUpdate bool) ([]*record.DungeonCharacter, error) {

	m.Log.Info("Getting game records params >%s<", params)

	r := m.DungeonCharacterRepository()

	return r.GetMany(params, operators, forUpdate)
}

// GetDungeonCharacterRec -
func (m *Model) GetDungeonCharacterRec(recID string, forUpdate bool) (*record.DungeonCharacter, error) {

	m.Log.Info("Getting game rec ID >%s<", recID)

	r := m.DungeonCharacterRepository()

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

// CreateDungeonCharacterRec -
func (m *Model) CreateDungeonCharacterRec(rec *record.DungeonCharacter) error {

	m.Log.Info("Creating game rec >%#v<", rec)

	r := m.DungeonCharacterRepository()

	err := m.ValidateDungeonCharacterRec(rec)
	if err != nil {
		m.Log.Info("Failed model validation >%v<", err)
		return err
	}

	return r.CreateOne(rec)
}

// UpdateDungeonCharacterRec -
func (m *Model) UpdateDungeonCharacterRec(rec *record.DungeonCharacter) error {

	m.Log.Info("Updating game rec >%#v<", rec)

	r := m.DungeonCharacterRepository()

	err := m.ValidateDungeonCharacterRec(rec)
	if err != nil {
		m.Log.Info("Failed model validation >%v<", err)
		return err
	}

	return r.UpdateOne(rec)
}

// DeleteDungeonCharacterRec -
func (m *Model) DeleteDungeonCharacterRec(recID string) error {

	m.Log.Info("Deleting game rec ID >%s<", recID)

	r := m.DungeonCharacterRepository()

	// validate UUID
	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.ValidateDeleteDungeonCharacterRec(recID)
	if err != nil {
		m.Log.Info("Failed model validation >%v<", err)
		return err
	}

	return r.DeleteOne(recID)
}

// RemoveDungeonCharacterRec -
func (m *Model) RemoveDungeonCharacterRec(recID string) error {

	m.Log.Info("Removing game rec ID >%s<", recID)

	r := m.DungeonCharacterRepository()

	// validate UUID
	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.ValidateDeleteDungeonCharacterRec(recID)
	if err != nil {
		m.Log.Info("Failed model validation >%v<", err)
		return err
	}

	return r.RemoveOne(recID)
}
