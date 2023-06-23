package model

import (
	"database/sql"
	"fmt"

	coresql "gitlab.com/alienspaces/go-mud/backend/core/sql"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/record"
)

// GetActionCharacterObjectRecs -
func (m *Model) GetActionCharacterObjectRecs(opts *coresql.Options) ([]*record.ActionCharacterObject, error) {

	l := m.Logger("GetActionCharacterObjectRecs")

	l.Debug("Getting dungeon action character object records opts >%#v<", opts)

	r := m.ActionCharacterObjectRepository()

	return r.GetMany(opts)
}

// GetActionCharacterObjectRec -
func (m *Model) GetActionCharacterObjectRec(recID string, lock *coresql.Lock) (*record.ActionCharacterObject, error) {

	l := m.Logger("GetActionCharacterObjectRec")

	l.Debug("Getting dungeon action character object rec ID >%s<", recID)

	r := m.ActionCharacterObjectRepository()

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

// CreateActionCharacterObjectRec -
func (m *Model) CreateActionCharacterObjectRec(rec *record.ActionCharacterObject) error {

	l := m.Logger("CreateActionCharacterObjectRec")

	l.Debug("Creating dungeon action character object record >%#v<", rec)

	r := m.ActionCharacterObjectRepository()

	err := m.validateActionCharacterObjectRec(rec)
	if err != nil {
		l.Warn("failed model validation >%v<", err)
		return err
	}

	return r.CreateOne(rec)
}

// UpdateActionCharacterObjectRec -
func (m *Model) UpdateActionCharacterObjectRec(rec *record.ActionCharacterObject) error {

	l := m.Logger("UpdateActionCharacterObjectRec")

	l.Debug("Updating dungeon action character object record >%#v<", rec)

	r := m.ActionCharacterObjectRepository()

	err := m.validateActionCharacterObjectRec(rec)
	if err != nil {
		l.Warn("failed model validation >%v<", err)
		return err
	}

	return r.UpdateOne(rec)
}

// DeleteActionCharacterObjectRec -
func (m *Model) DeleteActionCharacterObjectRec(recID string) error {

	l := m.Logger("DeleteActionCharacterObjectRec")

	l.Debug("Deleting dungeon action character object rec ID >%s<", recID)

	r := m.ActionCharacterObjectRepository()

	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.validateDeleteActionCharacterObjectRec(recID)
	if err != nil {
		l.Warn("failed model validation >%v<", err)
		return err
	}

	return r.DeleteOne(recID)
}

// RemoveActionCharacterObjectRec -
func (m *Model) RemoveActionCharacterObjectRec(recID string) error {

	l := m.Logger("RemoveActionCharacterObjectRec")

	l.Debug("Removing dungeon action character object rec ID >%s<", recID)

	r := m.ActionCharacterObjectRepository()

	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.validateDeleteActionCharacterObjectRec(recID)
	if err != nil {
		l.Warn("failed model validation >%v<", err)
		return err
	}

	return r.RemoveOne(recID)
}
