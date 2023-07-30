package model

import (
	"database/sql"
	"fmt"

	coresql "gitlab.com/alienspaces/go-mud/backend/core/sql"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/record"
)

// GetActionCharacterRecs -
func (m *Model) GetActionCharacterRecs(opts *coresql.Options) ([]*record.ActionCharacter, error) {

	l := m.loggerWithFunctionContext("GetActionCharacterRecs")

	l.Debug("Getting dungeon action character records opts >%#v<", opts)

	r := m.ActionCharacterRepository()

	return r.GetMany(opts)
}

// GetActionCharacterRec -
func (m *Model) GetActionCharacterRec(recID string, lock *coresql.Lock) (*record.ActionCharacter, error) {

	l := m.loggerWithFunctionContext("GetActionCharacterRec")

	l.Debug("Getting dungeon action character rec ID >%s<", recID)

	r := m.ActionCharacterRepository()

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

// CreateActionCharacterRec -
func (m *Model) CreateActionCharacterRec(rec *record.ActionCharacter) error {

	l := m.loggerWithFunctionContext("CreateActionCharacterRec")

	l.Debug("Creating dungeon action character record >%#v<", rec)

	r := m.ActionCharacterRepository()

	err := m.validateActionCharacterRec(rec)
	if err != nil {
		l.Warn("failed model validation >%v<", err)
		return err
	}

	return r.CreateOne(rec)
}

// UpdateActionCharacterRec -
func (m *Model) UpdateActionCharacterRec(rec *record.ActionCharacter) error {

	l := m.loggerWithFunctionContext("UpdateActionCharacterRec")

	l.Debug("Updating dungeon action character record >%#v<", rec)

	r := m.ActionCharacterRepository()

	err := m.validateActionCharacterRec(rec)
	if err != nil {
		l.Warn("failed model validation >%v<", err)
		return err
	}

	return r.UpdateOne(rec)
}

// DeleteActionCharacterRec -
func (m *Model) DeleteActionCharacterRec(recID string) error {

	l := m.loggerWithFunctionContext("DeleteActionCharacterRec")

	l.Debug("Deleting dungeon action character rec ID >%s<", recID)

	r := m.ActionCharacterRepository()

	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.validateDeleteActionCharacterRec(recID)
	if err != nil {
		l.Warn("failed model validation >%v<", err)
		return err
	}

	return r.DeleteOne(recID)
}

// RemoveActionCharacterRec -
func (m *Model) RemoveActionCharacterRec(recID string) error {

	l := m.loggerWithFunctionContext("RemoveActionCharacterRec")

	l.Debug("Removing dungeon action character rec ID >%s<", recID)

	r := m.ActionCharacterRepository()

	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.validateDeleteActionCharacterRec(recID)
	if err != nil {
		l.Warn("failed model validation >%v<", err)
		return err
	}

	return r.RemoveOne(recID)
}
