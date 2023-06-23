package model

import (
	"fmt"

	coreerror "gitlab.com/alienspaces/go-mud/backend/core/error"
	coresql "gitlab.com/alienspaces/go-mud/backend/core/sql"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/record"
)

// validateCreateCharacterInstanceRec - validates creating a character instance record
func (m *Model) validateCreateCharacterInstanceRec(rec *record.CharacterInstance) error {
	l := m.Logger("validateCreateCharacterInstanceRec")

	characterInstanceRecs, err := m.GetCharacterInstanceRecs(
		&coresql.Options{
			Params: []coresql.Param{
				{
					Col: "character_id",
					Val: rec.CharacterID,
				},
			},
		},
	)
	if err != nil {
		l.Warn("failed getting character instance record >%v<", err)
		return err
	}

	if len(characterInstanceRecs) > 0 {
		msg := fmt.Sprintf("character with ID >%s< is already inside a dungeon", rec.CharacterID)
		l.Warn(msg)
		err := coreerror.CreateInvalidError("character_id", msg)
		return err
	}

	err = m.validateCharacterInstanceAttributes(rec)
	if err != nil {
		l.Warn("failed validating character instance records >%v<", err)
		return err
	}

	return nil
}

// validateUpdateCharacterInstanceRec - validates updating a character instance record
func (m *Model) validateUpdateCharacterInstanceRec(rec *record.CharacterInstance) error {
	l := m.Logger("validateUpdateCharacterInstanceRec")

	err := m.validateCharacterInstanceAttributes(rec)
	if err != nil {
		l.Warn("failed validating character instance records >%v<", err)
		return err
	}

	return nil
}

// validateDeleteCharacterInstanceRec - validates deleting a character instance record
func (m *Model) validateDeleteCharacterInstanceRec(recID string) error {
	l := m.Logger("validateDeleteCharacterInstanceRec")

	rec, err := m.GetCharacterInstanceRec(recID, nil)
	if err != nil {
		l.Warn("failed getting character instance record >%v<", err)
		return err
	}

	if rec == nil {
		msg := fmt.Sprintf("failed validation, character instance ID >%s< does not exist", recID)
		l.Warn(msg)
		err := coreerror.CreateInvalidError("id", msg)
		return err
	}

	return nil
}

func (m *Model) validateCharacterInstanceAttributes(rec *record.CharacterInstance) error {

	// New character
	if rec.ID == "" {
		if rec.Strength+rec.Intelligence+rec.Dexterity > defaultAttributePoints {
			return fmt.Errorf("new character attributes exceeds allowed maximum of %d", defaultAttributePoints)
		}
		if rec.Strength == 0 {
			return fmt.Errorf("failed validation, Strength is empty")
		}
		if rec.Dexterity == 0 {
			return fmt.Errorf("failed validation, Dexterity is empty")
		}
		if rec.Intelligence == 0 {
			return fmt.Errorf("failed validation, Intelligence is empty")
		}
		if rec.Health == 0 {
			return fmt.Errorf("failed validation, Health is empty")
		}
		if rec.Fatigue == 0 {
			return fmt.Errorf("failed validation, Fatigue is empty")
		}
	}

	if rec.DungeonInstanceID == "" {
		return fmt.Errorf("failed validation, DungeonInstanceID is empty")
	}
	if rec.LocationInstanceID == "" {
		return fmt.Errorf("failed validation, LocationInstanceID is empty")
	}

	return nil
}
