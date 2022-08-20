package model

import (
	"database/sql"
	"fmt"

	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// GetMonsterRecs -
func (m *Model) GetMonsterRecs(params map[string]interface{}, operators map[string]string, forUpdate bool) ([]*record.Monster, error) {

	l := m.Logger("GetMonsterRecs")

	l.Debug("Getting dungeon monster records params >%s<", params)

	r := m.MonsterRepository()

	return r.GetMany(params, operators, forUpdate)
}

// GetMonsterRec -
func (m *Model) GetMonsterRec(recID string, forUpdate bool) (*record.Monster, error) {

	l := m.Logger("GetMonsterRec")

	l.Debug("Getting dungeon monster rec ID >%s<", recID)

	r := m.MonsterRepository()

	// validate UUID
	if !m.IsUUID(recID) {
		return nil, fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	rec, err := r.GetOne(recID, forUpdate)
	if err == sql.ErrNoRows {
		l.Warn("No record found ID >%s<", recID)
		return nil, nil
	}

	return rec, err
}

// CreateMonsterRec -
func (m *Model) CreateMonsterRec(rec *record.Monster) error {

	l := m.Logger("CreateMonsterRec")

	l.Debug("Creating dungeon monster record >%#v<", rec)

	r := m.MonsterRepository()

	rec.AttributePoints = defaultAttributePoints - (rec.Strength + rec.Dexterity + rec.Intelligence)
	rec.ExperiencePoints = defaultExperiencePoints
	rec.Health = m.calculateHealth(rec.Strength, rec.Dexterity)
	rec.Fatigue = m.calculateFatigue(rec.Strength, rec.Intelligence)
	rec.Coins = defaultCoins

	err := m.validateMonsterRec(rec)
	if err != nil {
		l.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.CreateOne(rec)
}

// UpdateMonsterRec -
func (m *Model) UpdateMonsterRec(rec *record.Monster) error {

	l := m.Logger("UpdateMonsterRec")

	l.Debug("Updating dungeon monster record >%#v<", rec)

	r := m.MonsterRepository()

	err := m.validateMonsterRec(rec)
	if err != nil {
		l.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.UpdateOne(rec)
}

// DeleteMonsterRec -
func (m *Model) DeleteMonsterRec(recID string) error {

	l := m.Logger("DeleteMonsterRec")

	l.Debug("Deleting dungeon monster rec ID >%s<", recID)

	r := m.MonsterRepository()

	// validate UUID
	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.validateDeleteMonsterRec(recID)
	if err != nil {
		l.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.DeleteOne(recID)
}

// RemoveMonsterRec -
func (m *Model) RemoveMonsterRec(recID string) error {

	l := m.Logger("RemoveMonsterRec")

	l.Debug("Removing dungeon monster rec ID >%s<", recID)

	r := m.MonsterRepository()

	// validate UUID
	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.validateDeleteMonsterRec(recID)
	if err != nil {
		l.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.RemoveOne(recID)
}
