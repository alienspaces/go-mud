package model

import (
	"fmt"

	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// TODO: Possibly hang methods off this structs to add, remove records etc
type CharacterInstanceRecordSet struct {
	CharacterInstanceRec *record.CharacterInstance
	ObjectInstanceRecs   []*record.ObjectInstance
}

// CharacterEnterDungeon -
func (m *Model) CharacterEnterDungeon(dungeonID, characterID string) (*DungeonInstanceRecordSet, error) {
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
	for _, locationInstanceRec := range dungeonInstance.LocationInstanceViewRecs {
		if locationInstanceRec.IsDefault {
			_, err := m.CreateCharacterInstance(dungeonInstance.DungeonInstanceViewRec.ID, locationInstanceRec.ID, characterID)
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

	return m.GetDungeonInstanceRecordSet(dungeonInstance.DungeonInstanceViewRec.ID)
}

// CharacterExitDungeon -
func (m *Model) CharacterExitDungeon(characterID string) error {

	//
	return nil
}

func (m *Model) CreateCharacterInstance(dungeonInstanceID string, locationInstanceID string, characterID string) (*CharacterInstanceRecordSet, error) {

	// TODO: Create character instance, return character instance

	return nil, nil
}

// GetCharacterInstanceObjectInstanceRecs -
func (m *Model) GetCharacterInstanceObjectInstanceRecs(characterID string) ([]*record.ObjectInstance, error) {

	m.Log.Info("Getting character ID >%s< object records", characterID)

	r := m.ObjectInstanceRepository()

	return r.GetMany(map[string]interface{}{
		"character_instance_id": characterID,
	}, nil, false)
}

// GetCharacterInstanceEquippedObjectInstanceRecs -
func (m *Model) GetCharacterInstanceEquippedObjectInstanceRecs(characterID string) ([]*record.ObjectInstance, error) {

	m.Log.Info("Getting character ID >%s< equipped object records", characterID)

	r := m.ObjectInstanceRepository()

	return r.GetMany(map[string]interface{}{
		"character_instance_id": characterID,
		"is_equipped":           true,
	}, nil, false)
}

// GetCharacterInstanceStashedObjectInstanceRecs -
func (m *Model) GetCharacterInstanceStashedObjectInstanceRecs(characterID string) ([]*record.ObjectInstance, error) {

	m.Log.Info("Getting character ID >%s< stashed object records", characterID)

	r := m.ObjectInstanceRepository()

	return r.GetMany(map[string]interface{}{
		"character_instance_id": characterID,
		"is_stashed":            true,
	}, nil, false)
}

// GetCharacterInstanceObjectInstanceViewRecs -
func (m *Model) GetCharacterInstanceObjectInstanceViewRecs(characterID string) ([]*record.ObjectInstanceView, error) {

	m.Log.Info("Getting character ID >%s< object records", characterID)

	r := m.ObjectInstanceViewRepository()

	return r.GetMany(map[string]interface{}{
		"character_instance_id": characterID,
	}, nil, false)
}

// GetCharacterInstanceEquippedObjectInstanceViewRecs -
func (m *Model) GetCharacterInstanceEquippedObjectInstanceViewRecs(characterID string) ([]*record.ObjectInstanceView, error) {

	m.Log.Info("Getting character ID >%s< equipped object records", characterID)

	r := m.ObjectInstanceViewRepository()

	return r.GetMany(map[string]interface{}{
		"character_instance_id": characterID,
		"is_equipped":           true,
	}, nil, false)
}

// GetCharacterInstanceStashedObjectInstanceViewRecs -
func (m *Model) GetCharacterInstanceStashedObjectInstanceViewRecs(characterID string) ([]*record.ObjectInstanceView, error) {

	m.Log.Info("Getting character ID >%s< stashed object records", characterID)

	r := m.ObjectInstanceViewRepository()

	return r.GetMany(map[string]interface{}{
		"character_instance_id": characterID,
		"is_stashed":            true,
	}, nil, false)
}
