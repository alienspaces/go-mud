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

func (m *Model) resolveAction(sentence string, records *DungeonLocationRecordSet) (*record.DungeonAction, error) {

	resolved, err := m.resolveCommand(sentence)
	if err != nil {
		m.Log.Warn("Failed resolving command >%v<", err)
		return nil, err
	}

	resolveFuncs := map[string]func(sentence string, records *DungeonLocationRecordSet) (*record.DungeonAction, error){
		"move": m.resolveMoveAction,
		"look": m.resolveLookAction,
		// "equip": m.resolveEquipAction,
		// "stash": m.resolveStashAction,
		// "drop":  m.resolveDropAction,
	}

	resolveFunc, ok := resolveFuncs[resolved.Command]
	if !ok {
		msg := fmt.Sprintf("Command >%s< could not be resolved", resolved.Command)
		m.Log.Warn(msg)
		return nil, fmt.Errorf(msg)
	}

	dungeonActionRec, err := resolveFunc(resolved.Sentence, records)
	if err != nil {
		m.Log.Warn("Failed resolver function for command >%s< >%v<", resolved.Command, err)
		return nil, err
	}

	m.Log.Info("Resolved dungeon action rec >%#v<", dungeonActionRec)

	return dungeonActionRec, nil
}

func (m *Model) resolveCommand(sentence string) (*ResolverSentence, error) {
	sentenceWords := strings.Split(sentence, " ")
	resolved := ResolverSentence{}

	m.Log.Info("Have sentence words >%v<", sentenceWords)

	for _, dungeonAction := range []string{"move", "look", "equip", "stash", "drop"} {
		m.Log.Info("Checking dungeon action >%s<", dungeonAction)

		// NOTE: The appended space is important
		if strings.Contains(sentence, dungeonAction+" ") {
			m.Log.Info("Sentence contains action >%s<", dungeonAction)
			sentence = strings.Replace(sentence, dungeonAction+" ", "", 1)
			resolved.Command = dungeonAction
			resolved.Sentence = sentence
		} else if sentence == dungeonAction {
			m.Log.Info("Sentence equals action >%s<", dungeonAction)
			sentence = strings.Replace(sentence, dungeonAction, "", 1)
			resolved.Command = dungeonAction
			resolved.Sentence = sentence
		}
	}

	m.Log.Info("Resolved command >%#v<", resolved)

	return &resolved, nil
}

func (m *Model) resolveMoveAction(sentence string, records *DungeonLocationRecordSet) (*record.DungeonAction, error) {

	// Resolve move direction
	var targetDungeonLocationID string
	var targetDungeonLocationDirection string

	if sentence != "" {
		if records.LocationRec.NortheastDungeonLocationID.Valid && strings.Contains(sentence, "northeast") {
			targetDungeonLocationID = records.LocationRec.NortheastDungeonLocationID.String
			targetDungeonLocationDirection = "northeast"
		} else if records.LocationRec.NorthwestDungeonLocationID.Valid && strings.Contains(sentence, "northwest") {
			targetDungeonLocationID = records.LocationRec.NorthwestDungeonLocationID.String
			targetDungeonLocationDirection = "northwest"
		} else if records.LocationRec.SoutheastDungeonLocationID.Valid && strings.Contains(sentence, "southeast") {
			targetDungeonLocationID = records.LocationRec.SoutheastDungeonLocationID.String
			targetDungeonLocationDirection = "southeast"
		} else if records.LocationRec.SoutheastDungeonLocationID.Valid && strings.Contains(sentence, "southeast") {
			targetDungeonLocationID = records.LocationRec.SoutheastDungeonLocationID.String
			targetDungeonLocationDirection = "southeast"
		} else if records.LocationRec.NorthDungeonLocationID.Valid && strings.Contains(sentence, "north") {
			targetDungeonLocationID = records.LocationRec.NorthDungeonLocationID.String
			targetDungeonLocationDirection = "north"
		} else if records.LocationRec.EastDungeonLocationID.Valid && strings.Contains(sentence, "east") {
			targetDungeonLocationID = records.LocationRec.EastDungeonLocationID.String
			targetDungeonLocationDirection = "east"
		} else if records.LocationRec.SouthDungeonLocationID.Valid && strings.Contains(sentence, "south") {
			targetDungeonLocationID = records.LocationRec.SouthDungeonLocationID.String
			targetDungeonLocationDirection = "south"
		} else if records.LocationRec.WestDungeonLocationID.Valid && strings.Contains(sentence, "west") {
			targetDungeonLocationID = records.LocationRec.WestDungeonLocationID.String
			targetDungeonLocationDirection = "west"
		} else if records.LocationRec.UpDungeonLocationID.Valid && strings.Contains(sentence, "up") {
			targetDungeonLocationID = records.LocationRec.UpDungeonLocationID.String
			targetDungeonLocationDirection = "up"
		} else if records.LocationRec.DownDungeonLocationID.Valid && strings.Contains(sentence, "down") {
			targetDungeonLocationID = records.LocationRec.DownDungeonLocationID.String
			targetDungeonLocationDirection = "down"
		}
	}

	var targetDungeonLocationName string
	if targetDungeonLocationID != "" && len(records.LocationRecs) > 0 {
		for _, locationRec := range records.LocationRecs {
			if locationRec.ID == targetDungeonLocationID {
				targetDungeonLocationName = locationRec.Name
				break
			}
		}
	}

	if targetDungeonLocationID == "" {
		msg := fmt.Sprintf("failed to resolve target dungeon location id with sentence >%s<", sentence)
		m.Log.Warn(msg)
		return nil, fmt.Errorf(msg)
	}

	dungeonActionRecord := record.DungeonAction{
		// TODO: Should DungeonID, DungeonLocationID, DungeonCharacterID be set to the target location
		// when the character is moving, or is another dungeon action record created when the character
		// is actually moved in the performActionXxxx function, or is this in fact not correct and
		// would result in weird replay behaviour?
		DungeonID:                              records.CharacterRec.DungeonID,
		DungeonLocationID:                      records.CharacterRec.DungeonLocationID,
		DungeonCharacterID:                     store.NullString(records.CharacterRec.ID),
		ResolvedCommand:                        "move",
		ResolvedTargetDungeonLocationDirection: store.NullString(targetDungeonLocationDirection),
		ResolvedTargetDungeonLocationName:      store.NullString(targetDungeonLocationName),
		ResolvedTargetDungeonLocationID:        store.NullString(targetDungeonLocationID),
	}

	return &dungeonActionRecord, nil
}

func (m *Model) resolveLookAction(sentence string, records *DungeonLocationRecordSet) (*record.DungeonAction, error) {

	// Resolve look direction
	var targetDungeonLocationID string
	var targetDungeonLocationDirection string

	if sentence != "" {
		if records.LocationRec.NortheastDungeonLocationID.Valid && strings.Contains(sentence, "northeast") {
			targetDungeonLocationID = records.LocationRec.NortheastDungeonLocationID.String
			targetDungeonLocationDirection = "northeast"
		} else if records.LocationRec.NorthwestDungeonLocationID.Valid && strings.Contains(sentence, "northwest") {
			targetDungeonLocationID = records.LocationRec.NorthwestDungeonLocationID.String
			targetDungeonLocationDirection = "northwest"
		} else if records.LocationRec.SoutheastDungeonLocationID.Valid && strings.Contains(sentence, "southeast") {
			targetDungeonLocationID = records.LocationRec.SoutheastDungeonLocationID.String
			targetDungeonLocationDirection = "southeast"
		} else if records.LocationRec.SoutheastDungeonLocationID.Valid && strings.Contains(sentence, "southeast") {
			targetDungeonLocationID = records.LocationRec.SoutheastDungeonLocationID.String
			targetDungeonLocationDirection = "southeast"
		} else if records.LocationRec.NorthDungeonLocationID.Valid && strings.Contains(sentence, "north") {
			targetDungeonLocationID = records.LocationRec.NorthDungeonLocationID.String
			targetDungeonLocationDirection = "north"
		} else if records.LocationRec.EastDungeonLocationID.Valid && strings.Contains(sentence, "east") {
			targetDungeonLocationID = records.LocationRec.EastDungeonLocationID.String
			targetDungeonLocationDirection = "east"
		} else if records.LocationRec.SouthDungeonLocationID.Valid && strings.Contains(sentence, "south") {
			targetDungeonLocationID = records.LocationRec.SouthDungeonLocationID.String
			targetDungeonLocationDirection = "south"
		} else if records.LocationRec.WestDungeonLocationID.Valid && strings.Contains(sentence, "west") {
			targetDungeonLocationID = records.LocationRec.WestDungeonLocationID.String
			targetDungeonLocationDirection = "west"
		} else if records.LocationRec.UpDungeonLocationID.Valid && strings.Contains(sentence, "up") {
			targetDungeonLocationID = records.LocationRec.UpDungeonLocationID.String
			targetDungeonLocationDirection = "up"
		} else if records.LocationRec.DownDungeonLocationID.Valid && strings.Contains(sentence, "down") {
			targetDungeonLocationID = records.LocationRec.DownDungeonLocationID.String
			targetDungeonLocationDirection = "down"
		}
	}

	var targetDungeonLocationName string
	if targetDungeonLocationID != "" && len(records.LocationRecs) > 0 {
		for _, locationRec := range records.LocationRecs {
			if locationRec.ID == targetDungeonLocationID {
				targetDungeonLocationName = locationRec.Name
				break
			}
		}
	}

	// TODO: When a location was not identified, try to identify an object, a monster, a character.

	// TODO: When nothing has been identified, assume we are just looking in the current room.
	if targetDungeonLocationID == "" {
		targetDungeonLocationID = records.LocationRec.ID
		targetDungeonLocationName = records.LocationRec.Name
		targetDungeonLocationDirection = ""
	}

	dungeonActionRecord := record.DungeonAction{
		DungeonID:                              records.CharacterRec.DungeonID,
		DungeonLocationID:                      records.CharacterRec.DungeonLocationID,
		DungeonCharacterID:                     store.NullString(records.CharacterRec.ID),
		ResolvedCommand:                        "look",
		ResolvedTargetDungeonLocationDirection: store.NullString(targetDungeonLocationDirection),
		ResolvedTargetDungeonLocationName:      store.NullString(targetDungeonLocationName),
		ResolvedTargetDungeonLocationID:        store.NullString(targetDungeonLocationID),
	}

	return &dungeonActionRecord, nil
}

// func (m *Model) resolveEquipAction(sentence: string, records: DungeonLocationRecordSet): DungeonActionRepositoryRecord {
// 	let dungeonActionRecord: Partial<DungeonActionRepositoryRecord> = {
// 		dungeon_id: records.character.dungeon_id,
// 		dungeon_location_id: records.character.dungeon_location_id,
// 		dungeon_character_id: records.character.id,
// 	};
// 	return null;
// }

// func (m *Model) resolveStashAction(sentence: string, records: DungeonLocationRecordSet): DungeonActionRepositoryRecord {
// 	let dungeonActionRecord: Partial<DungeonActionRepositoryRecord> = {
// 		dungeon_id: records.character.dungeon_id,
// 		dungeon_location_id: records.character.dungeon_location_id,
// 		dungeon_character_id: records.character.id,
// 	};
// 	return null;
// }

// func (m *Model) resolveDropAction(sentence: string, records: DungeonLocationRecordSet): DungeonActionRepositoryRecord {
// 	let dungeonActionRecord: Partial<DungeonActionRepositoryRecord> = {
// 		dungeon_id: records.character.dungeon_id,
// 		dungeon_location_id: records.character.dungeon_location_id,
// 		dungeon_character_id: records.character.id,
// 	};
// 	return null;
// }
