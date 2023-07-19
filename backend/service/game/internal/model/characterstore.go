package model

import (
	"database/sql"
	"fmt"
	"strings"

	coresql "gitlab.com/alienspaces/go-mud/backend/core/sql"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/calculator"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/record"
)

const defaultCoins = 100
const defaultExperiencePoints = 0
const defaultAttributePoints = 36

// GetCharacterRecs -
func (m *Model) GetCharacterRecs(opts *coresql.Options) ([]*record.Character, error) {

	l := m.Logger("GetCharacterRecs")

	l.Debug("Getting dungeon character records opts >%#v<", opts)

	r := m.CharacterRepository()

	return r.GetMany(opts)
}

// GetCharacterRec -
func (m *Model) GetCharacterRec(recID string, lock *coresql.Lock) (*record.Character, error) {

	l := m.Logger("GetCharacterRec")

	l.Debug("Getting dungeon character rec ID >%s<", recID)

	r := m.CharacterRepository()

	if !m.IsUUID(recID) {
		return nil, fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	rec, err := r.GetOne(recID, lock)
	if err == sql.ErrNoRows {
		l.Warn("No record found ID >%s<", recID)
		return nil, nil
	}

	return rec, err
}

// CreateCharacterRec -
func (m *Model) CreateCharacterRec(rec *record.Character) error {

	l := m.Logger("CreateCharacterRec")

	l.Debug("Creating dungeon character record >%#v<", rec)

	r := m.CharacterRepository()

	rec.AttributePoints = defaultAttributePoints - (rec.Strength + rec.Dexterity + rec.Intelligence)
	rec.ExperiencePoints = defaultExperiencePoints
	rec.Coins = defaultCoins

	rec, err := calculator.CalculateCharacterHealth(rec)
	if err != nil {
		l.Debug("Failed calculating character health >%v<", err)
		return err
	}

	rec, err = calculator.CalculateCharacterFatigue(rec)
	if err != nil {
		l.Debug("Failed calculating character fatigue >%v<", err)
		return err
	}

	err = m.validateCharacterRec(rec)
	if err != nil {
		l.Debug("Failed model validation >%v<", err)
		return err
	}

	err = r.CreateOne(rec)
	if err != nil {
		if strings.Contains(err.Error(), "character_name_uq") {
			return NewCharacterNameTakenError(rec)
		}
	}

	return nil
}

// UpdateCharacterRec -
func (m *Model) UpdateCharacterRec(rec *record.Character) error {

	l := m.Logger("UpdateCharacterRec")

	l.Debug("Updating dungeon character record >%#v<", rec)

	r := m.CharacterRepository()

	err := m.validateCharacterRec(rec)
	if err != nil {
		l.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.UpdateOne(rec)
}

// DeleteCharacterRec -
func (m *Model) DeleteCharacterRec(recID string) error {

	l := m.Logger("DeleteCharacterRec")

	l.Debug("Deleting dungeon character rec ID >%s<", recID)

	r := m.CharacterRepository()

	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.validateDeleteCharacterRec(recID)
	if err != nil {
		l.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.DeleteOne(recID)
}

// RemoveCharacterRec -
func (m *Model) RemoveCharacterRec(recID string) error {

	l := m.Logger("RemoveCharacterRec")

	l.Debug("Removing dungeon character rec ID >%s<", recID)

	r := m.CharacterRepository()

	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.validateDeleteCharacterRec(recID)
	if err != nil {
		l.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.RemoveOne(recID)
}
