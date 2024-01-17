package model

import (
	"database/sql"
	"fmt"

	coresql "gitlab.com/alienspaces/go-mud/backend/core/sql"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/record"
)

// GetActionMonsterRecs -
func (m *Model) GetActionMonsterRecs(opts *coresql.Options) ([]*record.ActionMonster, error) {

	l := m.loggerWithFunctionContext("GetActionMonsterRecs")

	l.Debug("Getting dungeon action monster records opts >%#v<", opts)

	r := m.ActionMonsterRepository()

	return r.GetMany(opts)
}

// GetActionMonsterRec -
func (m *Model) GetActionMonsterRec(recID string, lock *coresql.Lock) (*record.ActionMonster, error) {

	l := m.loggerWithFunctionContext("GetActionMonsterRec")

	l.Debug("Getting dungeon action monster rec ID >%s<", recID)

	r := m.ActionMonsterRepository()

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

// CreateActionMonsterRec -
func (m *Model) CreateActionMonsterRec(rec *record.ActionMonster) error {

	l := m.loggerWithFunctionContext("CreateActionMonsterRec")

	l.Debug("Creating dungeon action monster record >%#v<", rec)

	r := m.ActionMonsterRepository()

	err := m.validateActionMonsterRec(rec)
	if err != nil {
		l.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.CreateOne(rec)
}

// UpdateActionMonsterRec -
func (m *Model) UpdateActionMonsterRec(rec *record.ActionMonster) error {

	l := m.loggerWithFunctionContext("UpdateActionMonsterRec")

	l.Debug("Updating dungeon action monster record >%#v<", rec)

	r := m.ActionMonsterRepository()

	err := m.validateActionMonsterRec(rec)
	if err != nil {
		l.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.UpdateOne(rec)
}

// DeleteActionMonsterRec -
func (m *Model) DeleteActionMonsterRec(recID string) error {

	l := m.loggerWithFunctionContext("DeleteActionMonsterRec")

	l.Debug("Deleting dungeon action monster rec ID >%s<", recID)

	r := m.ActionMonsterRepository()

	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.validateDeleteActionMonsterRec(recID)
	if err != nil {
		l.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.DeleteOne(recID)
}

// RemoveActionMonsterRec -
func (m *Model) RemoveActionMonsterRec(recID string) error {

	l := m.loggerWithFunctionContext("RemoveActionMonsterRec")

	l.Debug("Removing dungeon action monster rec ID >%s<", recID)

	r := m.ActionMonsterRepository()

	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.validateDeleteActionMonsterRec(recID)
	if err != nil {
		l.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.RemoveOne(recID)
}
