package model

import (
	"database/sql"
	"fmt"

	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// GetDungeonActionMonsterRecs -
func (m *Model) GetDungeonActionMonsterRecs(params map[string]interface{}, operators map[string]string, forUpdate bool) ([]*record.DungeonActionMonster, error) {

	m.Log.Info("Getting game records params >%s<", params)

	r := m.DungeonActionMonsterRepository()

	return r.GetMany(params, operators, forUpdate)
}

// GetDungeonActionMonsterRec -
func (m *Model) GetDungeonActionMonsterRec(recID string, forUpdate bool) (*record.DungeonActionMonster, error) {

	m.Log.Info("Getting game rec ID >%s<", recID)

	r := m.DungeonActionMonsterRepository()

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

// CreateDungeonActionMonsterRec -
func (m *Model) CreateDungeonActionMonsterRec(rec *record.DungeonActionMonster) error {

	m.Log.Info("Creating game rec >%#v<", rec)

	r := m.DungeonActionMonsterRepository()

	err := m.ValidateDungeonActionMonsterRec(rec)
	if err != nil {
		m.Log.Info("Failed model validation >%v<", err)
		return err
	}

	return r.CreateOne(rec)
}

// UpdateDungeonActionMonsterRec -
func (m *Model) UpdateDungeonActionMonsterRec(rec *record.DungeonActionMonster) error {

	m.Log.Info("Updating game rec >%#v<", rec)

	r := m.DungeonActionMonsterRepository()

	err := m.ValidateDungeonActionMonsterRec(rec)
	if err != nil {
		m.Log.Info("Failed model validation >%v<", err)
		return err
	}

	return r.UpdateOne(rec)
}

// DeleteDungeonActionMonsterRec -
func (m *Model) DeleteDungeonActionMonsterRec(recID string) error {

	m.Log.Info("Deleting game rec ID >%s<", recID)

	r := m.DungeonActionMonsterRepository()

	// validate UUID
	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.ValidateDeleteDungeonActionMonsterRec(recID)
	if err != nil {
		m.Log.Info("Failed model validation >%v<", err)
		return err
	}

	return r.DeleteOne(recID)
}

// RemoveDungeonActionMonsterRec -
func (m *Model) RemoveDungeonActionMonsterRec(recID string) error {

	m.Log.Info("Removing game rec ID >%s<", recID)

	r := m.DungeonActionMonsterRepository()

	// validate UUID
	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.ValidateDeleteDungeonActionMonsterRec(recID)
	if err != nil {
		m.Log.Info("Failed model validation >%v<", err)
		return err
	}

	return r.RemoveOne(recID)
}
