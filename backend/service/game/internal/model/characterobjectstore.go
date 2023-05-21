package model

import (
	"database/sql"
	"fmt"

	coresql "gitlab.com/alienspaces/go-mud/backend/core/sql"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/record"
)

// GetCharacterObjectRecs -
func (m *Model) GetCharacterObjectRecs(opts *coresql.Options) ([]*record.CharacterObject, error) {

	l := m.Logger("GetCharacterObjectRecs")

	l.Debug("Getting dungeon character records opts >%#v<", opts)

	r := m.CharacterObjectRepository()

	return r.GetMany(opts)
}

// GetCharacterObjectRec -
func (m *Model) GetCharacterObjectRec(recID string, lock *coresql.Lock) (*record.CharacterObject, error) {

	l := m.Logger("GetCharacterObjectRec")

	l.Debug("Getting dungeon character rec ID >%s<", recID)

	r := m.CharacterObjectRepository()

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

// CreateCharacterObjectRec -
func (m *Model) CreateCharacterObjectRec(rec *record.CharacterObject) error {

	l := m.Logger("CreateCharacterObjectRec")

	l.Debug("Creating dungeon character record >%#v<", rec)

	r := m.CharacterObjectRepository()

	err := m.validateCharacterObjectRec(rec)
	if err != nil {
		l.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.CreateOne(rec)
}

// UpdateCharacterObjectRec -
func (m *Model) UpdateCharacterObjectRec(rec *record.CharacterObject) error {

	l := m.Logger("UpdateCharacterObjectRec")

	l.Debug("Updating dungeon character record >%#v<", rec)

	r := m.CharacterObjectRepository()

	err := m.validateCharacterObjectRec(rec)
	if err != nil {
		l.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.UpdateOne(rec)
}

// DeleteCharacterObjectRec -
func (m *Model) DeleteCharacterObjectRec(recID string) error {

	l := m.Logger("DeleteCharacterObjectRec")

	l.Debug("Deleting dungeon character rec ID >%s<", recID)

	r := m.CharacterObjectRepository()

	// validate UUID
	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.validateDeleteCharacterObjectRec(recID)
	if err != nil {
		l.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.DeleteOne(recID)
}

// RemoveCharacterObjectRec -
func (m *Model) RemoveCharacterObjectRec(recID string) error {

	l := m.Logger("RemoveCharacterObjectRec")

	l.Debug("Removing dungeon character rec ID >%s<", recID)

	r := m.CharacterObjectRepository()

	// validate UUID
	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.validateDeleteCharacterObjectRec(recID)
	if err != nil {
		l.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.RemoveOne(recID)
}
