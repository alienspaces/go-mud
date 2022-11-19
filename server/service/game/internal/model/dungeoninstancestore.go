package model

import (
	"database/sql"
	"fmt"

	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// GetDungeonInstanceViewRecs -
func (m *Model) GetDungeonInstanceViewRecs(params map[string]interface{}, operators map[string]string, forUpdate bool) ([]*record.DungeonInstanceView, error) {

	l := m.Logger("GetDungeonInstanceViewRecs")

	l.Debug("Getting dungeon records params >%s<", params)

	r := m.DungeonInstanceViewRepository()

	return r.GetMany(params, operators, forUpdate)
}

// GetDungeonInstanceViewRec -
func (m *Model) GetDungeonInstanceViewRec(recID string) (*record.DungeonInstanceView, error) {

	l := m.Logger("GetDungeonInstanceViewRec")

	l.Debug("Getting dungeon rec ID >%s<", recID)

	r := m.DungeonInstanceViewRepository()

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

// GetDungeonInstanceRecs -
func (m *Model) GetDungeonInstanceRecs(params map[string]interface{}, operators map[string]string, forUpdate bool) ([]*record.DungeonInstance, error) {

	l := m.Logger("GetDungeonInstanceRecs")

	l.Debug("Getting dungeon records params >%s<", params)

	r := m.DungeonInstanceRepository()

	return r.GetMany(params, operators, forUpdate)
}

// GetDungeonInstanceRec -
func (m *Model) GetDungeonInstanceRec(recID string, forUpdate bool) (*record.DungeonInstance, error) {

	l := m.Logger("GetDungeonInstanceRec")

	l.Debug("Getting dungeon rec ID >%s<", recID)

	r := m.DungeonInstanceRepository()

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

// CreateDungeonInstanceRec -
func (m *Model) CreateDungeonInstanceRec(rec *record.DungeonInstance) error {

	l := m.Logger("CreateDungeonInstanceRec")

	l.Debug("Creating dungeon record >%#v<", rec)

	r := m.DungeonInstanceRepository()

	err := m.validateDungeonInstanceRec(rec)
	if err != nil {
		l.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.CreateOne(rec)
}

// UpdateDungeonInstanceRec -
func (m *Model) UpdateDungeonInstanceRec(rec *record.DungeonInstance) error {

	l := m.Logger("UpdateDungeonInstanceRec")

	l.Debug("Updating dungeon record >%#v<", rec)

	r := m.DungeonInstanceRepository()

	err := m.validateDungeonInstanceRec(rec)
	if err != nil {
		l.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.UpdateOne(rec)
}

// DeleteDungeonInstanceRec -
func (m *Model) DeleteDungeonInstanceRec(recID string) error {

	l := m.Logger("DeleteDungeonInstanceRec")

	l.Debug("Deleting dungeon rec ID >%s<", recID)

	r := m.DungeonInstanceRepository()

	// validate UUID
	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.validateDeleteDungeonInstanceRec(recID)
	if err != nil {
		l.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.DeleteOne(recID)
}

// RemoveDungeonInstanceRec -
func (m *Model) RemoveDungeonInstanceRec(recID string) error {

	l := m.Logger("RemoveDungeonInstanceRec")

	l.Debug("Removing dungeon rec ID >%s<", recID)

	r := m.DungeonInstanceRepository()

	// validate UUID
	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.validateDeleteDungeonInstanceRec(recID)
	if err != nil {
		l.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.RemoveOne(recID)
}
