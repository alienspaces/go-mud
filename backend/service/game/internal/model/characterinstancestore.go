package model

import (
	"database/sql"
	"fmt"

	coreerror "gitlab.com/alienspaces/go-mud/backend/core/error"
	coresql "gitlab.com/alienspaces/go-mud/backend/core/sql"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/record"
)

// GetCharacterInstanceRecByCharacterID -
func (m *Model) GetCharacterInstanceRecByCharacterID(characterID string) (*record.CharacterInstance, error) {
	l := m.Logger("GetCharacterInstanceRecByCharacterID")

	characterInstanceRecs, err := m.GetCharacterInstanceRecs(
		&coresql.Options{
			Params: []coresql.Param{
				{
					Col: "character_id",
					Val: characterID,
				},
			},
		},
	)
	if err != nil {
		l.Warn("failed getting character ID >%s< instance records >%v<", characterID, err)
		return nil, err
	}

	if len(characterInstanceRecs) == 0 {
		l.Info("character with ID >%s< has no character instance record", characterID)
		return nil, nil
	}

	if len(characterInstanceRecs) > 1 {
		err := coreerror.NewInternalError("unexpected number of character instance records returned >%d<", len(characterInstanceRecs))
		l.Warn(err.Error())
		return nil, err
	}

	return characterInstanceRecs[0], nil
}

// GetCharacterInstanceRecs -
func (m *Model) GetCharacterInstanceRecs(opts *coresql.Options) ([]*record.CharacterInstance, error) {

	l := m.Logger("GetCharacterInstanceRecs")

	l.Debug("Getting character instance records opts >%#v<", opts)

	r := m.CharacterInstanceRepository()

	return r.GetMany(opts)
}

// GetCharacterInstanceRec -
func (m *Model) GetCharacterInstanceRec(recID string, lock *coresql.Lock) (*record.CharacterInstance, error) {
	l := m.Logger("GetCharacterInstanceRec")

	l.Debug("Getting character instance record ID >%s<", recID)

	r := m.CharacterInstanceRepository()

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

// GetCharacterInstanceViewRecByCharacterID -
func (m *Model) GetCharacterInstanceViewRecByCharacterID(characterID string) (*record.CharacterInstanceView, error) {
	l := m.Logger("GetCharacterInstanceViewRecByCharacterID")

	characterInstanceViewRecs, err := m.GetCharacterInstanceViewRecs(
		&coresql.Options{
			Params: []coresql.Param{
				{
					Col: "character_id",
					Val: characterID,
				},
			},
		},
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
		err := coreerror.NewInternalError("unexpected number of character instance view records returned >%d<", len(characterInstanceViewRecs))
		l.Warn(err.Error())
		return nil, err
	}

	return characterInstanceViewRecs[0], nil
}

// GetCharacterInstanceRecs -
func (m *Model) GetCharacterInstanceViewRecs(opts *coresql.Options) ([]*record.CharacterInstanceView, error) {

	l := m.Logger("GetCharacterInstanceViewRecs")

	l.Debug("Getting character instance view records opts >%#v<", opts)

	r := m.CharacterInstanceViewRepository()

	return r.GetMany(opts)
}

// GetCharacterInstanceRec -
func (m *Model) GetCharacterInstanceViewRec(recID string) (*record.CharacterInstanceView, error) {

	l := m.Logger("GetCharacterInstanceViewRec")

	l.Debug("Getting character instance view record ID >%s<", recID)

	r := m.CharacterInstanceViewRepository()

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
