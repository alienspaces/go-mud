package model

import (
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// GetCharacterEquippedObjectRecs -
func (m *Model) GetCharacterEquippedObjectRecs(characterID string) ([]*record.Object, error) {
	l := m.Logger("GetCharacterEquippedObjectRecs")

	l.Debug("Getting character ID >%s< equipped object records", characterID)

	cor := m.CharacterObjectRepository()
	or := m.ObjectRepository()

	characterObjectRecs, err := cor.GetMany(map[string]interface{}{
		"character_id": characterID,
		"is_equipped":  true,
	}, nil, false)
	if err != nil {
		return nil, err
	}

	objectRecs := []*record.Object{}
	for _, characterObjectRec := range characterObjectRecs {
		objectRec, err := or.GetOne(characterObjectRec.ObjectID, false)
		if err != nil {
			return nil, err
		}
		objectRecs = append(objectRecs, objectRec)
	}

	return objectRecs, nil
}

// GetCharacterStashedObjectRecs -
func (m *Model) GetCharacterStashedObjectRecs(characterID string) ([]*record.Object, error) {

	l := m.Logger("GetCharacterStashedObjectRecs")

	l.Debug("Getting character ID >%s< stashed object records", characterID)

	cor := m.CharacterObjectRepository()
	or := m.ObjectRepository()

	characterObjectRecs, err := cor.GetMany(map[string]interface{}{
		"character_id": characterID,
		"is_stashed":   true,
	}, nil, false)
	if err != nil {
		return nil, err
	}

	objectRecs := []*record.Object{}
	for _, characterObjectRec := range characterObjectRecs {
		objectRec, err := or.GetOne(characterObjectRec.ObjectID, false)
		if err != nil {
			return nil, err
		}
		objectRecs = append(objectRecs, objectRec)
	}

	return objectRecs, nil
}
