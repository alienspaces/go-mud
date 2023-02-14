package model

import (
	"database/sql"
	"fmt"
	"math/rand"

	"gitlab.com/alienspaces/go-mud/backend/core/nullstring"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/record"
)

// TODO: 9-implement-monster-actions
// NOTES: Monsters should ultimately behave based on their statistics and capabilities.
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
	MonsterInstanceViewRec    *record.MonsterInstanceView
	CharacterInstanceViewRec  *record.CharacterInstanceView
	LocationInstanceRecordSet *record.LocationInstanceViewRecordSet
}

func (m *Model) decideAction(args *DeciderArgs) (string, error) {
	l := m.Logger("decideAction")

	l.Info("Deciding action args >%#v<", args)

	if args.MonsterInstanceViewRec != nil {
		l.Info("Deciding action for monster name >%s<", args.MonsterInstanceViewRec.Name)
	} else {
		l.Info("Deciding action for character name >%s<", args.CharacterInstanceViewRec.Name)
	}

	// Decider functions are typically prioritised as attack if anything is
	// worth attacking, then grab anything thats worth grabbing, look into
	// other rooms to find someone to attack, and then the move if there's
	// somewhere worth moving to.
	deciderFuncs := []func(args *DeciderArgs) (string, error){
		m.decideActionAttack,
		m.decideActionStash,
		m.decideActionLook,
		m.decideActionMove,
	}

	var err error
	var sentence string

	for idx := range deciderFuncs {
		sentence, err = deciderFuncs[idx](args)
		if err != nil {
			return "", err
		}
		if sentence != "" {
			break
		}
	}

	return sentence, nil
}

func (m *Model) decideActionAttack(args *DeciderArgs) (string, error) {
	l := m.Logger("decideActionAttack")

	lirs := args.LocationInstanceRecordSet

	action := ""
	if args.MonsterInstanceViewRec != nil {
		l.Info("Character instance count >%d<", len(lirs.CharacterInstanceViewRecs))
		if len(lirs.CharacterInstanceViewRecs) != 0 {
			// TODO: This is currently a random target however it would be cool
			// if a monster "stuck" to a target unless something even better came
			// along. Define "better"!
			rIdx := rand.Intn(len(lirs.CharacterInstanceViewRecs))
			action = fmt.Sprintf("attack %s", lirs.CharacterInstanceViewRecs[rIdx].Name)
		}
	}

	l.Info("Returning action >%s<", action)

	return action, nil
}

// TODO: 9-implement-monster-actions
func (m *Model) decideActionStash(args *DeciderArgs) (string, error) {
	return "", nil
}

// TODO: 9-implement-monster-actions
func (m *Model) decideActionLook(args *DeciderArgs) (string, error) {

	return "", nil
}

func (m *Model) decideActionMove(args *DeciderArgs) (string, error) {
	l := m.Logger("decideActionAttack")

	lirs := args.LocationInstanceRecordSet

	var possibleDirections []string
	directions := map[string]sql.NullString{
		"north":     lirs.LocationInstanceViewRec.NorthLocationInstanceID,
		"northeast": lirs.LocationInstanceViewRec.NortheastLocationInstanceID,
		"east":      lirs.LocationInstanceViewRec.EastLocationInstanceID,
		"southeast": lirs.LocationInstanceViewRec.SoutheastLocationInstanceID,
		"south":     lirs.LocationInstanceViewRec.SouthLocationInstanceID,
		"southwest": lirs.LocationInstanceViewRec.SouthwestLocationInstanceID,
		"west":      lirs.LocationInstanceViewRec.WestLocationInstanceID,
		"northwest": lirs.LocationInstanceViewRec.NorthwestLocationInstanceID,
		"up":        lirs.LocationInstanceViewRec.UpLocationInstanceID,
		"down":      lirs.LocationInstanceViewRec.DownLocationInstanceID,
	}

	for direction := range directions {
		if nullstring.IsValid(directions[direction]) {
			possibleDirections = append(possibleDirections, direction)
		}
	}

	l.Info("Have possible directions >%v<", possibleDirections)

	action := ""
	if len(possibleDirections) != 0 {
		// TODO: This is currently a random direction however it would be cool
		// if a monster was chasing a target for it to look in various directions
		// for that target before deciding which direction to move.
		rIdx := rand.Intn(len(possibleDirections))
		action = fmt.Sprintf("move %s", possibleDirections[rIdx])
	}

	l.Info("Returning action >%s<", action)

	return action, nil
}
