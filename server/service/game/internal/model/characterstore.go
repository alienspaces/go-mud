package model

import (
	"database/sql"
	"fmt"

	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

const defaultCoins = 100
const defaultExperiencePoints = 0
const defaultAttributePoints = 36

// GetCharacterRecs -
func (m *Model) GetCharacterRecs(params map[string]interface{}, operators map[string]string, forUpdate bool) ([]*record.Character, error) {

	m.Log.Debug("Getting dungeon character records params >%s<", params)

	r := m.CharacterRepository()

	return r.GetMany(params, operators, forUpdate)
}

// GetCharacterRec -
func (m *Model) GetCharacterRec(recID string, forUpdate bool) (*record.Character, error) {

	m.Log.Debug("Getting dungeon character rec ID >%s<", recID)

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

	m.Log.Debug("Creating dungeon character rec >%#v<", rec)

	characterRepo := m.CharacterRepository()

	rec.AttributePoints = defaultAttributePoints - (rec.Strength + rec.Dexterity + rec.Intelligence)
	rec.ExperiencePoints = defaultExperiencePoints
	rec.Health = m.calculateHealth(rec.Strength, rec.Dexterity)
	rec.Fatigue = m.calculateFatigue(rec.Strength, rec.Intelligence)
	rec.Coins = defaultCoins

	err := m.ValidateCharacterRec(rec)
	if err != nil {
		m.Log.Debug("Failed model validation >%v<", err)
		return err
	}

	return characterRepo.CreateOne(rec)
}

// UpdateCharacterRec -
func (m *Model) UpdateCharacterRec(rec *record.Character) error {

	m.Log.Debug("Updating dungeon character rec >%#v<", rec)

	r := m.CharacterRepository()

	err := m.ValidateCharacterRec(rec)
	if err != nil {
		m.Log.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.UpdateOne(rec)
}

// DeleteCharacterRec -
func (m *Model) DeleteCharacterRec(recID string) error {

	m.Log.Debug("Deleting dungeon character rec ID >%s<", recID)

	r := m.CharacterRepository()

	// validate UUID
	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.ValidateDeleteCharacterRec(recID)
	if err != nil {
		m.Log.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.DeleteOne(recID)
}

// RemoveCharacterRec -
func (m *Model) RemoveCharacterRec(recID string) error {

	m.Log.Debug("Removing dungeon character rec ID >%s<", recID)

	r := m.CharacterRepository()

	// validate UUID
	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.ValidateDeleteCharacterRec(recID)
	if err != nil {
		m.Log.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.RemoveOne(recID)
}

func (m *Model) calculateHealth(strength int, dexterity int) int {
	return strength + dexterity*10
}

func (m *Model) calculateFatigue(strength int, intelligence int) int {
	return strength + intelligence*10
}
