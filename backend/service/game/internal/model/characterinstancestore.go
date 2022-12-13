package model

import (
	"database/sql"
	"fmt"

	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/record"
)

// GetCharacterInstanceRecs -
func (m *Model) GetCharacterInstanceRecs(params map[string]interface{}, operators map[string]string, forUpdate bool) ([]*record.CharacterInstance, error) {

	l := m.Logger("GetCharacterInstanceRecs")

	l.Debug("Getting character instance records params >%s<", params)

	r := m.CharacterInstanceRepository()

	return r.GetMany(params, operators, forUpdate)
}

// GetCharacterInstanceRec -
func (m *Model) GetCharacterInstanceRec(recID string, forUpdate bool) (*record.CharacterInstance, error) {

	l := m.Logger("GetCharacterInstanceRec")

	l.Debug("Getting character instance record ID >%s<", recID)

	r := m.CharacterInstanceRepository()

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

// GetCharacterInstanceRecs -
func (m *Model) GetCharacterInstanceViewRecs(params map[string]interface{}, operators map[string]string) ([]*record.CharacterInstanceView, error) {

	l := m.Logger("GetCharacterInstanceViewRecs")

	l.Debug("Getting character instance view records params >%s<", params)

	r := m.CharacterInstanceViewRepository()

	return r.GetMany(params, operators, false)
}

// GetCharacterInstanceRec -
func (m *Model) GetCharacterInstanceViewRec(recID string) (*record.CharacterInstanceView, error) {

	l := m.Logger("GetCharacterInstanceViewRec")

	l.Debug("Getting character instance view record ID >%s<", recID)

	r := m.CharacterInstanceViewRepository()

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

// CreateCharacterInstanceRec -
func (m *Model) CreateCharacterInstanceRec(rec *record.CharacterInstance) error {

	l := m.Logger("CreateCharacterInstanceRec")

	l.Debug("Creating character record >%#v<", rec)

	characterRepo := m.CharacterInstanceRepository()

	err := m.validateCreateCharacterInstanceRec(rec)
	if err != nil {
		l.Debug("Failed model validation >%v<", err)
		return err
	}

	return characterRepo.CreateOne(rec)
}

// UpdateCharacterInstanceRec -
func (m *Model) UpdateCharacterInstanceRec(rec *record.CharacterInstance) error {

	l := m.Logger("UpdateCharacterInstanceRec")

	l.Debug("Updating character record >%#v<", rec)

	r := m.CharacterInstanceRepository()

	err := m.validateUpdateCharacterInstanceRec(rec)
	if err != nil {
		l.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.UpdateOne(rec)
}

// DeleteCharacterInstanceRec -
func (m *Model) DeleteCharacterInstanceRec(recID string) error {

	l := m.Logger("DeleteCharacterInstanceRec")

	l.Debug("Deleting character rec ID >%s<", recID)

	r := m.CharacterInstanceRepository()

	// validate UUID
	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.validateDeleteCharacterInstanceRec(recID)
	if err != nil {
		l.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.DeleteOne(recID)
}

// RemoveCharacterInstanceRec -
func (m *Model) RemoveCharacterInstanceRec(recID string) error {

	l := m.Logger("RemoveCharacterInstanceRec")

	l.Debug("Removing character rec ID >%s<", recID)

	r := m.CharacterInstanceRepository()

	// validate UUID
	if !m.IsUUID(recID) {
		return fmt.Errorf("ID >%s< is not a valid UUID", recID)
	}

	err := m.validateDeleteCharacterInstanceRec(recID)
	if err != nil {
		l.Debug("Failed model validation >%v<", err)
		return err
	}

	return r.RemoveOne(recID)
}
