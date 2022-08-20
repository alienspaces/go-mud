package model

import (
	"database/sql"
	"fmt"

	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// GetMonsterObjectRecs -
func (m *Model) GetMonsterObjectRecs(params map[string]interface{}, operators map[string]string, forUpdate bool) ([]*record.MonsterObject, error) {

	l := m.Logger("GetMonsterObjectRecs")

	l.Debug("Getting dungeon monster records params >%s<", params)

	r := m.MonsterObjectRepository()

	return r.GetMany(params, operators, forUpdate)
}

// GetMonsterObjectRec -
func (m *Model) GetMonsterObjectRec(recID string, forUpdate bool) (*record.MonsterObject, error) {

	l := m.Logger("GetMonsterObjectRec")

	l.Debug("Getting dungeon monster rec ID >%s<", recID)

	r := m.MonsterObjectRepository()

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

// CreateMonsterObjectRec -
func (m *Model) CreateMonsterObjectRec(rec *record.MonsterObject) error {

	l := m.Logger("CreateMonsterObjectRec")

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

	l := m.Logger("UpdateMonsterObjectRec")

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

	l := m.Logger("DeleteMonsterObjectRec")

	l.Debug("Deleting dungeon monster rec ID >%s<", recID)

	r := m.MonsterObjectRepository()

	// validate UUID
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

	l := m.Logger("RemoveMonsterObjectRec")

	l.Debug("Removing dungeon monster rec ID >%s<", recID)

	r := m.MonsterObjectRepository()

	// validate UUID
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
