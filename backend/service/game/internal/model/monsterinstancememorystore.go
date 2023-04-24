package model

import (
	"database/sql"
	"fmt"

	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/record"
)

// GetMonsterInstanceMemoryRecs -
func (m *Model) GetMonsterInstanceMemoryRecs(params map[string]interface{}, operators map[string]string, forUpdate bool) ([]*record.MonsterInstanceMemory, error) {

	l := m.Logger("GetMonsterInstanceMemoryRecs")

	l.Debug("Getting monster instance records params >%s<", params)

	r := m.MonsterInstanceMemoryRepository()

	return r.GetMany(params, operators, forUpdate)
}

// GetMonsterInstanceMemoryRec -
func (m *Model) GetMonsterInstanceMemoryRec(recID string, forUpdate bool) (*record.MonsterInstanceMemory, error) {

	l := m.Logger("GetMonsterInstanceMemoryRec")

	l.Debug("Getting monster instance record ID >%s<", recID)

	r := m.MonsterInstanceMemoryRepository()

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

// CreateMonsterInstanceMemoryRec -
func (m *Model) CreateMonsterInstanceMemoryRec(rec *record.MonsterInstanceMemory) error {

	l := m.Logger("CreateMonsterInstanceMemoryRec")

	l.Debug("Creating monster record >%#v<", rec)

	r := m.MonsterInstanceMemoryRepository()

	err := m.validateMonsterInstanceMemoryRec(rec)
	if err != nil {
		l.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.CreateOne(rec)
}

// UpdateMonsterInstanceMemoryRec -
func (m *Model) UpdateMonsterInstanceMemoryRec(rec *record.MonsterInstanceMemory) error {

	l := m.Logger("UpdateMonsterInstanceMemoryRec")

	l.Debug("Updating monster record >%#v<", rec)

	r := m.MonsterInstanceMemoryRepository()

	err := m.validateMonsterInstanceMemoryRec(rec)
	if err != nil {
		l.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.UpdateOne(rec)
}

// DeleteMonsterInstanceMemoryRec -
func (m *Model) DeleteMonsterInstanceMemoryRec(recID string) error {

	l := m.Logger("DeleteMonsterInstanceMemoryRec")

	l.Debug("Deleting monster rec ID >%s<", recID)

	r := m.MonsterInstanceMemoryRepository()

	// validate UUID
	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.validateDeleteMonsterInstanceMemoryRec(recID)
	if err != nil {
		l.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.DeleteOne(recID)
}

// RemoveMonsterInstanceMemoryRec -
func (m *Model) RemoveMonsterInstanceMemoryRec(recID string) error {

	l := m.Logger("RemoveMonsterInstanceMemoryRec")

	l.Debug("Removing monster rec ID >%s<", recID)

	r := m.MonsterInstanceMemoryRepository()

	// validate UUID
	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.validateDeleteMonsterInstanceMemoryRec(recID)
	if err != nil {
		l.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.RemoveOne(recID)
}
