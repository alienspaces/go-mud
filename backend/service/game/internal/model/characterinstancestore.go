package model

import (
	"database/sql"
	"fmt"

	coreerror "gitlab.com/alienspaces/go-mud/backend/core/error"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/record"
)

// GetCharacterInstanceRecByCharacterID -
func (m *Model) GetCharacterInstanceRecByCharacterID(characterID string) (*record.CharacterInstance, error) {
	l := m.Logger("GetCharacterInstanceRecByCharacterID")

	characterInstanceRecs, err := m.GetCharacterInstanceRecs(
		map[string]interface{}{
			"character_id": characterID,
		}, nil, false,
	)
	if err != nil {
		l.Warn("failed getting character ID >%s< instance records >%v<", characterID, err)
		return nil, err
	}

	if len(characterInstanceRecs) == 0 {
		l.Warn("character with ID >%s< has no character instance record", characterID)
		return nil, nil
	}

	if len(characterInstanceRecs) > 1 {
		l.Warn("unexpected number of character instance records returned >%d<", len(characterInstanceRecs))
		err := coreerror.NewInternalError()
		return nil, err
	}

	return characterInstanceRecs[0], nil
}

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

// GetCharacterInstanceViewRecByCharacterID -
func (m *Model) GetCharacterInstanceViewRecByCharacterID(characterID string) (*record.CharacterInstanceView, error) {
	l := m.Logger("GetCharacterInstanceViewRecByCharacterID")

	characterInstanceViewRecs, err := m.GetCharacterInstanceViewRecs(
		map[string]interface{}{
			"character_id": characterID,
		}, nil,
	)
	if err != nil {
		l.Warn("failed getting character ID >%s< character instance view records >%v<", characterID, err)
		return nil, err
	}

	if len(characterInstanceViewRecs) == 0 {
		l.Warn("character with ID >%s< has no character instance view record", characterID)
		return nil, nil
	}

	if len(characterInstanceViewRecs) > 1 {
		l.Warn("unexpected number of character instance view records returned >%d<", len(characterInstanceViewRecs))
		err := coreerror.NewInternalError()
		return nil, err
	}

	return characterInstanceViewRecs[0], nil
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
