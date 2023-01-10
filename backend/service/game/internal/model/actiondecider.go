package model

import (
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/record"
)

// NOTES: Monsters should behave based on their statistics and capabilities.
// Examples:
// - Intelligent humanoid monsters with no equipped weapon should attempt to equip any
//   weapons left lying in the current location.
//   -- Requires item types
// - Intelligent monsters that are not defending anything other than themselves should
//   attempt to run away when losing a fight
//   -- Compare own health with health of characters in the room that have previously
//      attacked it
// - Intelligent monsters that are obviously out of their depth in a battle should run
//   away.

type DeciderArgs struct {
	EntityType                EntityType
	EntityInstanceID          string
	LocationInstanceRecordSet *record.LocationInstanceViewRecordSet
}

func (m *Model) decideAction(args *DeciderArgs) (string, error) {
	l := m.Logger("decideAction")

	l.Info("deciding action args >%#v<", args)

	return "", nil
}
