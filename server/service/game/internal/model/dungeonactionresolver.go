package model

import (
	"strings"

	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

type ResolverSentence struct {
	Command  string
	Sentence string
}

func (m *Model) resolveAction(sentence string, records DungeonLocationRecordSet) (*record.DungeonAction, error) {

	resolved, err := m.resolveCommand(sentence)
	if err != nil {
		m.Log.Warn("Failed resolving command >%v<", err)
		return nil, err
	}

	resolveFuncs := map[string]func(sentence string, records DungeonLocationRecordSet) (*record.DungeonAction, error){
		"move": m.resolveMoveAction,
		// "look":  m.resolveLookAction,
		// "equip": m.resolveEquipAction,
		// "stash": m.resolveStashAction,
		// "drop":  m.resolveDropAction,
	}

	dungeonActionRecord, err := resolveFuncs[resolved.Command](resolved.Sentence, records)
	if err != nil {
		m.Log.Warn("Failed resolver function for command >%s< >%v<", resolved.Command, err)
		return nil, err
	}

	return dungeonActionRecord, nil
}

func (m *Model) resolveCommand(sentence string) (*ResolverSentence, error) {
	parts := strings.Split(sentence, " ")
	resolved := ResolverSentence{}

	m.Log.Info("Have command parts >%v<", parts)

	for _, dungeonAction := range []string{"move", "look", "equip", "stash", "drop"} {
		m.Log.Info("Checking dungeon action >%s<", dungeonAction)
	}

	// 	const index = parts.indexOf(findAction);
	// 	if (index === -1) {
	// 		return;
	// 	}
	// 	resolved = {
	// 		command: findAction,
	// 		sentence: parts.length > index + 1 ? parts.splice(index + 1).join(' ') : undefined,
	// 	};
	// 	return true;
	// });

	return &resolved, nil
}

func (m *Model) resolveMoveAction(sentence string, records DungeonLocationRecordSet) (*record.DungeonAction, error) {
	var command string
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

	dungeonActionRecord := record.DungeonAction{
		DungeonID:                              records.CharacterRec.DungeonID,
		DungeonLocationID:                      records.CharacterRec.DungeonLocationID,
		DungeonCharacterID:                     records.CharacterRec.ID,
		ResolvedCommand:                        command,
		ResolvedTargetDungeonLocationDirection: targetDungeonLocationDirection,
		ResolvedTargetDungeonLocationName:      targetDungeonLocationName,
		ResolvedTargetDungeonLocationID:        targetDungeonLocationID,
	}

	return &dungeonActionRecord, nil
}

// func (m *Model) resolveLookAction(sentence: string, records: DungeonLocationRecordSet): DungeonActionRepositoryRecord {
// 	const command = 'look';
// 	let targetDungeonLocationId: string;
// 	let targetDungeonLocationDirection: string;

// 	// Looking in a direction where there is another location?
// 	if (sentence) {
// 		for (var prop in DIRECTION_MAP) {
// 			if (records.location[prop] && sentence.match(new RegExp(`\s?${DIRECTION_MAP[prop]}(?![A-Za-z]+)`))) {
// 				targetDungeonLocationId = records.location[prop];
// 				targetDungeonLocationDirection = DIRECTION_MAP[prop];
// 				break;
// 			}
// 		}
// 	}

// 	// When not a direction where there is a room exit the character
// 	// is looking at the current location.
// 	let targetDungeonLocationName: string;
// 	if (targetDungeonLocationId == null) {
// 		targetDungeonLocationId = records.location.id;
// 		targetDungeonLocationName = records.location.name;
// 	} else if (records.locations) {
// 		records.locations.some((location) => {
// 			if (location.id === targetDungeonLocationId) {
// 				targetDungeonLocationName = location.name;
// 				return true;
// 			}
// 		});
// 	}

// 	const dungeonActionRecord: DungeonActionRepositoryRecord = {
// 		dungeon_id: records.character.dungeon_id,
// 		dungeon_location_id: records.character.dungeon_location_id,
// 		dungeon_character_id: records.character.id,
// 		resolved_command: command,
// 		resolved_target_dungeon_location_direction: targetDungeonLocationDirection,
// 		resolved_target_dungeon_location_name: targetDungeonLocationName,
// 		resolved_target_dungeon_location_id: targetDungeonLocationId,
// 	};

// 	return dungeonActionRecord;
// }

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
