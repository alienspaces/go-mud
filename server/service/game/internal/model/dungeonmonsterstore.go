package model

import (
	"database/sql"
	"fmt"

	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// GetDungeonMonsterRecs -
func (m *Model) GetDungeonMonsterRecs(params map[string]interface{}, operators map[string]string, forUpdate bool) ([]*record.DungeonMonster, error) {

	m.Log.Debug("Getting dungeon monster records params >%s<", params)

	r := m.DungeonMonsterRepository()

	return r.GetMany(params, operators, forUpdate)
}

// GetDungeonMonsterRec -
func (m *Model) GetDungeonMonsterRec(recID string, forUpdate bool) (*record.DungeonMonster, error) {

	m.Log.Debug("Getting dungeon monster rec ID >%s<", recID)

	r := m.DungeonMonsterRepository()

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

// CreateDungeonMonsterRec -
func (m *Model) CreateDungeonMonsterRec(rec *record.DungeonMonster) error {

	m.Log.Debug("Creating dungeon monster rec >%#v<", rec)

	r := m.DungeonMonsterRepository()

	rec.AttributePoints = defaultAttributePoints - (rec.Strength + rec.Dexterity + rec.Intelligence)
	rec.ExperiencePoints = defaultExperiencePoints
	rec.Health = m.calculateHealth(rec.Strength, rec.Dexterity)
	rec.Fatigue = m.calculateFatigue(rec.Strength, rec.Intelligence)
	rec.Coins = defaultCoins

	err := m.ValidateDungeonMonsterRec(rec)
	if err != nil {
		m.Log.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.CreateOne(rec)
}

// UpdateDungeonMonsterRec -
func (m *Model) UpdateDungeonMonsterRec(rec *record.DungeonMonster) error {

	m.Log.Debug("Updating dungeon monster rec >%#v<", rec)

	r := m.DungeonMonsterRepository()

	err := m.ValidateDungeonMonsterRec(rec)
	if err != nil {
		m.Log.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.UpdateOne(rec)
}

// DeleteDungeonMonsterRec -
func (m *Model) DeleteDungeonMonsterRec(recID string) error {

	m.Log.Debug("Deleting dungeon monster rec ID >%s<", recID)

	r := m.DungeonMonsterRepository()

	// validate UUID
	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.ValidateDeleteDungeonMonsterRec(recID)
	if err != nil {
		m.Log.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.DeleteOne(recID)
}

// RemoveDungeonMonsterRec -
func (m *Model) RemoveDungeonMonsterRec(recID string) error {

	m.Log.Debug("Removing dungeon monster rec ID >%s<", recID)

	r := m.DungeonMonsterRepository()

	// validate UUID
	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.ValidateDeleteDungeonMonsterRec(recID)
	if err != nil {
		m.Log.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.RemoveOne(recID)
}
