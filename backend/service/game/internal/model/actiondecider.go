package model

import (
	"database/sql"
	"fmt"
	"math/rand"

	"gitlab.com/alienspaces/go-mud/backend/core/nullstring"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/record"
)

// TODO: 14-implement-smarter-monsters - Pass memories into the decider so smarter decisions can be made

type DeciderArgs struct {
	MonsterInstanceViewRec    *record.MonsterInstanceView
	CharacterInstanceViewRec  *record.CharacterInstanceView
	LocationInstanceRecordSet *record.LocationInstanceViewRecordSet
	// Memory action records are action records that the character or monster
	// has recently performed or action records that were performed where the
	// character or monster was the target.
	MemoryActionRecs []*record.Action
}

func (m *Model) decideAction(args *DeciderArgs) (string, error) {
	l := m.Logger("decideAction")

	l.Info("Deciding action args >%#v<", args)

	if args.MonsterInstanceViewRec != nil {
		l.Info("Deciding action for monster name >%s<", args.MonsterInstanceViewRec.Name)
	} else {
		l.Info("Deciding action for character name >%s<", args.CharacterInstanceViewRec.Name)
	}

	// TODO: 14-implement-smarter-monsters - Revise the following comment

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
			// TODO: 14-implement-smarter-monsters
			// This is currently a random target however it would be cool if a
			// monster "stuck" to a target unless something even better came
			// along. Define "better"!
			rIdx := rand.Intn(len(lirs.CharacterInstanceViewRecs))
			action = fmt.Sprintf("attack %s", lirs.CharacterInstanceViewRecs[rIdx].Name)
		}
	}

	l.Info("Returning action >%s<", action)

	return action, nil
}

func (m *Model) decideActionStash(args *DeciderArgs) (string, error) {
	l := m.Logger("decideActionStash")

	lirs := args.LocationInstanceRecordSet

	action := ""
	if len(lirs.ObjectInstanceViewRecs) != 0 {
		// TODO: 14-implement-smarter-monsters
		// This is currently a random object however it would be cool if a
		// monster which was smart enough took the time to look through
		// object in the room and decide which was better than what it has!
		rIdx := rand.Intn(len(lirs.ObjectInstanceViewRecs))
		action = fmt.Sprintf("stash %s", lirs.ObjectInstanceViewRecs[rIdx].Name)
	}

	l.Info("Returning action >%s<", action)

	return action, nil
}

func (m *Model) decideActionLook(args *DeciderArgs) (string, error) {
	l := m.Logger("decideActionLook")

	lirs := args.LocationInstanceRecordSet

	action := ""

	things := []string{}
	if len(lirs.ObjectInstanceViewRecs) != 0 {
		rIdx := rand.Intn(len(lirs.ObjectInstanceViewRecs))
		things = append(things, lirs.ObjectInstanceViewRecs[rIdx].Name)
	}

	// Exclude the monster looking from the list of monsters to possibly look at
	var mivRecs []*record.MonsterInstanceView
	if args.MonsterInstanceViewRec != nil && len(lirs.MonsterInstanceViewRecs) != 0 {
		for idx := range lirs.MonsterInstanceViewRecs {
			if lirs.MonsterInstanceViewRecs[idx].ID == args.MonsterInstanceViewRec.ID {
				continue
			}
			mivRecs = append(mivRecs, lirs.MonsterInstanceViewRecs[idx])
		}
	}

	if len(mivRecs) != 0 {
		rIdx := rand.Intn(len(mivRecs))
		things = append(things, mivRecs[rIdx].Name)
	}

	// Exclude the character looking from the list of characters to possibly look at
	var civRecs []*record.CharacterInstanceView
	if args.CharacterInstanceViewRec != nil && len(lirs.CharacterInstanceViewRecs) != 0 {
		for idx := range lirs.CharacterInstanceViewRecs {
			if lirs.CharacterInstanceViewRecs[idx].ID == args.CharacterInstanceViewRec.ID {
				continue
			}
			civRecs = append(civRecs, lirs.CharacterInstanceViewRecs[idx])
		}
	}

	if len(civRecs) != 0 {
		rIdx := rand.Intn(len(civRecs))
		things = append(things, civRecs[rIdx].Name)
	}

	if len(things) > 0 {
		rIdx := rand.Intn(len(things))
		action = fmt.Sprintf("look %s", things[rIdx])
	}

	l.Info("Returning action >%s<", action)

	return action, nil
}

func getPreviousLocationInstanceID(args *DeciderArgs) string {

	previousLocationInstanceID := ""
	for idx := range args.MemoryActionRecs {
		if args.MonsterInstanceViewRec != nil &&
			nullstring.ToString(args.MemoryActionRecs[idx].MonsterInstanceID) == args.MonsterInstanceViewRec.ID &&
			args.MemoryActionRecs[idx].LocationInstanceID != args.MonsterInstanceViewRec.LocationInstanceID {
			previousLocationInstanceID = args.MemoryActionRecs[idx].LocationInstanceID
			break
		}
		if args.CharacterInstanceViewRec != nil &&
			nullstring.ToString(args.MemoryActionRecs[idx].CharacterInstanceID) == args.CharacterInstanceViewRec.ID &&
			args.MemoryActionRecs[idx].LocationInstanceID != args.CharacterInstanceViewRec.LocationInstanceID {
			previousLocationInstanceID = args.MemoryActionRecs[idx].LocationInstanceID
			break
		}
	}

	return previousLocationInstanceID
}

func (m *Model) decideActionMove(args *DeciderArgs) (string, error) {
	l := m.Logger("decideActionAttack")

	lirs := args.LocationInstanceRecordSet

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

	var possibleDirections []string
	for direction := range directions {
		if nullstring.IsValid(directions[direction]) {
			possibleDirections = append(possibleDirections, direction)
		}
	}

	// Always prefer to not move back the direction we just came from by excluding that
	// direction when there are multiple directions to choose from.
	recentLocationInstanceID := getPreviousLocationInstanceID(args)

	var availableDirections []string
	if len(possibleDirections) > 1 && recentLocationInstanceID != "" {
		for idx := range possibleDirections {
			if possibleDirections[idx] != recentLocationInstanceID {
				availableDirections = append(availableDirections, possibleDirections[idx])
			}
		}
	} else {
		availableDirections = possibleDirections
	}

	l.Info("Have possible directions >%v<", possibleDirections)

	action := ""
	if len(availableDirections) != 0 {
		rIdx := rand.Intn(len(availableDirections))
		action = fmt.Sprintf("move %s", availableDirections[rIdx])
	}

	l.Info("Returning action >%s<", action)

	return action, nil
}
