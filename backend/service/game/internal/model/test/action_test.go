package test

// NOTE: model tests are run is the public space to avoid cyclic dependencies

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/go-mud/backend/core/null"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/dependencies"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/harness"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/model"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/record"
)

func TestConsecutiveProcessCharacterActions(t *testing.T) {

	c, l, s, err := dependencies.Default()
	require.NoError(t, err, "NewTesting returns without error")

	config := harness.DefaultDataConfig
	th, err := harness.NewTesting(c, l, s, config)
	require.NoError(t, err, "NewTesting returns without error")

	// harness commit data
	th.ShouldCommitData = true

	tests := []struct {
		name                string
		dungeonInstanceID   func(data harness.Data) string
		characterInstanceID func(data harness.Data) string
		sentence            func(data harness.Data) string
		incrementTurns      int
		expectError         error
	}{
		{
			name: "Too fast",
			dungeonInstanceID: func(data harness.Data) string {
				diRec, _ := data.GetDungeonInstanceRecByName(harness.DungeonNameCave)
				return diRec.ID
			},
			characterInstanceID: func(data harness.Data) string {
				ciRec, _ := data.GetCharacterInstanceRecByName(harness.CharacterNameBarricade)
				return ciRec.ID
			},
			sentence: func(data harness.Data) string {
				return "look"
			},
			incrementTurns: 0,
			expectError:    model.NewActionTooEarlyError(2, 2),
		},
		{
			name: "Too slow",
			dungeonInstanceID: func(data harness.Data) string {
				diRec, _ := data.GetDungeonInstanceRecByName(harness.DungeonNameCave)
				return diRec.ID
			},
			characterInstanceID: func(data harness.Data) string {
				ciRec, _ := data.GetCharacterInstanceRecByName(harness.CharacterNameBarricade)
				return ciRec.ID
			},
			sentence: func(data harness.Data) string {
				return "look"
			},
			incrementTurns: 2,
			expectError:    nil,
		},
		{
			name: "Just right",
			dungeonInstanceID: func(data harness.Data) string {
				diRec, _ := data.GetDungeonInstanceRecByName(harness.DungeonNameCave)
				return diRec.ID
			},
			characterInstanceID: func(data harness.Data) string {
				ciRec, _ := data.GetCharacterInstanceRecByName(harness.CharacterNameBarricade)
				return ciRec.ID
			},
			sentence: func(data harness.Data) string {
				return "look"
			},
			incrementTurns: 1,
			expectError:    nil,
		},
	}

	for _, tc := range tests {

		t.Logf("Run test >%s<", tc.name)

		t.Run(tc.name, func(t *testing.T) {

			// Test harness
			_, err = th.Setup()
			require.NoError(t, err, "Setup returns without error")
			defer func() {
				err = th.RollbackTx()
				require.NoError(t, err, "RollbackTx returns without error")
				err = th.Teardown()
				require.NoError(t, err, "Teardown returns without error")
			}()

			// init tx
			_, err = th.InitTx()
			require.NoError(t, err, "InitTx returns without error")

			dungeonInstanceID := tc.dungeonInstanceID(th.Data)
			characterInstanceID := tc.characterInstanceID(th.Data)

			sentence := tc.sentence(th.Data)
			t.Logf("Sentence >%s<", sentence)

			rslt, err := th.Model.(*model.Model).ProcessCharacterAction(dungeonInstanceID, characterInstanceID, sentence)
			require.NoError(t, err, "Initial ProcessCharacterAction returns without error")
			require.NotNil(t, rslt.ActionRec, "Initial ProcessCharacterAction returns ActionRecordSet with ActionRec")

			// Increment turn
			for inc := 0; inc < tc.incrementTurns; inc++ {
				t.Logf("Incrementing turn >%d<", inc+1)

				turnDuration := time.Duration(0) * time.Millisecond
				incrslt, err := th.Model.(*model.Model).IncrementDungeonInstanceTurn(&model.IncrementDungeonInstanceTurnArgs{
					DungeonInstanceID: dungeonInstanceID,
					TurnDuration:      &turnDuration,
				})
				require.NoError(t, err, "IncrementDungeonInstanceTurn returns without error")
				require.True(t, incrslt.Incremented, "IncrementDungeonInstanceTurn increments dungeon instance turn")
			}

			rslt, err = th.Model.(*model.Model).ProcessCharacterAction(dungeonInstanceID, characterInstanceID, sentence)
			require.Equal(t, tc.expectError, err, "Second ProcessCharacterAction error response equals expected")
			if err != nil {
				return
			}

			require.NotNil(t, rslt.ActionRec, "Second ProcessCharacterAction returns ActionRecordSet with ActionRec")
		})
	}
}

func TestConsecutiveProcessMonsterActions(t *testing.T) {

	c, l, s, err := dependencies.Default()
	require.NoError(t, err, "NewTesting returns without error")

	config := harness.DefaultDataConfig
	th, err := harness.NewTesting(c, l, s, config)
	require.NoError(t, err, "NewTesting returns without error")

	// harness commit data
	th.ShouldCommitData = true

	tests := []struct {
		name              string
		dungeonInstanceID func(data harness.Data) string
		monsterInstanceID func(data harness.Data) string
		sentence          func(data harness.Data) string
		incrementTurns    int
		expectError       error
	}{
		{
			name: "Too fast",
			dungeonInstanceID: func(data harness.Data) string {
				diRec, _ := data.GetDungeonInstanceRecByName(harness.DungeonNameCave)
				return diRec.ID
			},
			monsterInstanceID: func(data harness.Data) string {
				ciRec, _ := data.GetMonsterInstanceRecByName(harness.MonsterNameGrumpyDwarf)
				return ciRec.ID
			},
			sentence: func(data harness.Data) string {
				return "look"
			},
			incrementTurns: 0,
			expectError:    model.NewActionTooEarlyError(2, 2),
		},
		{
			name: "Too slow",
			dungeonInstanceID: func(data harness.Data) string {
				diRec, _ := data.GetDungeonInstanceRecByName(harness.DungeonNameCave)
				return diRec.ID
			},
			monsterInstanceID: func(data harness.Data) string {
				ciRec, _ := data.GetMonsterInstanceRecByName(harness.MonsterNameGrumpyDwarf)
				return ciRec.ID
			},
			sentence: func(data harness.Data) string {
				return "look"
			},
			incrementTurns: 2,
			expectError:    nil,
		},
		{
			name: "Just right",
			dungeonInstanceID: func(data harness.Data) string {
				diRec, _ := data.GetDungeonInstanceRecByName(harness.DungeonNameCave)
				return diRec.ID
			},
			monsterInstanceID: func(data harness.Data) string {
				ciRec, _ := data.GetMonsterInstanceRecByName(harness.MonsterNameGrumpyDwarf)
				return ciRec.ID
			},
			sentence: func(data harness.Data) string {
				return "look"
			},
			incrementTurns: 1,
			expectError:    nil,
		},
	}

	for _, tc := range tests {

		t.Logf("Run test >%s<", tc.name)

		t.Run(tc.name, func(t *testing.T) {

			// Test harness
			_, err = th.Setup()
			require.NoError(t, err, "Setup returns without error")
			defer func() {
				err = th.RollbackTx()
				require.NoError(t, err, "RollbackTx returns without error")
				err = th.Teardown()
				require.NoError(t, err, "Teardown returns without error")
			}()

			// init tx
			_, err = th.InitTx()
			require.NoError(t, err, "InitTx returns without error")

			dungeonInstanceID := tc.dungeonInstanceID(th.Data)
			monsterInstanceID := tc.monsterInstanceID(th.Data)

			sentence := tc.sentence(th.Data)
			t.Logf("Sentence >%s<", sentence)

			rslt, err := th.Model.(*model.Model).ProcessMonsterAction(dungeonInstanceID, monsterInstanceID, sentence)
			require.NoError(t, err, "Initial ProcessMonsterAction returns without error")
			require.NotNil(t, rslt.ActionRec, "Initial ProcessMonsterAction returns ActionRecordSet with ActionRec")

			// Increment turn
			for inc := 0; inc < tc.incrementTurns; inc++ {
				t.Logf("Incrementing turn >%d<", inc+1)

				turnDuration := time.Duration(0) * time.Millisecond
				incrslt, err := th.Model.(*model.Model).IncrementDungeonInstanceTurn(&model.IncrementDungeonInstanceTurnArgs{
					DungeonInstanceID: dungeonInstanceID,
					TurnDuration:      &turnDuration,
				})
				require.NoError(t, err, "IncrementDungeonInstanceTurn returns without error")
				require.True(t, incrslt.Incremented, "IncrementDungeonInstanceTurn increments dungeon instance turn")
			}

			rslt, err = th.Model.(*model.Model).ProcessMonsterAction(dungeonInstanceID, monsterInstanceID, sentence)
			require.Equal(t, tc.expectError, err, "Second ProcessMonsterAction error response equals expected")
			if err != nil {
				return
			}

			require.NotNil(t, rslt.ActionRec, "Second ProcessMonsterAction returns ActionRecordSet with ActionRec")
		})
	}
}

func TestProcessCharacterAction(t *testing.T) {

	// harness
	config := harness.DefaultDataConfig

	c, l, s, err := dependencies.Default()
	require.NoError(t, err, "NewTesting returns without error")

	th, err := harness.NewTesting(c, l, s, config)
	require.NoError(t, err, "NewTesting returns without error")

	// harness commit data
	th.ShouldCommitData = true

	tests := []struct {
		name                  string
		dungeonInstanceID     func(data harness.Data) string
		characterInstanceID   func(data harness.Data) string
		sentence              func(data harness.Data) string
		expectActionRecordSet func(data harness.Data) *record.ActionRecordSet
		expectError           bool
	}{
		{
			name: "look current location",
			dungeonInstanceID: func(data harness.Data) string {
				diRec, _ := data.GetDungeonInstanceRecByName("cave")
				return diRec.ID
			},
			characterInstanceID: func(data harness.Data) string {
				ciRec, _ := data.GetCharacterInstanceRecByName(harness.CharacterNameBarricade)
				return ciRec.ID
			},
			sentence: func(data harness.Data) string {
				return "look"
			},
			expectActionRecordSet: func(data harness.Data) *record.ActionRecordSet {

				cRec, _ := data.GetCharacterRecByName(harness.CharacterNameBarricade)
				ciRec, _ := data.GetCharacterInstanceRecByName(harness.CharacterNameBarricade)

				oiRecs := data.GetObjectInstanceRecsByCharacterInstanceID(ciRec.ID)
				acoRecs := []*record.ActionCharacterObject{}
				for _, oiRec := range oiRecs {
					oRec, _ := data.GetObjectRecByID(oiRec.ObjectID)
					acoRecs = append(acoRecs, &record.ActionCharacterObject{
						ObjectInstanceID: oiRec.ID,
						Name:             oRec.Name,
						IsStashed:        oiRec.IsStashed,
						IsEquipped:       oiRec.IsEquipped,
					})
				}

				diRec, _ := data.GetDungeonInstanceRecByName("cave")
				lRec, _ := data.GetLocationRecByName("cave entrance")
				liRec, _ := data.GetLocationInstanceRecByName("cave entrance")

				return &record.ActionRecordSet{
					ActionRec: &record.Action{
						DungeonInstanceID:   diRec.ID,
						LocationInstanceID:  liRec.ID,
						CharacterInstanceID: null.NullStringFromString(ciRec.ID),
					},
					ActionCharacterRec: &record.ActionCharacter{
						Name:                cRec.Name,
						Strength:            cRec.Strength,
						Dexterity:           cRec.Dexterity,
						Intelligence:        cRec.Intelligence,
						CurrentStrength:     ciRec.Strength,
						CurrentDexterity:    ciRec.Dexterity,
						CurrentIntelligence: ciRec.Intelligence,
						Health:              cRec.Health,
						Fatigue:             cRec.Fatigue,
						CurrentHealth:       ciRec.Health,
						CurrentFatigue:      ciRec.Fatigue,
					},
					ActionCharacterObjectRecs: acoRecs,
					TargetLocation: &record.ActionLocationRecordSet{
						LocationInstanceViewRec: &record.LocationInstanceView{
							Name: lRec.Name,
						},
					},
				}
			},
			expectError: false,
		},
		{
			name: "look valid direction",
			dungeonInstanceID: func(data harness.Data) string {
				diRec, _ := data.GetDungeonInstanceRecByName("cave")
				return diRec.ID
			},
			characterInstanceID: func(data harness.Data) string {
				ciRec, _ := data.GetCharacterInstanceRecByName(harness.CharacterNameBarricade)
				return ciRec.ID
			},
			sentence: func(data harness.Data) string {
				return "look north"
			},
			expectActionRecordSet: func(data harness.Data) *record.ActionRecordSet {

				cRec, _ := data.GetCharacterRecByName(harness.CharacterNameBarricade)
				ciRec, _ := data.GetCharacterInstanceRecByName(harness.CharacterNameBarricade)

				oiRecs := data.GetObjectInstanceRecsByCharacterInstanceID(ciRec.ID)
				acoRecs := []*record.ActionCharacterObject{}
				for _, oiRec := range oiRecs {
					oRec, _ := data.GetObjectRecByID(oiRec.ObjectID)
					acoRecs = append(acoRecs, &record.ActionCharacterObject{
						ObjectInstanceID: oiRec.ID,
						Name:             oRec.Name,
						IsStashed:        oiRec.IsStashed,
						IsEquipped:       oiRec.IsEquipped,
					})
				}

				diRec, _ := data.GetDungeonInstanceRecByName("cave")
				liRec, _ := data.GetLocationInstanceRecByName("cave entrance")
				tlRec, _ := data.GetLocationRecByName("cave tunnel")

				return &record.ActionRecordSet{
					ActionRec: &record.Action{
						DungeonInstanceID:   diRec.ID,
						LocationInstanceID:  liRec.ID,
						CharacterInstanceID: null.NullStringFromString(ciRec.ID),
					},
					ActionCharacterRec: &record.ActionCharacter{
						Name:                cRec.Name,
						Strength:            cRec.Strength,
						Dexterity:           cRec.Dexterity,
						Intelligence:        cRec.Intelligence,
						CurrentStrength:     ciRec.Strength,
						CurrentDexterity:    ciRec.Dexterity,
						CurrentIntelligence: ciRec.Intelligence,
						Health:              cRec.Health,
						Fatigue:             cRec.Fatigue,
						CurrentHealth:       ciRec.Health,
						CurrentFatigue:      ciRec.Fatigue,
					},
					ActionCharacterObjectRecs: acoRecs,
					TargetLocation: &record.ActionLocationRecordSet{
						LocationInstanceViewRec: &record.LocationInstanceView{
							Name: tlRec.Name,
						},
					},
				}
			},
			expectError: false,
		},
		{
			name: "look valid item",
			dungeonInstanceID: func(data harness.Data) string {
				diRec, _ := data.GetDungeonInstanceRecByName("cave")
				return diRec.ID
			},
			characterInstanceID: func(data harness.Data) string {
				ciRec, _ := data.GetCharacterInstanceRecByName(harness.CharacterNameBarricade)
				return ciRec.ID
			},
			sentence: func(data harness.Data) string {
				toRec, _ := data.GetObjectRecByName("rusted sword")
				return fmt.Sprintf("look %s", toRec.Name)
			},
			expectActionRecordSet: func(data harness.Data) *record.ActionRecordSet {

				cRec, _ := data.GetCharacterRecByName(harness.CharacterNameBarricade)
				ciRec, _ := data.GetCharacterInstanceRecByName(harness.CharacterNameBarricade)

				objectInstanceRecs := data.GetObjectInstanceRecsByCharacterInstanceID(ciRec.ID)
				acoRecs := []*record.ActionCharacterObject{}
				for _, objectInstanceRec := range objectInstanceRecs {
					oRec, _ := data.GetObjectRecByID(objectInstanceRec.ObjectID)
					acoRecs = append(acoRecs, &record.ActionCharacterObject{
						ObjectInstanceID: objectInstanceRec.ID,
						Name:             oRec.Name,
						IsStashed:        objectInstanceRec.IsStashed,
						IsEquipped:       objectInstanceRec.IsEquipped,
					})
				}

				diRec, _ := data.GetDungeonInstanceRecByName("cave")
				liRec, _ := data.GetLocationInstanceRecByName("cave entrance")
				toRec, _ := data.GetObjectRecByName("rusted sword")
				toiRec, _ := data.GetObjectInstanceRecByName("rusted sword")

				return &record.ActionRecordSet{
					ActionRec: &record.Action{
						DungeonInstanceID:   diRec.ID,
						LocationInstanceID:  liRec.ID,
						CharacterInstanceID: sql.NullString{String: ciRec.ID, Valid: true},
					},
					ActionCharacterRec: &record.ActionCharacter{
						Name:                cRec.Name,
						Strength:            cRec.Strength,
						Dexterity:           cRec.Dexterity,
						Intelligence:        cRec.Intelligence,
						CurrentStrength:     ciRec.Strength,
						CurrentDexterity:    ciRec.Dexterity,
						CurrentIntelligence: ciRec.Intelligence,
						Health:              cRec.Health,
						Fatigue:             cRec.Fatigue,
						CurrentHealth:       ciRec.Health,
						CurrentFatigue:      ciRec.Fatigue,
					},
					ActionCharacterObjectRecs: acoRecs,
					TargetActionObjectRec: &record.ActionObject{
						Name:        toRec.Name,
						Description: toRec.Description,
						IsStashed:   toiRec.IsStashed,
						IsEquipped:  toiRec.IsEquipped,
					},
				}
			},
			expectError: false,
		},
		{
			name: "look valid monster",
			dungeonInstanceID: func(data harness.Data) string {
				diRec, _ := data.GetDungeonInstanceRecByName("cave")
				return diRec.ID
			},
			characterInstanceID: func(data harness.Data) string {
				ciRec, _ := data.GetCharacterInstanceRecByName(harness.CharacterNameBarricade)
				return ciRec.ID
			},
			sentence: func(data harness.Data) string {
				tmRec, _ := data.GetMonsterRecByName("grumpy dwarf")
				return fmt.Sprintf("look %s", tmRec.Name)
			},
			expectActionRecordSet: func(data harness.Data) *record.ActionRecordSet {

				cRec, _ := data.GetCharacterRecByName(harness.CharacterNameBarricade)
				ciRec, _ := data.GetCharacterInstanceRecByName(harness.CharacterNameBarricade)

				oiRecs := data.GetObjectInstanceRecsByCharacterInstanceID(ciRec.ID)
				acoRecs := []*record.ActionCharacterObject{}
				for _, oiRec := range oiRecs {
					oRec, _ := data.GetObjectRecByID(oiRec.ObjectID)
					acoRecs = append(acoRecs, &record.ActionCharacterObject{
						ObjectInstanceID: oiRec.ID,
						Name:             oRec.Name,
						IsStashed:        oiRec.IsStashed,
						IsEquipped:       oiRec.IsEquipped,
					})
				}

				oiRecs = data.GetEquippedObjectInstanceRecsByMonsterInstanceID(data.MonsterInstanceRecs[0].ID)
				amoRecs := []*record.ActionMonsterObject{}
				for _, oiRec := range oiRecs {
					oRec, _ := data.GetObjectRecByID(oiRec.ObjectID)
					amoRecs = append(amoRecs, &record.ActionMonsterObject{
						ObjectInstanceID: oiRec.ID,
						Name:             oRec.Name,
						IsStashed:        oiRec.IsStashed,
						IsEquipped:       oiRec.IsEquipped,
					})
				}

				diRec, _ := data.GetDungeonInstanceRecByName("cave")
				liRec, _ := data.GetLocationInstanceRecByName("cave entrance")
				mRec, _ := data.GetMonsterRecByName("grumpy dwarf")
				miRec, _ := data.GetMonsterInstanceRecByName("grumpy dwarf")

				return &record.ActionRecordSet{
					ActionRec: &record.Action{
						DungeonInstanceID:   diRec.ID,
						LocationInstanceID:  liRec.ID,
						CharacterInstanceID: null.NullStringFromString(ciRec.ID),
					},
					ActionCharacterRec: &record.ActionCharacter{
						Name:                cRec.Name,
						Strength:            cRec.Strength,
						Dexterity:           cRec.Dexterity,
						Intelligence:        cRec.Intelligence,
						CurrentStrength:     ciRec.Strength,
						CurrentDexterity:    ciRec.Dexterity,
						CurrentIntelligence: ciRec.Intelligence,
						Health:              cRec.Health,
						Fatigue:             cRec.Fatigue,
						CurrentHealth:       ciRec.Health,
						CurrentFatigue:      ciRec.Fatigue,
					},
					ActionCharacterObjectRecs: acoRecs,
					TargetActionMonsterRec: &record.ActionMonster{
						Name:                mRec.Name,
						Strength:            mRec.Strength,
						Dexterity:           mRec.Dexterity,
						Intelligence:        mRec.Intelligence,
						CurrentStrength:     miRec.Strength,
						CurrentDexterity:    miRec.Dexterity,
						CurrentIntelligence: miRec.Intelligence,
						Health:              mRec.Health,
						Fatigue:             mRec.Fatigue,
						CurrentHealth:       miRec.Health,
						CurrentFatigue:      miRec.Fatigue,
					},
					TargetActionMonsterObjectRecs: amoRecs,
				}
			},
			expectError: false,
		},
		{
			name: "look valid character",
			dungeonInstanceID: func(data harness.Data) string {
				diRec, _ := data.GetDungeonInstanceRecByName("cave")
				return diRec.ID
			},
			characterInstanceID: func(data harness.Data) string {
				ciRec, _ := data.GetCharacterInstanceRecByName(harness.CharacterNameBarricade)
				return ciRec.ID
			},
			sentence: func(data harness.Data) string {
				tcRec, _ := data.GetCharacterRecByName(harness.CharacterNameBarricade)
				return fmt.Sprintf("look %s", tcRec.Name)
			},
			expectActionRecordSet: func(data harness.Data) *record.ActionRecordSet {

				cRec, _ := data.GetCharacterRecByName(harness.CharacterNameBarricade)
				ciRec, _ := data.GetCharacterInstanceRecByName(harness.CharacterNameBarricade)

				oiRecs := data.GetObjectInstanceRecsByCharacterInstanceID(ciRec.ID)
				acoRecs := []*record.ActionCharacterObject{}
				for _, oiRec := range oiRecs {
					oRec, _ := data.GetObjectRecByID(oiRec.ObjectID)
					acoRecs = append(acoRecs, &record.ActionCharacterObject{
						ObjectInstanceID: oiRec.ID,
						Name:             oRec.Name,
						IsStashed:        oiRec.IsStashed,
						IsEquipped:       oiRec.IsEquipped,
					})
				}

				oiRecs = data.GetEquippedObjectInstanceRecsByCharacterInstanceID(ciRec.ID)
				tacoRecs := []*record.ActionCharacterObject{}
				for _, oiRec := range oiRecs {
					oRec, _ := data.GetObjectRecByID(oiRec.ObjectID)
					tacoRecs = append(tacoRecs, &record.ActionCharacterObject{
						ObjectInstanceID: oiRec.ID,
						Name:             oRec.Name,
						IsStashed:        oiRec.IsStashed,
						IsEquipped:       oiRec.IsEquipped,
					})
				}

				diRec, _ := data.GetDungeonInstanceRecByName("cave")
				liRec, _ := data.GetLocationInstanceRecByName("cave entrance")

				return &record.ActionRecordSet{
					ActionRec: &record.Action{
						DungeonInstanceID:   diRec.ID,
						LocationInstanceID:  liRec.ID,
						CharacterInstanceID: null.NullStringFromString(ciRec.ID),
					},
					ActionCharacterRec: &record.ActionCharacter{
						Name:                cRec.Name,
						Strength:            cRec.Strength,
						Dexterity:           cRec.Dexterity,
						Intelligence:        cRec.Intelligence,
						CurrentStrength:     ciRec.Strength,
						CurrentDexterity:    ciRec.Dexterity,
						CurrentIntelligence: ciRec.Intelligence,
						Health:              cRec.Health,
						Fatigue:             cRec.Fatigue,
						CurrentHealth:       ciRec.Health,
						CurrentFatigue:      ciRec.Fatigue,
					},
					ActionCharacterObjectRecs: acoRecs,
					TargetActionCharacterRec: &record.ActionCharacter{
						Name:                cRec.Name,
						Strength:            cRec.Strength,
						Dexterity:           cRec.Dexterity,
						Intelligence:        cRec.Intelligence,
						CurrentStrength:     ciRec.Strength,
						CurrentDexterity:    ciRec.Dexterity,
						CurrentIntelligence: ciRec.Intelligence,
						Health:              cRec.Health,
						Fatigue:             cRec.Fatigue,
						CurrentHealth:       ciRec.Health,
						CurrentFatigue:      ciRec.Fatigue,
					},
					TargetActionCharacterObjectRecs: tacoRecs,
				}
			},
			expectError: false,
		},
		{
			name: "stash valid location item",
			dungeonInstanceID: func(data harness.Data) string {
				diRec, _ := data.GetDungeonInstanceRecByName("cave")
				return diRec.ID
			},
			characterInstanceID: func(data harness.Data) string {
				ciRec, _ := data.GetCharacterInstanceRecByName(harness.CharacterNameBarricade)
				return ciRec.ID
			},
			sentence: func(data harness.Data) string {
				toRec, _ := data.GetObjectRecByName("rusted sword")
				return fmt.Sprintf("stash %s", toRec.Name)
			},
			expectActionRecordSet: func(data harness.Data) *record.ActionRecordSet {

				cRec, _ := data.GetCharacterRecByName(harness.CharacterNameBarricade)
				ciRec, _ := data.GetCharacterInstanceRecByName(harness.CharacterNameBarricade)

				oiRecs := data.GetObjectInstanceRecsByCharacterInstanceID(ciRec.ID)
				oiRec, _ := data.GetObjectInstanceRecByName("rusted sword")
				oiRec.IsStashed = true
				oiRec.IsEquipped = false
				oiRecs = append(oiRecs, oiRec)

				oRec, _ := data.GetObjectRecByID(oiRec.ObjectID)

				acoRecs := []*record.ActionCharacterObject{}
				for _, oiRec := range oiRecs {
					oRec, _ := data.GetObjectRecByID(oiRec.ObjectID)
					acoRecs = append(acoRecs, &record.ActionCharacterObject{
						ObjectInstanceID: oiRec.ID,
						Name:             oRec.Name,
						IsStashed:        oiRec.IsStashed,
						IsEquipped:       oiRec.IsEquipped,
					})
				}

				diRec, _ := data.GetDungeonInstanceRecByName("cave")
				liRec, _ := data.GetLocationInstanceRecByName("cave entrance")

				return &record.ActionRecordSet{
					ActionRec: &record.Action{
						DungeonInstanceID:   diRec.ID,
						LocationInstanceID:  liRec.ID,
						CharacterInstanceID: null.NullStringFromString(ciRec.ID),
					},
					ActionCharacterRec: &record.ActionCharacter{
						Name:                cRec.Name,
						Strength:            cRec.Strength,
						Dexterity:           cRec.Dexterity,
						Intelligence:        cRec.Intelligence,
						CurrentStrength:     ciRec.Strength,
						CurrentDexterity:    ciRec.Dexterity,
						CurrentIntelligence: ciRec.Intelligence,
						Health:              cRec.Health,
						Fatigue:             cRec.Fatigue,
						CurrentHealth:       ciRec.Health,
						CurrentFatigue:      ciRec.Fatigue,
					},
					ActionCharacterObjectRecs: acoRecs,
					TargetActionObjectRec: &record.ActionObject{
						Name:        oRec.Name,
						Description: oRec.Description,
						IsStashed:   true,
						IsEquipped:  false,
					},
					StashedActionObjectRec: &record.ActionObject{
						Name:        oRec.Name,
						Description: oRec.Description,
						IsStashed:   true,
						IsEquipped:  false,
					},
				}
			},
			expectError: false,
		},
		{
			name: "stash valid equipped item",
			dungeonInstanceID: func(data harness.Data) string {
				diRec, _ := data.GetDungeonInstanceRecByName("cave")
				return diRec.ID
			},
			characterInstanceID: func(data harness.Data) string {
				ciRec, _ := data.GetCharacterInstanceRecByName(harness.CharacterNameBarricade)
				return ciRec.ID
			},
			sentence: func(data harness.Data) string {
				toRec, _ := data.GetObjectRecByName("dull bronze ring")
				return fmt.Sprintf("stash %s", toRec.Name)
			},
			expectActionRecordSet: func(data harness.Data) *record.ActionRecordSet {

				cRec, _ := data.GetCharacterRecByName(harness.CharacterNameBarricade)
				ciRec, _ := data.GetCharacterInstanceRecByName(harness.CharacterNameBarricade)

				oiRecs := data.GetObjectInstanceRecsByCharacterInstanceID(ciRec.ID)
				toRec, _ := data.GetObjectRecByName("dull bronze ring")
				for i := range oiRecs {
					if oiRecs[i].ObjectID == toRec.ID {
						oiRecs[i].IsStashed = true
						oiRecs[i].IsEquipped = false
						break
					}
				}

				acoRecs := []*record.ActionCharacterObject{}
				for _, oiRec := range oiRecs {
					oRec, _ := data.GetObjectRecByID(oiRec.ObjectID)
					acoRecs = append(acoRecs, &record.ActionCharacterObject{
						ObjectInstanceID: oiRec.ID,
						Name:             oRec.Name,
						IsStashed:        oiRec.IsStashed,
						IsEquipped:       oiRec.IsEquipped,
					})
				}

				diRec, _ := data.GetDungeonInstanceRecByName("cave")
				liRec, _ := data.GetLocationInstanceRecByName("cave entrance")

				return &record.ActionRecordSet{
					ActionRec: &record.Action{
						DungeonInstanceID:   diRec.ID,
						LocationInstanceID:  liRec.ID,
						CharacterInstanceID: null.NullStringFromString(ciRec.ID),
					},
					ActionCharacterRec: &record.ActionCharacter{
						Name:                cRec.Name,
						Strength:            cRec.Strength,
						Dexterity:           cRec.Dexterity,
						Intelligence:        cRec.Intelligence,
						CurrentStrength:     ciRec.Strength,
						CurrentDexterity:    ciRec.Dexterity,
						CurrentIntelligence: ciRec.Intelligence,
						Health:              cRec.Health,
						Fatigue:             cRec.Fatigue,
						CurrentHealth:       ciRec.Health,
						CurrentFatigue:      ciRec.Fatigue,
					},
					ActionCharacterObjectRecs: acoRecs,
					TargetActionObjectRec: &record.ActionObject{
						Name:        toRec.Name,
						Description: toRec.Description,
						IsStashed:   true,
						IsEquipped:  false,
					},
					StashedActionObjectRec: &record.ActionObject{
						Name:        toRec.Name,
						Description: toRec.Description,
						IsStashed:   true,
						IsEquipped:  false,
					},
				}
			},
			expectError: false,
		},
		{
			name: "equip valid location item",
			dungeonInstanceID: func(data harness.Data) string {
				diRec, _ := data.GetDungeonInstanceRecByName("cave")
				return diRec.ID
			},
			characterInstanceID: func(data harness.Data) string {
				ciRec, _ := data.GetCharacterInstanceRecByName(harness.CharacterNameBarricade)
				return ciRec.ID
			},
			sentence: func(data harness.Data) string {
				toRec, _ := data.GetObjectRecByName("rusted sword")
				return fmt.Sprintf("equip %s", toRec.Name)
			},
			expectActionRecordSet: func(data harness.Data) *record.ActionRecordSet {

				cRec, _ := data.GetCharacterRecByName(harness.CharacterNameBarricade)
				ciRec, _ := data.GetCharacterInstanceRecByName(harness.CharacterNameBarricade)

				oiRecs := data.GetObjectInstanceRecsByCharacterInstanceID(ciRec.ID)
				oiRec, _ := data.GetObjectInstanceRecByName("rusted sword")
				oiRec.IsStashed = false
				oiRec.IsEquipped = true
				oiRecs = append(oiRecs, oiRec)

				oRec, _ := data.GetObjectRecByID(oiRec.ObjectID)

				acoRecs := []*record.ActionCharacterObject{}
				for _, oiRec := range oiRecs {
					oRec, _ := data.GetObjectRecByID(oiRec.ObjectID)
					acoRecs = append(acoRecs, &record.ActionCharacterObject{
						ObjectInstanceID: oiRec.ID,
						Name:             oRec.Name,
						IsStashed:        oiRec.IsStashed,
						IsEquipped:       oiRec.IsEquipped,
					})
				}

				diRec, _ := data.GetDungeonInstanceRecByName("cave")
				liRec, _ := data.GetLocationInstanceRecByName("cave entrance")

				return &record.ActionRecordSet{
					ActionRec: &record.Action{
						DungeonInstanceID:   diRec.ID,
						LocationInstanceID:  liRec.ID,
						CharacterInstanceID: null.NullStringFromString(ciRec.ID),
					},
					ActionCharacterRec: &record.ActionCharacter{
						Name:                cRec.Name,
						Strength:            cRec.Strength,
						Dexterity:           cRec.Dexterity,
						Intelligence:        cRec.Intelligence,
						CurrentStrength:     ciRec.Strength,
						CurrentDexterity:    ciRec.Dexterity,
						CurrentIntelligence: ciRec.Intelligence,
						Health:              cRec.Health,
						Fatigue:             cRec.Fatigue,
						CurrentHealth:       ciRec.Health,
						CurrentFatigue:      ciRec.Fatigue,
					},
					ActionCharacterObjectRecs: acoRecs,
					TargetActionObjectRec: &record.ActionObject{
						Name:        oRec.Name,
						Description: oRec.Description,
						IsStashed:   false,
						IsEquipped:  true,
					},
					EquippedActionObjectRec: &record.ActionObject{
						Name:        oRec.Name,
						Description: oRec.Description,
						IsStashed:   false,
						IsEquipped:  true,
					},
				}
			},
			expectError: false,
		},
		{
			name: "equip valid stashed item",
			dungeonInstanceID: func(data harness.Data) string {
				diRec, _ := data.GetDungeonInstanceRecByName("cave")
				return diRec.ID
			},
			characterInstanceID: func(data harness.Data) string {
				ciRec, _ := data.GetCharacterInstanceRecByName(harness.CharacterNameBarricade)
				return ciRec.ID
			},
			sentence: func(data harness.Data) string {
				toRec, _ := data.GetObjectRecByName("blood stained pouch")
				return fmt.Sprintf("equip %s", toRec.Name)
			},
			expectActionRecordSet: func(data harness.Data) *record.ActionRecordSet {

				cRec, _ := data.GetCharacterRecByName(harness.CharacterNameBarricade)
				ciRec, _ := data.GetCharacterInstanceRecByName(harness.CharacterNameBarricade)

				oiRecs := data.GetObjectInstanceRecsByCharacterInstanceID(ciRec.ID)
				toRec, _ := data.GetObjectRecByName("blood stained pouch")
				for i := range oiRecs {
					if oiRecs[i].ObjectID == toRec.ID {
						oiRecs[i].IsStashed = false
						oiRecs[i].IsEquipped = true
						break
					}
				}

				acoRecs := []*record.ActionCharacterObject{}
				for _, oiRec := range oiRecs {
					oRec, _ := data.GetObjectRecByID(oiRec.ObjectID)
					acoRecs = append(acoRecs, &record.ActionCharacterObject{
						ObjectInstanceID: oiRec.ID,
						Name:             oRec.Name,
						IsStashed:        oiRec.IsStashed,
						IsEquipped:       oiRec.IsEquipped,
					})
				}

				diRec, _ := data.GetDungeonInstanceRecByName("cave")
				liRec, _ := data.GetLocationInstanceRecByName("cave entrance")

				return &record.ActionRecordSet{
					ActionRec: &record.Action{
						DungeonInstanceID:   diRec.ID,
						LocationInstanceID:  liRec.ID,
						CharacterInstanceID: null.NullStringFromString(ciRec.ID),
					},
					ActionCharacterRec: &record.ActionCharacter{
						Name:                cRec.Name,
						Strength:            cRec.Strength,
						Dexterity:           cRec.Dexterity,
						Intelligence:        cRec.Intelligence,
						CurrentStrength:     ciRec.Strength,
						CurrentDexterity:    ciRec.Dexterity,
						CurrentIntelligence: ciRec.Intelligence,
						Health:              cRec.Health,
						Fatigue:             cRec.Fatigue,
						CurrentHealth:       ciRec.Health,
						CurrentFatigue:      ciRec.Fatigue,
					},
					ActionCharacterObjectRecs: acoRecs,
					TargetActionObjectRec: &record.ActionObject{
						Name:        toRec.Name,
						Description: toRec.Description,
						IsStashed:   false,
						IsEquipped:  true,
					},
					EquippedActionObjectRec: &record.ActionObject{
						Name:        toRec.Name,
						Description: toRec.Description,
						IsStashed:   false,
						IsEquipped:  true,
					},
				}
			},
			expectError: false,
		},
		{
			name: "drop valid equipped item",
			dungeonInstanceID: func(data harness.Data) string {
				diRec, _ := data.GetDungeonInstanceRecByName("cave")
				return diRec.ID
			},
			characterInstanceID: func(data harness.Data) string {
				ciRec, _ := data.GetCharacterInstanceRecByName(harness.CharacterNameBarricade)
				return ciRec.ID
			},
			sentence: func(data harness.Data) string {
				toRec, _ := data.GetObjectRecByName("dull bronze ring")
				return fmt.Sprintf("drop %s", toRec.Name)
			},
			expectActionRecordSet: func(data harness.Data) *record.ActionRecordSet {

				cRec, _ := data.GetCharacterRecByName(harness.CharacterNameBarricade)
				ciRec, _ := data.GetCharacterInstanceRecByName(harness.CharacterNameBarricade)

				toRec, _ := data.GetObjectRecByName("dull bronze ring")
				oiRecs := data.GetObjectInstanceRecsByCharacterInstanceID(ciRec.ID)

				acoRecs := []*record.ActionCharacterObject{}
				for i := range oiRecs {
					if oiRecs[i].ObjectID == toRec.ID {
						continue
					}
					oRec, _ := data.GetObjectRecByID(oiRecs[i].ObjectID)
					acoRecs = append(acoRecs, &record.ActionCharacterObject{
						ObjectInstanceID: oiRecs[i].ID,
						Name:             oRec.Name,
						IsStashed:        oiRecs[i].IsStashed,
						IsEquipped:       oiRecs[i].IsEquipped,
					})
				}

				diRec, _ := data.GetDungeonInstanceRecByName("cave")
				liRec, _ := data.GetLocationInstanceRecByName("cave entrance")

				return &record.ActionRecordSet{
					ActionRec: &record.Action{
						DungeonInstanceID:   diRec.ID,
						LocationInstanceID:  liRec.ID,
						CharacterInstanceID: null.NullStringFromString(ciRec.ID),
					},
					ActionCharacterRec: &record.ActionCharacter{
						Name:                cRec.Name,
						Strength:            cRec.Strength,
						Dexterity:           cRec.Dexterity,
						Intelligence:        cRec.Intelligence,
						CurrentStrength:     ciRec.Strength,
						CurrentDexterity:    ciRec.Dexterity,
						CurrentIntelligence: ciRec.Intelligence,
						Health:              cRec.Health,
						Fatigue:             cRec.Fatigue,
						CurrentHealth:       ciRec.Health,
						CurrentFatigue:      ciRec.Fatigue,
					},
					ActionCharacterObjectRecs: acoRecs,
					TargetActionObjectRec: &record.ActionObject{
						Name:        toRec.Name,
						Description: toRec.Description,
						IsStashed:   false,
						IsEquipped:  false,
					},
					DroppedActionObjectRec: &record.ActionObject{
						Name:        toRec.Name,
						Description: toRec.Description,
						IsStashed:   false,
						IsEquipped:  false,
					},
				}
			},
			expectError: false,
		},
		{
			name: "drop valid stashed item",
			dungeonInstanceID: func(data harness.Data) string {
				diRec, _ := data.GetDungeonInstanceRecByName("cave")
				return diRec.ID
			},
			characterInstanceID: func(data harness.Data) string {
				ciRec, _ := data.GetCharacterInstanceRecByName(harness.CharacterNameBarricade)
				return ciRec.ID
			},
			sentence: func(data harness.Data) string {
				toRec, _ := data.GetObjectRecByName("blood stained pouch")
				return fmt.Sprintf("drop %s", toRec.Name)
			},
			expectActionRecordSet: func(data harness.Data) *record.ActionRecordSet {

				cRec, _ := data.GetCharacterRecByName(harness.CharacterNameBarricade)
				ciRec, _ := data.GetCharacterInstanceRecByName(harness.CharacterNameBarricade)

				toRec, _ := data.GetObjectRecByName("blood stained pouch")
				oiRecs := data.GetObjectInstanceRecsByCharacterInstanceID(ciRec.ID)

				acoRecs := []*record.ActionCharacterObject{}
				for i := range oiRecs {
					if oiRecs[i].ObjectID == toRec.ID {
						continue
					}
					oRec, _ := data.GetObjectRecByID(oiRecs[i].ObjectID)
					acoRecs = append(acoRecs, &record.ActionCharacterObject{
						ObjectInstanceID: oiRecs[i].ID,
						Name:             oRec.Name,
						IsStashed:        oiRecs[i].IsStashed,
						IsEquipped:       oiRecs[i].IsEquipped,
					})
				}

				diRec, _ := data.GetDungeonInstanceRecByName("cave")
				liRec, _ := data.GetLocationInstanceRecByName("cave entrance")

				return &record.ActionRecordSet{
					ActionRec: &record.Action{
						DungeonInstanceID:   diRec.ID,
						LocationInstanceID:  liRec.ID,
						CharacterInstanceID: null.NullStringFromString(ciRec.ID),
					},
					ActionCharacterRec: &record.ActionCharacter{
						Name:                cRec.Name,
						Strength:            cRec.Strength,
						Dexterity:           cRec.Dexterity,
						Intelligence:        cRec.Intelligence,
						CurrentStrength:     ciRec.Strength,
						CurrentDexterity:    ciRec.Dexterity,
						CurrentIntelligence: ciRec.Intelligence,
						Health:              cRec.Health,
						Fatigue:             cRec.Fatigue,
						CurrentHealth:       ciRec.Health,
						CurrentFatigue:      ciRec.Fatigue,
					},
					ActionCharacterObjectRecs: acoRecs,
					TargetActionObjectRec: &record.ActionObject{
						Name:        toRec.Name,
						Description: toRec.Description,
						IsStashed:   false,
						IsEquipped:  false,
					},
					DroppedActionObjectRec: &record.ActionObject{
						Name:        toRec.Name,
						Description: toRec.Description,
						IsStashed:   false,
						IsEquipped:  false,
					},
				}
			},
			expectError: false,
		},
		{
			name: "attack valid monster",
			dungeonInstanceID: func(data harness.Data) string {
				diRec, _ := data.GetDungeonInstanceRecByName("cave")
				return diRec.ID
			},
			characterInstanceID: func(data harness.Data) string {
				ciRec, _ := data.GetCharacterInstanceRecByName("legislate")
				return ciRec.ID
			},
			sentence: func(data harness.Data) string {
				tmRec, _ := data.GetMonsterRecByName("grumpy dwarf")
				return fmt.Sprintf("attack %s", tmRec.Name)
			},
			expectActionRecordSet: func(data harness.Data) *record.ActionRecordSet {

				cRec, _ := data.GetCharacterRecByName("legislate")
				ciRec, _ := data.GetCharacterInstanceRecByName("legislate")

				oiRecs := data.GetObjectInstanceRecsByCharacterInstanceID(ciRec.ID)
				acoRecs := []*record.ActionCharacterObject{}
				for _, oiRec := range oiRecs {
					oRec, _ := data.GetObjectRecByID(oiRec.ObjectID)
					acoRecs = append(acoRecs, &record.ActionCharacterObject{
						ObjectInstanceID: oiRec.ID,
						Name:             oRec.Name,
						IsStashed:        oiRec.IsStashed,
						IsEquipped:       oiRec.IsEquipped,
					})
				}

				oiRecs = data.GetEquippedObjectInstanceRecsByMonsterInstanceID(data.MonsterInstanceRecs[0].ID)
				amoRecs := []*record.ActionMonsterObject{}
				for _, oiRec := range oiRecs {
					oRec, _ := data.GetObjectRecByID(oiRec.ObjectID)
					amoRecs = append(amoRecs, &record.ActionMonsterObject{
						ObjectInstanceID: oiRec.ID,
						Name:             oRec.Name,
						IsStashed:        oiRec.IsStashed,
						IsEquipped:       oiRec.IsEquipped,
					})
				}

				diRec, _ := data.GetDungeonInstanceRecByName("cave")
				liRec, _ := data.GetLocationInstanceRecByName("cave entrance")
				mRec, _ := data.GetMonsterRecByName("grumpy dwarf")
				miRec, _ := data.GetMonsterInstanceRecByName("grumpy dwarf")

				return &record.ActionRecordSet{
					ActionRec: &record.Action{
						DungeonInstanceID:   diRec.ID,
						LocationInstanceID:  liRec.ID,
						CharacterInstanceID: null.NullStringFromString(ciRec.ID),
					},
					ActionCharacterRec: &record.ActionCharacter{
						Name:                cRec.Name,
						Strength:            cRec.Strength,
						Dexterity:           cRec.Dexterity,
						Intelligence:        cRec.Intelligence,
						CurrentStrength:     ciRec.Strength,
						CurrentDexterity:    ciRec.Dexterity,
						CurrentIntelligence: ciRec.Intelligence,
						Health:              cRec.Health,
						Fatigue:             cRec.Fatigue,
						CurrentHealth:       ciRec.Health,
						CurrentFatigue:      ciRec.Fatigue,
					},
					ActionCharacterObjectRecs: acoRecs,
					TargetActionMonsterRec: &record.ActionMonster{
						Name:                mRec.Name,
						Strength:            mRec.Strength,
						Dexterity:           mRec.Dexterity,
						Intelligence:        mRec.Intelligence,
						CurrentStrength:     miRec.Strength,
						CurrentDexterity:    miRec.Dexterity,
						CurrentIntelligence: miRec.Intelligence,
						Health:              mRec.Health,
						Fatigue:             mRec.Fatigue,
						CurrentHealth:       miRec.Health,
						CurrentFatigue:      miRec.Fatigue,
					},
					TargetActionMonsterObjectRecs: amoRecs,
				}
			},
			expectError: false,
		},
		{
			name: "attack valid monster with valid weapon",
			dungeonInstanceID: func(data harness.Data) string {
				diRec, _ := data.GetDungeonInstanceRecByName("cave")
				return diRec.ID
			},
			characterInstanceID: func(data harness.Data) string {
				ciRec, _ := data.GetCharacterInstanceRecByName("legislate")
				return ciRec.ID
			},
			sentence: func(data harness.Data) string {
				tmRec, _ := data.GetMonsterRecByName("grumpy dwarf")
				toRec, _ := data.GetObjectRecByName("stone mace")
				return fmt.Sprintf("attack %s with %s", tmRec.Name, toRec.Name)
			},
			expectActionRecordSet: func(data harness.Data) *record.ActionRecordSet {

				cRec, _ := data.GetCharacterRecByName("legislate")
				ciRec, _ := data.GetCharacterInstanceRecByName("legislate")

				eoRec, _ := data.GetObjectRecByName("stone mace")
				oiRecs := data.GetObjectInstanceRecsByCharacterInstanceID(ciRec.ID)

				acoRecs := []*record.ActionCharacterObject{}
				for idx := range oiRecs {
					oRec, _ := data.GetObjectRecByID(oiRecs[idx].ObjectID)
					acoRecs = append(acoRecs, &record.ActionCharacterObject{
						ObjectInstanceID: oiRecs[idx].ID,
						Name:             oRec.Name,
						IsStashed:        oiRecs[idx].IsStashed,
						IsEquipped:       oiRecs[idx].IsEquipped,
					})
				}

				oiRecs = data.GetEquippedObjectInstanceRecsByMonsterInstanceID(data.MonsterInstanceRecs[0].ID)
				amoRecs := []*record.ActionMonsterObject{}
				for _, oiRec := range oiRecs {
					oRec, _ := data.GetObjectRecByID(oiRec.ObjectID)
					amoRecs = append(amoRecs, &record.ActionMonsterObject{
						ObjectInstanceID: oiRec.ID,
						Name:             oRec.Name,
						IsStashed:        oiRec.IsStashed,
						IsEquipped:       oiRec.IsEquipped,
					})
				}

				diRec, _ := data.GetDungeonInstanceRecByName("cave")
				liRec, _ := data.GetLocationInstanceRecByName("cave entrance")
				mRec, _ := data.GetMonsterRecByName("grumpy dwarf")
				miRec, _ := data.GetMonsterInstanceRecByName("grumpy dwarf")

				return &record.ActionRecordSet{
					ActionRec: &record.Action{
						DungeonInstanceID:   diRec.ID,
						LocationInstanceID:  liRec.ID,
						CharacterInstanceID: null.NullStringFromString(ciRec.ID),
					},
					ActionCharacterRec: &record.ActionCharacter{
						Name:                cRec.Name,
						Strength:            cRec.Strength,
						Dexterity:           cRec.Dexterity,
						Intelligence:        cRec.Intelligence,
						CurrentStrength:     ciRec.Strength,
						CurrentDexterity:    ciRec.Dexterity,
						CurrentIntelligence: ciRec.Intelligence,
						Health:              cRec.Health,
						Fatigue:             cRec.Fatigue,
						CurrentHealth:       ciRec.Health,
						CurrentFatigue:      ciRec.Fatigue,
					},
					ActionCharacterObjectRecs: acoRecs,
					TargetActionMonsterRec: &record.ActionMonster{
						Name:                mRec.Name,
						Strength:            mRec.Strength,
						Dexterity:           mRec.Dexterity,
						Intelligence:        mRec.Intelligence,
						CurrentStrength:     miRec.Strength,
						CurrentDexterity:    miRec.Dexterity,
						CurrentIntelligence: miRec.Intelligence,
						Health:              mRec.Health,
						Fatigue:             mRec.Fatigue,
						CurrentHealth:       miRec.Health,
						CurrentFatigue:      miRec.Fatigue,
					},
					TargetActionMonsterObjectRecs: amoRecs,
					EquippedActionObjectRec: &record.ActionObject{
						Name:        eoRec.Name,
						Description: eoRec.Description,
						IsStashed:   false,
						IsEquipped:  true,
					},
				}
			},
			expectError: false,
		},
	}

	for _, tc := range tests {

		t.Logf("Run test >%s<", tc.name)

		t.Run(tc.name, func(t *testing.T) {

			// Test harness
			_, err = th.Setup()
			require.NoError(t, err, "Setup returns without error")
			defer func() {
				err = th.RollbackTx()
				require.NoError(t, err, "RollbackTx returns without error")
				err = th.Teardown()
				require.NoError(t, err, "Teardown returns without error")
			}()

			// init tx
			_, err = th.InitTx()
			require.NoError(t, err, "InitTx returns without error")

			dungeonInstanceID := tc.dungeonInstanceID(th.Data)
			characterInstanceID := tc.characterInstanceID(th.Data)

			sentence := tc.sentence(th.Data)
			t.Logf("Sentence >%s<", sentence)

			rslt, err := th.Model.(*model.Model).ProcessCharacterAction(dungeonInstanceID, characterInstanceID, sentence)
			if tc.expectError == true {
				require.Error(t, err, "ProcessCharacterAction returns error")
				return
			}
			require.NoError(t, err, "ProcessCharacterAction returns without error")
			require.NotNil(t, rslt.ActionRec, "ProcessCharacterAction returns ActionRecordSet with ActionRec")

			xrslt := tc.expectActionRecordSet(th.Data)
			if xrslt == nil {
				return
			}

			if xrslt.ActionRec != nil {
				t.Logf("Checking ActionRec >%s<", xrslt.ActionRec.DungeonInstanceID)
				require.NotNil(t, rslt.ActionRec, "ActionRec is not nil")
				require.Equal(t, xrslt.ActionRec.DungeonInstanceID, rslt.ActionRec.DungeonInstanceID, "ActionRec.DungeonInstanceID equals expected")
				require.Equal(t, xrslt.ActionRec.CharacterInstanceID, rslt.ActionRec.CharacterInstanceID, "ActionRec.CharacterInstanceID equals expected")
				require.Equal(t, xrslt.ActionRec.MonsterInstanceID, rslt.ActionRec.MonsterInstanceID, "ActionRec.MonsterInstanceID equals expected")
				require.Equal(t, xrslt.ActionRec.LocationInstanceID, rslt.ActionRec.LocationInstanceID, "ActionRec.LocationInstanceID equals expected")
			} else {
				require.Nil(t, rslt.ActionRec, "ActionRec is nil")
			}

			if xrslt.ActionCharacterRec != nil {
				t.Logf("Checking ActionCharacterRec >%s<", xrslt.ActionCharacterRec.Name)
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
					t.Logf("Checking ActionCharacterObject >%s<", xrObjectRec.Name)
					t.Logf("Checking ActionCharacterObject IsEquipped >%t<", xrObjectRec.IsEquipped)
					t.Logf("Checking ActionCharacterObject IsStashed >%t<", xrObjectRec.IsStashed)
					found := false
					for _, rObjectRec := range rslt.ActionCharacterObjectRecs {
						if xrObjectRec.Name == rObjectRec.Name {
							t.Logf("Have ActionCharacterObject Name >%s<", rObjectRec.Name)
							t.Logf("Have ActionCharacterObject IsEquipped >%t<", rObjectRec.IsEquipped)
							t.Logf("Have ActionCharacterObject IsStashed >%t<", rObjectRec.IsStashed)
							found = true
							require.Equal(t, xrObjectRec.IsEquipped, rObjectRec.IsEquipped, "ActionCharacterObjectRec IsEquipped equals expected")
							require.Equal(t, xrObjectRec.IsStashed, rObjectRec.IsStashed, "ActionCharacterObjectRec IsStashed equals expected")
							require.NotEmpty(t, rObjectRec.ActionCharacterID, "ActionCharacterObjectRec ActionCharacterID is not empty")
							require.NotEmpty(t, rObjectRec.ObjectInstanceID, "ActionCharacterObjectRec ObjectInstanceID is not empty")
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
							require.NotEmpty(t, rObjectRec.ActionMonsterID, "ActionMonsterObjectRec ActionMonsterID is not empty")
							require.NotEmpty(t, rObjectRec.ObjectInstanceID, "ActionMonsterObjectRec ObjectInstanceID is not empty")
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
							require.NotEmpty(t, rObjectRec.ActionCharacterID, "TargetActionCharacterObjectRec ActionCharacterID is not empty")
							require.NotEmpty(t, rObjectRec.ObjectInstanceID, "TargetActionCharacterObjectRec ObjectInstanceID is not empty")
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
							require.NotEmpty(t, rObjectRec.ActionMonsterID, "TargetActionMonsterObjectRec ActionMonsterID is not empty")
							require.NotEmpty(t, rObjectRec.ObjectInstanceID, "TargetActionMonsterObjectRec ObjectInstanceID is not empty")
						}
					}
					require.True(t, found, fmt.Sprintf("TargetActionMonsterObjectRec >%s< found", xrObjectRec.Name))
				}
			}

			if xrslt.TargetLocation != nil {
				require.NotNil(t, rslt.TargetLocation, "TargetLocation is not nil")
				if xrslt.TargetLocation.LocationInstanceViewRec != nil {
					require.NotNil(t, rslt.TargetLocation.LocationInstanceViewRec, "TargetLocation.LocationInstanceViewRec is not nil")
					require.Equal(t, xrslt.TargetLocation.LocationInstanceViewRec.Name, rslt.TargetLocation.LocationInstanceViewRec.Name, "TargetLocation.LocationInstanceViewRec equals expected")
				}

				// TODO: (game) Compare items, monsters and characters

			} else {
				require.Nil(t, rslt.TargetLocation, "TargetLocation is nil")
			}
		})
	}
}
