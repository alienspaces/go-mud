package model

import (
	"database/sql"
	"fmt"

	coresql "gitlab.com/alienspaces/go-mud/backend/core/sql"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/calculator"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/record"
)

// GetMonsterRecs -
func (m *Model) GetMonsterRecs(opts *coresql.Options) ([]*record.Monster, error) {

	l := m.loggerWithFunctionContext("GetMonsterRecs")

	l.Debug("Getting dungeon monster records opts >%#v<", opts)

	r := m.MonsterRepository()

	return r.GetMany(opts)
}

// GetMonsterRec -
func (m *Model) GetMonsterRec(recID string, lock *coresql.Lock) (*record.Monster, error) {

	l := m.loggerWithFunctionContext("GetMonsterRec")

	l.Debug("Getting dungeon monster rec ID >%s<", recID)

	r := m.MonsterRepository()

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

// CreateMonsterRec -
func (m *Model) CreateMonsterRec(rec *record.Monster) error {

	l := m.loggerWithFunctionContext("CreateMonsterRec")

	l.Debug("Creating dungeon monster record >%#v<", rec)

	r := m.MonsterRepository()

	rec.AttributePoints = defaultAttributePoints - (rec.Strength + rec.Dexterity + rec.Intelligence)
	rec.ExperiencePoints = defaultExperiencePoints
	rec.Coins = defaultCoins

	rec, err := calculator.CalculateMonsterHealth(rec)
	if err != nil {
		l.Debug("Failed calculating monster health >%v<", err)
		return err
	}

	rec, err = calculator.CalculateMonsterFatigue(rec)
	if err != nil {
		l.Debug("Failed calculating monster fatigue >%v<", err)
		return err
	}

	err = m.validateMonsterRec(rec)
	if err != nil {
		l.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.CreateOne(rec)
}

// UpdateMonsterRec -
func (m *Model) UpdateMonsterRec(rec *record.Monster) error {

	l := m.loggerWithFunctionContext("UpdateMonsterRec")

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

	l := m.loggerWithFunctionContext("DeleteMonsterRec")

	l.Debug("Deleting dungeon monster rec ID >%s<", recID)

	r := m.MonsterRepository()

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

	l := m.loggerWithFunctionContext("RemoveMonsterRec")

	l.Debug("Removing dungeon monster rec ID >%s<", recID)

	r := m.MonsterRepository()

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
