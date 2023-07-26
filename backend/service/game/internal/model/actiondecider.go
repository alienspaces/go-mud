package model

import (
	"fmt"
	"math/rand"

	"gitlab.com/alienspaces/go-mud/backend/core/null"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/record"
)

// Feature flags implemented

const (
	// TODO: 15-implement-monster-goals
	FeatureMonsterGoalsImplemented bool = false
)

type DeciderArgs struct {
	MonsterInstanceViewRec    *record.MonsterInstanceView
	CharacterInstanceViewRec  *record.CharacterInstanceView
	LocationInstanceRecordSet *record.LocationInstanceViewRecordSet
	// Memory action records are action records that the character or monster
	// has recently performed or action records that were performed where the
	// character or monster was the target.
	Memories []*Memory
}

func (m *Model) decideAction(args *DeciderArgs) (string, error) {
	l := m.Logger("decideAction")

	l.Info("Deciding action args >%#v<", args)

	if args.MonsterInstanceViewRec != nil {
		l.Info("Deciding action for monster name >%s<", args.MonsterInstanceViewRec.Name)
	} else {
		l.Info("Deciding action for character name >%s<", args.CharacterInstanceViewRec.Name)
	}

	// Decider functions are typically prioritised as attack if anything is worth
	// attacking, then grab anything thats worth grabbing, look into other rooms
	// to find something interesting to move towards, and then move if there's
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

// getPriorityAttackTargetIndex takes a list of action records and a monster
// or character target instance ID and returns a list of monster or character
// instance IDs that have attacked the provided monster or character instance ID.
func (m *Model) getPriorityAttackTargetIndex(memories []*Memory, targetInstanceID string) (map[string]struct{}, error) {

	iidx := map[string]struct{}{}

	for idx := range memories {
		memory := memories[idx]
		if memory.ActionRec.ResolvedCommand == record.ActionCommandAttack &&
			(null.NullStringToString(memory.ActionRec.ResolvedTargetCharacterInstanceID) == targetInstanceID ||
				null.NullStringToString(memory.ActionRec.ResolvedTargetMonsterInstanceID) == targetInstanceID) {
			if null.NullStringIsValid(memory.ActionRec.CharacterInstanceID) {
				iidx[null.NullStringToString(memory.ActionRec.CharacterInstanceID)] = struct{}{}
				continue
			}
			if null.NullStringIsValid(memory.ActionRec.MonsterInstanceID) {
				iidx[null.NullStringToString(memory.ActionRec.MonsterInstanceID)] = struct{}{}
				continue
			}
		}
	}

	return iidx, nil
}

func (m *Model) decideActionAttack(args *DeciderArgs) (string, error) {
	l := m.Logger("decideActionAttack")

	lirs := args.LocationInstanceRecordSet

	targetName := ""
	if args.MonsterInstanceViewRec != nil {
		l.Info("Memories count >%d<", len(args.Memories))

		// Prioritise attacking anything that has attacked the monster
		pidx, err := m.getPriorityAttackTargetIndex(args.Memories, args.MonsterInstanceViewRec.ID)
		if err != nil {
			return "", err
		}

		l.Info("CharacterInstanceViewRecs count >%d<", len(lirs.CharacterInstanceViewRecs))
		l.Info("MonsterInstanceViewRecs count >%d<", len(lirs.MonsterInstanceViewRecs))
		l.Info("PriorityIndex count >%d<", len(pidx))

		// Highest priority are characters that have attacked the monster
		if len(pidx) != 0 && len(lirs.CharacterInstanceViewRecs) != 0 {
			for idx := range lirs.CharacterInstanceViewRecs {

				// No point attacking an already dead character
				if lirs.CharacterInstanceViewRecs[idx].CurrentHealth <= 0 {
					continue
				}

				if _, ok := pidx[lirs.CharacterInstanceViewRecs[idx].ID]; ok {
					targetName = lirs.CharacterInstanceViewRecs[idx].Name
					l.Warn("Choosing priority character instance >%s<", targetName)
					break
				}
			}
		}

		// Second highest priority are other monster that have attacked the monster
		if targetName == "" && len(pidx) != 0 && len(lirs.MonsterInstanceViewRecs) != 0 {
			for idx := range lirs.MonsterInstanceViewRecs {

				// No point attacking an already dead monster
				if lirs.MonsterInstanceViewRecs[idx].CurrentHealth <= 0 {
					continue
				}

				if _, ok := pidx[lirs.MonsterInstanceViewRecs[idx].ID]; ok {
					targetName = lirs.MonsterInstanceViewRecs[idx].Name
					l.Warn("Choosing priority monster instance >%s<", targetName)
					break
				}
			}
		}

		// Third priority is simply any characters present in the room
		if targetName == "" && len(lirs.CharacterInstanceViewRecs) != 0 {

			// No point randomly attacking a dead person!
			civRecs := []record.CharacterInstanceView{}
			for idx := range lirs.CharacterInstanceViewRecs {
				if lirs.CharacterInstanceViewRecs[idx].CurrentHealth > 0 {
					civRecs = append(civRecs, *lirs.CharacterInstanceViewRecs[idx])
				}
			}

			if len(civRecs) > 0 {
				rIdx := rand.Intn(len(civRecs))
				targetName = civRecs[rIdx].Name
				l.Warn("Choosing random character instance >%s<", targetName)
			}
		}
	}

	action := ""
	if targetName != "" {
		action = fmt.Sprintf("attack %s", targetName)
	}
	l.Info("Returning action >%s<", action)

	return action, nil
}

func (m *Model) decideActionStash(args *DeciderArgs) (string, error) {
	l := m.Logger("decideActionStash")

	lirs := args.LocationInstanceRecordSet

	action := ""
	if len(lirs.ObjectInstanceViewRecs) != 0 {
		// TODO: 16-implement-intelligent-stashing
		rIdx := rand.Intn(len(lirs.ObjectInstanceViewRecs))
		action = fmt.Sprintf("stash %s", lirs.ObjectInstanceViewRecs[rIdx].Name)
	}

	// TODO: 16-implement-intelligent-stashing
	l.Info("Might have returned action >%s<", action)

	return "", nil
}

// getPriorityLookDirection
func (m *Model) getPriorityLookDirection(memories []*Memory, lirs *record.LocationInstanceViewRecordSet) (string, error) {

	// Check if we've already seen anything interesting while collecting an
	// index of all directions we've recently looked.
	sawSomethingInteresting := false
	lookedIndex := map[string]struct{}{}

	for idx := range memories {
		memory := memories[idx]

		// The memory occurred at the current location, the memory was a look command and
		// and the look command was looking a direction
		if memory.ActionRec.LocationInstanceID == lirs.LocationInstanceViewRec.ID &&
			memory.ActionRec.ResolvedCommand == record.ActionCommandLook &&
			null.NullStringIsValid(memory.ActionRec.ResolvedTargetLocationDirection) {

			// Register that we looked this direction
			lookedIndex[null.NullStringToString(memory.ActionRec.ResolvedTargetLocationDirection)] = struct{}{}

			// TODO: 15-implement-monster-goals: The concept of finding something interesting
			// could be extended to anything that looks like it might help fulfil a goal. For
			// now we'll simply assume characters are interesting.

			// There were characters this direction so that will be interesting enough
			// to influence the movement phase.
			if len(memory.ActionCharacterRecs) > 0 {
				sawSomethingInteresting = true
			}
		}
	}

	if sawSomethingInteresting {
		return "", nil
	}

	// Return the first direction they haven't looked
	directions := lirs.LocationDirections()
	unlookedDirection := ""
	for _, direction := range directions {
		if _, ok := lookedIndex[direction]; ok {
			continue
		}
		unlookedDirection = direction
		break
	}

	return unlookedDirection, nil
}

// getPriorityLookObject chooses an object in the current room that has not been
// looked at which might influence a future decision.
func (m *Model) getPriorityLookObject(memories []*Memory, lirs *record.LocationInstanceViewRecordSet) (string, error) {

	objectName := ""
	lookedIndex := map[string]struct{}{}

	for idx := range memories {
		memory := memories[idx]

		// The memory was a look command and the look command was looking at an object
		if memory.ActionRec.ResolvedCommand == record.ActionCommandLook &&
			null.NullStringIsValid(memory.ActionRec.ResolvedTargetObjectInstanceID) {

			// Register that we looked this direction
			lookedIndex[null.NullStringToString(memory.ActionRec.ResolvedTargetObjectInstanceID)] = struct{}{}
		}
	}

	// Return the first object they haven't looked at
	for idx := range lirs.ObjectInstanceViewRecs {
		if _, ok := lookedIndex[lirs.ObjectInstanceViewRecs[idx].ID]; ok {
			continue
		}
		objectName = lirs.ObjectInstanceViewRecs[idx].Name
	}

	return objectName, nil
}

func (m *Model) getPriorityLookMonster(args *DeciderArgs) (string, error) {

	memories := args.Memories
	lirs := args.LocationInstanceRecordSet

	monsterName := ""
	lookedIndex := map[string]struct{}{}

	for idx := range memories {
		memory := memories[idx]

		// The memory was a look command and the look command was looking at a monster
		if memory.ActionRec.ResolvedCommand == record.ActionCommandLook &&
			null.NullStringIsValid(memory.ActionRec.ResolvedTargetMonsterInstanceID) {

			// Register that we looked this direction
			lookedIndex[null.NullStringToString(memory.ActionRec.ResolvedTargetMonsterInstanceID)] = struct{}{}
		}
	}

	// Return the first monster they haven't looked at that isn't themselves
	for idx := range lirs.MonsterInstanceViewRecs {
		if _, ok := lookedIndex[lirs.MonsterInstanceViewRecs[idx].ID]; ok {
			continue
		}
		if args.MonsterInstanceViewRec != nil && lirs.MonsterInstanceViewRecs[idx].ID == args.MonsterInstanceViewRec.ID {
			continue
		}
		monsterName = lirs.MonsterInstanceViewRecs[idx].Name
	}

	return monsterName, nil
}

func (m *Model) getPriorityLookCharacter(args *DeciderArgs) (string, error) {

	memories := args.Memories
	lirs := args.LocationInstanceRecordSet

	characterName := ""
	lookedIndex := map[string]struct{}{}

	for idx := range memories {
		memory := memories[idx]

		// The memory was a look command and the look command was looking at a character
		if memory.ActionRec.ResolvedCommand == record.ActionCommandLook &&
			null.NullStringIsValid(memory.ActionRec.ResolvedTargetCharacterInstanceID) {

			// Register that we looked this direction
			lookedIndex[null.NullStringToString(memory.ActionRec.ResolvedTargetCharacterInstanceID)] = struct{}{}
		}
	}

	// Return the first character they haven't looked at that isn't themselves
	for idx := range lirs.CharacterInstanceViewRecs {
		if _, ok := lookedIndex[lirs.CharacterInstanceViewRecs[idx].ID]; ok {
			continue
		}
		if args.CharacterInstanceViewRec != nil && lirs.CharacterInstanceViewRecs[idx].ID == args.CharacterInstanceViewRec.ID {
			continue
		}
		characterName = lirs.CharacterInstanceViewRecs[idx].Name
	}

	return characterName, nil
}

func (m *Model) decideActionLook(args *DeciderArgs) (string, error) {
	l := m.Logger("decideActionLook")

	lirs := args.LocationInstanceRecordSet

	action := ""

	// Prioritise looking for interesting things in other rooms
	direction, err := m.getPriorityLookDirection(args.Memories, lirs)
	if err != nil {
		l.Warn("failed getting priority look direction >%v<", err)
		return "", err
	}

	if direction != "" {
		action = fmt.Sprintf("look %s", direction)
		return action, nil
	}

	// TODO: 15-implement-monster-goals: Unless a monster has a goal
	// to achieve there is little point at looking at things to get
	// more information.
	if FeatureMonsterGoalsImplemented {

		monsterName, err := m.getPriorityLookMonster(args)
		if err != nil {
			l.Warn("failed getting priority look monster >%v<", err)
			return "", err
		}

		if monsterName != "" {
			action = fmt.Sprintf("look %s", monsterName)
			return action, nil
		}

		characterName, err := m.getPriorityLookCharacter(args)
		if err != nil {
			l.Warn("failed getting priority look character >%v<", err)
			return "", err
		}

		if characterName != "" {
			action = fmt.Sprintf("look %s", characterName)
			return action, nil
		}

		objectName, err := m.getPriorityLookObject(args.Memories, lirs)
		if err != nil {
			l.Warn("failed getting priority look objectName >%v<", err)
			return "", err
		}

		if objectName != "" {
			action = fmt.Sprintf("look %s", objectName)
			return action, nil
		}
	}

	return "", nil
}

func getPreviousLocationInstanceID(args *DeciderArgs) string {

	previousLocationInstanceID := ""
	for idx := range args.Memories {
		if args.MonsterInstanceViewRec != nil &&
			null.NullStringToString(args.Memories[idx].ActionRec.MonsterInstanceID) == args.MonsterInstanceViewRec.ID &&
			args.Memories[idx].ActionRec.LocationInstanceID != args.MonsterInstanceViewRec.LocationInstanceID {
			previousLocationInstanceID = args.Memories[idx].ActionRec.LocationInstanceID
			break
		}
		if args.CharacterInstanceViewRec != nil &&
			null.NullStringToString(args.Memories[idx].ActionRec.CharacterInstanceID) == args.CharacterInstanceViewRec.ID &&
			args.Memories[idx].ActionRec.LocationInstanceID != args.CharacterInstanceViewRec.LocationInstanceID {
			previousLocationInstanceID = args.Memories[idx].ActionRec.LocationInstanceID
			break
		}
	}

	return previousLocationInstanceID
}

func (m *Model) getPriorityMoveDirections(memories []*Memory, lirs *record.LocationInstanceViewRecordSet) ([]string, error) {

	// Check if we've already seen anything interesting while collecting an
	// index of all directions we've recently looked.
	directions := []string{}

	for idx := range memories {
		memory := memories[idx]

		// The memory occurred at the current location, the memory was a look command and
		// and the look command was looking a direction
		if memory.ActionRec.LocationInstanceID == lirs.LocationInstanceViewRec.ID &&
			memory.ActionRec.ResolvedCommand == record.ActionCommandLook &&
			null.NullStringIsValid(memory.ActionRec.ResolvedTargetLocationDirection) {

			// TODO: 15-implement-monster-goals: The concept of finding something interesting
			// could be extended to anything that looks like it might help fulfil a goal. For
			// now we'll simply assume characters are interesting.

			// There were characters in this direction so we'll add this direction to the list
			// of possible directions we could move
			if len(memory.ActionCharacterRecs) > 0 {
				directions = append(directions, null.NullStringToString((memory.ActionRec.ResolvedTargetLocationDirection)))
			}
		}
	}

	return directions, nil
}

func (m *Model) decideActionMove(args *DeciderArgs) (string, error) {
	l := m.Logger("decideActionAttack")

	action := ""
	lirs := args.LocationInstanceRecordSet

	// Move an interesting direction if possible
	interestingDirections, err := m.getPriorityMoveDirections(args.Memories, args.LocationInstanceRecordSet)
	if err != nil {
		l.Warn("failed to get priority move directions >%v<", err)
		return "", err
	}

	if len(interestingDirections) > 0 {
		rIdx := rand.Intn(len(interestingDirections))
		action = fmt.Sprintf("move %s", interestingDirections[rIdx])
	}

	// Otherwise move a direction we have not just come from
	if action == "" {
		directions := lirs.LocationDirections()

		// Always prefer to not move back the direction we just came from by excluding that
		// direction when there are multiple directions to choose from.
		recentLocationInstanceID := getPreviousLocationInstanceID(args)

		var availableDirections []string
		if len(directions) > 1 && recentLocationInstanceID != "" {
			for idx := range directions {
				if directions[idx] != recentLocationInstanceID {
					availableDirections = append(availableDirections, directions[idx])
				}
			}
		} else {
			availableDirections = directions
		}

		l.Info("Have possible directions >%v<", availableDirections)

		if len(availableDirections) != 0 {
			rIdx := rand.Intn(len(availableDirections))
			action = fmt.Sprintf("move %s", availableDirections[rIdx])
		}
	}

	l.Info("Returning action >%s<", action)

	return action, nil
}
