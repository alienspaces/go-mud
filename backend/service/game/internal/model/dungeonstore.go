package model

import (
	"database/sql"
	"fmt"

	coresql "gitlab.com/alienspaces/go-mud/backend/core/sql"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/record"
)

// GetDungeonRecs -
func (m *Model) GetDungeonRecs(opts *coresql.Options) ([]*record.Dungeon, error) {

	l := m.loggerWithFunctionContext("GetDungeonRecs")

	l.Debug("Getting dungeon records opts >%#v<", opts)

	r := m.DungeonRepository()

	return r.GetMany(opts)
}

// GetDungeonRec -
func (m *Model) GetDungeonRec(recID string, lock *coresql.Lock) (*record.Dungeon, error) {

	l := m.loggerWithFunctionContext("GetDungeonRec")

	l.Debug("Getting dungeon rec ID >%s<", recID)

	r := m.DungeonRepository()

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

// CreateDungeonRec -
func (m *Model) CreateDungeonRec(rec *record.Dungeon) error {

	l := m.loggerWithFunctionContext("CreateDungeonRec")

	l.Debug("Creating dungeon record >%#v<", rec)

	r := m.DungeonRepository()

	err := m.validateDungeonRec(rec)
	if err != nil {
		l.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.CreateOne(rec)
}

// UpdateDungeonRec -
func (m *Model) UpdateDungeonRec(rec *record.Dungeon) error {

	l := m.loggerWithFunctionContext("UpdateDungeonRec")

	l.Debug("Updating dungeon record >%#v<", rec)

	r := m.DungeonRepository()

	err := m.validateDungeonRec(rec)
	if err != nil {
		l.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.UpdateOne(rec)
}

// DeleteDungeonRec -
func (m *Model) DeleteDungeonRec(recID string) error {

	l := m.loggerWithFunctionContext("DeleteDungeonRec")

	l.Debug("Deleting dungeon rec ID >%s<", recID)

	r := m.DungeonRepository()

	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.validateDeleteDungeonRec(recID)
	if err != nil {
		l.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.DeleteOne(recID)
}

// RemoveDungeonRec -
func (m *Model) RemoveDungeonRec(recID string) error {

	l := m.loggerWithFunctionContext("RemoveDungeonRec")

	l.Debug("Removing dungeon rec ID >%s<", recID)

	r := m.DungeonRepository()

	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.validateDeleteDungeonRec(recID)
	if err != nil {
		l.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.RemoveOne(recID)
}
