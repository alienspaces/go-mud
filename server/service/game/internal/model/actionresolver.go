package model

import (
	"fmt"
	"strings"

	"gitlab.com/alienspaces/go-mud/server/core/nullstring"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

var dungeonActionResolvedCommands []string = []string{
	record.ActionCommandMove,
	record.ActionCommandLook,
	record.ActionCommandStash,
	record.ActionCommandEquip,
	record.ActionCommandDrop,
	// record.ActionCommandAttack,
}

type ResolverSentence struct {
	Command  string
	Sentence string
}

func (m *Model) resolveAction(
	sentence string,
	characterInstanceViewRec *record.CharacterInstanceView,
	locationInstanceRecordSet *record.LocationInstanceViewRecordSet) (*record.Action, error) {

	l := m.Logger("resolveAction")

	resolved, err := m.resolveCommand(sentence)
	if err != nil {
		l.Warn("failed resolving command >%v<", err)
		return nil, err
	}

	resolveFuncs := map[string]func(sentence string, characterInstanceViewRec *record.CharacterInstanceView, locationInstanceRecordSet *record.LocationInstanceViewRecordSet) (*record.Action, error){
		record.ActionCommandMove:  m.resolveMoveAction,
		record.ActionCommandLook:  m.resolveLookAction,
		record.ActionCommandStash: m.resolveStashAction,
		record.ActionCommandEquip: m.resolveEquipAction,
		record.ActionCommandDrop:  m.resolveDropAction,
	}

	resolveFunc, ok := resolveFuncs[resolved.Command]
	if !ok {
		msg := fmt.Sprintf("Command >%s< could not be resolved", resolved.Command)
		l.Warn(msg)
		return nil, fmt.Errorf(msg)
	}

	dungeonActionRec, err := resolveFunc(resolved.Sentence, characterInstanceViewRec, locationInstanceRecordSet)
	if err != nil {
		l.Warn("failed resolver function for command >%s< >%v<", resolved.Command, err)
		return nil, err
	}

	l.Debug("Resolved dungeon action record >%#v<", dungeonActionRec)

	return dungeonActionRec, nil
}

func (m *Model) resolveCommand(sentence string) (*ResolverSentence, error) {
	sentenceWords := strings.Split(sentence, " ")
	resolved := ResolverSentence{}

	l := m.Logger("resolveCommand")

	l.Debug("Have sentence words >%v<", sentenceWords)

	for _, dungeonAction := range dungeonActionResolvedCommands {
		l.Debug("Checking dungeon action >%s<", dungeonAction)

		// NOTE: The appended space is important
		if strings.Contains(sentence, dungeonAction+" ") {
			l.Debug("Sentence contains action >%s<", dungeonAction)
			sentence = strings.Replace(sentence, dungeonAction+" ", "", 1)
			resolved.Command = dungeonAction
			resolved.Sentence = sentence
		} else if sentence == dungeonAction {
			l.Debug("Sentence equals action >%s<", dungeonAction)
			sentence = strings.Replace(sentence, dungeonAction, "", 1)
			resolved.Command = dungeonAction
			resolved.Sentence = sentence
		}
	}

	l.Debug("Resolved command >%#v<", resolved)

	return &resolved, nil
}

func (m *Model) resolveMoveAction(sentence string, characterInstanceViewRec *record.CharacterInstanceView, records *record.LocationInstanceViewRecordSet) (*record.Action, error) {

	l := m.Logger("resolveMoveAction")

	var err error
	var targetLocationInstanceID string
	var targetLocationDirection string

	if sentence != "" {
		targetLocationInstanceID, targetLocationDirection, err = m.resolveSentenceLocationDirection(sentence, records.LocationInstanceViewRec)
		if err != nil {
			l.Warn("failed to resolve sentence location direction >%v<", err)
			return nil, err
		}
	}

	if targetLocationInstanceID == "" {
		msg := fmt.Sprintf("failed to resolve target dungeon location id with sentence >%s<", sentence)
		l.Warn(msg)
		return nil, fmt.Errorf(msg)
	}

	dungeonActionRecord := record.Action{
		DungeonInstanceID:                characterInstanceViewRec.DungeonInstanceID,
		LocationInstanceID:               characterInstanceViewRec.LocationInstanceID,
		CharacterInstanceID:              nullstring.FromString(characterInstanceViewRec.ID),
		ResolvedCommand:                  "move",
		ResolvedTargetLocationDirection:  nullstring.FromString(targetLocationDirection),
		ResolvedTargetLocationInstanceID: nullstring.FromString(targetLocationInstanceID),
	}

	return &dungeonActionRecord, nil
}

func (m *Model) resolveLookAction(
	sentence string,
	characterInstanceViewRec *record.CharacterInstanceView,
	locationRecordSet *record.LocationInstanceViewRecordSet) (*record.Action, error) {

	l := m.Logger("resolveLookAction")

	var err error
	var targetLocationInstanceID string
	var targetLocationDirection string
	var targetObjectInstanceID string
	var targetMonsterInstanceID string
	var targetCharacterInstanceID string

	if sentence != "" {
		targetLocationInstanceID, targetLocationDirection, err = m.resolveSentenceLocationDirection(sentence, locationRecordSet.LocationInstanceViewRec)
		if err != nil {
			l.Warn("failed to resolve sentence location direction >%v<", err)
			return nil, err
		}

		if targetLocationInstanceID == "" {
			dungeonObjectRec, err := m.getObjectFromSentence(sentence, locationRecordSet.ObjectInstanceViewRecs)
			if err != nil {
				l.Warn("failed to resolve sentence object >%v<", err)
				return nil, err
			}
			if dungeonObjectRec != nil {
				targetObjectInstanceID = dungeonObjectRec.ID
			}
		}

		if targetLocationInstanceID == "" && targetObjectInstanceID == "" {
			dungeonMonsterRec, err := m.resolveSentenceMonster(sentence, locationRecordSet.MonsterInstanceViewRecs)
			if err != nil {
				l.Warn("failed to resolve sentence monster >%v<", err)
				return nil, err
			}
			if dungeonMonsterRec != nil {
				targetMonsterInstanceID = dungeonMonsterRec.ID
			}
		}

		if targetLocationInstanceID == "" && targetObjectInstanceID == "" && targetMonsterInstanceID == "" {
			characterInstanceViewRec, err := m.resolveSentenceCharacter(sentence, locationRecordSet.CharacterInstanceViewRecs)
			if err != nil {
				l.Warn("failed to resolve sentence character >%v<", err)
				return nil, err
			}
			if characterInstanceViewRec != nil {
				targetCharacterInstanceID = characterInstanceViewRec.ID
			}
		}
	}

	// When nothing has been identified, assume we are just looking in the current room.
	if targetLocationInstanceID == "" && targetObjectInstanceID == "" && targetMonsterInstanceID == "" && targetCharacterInstanceID == "" {
		targetLocationInstanceID = locationRecordSet.LocationInstanceViewRec.ID
		targetLocationDirection = ""
	}

	dungeonActionRecord := record.Action{
		DungeonInstanceID:                 characterInstanceViewRec.DungeonInstanceID,
		LocationInstanceID:                characterInstanceViewRec.LocationInstanceID,
		CharacterInstanceID:               nullstring.FromString(characterInstanceViewRec.ID),
		ResolvedCommand:                   "look",
		ResolvedTargetObjectInstanceID:    nullstring.FromString(targetObjectInstanceID),
		ResolvedTargetMonsterInstanceID:   nullstring.FromString(targetMonsterInstanceID),
		ResolvedTargetCharacterInstanceID: nullstring.FromString(targetCharacterInstanceID),
		ResolvedTargetLocationDirection:   nullstring.FromString(targetLocationDirection),
		ResolvedTargetLocationInstanceID:  nullstring.FromString(targetLocationInstanceID),
	}

	return &dungeonActionRecord, nil
}

func (m *Model) resolveStashAction(sentence string, characterInstanceViewRec *record.CharacterInstanceView, locationRecordSet *record.LocationInstanceViewRecordSet) (*record.Action, error) {

	l := m.Logger("resolveStashAction")

	var stashedObjectInstanceID string

	if sentence != "" {
		// Find object in room
		objectInstanceRec, err := m.getObjectFromSentence(sentence, locationRecordSet.ObjectInstanceViewRecs)
		if err != nil {
			l.Warn("failed to get location object from sentence >%v<", err)
			return nil, err
		}
		if objectInstanceRec == nil {
			// Find object equipped on character
			objectInstanceViewRecs, err := m.GetCharacterInstanceEquippedObjectInstanceViewRecs(characterInstanceViewRec.ID)
			if err != nil {
				l.Warn("failed to get character equipped objects >%v<", err)
				return nil, err
			}
			objectInstanceRec, err = m.getObjectFromSentence(sentence, objectInstanceViewRecs)
			if err != nil {
				l.Warn("failed to get character object from sentence >%v<", err)
				return nil, err
			}
		}
		if objectInstanceRec != nil {
			stashedObjectInstanceID = objectInstanceRec.ID
		}
	}

	dungeonActionRec := record.Action{
		DungeonInstanceID:               characterInstanceViewRec.DungeonInstanceID,
		LocationInstanceID:              characterInstanceViewRec.LocationInstanceID,
		CharacterInstanceID:             nullstring.FromString(characterInstanceViewRec.ID),
		ResolvedCommand:                 "stash",
		ResolvedTargetObjectInstanceID:  nullstring.FromString(stashedObjectInstanceID),
		ResolvedStashedObjectInstanceID: nullstring.FromString(stashedObjectInstanceID),
	}

	return &dungeonActionRec, nil
}

func (m *Model) resolveEquipAction(sentence string, characterInstanceViewRec *record.CharacterInstanceView, locationRecordSet *record.LocationInstanceViewRecordSet) (*record.Action, error) {

	l := m.Logger("resolveEquipAction")

	var equippedObjectInstanceID string

	if sentence != "" {
		// Find object in room
		dungeonObjectViewRec, err := m.getObjectFromSentence(sentence, locationRecordSet.ObjectInstanceViewRecs)
		if err != nil {
			l.Warn("failed to get location object from sentence >%v<", err)
			return nil, err
		}
		if dungeonObjectViewRec == nil {
			// Find object stashed on character
			objectInstanceViewRecs, err := m.GetCharacterInstanceStashedObjectInstanceViewRecs(characterInstanceViewRec.ID)
			if err != nil {
				l.Warn("failed to get character stashed objects >%v<", err)
				return nil, err
			}
			dungeonObjectViewRec, err = m.getObjectFromSentence(sentence, objectInstanceViewRecs)
			if err != nil {
				l.Warn("failed to get character object from sentence >%v<", err)
				return nil, err
			}
		}
		if dungeonObjectViewRec != nil {
			equippedObjectInstanceID = dungeonObjectViewRec.ID
		}
	}

	dungeonActionRec := record.Action{
		DungeonInstanceID:                characterInstanceViewRec.DungeonInstanceID,
		LocationInstanceID:               characterInstanceViewRec.LocationInstanceID,
		CharacterInstanceID:              nullstring.FromString(characterInstanceViewRec.ID),
		ResolvedCommand:                  "equip",
		ResolvedTargetObjectInstanceID:   nullstring.FromString(equippedObjectInstanceID),
		ResolvedEquippedObjectInstanceID: nullstring.FromString(equippedObjectInstanceID),
	}

	return &dungeonActionRec, nil
}

func (m *Model) resolveDropAction(sentence string, characterInstanceViewRec *record.CharacterInstanceView, locationRecordSet *record.LocationInstanceViewRecordSet) (*record.Action, error) {

	l := m.Logger("resolveDropAction")

	var droppedObjectInstanceID string

	if sentence != "" {
		// Find object stashed on character
		l.Debug("Finding object stashed on character")
		objectInstanceViewRecs, err := m.GetCharacterInstanceStashedObjectInstanceViewRecs(characterInstanceViewRec.ID)
		if err != nil {
			l.Warn("failed to get character stashed objects >%v<", err)
			return nil, err
		}
		dungeonObjectRec, err := m.getObjectFromSentence(sentence, objectInstanceViewRecs)
		if err != nil {
			l.Warn("failed to get character object from sentence >%v<", err)
			return nil, err
		}
		l.Debug("Found object >%v< stashed on character", dungeonObjectRec)
		if dungeonObjectRec == nil {
			// Find object equipped on character
			l.Debug("Finding object equipped on character")
			objectInstanceViewRecs, err := m.GetCharacterInstanceEquippedObjectInstanceViewRecs(characterInstanceViewRec.ID)
			if err != nil {
				l.Warn("failed to get character equipped objects >%v<", err)
				return nil, err
			}
			dungeonObjectRec, err = m.getObjectFromSentence(sentence, objectInstanceViewRecs)
			if err != nil {
				l.Warn("failed to get character object from sentence >%v<", err)
				return nil, err
			}
			l.Debug("Found object >%v< equipped on character", dungeonObjectRec)
		}
		if dungeonObjectRec != nil {
			droppedObjectInstanceID = dungeonObjectRec.ID
		}
	}

	dungeonActionRec := record.Action{
		DungeonInstanceID:               characterInstanceViewRec.DungeonInstanceID,
		LocationInstanceID:              characterInstanceViewRec.LocationInstanceID,
		CharacterInstanceID:             nullstring.FromString(characterInstanceViewRec.ID),
		ResolvedCommand:                 "drop",
		ResolvedTargetObjectInstanceID:  nullstring.FromString(droppedObjectInstanceID),
		ResolvedDroppedObjectInstanceID: nullstring.FromString(droppedObjectInstanceID),
	}

	return &dungeonActionRec, nil
}

func (m *Model) resolveSentenceLocationDirection(sentence string, locationInstanceRec *record.LocationInstanceView) (string, string, error) {

	var dungeonLocationInstanceID string
	var dungeonLocationDirection string

	if locationInstanceRec.NortheastLocationInstanceID.Valid && strings.Contains(sentence, "northeast") {
		dungeonLocationInstanceID = locationInstanceRec.NortheastLocationInstanceID.String
		dungeonLocationDirection = "northeast"
	} else if locationInstanceRec.NorthwestLocationInstanceID.Valid && strings.Contains(sentence, "northwest") {
		dungeonLocationInstanceID = locationInstanceRec.NorthwestLocationInstanceID.String
		dungeonLocationDirection = "northwest"
	} else if locationInstanceRec.SoutheastLocationInstanceID.Valid && strings.Contains(sentence, "southeast") {
		dungeonLocationInstanceID = locationInstanceRec.SoutheastLocationInstanceID.String
		dungeonLocationDirection = "southeast"
	} else if locationInstanceRec.SoutheastLocationInstanceID.Valid && strings.Contains(sentence, "southeast") {
		dungeonLocationInstanceID = locationInstanceRec.SoutheastLocationInstanceID.String
		dungeonLocationDirection = "southeast"
	} else if locationInstanceRec.NorthLocationInstanceID.Valid && strings.Contains(sentence, "north") {
		dungeonLocationInstanceID = locationInstanceRec.NorthLocationInstanceID.String
		dungeonLocationDirection = "north"
	} else if locationInstanceRec.EastLocationInstanceID.Valid && strings.Contains(sentence, "east") {
		dungeonLocationInstanceID = locationInstanceRec.EastLocationInstanceID.String
		dungeonLocationDirection = "east"
	} else if locationInstanceRec.SouthLocationInstanceID.Valid && strings.Contains(sentence, "south") {
		dungeonLocationInstanceID = locationInstanceRec.SouthLocationInstanceID.String
		dungeonLocationDirection = "south"
	} else if locationInstanceRec.WestLocationInstanceID.Valid && strings.Contains(sentence, "west") {
		dungeonLocationInstanceID = locationInstanceRec.WestLocationInstanceID.String
		dungeonLocationDirection = "west"
	} else if locationInstanceRec.UpLocationInstanceID.Valid && strings.Contains(sentence, "up") {
		dungeonLocationInstanceID = locationInstanceRec.UpLocationInstanceID.String
		dungeonLocationDirection = "up"
	} else if locationInstanceRec.DownLocationInstanceID.Valid && strings.Contains(sentence, "down") {
		dungeonLocationInstanceID = locationInstanceRec.DownLocationInstanceID.String
		dungeonLocationDirection = "down"
	}

	return dungeonLocationInstanceID, dungeonLocationDirection, nil
}

func (m *Model) getObjectFromSentence(sentence string, dungeonObjectViewRecs []*record.ObjectInstanceView) (*record.ObjectInstanceView, error) {
	for _, dungeonObjectViewRec := range dungeonObjectViewRecs {
		if strings.Contains(sentence, strings.ToLower(dungeonObjectViewRec.Name)) {
			return dungeonObjectViewRec, nil
		}
	}
	return nil, nil
}

func (m *Model) resolveSentenceMonster(sentence string, dungeonMonsterViewRecs []*record.MonsterInstanceView) (*record.MonsterInstanceView, error) {
	for _, dungeonMonsterViewRec := range dungeonMonsterViewRecs {
		if strings.Contains(sentence, strings.ToLower(dungeonMonsterViewRec.Name)) {
			return dungeonMonsterViewRec, nil
		}
	}
	return nil, nil
}

func (m *Model) resolveSentenceCharacter(sentence string, characterInstanceViewRecs []*record.CharacterInstanceView) (*record.CharacterInstanceView, error) {
	for _, characterInstanceViewRec := range characterInstanceViewRecs {
		if strings.Contains(sentence, strings.ToLower(characterInstanceViewRec.Name)) {
			return characterInstanceViewRec, nil
		}
	}
	return nil, nil
}
