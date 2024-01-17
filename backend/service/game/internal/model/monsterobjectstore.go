package model

import (
	"database/sql"
	"fmt"

	coresql "gitlab.com/alienspaces/go-mud/backend/core/sql"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/record"
)

// GetMonsterObjectRecs -
func (m *Model) GetMonsterObjectRecs(opts *coresql.Options) ([]*record.MonsterObject, error) {

	l := m.loggerWithFunctionContext("GetMonsterObjectRecs")

	l.Debug("Getting dungeon monster records opts >%#v<", opts)

	r := m.MonsterObjectRepository()

	return r.GetMany(opts)
}

// GetMonsterObjectRec -
func (m *Model) GetMonsterObjectRec(recID string, lock *coresql.Lock) (*record.MonsterObject, error) {

	l := m.loggerWithFunctionContext("GetMonsterObjectRec")

	l.Debug("Getting dungeon monster rec ID >%s<", recID)

	r := m.MonsterObjectRepository()

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

// CreateMonsterObjectRec -
func (m *Model) CreateMonsterObjectRec(rec *record.MonsterObject) error {

	l := m.loggerWithFunctionContext("CreateMonsterObjectRec")

	l.Debug("Creating dungeon monster record >%#v<", rec)

	r := m.MonsterObjectRepository()

	err := m.validateMonsterObjectRec(rec)
	if err != nil {
		l.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.CreateOne(rec)
}

// UpdateMonsterObjectRec -
func (m *Model) UpdateMonsterObjectRec(rec *record.MonsterObject) error {

	l := m.loggerWithFunctionContext("UpdateMonsterObjectRec")

	l.Debug("Updating dungeon monster record >%#v<", rec)

	r := m.MonsterObjectRepository()

	err := m.validateMonsterObjectRec(rec)
	if err != nil {
		l.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.UpdateOne(rec)
}

// DeleteMonsterObjectRec -
func (m *Model) DeleteMonsterObjectRec(recID string) error {

	l := m.loggerWithFunctionContext("DeleteMonsterObjectRec")

	l.Debug("Deleting dungeon monster rec ID >%s<", recID)

	r := m.MonsterObjectRepository()

	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.validateDeleteMonsterObjectRec(recID)
	if err != nil {
		l.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.DeleteOne(recID)
}

// RemoveMonsterObjectRec -
func (m *Model) RemoveMonsterObjectRec(recID string) error {

	l := m.loggerWithFunctionContext("RemoveMonsterObjectRec")

	l.Debug("Removing dungeon monster rec ID >%s<", recID)

	r := m.MonsterObjectRepository()

	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.validateDeleteMonsterObjectRec(recID)
	if err != nil {
		l.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.RemoveOne(recID)
}
