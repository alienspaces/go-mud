package model

import (
	"fmt"

	"gitlab.com/alienspaces/go-mud/backend/core/null"
	coresql "gitlab.com/alienspaces/go-mud/backend/core/sql"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/record"
)

type EntityType string

const EntityTypeMonster EntityType = "monster"
const EntityTypeCharacter EntityType = "character"

type Memory struct {
	ActionRec           *record.Action
	ActionCharacterRecs []*record.ActionCharacter
	ActionMonsterRecs   []*record.ActionMonster
}

// ProcessCharacterAction - Processes a submitted character action
func (m *Model) ProcessCharacterAction(dungeonInstanceID string, characterInstanceID string, sentence string) (*record.ActionRecordSet, error) {
	l := m.Logger("ProcessCharacterAction")

	l.Info("Processing character ID >%s< action command >%s<", characterInstanceID, sentence)

	// Verify the character performing the action exists within the specified dungeon
	civRec, err := m.GetCharacterInstanceViewRec(characterInstanceID)
	if err != nil {
		l.Warn("failed getting character record before performing action >%v<", err)
		return nil, err
	}
	if civRec == nil {
		msg := fmt.Sprintf("failed getting character record ID >%s< before performing action", characterInstanceID)
		l.Warn(msg)
		return nil, fmt.Errorf(msg)
	}

	if civRec.DungeonInstanceID != dungeonInstanceID {
		msg := fmt.Sprintf("character ID >%s< does not exist in dungeon ID >%s<", characterInstanceID, dungeonInstanceID)
		l.Warn(msg)
		return nil, fmt.Errorf(msg)
	}

	// TODO: 10-implement-effects:
	// Process any active effects that are still applied to the character.

	// TODO: 12-implement-death: Remove character instance when dead

	// Get the current dungeon location set of related records
	locationInstanceRecordSet, err := m.GetLocationInstanceViewRecordSet(civRec.LocationInstanceID, true)
	if err != nil {
		l.Warn("failed getting dungeon location record set before performing action >%v<", err)
		return nil, err
	}
	if locationInstanceRecordSet == nil {
		msg := fmt.Sprintf("failed getting dungeon location record ID >%s< set before performing action", civRec.LocationInstanceID)
		l.Warn(msg)
		return nil, fmt.Errorf(msg)
	}

	// Resolve the submitted character action
	actionRec, err := m.resolveAction(sentence, &ResolveActionArgs{
		EntityType:                EntityTypeCharacter,
		EntityInstanceID:          civRec.ID,
		LocationInstanceRecordSet: locationInstanceRecordSet,
	})
	if err != nil {
		l.Warn("failed resolving character action >%v<", err)
		return nil, err
	}

	// Resolve the initial action turn
	actionRec, err = m.resolveActionTurn(&ResolveActionTurnArgs{
		ActionRec:         actionRec,
		EntityType:        EntityTypeCharacter,
		EntityInstanceID:  civRec.ID,
		DungeonInstanceID: locationInstanceRecordSet.LocationInstanceViewRec.DungeonInstanceID,
	})
	if err != nil {
		l.Warn("failed resolving action turn >%v<", err)
		return nil, err
	}

	l.Info("Character ID >%s< Name >%s< Action record ID >%s< TurnNumber >%d<", civRec.CharacterID, civRec.Name, actionRec.ID, actionRec.TurnNumber)

	// Perform the submitted character action
	actionRec, err = m.performAction(
		&PerformActionArgs{
			ActionRec:                 actionRec,
			CharacterInstanceViewRec:  civRec,
			MonsterInstanceViewRec:    nil,
			LocationInstanceRecordSet: locationInstanceRecordSet,
		},
	)
	if err != nil {
		l.Warn("failed performing character action >%v<", err)
		return nil, err
	}

	// Create the resulting action event record
	err = m.CreateActionRec(actionRec)
	if err != nil {
		l.Warn("failed creating action record >%v<", err)
		return nil, err
	}

	l.Info("Created action record ID >%s< SerialNumber >%d<", actionRec.ID, null.NullInt16ToInt16(actionRec.SerialNumber))

	// TODO: (game) Maybe don't need to do this... Get the updated character record
	civRec, err = m.GetCharacterInstanceViewRec(characterInstanceID)
	if err != nil {
		l.Warn("failed getting character record after performing action >%v<", err)
		return nil, err
	}

	// Create action character record
	actionCharacterRec := record.ActionCharacter{
		RecordType:          record.ActionCharacterRecordTypeSource,
		ActionID:            actionRec.ID,
		LocationInstanceID:  actionRec.LocationInstanceID,
		CharacterInstanceID: civRec.ID,
		Name:                civRec.Name,
		Strength:            civRec.Strength,
		Dexterity:           civRec.Dexterity,
		Intelligence:        civRec.Intelligence,
		CurrentStrength:     civRec.CurrentStrength,
		CurrentDexterity:    civRec.CurrentDexterity,
		CurrentIntelligence: civRec.CurrentIntelligence,
		Health:              civRec.Health,
		Fatigue:             civRec.Fatigue,
		CurrentHealth:       civRec.CurrentHealth,
		CurrentFatigue:      civRec.CurrentFatigue,
	}

	// Create source action character record
	err = m.CreateActionCharacterRec(&actionCharacterRec)
	if err != nil {
		l.Warn("failed creating source action character record >%v<", err)
		return nil, err
	}

	// Create source action character object records
	oivRecs, err := m.GetCharacterInstanceObjectInstanceViewRecs(civRec.ID)
	if err != nil {
		l.Warn("failed getting source character object instance view records >%v<", err)
		return nil, err
	}

	actionCharacterObjectRecs := []*record.ActionCharacterObject{}
	for _, oivRec := range oivRecs {
		l.Info("Adding source character action object >%s<", oivRec.Name)
		dungeonCharacterObjectRec := record.ActionCharacterObject{
			ActionCharacterID: actionCharacterRec.ID,
			ObjectInstanceID:  oivRec.ID,
			Name:              oivRec.Name,
			IsEquipped:        oivRec.IsEquipped,
			IsStashed:         oivRec.IsStashed,
		}
		err := m.CreateActionCharacterObjectRec(&dungeonCharacterObjectRec)
		if err != nil {
			l.Warn("failed creating source action character object record >%v<", err)
			return nil, err
		}
		actionCharacterObjectRecs = append(actionCharacterObjectRecs, &dungeonCharacterObjectRec)
	}

	actionRecordSet, err := m.createActionRecordSetRecords(&record.ActionRecordSet{
		ActionRec:                 actionRec,
		ActionCharacterRec:        &actionCharacterRec,
		ActionCharacterObjectRecs: actionCharacterObjectRecs,
	})
	if err != nil {
		l.Warn("failed creating action record set records >%v<", err)
		return nil, err
	}

	return actionRecordSet, nil
}

type DecideMonsterActionResult struct {
	DungeonInstanceID string
	MonsterInstanceID string
	Sentence          string
}

// DecideMonsterAction -
func (m *Model) DecideMonsterAction(monsterInstanceID string) (*DecideMonsterActionResult, error) {
	l := m.Logger("DecideMonsterAction")

	l.Info("Deciding monster instance ID >%s< action", monsterInstanceID)

	rec, err := m.GetMonsterInstanceViewRec(monsterInstanceID)
	if err != nil {
		l.Warn("failed getting monster instance view record >%v<", err)
		return nil, err
	}

	// Get the current dungeon location set of related records
	locationInstanceRecordSet, err := m.GetLocationInstanceViewRecordSet(rec.LocationInstanceID, true)
	if err != nil {
		l.Warn("failed getting dungeon location record set before performing action >%v<", err)
		return nil, err
	}
	if locationInstanceRecordSet == nil {
		msg := fmt.Sprintf("failed getting dungeon location record ID >%s< set before performing action", rec.LocationInstanceID)
		l.Warn(msg)
		return nil, fmt.Errorf(msg)
	}

	l.Info("Location instance ID >%s< name >%s<", locationInstanceRecordSet.LocationInstanceViewRec.ID, locationInstanceRecordSet.LocationInstanceViewRec.Name)
	l.Info("Location instance ID >%s< record set character instance recs >%d<", locationInstanceRecordSet.LocationInstanceViewRec.ID, len(locationInstanceRecordSet.CharacterInstanceViewRecs))
	l.Info("Location instance ID >%s< record set monster instance recs >%d<", locationInstanceRecordSet.LocationInstanceViewRec.ID, len(locationInstanceRecordSet.MonsterInstanceViewRecs))
	l.Info("Location instance ID >%s< record set object instance recs >%d<", locationInstanceRecordSet.LocationInstanceViewRec.ID, len(locationInstanceRecordSet.ObjectInstanceViewRecs))

	// Get the monsters memory action records
	memories, err := m.GetMonsterInstanceMemories(rec)
	if err != nil {
		l.Warn("failed getting monster instance action memory records >%v<", err)
		return nil, err
	}

	sentence, err := m.decideAction(&DeciderArgs{
		MonsterInstanceViewRec:    rec,
		LocationInstanceRecordSet: locationInstanceRecordSet,
		Memories:                  memories,
	})
	if err != nil {
		l.Warn("failed deciding action >%v<", err)
		return nil, err
	}

	return &DecideMonsterActionResult{
		DungeonInstanceID: locationInstanceRecordSet.LocationInstanceViewRec.DungeonInstanceID,
		MonsterInstanceID: monsterInstanceID,
		Sentence:          sentence,
	}, nil
}

// ProcessMonsterAction - Processes a submitted character action
func (m *Model) ProcessMonsterAction(dungeonInstanceID string, monsterInstanceID string, sentence string) (*record.ActionRecordSet, error) {
	l := m.Logger("ProcessMonsterAction")

	l.Info("Processing monster ID >%s< action command >%s<", monsterInstanceID, sentence)

	// Verify the monster performing the action exists within the specified dungeon
	mivRec, err := m.GetMonsterInstanceViewRec(monsterInstanceID)
	if err != nil {
		l.Warn("failed getting monster record before performing action >%v<", err)
		return nil, err
	}
	if mivRec == nil {
		msg := fmt.Sprintf("failed getting monster record ID >%s< before performing action", monsterInstanceID)
		l.Warn(msg)
		return nil, fmt.Errorf(msg)
	}

	if mivRec.DungeonInstanceID != dungeonInstanceID {
		msg := fmt.Sprintf("monster ID >%s< does not exist in dungeon ID >%s<", monsterInstanceID, dungeonInstanceID)
		l.Warn(msg)
		return nil, fmt.Errorf(msg)
	}

	// TODO: 10-implement-effects:
	// Process any active effects that are still applied to the monster.

	// TODO: 12-implement-death: Remove monster instance when dead

	// Get the current dungeon location set of related records
	locationInstanceRecordSet, err := m.GetLocationInstanceViewRecordSet(mivRec.LocationInstanceID, true)
	if err != nil {
		l.Warn("failed getting dungeon location record set before performing action >%v<", err)
		return nil, err
	}
	if locationInstanceRecordSet == nil {
		msg := fmt.Sprintf("failed getting dungeon location record ID >%s< set before performing action", mivRec.LocationInstanceID)
		l.Warn(msg)
		return nil, fmt.Errorf(msg)
	}

	// Resolve the submitted monster action
	actionRec, err := m.resolveAction(sentence, &ResolveActionArgs{
		EntityType:                EntityTypeMonster,
		EntityInstanceID:          mivRec.ID,
		LocationInstanceRecordSet: locationInstanceRecordSet,
	})
	if err != nil {
		l.Warn("failed resolving monster action >%v<", err)
		return nil, err
	}

	// Resolve the initial action turn
	actionRec, err = m.resolveActionTurn(&ResolveActionTurnArgs{
		ActionRec:         actionRec,
		EntityType:        EntityTypeMonster,
		EntityInstanceID:  mivRec.ID,
		DungeonInstanceID: locationInstanceRecordSet.LocationInstanceViewRec.DungeonInstanceID,
	})
	if err != nil {
		l.Warn("failed resolving action turn >%v<", err)
		return nil, err
	}

	// Perform the submitted monster action
	actionRec, err = m.performAction(&PerformActionArgs{
		ActionRec:                 actionRec,
		CharacterInstanceViewRec:  nil,
		MonsterInstanceViewRec:    mivRec,
		LocationInstanceRecordSet: locationInstanceRecordSet,
	})
	if err != nil {
		l.Warn("failed performing monster action >%v<", err)
		return nil, err
	}

	// Create the resulting action event record
	err = m.CreateActionRec(actionRec)
	if err != nil {
		l.Warn("failed creating action record >%v<", err)
		return nil, err
	}

	// Refetch the resulting action event record so we ave its serial number
	actionRec, err = m.GetActionRec(actionRec.ID, nil)
	if err != nil {
		l.Warn("failed refetching action record >%v<", err)
		return nil, err
	}

	l.Info("Created action record ID >%s< SerialNumber >%d<", actionRec.ID, null.NullInt16ToInt16(actionRec.SerialNumber))

	// Get the updated monster record
	mivRec, err = m.GetMonsterInstanceViewRec(monsterInstanceID)
	if err != nil {
		l.Warn("failed getting monster record after performing action >%v<", err)
		return nil, err
	}

	// Create action monster record
	actionMonsterRec := record.ActionMonster{
		RecordType:          record.ActionMonsterRecordTypeSource,
		ActionID:            actionRec.ID,
		LocationInstanceID:  actionRec.LocationInstanceID,
		MonsterInstanceID:   mivRec.ID,
		Name:                mivRec.Name,
		Strength:            mivRec.Strength,
		Dexterity:           mivRec.Dexterity,
		Intelligence:        mivRec.Intelligence,
		CurrentStrength:     mivRec.CurrentStrength,
		CurrentDexterity:    mivRec.CurrentDexterity,
		CurrentIntelligence: mivRec.CurrentIntelligence,
		Health:              mivRec.Health,
		Fatigue:             mivRec.Fatigue,
		CurrentHealth:       mivRec.CurrentHealth,
		CurrentFatigue:      mivRec.CurrentFatigue,
	}

	err = m.CreateActionMonsterRec(&actionMonsterRec)
	if err != nil {
		l.Warn("failed creating source action monster record >%v<", err)
		return nil, err
	}

	// Create action monster object records
	oivRecs, err := m.GetMonsterInstanceObjectInstanceViewRecs(mivRec.ID)
	if err != nil {
		l.Warn("failed getting source monster object instance view records >%v<", err)
		return nil, err
	}

	actionMonsterObjectRecs := []*record.ActionMonsterObject{}
	for _, oivRec := range oivRecs {
		l.Info("Adding monster action object record >%#v<", oivRec)
		dungeonMonsterObjectRec := record.ActionMonsterObject{
			ActionMonsterID:  actionMonsterRec.ID,
			ObjectInstanceID: oivRec.ID,
			Name:             oivRec.Name,
			IsEquipped:       oivRec.IsEquipped,
			IsStashed:        oivRec.IsStashed,
		}
		err := m.CreateActionMonsterObjectRec(&dungeonMonsterObjectRec)
		if err != nil {
			l.Warn("failed creating source action character object record >%v<", err)
			return nil, err
		}
		actionMonsterObjectRecs = append(actionMonsterObjectRecs, &dungeonMonsterObjectRec)
	}

	actionRecordSet, err := m.createActionRecordSetRecords(&record.ActionRecordSet{
		ActionRec:               actionRec,
		ActionMonsterRec:        &actionMonsterRec,
		ActionMonsterObjectRecs: actionMonsterObjectRecs,
	})
	if err != nil {
		l.Warn("failed creating action record set records >%v<", err)
		return nil, err
	}

	return actionRecordSet, nil
}

func (m *Model) GetActionRecordSet(actionID string) (*record.ActionRecordSet, error) {
	l := m.Logger("GetActionRecordSet")

	actionRecordSet := record.ActionRecordSet{}

	actionRec, err := m.GetActionRec(actionID, nil)
	if err != nil {
		l.Warn("failed getting action record >%v<", err)
		return nil, err
	}
	actionRecordSet.ActionRec = actionRec

	// Add the source action character record that performed the action.
	if actionRec.CharacterInstanceID.Valid {
		actionCharacterRecs, err := m.GetActionCharacterRecs(
			&coresql.Options{
				Params: []coresql.Param{
					{
						Col: "record_type",
						Val: record.ActionCharacterRecordTypeSource,
						Op:  coresql.OpEqual,
					},
					{
						Col: "action_id",
						Val: actionID,
						Op:  coresql.OpEqual,
					},
					{
						Col: "character_instance_id",
						Val: actionRec.CharacterInstanceID.String,
						Op:  coresql.OpEqual,
					},
				},
			},
		)
		if err != nil {
			l.Warn("failed getting action character records >%v<", err)
			return nil, err
		}
		if len(actionCharacterRecs) != 1 {
			msg := fmt.Sprintf("Unexpected number of action character records returned >%d<", len(actionCharacterRecs))
			l.Warn(msg)
			return nil, fmt.Errorf(msg)
		}
		actionRecordSet.ActionCharacterRec = actionCharacterRecs[0]

		actionCharacterObjectRecs, err := m.GetActionCharacterObjectRecs(
			&coresql.Options{
				Params: []coresql.Param{
					{
						Col: "action_character_id",
						Val: actionRecordSet.ActionCharacterRec.ID,
					},
				},
			},
		)
		if err != nil {
			l.Warn("failed getting action character object records >%v<", err)
			return nil, err
		}
		actionRecordSet.ActionCharacterObjectRecs = actionCharacterObjectRecs
	}

	// Add the source action monster record that performed the action.
	if actionRec.MonsterInstanceID.Valid {
		actionMonsterRecs, err := m.GetActionMonsterRecs(
			&coresql.Options{
				Params: []coresql.Param{
					{
						Col: "record_type",
						Val: record.ActionMonsterRecordTypeSource,
					},
					{
						Col: "action_id",
						Val: actionID,
					},
					{
						Col: "monster_instance_id",
						Val: actionRec.MonsterInstanceID.String,
					},
				},
			},
		)
		if err != nil {
			l.Warn("failed getting action monster records >%v<", err)
			return nil, err
		}
		if len(actionMonsterRecs) != 1 {
			msg := fmt.Sprintf("Unexpected number of action monster records returned >%d<", len(actionMonsterRecs))
			l.Warn(msg)
			return nil, fmt.Errorf(msg)
		}
		actionRecordSet.ActionMonsterRec = actionMonsterRecs[0]

		actionMonsterObjectRecs, err := m.GetActionMonsterObjectRecs(
			&coresql.Options{
				Params: []coresql.Param{
					{
						Col: "action_monster_id",
						Val: actionRecordSet.ActionMonsterRec.ID,
					},
				},
			},
		)
		if err != nil {
			l.Warn("failed getting action monster object records >%v<", err)
			return nil, err
		}
		actionRecordSet.ActionMonsterObjectRecs = actionMonsterObjectRecs
	}

	// Add the current location record set where the action was performed.
	locationInstanceViewRec, err := m.GetLocationInstanceViewRec(actionRec.LocationInstanceID)
	if err != nil {
		l.Warn("failed getting location instance view record >%v<", err)
		return nil, err
	}

	currentLocationRecordSet := record.ActionLocationRecordSet{
		LocationInstanceViewRec: locationInstanceViewRec,
		ActionCharacterRecs:     []*record.ActionCharacter{},
		ActionMonsterRecs:       []*record.ActionMonster{},
		ActionObjectRecs:        []*record.ActionObject{},
	}

	// Add the current location action character records
	actionCharacterRecs, err := m.GetActionCharacterRecs(
		&coresql.Options{
			Params: []coresql.Param{
				{
					Col: "record_type",
					Val: record.ActionCharacterRecordTypeCurrentLocation,
				},
				{
					Col: "action_id",
					Val: actionID,
				},
				{
					Col: "location_instance_id",
					Val: locationInstanceViewRec.ID,
				},
			},
		},
	)
	if err != nil {
		l.Warn("failed getting current location occupant action character records >%v<", err)
		return nil, err
	}

	currentLocationRecordSet.ActionCharacterRecs = actionCharacterRecs

	// Add the current location action monster records
	actionMonsterRecs, err := m.GetActionMonsterRecs(
		&coresql.Options{
			Params: []coresql.Param{
				{
					Col: "record_type",
					Val: record.ActionMonsterRecordTypeCurrentLocation,
				},
				{
					Col: "action_id",
					Val: actionID,
				},
				{
					Col: "location_instance_id",
					Val: locationInstanceViewRec.ID,
				},
			},
		},
	)
	if err != nil {
		l.Warn("failed getting current location occupant action monster records >%v<", err)
		return nil, err
	}

	currentLocationRecordSet.ActionMonsterRecs = actionMonsterRecs

	// Add the current location action object records
	dungeonActionObjectRecs, err := m.GetActionObjectRecs(
		&coresql.Options{
			Params: []coresql.Param{
				{
					Col: "record_type",
					Val: record.ActionObjectRecordTypeCurrentLocation,
				},
				{
					Col: "action_id",
					Val: actionID,
				},
				{
					Col: "location_instance_id",
					Val: locationInstanceViewRec.ID,
				},
			},
		},
	)
	if err != nil {
		l.Warn("failed getting current location occupant action monster records >%v<", err)
		return nil, err
	}
	currentLocationRecordSet.ActionObjectRecs = dungeonActionObjectRecs

	actionRecordSet.CurrentLocation = &currentLocationRecordSet

	// Get the target dungeon location record set when set
	if actionRec.ResolvedTargetLocationInstanceID.Valid {

		// Add the target location record set when the action was performed.
		locationInstanceViewRec, err := m.GetLocationInstanceViewRec(actionRec.ResolvedTargetLocationInstanceID.String)
		if err != nil {
			l.Warn("failed getting target location instance view record >%v<", err)
			return nil, err
		}

		targetLocationRecordSet := record.ActionLocationRecordSet{
			LocationInstanceViewRec: locationInstanceViewRec,
			ActionCharacterRecs:     []*record.ActionCharacter{},
			ActionMonsterRecs:       []*record.ActionMonster{},
			ActionObjectRecs:        []*record.ActionObject{},
		}

		// Add the target location occupant action character records
		actionCharacterRecs, err := m.GetActionCharacterRecs(
			&coresql.Options{
				Params: []coresql.Param{
					{
						Col: "record_type",
						Val: record.ActionCharacterRecordTypeTargetLocation,
					},
					{
						Col: "action_id",
						Val: actionID,
					},
					{
						Col: "location_instance_id",
						Val: locationInstanceViewRec.ID,
					},
				},
			},
		)
		if err != nil {
			l.Warn("failed getting target location occupant action character records >%v<", err)
			return nil, err
		}
		targetLocationRecordSet.ActionCharacterRecs = actionCharacterRecs

		// Add the target location occupant action monster records
		actionMonsterRecs, err := m.GetActionMonsterRecs(
			&coresql.Options{
				Params: []coresql.Param{
					{
						Col: "record_type",
						Val: record.ActionMonsterRecordTypeTargetLocation,
					},
					{
						Col: "action_id",
						Val: actionID,
					},
					{
						Col: "location_instance_id",
						Val: locationInstanceViewRec.ID,
					},
				},
			},
		)
		if err != nil {
			l.Warn("failed getting target location occupant action monster records >%v<", err)
			return nil, err
		}
		targetLocationRecordSet.ActionMonsterRecs = actionMonsterRecs

		// Add the target location occupant action object records
		actionObjectRecs, err := m.GetActionObjectRecs(
			&coresql.Options{
				Params: []coresql.Param{
					{
						Col: "record_type",
						Val: record.ActionObjectRecordTypeTargetLocation,
					},
					{
						Col: "action_id",
						Val: actionID,
					},
					{
						Col: "location_instance_id",
						Val: locationInstanceViewRec.ID,
					},
				},
			},
		)
		if err != nil {
			l.Warn("failed getting target location occupant action monster records >%v<", err)
			return nil, err
		}
		targetLocationRecordSet.ActionObjectRecs = actionObjectRecs

		actionRecordSet.TargetLocation = &targetLocationRecordSet
	}

	// Get the target character action record
	if actionRec.ResolvedTargetCharacterInstanceID.Valid {
		actionCharacterRecs, err := m.GetActionCharacterRecs(
			&coresql.Options{
				Params: []coresql.Param{
					{
						Col: "record_type",
						Val: record.ActionCharacterRecordTypeTarget,
					},
					{
						Col: "action_id",
						Val: actionID,
					},
					{
						Col: "character_instance_id",
						Val: actionRec.ResolvedTargetCharacterInstanceID.String,
					},
				},
			},
		)
		if err != nil {
			l.Warn("failed getting target action character record >%v<", err)
			return nil, err
		}
		if len(actionCharacterRecs) != 1 {
			msg := fmt.Sprintf("Unexpected number of action character records returned >%d<", len(actionCharacterRecs))
			l.Warn(msg)
			return nil, fmt.Errorf(msg)
		}
		actionRecordSet.TargetActionCharacterRec = actionCharacterRecs[0]

		actionCharacterObjectRecs, err := m.GetActionCharacterObjectRecs(
			&coresql.Options{
				Params: []coresql.Param{
					{
						Col: "action_character_id",
						Val: actionRecordSet.TargetActionCharacterRec.ID,
					},
				},
			},
		)
		if err != nil {
			l.Warn("failed getting target character object records >%v<", err)
			return nil, err
		}
		actionRecordSet.TargetActionCharacterObjectRecs = actionCharacterObjectRecs
	}

	// Get the target dungeon monster action record
	if actionRec.ResolvedTargetMonsterInstanceID.Valid {
		actionMonsterRecs, err := m.GetActionMonsterRecs(
			&coresql.Options{
				Params: []coresql.Param{
					{
						Col: "record_type",
						Val: record.ActionMonsterRecordTypeTarget,
					},
					{
						Col: "action_id",
						Val: actionID,
					},
					{
						Col: "monster_instance_id",
						Val: actionRec.ResolvedTargetMonsterInstanceID.String,
					},
				},
			},
		)
		if err != nil {
			l.Warn("failed getting target action monster record >%v<", err)
			return nil, err
		}
		if len(actionMonsterRecs) != 1 {
			msg := fmt.Sprintf("Unexpected number of action monster records returned >%d<", len(actionMonsterRecs))
			l.Warn(msg)
			return nil, fmt.Errorf(msg)
		}
		actionRecordSet.TargetActionMonsterRec = actionMonsterRecs[0]

		actionMonsterObjectRecs, err := m.GetActionMonsterObjectRecs(
			&coresql.Options{
				Params: []coresql.Param{
					{
						Col: "action_monster_id",
						Val: actionRecordSet.TargetActionMonsterRec.ID,
					},
				},
			},
		)
		if err != nil {
			l.Warn("failed getting target monster object records >%v<", err)
			return nil, err
		}
		actionRecordSet.TargetActionMonsterObjectRecs = actionMonsterObjectRecs
	}

	// Get the target dungeon object action record
	if actionRec.ResolvedTargetObjectInstanceID.Valid {
		dungeonActionObjectRecs, err := m.GetActionObjectRecs(
			&coresql.Options{
				Params: []coresql.Param{
					{
						Col: "record_type",
						Val: record.ActionObjectRecordTypeTarget,
					},
					{
						Col: "action_id",
						Val: actionID,
					},
					{
						Col: "object_instance_id",
						Val: actionRec.ResolvedTargetObjectInstanceID.String,
					},
				},
			},
		)
		if err != nil {
			l.Warn("failed getting target action object record >%v<", err)
			return nil, err
		}
		if len(dungeonActionObjectRecs) != 1 {
			msg := fmt.Sprintf("Unexpected number of action object records returned >%d<", len(dungeonActionObjectRecs))
			l.Warn(msg)
			return nil, fmt.Errorf(msg)
		}
		actionRecordSet.TargetActionObjectRec = dungeonActionObjectRecs[0]
	}

	// Get the stashed dungeon object action record
	if actionRec.ResolvedStashedObjectInstanceID.Valid {
		dungeonActionObjectRecs, err := m.GetActionObjectRecs(
			&coresql.Options{
				Params: []coresql.Param{
					{
						Col: "record_type",
						Val: record.ActionObjectRecordTypeStashed,
					},
					{
						Col: "action_id",
						Val: actionID,
					},
					{
						Col: "object_instance_id",
						Val: actionRec.ResolvedStashedObjectInstanceID.String,
					},
				},
			},
		)
		if err != nil {
			l.Warn("failed getting stashed action object record >%v<", err)
			return nil, err
		}
		if len(dungeonActionObjectRecs) != 1 {
			msg := fmt.Sprintf("Unexpected number of action object records returned >%d<", len(dungeonActionObjectRecs))
			l.Warn(msg)
			return nil, fmt.Errorf(msg)
		}
		actionRecordSet.StashedActionObjectRec = dungeonActionObjectRecs[0]
	}

	// Get the equipped dungeon object action record
	if actionRec.ResolvedEquippedObjectInstanceID.Valid {
		dungeonActionObjectRecs, err := m.GetActionObjectRecs(
			&coresql.Options{
				Params: []coresql.Param{
					{
						Col: "record_type",
						Val: record.ActionObjectRecordTypeEquipped,
					},
					{
						Col: "action_id",
						Val: actionID,
					},
					{
						Col: "object_instance_id",
						Val: actionRec.ResolvedEquippedObjectInstanceID.String,
					},
				},
			},
		)
		if err != nil {
			l.Warn("failed getting equipped action object record >%v<", err)
			return nil, err
		}
		if len(dungeonActionObjectRecs) != 1 {
			msg := fmt.Sprintf("Unexpected number of action object records returned >%d<", len(dungeonActionObjectRecs))
			l.Warn(msg)
			return nil, fmt.Errorf(msg)
		}
		actionRecordSet.EquippedActionObjectRec = dungeonActionObjectRecs[0]
	}

	// Get the dropped dungeon object action record
	if actionRec.ResolvedDroppedObjectInstanceID.Valid {
		dungeonActionObjectRecs, err := m.GetActionObjectRecs(
			&coresql.Options{
				Params: []coresql.Param{
					{
						Col: "record_type",
						Val: record.ActionObjectRecordTypeDropped,
					},
					{
						Col: "action_id",
						Val: actionID,
					},
					{
						Col: "object_instance_id",
						Val: actionRec.ResolvedDroppedObjectInstanceID.String,
					},
				},
			},
		)
		if err != nil {
			l.Warn("failed getting dropped action object record >%v<", err)
			return nil, err
		}
		if len(dungeonActionObjectRecs) != 1 {
			msg := fmt.Sprintf("Unexpected number of action object records returned >%d<", len(dungeonActionObjectRecs))
			l.Warn(msg)
			return nil, fmt.Errorf(msg)
		}
		actionRecordSet.DroppedActionObjectRec = dungeonActionObjectRecs[0]
	}

	return &actionRecordSet, nil
}

func (m *Model) GetLocationInstanceViewRecordSet(locationInstanceID string, forUpdate bool) (*record.LocationInstanceViewRecordSet, error) {
	l := m.Logger("GetLocationInstanceViewRecordSet")

	locationInstanceRecordSet := &record.LocationInstanceViewRecordSet{}

	// Location record
	locationInstanceViewRec, err := m.GetLocationInstanceViewRec(locationInstanceID)
	if err != nil {
		l.Warn("failed to get dungeon location record >%v<", err)
		return nil, err
	}
	locationInstanceRecordSet.LocationInstanceViewRec = locationInstanceViewRec

	// All characters at location
	characterInstanceViewRecs, err := m.GetCharacterInstanceViewRecs(
		&coresql.Options{
			Params: []coresql.Param{
				{
					Col: "location_instance_id",
					Val: locationInstanceViewRec.ID,
				},
			},
		},
	)
	if err != nil {
		l.Warn("failed to get dungeon location character records >%v<", err)
		return nil, err
	}
	locationInstanceRecordSet.CharacterInstanceViewRecs = characterInstanceViewRecs

	// All monsters at location
	monsterInstanceViewRecs, err := m.GetMonsterInstanceViewRecs(
		&coresql.Options{
			Params: []coresql.Param{
				{
					Col: "location_instance_id",
					Val: locationInstanceViewRec.ID,
				},
			},
		},
	)
	if err != nil {
		l.Warn("failed to get dungeon location monster records >%v<", err)
		return nil, err
	}
	locationInstanceRecordSet.MonsterInstanceViewRecs = monsterInstanceViewRecs

	// All objects at location
	objectInstanceViewRecs, err := m.GetObjectInstanceViewRecs(
		&coresql.Options{
			Params: []coresql.Param{
				{
					Col: "location_instance_id",
					Val: locationInstanceViewRec.ID,
				},
			},
		},
	)
	if err != nil {
		l.Warn("failed to get dungeon location object records >%v<", err)
		return nil, err
	}
	locationInstanceRecordSet.ObjectInstanceViewRecs = objectInstanceViewRecs

	locationInstanceIDs := []any{}
	if locationInstanceViewRec.NorthLocationInstanceID.Valid {
		locationInstanceIDs = append(locationInstanceIDs, locationInstanceViewRec.NorthLocationInstanceID.String)
	}
	if locationInstanceViewRec.NortheastLocationInstanceID.Valid {
		locationInstanceIDs = append(locationInstanceIDs, locationInstanceViewRec.NortheastLocationInstanceID.String)
	}
	if locationInstanceViewRec.EastLocationInstanceID.Valid {
		locationInstanceIDs = append(locationInstanceIDs, locationInstanceViewRec.EastLocationInstanceID.String)
	}
	if locationInstanceViewRec.SoutheastLocationInstanceID.Valid {
		locationInstanceIDs = append(locationInstanceIDs, locationInstanceViewRec.SoutheastLocationInstanceID.String)
	}
	if locationInstanceViewRec.SouthLocationInstanceID.Valid {
		locationInstanceIDs = append(locationInstanceIDs, locationInstanceViewRec.SouthLocationInstanceID.String)
	}
	if locationInstanceViewRec.SouthwestLocationInstanceID.Valid {
		locationInstanceIDs = append(locationInstanceIDs, locationInstanceViewRec.SouthwestLocationInstanceID.String)
	}
	if locationInstanceViewRec.WestLocationInstanceID.Valid {
		locationInstanceIDs = append(locationInstanceIDs, locationInstanceViewRec.WestLocationInstanceID.String)
	}
	if locationInstanceViewRec.NorthwestLocationInstanceID.Valid {
		locationInstanceIDs = append(locationInstanceIDs, locationInstanceViewRec.NorthwestLocationInstanceID.String)
	}
	if locationInstanceViewRec.UpLocationInstanceID.Valid {
		locationInstanceIDs = append(locationInstanceIDs, locationInstanceViewRec.UpLocationInstanceID.String)
	}
	if locationInstanceViewRec.DownLocationInstanceID.Valid {
		locationInstanceIDs = append(locationInstanceIDs, locationInstanceViewRec.DownLocationInstanceID.String)
	}

	locationInstanceViewRecs, err := m.GetLocationInstanceViewRecs(
		&coresql.Options{
			Params: []coresql.Param{
				{
					Col:   "id",
					Array: locationInstanceIDs,
				},
			},
		},
	)
	if err != nil {
		l.Warn("failed to get dungeon location direction location records >%v<", err)
		return nil, err
	}

	locationInstanceRecordSet.LocationInstanceViewRecs = locationInstanceViewRecs

	return locationInstanceRecordSet, nil
}

// Returns all action records that occurred at the location of the previous action
// for the entity associated with the given action record. The given action record
// is then appended to the result providing the full list.
func (m *Model) GetActionRecsSincePreviousAction(rec *record.Action) ([]*record.Action, error) {
	l := m.Logger("GetActionRecsSincePreviousAction")

	if rec == nil {
		return nil, fmt.Errorf("missing action record argument, cannot get action record since previous action")
	}

	if !null.NullStringIsValid(rec.CharacterInstanceID) {
		return nil, nil
	}

	l.Info("Current action record ID >%s<", rec.ID)
	l.Info("Current action record location instance ID >%s<", rec.LocationInstanceID)
	l.Info("Current action record turn number >%d<", rec.TurnNumber)
	l.Info("Current action record serial number >%d<", null.NullInt16ToInt16(rec.SerialNumber))
	l.Info("Current action record character instance ID >%s<", null.NullStringToString(rec.CharacterInstanceID))
	l.Info("Current action record monster instance ID >%s<", null.NullStringToString(rec.MonsterInstanceID))

	actionRecs, err := m.GetActionRecs(
		&coresql.Options{
			Params: []coresql.Param{
				{
					Col: "character_instance_id",
					Val: null.NullStringToString(rec.CharacterInstanceID),
				},
				{
					Col: "turn_number",
					Val: rec.TurnNumber,
					Op:  coresql.OpLessThan,
				},
			},
			OrderBy: []coresql.OrderBy{
				{
					Col:       "turn_number",
					Direction: coresql.OrderDirectionDESC,
				},
			},
			Limit: 1,
		},
	)
	if err != nil {
		l.Warn("failed getting previous action record >%v<", err)
		return nil, err
	}

	if len(actionRecs) != 1 {
		l.Info("Character instance ID >%s< has no previous action records", null.NullStringToString(rec.CharacterInstanceID))
		actionRecs = append(actionRecs, rec)
		return actionRecs, nil
	}

	prevActionRec := actionRecs[0]

	l.Info("Previous action record ID >%s<", prevActionRec.ID)
	l.Info("Previous action record location instance ID >%s<", prevActionRec.LocationInstanceID)
	l.Info("Previous action record turn number >%d<", prevActionRec.TurnNumber)
	l.Info("Previous action record serial number >%d<", null.NullInt16ToInt16(prevActionRec.SerialNumber))
	l.Info("Previous action record character instance ID >%s<", null.NullStringToString(prevActionRec.CharacterInstanceID))
	l.Info("Previous action record monster instance ID >%s<", null.NullStringToString(prevActionRec.MonsterInstanceID))

	// We add one to the previous action serial number and subtract one from the current action
	// serial number so we exclude those specific records when looking between.
	var adjustAmount int16 = 1
	actionRecs, err = m.GetActionRecs(
		&coresql.Options{
			Params: []coresql.Param{
				{
					Col: "location_instance_id",
					Val: prevActionRec.LocationInstanceID,
				},
				{
					Col:  "serial_number",
					Val:  fmt.Sprintf("%d", null.NullInt16ToInt16(prevActionRec.SerialNumber)+adjustAmount),
					ValB: fmt.Sprintf("%d", null.NullInt16ToInt16(rec.SerialNumber)-adjustAmount),
					Op:   coresql.OpBetween,
				},
			},
		},
	)
	if err != nil {
		l.Warn("failed getting action records since serial number >%s< >%v<", err)
		return nil, err
	}

	// Append current action
	actionRecs = append(actionRecs, rec)

	return actionRecs, nil
}

// TODO: We need more than just the action records, we also need the characters and monsters
// that were at the location the action occurred, as those details are valid memories that
// need to be referenced when deciding what action to take..

func (m *Model) GetMonsterInstanceMemories(rec *record.MonsterInstanceView) ([]*Memory, error) {
	l := m.Logger("GetMonsterInstanceMemories")

	var memories []*Memory

	// Maybe need a view here to make it easier to union monster actions and other actions?
	maRecs, err := m.GetActionRecs(
		&coresql.Options{
			Params: []coresql.Param{
				{
					Col: "monster_instance_id",
					Val: rec.ID,
				},
			},
			Limit: rec.CurrentIntelligence,
			OrderBy: []coresql.OrderBy{
				{
					Col:       "created_at",
					Direction: coresql.OrderDirectionDESC,
				},
			},
		},
	)
	if err != nil {
		l.Warn("failed fetching monster action recs >%v<", err)
		return nil, err
	}

	oaRecs, err := m.GetActionRecs(
		&coresql.Options{
			Params: []coresql.Param{
				{
					Col: "resolved_target_monster_instance_id",
					Val: rec.ID,
				},
			},
			Limit: rec.CurrentIntelligence,
			OrderBy: []coresql.OrderBy{
				{
					Col:       "created_at",
					Direction: coresql.OrderDirectionDESC,
				},
			},
		},
	)
	if err != nil {
		l.Warn("failed fetching resolved monster target action recs >%v<", err)
		return nil, err
	}

	l.Info("Have >%d< monster action records", len(maRecs))
	l.Info("Have >%d< resolved monster target action records", len(oaRecs))

	for len(maRecs) > 0 {
		if len(oaRecs) > 0 {
			if null.NullInt16ToInt16(maRecs[0].SerialNumber) > null.NullInt16ToInt16(oaRecs[0].SerialNumber) {
				memories = append(memories, &Memory{
					ActionRec: maRecs[0],
				})
				maRecs = maRecs[1:]
			} else {
				memories = append(memories, &Memory{
					ActionRec: oaRecs[0],
				})
				oaRecs = oaRecs[1:]
			}
			continue
		}
		memories = append(memories, &Memory{
			ActionRec: maRecs[0],
		})
		maRecs = maRecs[1:]
	}

	for idx := range memories {

		actionCharacterRecs, err := m.GetActionCharacterRecs(
			&coresql.Options{
				Params: []coresql.Param{
					{
						Col: "action_id",
						Val: memories[idx].ActionRec.ID,
					},
				},
			},
		)
		if err != nil {
			l.Warn("failed getting action ID >%s< action character records >%v<", memories[idx].ActionRec.ID, err)
			return nil, err
		}

		memories[idx].ActionCharacterRecs = actionCharacterRecs

		actionMonsterRecs, err := m.GetActionMonsterRecs(
			&coresql.Options{
				Params: []coresql.Param{
					{
						Col: "action_id",
						Val: memories[idx].ActionRec.ID,
					},
				},
			},
		)
		if err != nil {
			l.Warn("failed getting action ID >%s< action monster records >%v<", memories[idx].ActionRec.ID, err)
			return nil, err
		}

		memories[idx].ActionMonsterRecs = actionMonsterRecs
	}

	// TODO: 15-implement-monster-goals
	// Add action object records so monsters can look for specific items
	// as one of their goals.
	l.Info("Have >%d< total memories", len(memories))

	return memories, nil
}

func (m *Model) createActionRecordSetRecords(actionRecordSet *record.ActionRecordSet) (*record.ActionRecordSet, error) {
	l := m.Logger("createActionRecordSetRecords")

	actionRec := actionRecordSet.ActionRec

	// Create current location record set
	currentLocationRecordSet, err := m.createCurrentActionLocationRecordSet(actionRec.ID, actionRec.LocationInstanceID)
	if err != nil {
		l.Warn("failed creating action location record set >%v<", err)
		return nil, err
	}
	actionRecordSet.CurrentLocation = currentLocationRecordSet

	// Create target location record set
	if actionRec.ResolvedTargetLocationInstanceID.Valid {
		targetLocationRecordSet, err := m.createTargetActionLocationRecordSet(actionRec.ID, null.NullStringToString(actionRec.ResolvedTargetLocationInstanceID))
		if err != nil {
			l.Warn("failed creating target action location record set >%v<", err)
			return nil, err
		}
		actionRecordSet.TargetLocation = targetLocationRecordSet
	}

	// Create the target character action record
	if actionRec.ResolvedTargetCharacterInstanceID.Valid {
		actionCharacterRec, actionCharacterObjectRecs, err := m.createActionTargetCharacterRecs(actionRec.ID, actionRec.LocationInstanceID, null.NullStringToString(actionRec.ResolvedTargetCharacterInstanceID))
		if err != nil {
			l.Warn("failed create action character records >%v<", err)
			return nil, err
		}
		actionRecordSet.TargetActionCharacterRec = actionCharacterRec
		actionRecordSet.TargetActionCharacterObjectRecs = actionCharacterObjectRecs
	}

	// Create the target dungeon monster action record
	if actionRec.ResolvedTargetMonsterInstanceID.Valid {
		actionMonsterRec, actionMonsterObjectRecs, err := m.createActionTargetMonsterRecs(actionRec.ID, actionRec.LocationInstanceID, null.NullStringToString(actionRec.ResolvedTargetMonsterInstanceID))
		if err != nil {
			l.Warn("failed create action character records >%v<", err)
			return nil, err
		}
		actionRecordSet.TargetActionMonsterRec = actionMonsterRec
		actionRecordSet.TargetActionMonsterObjectRecs = actionMonsterObjectRecs
	}

	// Create the target dungeon object action record
	if actionRec.ResolvedTargetObjectInstanceID.Valid {
		actionObjectRec, err := m.createActionObjectRec(
			actionRec.ID,
			actionRec.LocationInstanceID,
			null.NullStringToString(actionRec.ResolvedTargetObjectInstanceID),
			record.ActionObjectRecordTypeTarget,
		)
		if err != nil {
			l.Warn("failed creating action target object record >%v<", err)
			return nil, err
		}
		actionRecordSet.TargetActionObjectRec = actionObjectRec
	}

	// Create the stashed dungeon object action record
	if actionRec.ResolvedStashedObjectInstanceID.Valid {
		actionObjectRec, err := m.createActionObjectRec(
			actionRec.ID,
			actionRec.LocationInstanceID,
			null.NullStringToString(actionRec.ResolvedStashedObjectInstanceID),
			record.ActionObjectRecordTypeStashed,
		)
		if err != nil {
			l.Warn("failed creating action stashed object record >%v<", err)
			return nil, err
		}
		actionRecordSet.StashedActionObjectRec = actionObjectRec
	}

	// Create the equipped dungeon object action record
	if actionRec.ResolvedEquippedObjectInstanceID.Valid {
		actionObjectRec, err := m.createActionObjectRec(
			actionRec.ID,
			actionRec.LocationInstanceID,
			null.NullStringToString(actionRec.ResolvedEquippedObjectInstanceID),
			record.ActionObjectRecordTypeEquipped,
		)
		if err != nil {
			l.Warn("failed creating action equipped object record >%v<", err)
			return nil, err
		}
		actionRecordSet.EquippedActionObjectRec = actionObjectRec
	}

	// Create the dropped dungeon object action record
	if actionRec.ResolvedDroppedObjectInstanceID.Valid {
		actionObjectRec, err := m.createActionObjectRec(
			actionRec.ID,
			actionRec.LocationInstanceID,
			null.NullStringToString(actionRec.ResolvedDroppedObjectInstanceID),
			record.ActionObjectRecordTypeDropped,
		)
		if err != nil {
			l.Warn("failed creating action equipped object record >%v<", err)
			return nil, err
		}
		actionRecordSet.DroppedActionObjectRec = actionObjectRec
	}

	return actionRecordSet, nil
}

type LocationType string

const LocationTypeCurrent LocationType = "current"
const LocationTypeTarget LocationType = "target"

func (m *Model) createCurrentActionLocationRecordSet(actionID, locationInstanceID string) (*record.ActionLocationRecordSet, error) {
	return m.createActionLocationRecordSet(actionID, locationInstanceID, LocationTypeCurrent)
}

func (m *Model) createTargetActionLocationRecordSet(actionID, locationInstanceID string) (*record.ActionLocationRecordSet, error) {
	return m.createActionLocationRecordSet(actionID, locationInstanceID, LocationTypeTarget)
}

func (m *Model) createActionLocationRecordSet(actionID, locationInstanceID string, locationType LocationType) (*record.ActionLocationRecordSet, error) {
	l := m.Logger("createActionLocationRecordSet")

	// TODO: Think we should be using turn_number here to get relevant records
	locationInstanceRecordSet, err := m.GetLocationInstanceViewRecordSet(locationInstanceID, true)
	if err != nil {
		l.Warn("failed getting dungeon location record set after performing action >%v<", err)
		return nil, err
	}
	locationInstanceViewRec := locationInstanceRecordSet.LocationInstanceViewRec

	l.Info("Dungeon location record set location name >%s<", locationInstanceRecordSet.LocationInstanceViewRec.Name)
	l.Info("Dungeon location record set characters >%d<", len(locationInstanceRecordSet.CharacterInstanceViewRecs))
	l.Info("Dungeon location record set monsters >%d<", len(locationInstanceRecordSet.MonsterInstanceViewRecs))
	l.Info("Dungeon location record set objects >%d<", len(locationInstanceRecordSet.ObjectInstanceViewRecs))

	currentLocationRecordSet := record.ActionLocationRecordSet{
		LocationInstanceViewRec: locationInstanceViewRec,
		ActionCharacterRecs:     []*record.ActionCharacter{},
		ActionMonsterRecs:       []*record.ActionMonster{},
		ActionObjectRecs:        []*record.ActionObject{},
	}

	// Character Occupants: Create the action character record for each character now at the current location
	if len(locationInstanceRecordSet.CharacterInstanceViewRecs) > 0 {

		actionRecordType := record.ActionCharacterRecordTypeCurrentLocation
		if locationType == LocationTypeTarget {
			actionRecordType = record.ActionCharacterRecordTypeTargetLocation
		}

		for _, characterInstanceViewRec := range locationInstanceRecordSet.CharacterInstanceViewRecs {
			actionCharacterRec := record.ActionCharacter{
				RecordType:          actionRecordType,
				ActionID:            actionID,
				LocationInstanceID:  locationInstanceViewRec.ID,
				CharacterInstanceID: characterInstanceViewRec.ID,
				Name:                characterInstanceViewRec.Name,
				Strength:            characterInstanceViewRec.Strength,
				Dexterity:           characterInstanceViewRec.Dexterity,
				Intelligence:        characterInstanceViewRec.Intelligence,
				CurrentStrength:     characterInstanceViewRec.CurrentStrength,
				CurrentDexterity:    characterInstanceViewRec.CurrentDexterity,
				CurrentIntelligence: characterInstanceViewRec.CurrentIntelligence,
				Health:              characterInstanceViewRec.Health,
				Fatigue:             characterInstanceViewRec.Fatigue,
				CurrentHealth:       characterInstanceViewRec.CurrentHealth,
				CurrentFatigue:      characterInstanceViewRec.CurrentFatigue,
			}

			err := m.CreateActionCharacterRec(&actionCharacterRec)
			if err != nil {
				l.Warn("failed creating current location action character record >%v<", err)
				return nil, err
			}

			l.Info("Created current location action character record ID >%s<", actionCharacterRec.ID)
			currentLocationRecordSet.ActionCharacterRecs = append(currentLocationRecordSet.ActionCharacterRecs, &actionCharacterRec)
		}
	}

	// Monster Occupants: Create the action monster record for each monster now at the current location
	if len(locationInstanceRecordSet.MonsterInstanceViewRecs) > 0 {

		actionRecordType := record.ActionMonsterRecordTypeCurrentLocation
		if locationType == LocationTypeTarget {
			actionRecordType = record.ActionMonsterRecordTypeTargetLocation
		}

		for _, monsterInstanceViewRec := range locationInstanceRecordSet.MonsterInstanceViewRecs {
			actionMonsterRec := record.ActionMonster{
				RecordType:          actionRecordType,
				ActionID:            actionID,
				LocationInstanceID:  locationInstanceViewRec.ID,
				MonsterInstanceID:   monsterInstanceViewRec.ID,
				Name:                monsterInstanceViewRec.Name,
				Strength:            monsterInstanceViewRec.Strength,
				Dexterity:           monsterInstanceViewRec.Dexterity,
				Intelligence:        monsterInstanceViewRec.Intelligence,
				CurrentStrength:     monsterInstanceViewRec.CurrentStrength,
				CurrentDexterity:    monsterInstanceViewRec.CurrentDexterity,
				CurrentIntelligence: monsterInstanceViewRec.CurrentIntelligence,
				Health:              monsterInstanceViewRec.Health,
				Fatigue:             monsterInstanceViewRec.Fatigue,
				CurrentHealth:       monsterInstanceViewRec.CurrentHealth,
				CurrentFatigue:      monsterInstanceViewRec.CurrentFatigue,
			}
			err := m.CreateActionMonsterRec(&actionMonsterRec)
			if err != nil {
				l.Warn("failed creating current location action monster record >%v<", err)
				return nil, err
			}

			l.Info("Created current location action monster record ID >%s<", actionMonsterRec.ID)
			currentLocationRecordSet.ActionMonsterRecs = append(currentLocationRecordSet.ActionMonsterRecs, &actionMonsterRec)
		}
	}

	// Object Occupants: Create the action object record for each object now at the current location
	if len(locationInstanceRecordSet.ObjectInstanceViewRecs) > 0 {

		actionRecordType := record.ActionObjectRecordTypeCurrentLocation
		if locationType == LocationTypeTarget {
			actionRecordType = record.ActionObjectRecordTypeTargetLocation
		}

		for _, objectInstanceViewRec := range locationInstanceRecordSet.ObjectInstanceViewRecs {
			dungeonActionObjectRec := record.ActionObject{
				RecordType:         actionRecordType,
				ActionID:           actionID,
				LocationInstanceID: locationInstanceViewRec.ID,
				ObjectInstanceID:   objectInstanceViewRec.ID,
				Name:               objectInstanceViewRec.Name,
				Description:        objectInstanceViewRec.Description,
				IsStashed:          objectInstanceViewRec.IsStashed,
				IsEquipped:         objectInstanceViewRec.IsEquipped,
			}
			err := m.CreateActionObjectRec(&dungeonActionObjectRec)
			if err != nil {
				l.Warn("failed creating current location action object record >%v<", err)
				return nil, err
			}

			l.Info("Created current location action object record ID >%s<", dungeonActionObjectRec.ID)
			currentLocationRecordSet.ActionObjectRecs = append(currentLocationRecordSet.ActionObjectRecs, &dungeonActionObjectRec)
		}
	}

	return &currentLocationRecordSet, nil
}

func (m *Model) createActionTargetCharacterRecs(actionID, locationInstanceID, characterInstanceID string) (*record.ActionCharacter, []*record.ActionCharacterObject, error) {
	l := m.Logger("createActionTargetCharacterRecs")

	l.Info("Creating action ID >%s< character instance ID >%s< records", actionID, characterInstanceID)

	targetCharacterInstanceViewRec, err := m.GetCharacterInstanceViewRec(characterInstanceID)
	if err != nil {
		l.Warn("failed getting target character instance view record >%v<", err)
		return nil, nil, err
	}

	rec := &record.ActionCharacter{
		RecordType:          record.ActionCharacterRecordTypeTarget,
		ActionID:            actionID,
		LocationInstanceID:  locationInstanceID,
		CharacterInstanceID: targetCharacterInstanceViewRec.ID,
		Name:                targetCharacterInstanceViewRec.Name,
		Strength:            targetCharacterInstanceViewRec.Strength,
		Dexterity:           targetCharacterInstanceViewRec.Dexterity,
		Intelligence:        targetCharacterInstanceViewRec.Intelligence,
		CurrentStrength:     targetCharacterInstanceViewRec.CurrentStrength,
		CurrentDexterity:    targetCharacterInstanceViewRec.CurrentDexterity,
		CurrentIntelligence: targetCharacterInstanceViewRec.CurrentIntelligence,
		Health:              targetCharacterInstanceViewRec.Health,
		Fatigue:             targetCharacterInstanceViewRec.Fatigue,
		CurrentHealth:       targetCharacterInstanceViewRec.CurrentHealth,
		CurrentFatigue:      targetCharacterInstanceViewRec.CurrentFatigue,
	}

	err = m.CreateActionCharacterRec(rec)
	if err != nil {
		l.Warn("failed creating target action character record >%v<", err)
		return nil, nil, err
	}

	// Create action character object records
	objectInstanceViewRecs, err := m.GetCharacterInstanceEquippedObjectInstanceViewRecs(targetCharacterInstanceViewRec.ID)
	if err != nil {
		l.Warn("failed getting target character object records >%v<", err)
		return nil, nil, err
	}

	l.Info("Adding >%d< target character object records", len(objectInstanceViewRecs))

	targetCharacterObjectRecs := []*record.ActionCharacterObject{}
	for _, objectInstanceViewRec := range objectInstanceViewRecs {
		l.Info("Adding target character object record >%v<", objectInstanceViewRecs)
		dungeonCharacterObjectRec := record.ActionCharacterObject{
			ActionCharacterID: rec.ID,
			ObjectInstanceID:  objectInstanceViewRec.ID,
			Name:              objectInstanceViewRec.Name,
			IsEquipped:        objectInstanceViewRec.IsEquipped,
			IsStashed:         objectInstanceViewRec.IsStashed,
		}
		err := m.CreateActionCharacterObjectRec(&dungeonCharacterObjectRec)
		if err != nil {
			l.Warn("failed creating source action character object record >%v<", err)
			return nil, nil, err
		}
		targetCharacterObjectRecs = append(targetCharacterObjectRecs, &dungeonCharacterObjectRec)
	}

	return rec, targetCharacterObjectRecs, nil
}

func (m *Model) createActionTargetMonsterRecs(actionID, locationInstanceID, monsterInstanceID string) (*record.ActionMonster, []*record.ActionMonsterObject, error) {
	l := m.Logger("createActionTargetMonsterRecs")

	l.Info("Creating action ID >%s< monster instance ID >%s< records", actionID, monsterInstanceID)

	targetMonsterInstanceViewRec, err := m.GetMonsterInstanceViewRec(monsterInstanceID)
	if err != nil {
		l.Warn("failed getting target monster instance view record >%v<", err)
		return nil, nil, err
	}

	rec := &record.ActionMonster{
		RecordType:          record.ActionMonsterRecordTypeTarget,
		ActionID:            actionID,
		LocationInstanceID:  locationInstanceID,
		MonsterInstanceID:   targetMonsterInstanceViewRec.ID,
		Name:                targetMonsterInstanceViewRec.Name,
		Strength:            targetMonsterInstanceViewRec.Strength,
		Dexterity:           targetMonsterInstanceViewRec.Dexterity,
		Intelligence:        targetMonsterInstanceViewRec.Intelligence,
		CurrentStrength:     targetMonsterInstanceViewRec.CurrentStrength,
		CurrentDexterity:    targetMonsterInstanceViewRec.CurrentDexterity,
		CurrentIntelligence: targetMonsterInstanceViewRec.CurrentIntelligence,
		Health:              targetMonsterInstanceViewRec.Health,
		Fatigue:             targetMonsterInstanceViewRec.Fatigue,
		CurrentHealth:       targetMonsterInstanceViewRec.CurrentHealth,
		CurrentFatigue:      targetMonsterInstanceViewRec.CurrentFatigue,
	}

	err = m.CreateActionMonsterRec(rec)
	if err != nil {
		l.Warn("failed creating target action monster record >%v<", err)
		return nil, nil, err
	}

	// Create action monster object records
	objectInstanceViewRecs, err := m.GetMonsterInstanceEquippedObjectInstanceViewRecs(targetMonsterInstanceViewRec.ID)
	if err != nil {
		l.Warn("failed getting target monster object records >%v<", err)
		return nil, nil, err
	}

	l.Info("Adding >%d< target monster object records", len(objectInstanceViewRecs))

	targetMonsterObjectRecs := []*record.ActionMonsterObject{}
	for _, objectInstanceViewRec := range objectInstanceViewRecs {
		l.Info("Adding target monster object record >%v<", objectInstanceViewRecs)
		dungeonMonsterObjectRec := record.ActionMonsterObject{
			ActionMonsterID:  rec.ID,
			ObjectInstanceID: objectInstanceViewRec.ID,
			Name:             objectInstanceViewRec.Name,
			IsEquipped:       objectInstanceViewRec.IsEquipped,
			IsStashed:        objectInstanceViewRec.IsStashed,
		}
		err := m.CreateActionMonsterObjectRec(&dungeonMonsterObjectRec)
		if err != nil {
			l.Warn("failed creating source action monster object record >%v<", err)
			return nil, nil, err
		}
		targetMonsterObjectRecs = append(targetMonsterObjectRecs, &dungeonMonsterObjectRec)
	}

	return rec, targetMonsterObjectRecs, nil
}

func (m *Model) createActionObjectRec(actionID, locationInstanceID, objectInstanceID, recordType string) (*record.ActionObject, error) {
	l := m.Logger("createActionObjectRec")

	targetObjectInstanceViewRec, err := m.GetObjectInstanceViewRec(objectInstanceID)
	if err != nil {
		l.Warn("failed getting target object instance view record >%v<", err)
		return nil, err
	}

	rec := &record.ActionObject{
		RecordType:         recordType,
		ActionID:           actionID,
		LocationInstanceID: locationInstanceID,
		ObjectInstanceID:   targetObjectInstanceViewRec.ID,
		Name:               targetObjectInstanceViewRec.Name,
		Description:        targetObjectInstanceViewRec.Description,
		IsStashed:          targetObjectInstanceViewRec.IsStashed,
		IsEquipped:         targetObjectInstanceViewRec.IsEquipped,
	}

	err = m.CreateActionObjectRec(rec)
	if err != nil {
		l.Warn("failed creating target action object record >%v<", err)
		return nil, err
	}

	return rec, nil
}
