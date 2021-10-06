package model

import (
	"database/sql"
	"fmt"

	"gitlab.com/alienspaces/go-boilerplate/server/service/character/internal/record"
)

// GetCharacterRecs -
func (m *Model) GetCharacterRecs(params map[string]interface{}, operators map[string]string, forUpdate bool) ([]*record.Character, error) {

	m.Log.Info("Getting character records params >%s<", params)

	r := m.CharacterRepository()

	return r.GetMany(params, operators, forUpdate)
}

// GetCharacterRec -
func (m *Model) GetCharacterRec(recID string, forUpdate bool) (*record.Character, error) {

	m.Log.Info("Getting character rec ID >%s<", recID)

	r := m.CharacterRepository()

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

// CreateCharacterRec -
func (m *Model) CreateCharacterRec(rec *record.Character) error {

	m.Log.Info("Creating character rec >%#v<", rec)

	r := m.CharacterRepository()

	// Defaults
	rec.AttributePoints = startingAttributePoints

	err := m.ValidateCharacterRec(rec)
	if err != nil {
		m.Log.Info("Failed model validation >%v<", err)
		return err
	}

	return r.CreateOne(rec)
}

// UpdateCharacterRec -
func (m *Model) UpdateCharacterRec(rec *record.Character) error {

	m.Log.Info("Updating character rec >%#v<", rec)

	r := m.CharacterRepository()

	err := m.ValidateCharacterRec(rec)
	if err != nil {
		m.Log.Info("Failed model validation >%v<", err)
		return err
	}

	return r.UpdateOne(rec)
}

// DeleteCharacterRec -
func (m *Model) DeleteCharacterRec(recID string) error {

	m.Log.Info("Deleting character rec ID >%s<", recID)

	r := m.CharacterRepository()

	// validate UUID
	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.ValidateDeleteCharacterRec(recID)
	if err != nil {
		m.Log.Info("Failed model validation >%v<", err)
		return err
	}

	return r.DeleteOne(recID)
}

// RemoveCharacterRec -
func (m *Model) RemoveCharacterRec(recID string) error {

	m.Log.Info("Removing character rec ID >%s<", recID)

	r := m.CharacterRepository()

	// validate UUID
	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.ValidateDeleteCharacterRec(recID)
	if err != nil {
		m.Log.Info("Failed model validation >%v<", err)
		return err
	}

	return r.RemoveOne(recID)
}
