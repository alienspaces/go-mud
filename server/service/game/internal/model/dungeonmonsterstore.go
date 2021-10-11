package model

import (
	"database/sql"
	"fmt"

	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// GetDungeonMonsterRecs -
func (m *Model) GetDungeonMonsterRecs(params map[string]interface{}, operators map[string]string, forUpdate bool) ([]*record.DungeonMonster, error) {

	m.Log.Info("Getting game records params >%s<", params)

	r := m.DungeonMonsterRepository()

	return r.GetMany(params, operators, forUpdate)
}

// GetDungeonMonsterRec -
func (m *Model) GetDungeonMonsterRec(recID string, forUpdate bool) (*record.DungeonMonster, error) {

	m.Log.Info("Getting game rec ID >%s<", recID)

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

	m.Log.Info("Creating game rec >%#v<", rec)

	r := m.DungeonMonsterRepository()

	err := m.ValidateDungeonMonsterRec(rec)
	if err != nil {
		m.Log.Info("Failed model validation >%v<", err)
		return err
	}

	return r.CreateOne(rec)
}

// UpdateDungeonMonsterRec -
func (m *Model) UpdateDungeonMonsterRec(rec *record.DungeonMonster) error {

	m.Log.Info("Updating game rec >%#v<", rec)

	r := m.DungeonMonsterRepository()

	err := m.ValidateDungeonMonsterRec(rec)
	if err != nil {
		m.Log.Info("Failed model validation >%v<", err)
		return err
	}

	return r.UpdateOne(rec)
}

// DeleteDungeonMonsterRec -
func (m *Model) DeleteDungeonMonsterRec(recID string) error {

	m.Log.Info("Deleting game rec ID >%s<", recID)

	r := m.DungeonMonsterRepository()

	// validate UUID
	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.ValidateDeleteDungeonMonsterRec(recID)
	if err != nil {
		m.Log.Info("Failed model validation >%v<", err)
		return err
	}

	return r.DeleteOne(recID)
}

// RemoveDungeonMonsterRec -
func (m *Model) RemoveDungeonMonsterRec(recID string) error {

	m.Log.Info("Removing game rec ID >%s<", recID)

	r := m.DungeonMonsterRepository()

	// validate UUID
	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.ValidateDeleteDungeonMonsterRec(recID)
	if err != nil {
		m.Log.Info("Failed model validation >%v<", err)
		return err
	}

	return r.RemoveOne(recID)
}
