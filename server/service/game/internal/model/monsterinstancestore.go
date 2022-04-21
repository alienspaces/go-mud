package model

import (
	"database/sql"
	"fmt"

	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// GetMonsterInstanceRecs -
func (m *Model) GetMonsterInstanceRecs(params map[string]interface{}, operators map[string]string, forUpdate bool) ([]*record.MonsterInstance, error) {

	l := m.Logger("GetMonsterInstanceRecs")

	l.Debug("Getting monster instance records params >%s<", params)

	r := m.MonsterInstanceRepository()

	return r.GetMany(params, operators, forUpdate)
}

// GetMonsterInstanceRec -
func (m *Model) GetMonsterInstanceRec(recID string, forUpdate bool) (*record.MonsterInstance, error) {

	l := m.Logger("GetMonsterInstanceRec")

	l.Debug("Getting monster instance record ID >%s<", recID)

	r := m.MonsterInstanceRepository()

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

// GetMonsterInstanceViewRecs -
func (m *Model) GetMonsterInstanceViewRecs(params map[string]interface{}, operators map[string]string) ([]*record.MonsterInstanceView, error) {

	l := m.Logger("GetMonsterInstanceViewRecs")

	l.Debug("Getting monster instance view records params >%s<", params)

	r := m.MonsterInstanceViewRepository()

	return r.GetMany(params, operators, false)
}

// GetMonsterInstanceViewRec -
func (m *Model) GetMonsterInstanceViewRec(recID string) (*record.MonsterInstanceView, error) {

	l := m.Logger("GetMonsterInstanceViewRec")

	l.Debug("Getting monster instance view record ID >%s<", recID)

	r := m.MonsterInstanceViewRepository()

	// validate UUID
	if !m.IsUUID(recID) {
		return nil, fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	rec, err := r.GetOne(recID, false)
	if err == sql.ErrNoRows {
		l.Warn("No record found ID >%s<", recID)
		return nil, nil
	}

	return rec, err
}

// CreateMonsterInstanceRec -
func (m *Model) CreateMonsterInstanceRec(rec *record.MonsterInstance) error {

	l := m.Logger("CreateMonsterInstanceRec")

	l.Debug("Creating monster rec >%#v<", rec)

	r := m.MonsterInstanceRepository()

	err := m.ValidateMonsterInstanceRec(rec)
	if err != nil {
		l.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.CreateOne(rec)
}

// UpdateMonsterInstanceRec -
func (m *Model) UpdateMonsterInstanceRec(rec *record.MonsterInstance) error {

	l := m.Logger("UpdateMonsterInstanceRec")

	l.Debug("Updating monster rec >%#v<", rec)

	r := m.MonsterInstanceRepository()

	err := m.ValidateMonsterInstanceRec(rec)
	if err != nil {
		l.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.UpdateOne(rec)
}

// DeleteMonsterInstanceRec -
func (m *Model) DeleteMonsterInstanceRec(recID string) error {

	l := m.Logger("DeleteMonsterInstanceRec")

	l.Debug("Deleting monster rec ID >%s<", recID)

	r := m.MonsterInstanceRepository()

	// validate UUID
	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.ValidateDeleteMonsterInstanceRec(recID)
	if err != nil {
		l.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.DeleteOne(recID)
}

// RemoveMonsterInstanceRec -
func (m *Model) RemoveMonsterInstanceRec(recID string) error {

	l := m.Logger("RemoveMonsterInstanceRec")

	l.Debug("Removing monster rec ID >%s<", recID)

	r := m.MonsterInstanceRepository()

	// validate UUID
	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.ValidateDeleteMonsterInstanceRec(recID)
	if err != nil {
		l.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.RemoveOne(recID)
}
