package model

import (
	"fmt"
	"strings"

	"gitlab.com/alienspaces/go-mud/backend/core/nullstring"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/record"
)

var validActionCommands []string = []string{
	record.ActionCommandMove,
	record.ActionCommandLook,
	record.ActionCommandUse,
	record.ActionCommandStash,
	record.ActionCommandEquip,
	record.ActionCommandDrop,
	record.ActionCommandAttack,
}

type ResolveActionArgs struct {
	EntityType                EntityType
	EntityInstanceID          string
	LocationInstanceRecordSet *record.LocationInstanceViewRecordSet
}

type ResolverSentence struct {
	Command  string
	Sentence string
}

func (m *Model) resolveAction(sentence string, args *ResolveActionArgs) (*record.Action, error) {
	l := m.Logger("resolveAction")

	resolved, err := m.resolveCommand(sentence)
	if err != nil {
		l.Warn("failed resolving command >%v<", err)
		return nil, err
	}

	resolveFuncs := map[string]func(sentence string, args *ResolveActionArgs) (*record.Action, error){
		record.ActionCommandMove:   m.resolveActionMove,
		record.ActionCommandLook:   m.resolveActionLook,
		record.ActionCommandUse:    m.resolveActionUse,
		record.ActionCommandStash:  m.resolveActionStash,
		record.ActionCommandEquip:  m.resolveActionEquip,
		record.ActionCommandDrop:   m.resolveActionDrop,
		record.ActionCommandAttack: m.resolveActionAttack,
	}

	resolveFunc, ok := resolveFuncs[resolved.Command]
	if !ok {
		msg := fmt.Sprintf("command >%s< could not be resolved", resolved.Command)
		l.Warn(msg)
		return nil, fmt.Errorf(msg)
	}

	actionRec, err := resolveFunc(resolved.Sentence, args)
	if err != nil {
		l.Warn("failed resolver function for command >%s< >%v<", resolved.Command, err)
		return nil, err
	}

	l.Debug("Resolved dungeon action record >%#v<", actionRec)

	return actionRec, nil
}

// IDEA: Some commands require an additional argument, we can probably short
// circuit that check here...
func (m *Model) resolveCommand(sentence string) (*ResolverSentence, error) {
	l := m.Logger("resolveCommand")

	sentenceWords := strings.Split(sentence, " ")
	resolved := ResolverSentence{}

	l.Debug("Have sentence words >%v<", sentenceWords)

	for _, actionCommand := range validActionCommands {
		l.Debug("Checking dungeon action >%s<", actionCommand)
		// NOTE: The appended space is important
		if strings.Contains(sentence, actionCommand+" ") {
			l.Debug("Sentence contains action >%s<", actionCommand)
			sentence = strings.Replace(sentence, actionCommand+" ", "", 1)
			resolved.Command = actionCommand
			resolved.Sentence = sentence
		} else if sentence == actionCommand {
			l.Debug("Sentence equals action >%s<", actionCommand)
			sentence = strings.Replace(sentence, actionCommand, "", 1)
			resolved.Command = actionCommand
			resolved.Sentence = sentence
		}
	}

	l.Debug("Resolved command >%#v<", resolved)

	if resolved.Command == "" {
		return nil, NewActionInvalidError("command empty or not recognised")
	}

	return &resolved, nil
}

func (m *Model) resolveActionMove(sentence string, args *ResolveActionArgs) (*record.Action, error) {
	l := m.Logger("resolveActionMove")

	var err error
	var targetLocationInstanceID string
	var targetLocationDirection string

	locationRecordSet := args.LocationInstanceRecordSet
	locationInstanceRec := locationRecordSet.LocationInstanceViewRec

	if sentence != "" {
		targetLocationInstanceID, targetLocationDirection, err = m.resolveSentenceLocationDirection(sentence, locationInstanceRec)
		if err != nil {
			l.Warn("failed to resolve sentence location direction >%v<", err)
			return nil, err
		}
	}

	if targetLocationInstanceID == "" || targetLocationDirection == "" {
		return nil, NewActionInvalidDirectionError("you cannot move that direction")
	}

	dungeonActionRec := record.Action{
		DungeonInstanceID:                locationInstanceRec.DungeonInstanceID,
		LocationInstanceID:               locationInstanceRec.ID,
		ResolvedCommand:                  "move",
		ResolvedTargetLocationDirection:  nullstring.FromString(targetLocationDirection),
		ResolvedTargetLocationInstanceID: nullstring.FromString(targetLocationInstanceID),
	}

	if args.EntityType == EntityTypeCharacter {
		dungeonActionRec.CharacterInstanceID = nullstring.FromString(args.EntityInstanceID)
	} else if args.EntityType == EntityTypeMonster {
		dungeonActionRec.MonsterInstanceID = nullstring.FromString(args.EntityInstanceID)
	}

	return &dungeonActionRec, nil
}

// TODO: (game) A fair amount of common code between action use and action look, consider building
// some shared functions for identifying objects, characters or monsters at a location.
func (m *Model) resolveActionLook(sentence string, args *ResolveActionArgs) (*record.Action, error) {
	l := m.Logger("resolveActionLook")

	var err error
	var targetLocationInstanceID string
	var targetLocationDirection string
	var targetObjectInstanceID string
	var targetMonsterInstanceID string
	var targetCharacterInstanceID string

	locationRecordSet := args.LocationInstanceRecordSet
	locationInstanceRec := locationRecordSet.LocationInstanceViewRec

	if sentence != "" {
		targetLocationInstanceID, targetLocationDirection, err = m.resolveSentenceLocationDirection(sentence, locationRecordSet.LocationInstanceViewRec)
		if err != nil {
			l.Warn("failed to resolve sentence location direction >%v<", err)
			return nil, err
		}

		if targetLocationInstanceID == "" {
			l.Info("Looking for object >%s<", sentence)
			objectInstanceViewRec, err := m.getObjectFromSentence(sentence, locationRecordSet.ObjectInstanceViewRecs)
			if err != nil {
				l.Warn("failed to resolve sentence object >%v<", err)
				return nil, err
			}
			if objectInstanceViewRec != nil {
				targetObjectInstanceID = objectInstanceViewRec.ID
			} else {
				if args.EntityType == EntityTypeCharacter {
					objectInstanceViewRecs, err := m.GetCharacterInstanceEquippedObjectInstanceViewRecs(args.EntityInstanceID)
					if err != nil {
						l.Warn("failed to character equipped object instance records >%v<", err)
						return nil, err
					}

					objectInstanceViewRec, err = m.getObjectFromSentence(sentence, objectInstanceViewRecs)
					if err != nil {
						l.Warn("failed to get location object from sentence >%v<", err)
						return nil, err
					}
					if objectInstanceViewRec != nil {
						targetObjectInstanceID = objectInstanceViewRec.ID
					}
				} else if args.EntityType == EntityTypeMonster {
					objectInstanceViewRecs, err := m.GetMonsterInstanceEquippedObjectInstanceViewRecs(args.EntityInstanceID)
					if err != nil {
						l.Warn("failed to monster equipped object instance records >%v<", err)
						return nil, err
					}
					objectInstanceViewRec, err = m.getObjectFromSentence(sentence, objectInstanceViewRecs)
					if err != nil {
						l.Warn("failed to get location object from sentence >%v<", err)
						return nil, err
					}
					if objectInstanceViewRec != nil {
						targetObjectInstanceID = objectInstanceViewRec.ID
					}
				}
			}
		}

		if targetLocationInstanceID == "" && targetObjectInstanceID == "" && targetCharacterInstanceID == "" {
			l.Info("Looking for monster >%s<", sentence)
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
			l.Info("Looking for character >%s<", sentence)
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

	dungeonActionRec := record.Action{
		DungeonInstanceID:                 locationInstanceRec.DungeonInstanceID,
		LocationInstanceID:                locationInstanceRec.ID,
		ResolvedCommand:                   "look",
		ResolvedTargetObjectInstanceID:    nullstring.FromString(targetObjectInstanceID),
		ResolvedTargetMonsterInstanceID:   nullstring.FromString(targetMonsterInstanceID),
		ResolvedTargetCharacterInstanceID: nullstring.FromString(targetCharacterInstanceID),
		ResolvedTargetLocationDirection:   nullstring.FromString(targetLocationDirection),
		ResolvedTargetLocationInstanceID:  nullstring.FromString(targetLocationInstanceID),
	}

	if args.EntityType == EntityTypeCharacter {
		dungeonActionRec.CharacterInstanceID = nullstring.FromString(args.EntityInstanceID)
	} else if args.EntityType == EntityTypeMonster {
		dungeonActionRec.MonsterInstanceID = nullstring.FromString(args.EntityInstanceID)
	}

	return &dungeonActionRec, nil
}

// TODO: (game) A fair amount of common code between action use and action look, consider building
// some shared functions for identifying objects, characters or monsters at a location.
func (m *Model) resolveActionUse(sentence string, args *ResolveActionArgs) (*record.Action, error) {
	l := m.Logger("resolveActionUse")

	var err error
	var targetObjectInstanceID string
	var targetMonsterInstanceID string
	var targetCharacterInstanceID string

	locationRecordSet := args.LocationInstanceRecordSet
	locationInstanceRec := locationRecordSet.LocationInstanceViewRec

	if sentence == "" {
		err := fmt.Errorf("missing sentence, cannot resolve use action")
		l.Warn(err.Error())
		return nil, err
	}

	l.Info("Looking for object >%s<", sentence)

	objectInstanceViewRec, err := m.getObjectFromSentence(sentence, locationRecordSet.ObjectInstanceViewRecs)
	if err != nil {
		l.Warn("failed to get location object from sentence >%v<", err)
		return nil, err
	}
	if objectInstanceViewRec != nil {
		targetObjectInstanceID = objectInstanceViewRec.ID
	} else {
		if args.EntityType == EntityTypeCharacter {
			objectInstanceViewRecs, err := m.GetCharacterInstanceEquippedObjectInstanceViewRecs(args.EntityInstanceID)
			if err != nil {
				l.Warn("failed to character equipped object instance records >%v<", err)
				return nil, err
			}

			objectInstanceViewRec, err = m.getObjectFromSentence(sentence, objectInstanceViewRecs)
			if err != nil {
				l.Warn("failed to get location object from sentence >%v<", err)
				return nil, err
			}
			if objectInstanceViewRec != nil {
				targetObjectInstanceID = objectInstanceViewRec.ID
			}
		} else if args.EntityType == EntityTypeMonster {
			objectInstanceViewRecs, err := m.GetMonsterInstanceEquippedObjectInstanceViewRecs(args.EntityInstanceID)
			if err != nil {
				l.Warn("failed to monster equipped object instance records >%v<", err)
				return nil, err
			}
			objectInstanceViewRec, err = m.getObjectFromSentence(sentence, objectInstanceViewRecs)
			if err != nil {
				l.Warn("failed to get location object from sentence >%v<", err)
				return nil, err
			}
			if objectInstanceViewRec != nil {
				targetObjectInstanceID = objectInstanceViewRec.ID
			}
		}
	}

	// No location or equipped object found
	if targetObjectInstanceID == "" {
		return nil, NewActionInvalidTargetError("failed to get object from location or equipped objects, cannot resolve use action")
	}

	// Use object on ... a monster?
	l.Info("Looking for monster >%s<", sentence)
	dungeonMonsterRec, err := m.resolveSentenceMonster(sentence, locationRecordSet.MonsterInstanceViewRecs)
	if err != nil {
		l.Warn("failed to resolve sentence monster >%v<", err)
		return nil, err
	}
	if dungeonMonsterRec != nil {
		targetMonsterInstanceID = dungeonMonsterRec.ID
	}

	// Use object on ... a character perhaps?
	if targetMonsterInstanceID == "" {
		l.Info("Looking for character >%s<", sentence)
		characterInstanceViewRec, err := m.resolveSentenceCharacter(sentence, locationRecordSet.CharacterInstanceViewRecs)
		if err != nil {
			l.Warn("failed to resolve sentence character >%v<", err)
			return nil, err
		}
		if characterInstanceViewRec != nil {
			targetCharacterInstanceID = characterInstanceViewRec.ID
		}
	}

	dungeonActionRec := record.Action{
		DungeonInstanceID:                 locationInstanceRec.DungeonInstanceID,
		LocationInstanceID:                locationInstanceRec.ID,
		ResolvedCommand:                   "use",
		ResolvedTargetObjectInstanceID:    nullstring.FromString(targetObjectInstanceID),
		ResolvedTargetMonsterInstanceID:   nullstring.FromString(targetMonsterInstanceID),
		ResolvedTargetCharacterInstanceID: nullstring.FromString(targetCharacterInstanceID),
	}

	if args.EntityType == EntityTypeCharacter {
		dungeonActionRec.CharacterInstanceID = nullstring.FromString(args.EntityInstanceID)
	} else if args.EntityType == EntityTypeMonster {
		dungeonActionRec.MonsterInstanceID = nullstring.FromString(args.EntityInstanceID)
	}

	return &dungeonActionRec, nil
}

// TODO: 11-implement-safe-locations: Check whether the current location is a safe location
// or not and disallow the attack action when it is.
func (m *Model) resolveActionAttack(sentence string, args *ResolveActionArgs) (*record.Action, error) {
	l := m.Logger("resolveActionAttack")

	var targetMonsterInstanceID string
	var targetCharacterInstanceID string

	locationRecordSet := args.LocationInstanceRecordSet
	locationInstanceRec := locationRecordSet.LocationInstanceViewRec

	if sentence != "" {

		// Attacking a monster
		dungeonMonsterRec, err := m.resolveSentenceMonster(sentence, locationRecordSet.MonsterInstanceViewRecs)
		if err != nil {
			l.Warn("failed to resolve sentence monster >%v<", err)
			return nil, err
		}
		if dungeonMonsterRec != nil {
			targetMonsterInstanceID = dungeonMonsterRec.ID
		}

		// Attacking a character
		if targetMonsterInstanceID == "" {
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

	if targetMonsterInstanceID == "" && targetCharacterInstanceID == "" {
		return nil, NewActionInvalidTargetError("failed to find target monster or character, cannot resolve attack action")
	}

	dungeonActionRec := record.Action{
		DungeonInstanceID:                 locationInstanceRec.DungeonInstanceID,
		LocationInstanceID:                locationInstanceRec.ID,
		ResolvedCommand:                   "attack",
		ResolvedTargetMonsterInstanceID:   nullstring.FromString(targetMonsterInstanceID),
		ResolvedTargetCharacterInstanceID: nullstring.FromString(targetCharacterInstanceID),
	}

	if args.EntityType == EntityTypeCharacter {
		dungeonActionRec.CharacterInstanceID = nullstring.FromString(args.EntityInstanceID)

		// Attacking with a specific weapon
		objectInstanceViewRecs, err := m.GetCharacterInstanceEquippedObjectInstanceViewRecs(args.EntityInstanceID)
		if err != nil {
			l.Warn("failed to get character equipped objects >%v<", err)
			return nil, err
		}

		objectInstanceRec, err := m.getObjectFromSentence(sentence, objectInstanceViewRecs)
		if err != nil {
			l.Warn("failed to get character object from sentence >%v<", err)
			return nil, err
		}

		if objectInstanceRec != nil {
			dungeonActionRec.ResolvedEquippedObjectInstanceID = nullstring.FromString(objectInstanceRec.ID)
		}

	} else if args.EntityType == EntityTypeMonster {
		dungeonActionRec.MonsterInstanceID = nullstring.FromString(args.EntityInstanceID)

		// Attacking with a specific weapon
		objectInstanceViewRecs, err := m.GetMonsterInstanceEquippedObjectInstanceViewRecs(args.EntityInstanceID)
		if err != nil {
			l.Warn("failed to get monster equipped objects >%v<", err)
			return nil, err
		}

		objectInstanceRec, err := m.getObjectFromSentence(sentence, objectInstanceViewRecs)
		if err != nil {
			l.Warn("failed to get monster object from sentence >%v<", err)
			return nil, err
		}

		if objectInstanceRec != nil {
			dungeonActionRec.ResolvedEquippedObjectInstanceID = nullstring.FromString(objectInstanceRec.ID)
		}
	}

	return &dungeonActionRec, nil
}

func (m *Model) resolveActionStash(sentence string, args *ResolveActionArgs) (*record.Action, error) {
	l := m.Logger("resolveActionStash")

	var stashedObjectInstanceID string

	locationRecordSet := args.LocationInstanceRecordSet
	locationInstanceRec := locationRecordSet.LocationInstanceViewRec

	if sentence != "" {
		// Find object in room
		objectInstanceRec, err := m.getObjectFromSentence(sentence, locationRecordSet.ObjectInstanceViewRecs)
		if err != nil {
			l.Warn("failed to get location object from sentence >%v<", err)
			return nil, err
		}
		if objectInstanceRec == nil {

			var objectInstanceViewRecs []*record.ObjectInstanceView
			if args.EntityType == EntityTypeCharacter {
				// Find object equipped on character
				objectInstanceViewRecs, err = m.GetCharacterInstanceEquippedObjectInstanceViewRecs(args.EntityInstanceID)
				if err != nil {
					l.Warn("failed to get character equipped objects >%v<", err)
					return nil, err
				}
			} else if args.EntityType == EntityTypeMonster {
				// Find object equipped on monster
				objectInstanceViewRecs, err = m.GetMonsterInstanceEquippedObjectInstanceViewRecs(args.EntityInstanceID)
				if err != nil {
					l.Warn("failed to get character equipped objects >%v<", err)
					return nil, err
				}
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

	if stashedObjectInstanceID == "" {
		return nil, NewActionInvalidTargetError("failed to identify object to stash, cannot resolve stash action")
	}

	dungeonActionRec := record.Action{
		DungeonInstanceID:               locationInstanceRec.DungeonInstanceID,
		LocationInstanceID:              locationInstanceRec.ID,
		ResolvedCommand:                 "stash",
		ResolvedTargetObjectInstanceID:  nullstring.FromString(stashedObjectInstanceID),
		ResolvedStashedObjectInstanceID: nullstring.FromString(stashedObjectInstanceID),
	}

	if args.EntityType == EntityTypeCharacter {
		dungeonActionRec.CharacterInstanceID = nullstring.FromString(args.EntityInstanceID)
	} else if args.EntityType == EntityTypeMonster {
		dungeonActionRec.MonsterInstanceID = nullstring.FromString(args.EntityInstanceID)
	}

	return &dungeonActionRec, nil
}

func (m *Model) resolveActionEquip(sentence string, args *ResolveActionArgs) (*record.Action, error) {
	l := m.Logger("resolveActionEquip")

	var equippedObjectInstanceID string

	locationRecordSet := args.LocationInstanceRecordSet
	locationInstanceRec := locationRecordSet.LocationInstanceViewRec

	if sentence != "" {
		// Find object in room
		dungeonObjectViewRec, err := m.getObjectFromSentence(sentence, locationRecordSet.ObjectInstanceViewRecs)
		if err != nil {
			l.Warn("failed to get location object from sentence >%v<", err)
			return nil, err
		}
		if dungeonObjectViewRec == nil {

			var objectInstanceViewRecs []*record.ObjectInstanceView
			if args.EntityType == EntityTypeCharacter {
				// Find object equipped on character
				objectInstanceViewRecs, err = m.GetCharacterInstanceStashedObjectInstanceViewRecs(args.EntityInstanceID)
				if err != nil {
					l.Warn("failed to get character equipped objects >%v<", err)
					return nil, err
				}
			} else if args.EntityType == EntityTypeMonster {
				// Find object equipped on monster
				objectInstanceViewRecs, err = m.GetMonsterInstanceStashedObjectInstanceViewRecs(args.EntityInstanceID)
				if err != nil {
					l.Warn("failed to get character equipped objects >%v<", err)
					return nil, err
				}
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

	if equippedObjectInstanceID == "" {
		return nil, NewActionInvalidTargetError("failed to identify object to equip, cannot resolve equip action")
	}

	dungeonActionRec := record.Action{
		DungeonInstanceID:                locationInstanceRec.DungeonInstanceID,
		LocationInstanceID:               locationInstanceRec.ID,
		ResolvedCommand:                  "equip",
		ResolvedTargetObjectInstanceID:   nullstring.FromString(equippedObjectInstanceID),
		ResolvedEquippedObjectInstanceID: nullstring.FromString(equippedObjectInstanceID),
	}

	if args.EntityType == EntityTypeCharacter {
		dungeonActionRec.CharacterInstanceID = nullstring.FromString(args.EntityInstanceID)
	} else if args.EntityType == EntityTypeMonster {
		dungeonActionRec.MonsterInstanceID = nullstring.FromString(args.EntityInstanceID)
	}

	return &dungeonActionRec, nil
}

func (m *Model) resolveActionDrop(sentence string, args *ResolveActionArgs) (*record.Action, error) {
	l := m.Logger("resolveActionDrop")

	var droppedObjectInstanceID string

	locationRecordSet := args.LocationInstanceRecordSet
	locationInstanceRec := locationRecordSet.LocationInstanceViewRec

	if sentence != "" {

		var err error
		var objectInstanceViewRecs []*record.ObjectInstanceView
		if args.EntityType == EntityTypeCharacter {
			// Find object equipped on character
			objectInstanceViewRecs, err = m.GetCharacterInstanceStashedObjectInstanceViewRecs(args.EntityInstanceID)
			if err != nil {
				l.Warn("failed to get character equipped objects >%v<", err)
				return nil, err
			}
		} else if args.EntityType == EntityTypeMonster {
			// Find object equipped on monster
			objectInstanceViewRecs, err = m.GetMonsterInstanceStashedObjectInstanceViewRecs(args.EntityInstanceID)
			if err != nil {
				l.Warn("failed to get character equipped objects >%v<", err)
				return nil, err
			}
		}

		dungeonObjectRec, err := m.getObjectFromSentence(sentence, objectInstanceViewRecs)
		if err != nil {
			l.Warn("failed to get character object from sentence >%v<", err)
			return nil, err
		}

		l.Debug("Found object >%v< stashed on character", dungeonObjectRec)

		if dungeonObjectRec == nil {

			if args.EntityType == EntityTypeCharacter {
				// Find object equipped on character
				objectInstanceViewRecs, err = m.GetCharacterInstanceEquippedObjectInstanceViewRecs(args.EntityInstanceID)
				if err != nil {
					l.Warn("failed to get character equipped objects >%v<", err)
					return nil, err
				}
			} else if args.EntityType == EntityTypeMonster {
				// Find object equipped on monster
				objectInstanceViewRecs, err = m.GetMonsterInstanceEquippedObjectInstanceViewRecs(args.EntityInstanceID)
				if err != nil {
					l.Warn("failed to get character equipped objects >%v<", err)
					return nil, err
				}
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

	if droppedObjectInstanceID == "" {
		return nil, NewActionInvalidTargetError("failed to identify object to drop, cannot resolve drop action")
	}

	dungeonActionRec := record.Action{
		DungeonInstanceID:               locationInstanceRec.DungeonInstanceID,
		LocationInstanceID:              locationInstanceRec.ID,
		ResolvedCommand:                 "drop",
		ResolvedTargetObjectInstanceID:  nullstring.FromString(droppedObjectInstanceID),
		ResolvedDroppedObjectInstanceID: nullstring.FromString(droppedObjectInstanceID),
	}

	if args.EntityType == EntityTypeCharacter {
		dungeonActionRec.CharacterInstanceID = nullstring.FromString(args.EntityInstanceID)
	} else if args.EntityType == EntityTypeMonster {
		dungeonActionRec.MonsterInstanceID = nullstring.FromString(args.EntityInstanceID)
	}

	return &dungeonActionRec, nil
}

// TODO: (game) The following functions do not need to return an error
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

func (m *Model) getObjectFromSentence(sentence string, objectInstanceViewRecs []*record.ObjectInstanceView) (*record.ObjectInstanceView, error) {
	l := m.Logger("getObjectFromSentence")

	for _, objectInstanceViewRec := range objectInstanceViewRecs {
		l.Info("Sentence >%s< contains >%s<", sentence, strings.ToLower(objectInstanceViewRec.Name))
		if strings.Contains(strings.ToLower(sentence), strings.ToLower(objectInstanceViewRec.Name)) {
			return objectInstanceViewRec, nil
		}
	}
	return nil, nil
}

func (m *Model) resolveSentenceMonster(sentence string, monsterInstanceViewRecs []*record.MonsterInstanceView) (*record.MonsterInstanceView, error) {
	l := m.Logger("resolveSentenceMonster")

	for _, monsterInstanceViewRec := range monsterInstanceViewRecs {
		l.Info("Sentence >%s< contains >%s<", strings.ToLower(monsterInstanceViewRec.Name))
		if strings.Contains(strings.ToLower(sentence), strings.ToLower(monsterInstanceViewRec.Name)) {
			return monsterInstanceViewRec, nil
		}
	}
	return nil, nil
}

func (m *Model) resolveSentenceCharacter(sentence string, characterInstanceViewRecs []*record.CharacterInstanceView) (*record.CharacterInstanceView, error) {
	l := m.Logger("resolveSentenceCharacter")

	for _, characterInstanceViewRec := range characterInstanceViewRecs {
		l.Info("Sentence >%s< contains >%s<", sentence, strings.ToLower(characterInstanceViewRec.Name))
		if strings.Contains(strings.ToLower(sentence), strings.ToLower(characterInstanceViewRec.Name)) {
			return characterInstanceViewRec, nil
		}
	}
	return nil, nil
}
