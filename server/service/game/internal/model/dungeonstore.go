package model

import (
	"database/sql"
	"fmt"

	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// GetDungeonRecs -
func (m *Model) GetDungeonRecs(params map[string]interface{}, operators map[string]string, forUpdate bool) ([]*record.Dungeon, error) {

	l := m.Logger("GetDungeonRecs")

	l.Debug("Getting dungeon records params >%s<", params)

	r := m.DungeonRepository()

	return r.GetMany(params, operators, forUpdate)
}

// GetDungeonRec -
func (m *Model) GetDungeonRec(recID string, forUpdate bool) (*record.Dungeon, error) {

	l := m.Logger("GetDungeonRec")

	l.Debug("Getting dungeon rec ID >%s<", recID)

	r := m.DungeonRepository()

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

// CreateDungeonRec -
func (m *Model) CreateDungeonRec(rec *record.Dungeon) error {

	l := m.Logger("CreateDungeonRec")

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

	l := m.Logger("UpdateDungeonRec")

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

	l := m.Logger("DeleteDungeonRec")

	l.Debug("Deleting dungeon rec ID >%s<", recID)

	r := m.DungeonRepository()

	// validate UUID
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

	l := m.Logger("RemoveDungeonRec")

	l.Debug("Removing dungeon rec ID >%s<", recID)

	r := m.DungeonRepository()

	// validate UUID
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
