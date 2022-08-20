package model

import (
	"database/sql"
	"fmt"

	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// GetActionMonsterObjectRecs -
func (m *Model) GetActionMonsterObjectRecs(params map[string]interface{}, operators map[string]string, forUpdate bool) ([]*record.ActionMonsterObject, error) {

	l := m.Logger("GetActionMonsterObjectRecs")

	l.Debug("Getting dungeon action monster object records params >%#v<", params)

	r := m.ActionMonsterObjectRepository()

	return r.GetMany(params, operators, forUpdate)
}

// GetActionMonsterObjectRec -
func (m *Model) GetActionMonsterObjectRec(recID string, forUpdate bool) (*record.ActionMonsterObject, error) {

	l := m.Logger("GetActionMonsterObjectRec")

	l.Debug("Getting dungeon action monster object rec ID >%s<", recID)

	r := m.ActionMonsterObjectRepository()

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

// CreateActionMonsterObjectRec -
func (m *Model) CreateActionMonsterObjectRec(rec *record.ActionMonsterObject) error {

	l := m.Logger("CreateActionMonsterObjectRec")

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

	l := m.Logger("UpdateActionMonsterObjectRec")

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

	l := m.Logger("DeleteActionMonsterObjectRec")

	l.Debug("Deleting dungeon action monster object rec ID >%s<", recID)

	r := m.ActionMonsterObjectRepository()

	// validate UUID
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

	l := m.Logger("RemoveActionMonsterObjectRec")

	l.Debug("Removing dungeon action monster object rec ID >%s<", recID)

	r := m.ActionMonsterObjectRepository()

	// validate UUID
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
