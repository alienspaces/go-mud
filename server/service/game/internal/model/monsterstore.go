package model

import (
	"database/sql"
	"fmt"

	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// GetMonsterRecs -
func (m *Model) GetMonsterRecs(params map[string]interface{}, operators map[string]string, forUpdate bool) ([]*record.Monster, error) {

	m.Log.Debug("Getting dungeon monster records params >%s<", params)

	r := m.MonsterRepository()

	return r.GetMany(params, operators, forUpdate)
}

// GetMonsterRec -
func (m *Model) GetMonsterRec(recID string, forUpdate bool) (*record.Monster, error) {

	m.Log.Debug("Getting dungeon monster rec ID >%s<", recID)

	r := m.MonsterRepository()

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

// CreateMonsterRec -
func (m *Model) CreateMonsterRec(rec *record.Monster) error {

	m.Log.Debug("Creating dungeon monster rec >%#v<", rec)

	r := m.MonsterRepository()

	rec.AttributePoints = defaultAttributePoints - (rec.Strength + rec.Dexterity + rec.Intelligence)
	rec.ExperiencePoints = defaultExperiencePoints
	rec.Health = m.calculateHealth(rec.Strength, rec.Dexterity)
	rec.Fatigue = m.calculateFatigue(rec.Strength, rec.Intelligence)
	rec.Coins = defaultCoins

	err := m.ValidateMonsterRec(rec)
	if err != nil {
		m.Log.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.CreateOne(rec)
}

// UpdateMonsterRec -
func (m *Model) UpdateMonsterRec(rec *record.Monster) error {

	m.Log.Debug("Updating dungeon monster rec >%#v<", rec)

	r := m.MonsterRepository()

	err := m.ValidateMonsterRec(rec)
	if err != nil {
		m.Log.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.UpdateOne(rec)
}

// DeleteMonsterRec -
func (m *Model) DeleteMonsterRec(recID string) error {

	m.Log.Debug("Deleting dungeon monster rec ID >%s<", recID)

	r := m.MonsterRepository()

	// validate UUID
	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.ValidateDeleteMonsterRec(recID)
	if err != nil {
		m.Log.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.DeleteOne(recID)
}

// RemoveMonsterRec -
func (m *Model) RemoveMonsterRec(recID string) error {

	m.Log.Debug("Removing dungeon monster rec ID >%s<", recID)

	r := m.MonsterRepository()

	// validate UUID
	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.ValidateDeleteMonsterRec(recID)
	if err != nil {
		m.Log.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.RemoveOne(recID)
}
