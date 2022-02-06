package model

import (
	"fmt"
	"strings"

	"gitlab.com/alienspaces/go-mud/server/core/store"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

type ResolverSentence struct {
	Command  string
	Sentence string
}

func (m *Model) resolveAction(sentence string, dungeonCharacterRec *record.DungeonCharacter, dungeonLocationRecordSet *DungeonLocationRecordSet) (*record.DungeonAction, error) {

	resolved, err := m.resolveCommand(sentence)
	if err != nil {
		m.Log.Warn("Failed resolving command >%v<", err)
		return nil, err
	}

	resolveFuncs := map[string]func(sentence string, dungeonCharacterRec *record.DungeonCharacter, dungeonLocationRecordSet *DungeonLocationRecordSet) (*record.DungeonAction, error){
		"move":  m.resolveMoveAction,
		"look":  m.resolveLookAction,
		"stash": m.resolveStashAction,
		"equip": m.resolveEquipAction,
		"drop":  m.resolveDropAction,
	}

	resolveFunc, ok := resolveFuncs[resolved.Command]
	if !ok {
		msg := fmt.Sprintf("Command >%s< could not be resolved", resolved.Command)
		m.Log.Warn(msg)
		return nil, fmt.Errorf(msg)
	}

	dungeonActionRec, err := resolveFunc(resolved.Sentence, dungeonCharacterRec, dungeonLocationRecordSet)
	if err != nil {
		m.Log.Warn("Failed resolver function for command >%s< >%v<", resolved.Command, err)
		return nil, err
	}

	m.Log.Warn("Resolved dungeon action rec >%#v<", dungeonActionRec)

	return dungeonActionRec, nil
}

func (m *Model) resolveCommand(sentence string) (*ResolverSentence, error) {
	sentenceWords := strings.Split(sentence, " ")
	resolved := ResolverSentence{}

	m.Log.Debug("Have sentence words >%v<", sentenceWords)

	for _, dungeonAction := range []string{"move", "look", "equip", "stash", "drop"} {
		m.Log.Debug("Checking dungeon action >%s<", dungeonAction)

		// NOTE: The appended space is important
		if strings.Contains(sentence, dungeonAction+" ") {
			m.Log.Debug("Sentence contains action >%s<", dungeonAction)
			sentence = strings.Replace(sentence, dungeonAction+" ", "", 1)
			resolved.Command = dungeonAction
			resolved.Sentence = sentence
		} else if sentence == dungeonAction {
			m.Log.Debug("Sentence equals action >%s<", dungeonAction)
			sentence = strings.Replace(sentence, dungeonAction, "", 1)
			resolved.Command = dungeonAction
			resolved.Sentence = sentence
		}
	}

	m.Log.Info("Resolved command >%#v<", resolved)

	return &resolved, nil
}

func (m *Model) resolveMoveAction(sentence string, dungeonCharacterRec *record.DungeonCharacter, records *DungeonLocationRecordSet) (*record.DungeonAction, error) {

	var err error
	var targetDungeonLocationID string
	var targetDungeonLocationDirection string

	if sentence != "" {
		targetDungeonLocationID, targetDungeonLocationDirection, err = m.resolveSentenceLocationDirection(sentence, records.LocationRec)
		if err != nil {
			m.Log.Warn("Failed to resolve sentence location direction >%v<", err)
			return nil, err
		}
	}

	if targetDungeonLocationID == "" {
		msg := fmt.Sprintf("failed to resolve target dungeon location id with sentence >%s<", sentence)
		m.Log.Warn(msg)
		return nil, fmt.Errorf(msg)
	}

	dungeonActionRecord := record.DungeonAction{
		DungeonID:                              dungeonCharacterRec.DungeonID,
		DungeonLocationID:                      dungeonCharacterRec.DungeonLocationID,
		DungeonCharacterID:                     store.NullString(dungeonCharacterRec.ID),
		ResolvedCommand:                        "move",
		ResolvedTargetDungeonLocationDirection: store.NullString(targetDungeonLocationDirection),
		ResolvedTargetDungeonLocationID:        store.NullString(targetDungeonLocationID),
	}

	return &dungeonActionRecord, nil
}

func (m *Model) resolveLookAction(sentence string, dungeonCharacterRec *record.DungeonCharacter, locationRecordSet *DungeonLocationRecordSet) (*record.DungeonAction, error) {

	var err error
	var targetDungeonLocationID string
	var targetDungeonLocationDirection string
	var targetDungeonObjectID string
	var targetDungeonMonsterID string
	var targetDungeonCharacterID string

	if sentence != "" {
		targetDungeonLocationID, targetDungeonLocationDirection, err = m.resolveSentenceLocationDirection(sentence, locationRecordSet.LocationRec)
		if err != nil {
			m.Log.Warn("Failed to resolve sentence location direction >%v<", err)
			return nil, err
		}

		if targetDungeonLocationID == "" {
			dungeonObjectRec, err := m.getObjectFromSentence(sentence, locationRecordSet.ObjectRecs)
			if err != nil {
				m.Log.Warn("Failed to resolve sentence object >%v<", err)
				return nil, err
			}
			if dungeonObjectRec != nil {
				targetDungeonObjectID = dungeonObjectRec.ID
			}
		}

		if targetDungeonLocationID == "" && targetDungeonObjectID == "" {
			dungeonMonsterRec, err := m.resolveSentenceMonster(sentence, locationRecordSet.MonsterRecs)
			if err != nil {
				m.Log.Warn("Failed to resolve sentence monster >%v<", err)
				return nil, err
			}
			if dungeonMonsterRec != nil {
				targetDungeonMonsterID = dungeonMonsterRec.ID
			}
		}

		if targetDungeonLocationID == "" && targetDungeonObjectID == "" && targetDungeonMonsterID == "" {
			dungeonCharacterRec, err := m.resolveSentenceCharacter(sentence, locationRecordSet.CharacterRecs)
			if err != nil {
				m.Log.Warn("Failed to resolve sentence character >%v<", err)
				return nil, err
			}
			if dungeonCharacterRec != nil {
				targetDungeonCharacterID = dungeonCharacterRec.ID
			}
		}
	}

	// When nothing has been identified, assume we are just looking in the current room.
	if targetDungeonLocationID == "" && targetDungeonObjectID == "" && targetDungeonMonsterID == "" && targetDungeonCharacterID == "" {
		targetDungeonLocationID = locationRecordSet.LocationRec.ID
		targetDungeonLocationDirection = ""
	}

	dungeonActionRecord := record.DungeonAction{
		DungeonID:                              dungeonCharacterRec.DungeonID,
		DungeonLocationID:                      dungeonCharacterRec.DungeonLocationID,
		DungeonCharacterID:                     store.NullString(dungeonCharacterRec.ID),
		ResolvedCommand:                        "look",
		ResolvedTargetDungeonObjectID:          store.NullString(targetDungeonObjectID),
		ResolvedTargetDungeonMonsterID:         store.NullString(targetDungeonMonsterID),
		ResolvedTargetDungeonCharacterID:       store.NullString(targetDungeonCharacterID),
		ResolvedTargetDungeonLocationDirection: store.NullString(targetDungeonLocationDirection),
		ResolvedTargetDungeonLocationID:        store.NullString(targetDungeonLocationID),
	}

	return &dungeonActionRecord, nil
}

func (m *Model) resolveStashAction(sentence string, dungeonCharacterRec *record.DungeonCharacter, locationRecordSet *DungeonLocationRecordSet) (*record.DungeonAction, error) {

	var stashedDungeonObjectID string

	if sentence != "" {
		// Find object in room
		dungeonObjectRec, err := m.getObjectFromSentence(sentence, locationRecordSet.ObjectRecs)
		if err != nil {
			m.Log.Warn("Failed to get location object from sentence >%v<", err)
			return nil, err
		}
		if dungeonObjectRec == nil {
			// Find object equipped on character
			dungeonObjectRecs, err := m.GetCharacterEquippedDungeonObjectRecs(dungeonCharacterRec.ID)
			if err != nil {
				m.Log.Warn("Failed to get character equipped objects >%v<", err)
				return nil, err
			}
			dungeonObjectRec, err = m.getObjectFromSentence(sentence, dungeonObjectRecs)
			if err != nil {
				m.Log.Warn("Failed to get character object from sentence >%v<", err)
				return nil, err
			}
		}
		if dungeonObjectRec != nil {
			stashedDungeonObjectID = dungeonObjectRec.ID
		}
	}

	dungeonActionRec := record.DungeonAction{
		DungeonID:                      dungeonCharacterRec.DungeonID,
		DungeonLocationID:              dungeonCharacterRec.DungeonLocationID,
		DungeonCharacterID:             store.NullString(dungeonCharacterRec.ID),
		ResolvedCommand:                "stash",
		ResolvedTargetDungeonObjectID:  store.NullString(stashedDungeonObjectID),
		ResolvedStashedDungeonObjectID: store.NullString(stashedDungeonObjectID),
	}

	return &dungeonActionRec, nil
}

func (m *Model) resolveEquipAction(sentence string, dungeonCharacterRec *record.DungeonCharacter, locationRecordSet *DungeonLocationRecordSet) (*record.DungeonAction, error) {

	var equippedDungeonObjectID string

	if sentence != "" {
		// Find object in room
		dungeonObjectRec, err := m.getObjectFromSentence(sentence, locationRecordSet.ObjectRecs)
		if err != nil {
			m.Log.Warn("Failed to get location object from sentence >%v<", err)
			return nil, err
		}
		if dungeonObjectRec == nil {
			// Find object stashed on character
			dungeonObjectRecs, err := m.GetCharacterStashedDungeonObjectRecs(dungeonCharacterRec.ID)
			if err != nil {
				m.Log.Warn("Failed to get character stashed objects >%v<", err)
				return nil, err
			}
			dungeonObjectRec, err = m.getObjectFromSentence(sentence, dungeonObjectRecs)
			if err != nil {
				m.Log.Warn("Failed to get character object from sentence >%v<", err)
				return nil, err
			}
		}
		if dungeonObjectRec != nil {
			equippedDungeonObjectID = dungeonObjectRec.ID
		}
	}

	dungeonActionRec := record.DungeonAction{
		DungeonID:                       dungeonCharacterRec.DungeonID,
		DungeonLocationID:               dungeonCharacterRec.DungeonLocationID,
		DungeonCharacterID:              store.NullString(dungeonCharacterRec.ID),
		ResolvedCommand:                 "equip",
		ResolvedTargetDungeonObjectID:   store.NullString(equippedDungeonObjectID),
		ResolvedEquippedDungeonObjectID: store.NullString(equippedDungeonObjectID),
	}

	return &dungeonActionRec, nil
}

func (m *Model) resolveDropAction(sentence string, dungeonCharacterRec *record.DungeonCharacter, locationRecordSet *DungeonLocationRecordSet) (*record.DungeonAction, error) {

	var droppedDungeonObjectID string

	if sentence != "" {
		// Find object stashed on character
		m.Log.Debug("Finding object stashed on character")
		dungeonObjectRecs, err := m.GetCharacterStashedDungeonObjectRecs(dungeonCharacterRec.ID)
		if err != nil {
			m.Log.Warn("Failed to get character stashed objects >%v<", err)
			return nil, err
		}
		dungeonObjectRec, err := m.getObjectFromSentence(sentence, dungeonObjectRecs)
		if err != nil {
			m.Log.Warn("Failed to get character object from sentence >%v<", err)
			return nil, err
		}
		m.Log.Debug("Found object >%v< stashed on character", dungeonObjectRec)
		if dungeonObjectRec == nil {
			// Find object equipped on character
			m.Log.Debug("Finding object equipped on character")
			dungeonObjectRecs, err := m.GetCharacterEquippedDungeonObjectRecs(dungeonCharacterRec.ID)
			if err != nil {
				m.Log.Warn("Failed to get character equipped objects >%v<", err)
				return nil, err
			}
			dungeonObjectRec, err = m.getObjectFromSentence(sentence, dungeonObjectRecs)
			if err != nil {
				m.Log.Warn("Failed to get character object from sentence >%v<", err)
				return nil, err
			}
			m.Log.Debug("Found object >%v< equipped on character", dungeonObjectRec)
		}
		if dungeonObjectRec != nil {
			droppedDungeonObjectID = dungeonObjectRec.ID
		}
	}

	dungeonActionRec := record.DungeonAction{
		DungeonID:                      dungeonCharacterRec.DungeonID,
		DungeonLocationID:              dungeonCharacterRec.DungeonLocationID,
		DungeonCharacterID:             store.NullString(dungeonCharacterRec.ID),
		ResolvedCommand:                "drop",
		ResolvedTargetDungeonObjectID:  store.NullString(droppedDungeonObjectID),
		ResolvedDroppedDungeonObjectID: store.NullString(droppedDungeonObjectID),
	}

	return &dungeonActionRec, nil
}

func (m *Model) resolveSentenceLocationDirection(sentence string, dungeonLocationRec *record.DungeonLocation) (string, string, error) {

	var dungeonLocationID string
	var dungeonLocationDirection string

	if dungeonLocationRec.NortheastDungeonLocationID.Valid && strings.Contains(sentence, "northeast") {
		dungeonLocationID = dungeonLocationRec.NortheastDungeonLocationID.String
		dungeonLocationDirection = "northeast"
	} else if dungeonLocationRec.NorthwestDungeonLocationID.Valid && strings.Contains(sentence, "northwest") {
		dungeonLocationID = dungeonLocationRec.NorthwestDungeonLocationID.String
		dungeonLocationDirection = "northwest"
	} else if dungeonLocationRec.SoutheastDungeonLocationID.Valid && strings.Contains(sentence, "southeast") {
		dungeonLocationID = dungeonLocationRec.SoutheastDungeonLocationID.String
		dungeonLocationDirection = "southeast"
	} else if dungeonLocationRec.SoutheastDungeonLocationID.Valid && strings.Contains(sentence, "southeast") {
		dungeonLocationID = dungeonLocationRec.SoutheastDungeonLocationID.String
		dungeonLocationDirection = "southeast"
	} else if dungeonLocationRec.NorthDungeonLocationID.Valid && strings.Contains(sentence, "north") {
		dungeonLocationID = dungeonLocationRec.NorthDungeonLocationID.String
		dungeonLocationDirection = "north"
	} else if dungeonLocationRec.EastDungeonLocationID.Valid && strings.Contains(sentence, "east") {
		dungeonLocationID = dungeonLocationRec.EastDungeonLocationID.String
		dungeonLocationDirection = "east"
	} else if dungeonLocationRec.SouthDungeonLocationID.Valid && strings.Contains(sentence, "south") {
		dungeonLocationID = dungeonLocationRec.SouthDungeonLocationID.String
		dungeonLocationDirection = "south"
	} else if dungeonLocationRec.WestDungeonLocationID.Valid && strings.Contains(sentence, "west") {
		dungeonLocationID = dungeonLocationRec.WestDungeonLocationID.String
		dungeonLocationDirection = "west"
	} else if dungeonLocationRec.UpDungeonLocationID.Valid && strings.Contains(sentence, "up") {
		dungeonLocationID = dungeonLocationRec.UpDungeonLocationID.String
		dungeonLocationDirection = "up"
	} else if dungeonLocationRec.DownDungeonLocationID.Valid && strings.Contains(sentence, "down") {
		dungeonLocationID = dungeonLocationRec.DownDungeonLocationID.String
		dungeonLocationDirection = "down"
	}

	return dungeonLocationID, dungeonLocationDirection, nil
}

func (m *Model) getObjectFromSentence(sentence string, dungeonObjectRecs []*record.DungeonObject) (*record.DungeonObject, error) {
	for _, dungeonObjectRec := range dungeonObjectRecs {
		if strings.Contains(sentence, strings.ToLower(dungeonObjectRec.Name)) {
			return dungeonObjectRec, nil
		}
	}
	return nil, nil
}

func (m *Model) resolveSentenceMonster(sentence string, dungeonMonsterRecs []*record.DungeonMonster) (*record.DungeonMonster, error) {
	for _, dungeonMonsterRec := range dungeonMonsterRecs {
		if strings.Contains(sentence, strings.ToLower(dungeonMonsterRec.Name)) {
			return dungeonMonsterRec, nil
		}
	}
	return nil, nil
}

func (m *Model) resolveSentenceCharacter(sentence string, dungeonCharacterRecs []*record.DungeonCharacter) (*record.DungeonCharacter, error) {
	for _, dungeonCharacterRec := range dungeonCharacterRecs {
		if strings.Contains(sentence, strings.ToLower(dungeonCharacterRec.Name)) {
			return dungeonCharacterRec, nil
		}
	}
	return nil, nil
}
