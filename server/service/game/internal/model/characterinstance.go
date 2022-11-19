package model

import (
	"fmt"

	"gitlab.com/alienspaces/go-mud/server/core/nullstring"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
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

	// TODO: (game) Implement exiting a dungeon instance..
	l.Warn("Not implemented")

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
