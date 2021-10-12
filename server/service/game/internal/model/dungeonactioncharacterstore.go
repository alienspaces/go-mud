package model

import (
	"database/sql"
	"fmt"

	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// GetDungeonActionCharacterRecs -
func (m *Model) GetDungeonActionCharacterRecs(params map[string]interface{}, operators map[string]string, forUpdate bool) ([]*record.DungeonActionCharacter, error) {

	m.Log.Info("Getting game records params >%s<", params)

	r := m.DungeonActionCharacterRepository()

	return r.GetMany(params, operators, forUpdate)
}

// GetDungeonActionCharacterRec -
func (m *Model) GetDungeonActionCharacterRec(recID string, forUpdate bool) (*record.DungeonActionCharacter, error) {

	m.Log.Info("Getting game rec ID >%s<", recID)

	r := m.DungeonActionCharacterRepository()

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

// CreateDungeonActionCharacterRec -
func (m *Model) CreateDungeonActionCharacterRec(rec *record.DungeonActionCharacter) error {

	m.Log.Info("Creating game rec >%#v<", rec)

	r := m.DungeonActionCharacterRepository()

	err := m.ValidateDungeonActionCharacterRec(rec)
	if err != nil {
		m.Log.Info("Failed model validation >%v<", err)
		return err
	}

	return r.CreateOne(rec)
}

// UpdateDungeonActionCharacterRec -
func (m *Model) UpdateDungeonActionCharacterRec(rec *record.DungeonActionCharacter) error {

	m.Log.Info("Updating game rec >%#v<", rec)

	r := m.DungeonActionCharacterRepository()

	err := m.ValidateDungeonActionCharacterRec(rec)
	if err != nil {
		m.Log.Info("Failed model validation >%v<", err)
		return err
	}

	return r.UpdateOne(rec)
}

// DeleteDungeonActionCharacterRec -
func (m *Model) DeleteDungeonActionCharacterRec(recID string) error {

	m.Log.Info("Deleting game rec ID >%s<", recID)

	r := m.DungeonActionCharacterRepository()

	// validate UUID
	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.ValidateDeleteDungeonActionCharacterRec(recID)
	if err != nil {
		m.Log.Info("Failed model validation >%v<", err)
		return err
	}

	return r.DeleteOne(recID)
}

// RemoveDungeonActionCharacterRec -
func (m *Model) RemoveDungeonActionCharacterRec(recID string) error {

	m.Log.Info("Removing game rec ID >%s<", recID)

	r := m.DungeonActionCharacterRepository()

	// validate UUID
	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.ValidateDeleteDungeonActionCharacterRec(recID)
	if err != nil {
		m.Log.Info("Failed model validation >%v<", err)
		return err
	}

	return r.RemoveOne(recID)
}
