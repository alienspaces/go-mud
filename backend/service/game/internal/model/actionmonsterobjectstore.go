package model

import (
	"database/sql"
	"fmt"

	coresql "gitlab.com/alienspaces/go-mud/backend/core/sql"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/record"
)

// GetActionMonsterObjectRecs -
func (m *Model) GetActionMonsterObjectRecs(opts *coresql.Options) ([]*record.ActionMonsterObject, error) {

	l := m.loggerWithContext("GetActionMonsterObjectRecs")

	l.Debug("Getting dungeon action monster object records opts >%#v<", opts)

	r := m.ActionMonsterObjectRepository()

	return r.GetMany(opts)
}

// GetActionMonsterObjectRec -
func (m *Model) GetActionMonsterObjectRec(recID string, lock *coresql.Lock) (*record.ActionMonsterObject, error) {

	l := m.loggerWithContext("GetActionMonsterObjectRec")

	l.Debug("Getting dungeon action monster object rec ID >%s<", recID)

	r := m.ActionMonsterObjectRepository()

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

// CreateActionMonsterObjectRec -
func (m *Model) CreateActionMonsterObjectRec(rec *record.ActionMonsterObject) error {

	l := m.loggerWithContext("CreateActionMonsterObjectRec")

	l.Debug("Creating dungeon action monster object record >%#v<", rec)

	r := m.ActionMonsterObjectRepository()

	err := m.validateActionMonsterObjectRec(rec)
	if err != nil {
		l.Warn("failed model validation >%v<", err)
		return err
	}

	return r.CreateOne(rec)
}

// UpdateActionMonsterObjectRec -
func (m *Model) UpdateActionMonsterObjectRec(rec *record.ActionMonsterObject) error {

	l := m.loggerWithContext("UpdateActionMonsterObjectRec")

	l.Debug("Updating dungeon action monster object record >%#v<", rec)

	r := m.ActionMonsterObjectRepository()

	err := m.validateActionMonsterObjectRec(rec)
	if err != nil {
		l.Warn("failed model validation >%v<", err)
		return err
	}

	return r.UpdateOne(rec)
}

// DeleteActionMonsterObjectRec -
func (m *Model) DeleteActionMonsterObjectRec(recID string) error {

	l := m.loggerWithContext("DeleteActionMonsterObjectRec")

	l.Debug("Deleting dungeon action monster object rec ID >%s<", recID)

	r := m.ActionMonsterObjectRepository()

	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.validateDeleteActionMonsterObjectRec(recID)
	if err != nil {
		l.Warn("failed model validation >%v<", err)
		return err
	}

	return r.DeleteOne(recID)
}

// RemoveActionMonsterObjectRec -
func (m *Model) RemoveActionMonsterObjectRec(recID string) error {

	l := m.loggerWithContext("RemoveActionMonsterObjectRec")

	l.Debug("Removing dungeon action monster object rec ID >%s<", recID)

	r := m.ActionMonsterObjectRepository()

	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.validateDeleteActionMonsterObjectRec(recID)
	if err != nil {
		l.Warn("failed model validation >%v<", err)
		return err
	}

	return r.RemoveOne(recID)
}
