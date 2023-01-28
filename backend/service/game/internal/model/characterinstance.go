package model

import (
	"fmt"

	coreerror "gitlab.com/alienspaces/go-mud/backend/core/error"
	"gitlab.com/alienspaces/go-mud/backend/core/nullstring"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/record"
)

type CharacterInstanceRecordSet struct {
	CharacterInstanceRec *record.CharacterInstance
	ObjectInstanceRecs   []*record.ObjectInstance
}

// CharacterEnterDungeon -
func (m *Model) CharacterEnterDungeon(dungeonID, characterID string) (*CharacterInstanceRecordSet, error) {
	l := m.Logger("CharacterEnterDungeon")

	dungeonInstance, err := m.GetAvailableDungeonInstanceViewRecordSet(dungeonID)
	if err != nil {
		l.Warn("failed getting an available dungeon instance >%v<", err)
		return nil, err
	}

	if dungeonInstance == nil {
		msg := "dungeon instance is nil, failed getting an available dungeon instance"
		l.Warn(msg)
		return nil, fmt.Errorf(msg)
	}

	// Find an entrance
	characterEntered := false
	characterInstanceRecordSet := &CharacterInstanceRecordSet{}
	for _, locationInstanceRec := range dungeonInstance.LocationInstanceViewRecs {
		if locationInstanceRec.IsDefault {
			characterInstanceRecordSet, err = m.CreateCharacterInstance(locationInstanceRec.ID, characterID)
			if err != nil {
				l.Warn("failed creating character instance >%v<", err)
				return nil, err
			}
			characterEntered = true
		}
	}

	if !characterEntered {
		msg := "failed to enter character into dungeon"
		l.Warn(msg)
		return nil, fmt.Errorf(msg)
	}

	return characterInstanceRecordSet, nil
}

// CharacterExitDungeon -
func (m *Model) CharacterExitDungeon(characterID string) error {
	l := m.Logger("CharacterExitDungeon")

	// Update character record
	characterInstanceRec, err := m.GetCharacterInstanceRecByCharacterID(characterID)
	if err != nil {
		l.Warn("failed to get character ID >%s< instance view record >%v<", characterID, err)
		return err
	}

	if characterInstanceRec == nil {
		l.Warn("character instance record is nil")
		err := coreerror.NewServerInternalError()
		return err
	}

	characterRec, err := m.GetCharacterRec(characterID, true)
	if err != nil {
		l.Warn("failed to get character ID >%s< record >%v<", characterID, err)
		return err
	}

	characterRec.ExperiencePoints = characterInstanceRec.ExperiencePoints
	characterRec.AttributePoints = characterInstanceRec.AttributePoints
	characterRec.Coins = characterInstanceRec.Coins

	err = m.UpdateCharacterRec(characterRec)
	if err != nil {
		l.Warn("failed to update character ID >%s< record >%v<", characterID, err)
		return err
	}

	// Replace character object records
	characterObjectRecs, err := m.GetCharacterObjectRecs(
		map[string]interface{}{
			"character_id": characterID,
		}, nil, false,
	)
	if err != nil {
		l.Warn("failed to get character ID >%s< object records >%v<", characterID, err)
		return err
	}

	characterObjectInstanceRecs, err := m.GetCharacterInstanceObjectInstanceRecs(characterInstanceRec.ID)
	if err != nil {
		l.Warn("failed to get character ID >%s< object instance records >%v<", characterID, err)
		return err
	}

	// Delete character object records for objects the character no longer has
CHARACTER_OBJECT_RECS:
	for idx := range characterObjectRecs {
		for iidx := range characterObjectInstanceRecs {
			if characterObjectRecs[idx].ObjectID == characterObjectInstanceRecs[iidx].ObjectID {
				continue CHARACTER_OBJECT_RECS
			}
		}
		l.Info("Deleting character object record ID >%s<", characterObjectRecs[idx].ID)
		err := m.DeleteCharacterObjectRec(characterObjectRecs[idx].ID)
		if err != nil {
			l.Warn("failed deleting character ID >%s< object ID >%s< record >%v<", characterObjectRecs[idx].CharacterID, characterObjectRecs[idx].ObjectID, err)
			return err
		}
	}

	// Update character object reocrds or create missing character object records the character now has
CHARACTER_OBJECT_INSTANCE_RECS:
	for iidx := range characterObjectInstanceRecs {
		for idx := range characterObjectRecs {
			if characterObjectRecs[idx].ObjectID == characterObjectInstanceRecs[iidx].ObjectID {
				l.Info("Updating character object record ID >%s<", characterObjectRecs[idx].ID)
				characterObjectRec := characterObjectRecs[idx]
				characterObjectRec.IsStashed = characterObjectInstanceRecs[iidx].IsStashed
				characterObjectRec.IsEquipped = characterObjectInstanceRecs[iidx].IsEquipped
				err := m.UpdateCharacterObjectRec(characterObjectRec)
				if err != nil {
					l.Warn("failed updating character ID >%s< object ID >%s< record >%v<", characterObjectRec.CharacterID, characterObjectRec.ObjectID, err)
					return err
				}
				continue CHARACTER_OBJECT_INSTANCE_RECS
			}
		}

		l.Info("Creating character ID >5s< object record ID >%s<", characterID, characterObjectInstanceRecs[iidx].ObjectID)
		characterObjectRec := record.CharacterObject{
			CharacterID: characterID,
			ObjectID:    characterObjectInstanceRecs[iidx].ObjectID,
			IsStashed:   characterObjectInstanceRecs[iidx].IsStashed,
			IsEquipped:  characterObjectInstanceRecs[iidx].IsEquipped,
		}
		err := m.CreateCharacterObjectRec(&characterObjectRec)
		if err != nil {
			l.Warn("failed creating character ID >%s< object ID >%s< record >%v<", characterObjectRec.CharacterID, characterObjectRec.ObjectID, err)
			return err
		}
	}

	// Delete character instance
	err = m.DeleteCharacterInstance(characterID)
	if err != nil {
		l.Warn("failed to delete character ID >%s< instance records >%v<", characterID, err)
		return err
	}

	return nil
}

// CreateCharacterInstance -
func (m *Model) CreateCharacterInstance(locationInstanceID string, characterID string) (*CharacterInstanceRecordSet, error) {
	l := m.Logger("CreateCharacterInstance")

	locationInstanceRec, err := m.GetLocationInstanceRec(locationInstanceID, false)
	if err != nil {
		l.Warn("failed getting location instance record >%v<", err)
		return nil, err
	}

	characterRec, err := m.GetCharacterRec(characterID, false)
	if err != nil {
		l.Warn("failed getting character record >%v<", err)
		return nil, err
	}

	characterInstanceRec := &record.CharacterInstance{
		CharacterID:        characterRec.ID,
		DungeonInstanceID:  locationInstanceRec.DungeonInstanceID,
		LocationInstanceID: locationInstanceRec.ID,
		Strength:           characterRec.Strength,
		Dexterity:          characterRec.Dexterity,
		Intelligence:       characterRec.Intelligence,
		Health:             characterRec.Health,
		Fatigue:            characterRec.Fatigue,
		Coins:              characterRec.Coins,
		ExperiencePoints:   characterRec.ExperiencePoints,
		AttributePoints:    characterRec.AttributePoints,
	}

	err = m.validateCreateCharacterInstanceRec(characterInstanceRec)
	if err != nil {
		l.Warn("failed validating character instance record >%v<", err)
		return nil, err
	}

	err = m.CreateCharacterInstanceRec(characterInstanceRec)
	if err != nil {
		l.Warn("failed creating character instance record >%v<", err)
		return nil, err
	}

	characterObjectRecs, err := m.GetCharacterObjectRecs(
		map[string]interface{}{
			"character_id": characterID,
		}, nil, false,
	)
	if err != nil {
		l.Warn("failed getting character object records >%v<", err)
		return nil, err
	}

	objectInstanceRecs := []*record.ObjectInstance{}

	for _, characterObjectRec := range characterObjectRecs {

		objectInstanceRec := &record.ObjectInstance{
			ObjectID:            characterObjectRec.ObjectID,
			DungeonInstanceID:   locationInstanceRec.DungeonInstanceID,
			CharacterInstanceID: nullstring.FromString(characterInstanceRec.ID),
			IsStashed:           characterObjectRec.IsStashed,
			IsEquipped:          characterObjectRec.IsEquipped,
		}

		err := m.CreateObjectInstanceRec(objectInstanceRec)
		if err != nil {
			l.Warn("failed creating object instance record >%v<", err)
			return nil, err
		}

		objectInstanceRecs = append(objectInstanceRecs, objectInstanceRec)
	}

	characterInstanceRecordSet := CharacterInstanceRecordSet{
		CharacterInstanceRec: characterInstanceRec,
		ObjectInstanceRecs:   objectInstanceRecs,
	}

	return &characterInstanceRecordSet, nil
}

// DeleteCharacterInstance -
func (m *Model) DeleteCharacterInstance(characterID string) error {
	l := m.Logger("DeleteCharacterInstance")

	characterInstanceRec, err := m.GetCharacterInstanceViewRecByCharacterID(characterID)
	if err != nil {
		l.Warn("failed getting character instance view record >%v<", err)
		return err
	}

	characterObjectInstanceRecs, err := m.GetCharacterInstanceObjectInstanceRecs(characterInstanceRec.CharacterID)
	if err != nil {
		l.Warn("failed getting character object instance view records >%v<", err)
		return err
	}

	for idx := range characterObjectInstanceRecs {
		err := m.DeleteObjectInstanceRec(characterObjectInstanceRecs[idx].ID)
		if err != nil {
			l.Warn("failed deleting character object instance record >%v<", err)
			return err
		}
	}

	err = m.DeleteCharacterInstanceRec(characterInstanceRec.ID)
	if err != nil {
		l.Warn("failed deleting character instance record >%v<", err)
		return err
	}

	return nil
}

// GetCharacterInstanceObjectInstanceRecs -
func (m *Model) GetCharacterInstanceObjectInstanceRecs(characterInstanceID string) ([]*record.ObjectInstance, error) {

	m.Log.Info("Getting character instance ID >%s< object records", characterInstanceID)

	r := m.ObjectInstanceRepository()

	return r.GetMany(map[string]interface{}{
		"character_instance_id": characterInstanceID,
	}, nil, false)
}

// GetCharacterInstanceEquippedObjectInstanceRecs -
func (m *Model) GetCharacterInstanceEquippedObjectInstanceRecs(characterInstanceID string) ([]*record.ObjectInstance, error) {
	l := m.Logger("GetCharacterInstanceEquippedObjectInstanceRecs")

	l.Info("Getting character instance ID >%s< equipped object records", characterInstanceID)

	r := m.ObjectInstanceRepository()

	return r.GetMany(map[string]interface{}{
		"character_instance_id": characterInstanceID,
		"is_equipped":           true,
	}, nil, false)
}

// GetCharacterInstanceStashedObjectInstanceRecs -
func (m *Model) GetCharacterInstanceStashedObjectInstanceRecs(characterInstanceID string) ([]*record.ObjectInstance, error) {
	l := m.Logger("GetCharacterInstanceStashedObjectInstanceRecs")

	l.Info("Getting character instance ID >%s< stashed object records", characterInstanceID)

	r := m.ObjectInstanceRepository()

	return r.GetMany(map[string]interface{}{
		"character_instance_id": characterInstanceID,
		"is_stashed":            true,
	}, nil, false)
}

// GetCharacterInstanceObjectInstanceViewRecs -
func (m *Model) GetCharacterInstanceObjectInstanceViewRecs(characterInstanceID string) ([]*record.ObjectInstanceView, error) {
	l := m.Logger("GetCharacterInstanceObjectInstanceViewRecs")

	l.Info("Getting character instance ID >%s< object records", characterInstanceID)

	r := m.ObjectInstanceViewRepository()

	return r.GetMany(map[string]interface{}{
		"character_instance_id": characterInstanceID,
	}, nil, false)
}

// GetCharacterInstanceEquippedObjectInstanceViewRecs -
func (m *Model) GetCharacterInstanceEquippedObjectInstanceViewRecs(characterInstanceID string) ([]*record.ObjectInstanceView, error) {
	l := m.Logger("GetCharacterInstanceEquippedObjectInstanceViewRecs")

	l.Info("Getting character instance ID >%s< equipped object records", characterInstanceID)

	r := m.ObjectInstanceViewRepository()

	return r.GetMany(map[string]interface{}{
		"character_instance_id": characterInstanceID,
		"is_equipped":           true,
	}, nil, false)
}

// GetCharacterInstanceStashedObjectInstanceViewRecs -
func (m *Model) GetCharacterInstanceStashedObjectInstanceViewRecs(characterInstanceID string) ([]*record.ObjectInstanceView, error) {
	l := m.Logger("GetCharacterInstanceStashedObjectInstanceViewRecs")

	l.Info("Getting character instance ID >%s< stashed object records", characterInstanceID)

	r := m.ObjectInstanceViewRepository()

	return r.GetMany(map[string]interface{}{
		"character_instance_id": characterInstanceID,
		"is_stashed":            true,
	}, nil, false)
}
