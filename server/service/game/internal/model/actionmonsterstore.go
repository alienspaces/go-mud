package model

import (
	"database/sql"
	"fmt"

	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// GetActionMonsterRecs -
func (m *Model) GetActionMonsterRecs(params map[string]interface{}, operators map[string]string, forUpdate bool) ([]*record.ActionMonster, error) {

	m.Log.Debug("Getting dungeon action monster records params >%s<", params)

	r := m.ActionMonsterRepository()

	return r.GetMany(params, operators, forUpdate)
}

// GetActionMonsterRec -
func (m *Model) GetActionMonsterRec(recID string, forUpdate bool) (*record.ActionMonster, error) {

	m.Log.Debug("Getting dungeon action monster rec ID >%s<", recID)

	r := m.ActionMonsterRepository()

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

// CreateActionMonsterRec -
func (m *Model) CreateActionMonsterRec(rec *record.ActionMonster) error {

	m.Log.Info("Creating dungeon action monster rec >%#v<", rec)

	r := m.ActionMonsterRepository()

	err := m.ValidateActionMonsterRec(rec)
	if err != nil {
		m.Log.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.CreateOne(rec)
}

// UpdateActionMonsterRec -
func (m *Model) UpdateActionMonsterRec(rec *record.ActionMonster) error {

	m.Log.Debug("Updating dungeon action monster rec >%#v<", rec)

	r := m.ActionMonsterRepository()

	err := m.ValidateActionMonsterRec(rec)
	if err != nil {
		m.Log.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.UpdateOne(rec)
}

// DeleteActionMonsterRec -
func (m *Model) DeleteActionMonsterRec(recID string) error {

	m.Log.Debug("Deleting dungeon action monster rec ID >%s<", recID)

	r := m.ActionMonsterRepository()

	// validate UUID
	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.ValidateDeleteActionMonsterRec(recID)
	if err != nil {
		m.Log.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.DeleteOne(recID)
}

// RemoveActionMonsterRec -
func (m *Model) RemoveActionMonsterRec(recID string) error {

	m.Log.Debug("Removing dungeon action monster rec ID >%s<", recID)

	r := m.ActionMonsterRepository()

	// validate UUID
	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.ValidateDeleteActionMonsterRec(recID)
	if err != nil {
		m.Log.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.RemoveOne(recID)
}
