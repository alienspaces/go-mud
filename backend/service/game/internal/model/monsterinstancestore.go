package model

import (
	"database/sql"
	"fmt"

	coresql "gitlab.com/alienspaces/go-mud/backend/core/sql"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/record"
)

// GetMonsterInstanceRecs -
func (m *Model) GetMonsterInstanceRecs(opts *coresql.Options) ([]*record.MonsterInstance, error) {

	l := m.loggerWithContext("GetMonsterInstanceRecs")

	l.Debug("Getting monster instance records opts >%#v<", opts)

	r := m.MonsterInstanceRepository()

	return r.GetMany(opts)
}

// GetMonsterInstanceRec -
func (m *Model) GetMonsterInstanceRec(recID string, lock *coresql.Lock) (*record.MonsterInstance, error) {

	l := m.loggerWithContext("GetMonsterInstanceRec")

	l.Debug("Getting monster instance record ID >%s<", recID)

	r := m.MonsterInstanceRepository()

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

// GetMonsterInstanceViewRecs -
func (m *Model) GetMonsterInstanceViewRecs(opts *coresql.Options) ([]*record.MonsterInstanceView, error) {

	l := m.loggerWithContext("GetMonsterInstanceViewRecs")

	l.Debug("Getting monster instance view records opts >%#v<", opts)

	r := m.MonsterInstanceViewRepository()

	return r.GetMany(opts)
}

// GetMonsterInstanceViewRec -
func (m *Model) GetMonsterInstanceViewRec(recID string) (*record.MonsterInstanceView, error) {

	l := m.loggerWithContext("GetMonsterInstanceViewRec")

	l.Debug("Getting monster instance view record ID >%s<", recID)

	r := m.MonsterInstanceViewRepository()

	if !m.IsUUID(recID) {
		return nil, fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	rec, err := r.GetOne(recID, nil)
	if err == sql.ErrNoRows {
		l.Warn("No record found ID >%s<", recID)
		return nil, nil
	}

	return rec, err
}

// CreateMonsterInstanceRec -
func (m *Model) CreateMonsterInstanceRec(rec *record.MonsterInstance) error {

	l := m.loggerWithContext("CreateMonsterInstanceRec")

	l.Debug("Creating monster record >%#v<", rec)

	r := m.MonsterInstanceRepository()

	err := m.validateMonsterInstanceRec(rec)
	if err != nil {
		l.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.CreateOne(rec)
}

// UpdateMonsterInstanceRec -
func (m *Model) UpdateMonsterInstanceRec(rec *record.MonsterInstance) error {

	l := m.loggerWithContext("UpdateMonsterInstanceRec")

	l.Debug("Updating monster record >%#v<", rec)

	r := m.MonsterInstanceRepository()

	err := m.validateMonsterInstanceRec(rec)
	if err != nil {
		l.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.UpdateOne(rec)
}

// DeleteMonsterInstanceRec -
func (m *Model) DeleteMonsterInstanceRec(recID string) error {

	l := m.loggerWithContext("DeleteMonsterInstanceRec")

	l.Debug("Deleting monster rec ID >%s<", recID)

	r := m.MonsterInstanceRepository()

	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.validateDeleteMonsterInstanceRec(recID)
	if err != nil {
		l.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.DeleteOne(recID)
}

// RemoveMonsterInstanceRec -
func (m *Model) RemoveMonsterInstanceRec(recID string) error {

	l := m.loggerWithContext("RemoveMonsterInstanceRec")

	l.Debug("Removing monster rec ID >%s<", recID)

	r := m.MonsterInstanceRepository()

	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.validateDeleteMonsterInstanceRec(recID)
	if err != nil {
		l.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.RemoveOne(recID)
}
