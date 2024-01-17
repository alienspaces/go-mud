package model

import (
	"fmt"
	"strings"

	"gitlab.com/alienspaces/go-mud/backend/core/null"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/record"
)

type ResolveCommandArgs struct {
	Sentence                  string
	EntityType                EntityType
	EntityInstanceID          string
	LocationInstanceRecordSet *record.LocationInstanceViewRecordSet
}

type ResolvedCommand struct {
	Command  string
	Sentence string
}

func (m *Model) resolveCommand(args *ResolveCommandArgs) (*ResolvedCommand, error) {
	l := m.loggerWithFunctionContext("resolveCommand")

	if args == nil {
		return nil, NewInternalError("missing resolve command arguments")
	}

	// You cannot do anything when you are dead
	if args.EntityType == EntityTypeCharacter {
		for idx := range args.LocationInstanceRecordSet.CharacterInstanceViewRecs {
			if args.LocationInstanceRecordSet.CharacterInstanceViewRecs[idx].ID == args.EntityInstanceID &&
				args.LocationInstanceRecordSet.CharacterInstanceViewRecs[idx].CurrentHealth <= 0 {
				err := NewInvalidActionError("character name >%s< has died", args.LocationInstanceRecordSet.CharacterInstanceViewRecs[idx].Name)
				l.Warn(err.Error())
				return nil, err
			}
		}
	} else if args.EntityType == EntityTypeMonster {
		for idx := range args.LocationInstanceRecordSet.MonsterInstanceViewRecs {
			if args.LocationInstanceRecordSet.MonsterInstanceViewRecs[idx].ID == args.EntityInstanceID &&
				args.LocationInstanceRecordSet.MonsterInstanceViewRecs[idx].CurrentHealth <= 0 {
				err := NewInvalidActionError("monster name >%s< has died", args.LocationInstanceRecordSet.MonsterInstanceViewRecs[idx].Name)
				l.Warn(err.Error())
				return nil, err
			}
		}
	}

	sentence := args.Sentence
	sentenceWords := strings.Split(sentence, " ")
	resolved := ResolvedCommand{}

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
		err := NewInvalidActionError("command empty or not recognised, could not resolve command from >%#v<", args)
		l.Warn(err.Error())
		return nil, err
	}

	return &resolved, nil
}

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

func (m *Model) resolveAction(resolved *ResolvedCommand, args *ResolveActionArgs) (*record.Action, error) {
	l := m.loggerWithFunctionContext("resolveAction")

	if args == nil {
		return nil, NewInternalError("missing resolve action arguments")
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

func (m *Model) resolveActionMove(sentence string, args *ResolveActionArgs) (*record.Action, error) {
	l := m.loggerWithFunctionContext("resolveActionMove")

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
		return nil, NewInvalidDirectionError("you cannot move that direction")
	}

	dungeonActionRec := record.Action{
		DungeonInstanceID:                locationInstanceRec.DungeonInstanceID,
		LocationInstanceID:               locationInstanceRec.ID,
		ResolvedCommand:                  "move",
		ResolvedTargetLocationDirection:  null.NullStringFromString(targetLocationDirection),
		ResolvedTargetLocationInstanceID: null.NullStringFromString(targetLocationInstanceID),
	}

	if args.EntityType == EntityTypeCharacter {
		dungeonActionRec.CharacterInstanceID = null.NullStringFromString(args.EntityInstanceID)
	} else if args.EntityType == EntityTypeMonster {
		dungeonActionRec.MonsterInstanceID = null.NullStringFromString(args.EntityInstanceID)
	}

	return &dungeonActionRec, nil
}

// TODO: (game) A fair amount of common code between action use and action look, consider building
// some shared functions for identifying objects, characters or monsters at a location.
func (m *Model) resolveActionLook(sentence string, args *ResolveActionArgs) (*record.Action, error) {
	l := m.loggerWithFunctionContext("resolveActionLook")

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
		ResolvedTargetObjectInstanceID:    null.NullStringFromString(targetObjectInstanceID),
		ResolvedTargetMonsterInstanceID:   null.NullStringFromString(targetMonsterInstanceID),
		ResolvedTargetCharacterInstanceID: null.NullStringFromString(targetCharacterInstanceID),
		ResolvedTargetLocationDirection:   null.NullStringFromString(targetLocationDirection),
		ResolvedTargetLocationInstanceID:  null.NullStringFromString(targetLocationInstanceID),
	}

	if args.EntityType == EntityTypeCharacter {
		dungeonActionRec.CharacterInstanceID = null.NullStringFromString(args.EntityInstanceID)
	} else if args.EntityType == EntityTypeMonster {
		dungeonActionRec.MonsterInstanceID = null.NullStringFromString(args.EntityInstanceID)
	}

	return &dungeonActionRec, nil
}

// TODO: (game) A fair amount of common code between action use and action look, consider building
// some shared functions for identifying objects, characters or monsters at a location.
func (m *Model) resolveActionUse(sentence string, args *ResolveActionArgs) (*record.Action, error) {
	l := m.loggerWithFunctionContext("resolveActionUse")

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
		return nil, NewInvalidTargetError("failed to get object from location or equipped objects, cannot resolve use action")
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
		ResolvedTargetObjectInstanceID:    null.NullStringFromString(targetObjectInstanceID),
		ResolvedTargetMonsterInstanceID:   null.NullStringFromString(targetMonsterInstanceID),
		ResolvedTargetCharacterInstanceID: null.NullStringFromString(targetCharacterInstanceID),
	}

	if args.EntityType == EntityTypeCharacter {
		dungeonActionRec.CharacterInstanceID = null.NullStringFromString(args.EntityInstanceID)
	} else if args.EntityType == EntityTypeMonster {
		dungeonActionRec.MonsterInstanceID = null.NullStringFromString(args.EntityInstanceID)
	}

	return &dungeonActionRec, nil
}

// TODO: 11-implement-safe-locations: Check whether the current location is a safe location
// or not and disallow the attack action when it is.
func (m *Model) resolveActionAttack(sentence string, args *ResolveActionArgs) (*record.Action, error) {
	l := m.loggerWithFunctionContext("resolveActionAttack")

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
		return nil, NewInvalidTargetError("failed to find target monster or character, cannot resolve attack action")
	}

	dungeonActionRec := record.Action{
		DungeonInstanceID:                 locationInstanceRec.DungeonInstanceID,
		LocationInstanceID:                locationInstanceRec.ID,
		ResolvedCommand:                   "attack",
		ResolvedTargetMonsterInstanceID:   null.NullStringFromString(targetMonsterInstanceID),
		ResolvedTargetCharacterInstanceID: null.NullStringFromString(targetCharacterInstanceID),
	}

	if args.EntityType == EntityTypeCharacter {
		dungeonActionRec.CharacterInstanceID = null.NullStringFromString(args.EntityInstanceID)

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
			dungeonActionRec.ResolvedEquippedObjectInstanceID = null.NullStringFromString(objectInstanceRec.ID)
		}

	} else if args.EntityType == EntityTypeMonster {
		dungeonActionRec.MonsterInstanceID = null.NullStringFromString(args.EntityInstanceID)

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
			dungeonActionRec.ResolvedEquippedObjectInstanceID = null.NullStringFromString(objectInstanceRec.ID)
		}
	}

	return &dungeonActionRec, nil
}

func (m *Model) resolveActionStash(sentence string, args *ResolveActionArgs) (*record.Action, error) {
	l := m.loggerWithFunctionContext("resolveActionStash")

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
		return nil, NewInvalidTargetError("failed to identify object to stash, cannot resolve stash action")
	}

	dungeonActionRec := record.Action{
		DungeonInstanceID:               locationInstanceRec.DungeonInstanceID,
		LocationInstanceID:              locationInstanceRec.ID,
		ResolvedCommand:                 "stash",
		ResolvedTargetObjectInstanceID:  null.NullStringFromString(stashedObjectInstanceID),
		ResolvedStashedObjectInstanceID: null.NullStringFromString(stashedObjectInstanceID),
	}

	if args.EntityType == EntityTypeCharacter {
		dungeonActionRec.CharacterInstanceID = null.NullStringFromString(args.EntityInstanceID)
	} else if args.EntityType == EntityTypeMonster {
		dungeonActionRec.MonsterInstanceID = null.NullStringFromString(args.EntityInstanceID)
	}

	return &dungeonActionRec, nil
}

func (m *Model) resolveActionEquip(sentence string, args *ResolveActionArgs) (*record.Action, error) {
	l := m.loggerWithFunctionContext("resolveActionEquip")

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
		return nil, NewInvalidTargetError("failed to identify object to equip, cannot resolve equip action")
	}

	dungeonActionRec := record.Action{
		DungeonInstanceID:                locationInstanceRec.DungeonInstanceID,
		LocationInstanceID:               locationInstanceRec.ID,
		ResolvedCommand:                  "equip",
		ResolvedTargetObjectInstanceID:   null.NullStringFromString(equippedObjectInstanceID),
		ResolvedEquippedObjectInstanceID: null.NullStringFromString(equippedObjectInstanceID),
	}

	if args.EntityType == EntityTypeCharacter {
		dungeonActionRec.CharacterInstanceID = null.NullStringFromString(args.EntityInstanceID)
	} else if args.EntityType == EntityTypeMonster {
		dungeonActionRec.MonsterInstanceID = null.NullStringFromString(args.EntityInstanceID)
	}

	return &dungeonActionRec, nil
}

func (m *Model) resolveActionDrop(sentence string, args *ResolveActionArgs) (*record.Action, error) {
	l := m.loggerWithFunctionContext("resolveActionDrop")

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
		return nil, NewInvalidTargetError("failed to identify object to drop, cannot resolve drop action")
	}

	dungeonActionRec := record.Action{
		DungeonInstanceID:               locationInstanceRec.DungeonInstanceID,
		LocationInstanceID:              locationInstanceRec.ID,
		ResolvedCommand:                 "drop",
		ResolvedTargetObjectInstanceID:  null.NullStringFromString(droppedObjectInstanceID),
		ResolvedDroppedObjectInstanceID: null.NullStringFromString(droppedObjectInstanceID),
	}

	if args.EntityType == EntityTypeCharacter {
		dungeonActionRec.CharacterInstanceID = null.NullStringFromString(args.EntityInstanceID)
	} else if args.EntityType == EntityTypeMonster {
		dungeonActionRec.MonsterInstanceID = null.NullStringFromString(args.EntityInstanceID)
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
	l := m.loggerWithFunctionContext("getObjectFromSentence")

	for _, objectInstanceViewRec := range objectInstanceViewRecs {
		l.Info("Sentence >%s< contains >%s<", sentence, strings.ToLower(objectInstanceViewRec.Name))
		if strings.Contains(strings.ToLower(sentence), strings.ToLower(objectInstanceViewRec.Name)) {
			return objectInstanceViewRec, nil
		}
	}
	return nil, nil
}

func (m *Model) resolveSentenceMonster(sentence string, monsterInstanceViewRecs []*record.MonsterInstanceView) (*record.MonsterInstanceView, error) {
	l := m.loggerWithFunctionContext("resolveSentenceMonster")

	for idx := range monsterInstanceViewRecs {
		mivr := monsterInstanceViewRecs[idx]
		l.Info("Sentence >%s< contains >%s<", sentence, strings.ToLower(mivr.Name))
		if strings.Contains(strings.ToLower(sentence), strings.ToLower(mivr.Name)) {
			return mivr, nil
		}
	}
	return nil, nil
}

func (m *Model) resolveSentenceCharacter(sentence string, characterInstanceViewRecs []*record.CharacterInstanceView) (*record.CharacterInstanceView, error) {
	l := m.loggerWithFunctionContext("resolveSentenceCharacter")

	for idx := range characterInstanceViewRecs {
		civr := characterInstanceViewRecs[idx]
		l.Info("Sentence >%s< contains >%s<", sentence, strings.ToLower(civr.Name))
		if strings.Contains(strings.ToLower(sentence), strings.ToLower(civr.Name)) {
			return civr, nil
		}
	}
	return nil, nil
}
