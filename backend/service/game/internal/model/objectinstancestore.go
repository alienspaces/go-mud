package model

import (
	"database/sql"
	"fmt"

	coresql "gitlab.com/alienspaces/go-mud/backend/core/sql"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/record"
)

// GetObjectInstanceRecs -
func (m *Model) GetObjectInstanceRecs(opts *coresql.Options) ([]*record.ObjectInstance, error) {

	l := m.Logger("GetObjectInstanceRecs")

	l.Debug("Getting object instance records opts >%#v<", opts)

	r := m.ObjectInstanceRepository()

	return r.GetMany(opts)
}

// GetObjectInstanceRec -
func (m *Model) GetObjectInstanceRec(recID string, lock *coresql.Lock) (*record.ObjectInstance, error) {

	l := m.Logger("GetObjectInstanceRec")

	l.Debug("Getting object instance record ID >%s<", recID)

	r := m.ObjectInstanceRepository()

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

// GetObjectInstanceViewRecs -
func (m *Model) GetObjectInstanceViewRecs(opts *coresql.Options) ([]*record.ObjectInstanceView, error) {

	l := m.Logger("GetObjectInstanceViewRecs")

	l.Debug("Getting object instance view records opts >%#v<", opts)

	r := m.ObjectInstanceViewRepository()

	return r.GetMany(opts)
}

// GetObjectInstanceViewRec -
func (m *Model) GetObjectInstanceViewRec(recID string) (*record.ObjectInstanceView, error) {

	l := m.Logger("GetObjectInstanceViewRec")

	l.Debug("Getting object instance view record ID >%s<", recID)

	r := m.ObjectInstanceViewRepository()

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

// CreateObjectInstanceRec -
func (m *Model) CreateObjectInstanceRec(rec *record.ObjectInstance) error {

	l := m.Logger("CreateObjectInstanceRec")

	l.Debug("Creating object instance record >%#v<", rec)

	r := m.ObjectInstanceRepository()

	err := m.validateObjectInstanceRec(rec)
	if err != nil {
		l.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.CreateOne(rec)
}

// UpdateObjectInstanceRec -
func (m *Model) UpdateObjectInstanceRec(rec *record.ObjectInstance) error {

	l := m.Logger("UpdateObjectInstanceRec")

	l.Warn("Updating object instance record >%#v<", rec)

	r := m.ObjectInstanceRepository()

	err := m.validateObjectInstanceRec(rec)
	if err != nil {
		l.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.UpdateOne(rec)
}

// DeleteObjectInstanceRec -
func (m *Model) DeleteObjectInstanceRec(recID string) error {

	l := m.Logger("DeleteObjectInstanceRec")

	l.Debug("Deleting object instance record ID >%s<", recID)

	r := m.ObjectInstanceRepository()

	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.validateDeleteObjectInstanceRec(recID)
	if err != nil {
		l.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.DeleteOne(recID)
}

// RemoveObjectInstanceRec -
func (m *Model) RemoveObjectInstanceRec(recID string) error {

	l := m.Logger("RemoveObjectInstanceRec")

	l.Debug("Removing object instance record ID >%s<", recID)

	r := m.ObjectInstanceRepository()

	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.validateDeleteObjectInstanceRec(recID)
	if err != nil {
		l.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.RemoveOne(recID)
}
