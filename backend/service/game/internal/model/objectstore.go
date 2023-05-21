package model

import (
	"database/sql"
	"fmt"

	coresql "gitlab.com/alienspaces/go-mud/backend/core/sql"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/record"
)

// GetObjectRecs -
func (m *Model) GetObjectRecs(opts *coresql.Options) ([]*record.Object, error) {

	l := m.Logger("GetObjectRecs")

	l.Debug("Getting dungeon object records opts >%#v<", opts)

	r := m.ObjectRepository()

	return r.GetMany(opts)
}

// GetObjectRec -
func (m *Model) GetObjectRec(recID string, lock *coresql.Lock) (*record.Object, error) {

	l := m.Logger("GetObjectRec")

	r := m.ObjectRepository()

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

// CreateObjectRec -
func (m *Model) CreateObjectRec(rec *record.Object) error {

	l := m.Logger("CreateObjectRec")

	l.Debug("Creating dungeon object record >%#v<", rec)

	r := m.ObjectRepository()

	err := m.validateObjectRec(rec)
	if err != nil {
		l.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.CreateOne(rec)
}

// UpdateObjectRec -
func (m *Model) UpdateObjectRec(rec *record.Object) error {

	l := m.Logger("UpdateObjectRec")

	l.Debug("Updating dungeon object record >%#v<", rec)

	r := m.ObjectRepository()

	err := m.validateObjectRec(rec)
	if err != nil {
		l.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.UpdateOne(rec)
}

// DeleteObjectRec -
func (m *Model) DeleteObjectRec(recID string) error {

	l := m.Logger("DeleteObjectRec")

	l.Debug("Deleting dungeon object rec ID >%s<", recID)

	r := m.ObjectRepository()

	// validate UUID
	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.validateDeleteObjectRec(recID)
	if err != nil {
		l.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.DeleteOne(recID)
}

// RemoveObjectRec -
func (m *Model) RemoveObjectRec(recID string) error {

	l := m.Logger("RemoveObjectRec")

	l.Debug("Removing dungeon object rec ID >%s<", recID)

	r := m.ObjectRepository()

	// validate UUID
	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.validateDeleteObjectRec(recID)
	if err != nil {
		l.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.RemoveOne(recID)
}
