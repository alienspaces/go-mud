package model

import (
	coresql "gitlab.com/alienspaces/go-mud/backend/core/sql"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/record"
)

// GetCharacterEquippedObjectRecs -
func (m *Model) GetCharacterEquippedObjectRecs(characterID string) ([]*record.Object, error) {
	l := m.loggerWithContext("GetCharacterEquippedObjectRecs")

	l.Debug("Getting character ID >%s< equipped object records", characterID)

	cor := m.CharacterObjectRepository()
	or := m.ObjectRepository()

	characterObjectRecs, err := cor.GetMany(
		&coresql.Options{
			Params: []coresql.Param{
				{
					Col: "character_id",
					Val: characterID,
				},
				{
					Col: "is_equipped",
					Val: true,
				},
			},
		},
	)
	if err != nil {
		return nil, err
	}

	objectRecs := []*record.Object{}
	for _, characterObjectRec := range characterObjectRecs {
		objectRec, err := or.GetOne(characterObjectRec.ObjectID, nil)
		if err != nil {
			return nil, err
		}
		objectRecs = append(objectRecs, objectRec)
	}

	return objectRecs, nil
}

// GetCharacterStashedObjectRecs -
func (m *Model) GetCharacterStashedObjectRecs(characterID string) ([]*record.Object, error) {

	l := m.loggerWithContext("GetCharacterStashedObjectRecs")

	l.Debug("Getting character ID >%s< stashed object records", characterID)

	cor := m.CharacterObjectRepository()
	or := m.ObjectRepository()

	characterObjectRecs, err := cor.GetMany(
		&coresql.Options{
			Params: []coresql.Param{
				{
					Col: "character_id",
					Val: characterID,
				},
				{
					Col: "is_stashed",
					Val: true,
				},
			},
		},
	)
	if err != nil {
		return nil, err
	}

	objectRecs := []*record.Object{}
	for _, characterObjectRec := range characterObjectRecs {
		objectRec, err := or.GetOne(characterObjectRec.ObjectID, nil)
		if err != nil {
			return nil, err
		}
		objectRecs = append(objectRecs, objectRec)
	}

	return objectRecs, nil
}
