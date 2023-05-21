package model

import (
	"database/sql"
	"fmt"

	coresql "gitlab.com/alienspaces/go-mud/backend/core/sql"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/record"
)

// GetDungeonInstanceViewRecs -
func (m *Model) GetDungeonInstanceViewRecs(opts *coresql.Options) ([]*record.DungeonInstanceView, error) {

	l := m.Logger("GetDungeonInstanceViewRecs")

	l.Debug("Getting dungeon records opts >%#v<", opts)

	r := m.DungeonInstanceViewRepository()

	return r.GetMany(opts)
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

	rec, err := r.GetOne(recID, nil)
	if err == sql.ErrNoRows {
		l.Warn("No record found ID >%s<", recID)
		return nil, nil
	}

	return rec, err
}

// GetDungeonInstanceRecs -
func (m *Model) GetDungeonInstanceRecs(opts *coresql.Options) ([]*record.DungeonInstance, error) {

	l := m.Logger("GetDungeonInstanceRecs")

	l.Debug("Getting dungeon records opts >%#v<", opts)

	r := m.DungeonInstanceRepository()

	return r.GetMany(opts)
}

// GetDungeonInstanceRec -
func (m *Model) GetDungeonInstanceRec(recID string, lock *coresql.Lock) (*record.DungeonInstance, error) {

	l := m.Logger("GetDungeonInstanceRec")

	l.Debug("Getting dungeon rec ID >%s<", recID)

	r := m.DungeonInstanceRepository()

	// validate UUID
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
