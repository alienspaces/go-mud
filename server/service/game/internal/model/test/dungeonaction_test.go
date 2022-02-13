package test

// NOTE: model tests are run is the public space to avoid cyclic dependencies

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/go-mud/server/service/game/internal/dependencies"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/harness"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/model"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

func TestProcessDungeonCharacterAction(t *testing.T) {

	// harness
	config := harness.DefaultDataConfig

	c, l, s, m, err := dependencies.Default()
	require.NoError(t, err, "NewTesting returns without error")

	th, err := harness.NewTesting(c, l, s, m, config)
	require.NoError(t, err, "NewTesting returns without error")

	// harness commit data
	th.CommitData = true

	// TODO: Add current location checks

	tests := []struct {
		name                         string
		dungeonID                    func(data harness.Data) string
		dungeonCharacterID           func(data harness.Data) string
		sentence                     string
		expectDungeonActionRecordSet func(data harness.Data) *model.DungeonActionRecordSet
		expectError                  bool
	}{
		{
			name: "Look current location",
			dungeonID: func(data harness.Data) string {
				return data.DungeonRecs[0].ID
			},
			dungeonCharacterID: func(data harness.Data) string {
				return data.DungeonCharacterRecs[0].ID
			},
			sentence: "look",
			expectDungeonActionRecordSet: func(data harness.Data) *model.DungeonActionRecordSet {
				return &model.DungeonActionRecordSet{
					ActionRec: &record.DungeonAction{
						DungeonID:          data.DungeonRecs[0].ID,
						DungeonLocationID:  data.DungeonLocationRecs[0].ID,
						DungeonCharacterID: sql.NullString{String: data.DungeonCharacterRecs[0].ID, Valid: true},
					},
					ActionCharacterRec: &record.DungeonActionCharacter{
						Name:                data.DungeonCharacterRecs[0].Name,
						Strength:            data.DungeonCharacterRecs[0].Strength,
						Dexterity:           data.DungeonCharacterRecs[0].Dexterity,
						Intelligence:        data.DungeonCharacterRecs[0].Intelligence,
						CurrentStrength:     data.DungeonCharacterRecs[0].CurrentStrength,
						CurrentDexterity:    data.DungeonCharacterRecs[0].CurrentDexterity,
						CurrentIntelligence: data.DungeonCharacterRecs[0].CurrentIntelligence,
						Health:              data.DungeonCharacterRecs[0].Health,
						Fatigue:             data.DungeonCharacterRecs[0].Fatigue,
					},
					ActionCharacterObjectRecs: []*record.DungeonActionCharacterObject{
						{
							DungeonCharacterID: data.DungeonCharacterRecs[0].ID,
							DungeonObjectID:    data.DungeonObjectRecs[2].ID,
							Name:               data.DungeonObjectRecs[2].Name,
							IsStashed:          data.DungeonObjectRecs[2].IsStashed,
							IsEquipped:         data.DungeonObjectRecs[2].IsEquipped,
						},
						{
							DungeonCharacterID: data.DungeonCharacterRecs[0].ID,
							DungeonObjectID:    data.DungeonObjectRecs[3].ID,
							Name:               data.DungeonObjectRecs[3].Name,
							IsStashed:          data.DungeonObjectRecs[3].IsStashed,
							IsEquipped:         data.DungeonObjectRecs[3].IsEquipped,
						},
					},
					TargetLocation: &model.DungeonActionLocationRecordSet{
						LocationRec: &record.DungeonLocation{
							Name: data.DungeonLocationRecs[0].Name,
						},
					},
				}
			},
			expectError: false,
		},
		{
			name: "Look valid direction",
			dungeonID: func(data harness.Data) string {
				return data.DungeonRecs[0].ID
			},
			dungeonCharacterID: func(data harness.Data) string {
				return data.DungeonCharacterRecs[0].ID
			},
			sentence: "look north",
			expectDungeonActionRecordSet: func(data harness.Data) *model.DungeonActionRecordSet {
				return &model.DungeonActionRecordSet{
					ActionRec: &record.DungeonAction{
						DungeonID:          data.DungeonRecs[0].ID,
						DungeonLocationID:  data.DungeonLocationRecs[0].ID,
						DungeonCharacterID: sql.NullString{String: data.DungeonCharacterRecs[0].ID, Valid: true},
					},
					ActionCharacterRec: &record.DungeonActionCharacter{
						Name:                data.DungeonCharacterRecs[0].Name,
						Strength:            data.DungeonCharacterRecs[0].Strength,
						Dexterity:           data.DungeonCharacterRecs[0].Dexterity,
						Intelligence:        data.DungeonCharacterRecs[0].Intelligence,
						CurrentStrength:     data.DungeonCharacterRecs[0].CurrentStrength,
						CurrentDexterity:    data.DungeonCharacterRecs[0].CurrentDexterity,
						CurrentIntelligence: data.DungeonCharacterRecs[0].CurrentIntelligence,
						Health:              data.DungeonCharacterRecs[0].Health,
						Fatigue:             data.DungeonCharacterRecs[0].Fatigue,
					},
					ActionCharacterObjectRecs: []*record.DungeonActionCharacterObject{
						{
							DungeonCharacterID: data.DungeonCharacterRecs[0].ID,
							DungeonObjectID:    data.DungeonObjectRecs[2].ID,
							Name:               data.DungeonObjectRecs[2].Name,
							IsStashed:          data.DungeonObjectRecs[2].IsStashed,
							IsEquipped:         data.DungeonObjectRecs[2].IsEquipped,
						},
						{
							DungeonCharacterID: data.DungeonCharacterRecs[0].ID,
							DungeonObjectID:    data.DungeonObjectRecs[3].ID,
							Name:               data.DungeonObjectRecs[3].Name,
							IsStashed:          data.DungeonObjectRecs[3].IsStashed,
							IsEquipped:         data.DungeonObjectRecs[3].IsEquipped,
						},
					},
					TargetLocation: &model.DungeonActionLocationRecordSet{
						LocationRec: &record.DungeonLocation{
							Name: data.DungeonLocationRecs[1].Name,
						},
					},
				}
			},
			expectError: false,
		},
		{
			name: "Look valid item",
			dungeonID: func(data harness.Data) string {
				return data.DungeonRecs[0].ID
			},
			dungeonCharacterID: func(data harness.Data) string {
				return data.DungeonCharacterRecs[0].ID
			},
			sentence: "look rusted sword",
			expectDungeonActionRecordSet: func(data harness.Data) *model.DungeonActionRecordSet {
				return &model.DungeonActionRecordSet{
					ActionRec: &record.DungeonAction{
						DungeonID:          data.DungeonRecs[0].ID,
						DungeonLocationID:  data.DungeonLocationRecs[0].ID,
						DungeonCharacterID: sql.NullString{String: data.DungeonCharacterRecs[0].ID, Valid: true},
					},
					ActionCharacterRec: &record.DungeonActionCharacter{
						Name:                data.DungeonCharacterRecs[0].Name,
						Strength:            data.DungeonCharacterRecs[0].Strength,
						Dexterity:           data.DungeonCharacterRecs[0].Dexterity,
						Intelligence:        data.DungeonCharacterRecs[0].Intelligence,
						CurrentStrength:     data.DungeonCharacterRecs[0].CurrentStrength,
						CurrentDexterity:    data.DungeonCharacterRecs[0].CurrentDexterity,
						CurrentIntelligence: data.DungeonCharacterRecs[0].CurrentIntelligence,
						Health:              data.DungeonCharacterRecs[0].Health,
						Fatigue:             data.DungeonCharacterRecs[0].Fatigue,
					},
					ActionCharacterObjectRecs: []*record.DungeonActionCharacterObject{
						{
							DungeonCharacterID: data.DungeonCharacterRecs[0].ID,
							DungeonObjectID:    data.DungeonObjectRecs[2].ID,
							Name:               data.DungeonObjectRecs[2].Name,
							IsStashed:          data.DungeonObjectRecs[2].IsStashed,
							IsEquipped:         data.DungeonObjectRecs[2].IsEquipped,
						},
						{
							DungeonCharacterID: data.DungeonCharacterRecs[0].ID,
							DungeonObjectID:    data.DungeonObjectRecs[3].ID,
							Name:               data.DungeonObjectRecs[3].Name,
							IsStashed:          data.DungeonObjectRecs[3].IsStashed,
							IsEquipped:         data.DungeonObjectRecs[3].IsEquipped,
						},
					},
					TargetActionObjectRec: &record.DungeonActionObject{
						Name:        data.DungeonObjectRecs[0].Name,
						Description: data.DungeonObjectRecs[0].Description,
						IsStashed:   data.DungeonObjectRecs[0].IsStashed,
						IsEquipped:  data.DungeonObjectRecs[0].IsEquipped,
					},
				}
			},
			expectError: false,
		},
		{
			name: "Look valid monster",
			dungeonID: func(data harness.Data) string {
				return data.DungeonRecs[0].ID
			},
			dungeonCharacterID: func(data harness.Data) string {
				return data.DungeonCharacterRecs[0].ID
			},
			sentence: "look grumpy dwarf",
			expectDungeonActionRecordSet: func(data harness.Data) *model.DungeonActionRecordSet {
				return &model.DungeonActionRecordSet{
					ActionRec: &record.DungeonAction{
						DungeonID:          data.DungeonRecs[0].ID,
						DungeonLocationID:  data.DungeonLocationRecs[0].ID,
						DungeonCharacterID: sql.NullString{String: data.DungeonCharacterRecs[0].ID, Valid: true},
					},
					ActionCharacterRec: &record.DungeonActionCharacter{
						Name:                data.DungeonCharacterRecs[0].Name,
						Strength:            data.DungeonCharacterRecs[0].Strength,
						Dexterity:           data.DungeonCharacterRecs[0].Dexterity,
						Intelligence:        data.DungeonCharacterRecs[0].Intelligence,
						CurrentStrength:     data.DungeonCharacterRecs[0].CurrentStrength,
						CurrentDexterity:    data.DungeonCharacterRecs[0].CurrentDexterity,
						CurrentIntelligence: data.DungeonCharacterRecs[0].CurrentIntelligence,
						Health:              data.DungeonCharacterRecs[0].Health,
						Fatigue:             data.DungeonCharacterRecs[0].Fatigue,
					},
					ActionCharacterObjectRecs: []*record.DungeonActionCharacterObject{
						{
							DungeonCharacterID: data.DungeonCharacterRecs[0].ID,
							DungeonObjectID:    data.DungeonObjectRecs[2].ID,
							Name:               data.DungeonObjectRecs[2].Name,
							IsStashed:          data.DungeonObjectRecs[2].IsStashed,
							IsEquipped:         data.DungeonObjectRecs[2].IsEquipped,
						},
						{
							DungeonCharacterID: data.DungeonCharacterRecs[0].ID,
							DungeonObjectID:    data.DungeonObjectRecs[3].ID,
							Name:               data.DungeonObjectRecs[3].Name,
							IsStashed:          data.DungeonObjectRecs[3].IsStashed,
							IsEquipped:         data.DungeonObjectRecs[3].IsEquipped,
						},
					},
					TargetActionMonsterRec: &record.DungeonActionMonster{
						Name:                data.DungeonMonsterRecs[0].Name,
						Strength:            data.DungeonMonsterRecs[0].Strength,
						Dexterity:           data.DungeonMonsterRecs[0].Dexterity,
						Intelligence:        data.DungeonMonsterRecs[0].Intelligence,
						CurrentStrength:     data.DungeonMonsterRecs[0].CurrentStrength,
						CurrentDexterity:    data.DungeonMonsterRecs[0].CurrentDexterity,
						CurrentIntelligence: data.DungeonMonsterRecs[0].CurrentIntelligence,
						Health:              data.DungeonMonsterRecs[0].Health,
						Fatigue:             data.DungeonMonsterRecs[0].Fatigue,
					},
					TargetActionMonsterObjectRecs: []*record.DungeonActionMonsterObject{
						{
							DungeonMonsterID: data.DungeonMonsterRecs[0].ID,
							DungeonObjectID:  data.DungeonObjectRecs[4].ID,
							Name:             data.DungeonObjectRecs[4].Name,
							IsStashed:        data.DungeonObjectRecs[4].IsStashed,
							IsEquipped:       data.DungeonObjectRecs[4].IsEquipped,
						},
					},
				}
			},
			expectError: false,
		},
		{
			name: "Look valid character",
			dungeonID: func(data harness.Data) string {
				return data.DungeonRecs[0].ID
			},
			dungeonCharacterID: func(data harness.Data) string {
				return data.DungeonCharacterRecs[0].ID
			},
			sentence: "look barricade",
			expectDungeonActionRecordSet: func(data harness.Data) *model.DungeonActionRecordSet {
				return &model.DungeonActionRecordSet{
					ActionRec: &record.DungeonAction{
						DungeonID:          data.DungeonRecs[0].ID,
						DungeonLocationID:  data.DungeonLocationRecs[0].ID,
						DungeonCharacterID: sql.NullString{String: data.DungeonCharacterRecs[0].ID, Valid: true},
					},
					ActionCharacterRec: &record.DungeonActionCharacter{
						Name:                data.DungeonCharacterRecs[0].Name,
						Strength:            data.DungeonCharacterRecs[0].Strength,
						Dexterity:           data.DungeonCharacterRecs[0].Dexterity,
						Intelligence:        data.DungeonCharacterRecs[0].Intelligence,
						CurrentStrength:     data.DungeonCharacterRecs[0].CurrentStrength,
						CurrentDexterity:    data.DungeonCharacterRecs[0].CurrentDexterity,
						CurrentIntelligence: data.DungeonCharacterRecs[0].CurrentIntelligence,
						Health:              data.DungeonCharacterRecs[0].Health,
						Fatigue:             data.DungeonCharacterRecs[0].Fatigue,
					},
					ActionCharacterObjectRecs: []*record.DungeonActionCharacterObject{
						{
							DungeonCharacterID: data.DungeonCharacterRecs[0].ID,
							DungeonObjectID:    data.DungeonObjectRecs[2].ID,
							Name:               data.DungeonObjectRecs[2].Name,
							IsStashed:          data.DungeonObjectRecs[2].IsStashed,
							IsEquipped:         data.DungeonObjectRecs[2].IsEquipped,
						},
						{
							DungeonCharacterID: data.DungeonCharacterRecs[0].ID,
							DungeonObjectID:    data.DungeonObjectRecs[3].ID,
							Name:               data.DungeonObjectRecs[3].Name,
							IsStashed:          data.DungeonObjectRecs[3].IsStashed,
							IsEquipped:         data.DungeonObjectRecs[3].IsEquipped,
						},
					},
					TargetActionCharacterRec: &record.DungeonActionCharacter{
						Name:                data.DungeonCharacterRecs[0].Name,
						Strength:            data.DungeonCharacterRecs[0].Strength,
						Dexterity:           data.DungeonCharacterRecs[0].Dexterity,
						Intelligence:        data.DungeonCharacterRecs[0].Intelligence,
						CurrentStrength:     data.DungeonCharacterRecs[0].CurrentStrength,
						CurrentDexterity:    data.DungeonCharacterRecs[0].CurrentDexterity,
						CurrentIntelligence: data.DungeonCharacterRecs[0].CurrentIntelligence,
						Health:              data.DungeonCharacterRecs[0].Health,
						Fatigue:             data.DungeonCharacterRecs[0].Fatigue,
					},
					TargetActionCharacterObjectRecs: []*record.DungeonActionCharacterObject{
						{
							DungeonCharacterID: data.DungeonCharacterRecs[0].ID,
							DungeonObjectID:    data.DungeonObjectRecs[2].ID,
							Name:               data.DungeonObjectRecs[2].Name,
							IsStashed:          data.DungeonObjectRecs[2].IsStashed,
							IsEquipped:         data.DungeonObjectRecs[2].IsEquipped,
						},
					},
				}
			},
			expectError: false,
		},
		{
			name: "Stash valid location item",
			dungeonID: func(data harness.Data) string {
				return data.DungeonRecs[0].ID
			},
			dungeonCharacterID: func(data harness.Data) string {
				return data.DungeonCharacterRecs[0].ID
			},
			sentence: "stash rusted sword",
			expectDungeonActionRecordSet: func(data harness.Data) *model.DungeonActionRecordSet {
				return &model.DungeonActionRecordSet{
					ActionRec: &record.DungeonAction{
						DungeonID:          data.DungeonRecs[0].ID,
						DungeonLocationID:  data.DungeonLocationRecs[0].ID,
						DungeonCharacterID: sql.NullString{String: data.DungeonCharacterRecs[0].ID, Valid: true},
					},
					ActionCharacterRec: &record.DungeonActionCharacter{
						Name:                data.DungeonCharacterRecs[0].Name,
						Strength:            data.DungeonCharacterRecs[0].Strength,
						Dexterity:           data.DungeonCharacterRecs[0].Dexterity,
						Intelligence:        data.DungeonCharacterRecs[0].Intelligence,
						CurrentStrength:     data.DungeonCharacterRecs[0].CurrentStrength,
						CurrentDexterity:    data.DungeonCharacterRecs[0].CurrentDexterity,
						CurrentIntelligence: data.DungeonCharacterRecs[0].CurrentIntelligence,
						Health:              data.DungeonCharacterRecs[0].Health,
						Fatigue:             data.DungeonCharacterRecs[0].Fatigue,
					},
					ActionCharacterObjectRecs: []*record.DungeonActionCharacterObject{
						{
							DungeonCharacterID: data.DungeonCharacterRecs[0].ID,
							DungeonObjectID:    data.DungeonObjectRecs[0].ID,
							Name:               data.DungeonObjectRecs[0].Name,
							IsStashed:          true,
							IsEquipped:         false,
						},
						{
							DungeonCharacterID: data.DungeonCharacterRecs[0].ID,
							DungeonObjectID:    data.DungeonObjectRecs[2].ID,
							Name:               data.DungeonObjectRecs[2].Name,
							IsStashed:          data.DungeonObjectRecs[2].IsStashed,
							IsEquipped:         data.DungeonObjectRecs[2].IsEquipped,
						},
						{
							DungeonCharacterID: data.DungeonCharacterRecs[0].ID,
							DungeonObjectID:    data.DungeonObjectRecs[3].ID,
							Name:               data.DungeonObjectRecs[3].Name,
							IsStashed:          data.DungeonObjectRecs[3].IsStashed,
							IsEquipped:         data.DungeonObjectRecs[3].IsEquipped,
						},
					},
					TargetActionObjectRec: &record.DungeonActionObject{
						Name:        data.DungeonObjectRecs[0].Name,
						Description: data.DungeonObjectRecs[0].Description,
						IsStashed:   true,
						IsEquipped:  false,
					},
					StashedActionObjectRec: &record.DungeonActionObject{
						Name:        data.DungeonObjectRecs[0].Name,
						Description: data.DungeonObjectRecs[0].Description,
						IsStashed:   true,
						IsEquipped:  false,
					},
				}
			},
			expectError: false,
		},
		{
			name: "Stash valid equipped item",
			dungeonID: func(data harness.Data) string {
				return data.DungeonRecs[0].ID
			},
			dungeonCharacterID: func(data harness.Data) string {
				return data.DungeonCharacterRecs[0].ID
			},
			sentence: "stash dull bronze ring",
			expectDungeonActionRecordSet: func(data harness.Data) *model.DungeonActionRecordSet {
				return &model.DungeonActionRecordSet{
					ActionRec: &record.DungeonAction{
						DungeonID:          data.DungeonRecs[0].ID,
						DungeonLocationID:  data.DungeonLocationRecs[0].ID,
						DungeonCharacterID: sql.NullString{String: data.DungeonCharacterRecs[0].ID, Valid: true},
					},
					ActionCharacterRec: &record.DungeonActionCharacter{
						Name:                data.DungeonCharacterRecs[0].Name,
						Strength:            data.DungeonCharacterRecs[0].Strength,
						Dexterity:           data.DungeonCharacterRecs[0].Dexterity,
						Intelligence:        data.DungeonCharacterRecs[0].Intelligence,
						CurrentStrength:     data.DungeonCharacterRecs[0].CurrentStrength,
						CurrentDexterity:    data.DungeonCharacterRecs[0].CurrentDexterity,
						CurrentIntelligence: data.DungeonCharacterRecs[0].CurrentIntelligence,
						Health:              data.DungeonCharacterRecs[0].Health,
						Fatigue:             data.DungeonCharacterRecs[0].Fatigue,
					},
					ActionCharacterObjectRecs: []*record.DungeonActionCharacterObject{
						{
							DungeonCharacterID: data.DungeonCharacterRecs[0].ID,
							DungeonObjectID:    data.DungeonObjectRecs[2].ID,
							Name:               data.DungeonObjectRecs[2].Name,
							IsStashed:          true,
							IsEquipped:         false,
						},
						{
							DungeonCharacterID: data.DungeonCharacterRecs[0].ID,
							DungeonObjectID:    data.DungeonObjectRecs[3].ID,
							Name:               data.DungeonObjectRecs[3].Name,
							IsStashed:          data.DungeonObjectRecs[3].IsStashed,
							IsEquipped:         data.DungeonObjectRecs[3].IsEquipped,
						},
					},
					TargetActionObjectRec: &record.DungeonActionObject{
						Name:        data.DungeonObjectRecs[2].Name,
						Description: data.DungeonObjectRecs[2].Description,
						IsStashed:   true,
						IsEquipped:  false,
					},
					StashedActionObjectRec: &record.DungeonActionObject{
						Name:        data.DungeonObjectRecs[2].Name,
						Description: data.DungeonObjectRecs[2].Description,
						IsStashed:   true,
						IsEquipped:  false,
					},
				}
			},
			expectError: false,
		},
		{
			name: "Equip valid location item",
			dungeonID: func(data harness.Data) string {
				return data.DungeonRecs[0].ID
			},
			dungeonCharacterID: func(data harness.Data) string {
				return data.DungeonCharacterRecs[0].ID
			},
			sentence: "equip rusted sword",
			expectDungeonActionRecordSet: func(data harness.Data) *model.DungeonActionRecordSet {
				return &model.DungeonActionRecordSet{
					ActionRec: &record.DungeonAction{
						DungeonID:          data.DungeonRecs[0].ID,
						DungeonLocationID:  data.DungeonLocationRecs[0].ID,
						DungeonCharacterID: sql.NullString{String: data.DungeonCharacterRecs[0].ID, Valid: true},
					},
					ActionCharacterRec: &record.DungeonActionCharacter{
						Name:                data.DungeonCharacterRecs[0].Name,
						Strength:            data.DungeonCharacterRecs[0].Strength,
						Dexterity:           data.DungeonCharacterRecs[0].Dexterity,
						Intelligence:        data.DungeonCharacterRecs[0].Intelligence,
						CurrentStrength:     data.DungeonCharacterRecs[0].CurrentStrength,
						CurrentDexterity:    data.DungeonCharacterRecs[0].CurrentDexterity,
						CurrentIntelligence: data.DungeonCharacterRecs[0].CurrentIntelligence,
						Health:              data.DungeonCharacterRecs[0].Health,
						Fatigue:             data.DungeonCharacterRecs[0].Fatigue,
					},
					ActionCharacterObjectRecs: []*record.DungeonActionCharacterObject{
						{
							DungeonCharacterID: data.DungeonCharacterRecs[0].ID,
							DungeonObjectID:    data.DungeonObjectRecs[0].ID,
							Name:               data.DungeonObjectRecs[0].Name,
							IsStashed:          false,
							IsEquipped:         true,
						},
						{
							DungeonCharacterID: data.DungeonCharacterRecs[0].ID,
							DungeonObjectID:    data.DungeonObjectRecs[2].ID,
							Name:               data.DungeonObjectRecs[2].Name,
							IsStashed:          data.DungeonObjectRecs[2].IsStashed,
							IsEquipped:         data.DungeonObjectRecs[2].IsEquipped,
						},
						{
							DungeonCharacterID: data.DungeonCharacterRecs[0].ID,
							DungeonObjectID:    data.DungeonObjectRecs[3].ID,
							Name:               data.DungeonObjectRecs[3].Name,
							IsStashed:          data.DungeonObjectRecs[3].IsStashed,
							IsEquipped:         data.DungeonObjectRecs[3].IsEquipped,
						},
					},
					TargetActionObjectRec: &record.DungeonActionObject{
						Name:        data.DungeonObjectRecs[0].Name,
						Description: data.DungeonObjectRecs[0].Description,
						IsStashed:   false,
						IsEquipped:  true,
					},
					EquippedActionObjectRec: &record.DungeonActionObject{
						Name:        data.DungeonObjectRecs[0].Name,
						Description: data.DungeonObjectRecs[0].Description,
						IsStashed:   false,
						IsEquipped:  true,
					},
				}
			},
			expectError: false,
		},
		{
			name: "Equip valid stashed item",
			dungeonID: func(data harness.Data) string {
				return data.DungeonRecs[0].ID
			},
			dungeonCharacterID: func(data harness.Data) string {
				return data.DungeonCharacterRecs[0].ID
			},
			sentence: "equip blood stained pouch",
			expectDungeonActionRecordSet: func(data harness.Data) *model.DungeonActionRecordSet {
				return &model.DungeonActionRecordSet{
					ActionRec: &record.DungeonAction{
						DungeonID:          data.DungeonRecs[0].ID,
						DungeonLocationID:  data.DungeonLocationRecs[0].ID,
						DungeonCharacterID: sql.NullString{String: data.DungeonCharacterRecs[0].ID, Valid: true},
					},
					ActionCharacterRec: &record.DungeonActionCharacter{
						Name:                data.DungeonCharacterRecs[0].Name,
						Strength:            data.DungeonCharacterRecs[0].Strength,
						Dexterity:           data.DungeonCharacterRecs[0].Dexterity,
						Intelligence:        data.DungeonCharacterRecs[0].Intelligence,
						CurrentStrength:     data.DungeonCharacterRecs[0].CurrentStrength,
						CurrentDexterity:    data.DungeonCharacterRecs[0].CurrentDexterity,
						CurrentIntelligence: data.DungeonCharacterRecs[0].CurrentIntelligence,
						Health:              data.DungeonCharacterRecs[0].Health,
						Fatigue:             data.DungeonCharacterRecs[0].Fatigue,
					},
					ActionCharacterObjectRecs: []*record.DungeonActionCharacterObject{
						{
							DungeonCharacterID: data.DungeonCharacterRecs[0].ID,
							DungeonObjectID:    data.DungeonObjectRecs[2].ID,
							Name:               data.DungeonObjectRecs[2].Name,
							IsStashed:          data.DungeonObjectRecs[2].IsStashed,
							IsEquipped:         data.DungeonObjectRecs[2].IsEquipped,
						},
						{
							DungeonCharacterID: data.DungeonCharacterRecs[0].ID,
							DungeonObjectID:    data.DungeonObjectRecs[3].ID,
							Name:               data.DungeonObjectRecs[3].Name,
							IsStashed:          false,
							IsEquipped:         true,
						},
					},
					TargetActionObjectRec: &record.DungeonActionObject{
						Name:        data.DungeonObjectRecs[3].Name,
						Description: data.DungeonObjectRecs[3].Description,
						IsStashed:   false,
						IsEquipped:  true,
					},
					EquippedActionObjectRec: &record.DungeonActionObject{
						Name:        data.DungeonObjectRecs[3].Name,
						Description: data.DungeonObjectRecs[3].Description,
						IsStashed:   false,
						IsEquipped:  true,
					},
				}
			},
			expectError: false,
		},
		{
			name: "Drop valid equipped item",
			dungeonID: func(data harness.Data) string {
				return data.DungeonRecs[0].ID
			},
			dungeonCharacterID: func(data harness.Data) string {
				return data.DungeonCharacterRecs[0].ID
			},
			sentence: "drop dull bronze ring",
			expectDungeonActionRecordSet: func(data harness.Data) *model.DungeonActionRecordSet {
				return &model.DungeonActionRecordSet{
					ActionRec: &record.DungeonAction{
						DungeonID:          data.DungeonRecs[0].ID,
						DungeonLocationID:  data.DungeonLocationRecs[0].ID,
						DungeonCharacterID: sql.NullString{String: data.DungeonCharacterRecs[0].ID, Valid: true},
					},
					ActionCharacterRec: &record.DungeonActionCharacter{
						Name:                data.DungeonCharacterRecs[0].Name,
						Strength:            data.DungeonCharacterRecs[0].Strength,
						Dexterity:           data.DungeonCharacterRecs[0].Dexterity,
						Intelligence:        data.DungeonCharacterRecs[0].Intelligence,
						CurrentStrength:     data.DungeonCharacterRecs[0].CurrentStrength,
						CurrentDexterity:    data.DungeonCharacterRecs[0].CurrentDexterity,
						CurrentIntelligence: data.DungeonCharacterRecs[0].CurrentIntelligence,
						Health:              data.DungeonCharacterRecs[0].Health,
						Fatigue:             data.DungeonCharacterRecs[0].Fatigue,
					},
					ActionCharacterObjectRecs: []*record.DungeonActionCharacterObject{
						{
							DungeonCharacterID: data.DungeonCharacterRecs[0].ID,
							DungeonObjectID:    data.DungeonObjectRecs[3].ID,
							Name:               data.DungeonObjectRecs[3].Name,
							IsStashed:          data.DungeonObjectRecs[3].IsStashed,
							IsEquipped:         data.DungeonObjectRecs[3].IsEquipped,
						},
					},
					TargetActionObjectRec: &record.DungeonActionObject{
						Name:        data.DungeonObjectRecs[2].Name,
						Description: data.DungeonObjectRecs[2].Description,
						IsStashed:   false,
						IsEquipped:  false,
					},
					DroppedActionObjectRec: &record.DungeonActionObject{
						Name:        data.DungeonObjectRecs[2].Name,
						Description: data.DungeonObjectRecs[2].Description,
						IsStashed:   false,
						IsEquipped:  false,
					},
				}
			},
			expectError: false,
		},
		{
			name: "Drop valid stashed item",
			dungeonID: func(data harness.Data) string {
				return data.DungeonRecs[0].ID
			},
			dungeonCharacterID: func(data harness.Data) string {
				return data.DungeonCharacterRecs[0].ID
			},
			sentence: "drop blood stained pouch",
			expectDungeonActionRecordSet: func(data harness.Data) *model.DungeonActionRecordSet {
				return &model.DungeonActionRecordSet{
					ActionRec: &record.DungeonAction{
						DungeonID:          data.DungeonRecs[0].ID,
						DungeonLocationID:  data.DungeonLocationRecs[0].ID,
						DungeonCharacterID: sql.NullString{String: data.DungeonCharacterRecs[0].ID, Valid: true},
					},
					ActionCharacterRec: &record.DungeonActionCharacter{
						Name:                data.DungeonCharacterRecs[0].Name,
						Strength:            data.DungeonCharacterRecs[0].Strength,
						Dexterity:           data.DungeonCharacterRecs[0].Dexterity,
						Intelligence:        data.DungeonCharacterRecs[0].Intelligence,
						CurrentStrength:     data.DungeonCharacterRecs[0].CurrentStrength,
						CurrentDexterity:    data.DungeonCharacterRecs[0].CurrentDexterity,
						CurrentIntelligence: data.DungeonCharacterRecs[0].CurrentIntelligence,
						Health:              data.DungeonCharacterRecs[0].Health,
						Fatigue:             data.DungeonCharacterRecs[0].Fatigue,
					},
					ActionCharacterObjectRecs: []*record.DungeonActionCharacterObject{
						{
							DungeonCharacterID: data.DungeonCharacterRecs[0].ID,
							DungeonObjectID:    data.DungeonObjectRecs[2].ID,
							Name:               data.DungeonObjectRecs[2].Name,
							IsStashed:          data.DungeonObjectRecs[2].IsStashed,
							IsEquipped:         data.DungeonObjectRecs[2].IsEquipped,
						},
					},
					TargetActionObjectRec: &record.DungeonActionObject{
						Name:        data.DungeonObjectRecs[3].Name,
						Description: data.DungeonObjectRecs[3].Description,
						IsStashed:   false,
						IsEquipped:  false,
					},
					DroppedActionObjectRec: &record.DungeonActionObject{
						Name:        data.DungeonObjectRecs[3].Name,
						Description: data.DungeonObjectRecs[3].Description,
						IsStashed:   false,
						IsEquipped:  false,
					},
				}
			},
			expectError: false,
		},
	}

	for _, tc := range tests {

		t.Logf("Run test >%s< - >%s<", tc.name, tc.sentence)

		func() {

			// Test harness
			err = th.Setup()
			require.NoError(t, err, "Setup returns without error")
			defer func() {
				err = th.RollbackTx()
				require.NoError(t, err, "RollbackTx returns without error")
				err = th.Teardown()
				require.NoError(t, err, "Teardown returns without error")
			}()

			// init tx
			err = th.InitTx(nil)
			require.NoError(t, err, "InitTx returns without error")

			dungeonID := tc.dungeonID(th.Data)
			dungeonCharacterID := tc.dungeonCharacterID(th.Data)

			rslt, err := th.Model.(*model.Model).ProcessDungeonCharacterAction(dungeonID, dungeonCharacterID, tc.sentence)
			if tc.expectError == true {
				require.Error(t, err, "CreateDungeonObjectRec returns error")
				return
			}
			require.NoError(t, err, "ProcessDungeonCharacterAction returns without error")
			require.NotNil(t, rslt.ActionRec, "ProcessDungeonCharacterAction returns DungeonActionRecordSet with ActionRec")

			xrslt := tc.expectDungeonActionRecordSet(th.Data)
			if xrslt == nil {
				return
			}

			if xrslt.ActionRec != nil {
				require.NotNil(t, rslt.ActionRec, "ActionRec is not nil")
				require.Equal(t, xrslt.ActionRec.DungeonID, rslt.ActionRec.DungeonID, "ActionRec.DungeonID equals expected")
				require.Equal(t, xrslt.ActionRec.DungeonCharacterID, rslt.ActionRec.DungeonCharacterID, "ActionRec.DungeonCharacterID equals expected")
				require.Equal(t, xrslt.ActionRec.DungeonMonsterID, rslt.ActionRec.DungeonMonsterID, "ActionRec.DungeonMonsterID equals expected")
				require.Equal(t, xrslt.ActionRec.DungeonLocationID, rslt.ActionRec.DungeonLocationID, "ActionRec.DungeonLocationID equals expected")
			} else {
				require.Nil(t, rslt.ActionRec, "ActionRec is nil")
			}

			if xrslt.ActionCharacterRec != nil {
				require.NotNil(t, rslt.ActionCharacterRec, "ActionCharacterRec is not nil")
				require.Equal(t, xrslt.ActionCharacterRec.Name, rslt.ActionCharacterRec.Name, "ActionCharacterRec.Name equals expected")
				require.Equal(t, xrslt.ActionCharacterRec.Strength, rslt.ActionCharacterRec.Strength, "ActionCharacterRec.Strength equals expected")
				require.Equal(t, xrslt.ActionCharacterRec.Dexterity, rslt.ActionCharacterRec.Dexterity, "ActionCharacterRec.Dexterity equals expected")
				require.Equal(t, xrslt.ActionCharacterRec.Intelligence, rslt.ActionCharacterRec.Intelligence, "ActionCharacterRec.Intelligence equals expected")

				require.Equal(t, xrslt.ActionCharacterRec.CurrentStrength, rslt.ActionCharacterRec.CurrentStrength, "ActionCharacterRec.CurrentStrength equals expected")
				require.Equal(t, xrslt.ActionCharacterRec.CurrentDexterity, rslt.ActionCharacterRec.CurrentDexterity, "ActionCharacterRec.CurrentDexterity equals expected")
				require.Equal(t, xrslt.ActionCharacterRec.CurrentIntelligence, rslt.ActionCharacterRec.CurrentIntelligence, "ActionCharacterRec.CurrentIntelligence equals expected")

				require.Equal(t, xrslt.ActionCharacterRec.Health, rslt.ActionCharacterRec.Health, "ActionCharacterRec.Health equals expected")
				require.Equal(t, xrslt.ActionCharacterRec.Fatigue, rslt.ActionCharacterRec.Fatigue, "ActionCharacterRec.Fatigue equals expected")
			} else {
				require.Nil(t, rslt.ActionCharacterRec, "ActionCharacterRec is nil")
			}

			require.Equal(t, len(xrslt.ActionCharacterObjectRecs), len(rslt.ActionCharacterObjectRecs), "ActionCharacterObjectRecs count equals expected")
			if len(xrslt.ActionCharacterObjectRecs) > 0 {
				for _, xrObjectRec := range xrslt.ActionCharacterObjectRecs {
					found := false
					for _, rObjectRec := range rslt.ActionCharacterObjectRecs {
						if xrObjectRec.Name == rObjectRec.Name {
							found = true
							require.Equal(t, xrObjectRec.IsEquipped, rObjectRec.IsEquipped, "ActionCharacterObjectRec IsEquipped equals expected")
							require.Equal(t, xrObjectRec.IsStashed, rObjectRec.IsStashed, "ActionCharacterObjectRec IsStashed equals expected")
						}
					}
					require.True(t, found, fmt.Sprintf("ActionCharacterObjectRec >%s< found", xrObjectRec.Name))
				}
			}

			if xrslt.ActionMonsterRec != nil {
				require.NotNil(t, rslt.ActionMonsterRec, "ActionMonsterRec is not nil")
				require.Equal(t, xrslt.ActionMonsterRec.Name, rslt.ActionMonsterRec.Name, "ActionMonsterRec.Name equals expected")
				require.Equal(t, xrslt.ActionMonsterRec.Strength, rslt.ActionMonsterRec.Strength, "ActionMonsterRec.Strength equals expected")
				require.Equal(t, xrslt.ActionMonsterRec.Dexterity, rslt.ActionMonsterRec.Dexterity, "ActionMonsterRec.Dexterity equals expected")
				require.Equal(t, xrslt.ActionMonsterRec.Intelligence, rslt.ActionMonsterRec.Intelligence, "ActionMonsterRec.Intelligence equals expected")

				require.Equal(t, xrslt.ActionMonsterRec.CurrentStrength, rslt.ActionMonsterRec.CurrentStrength, "ActionMonsterRec.CurrentStrength equals expected")
				require.Equal(t, xrslt.ActionMonsterRec.CurrentDexterity, rslt.ActionMonsterRec.CurrentDexterity, "ActionMonsterRec.CurrentDexterity equals expected")
				require.Equal(t, xrslt.ActionMonsterRec.CurrentIntelligence, rslt.ActionMonsterRec.CurrentIntelligence, "ActionMonsterRec.CurrentIntelligence equals expected")

				require.Equal(t, xrslt.ActionMonsterRec.Health, rslt.ActionMonsterRec.Health, "ActionMonsterRec.Health equals expected")
				require.Equal(t, xrslt.ActionMonsterRec.Fatigue, rslt.ActionMonsterRec.Fatigue, "ActionMonsterRec.Fatigue equals expected")
			} else {
				require.Nil(t, rslt.ActionMonsterRec, "ActionMonsterRec is nil")
			}

			require.Equal(t, len(xrslt.ActionMonsterObjectRecs), len(rslt.ActionMonsterObjectRecs), "ActionMonsterObjectRecs count equals expected")
			if len(xrslt.ActionMonsterObjectRecs) > 0 {
				for _, xrObjectRec := range xrslt.ActionMonsterObjectRecs {
					found := false
					for _, rObjectRec := range rslt.ActionMonsterObjectRecs {
						if xrObjectRec.Name == rObjectRec.Name {
							found = true
							require.Equal(t, xrObjectRec.IsEquipped, rObjectRec.IsEquipped, "ActionMonsterObjectRec IsEquipped equals expected")
							require.Equal(t, xrObjectRec.IsStashed, rObjectRec.IsStashed, "ActionMonsterObjectRec IsStashed equals expected")
						}
					}
					require.True(t, found, fmt.Sprintf("ActionMonsterObjectRec >%s< found", xrObjectRec.Name))
				}
			}

			if xrslt.EquippedActionObjectRec != nil {
				require.NotNil(t, rslt.EquippedActionObjectRec, "EquippedActionObjectRec is not nil")
				require.Equal(t, xrslt.EquippedActionObjectRec.Name, rslt.EquippedActionObjectRec.Name, "EquippedActionObjectRec.Name equals expected")
				require.Equal(t, xrslt.EquippedActionObjectRec.Description, rslt.EquippedActionObjectRec.Description, "EquippedActionObjectRec.Description equals expected")
				require.Equal(t, xrslt.EquippedActionObjectRec.IsStashed, rslt.EquippedActionObjectRec.IsStashed, "EquippedActionObjectRec.IsStashed equals expected")
				require.Equal(t, xrslt.EquippedActionObjectRec.IsEquipped, rslt.EquippedActionObjectRec.IsEquipped, "EquippedActionObjectRec.IsEquipped equals expected")
			} else {
				require.Nil(t, rslt.EquippedActionObjectRec, "EquippedActionObjectRec is nil")
			}

			if xrslt.StashedActionObjectRec != nil {
				require.NotNil(t, rslt.StashedActionObjectRec, "StashedActionObjectRec is not nil")
				require.Equal(t, xrslt.StashedActionObjectRec.Name, rslt.StashedActionObjectRec.Name, "StashedActionObjectRec.Name equals expected")
				require.Equal(t, xrslt.StashedActionObjectRec.Description, rslt.StashedActionObjectRec.Description, "StashedActionObjectRec.Description equals expected")
				require.Equal(t, xrslt.StashedActionObjectRec.IsStashed, rslt.StashedActionObjectRec.IsStashed, "StashedActionObjectRec.IsStashed equals expected")
				require.Equal(t, xrslt.StashedActionObjectRec.IsEquipped, rslt.StashedActionObjectRec.IsEquipped, "StashedActionObjectRec.IsEquipped equals expected")
			} else {
				require.Nil(t, rslt.StashedActionObjectRec, "StashedActionObjectRec is nil")
			}

			if xrslt.DroppedActionObjectRec != nil {
				require.NotNil(t, rslt.DroppedActionObjectRec, "DroppedActionObjectRec is not nil")
				require.Equal(t, xrslt.DroppedActionObjectRec.Name, rslt.DroppedActionObjectRec.Name, "DroppedActionObjectRec.Name equals expected")
				require.Equal(t, xrslt.DroppedActionObjectRec.Description, rslt.DroppedActionObjectRec.Description, "DroppedActionObjectRec.Description equals expected")
				require.Equal(t, xrslt.DroppedActionObjectRec.IsStashed, rslt.DroppedActionObjectRec.IsStashed, "DroppedActionObjectRec.IsStashed equals expected")
				require.Equal(t, xrslt.DroppedActionObjectRec.IsEquipped, rslt.DroppedActionObjectRec.IsEquipped, "DroppedActionObjectRec.IsEquipped equals expected")
			} else {
				require.Nil(t, rslt.DroppedActionObjectRec, "DroppedActionObjectRec is nil")
			}

			if xrslt.TargetActionObjectRec != nil {
				require.NotNil(t, rslt.TargetActionObjectRec, "TargetActionObjectRec is not nil")
				require.Equal(t, xrslt.TargetActionObjectRec.Name, rslt.TargetActionObjectRec.Name, "TargetActionObjectRec.Name equals expected")
				require.Equal(t, xrslt.TargetActionObjectRec.Description, rslt.TargetActionObjectRec.Description, "TargetActionObjectRec.Description equals expected")
				require.Equal(t, xrslt.TargetActionObjectRec.IsStashed, rslt.TargetActionObjectRec.IsStashed, "TargetActionObjectRec.IsStashed equals expected")
				require.Equal(t, xrslt.TargetActionObjectRec.IsEquipped, rslt.TargetActionObjectRec.IsEquipped, "TargetActionObjectRec.IsEquipped equals expected")
			} else {
				require.Nil(t, rslt.TargetActionObjectRec, "TargetActionObjectRec is nil")
			}

			if xrslt.TargetActionCharacterRec != nil {
				require.NotNil(t, rslt.TargetActionCharacterRec, "TargetActionCharacterRec is not nil")
				require.Equal(t, xrslt.TargetActionCharacterRec.Name, rslt.TargetActionCharacterRec.Name, "TargetActionCharacterRec.Name equals expected")
				require.Equal(t, xrslt.TargetActionCharacterRec.Strength, rslt.TargetActionCharacterRec.Strength, "TargetActionCharacterRec.Strength equals expected")
				require.Equal(t, xrslt.TargetActionCharacterRec.Dexterity, rslt.TargetActionCharacterRec.Dexterity, "TargetActionCharacterRec.Dexterity equals expected")
				require.Equal(t, xrslt.TargetActionCharacterRec.Intelligence, rslt.TargetActionCharacterRec.Intelligence, "TargetActionCharacterRec.Intelligence equals expected")

				require.Equal(t, xrslt.TargetActionCharacterRec.CurrentStrength, rslt.TargetActionCharacterRec.CurrentStrength, "TargetActionCharacterRec.CurrentStrength equals expected")
				require.Equal(t, xrslt.TargetActionCharacterRec.CurrentDexterity, rslt.TargetActionCharacterRec.CurrentDexterity, "TargetActionCharacterRec.CurrentDexterity equals expected")
				require.Equal(t, xrslt.TargetActionCharacterRec.CurrentIntelligence, rslt.TargetActionCharacterRec.CurrentIntelligence, "TargetActionCharacterRec.CurrentIntelligence equals expected")

				require.Equal(t, xrslt.TargetActionCharacterRec.Health, rslt.TargetActionCharacterRec.Health, "TargetActionCharacterRec.Health equals expected")
				require.Equal(t, xrslt.TargetActionCharacterRec.Fatigue, rslt.TargetActionCharacterRec.Fatigue, "TargetActionCharacterRec.Fatigue equals expected")
			} else {
				require.Nil(t, rslt.TargetActionCharacterRec, "TargetActionCharacterRec is nil")
			}

			require.Equal(t, len(xrslt.TargetActionCharacterObjectRecs), len(rslt.TargetActionCharacterObjectRecs), "TargetActionCharacterObjectRecs count equals expected")
			if len(xrslt.TargetActionCharacterObjectRecs) > 0 {
				for _, xrObjectRec := range xrslt.TargetActionCharacterObjectRecs {
					found := false
					for _, rObjectRec := range rslt.TargetActionCharacterObjectRecs {
						if xrObjectRec.Name == rObjectRec.Name {
							found = true
							require.Equal(t, xrObjectRec.IsEquipped, rObjectRec.IsEquipped, "TargetActionCharacterObjectRec IsEquipped equals expected")
							require.Equal(t, xrObjectRec.IsStashed, rObjectRec.IsStashed, "TargetActionCharacterObjectRec IsStashed equals expected")
						}
					}
					require.True(t, found, fmt.Sprintf("TargetActionCharacterObjectRec >%s< found", xrObjectRec.Name))
				}
			}

			if xrslt.TargetActionMonsterRec != nil {
				require.NotNil(t, rslt.TargetActionMonsterRec, "TargetActionMonsterRec is not nil")
				require.Equal(t, xrslt.TargetActionMonsterRec.Name, rslt.TargetActionMonsterRec.Name, "TargetActionMonsterRec.Name equals expected")
				require.Equal(t, xrslt.TargetActionMonsterRec.Strength, rslt.TargetActionMonsterRec.Strength, "TargetActionMonsterRec.Strength equals expected")
				require.Equal(t, xrslt.TargetActionMonsterRec.Dexterity, rslt.TargetActionMonsterRec.Dexterity, "TargetActionMonsterRec.Dexterity equals expected")
				require.Equal(t, xrslt.TargetActionMonsterRec.Intelligence, rslt.TargetActionMonsterRec.Intelligence, "TargetActionMonsterRec.Intelligence equals expected")

				require.Equal(t, xrslt.TargetActionMonsterRec.CurrentStrength, rslt.TargetActionMonsterRec.CurrentStrength, "TargetActionMonsterRec.CurrentStrength equals expected")
				require.Equal(t, xrslt.TargetActionMonsterRec.CurrentDexterity, rslt.TargetActionMonsterRec.CurrentDexterity, "TargetActionMonsterRec.CurrentDexterity equals expected")
				require.Equal(t, xrslt.TargetActionMonsterRec.CurrentIntelligence, rslt.TargetActionMonsterRec.CurrentIntelligence, "TargetActionMonsterRec.CurrentIntelligence equals expected")

				require.Equal(t, xrslt.TargetActionMonsterRec.Health, rslt.TargetActionMonsterRec.Health, "TargetActionMonsterRec.Health equals expected")
				require.Equal(t, xrslt.TargetActionMonsterRec.Fatigue, rslt.TargetActionMonsterRec.Fatigue, "TargetActionMonsterRec.Fatigue equals expected")
			} else {
				require.Nil(t, rslt.TargetActionMonsterRec, "TargetActionMonsterRec is nil")
			}

			require.Equal(t, len(xrslt.TargetActionMonsterObjectRecs), len(rslt.TargetActionMonsterObjectRecs), "TargetActionMonsterObjectRecs count equals expected")
			if len(xrslt.TargetActionMonsterObjectRecs) > 0 {
				for _, xrObjectRec := range xrslt.TargetActionMonsterObjectRecs {
					found := false
					for _, rObjectRec := range rslt.TargetActionMonsterObjectRecs {
						if xrObjectRec.Name == rObjectRec.Name {
							found = true
							require.Equal(t, xrObjectRec.IsEquipped, rObjectRec.IsEquipped, "TargetActionMonsterObjectRec IsEquipped equals expected")
							require.Equal(t, xrObjectRec.IsStashed, rObjectRec.IsStashed, "TargetActionMonsterObjectRec IsStashed equals expected")
						}
					}
					require.True(t, found, fmt.Sprintf("TargetActionMonsterObjectRec >%s< found", xrObjectRec.Name))
				}
			}

			if xrslt.TargetLocation != nil {
				require.NotNil(t, rslt.TargetLocation, "TargetLocation is not nil")
				if xrslt.TargetLocation.LocationRec != nil {
					require.NotNil(t, rslt.TargetLocation.LocationRec, "TargetLocation.LocationRec is not nil")
					require.Equal(t, xrslt.TargetLocation.LocationRec.Name, rslt.TargetLocation.LocationRec.Name, "TargetLocation.LocationRec equals expected")
				}

				// TODO: Compare items, monsters and characters

			} else {
				require.Nil(t, rslt.TargetLocation, "TargetLocation is nil")
			}
		}()
	}
}
