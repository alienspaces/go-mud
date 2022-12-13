package runner

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"gitlab.com/alienspaces/go-mud/server/schema"

	"github.com/stretchr/testify/require"
	"gitlab.com/alienspaces/go-mud/server/core/auth"
	"gitlab.com/alienspaces/go-mud/server/core/server"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/harness"
)

func TestPostActionHandler(t *testing.T) {
	t.Parallel()

	// Test harness
	th, err := NewTestHarness()
	require.NoError(t, err, "New test data returns without error")

	type testCase struct {
		TestCase
		expectResponseBody func(data harness.Data) *schema.ActionResponse
	}

	// validAuthToken - Generate a valid authentication token for this handler
	validAuthToken := func() string {
		authen, _ := auth.NewAuth(th.Config, th.Log)
		token, _ := authen.EncodeJWT(&auth.Claims{
			Roles:    []string{},
			Identity: map[string]interface{}{},
		})
		return token
	}

	testCaseHandlerConfig := func(rnr *Runner) server.HandlerConfig {
		return rnr.HandlerConfig[postAction]
	}

	testCaseRequestHeaders := func(data harness.Data) map[string]string {
		headers := map[string]string{
			"Authorization": "Bearer " + validAuthToken(),
			"X-Tx-Rollback": "true",
		}
		return headers
	}

	// All actions are performed by "Barricade" in the "Cave"
	testCaseRequestPathParams := func(data harness.Data) map[string]string {
		dRec, _ := data.GetDungeonRecByName("cave")
		cRec, _ := data.GetCharacterRecByName("barricade")

		params := map[string]string{
			":dungeon_id":   dRec.ID,
			":character_id": cRec.ID,
		}
		return params
	}

	testCaseResponseBody := func(body io.Reader) (interface{}, error) {
		var responseBody *schema.ActionResponse
		err = json.NewDecoder(body).Decode(&responseBody)
		return responseBody, err
	}

	testCases := []testCase{
		{
			// Look at the current room
			TestCase: TestCase{
				Name:              "Look at the current room",
				HandlerConfig:     testCaseHandlerConfig,
				RequestHeaders:    testCaseRequestHeaders,
				RequestPathParams: testCaseRequestPathParams,
				RequestBody: func(data harness.Data) interface{} {
					res := schema.ActionRequest{
						Data: schema.ActionRequestData{
							Sentence: "look",
						},
					}
					return &res
				},
				ResponseBody: testCaseResponseBody,
				ResponseCode: http.StatusOK,
			},
			expectResponseBody: func(data harness.Data) *schema.ActionResponse {

				cRec, _ := data.GetCharacterRecByName("barricade")
				ceoRec, _ := data.GetObjectRecByName("dull bronze ring")
				csoRec, _ := data.GetObjectRecByName("blood stained pouch")
				ciRec, _ := data.GetCharacterInstanceRecByName("barricade")
				lRec, _ := data.GetLocationRecByName("cave entrance")
				lcRec, _ := data.GetCharacterRecByName("legislate")
				loRec, _ := data.GetObjectRecByName("rusted sword")
				lmRec, _ := data.GetMonsterRecByName("grumpy dwarf")

				res := schema.ActionResponse{
					Data: []schema.ActionResponseData{
						{
							ActionCommand:   "look",
							ActionNarrative: fmt.Sprintf("%s looks", cRec.Name),
							ActionLocation: schema.ActionLocation{
								LocationName:        lRec.Name,
								LocationDescription: lRec.Description,
								LocationDirections:  []string{"north"},
								LocationCharacters: []schema.ActionLocationCharacter{
									{
										CharacterName: cRec.Name,
									},
									{
										CharacterName: lcRec.Name,
									},
								},
								LocationMonsters: []schema.ActionLocationMonster{
									{
										MonsterName: lmRec.Name,
									},
								},
								LocationObjects: []schema.ActionLocationObject{
									{
										ObjectName: loRec.Name,
									},
								},
							},
							ActionCharacter: &schema.ActionCharacter{
								CharacterName:                cRec.Name,
								CharacterStrength:            cRec.Strength,
								CharacterDexterity:           cRec.Dexterity,
								CharacterIntelligence:        cRec.Intelligence,
								CharacterHealth:              cRec.Health,
								CharacterFatigue:             cRec.Fatigue,
								CharacterCurrentStrength:     ciRec.Strength,
								CharacterCurrentDexterity:    ciRec.Dexterity,
								CharacterCurrentIntelligence: ciRec.Intelligence,
								CharacterCurrentHealth:       ciRec.Health,
								CharacterCurrentFatigue:      ciRec.Fatigue,
								CharacterEquippedObjects: []schema.ActionObject{
									{
										ObjectName: ceoRec.Name,
									},
								},
								CharacterStashedObjects: []schema.ActionObject{
									{
										ObjectName: csoRec.Name,
									},
								},
							},
							ActionMonster:         nil,
							ActionEquippedObject:  nil,
							ActionStashedObject:   nil,
							ActionTargetObject:    nil,
							ActionTargetCharacter: nil,
							ActionTargetMonster:   nil,
							ActionTargetLocation: &schema.ActionLocation{
								LocationName:        lRec.Name,
								LocationDescription: lRec.Description,
								LocationDirections:  []string{"north"},
								LocationCharacters: []schema.ActionLocationCharacter{
									{
										CharacterName: cRec.Name,
									},
									{
										CharacterName: lcRec.Name,
									},
								},
								LocationMonsters: []schema.ActionLocationMonster{
									{
										MonsterName: lmRec.Name,
									},
								},
								LocationObjects: []schema.ActionLocationObject{
									{
										ObjectName: loRec.Name,
									},
								},
							},
						},
					},
				}
				return &res
			},
		},
		{
			// Move north from the current room
			TestCase: TestCase{
				Name:              "Move north from the current room",
				HandlerConfig:     testCaseHandlerConfig,
				RequestHeaders:    testCaseRequestHeaders,
				RequestPathParams: testCaseRequestPathParams,
				RequestBody: func(data harness.Data) interface{} {
					res := schema.ActionRequest{
						Data: schema.ActionRequestData{
							Sentence: "move north",
						},
					}
					return &res
				},
				ResponseBody: testCaseResponseBody,
				ResponseCode: http.StatusOK,
			},
			expectResponseBody: func(data harness.Data) *schema.ActionResponse {

				cRec, _ := data.GetCharacterRecByName("barricade")
				ceoRec, _ := data.GetObjectRecByName("dull bronze ring")
				csoRec, _ := data.GetObjectRecByName("blood stained pouch")
				ciRec, _ := data.GetCharacterInstanceRecByName("barricade")
				lRec, _ := data.GetLocationRecByName("cave tunnel")
				loRec, _ := data.GetObjectRecByName("rusted helmet")
				lmRec, _ := data.GetMonsterRecByName("angry goblin")

				res := schema.ActionResponse{
					Data: []schema.ActionResponseData{
						{
							ActionCommand:   "move",
							ActionNarrative: fmt.Sprintf("%s moves north", cRec.Name),
							ActionLocation: schema.ActionLocation{
								LocationName:        lRec.Name,
								LocationDescription: lRec.Description,
								LocationDirections:  []string{"north", "south", "northwest"},
								LocationCharacters: []schema.ActionLocationCharacter{
									{
										CharacterName: cRec.Name,
									},
								},
								LocationMonsters: []schema.ActionLocationMonster{
									{
										MonsterName: lmRec.Name,
									},
								},
								LocationObjects: []schema.ActionLocationObject{
									{
										ObjectName: loRec.Name,
									},
								},
							},
							ActionCharacter: &schema.ActionCharacter{
								CharacterName:                cRec.Name,
								CharacterStrength:            cRec.Strength,
								CharacterDexterity:           cRec.Dexterity,
								CharacterIntelligence:        cRec.Intelligence,
								CharacterHealth:              cRec.Health,
								CharacterFatigue:             cRec.Fatigue,
								CharacterCurrentStrength:     ciRec.Strength,
								CharacterCurrentDexterity:    ciRec.Dexterity,
								CharacterCurrentIntelligence: ciRec.Intelligence,
								CharacterCurrentHealth:       ciRec.Health,
								CharacterCurrentFatigue:      ciRec.Fatigue,
								CharacterEquippedObjects: []schema.ActionObject{
									{
										ObjectName: ceoRec.Name,
									},
								},
								CharacterStashedObjects: []schema.ActionObject{
									{
										ObjectName: csoRec.Name,
									},
								},
							},
							ActionMonster:         nil,
							ActionEquippedObject:  nil,
							ActionStashedObject:   nil,
							ActionTargetObject:    nil,
							ActionTargetCharacter: nil,
							ActionTargetMonster:   nil,
							ActionTargetLocation: &schema.ActionLocation{
								LocationName:        lRec.Name,
								LocationDescription: lRec.Description,
								LocationDirection:   "north",
								LocationDirections:  []string{"north", "south", "northwest"},
								LocationCharacters: []schema.ActionLocationCharacter{
									{
										CharacterName: cRec.Name,
									},
								},
								LocationMonsters: []schema.ActionLocationMonster{
									{
										MonsterName: lmRec.Name,
									},
								},
								LocationObjects: []schema.ActionLocationObject{
									{
										ObjectName: loRec.Name,
									},
								},
							},
						},
					},
				}
				return &res
			},
		},
		{
			// Look north from the current room
			TestCase: TestCase{
				Name:              "Look north from the current room",
				HandlerConfig:     testCaseHandlerConfig,
				RequestHeaders:    testCaseRequestHeaders,
				RequestPathParams: testCaseRequestPathParams,
				RequestBody: func(data harness.Data) interface{} {
					res := schema.ActionRequest{
						Data: schema.ActionRequestData{
							Sentence: "look north",
						},
					}
					return &res
				},
				ResponseBody: testCaseResponseBody,
				ResponseCode: http.StatusOK,
			},
			expectResponseBody: func(data harness.Data) *schema.ActionResponse {

				cRec, _ := data.GetCharacterRecByName("barricade")
				ceoRec, _ := data.GetObjectRecByName("dull bronze ring")
				csoRec, _ := data.GetObjectRecByName("blood stained pouch")
				ciRec, _ := data.GetCharacterInstanceRecByName("barricade")
				lRec, _ := data.GetLocationRecByName("cave entrance")
				lcRec, _ := data.GetCharacterRecByName("legislate")
				loRec, _ := data.GetObjectRecByName("rusted sword")
				lmRec, _ := data.GetMonsterRecByName("grumpy dwarf")
				tlRec, _ := data.GetLocationRecByName("cave tunnel")
				tloRec, _ := data.GetObjectRecByName("rusted helmet")
				tlmRec, _ := data.GetMonsterRecByName("angry goblin")

				res := schema.ActionResponse{
					Data: []schema.ActionResponseData{
						{
							ActionCommand:   "look",
							ActionNarrative: fmt.Sprintf("%s looks north", cRec.Name),
							ActionLocation: schema.ActionLocation{
								LocationName:        lRec.Name,
								LocationDescription: lRec.Description,
								LocationDirections:  []string{"north"},
								LocationCharacters: []schema.ActionLocationCharacter{
									{
										CharacterName: cRec.Name,
									},
									{
										CharacterName: lcRec.Name,
									},
								},
								LocationMonsters: []schema.ActionLocationMonster{
									{
										MonsterName: lmRec.Name,
									},
								},
								LocationObjects: []schema.ActionLocationObject{
									{
										ObjectName: loRec.Name,
									},
								},
							},
							ActionCharacter: &schema.ActionCharacter{
								CharacterName:                cRec.Name,
								CharacterStrength:            cRec.Strength,
								CharacterDexterity:           cRec.Dexterity,
								CharacterIntelligence:        cRec.Intelligence,
								CharacterHealth:              cRec.Health,
								CharacterFatigue:             cRec.Fatigue,
								CharacterCurrentStrength:     ciRec.Strength,
								CharacterCurrentDexterity:    ciRec.Dexterity,
								CharacterCurrentIntelligence: ciRec.Intelligence,
								CharacterCurrentHealth:       ciRec.Health,
								CharacterCurrentFatigue:      ciRec.Fatigue,
								CharacterEquippedObjects: []schema.ActionObject{
									{
										ObjectName: ceoRec.Name,
									},
								},
								CharacterStashedObjects: []schema.ActionObject{
									{
										ObjectName: csoRec.Name,
									},
								},
							},
							ActionMonster:         nil,
							ActionEquippedObject:  nil,
							ActionStashedObject:   nil,
							ActionTargetObject:    nil,
							ActionTargetCharacter: nil,
							ActionTargetMonster:   nil,
							ActionTargetLocation: &schema.ActionLocation{
								LocationName:        tlRec.Name,
								LocationDescription: tlRec.Description,
								LocationDirection:   "north",
								LocationDirections:  []string{"north", "south", "northwest"},
								LocationMonsters: []schema.ActionLocationMonster{
									{
										MonsterName: tlmRec.Name,
									},
								},
								LocationObjects: []schema.ActionLocationObject{
									{
										ObjectName: tloRec.Name,
									},
								},
							},
						},
					},
				}
				return &res
			},
		},
		{
			// Look at an item in the current room
			TestCase: TestCase{
				Name:              "Look at an item in the current room",
				HandlerConfig:     testCaseHandlerConfig,
				RequestHeaders:    testCaseRequestHeaders,
				RequestPathParams: testCaseRequestPathParams,
				RequestBody: func(data harness.Data) interface{} {
					toRec, _ := data.GetObjectRecByName("rusted sword")
					res := schema.ActionRequest{
						Data: schema.ActionRequestData{
							Sentence: fmt.Sprintf("look %s", toRec.Name),
						},
					}
					return &res
				},
				ResponseBody: testCaseResponseBody,
				ResponseCode: http.StatusOK,
			},
			expectResponseBody: func(data harness.Data) *schema.ActionResponse {

				cRec, _ := data.GetCharacterRecByName("barricade")
				ceoRec, _ := data.GetObjectRecByName("dull bronze ring")
				csoRec, _ := data.GetObjectRecByName("blood stained pouch")
				ciRec, _ := data.GetCharacterInstanceRecByName("barricade")
				lRec, _ := data.GetLocationRecByName("cave entrance")
				lcRec, _ := data.GetCharacterRecByName("legislate")
				loRec, _ := data.GetObjectRecByName("rusted sword")
				lmRec, _ := data.GetMonsterRecByName("grumpy dwarf")
				toRec, _ := data.GetObjectRecByName("rusted sword")

				res := schema.ActionResponse{
					Data: []schema.ActionResponseData{
						{
							ActionCommand:   "look",
							ActionNarrative: fmt.Sprintf("%s looks %s", cRec.Name, loRec.Name),
							ActionLocation: schema.ActionLocation{
								LocationName:        lRec.Name,
								LocationDescription: lRec.Description,
								LocationDirections:  []string{"north"},
								LocationCharacters: []schema.ActionLocationCharacter{
									{
										CharacterName: cRec.Name,
									},
									{
										CharacterName: lcRec.Name,
									},
								},
								LocationMonsters: []schema.ActionLocationMonster{
									{
										MonsterName: lmRec.Name,
									},
								},
								LocationObjects: []schema.ActionLocationObject{
									{
										ObjectName: loRec.Name,
									},
								},
							},
							ActionCharacter: &schema.ActionCharacter{
								CharacterName:                cRec.Name,
								CharacterStrength:            cRec.Strength,
								CharacterDexterity:           cRec.Dexterity,
								CharacterIntelligence:        cRec.Intelligence,
								CharacterHealth:              cRec.Health,
								CharacterFatigue:             cRec.Fatigue,
								CharacterCurrentStrength:     ciRec.Strength,
								CharacterCurrentDexterity:    ciRec.Dexterity,
								CharacterCurrentIntelligence: ciRec.Intelligence,
								CharacterCurrentHealth:       ciRec.Health,
								CharacterCurrentFatigue:      ciRec.Fatigue,
								CharacterEquippedObjects: []schema.ActionObject{
									{
										ObjectName: ceoRec.Name,
									},
								},
								CharacterStashedObjects: []schema.ActionObject{
									{
										ObjectName: csoRec.Name,
									},
								},
							},
							ActionMonster:        nil,
							ActionEquippedObject: nil,
							ActionStashedObject:  nil,
							ActionTargetObject: &schema.ActionObject{
								ObjectName:        toRec.Name,
								ObjectDescription: toRec.Description,
							},
							ActionTargetCharacter: nil,
							ActionTargetMonster:   nil,
							ActionTargetLocation:  nil,
						},
					},
				}
				return &res
			},
		},
		{
			// Look at a monster in the current room
			TestCase: TestCase{
				Name:              "Look at a monster in the current room",
				HandlerConfig:     testCaseHandlerConfig,
				RequestHeaders:    testCaseRequestHeaders,
				RequestPathParams: testCaseRequestPathParams,
				RequestBody: func(data harness.Data) interface{} {
					tmRec, _ := data.GetMonsterRecByName("grumpy dwarf")
					res := schema.ActionRequest{
						Data: schema.ActionRequestData{
							Sentence: fmt.Sprintf("look %s", tmRec.Name),
						},
					}
					return &res
				},
				ResponseBody: testCaseResponseBody,
				ResponseCode: http.StatusOK,
			},
			expectResponseBody: func(data harness.Data) *schema.ActionResponse {

				cRec, _ := data.GetCharacterRecByName("barricade")
				ceoRec, _ := data.GetObjectRecByName("dull bronze ring")
				csoRec, _ := data.GetObjectRecByName("blood stained pouch")
				ciRec, _ := data.GetCharacterInstanceRecByName("barricade")
				lRec, _ := data.GetLocationRecByName("cave entrance")
				lcRec, _ := data.GetCharacterRecByName("legislate")
				loRec, _ := data.GetObjectRecByName("rusted sword")
				lmRec, _ := data.GetMonsterRecByName("grumpy dwarf")
				tmRec, _ := data.GetMonsterRecByName("grumpy dwarf")
				tmiRec, _ := data.GetMonsterInstanceRecByName("grumpy dwarf")
				tmeoRec, _ := data.GetObjectRecByName("bone dagger")

				res := schema.ActionResponse{
					Data: []schema.ActionResponseData{
						{
							ActionCommand:   "look",
							ActionNarrative: fmt.Sprintf("%s looks %s", cRec.Name, tmRec.Name),
							ActionLocation: schema.ActionLocation{
								LocationName:        lRec.Name,
								LocationDescription: lRec.Description,
								LocationDirections:  []string{"north"},
								LocationCharacters: []schema.ActionLocationCharacter{
									{
										CharacterName: cRec.Name,
									},
									{
										CharacterName: lcRec.Name,
									},
								},
								LocationMonsters: []schema.ActionLocationMonster{
									{
										MonsterName: lmRec.Name,
									},
								},
								LocationObjects: []schema.ActionLocationObject{
									{
										ObjectName: loRec.Name,
									},
								},
							},
							ActionCharacter: &schema.ActionCharacter{
								CharacterName:                cRec.Name,
								CharacterStrength:            cRec.Strength,
								CharacterDexterity:           cRec.Dexterity,
								CharacterIntelligence:        cRec.Intelligence,
								CharacterHealth:              cRec.Health,
								CharacterFatigue:             cRec.Fatigue,
								CharacterCurrentStrength:     ciRec.Strength,
								CharacterCurrentDexterity:    ciRec.Dexterity,
								CharacterCurrentIntelligence: ciRec.Intelligence,
								CharacterCurrentHealth:       ciRec.Health,
								CharacterCurrentFatigue:      ciRec.Fatigue,
								CharacterEquippedObjects: []schema.ActionObject{
									{
										ObjectName: ceoRec.Name,
									},
								},
								CharacterStashedObjects: []schema.ActionObject{
									{
										ObjectName: csoRec.Name,
									},
								},
							},
							ActionMonster:         nil,
							ActionEquippedObject:  nil,
							ActionStashedObject:   nil,
							ActionTargetObject:    nil,
							ActionTargetCharacter: nil,
							ActionTargetMonster: &schema.ActionMonster{
								MonsterName:                tmRec.Name,
								MonsterStrength:            tmRec.Strength,
								MonsterDexterity:           tmRec.Dexterity,
								MonsterIntelligence:        tmRec.Intelligence,
								MonsterHealth:              tmRec.Health,
								MonsterFatigue:             tmRec.Fatigue,
								MonsterCurrentStrength:     tmiRec.Strength,
								MonsterCurrentDexterity:    tmiRec.Dexterity,
								MonsterCurrentIntelligence: tmiRec.Intelligence,
								MonsterCurrentHealth:       tmiRec.Health,
								MonsterCurrentFatigue:      tmiRec.Fatigue,
								MonsterEquippedObjects: []schema.ActionObject{
									{
										ObjectName: tmeoRec.Name,
									},
								},
							},
							ActionTargetLocation: nil,
						},
					},
				}
				return &res
			},
		},
		{
			// Look at themselves in the current room
			TestCase: TestCase{
				Name:              "Look at themselves in the current room",
				HandlerConfig:     testCaseHandlerConfig,
				RequestHeaders:    testCaseRequestHeaders,
				RequestPathParams: testCaseRequestPathParams,
				RequestBody: func(data harness.Data) interface{} {
					tcRec, _ := data.GetCharacterRecByName("barricade")
					res := schema.ActionRequest{
						Data: schema.ActionRequestData{
							Sentence: fmt.Sprintf("look %s", tcRec.Name),
						},
					}
					return &res
				},
				ResponseBody: testCaseResponseBody,
				ResponseCode: http.StatusOK,
			},
			expectResponseBody: func(data harness.Data) *schema.ActionResponse {

				cRec, _ := data.GetCharacterRecByName("barricade")
				ceoRec, _ := data.GetObjectRecByName("dull bronze ring")
				csoRec, _ := data.GetObjectRecByName("blood stained pouch")
				ciRec, _ := data.GetCharacterInstanceRecByName("barricade")
				lRec, _ := data.GetLocationRecByName("cave entrance")
				lcRec, _ := data.GetCharacterRecByName("legislate")
				loRec, _ := data.GetObjectRecByName("rusted sword")
				lmRec, _ := data.GetMonsterRecByName("grumpy dwarf")
				tcRec, _ := data.GetCharacterRecByName("barricade")
				tciRec, _ := data.GetCharacterInstanceRecByName("barricade")
				tceoRec, _ := data.GetObjectRecByName("dull bronze ring")

				res := schema.ActionResponse{
					Data: []schema.ActionResponseData{
						{
							ActionCommand:   "look",
							ActionNarrative: fmt.Sprintf("%s looks %s", cRec.Name, tcRec.Name),
							ActionLocation: schema.ActionLocation{
								LocationName:        lRec.Name,
								LocationDescription: lRec.Description,
								LocationDirections:  []string{"north"},
								LocationCharacters: []schema.ActionLocationCharacter{
									{
										CharacterName: cRec.Name,
									},
									{
										CharacterName: lcRec.Name,
									},
								},
								LocationMonsters: []schema.ActionLocationMonster{
									{
										MonsterName: lmRec.Name,
									},
								},
								LocationObjects: []schema.ActionLocationObject{
									{
										ObjectName: loRec.Name,
									},
								},
							},
							ActionCharacter: &schema.ActionCharacter{
								CharacterName:                cRec.Name,
								CharacterStrength:            cRec.Strength,
								CharacterDexterity:           cRec.Dexterity,
								CharacterIntelligence:        cRec.Intelligence,
								CharacterHealth:              cRec.Health,
								CharacterFatigue:             cRec.Fatigue,
								CharacterCurrentStrength:     ciRec.Strength,
								CharacterCurrentDexterity:    ciRec.Dexterity,
								CharacterCurrentIntelligence: ciRec.Intelligence,
								CharacterCurrentHealth:       ciRec.Health,
								CharacterCurrentFatigue:      ciRec.Fatigue,
								CharacterEquippedObjects: []schema.ActionObject{
									{
										ObjectName: ceoRec.Name,
									},
								},
								CharacterStashedObjects: []schema.ActionObject{
									{
										ObjectName: csoRec.Name,
									},
								},
							},
							ActionMonster:        nil,
							ActionEquippedObject: nil,
							ActionStashedObject:  nil,
							ActionTargetObject:   nil,
							ActionTargetCharacter: &schema.ActionCharacter{
								CharacterName:                tcRec.Name,
								CharacterStrength:            tcRec.Strength,
								CharacterDexterity:           tcRec.Dexterity,
								CharacterIntelligence:        tcRec.Intelligence,
								CharacterHealth:              tcRec.Health,
								CharacterFatigue:             tcRec.Fatigue,
								CharacterCurrentStrength:     tciRec.Strength,
								CharacterCurrentDexterity:    tciRec.Dexterity,
								CharacterCurrentIntelligence: tciRec.Intelligence,
								CharacterCurrentHealth:       tciRec.Health,
								CharacterCurrentFatigue:      tciRec.Fatigue,
								CharacterEquippedObjects: []schema.ActionObject{
									{
										ObjectName: tceoRec.Name,
									},
								},
							},
							ActionTargetMonster:  nil,
							ActionTargetLocation: nil,
						},
					},
				}
				return &res
			},
		},
		{
			// Stash object that is in the current room
			TestCase: TestCase{
				Name:              "Stash object that is in the current room",
				HandlerConfig:     testCaseHandlerConfig,
				RequestHeaders:    testCaseRequestHeaders,
				RequestPathParams: testCaseRequestPathParams,
				RequestBody: func(data harness.Data) interface{} {
					toRec, _ := data.GetObjectRecByName("rusted sword")
					res := schema.ActionRequest{
						Data: schema.ActionRequestData{
							Sentence: fmt.Sprintf("stash %s", toRec.Name),
						},
					}
					return &res
				},
				ResponseBody: testCaseResponseBody,
				ResponseCode: http.StatusOK,
			},
			expectResponseBody: func(data harness.Data) *schema.ActionResponse {

				cRec, _ := data.GetCharacterRecByName("barricade")
				ceoRec, _ := data.GetObjectRecByName("dull bronze ring")
				csoRec, _ := data.GetObjectRecByName("blood stained pouch")
				ciRec, _ := data.GetCharacterInstanceRecByName("barricade")
				lRec, _ := data.GetLocationRecByName("cave entrance")
				lcRec, _ := data.GetCharacterRecByName("legislate")
				lmRec, _ := data.GetMonsterRecByName("grumpy dwarf")
				toRec, _ := data.GetObjectRecByName("rusted sword")

				res := schema.ActionResponse{
					Data: []schema.ActionResponseData{
						{
							ActionCommand:   "stash",
							ActionNarrative: fmt.Sprintf("%s stashes %s", cRec.Name, toRec.Name),
							ActionLocation: schema.ActionLocation{
								LocationName:        lRec.Name,
								LocationDescription: lRec.Description,
								LocationDirections:  []string{"north"},
								LocationCharacters: []schema.ActionLocationCharacter{
									{
										CharacterName: cRec.Name,
									},
									{
										CharacterName: lcRec.Name,
									},
								},
								LocationMonsters: []schema.ActionLocationMonster{
									{
										MonsterName: lmRec.Name,
									},
								},
								LocationObjects: []schema.ActionLocationObject{},
							},
							ActionCharacter: &schema.ActionCharacter{
								CharacterName:                cRec.Name,
								CharacterStrength:            cRec.Strength,
								CharacterDexterity:           cRec.Dexterity,
								CharacterIntelligence:        cRec.Intelligence,
								CharacterHealth:              cRec.Health,
								CharacterFatigue:             cRec.Fatigue,
								CharacterCurrentStrength:     ciRec.Strength,
								CharacterCurrentDexterity:    ciRec.Dexterity,
								CharacterCurrentIntelligence: ciRec.Intelligence,
								CharacterCurrentHealth:       ciRec.Health,
								CharacterCurrentFatigue:      ciRec.Fatigue,
								CharacterEquippedObjects: []schema.ActionObject{
									{
										ObjectName: ceoRec.Name,
									},
								},
								CharacterStashedObjects: []schema.ActionObject{
									{
										ObjectName: toRec.Name,
									},
									{
										ObjectName: csoRec.Name,
									},
								},
							},
							ActionMonster:        nil,
							ActionEquippedObject: nil,
							ActionStashedObject: &schema.ActionObject{
								ObjectName:        toRec.Name,
								ObjectDescription: toRec.Description,
								ObjectIsEquipped:  false,
								ObjectIsStashed:   true,
							},
							ActionDroppedObject: nil,
							ActionTargetObject: &schema.ActionObject{
								ObjectName:        toRec.Name,
								ObjectDescription: toRec.Description,
								ObjectIsEquipped:  false,
								ObjectIsStashed:   true,
							},
							ActionTargetCharacter: nil,
							ActionTargetMonster:   nil,
							ActionTargetLocation:  nil,
						},
					},
				}
				return &res
			},
		},
		{
			// Equip object that is in the current room
			TestCase: TestCase{
				Name:              "Equip object that is in the current room",
				HandlerConfig:     testCaseHandlerConfig,
				RequestHeaders:    testCaseRequestHeaders,
				RequestPathParams: testCaseRequestPathParams,
				RequestBody: func(data harness.Data) interface{} {
					toRec, _ := data.GetObjectRecByName("rusted sword")
					res := schema.ActionRequest{
						Data: schema.ActionRequestData{
							Sentence: fmt.Sprintf("equip %s", toRec.Name),
						},
					}
					return &res
				},
				ResponseBody: testCaseResponseBody,
				ResponseCode: http.StatusOK,
			},
			expectResponseBody: func(data harness.Data) *schema.ActionResponse {

				cRec, _ := data.GetCharacterRecByName("barricade")
				ceoRec, _ := data.GetObjectRecByName("dull bronze ring")
				csoRec, _ := data.GetObjectRecByName("blood stained pouch")
				ciRec, _ := data.GetCharacterInstanceRecByName("barricade")
				lRec, _ := data.GetLocationRecByName("cave entrance")
				lcRec, _ := data.GetCharacterRecByName("legislate")
				lmRec, _ := data.GetMonsterRecByName("grumpy dwarf")
				toRec, _ := data.GetObjectRecByName("rusted sword")

				res := schema.ActionResponse{
					Data: []schema.ActionResponseData{
						{
							ActionCommand:   "equip",
							ActionNarrative: fmt.Sprintf("%s equips %s", cRec.Name, toRec.Name),
							ActionLocation: schema.ActionLocation{
								LocationName:        lRec.Name,
								LocationDescription: lRec.Description,
								LocationDirections:  []string{"north"},
								LocationCharacters: []schema.ActionLocationCharacter{
									{
										CharacterName: cRec.Name,
									},
									{
										CharacterName: lcRec.Name,
									},
								},
								LocationMonsters: []schema.ActionLocationMonster{
									{
										MonsterName: lmRec.Name,
									},
								},
								LocationObjects: []schema.ActionLocationObject{},
							},
							ActionCharacter: &schema.ActionCharacter{
								CharacterName:                cRec.Name,
								CharacterStrength:            cRec.Strength,
								CharacterDexterity:           cRec.Dexterity,
								CharacterIntelligence:        cRec.Intelligence,
								CharacterHealth:              cRec.Health,
								CharacterFatigue:             cRec.Fatigue,
								CharacterCurrentStrength:     ciRec.Strength,
								CharacterCurrentDexterity:    ciRec.Dexterity,
								CharacterCurrentIntelligence: ciRec.Intelligence,
								CharacterCurrentHealth:       ciRec.Health,
								CharacterCurrentFatigue:      ciRec.Fatigue,
								CharacterEquippedObjects: []schema.ActionObject{
									{
										ObjectName: toRec.Name,
									},
									{
										ObjectName: ceoRec.Name,
									},
								},
								CharacterStashedObjects: []schema.ActionObject{
									{
										ObjectName: csoRec.Name,
									},
								},
							},
							ActionMonster: nil,
							ActionEquippedObject: &schema.ActionObject{
								ObjectName:        toRec.Name,
								ObjectDescription: toRec.Description,
								ObjectIsEquipped:  true,
								ObjectIsStashed:   false,
							},
							ActionStashedObject: nil,
							ActionDroppedObject: nil,
							ActionTargetObject: &schema.ActionObject{
								ObjectName:        toRec.Name,
								ObjectDescription: toRec.Description,
								ObjectIsEquipped:  true,
								ObjectIsStashed:   false,
							},
							ActionTargetCharacter: nil,
							ActionTargetMonster:   nil,
							ActionTargetLocation:  nil,
						},
					},
				}
				return &res
			},
		},
		{
			// Drop object that is equipped
			TestCase: TestCase{
				Name:              "Drop object that is equipped",
				HandlerConfig:     testCaseHandlerConfig,
				RequestHeaders:    testCaseRequestHeaders,
				RequestPathParams: testCaseRequestPathParams,
				RequestBody: func(data harness.Data) interface{} {
					toRec, _ := data.GetObjectRecByName("dull bronze ring")
					res := schema.ActionRequest{
						Data: schema.ActionRequestData{
							Sentence: fmt.Sprintf("drop %s", toRec.Name),
						},
					}
					return &res
				},
				ResponseBody: testCaseResponseBody,
				ResponseCode: http.StatusOK,
			},
			expectResponseBody: func(data harness.Data) *schema.ActionResponse {

				cRec, _ := data.GetCharacterRecByName("barricade")
				csoRec, _ := data.GetObjectRecByName("blood stained pouch")
				ciRec, _ := data.GetCharacterInstanceRecByName("barricade")
				lRec, _ := data.GetLocationRecByName("cave entrance")
				lcRec, _ := data.GetCharacterRecByName("legislate")
				loRec, _ := data.GetObjectRecByName("rusted sword")
				lmRec, _ := data.GetMonsterRecByName("grumpy dwarf")
				toRec, _ := data.GetObjectRecByName("dull bronze ring")

				res := schema.ActionResponse{
					Data: []schema.ActionResponseData{
						{
							ActionCommand:   "drop",
							ActionNarrative: fmt.Sprintf("%s drops %s", cRec.Name, toRec.Name),
							ActionLocation: schema.ActionLocation{
								LocationName:        lRec.Name,
								LocationDescription: lRec.Description,
								LocationDirections:  []string{"north"},
								LocationCharacters: []schema.ActionLocationCharacter{
									{
										CharacterName: cRec.Name,
									},
									{
										CharacterName: lcRec.Name,
									},
								},
								LocationMonsters: []schema.ActionLocationMonster{
									{
										MonsterName: lmRec.Name,
									},
								},
								LocationObjects: []schema.ActionLocationObject{
									{
										ObjectName: loRec.Name,
									},
									{
										ObjectName: toRec.Name,
									},
								},
							},
							ActionCharacter: &schema.ActionCharacter{
								CharacterName:                cRec.Name,
								CharacterStrength:            cRec.Strength,
								CharacterDexterity:           cRec.Dexterity,
								CharacterIntelligence:        cRec.Intelligence,
								CharacterHealth:              cRec.Health,
								CharacterFatigue:             cRec.Fatigue,
								CharacterCurrentStrength:     ciRec.Strength,
								CharacterCurrentDexterity:    ciRec.Dexterity,
								CharacterCurrentIntelligence: ciRec.Intelligence,
								CharacterCurrentHealth:       ciRec.Health,
								CharacterCurrentFatigue:      ciRec.Fatigue,
								CharacterEquippedObjects:     []schema.ActionObject{},
								CharacterStashedObjects: []schema.ActionObject{
									{
										ObjectName: csoRec.Name,
									},
								},
							},
							ActionMonster:        nil,
							ActionEquippedObject: nil,
							ActionStashedObject:  nil,
							ActionDroppedObject: &schema.ActionObject{
								ObjectName:        toRec.Name,
								ObjectDescription: toRec.Description,
								ObjectIsEquipped:  false,
								ObjectIsStashed:   false,
							},
							ActionTargetObject: &schema.ActionObject{
								ObjectName:        toRec.Name,
								ObjectDescription: toRec.Description,
								ObjectIsEquipped:  false,
								ObjectIsStashed:   false,
							},
							ActionTargetCharacter: nil,
							ActionTargetMonster:   nil,
							ActionTargetLocation:  nil,
						},
					},
				}
				return &res
			},
		},
		{
			// Submit nothing
			TestCase: TestCase{
				Name:              "Submit nothing",
				HandlerConfig:     testCaseHandlerConfig,
				RequestHeaders:    testCaseRequestHeaders,
				RequestPathParams: testCaseRequestPathParams,
				RequestBody: func(data harness.Data) interface{} {
					res := schema.ActionRequest{
						Data: schema.ActionRequestData{
							Sentence: "",
						},
					}
					return &res
				},
				ResponseBody: testCaseResponseBody,
				ResponseCode: http.StatusBadRequest,
			},
		},
	}

	isCharacterNil := func(c *schema.ActionCharacter) bool {
		return c == nil
	}
	isMonsterNil := func(m *schema.ActionMonster) bool {
		return m == nil
	}
	isObjectNil := func(o *schema.ActionObject) bool {
		return o == nil
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Logf("Running test >%s<", testCase.Name)

			testFunc := func(method string, body interface{}) {

				if testCase.TestResponseCode() != http.StatusOK {
					return
				}

				var responseBody *schema.ActionResponse
				if body != nil {
					responseBody = body.(*schema.ActionResponse)
				}

				// Validate response body
				if testCase.expectResponseBody != nil {
					require.NotNil(t, responseBody, "Response body is not nil")
					require.GreaterOrEqual(t, len(responseBody.Data), 0, "Response body data ")

					expectResponseBody := testCase.expectResponseBody(th.Data)

					// Validate response body data
					for idx, expectData := range expectResponseBody.Data {
						require.NotNil(t, responseBody.Data[idx], "Response body index is not empty")

						// Command
						require.Equal(t, expectData.ActionCommand, responseBody.Data[idx].ActionCommand, "Response action command equals expected")

						// Narrative
						require.Equal(t, expectData.ActionNarrative, responseBody.Data[idx].ActionNarrative, "Response action narrative equals expected")

						// Current location
						t.Logf("Checking location name >%s< >%s<", expectData.ActionLocation.LocationName, responseBody.Data[idx].ActionLocation.LocationName)
						require.Equal(t, expectData.ActionLocation.LocationName, responseBody.Data[idx].ActionLocation.LocationName)
						t.Logf("Checking location description >%s< >%s<", expectData.ActionLocation.LocationDescription, responseBody.Data[idx].ActionLocation.LocationDescription)
						require.Equal(t, expectData.ActionLocation.LocationDescription, responseBody.Data[idx].ActionLocation.LocationDescription)
						t.Logf("Checking location directions >%s< >%s<", expectData.ActionLocation.LocationDirections, responseBody.Data[idx].ActionLocation.LocationDirections)
						require.Equal(t, expectData.ActionLocation.LocationDirections, responseBody.Data[idx].ActionLocation.LocationDirections)

						// Current location characters
						t.Logf("Checking character count >%d< >%d<", len(expectData.ActionLocation.LocationCharacters), len(responseBody.Data[idx].ActionLocation.LocationCharacters))
						require.Equal(t, len(expectData.ActionLocation.LocationCharacters), len(responseBody.Data[idx].ActionLocation.LocationCharacters), "Response action location characters count equals expected")
						if len(expectData.ActionLocation.LocationCharacters) > 0 {
							for _, character := range expectData.ActionLocation.LocationCharacters {
								t.Logf("Checking action location character name >%s<", character.CharacterName)
								require.True(t, responseBody.Data[idx].ActionLocation.LocationCharacters.HasCharacterWithName(character.CharacterName), "Response action location characters has character with name ")
							}
						}
						if len(expectData.ActionLocation.LocationCharacters) == 0 {
							require.Equal(t, 0, len(responseBody.Data[idx].ActionLocation.LocationCharacters), "Location characters length is 0")
						}

						// Current location monsters
						t.Logf("Checking monster count >%d< >%d<", len(expectData.ActionLocation.LocationMonsters), len(responseBody.Data[idx].ActionLocation.LocationMonsters))
						require.Equal(t, len(expectData.ActionLocation.LocationMonsters), len(responseBody.Data[idx].ActionLocation.LocationMonsters), "Response action location monsters count equals expected")
						if len(expectData.ActionLocation.LocationMonsters) > 0 {
							for _, monster := range expectData.ActionLocation.LocationMonsters {
								t.Logf("Checking action location monster name >%s<", monster.MonsterName)
								require.True(t, responseBody.Data[idx].ActionLocation.LocationMonsters.HasMonsterWithName(monster.MonsterName))
							}
						}
						if len(expectData.ActionLocation.LocationMonsters) == 0 {
							require.Equal(t, 0, len(responseBody.Data[idx].ActionLocation.LocationMonsters), "Location monsters length is 0")
						}

						// Current location objects (any order)
						t.Logf("Checking object count >%d< >%d<", len(expectData.ActionLocation.LocationObjects), len(responseBody.Data[idx].ActionLocation.LocationObjects))
						require.Equal(t, len(expectData.ActionLocation.LocationObjects), len(responseBody.Data[idx].ActionLocation.LocationObjects), "Response action location objects count equals expected")
						if len(expectData.ActionLocation.LocationObjects) > 0 {
							for _, object := range expectData.ActionLocation.LocationObjects {
								t.Logf("Checking action location object name >%s<", object.ObjectName)
								require.True(t, responseBody.Data[idx].ActionLocation.LocationObjects.HasObjectWithName(object.ObjectName))
							}
						}

						// Target location
						if expectData.ActionTargetLocation != nil {
							require.NotNil(t, responseBody.Data[idx].ActionTargetLocation, "Response target location is not empty")
							t.Logf("Checking location name >%s< >%s<", expectData.ActionTargetLocation.LocationName, responseBody.Data[idx].ActionTargetLocation.LocationName)
							require.Equal(t, expectData.ActionTargetLocation.LocationName, responseBody.Data[idx].ActionTargetLocation.LocationName)
							t.Logf("Checking location description >%s< >%s<", expectData.ActionTargetLocation.LocationDescription, responseBody.Data[idx].ActionTargetLocation.LocationDescription)
							require.Equal(t, expectData.ActionTargetLocation.LocationDescription, responseBody.Data[idx].ActionTargetLocation.LocationDescription)
							t.Logf("Checking location direction >%s< >%s<", expectData.ActionTargetLocation.LocationDirection, responseBody.Data[idx].ActionTargetLocation.LocationDirection)
							require.Equal(t, expectData.ActionTargetLocation.LocationDirection, responseBody.Data[idx].ActionTargetLocation.LocationDirection)
							t.Logf("Checking location directions >%s< >%s<", expectData.ActionTargetLocation.LocationDirections, responseBody.Data[idx].ActionTargetLocation.LocationDirections)
							require.Equal(t, expectData.ActionTargetLocation.LocationDirections, responseBody.Data[idx].ActionTargetLocation.LocationDirections)

							// Target location characters
							t.Logf("Checking character count >%d< >%d<", len(expectData.ActionTargetLocation.LocationCharacters), len(responseBody.Data[idx].ActionTargetLocation.LocationCharacters))
							require.Equal(t, len(expectData.ActionTargetLocation.LocationCharacters), len(responseBody.Data[idx].ActionTargetLocation.LocationCharacters), "Response action target location characters count equals expected")
							if len(expectData.ActionTargetLocation.LocationCharacters) > 0 {
								for _, character := range expectData.ActionTargetLocation.LocationCharacters {
									t.Logf("Checking action target location character name >%s<", character.CharacterName)
									require.True(t, responseBody.Data[idx].ActionTargetLocation.LocationCharacters.HasCharacterWithName(character.CharacterName))
								}
							}

							// Target location monsters
							t.Logf("Checking monster count >%d< >%d<", len(expectData.ActionTargetLocation.LocationMonsters), len(responseBody.Data[idx].ActionTargetLocation.LocationMonsters))
							require.Equal(t, len(expectData.ActionTargetLocation.LocationMonsters), len(responseBody.Data[idx].ActionTargetLocation.LocationMonsters), "Response action target location monsters count equals expected")
							if len(expectData.ActionTargetLocation.LocationMonsters) > 0 {
								for _, monster := range expectData.ActionTargetLocation.LocationMonsters {
									t.Logf("Checking action target location monster name >%s<", monster.MonsterName)
									require.True(t, responseBody.Data[idx].ActionTargetLocation.LocationMonsters.HasMonsterWithName(monster.MonsterName))
								}
							}

							// Target location objects
							t.Logf("Checking object count >%d< >%d<", len(expectData.ActionTargetLocation.LocationObjects), len(responseBody.Data[idx].ActionTargetLocation.LocationObjects))
							require.Equal(t, len(expectData.ActionTargetLocation.LocationObjects), len(responseBody.Data[idx].ActionTargetLocation.LocationObjects), "Response action target location objects count equals expected")
							if len(expectData.ActionTargetLocation.LocationObjects) > 0 {
								for _, object := range expectData.ActionTargetLocation.LocationObjects {
									t.Logf("Checking action target location object name >%s<", object.ObjectName)
									require.True(t, responseBody.Data[idx].ActionTargetLocation.LocationObjects.HasObjectWithName(object.ObjectName))
								}
							}
						}

						// Response character
						t.Logf("Checking character nil >%t< >%t<", isCharacterNil(expectData.ActionCharacter), isCharacterNil(responseBody.Data[idx].ActionCharacter))
						require.Equal(t, isCharacterNil(expectData.ActionCharacter), isCharacterNil(responseBody.Data[idx].ActionCharacter), "Response character is nil or not nil as expected")
						if !isCharacterNil(expectData.ActionCharacter) {
							t.Logf("Checking action character name >%s< >%s<", expectData.ActionCharacter.CharacterName, responseBody.Data[idx].ActionCharacter.CharacterName)
							require.NotEmpty(t, responseBody.Data[idx].ActionCharacter.CharacterName, "Response character name is not empty")
							require.Equal(t, expectData.ActionCharacter.CharacterName, responseBody.Data[idx].ActionCharacter.CharacterName, "Response character name equals expected")
							t.Logf("Checking action character strenth >%d< >%d<", expectData.ActionCharacter.CharacterStrength, responseBody.Data[idx].ActionCharacter.CharacterStrength)
							require.NotEmpty(t, responseBody.Data[idx].ActionCharacter.CharacterStrength, "Response character strenth is not empty")
							require.Equal(t, expectData.ActionCharacter.CharacterStrength, responseBody.Data[idx].ActionCharacter.CharacterStrength, "Response character strength equals expected")
							t.Logf("Checking action character dexterity >%d< >%d<", expectData.ActionCharacter.CharacterDexterity, responseBody.Data[idx].ActionCharacter.CharacterDexterity)
							require.NotEmpty(t, responseBody.Data[idx].ActionCharacter.CharacterDexterity, "Response character dexterity is not empty")
							require.Equal(t, expectData.ActionCharacter.CharacterDexterity, responseBody.Data[idx].ActionCharacter.CharacterDexterity, "Response character dexterity equals expected")
							t.Logf("Checking action character intelligence >%d< >%d<", expectData.ActionCharacter.CharacterIntelligence, responseBody.Data[idx].ActionCharacter.CharacterIntelligence)
							require.NotEmpty(t, responseBody.Data[idx].ActionCharacter.CharacterIntelligence, "Response character intelligence is not empty")
							require.Equal(t, expectData.ActionCharacter.CharacterIntelligence, responseBody.Data[idx].ActionCharacter.CharacterIntelligence, "Response character intelligence equals expected")

							t.Logf("Checking action character current strenth >%d< >%d<", expectData.ActionCharacter.CharacterCurrentStrength, responseBody.Data[idx].ActionCharacter.CharacterCurrentStrength)
							require.NotEmpty(t, responseBody.Data[idx].ActionCharacter.CharacterCurrentStrength, "Response character current strenth is not empty")
							require.Equal(t, expectData.ActionCharacter.CharacterCurrentStrength, responseBody.Data[idx].ActionCharacter.CharacterCurrentStrength, "Response character current strength equals expected")
							t.Logf("Checking action character current dexterity >%d< >%d<", expectData.ActionCharacter.CharacterCurrentDexterity, responseBody.Data[idx].ActionCharacter.CharacterCurrentDexterity)
							require.NotEmpty(t, responseBody.Data[idx].ActionCharacter.CharacterCurrentDexterity, "Response character dexterity is not empty")
							require.Equal(t, expectData.ActionCharacter.CharacterCurrentDexterity, responseBody.Data[idx].ActionCharacter.CharacterCurrentDexterity, "Response character current dexterity equals expected")
							t.Logf("Checking action character current intelligence >%d< >%d<", expectData.ActionCharacter.CharacterCurrentIntelligence, responseBody.Data[idx].ActionCharacter.CharacterCurrentIntelligence)
							require.NotEmpty(t, responseBody.Data[idx].ActionCharacter.CharacterCurrentIntelligence, "Response character intelligence is not empty")
							require.Equal(t, expectData.ActionCharacter.CharacterCurrentIntelligence, responseBody.Data[idx].ActionCharacter.CharacterCurrentIntelligence, "Response character current intelligence equals expected")

							t.Logf("Checking action character health >%d< >%d<", expectData.ActionCharacter.CharacterHealth, responseBody.Data[idx].ActionCharacter.CharacterHealth)
							require.Equal(t, expectData.ActionCharacter.CharacterHealth, responseBody.Data[idx].ActionCharacter.CharacterHealth, "Response character health equals expected")
							t.Logf("Checking action character fatigue >%d< >%d<", expectData.ActionCharacter.CharacterFatigue, responseBody.Data[idx].ActionCharacter.CharacterFatigue)
							require.Equal(t, expectData.ActionCharacter.CharacterFatigue, responseBody.Data[idx].ActionCharacter.CharacterFatigue, "Response character fatigue equals expected")

							t.Logf("Checking action character current health >%d< >%d<", expectData.ActionCharacter.CharacterCurrentHealth, responseBody.Data[idx].ActionCharacter.CharacterCurrentHealth)
							require.Equal(t, expectData.ActionCharacter.CharacterCurrentHealth, responseBody.Data[idx].ActionCharacter.CharacterCurrentHealth, "Response character current health equals expected")
							t.Logf("Checking action character current fatigue >%d< >%d<", expectData.ActionCharacter.CharacterCurrentFatigue, responseBody.Data[idx].ActionCharacter.CharacterCurrentFatigue)
							require.Equal(t, expectData.ActionCharacter.CharacterCurrentFatigue, responseBody.Data[idx].ActionCharacter.CharacterCurrentFatigue, "Response character current fatigue equals expected")

							t.Logf("Checking action character equipped objects >%d< >%d<", len(expectData.ActionCharacter.CharacterEquippedObjects), len(responseBody.Data[idx].ActionCharacter.CharacterEquippedObjects))
							require.Equal(t, len(expectData.ActionCharacter.CharacterEquippedObjects), len(responseBody.Data[idx].ActionCharacter.CharacterEquippedObjects), "Response character equipped object count equals expected")

							t.Logf("Checking action character stashed objects >%d< >%d<", len(expectData.ActionCharacter.CharacterStashedObjects), len(responseBody.Data[idx].ActionCharacter.CharacterStashedObjects))
							require.Equal(t, len(expectData.ActionCharacter.CharacterStashedObjects), len(responseBody.Data[idx].ActionCharacter.CharacterStashedObjects), "Response character stashed object count equals expected")
						}

						// Response target character
						t.Logf("Checking target character nil >%t< >%t<", isCharacterNil(expectData.ActionTargetCharacter), isCharacterNil(responseBody.Data[idx].ActionTargetCharacter))
						require.Equal(t, isCharacterNil(expectData.ActionTargetCharacter), isCharacterNil(responseBody.Data[idx].ActionTargetCharacter), "Response target character is nil or not nil as expected")
						if !isCharacterNil(expectData.ActionTargetCharacter) {
							t.Logf("Checking action target character name >%s< >%s<", expectData.ActionTargetCharacter.CharacterName, responseBody.Data[idx].ActionTargetCharacter.CharacterName)
							require.Equal(t, expectData.ActionTargetCharacter.CharacterName, responseBody.Data[idx].ActionTargetCharacter.CharacterName, "Response target character name equals expected")

							t.Logf("Checking action target character strenth >%d< >%d<", expectData.ActionTargetCharacter.CharacterStrength, responseBody.Data[idx].ActionTargetCharacter.CharacterStrength)
							require.NotEmpty(t, responseBody.Data[idx].ActionTargetCharacter.CharacterStrength, "Response target character strength is not empty")
							require.Equal(t, expectData.ActionTargetCharacter.CharacterStrength, responseBody.Data[idx].ActionTargetCharacter.CharacterStrength, "Response target character strength equals expected")
							t.Logf("Checking action target character dexterity >%d< >%d<", expectData.ActionTargetCharacter.CharacterDexterity, responseBody.Data[idx].ActionTargetCharacter.CharacterDexterity)
							require.NotEmpty(t, responseBody.Data[idx].ActionTargetCharacter.CharacterDexterity, "Response target character dexterity is not empty")
							require.Equal(t, expectData.ActionTargetCharacter.CharacterDexterity, responseBody.Data[idx].ActionTargetCharacter.CharacterDexterity, "Response target character dexterity equals expected")
							t.Logf("Checking action target character intelligence >%d< >%d<", expectData.ActionTargetCharacter.CharacterIntelligence, responseBody.Data[idx].ActionTargetCharacter.CharacterIntelligence)
							require.NotEmpty(t, responseBody.Data[idx].ActionTargetCharacter.CharacterIntelligence, "Response target character intelligence is not empty")
							require.Equal(t, expectData.ActionTargetCharacter.CharacterIntelligence, responseBody.Data[idx].ActionTargetCharacter.CharacterIntelligence, "Response target character intelligence equals expected")

							t.Logf("Checking target character current strenth >%d< >%d<", expectData.ActionTargetCharacter.CharacterCurrentStrength, responseBody.Data[idx].ActionTargetCharacter.CharacterCurrentStrength)
							require.NotEmpty(t, responseBody.Data[idx].ActionTargetCharacter.CharacterCurrentStrength, "Response target character current strength is not empty")
							require.Equal(t, expectData.ActionTargetCharacter.CharacterCurrentStrength, responseBody.Data[idx].ActionTargetCharacter.CharacterCurrentStrength, "Response target character current strength equals expected")
							t.Logf("Checking target character current dexterity >%d< >%d<", expectData.ActionTargetCharacter.CharacterCurrentDexterity, responseBody.Data[idx].ActionTargetCharacter.CharacterCurrentDexterity)
							require.NotEmpty(t, responseBody.Data[idx].ActionTargetCharacter.CharacterCurrentDexterity, "Response target character current dexterity is not empty")
							require.Equal(t, expectData.ActionTargetCharacter.CharacterCurrentDexterity, responseBody.Data[idx].ActionTargetCharacter.CharacterCurrentDexterity, "Response target character current dexterity equals expected")
							t.Logf("Checking target character current intelligence >%d< >%d<", expectData.ActionTargetCharacter.CharacterCurrentIntelligence, responseBody.Data[idx].ActionTargetCharacter.CharacterCurrentIntelligence)
							require.NotEmpty(t, responseBody.Data[idx].ActionTargetCharacter.CharacterCurrentIntelligence, "Response target character current intelligence is not empty")
							require.Equal(t, expectData.ActionTargetCharacter.CharacterCurrentIntelligence, responseBody.Data[idx].ActionTargetCharacter.CharacterCurrentIntelligence, "Response target character current intelligence equals expected")

							t.Logf("Checking target character health >%d< >%d<", expectData.ActionTargetCharacter.CharacterHealth, responseBody.Data[idx].ActionTargetCharacter.CharacterHealth)
							require.Equal(t, expectData.ActionTargetCharacter.CharacterHealth, responseBody.Data[idx].ActionTargetCharacter.CharacterHealth, "Response target character health equals expected")
							t.Logf("Checking target character fatigue >%d< >%d<", expectData.ActionTargetCharacter.CharacterFatigue, responseBody.Data[idx].ActionTargetCharacter.CharacterFatigue)
							require.Equal(t, expectData.ActionTargetCharacter.CharacterFatigue, responseBody.Data[idx].ActionTargetCharacter.CharacterFatigue, "Response target character fatigue equals expected")

							t.Logf("Checking target character current health >%d< >%d<", expectData.ActionTargetCharacter.CharacterCurrentHealth, responseBody.Data[idx].ActionTargetCharacter.CharacterCurrentHealth)
							require.Equal(t, expectData.ActionTargetCharacter.CharacterCurrentHealth, responseBody.Data[idx].ActionTargetCharacter.CharacterCurrentHealth, "Response target character current health equals expected")
							t.Logf("Checking target character current fatigue >%d< >%d<", expectData.ActionTargetCharacter.CharacterCurrentFatigue, responseBody.Data[idx].ActionTargetCharacter.CharacterCurrentFatigue)
							require.Equal(t, expectData.ActionTargetCharacter.CharacterCurrentFatigue, responseBody.Data[idx].ActionTargetCharacter.CharacterCurrentFatigue, "Response target character current fatigue equals expected")

							t.Logf("Checking target character equipped objects >%d< >%d<", len(expectData.ActionTargetCharacter.CharacterEquippedObjects), len(responseBody.Data[idx].ActionTargetCharacter.CharacterEquippedObjects))
							require.Equal(t, len(expectData.ActionTargetCharacter.CharacterEquippedObjects), len(responseBody.Data[idx].ActionTargetCharacter.CharacterEquippedObjects), "Response target character equipped object count equals expected")

							t.Logf("Checking target character stashed objects >%d< >%d<", len(expectData.ActionTargetCharacter.CharacterStashedObjects), len(responseBody.Data[idx].ActionTargetCharacter.CharacterStashedObjects))
							require.Equal(t, len(expectData.ActionTargetCharacter.CharacterStashedObjects), len(responseBody.Data[idx].ActionTargetCharacter.CharacterStashedObjects), "Response target character stashed object count equals expected")
						}

						// Response monster
						t.Logf("Checking monster nil >%t< >%t<", isMonsterNil(expectData.ActionMonster), isMonsterNil(responseBody.Data[idx].ActionMonster))
						require.Equal(t, isMonsterNil(expectData.ActionMonster), isMonsterNil(responseBody.Data[idx].ActionMonster), "Response monster is nil or not nil as expected")
						if !isMonsterNil(expectData.ActionMonster) {
							t.Logf("Checking monster name >%s< >%s<", expectData.ActionMonster.MonsterName, responseBody.Data[idx].ActionMonster.MonsterName)
							require.Equal(t, expectData.ActionMonster.MonsterName, responseBody.Data[idx].ActionMonster.MonsterName, "Response monster name equals expected")

							t.Logf("Checking monster strenth >%d< >%d<", expectData.ActionMonster.MonsterStrength, responseBody.Data[idx].ActionMonster.MonsterStrength)
							require.NotEmpty(t, responseBody.Data[idx].ActionMonster.MonsterStrength, "Response monster strength is not empty")
							require.Equal(t, expectData.ActionMonster.MonsterStrength, responseBody.Data[idx].ActionMonster.MonsterStrength, "Response monster strength equals expected")
							t.Logf("Checking monster dexterity >%d< >%d<", expectData.ActionMonster.MonsterDexterity, responseBody.Data[idx].ActionMonster.MonsterDexterity)
							require.NotEmpty(t, responseBody.Data[idx].ActionMonster.MonsterDexterity, "Response monster dexterity is not empty")
							require.Equal(t, expectData.ActionMonster.MonsterDexterity, responseBody.Data[idx].ActionMonster.MonsterDexterity, "Response monster dexterity equals expected")
							t.Logf("Checking monster intelligence >%d< >%d<", expectData.ActionMonster.MonsterIntelligence, responseBody.Data[idx].ActionMonster.MonsterIntelligence)
							require.NotEmpty(t, responseBody.Data[idx].ActionMonster.MonsterIntelligence, "Response monster intelligence is not empty")
							require.Equal(t, expectData.ActionMonster.MonsterIntelligence, responseBody.Data[idx].ActionMonster.MonsterIntelligence, "Response monster intelligence equals expected")

							t.Logf("Checking monster current strenth >%d< >%d<", expectData.ActionMonster.MonsterCurrentStrength, responseBody.Data[idx].ActionMonster.MonsterCurrentStrength)
							require.NotEmpty(t, responseBody.Data[idx].ActionMonster.MonsterCurrentStrength, "Response monster current strength is not empty")
							require.Equal(t, expectData.ActionMonster.MonsterCurrentStrength, responseBody.Data[idx].ActionMonster.MonsterCurrentStrength, "Response monster current strength equals expected")
							t.Logf("Checking monster current dexterity >%d< >%d<", expectData.ActionMonster.MonsterCurrentDexterity, responseBody.Data[idx].ActionMonster.MonsterCurrentDexterity)
							require.NotEmpty(t, responseBody.Data[idx].ActionMonster.MonsterCurrentDexterity, "Response monster current dexterity is not empty")
							require.Equal(t, expectData.ActionMonster.MonsterCurrentDexterity, responseBody.Data[idx].ActionMonster.MonsterCurrentDexterity, "Response monster current dexterity equals expected")
							t.Logf("Checking monster current intelligence >%d< >%d<", expectData.ActionMonster.MonsterCurrentIntelligence, responseBody.Data[idx].ActionMonster.MonsterCurrentIntelligence)
							require.NotEmpty(t, responseBody.Data[idx].ActionMonster.MonsterCurrentIntelligence, "Response monster current intelligence is not empty")
							require.Equal(t, expectData.ActionMonster.MonsterCurrentIntelligence, responseBody.Data[idx].ActionMonster.MonsterCurrentIntelligence, "Response monster current intelligence equals expected")

							t.Logf("Checking monster health >%d< >%d<", expectData.ActionMonster.MonsterHealth, responseBody.Data[idx].ActionMonster.MonsterHealth)
							require.Equal(t, expectData.ActionMonster.MonsterHealth, responseBody.Data[idx].ActionMonster.MonsterHealth, "Response monster health equals expected")
							t.Logf("Checking monster fatigue >%d< >%d<", expectData.ActionMonster.MonsterFatigue, responseBody.Data[idx].ActionMonster.MonsterFatigue)
							require.Equal(t, expectData.ActionMonster.MonsterFatigue, responseBody.Data[idx].ActionMonster.MonsterFatigue, "Response monster fatigue equals expected")

							t.Logf("Checking monster current health >%d< >%d<", expectData.ActionMonster.MonsterCurrentHealth, responseBody.Data[idx].ActionMonster.MonsterCurrentHealth)
							require.Equal(t, expectData.ActionMonster.MonsterCurrentHealth, responseBody.Data[idx].ActionMonster.MonsterCurrentHealth, "Response monster current health equals expected")
							t.Logf("Checking monster current fatigue >%d< >%d<", expectData.ActionMonster.MonsterCurrentFatigue, responseBody.Data[idx].ActionMonster.MonsterCurrentFatigue)
							require.Equal(t, expectData.ActionMonster.MonsterCurrentFatigue, responseBody.Data[idx].ActionMonster.MonsterCurrentFatigue, "Response monster current fatigue equals expected")
						}

						// Response target monster
						t.Logf("Checking target monster nil >%t< >%t<", isMonsterNil(expectData.ActionTargetMonster), isMonsterNil(responseBody.Data[idx].ActionTargetMonster))
						require.Equal(t, isMonsterNil(expectData.ActionTargetMonster), isMonsterNil(responseBody.Data[idx].ActionTargetMonster), "Response target monster is nil or not nil as expected")
						if !isMonsterNil(expectData.ActionTargetMonster) {
							t.Logf("Checking target monster name >%s< >%s<", expectData.ActionTargetMonster.MonsterName, responseBody.Data[idx].ActionTargetMonster.MonsterName)
							require.Equal(t, expectData.ActionTargetMonster.MonsterName, responseBody.Data[idx].ActionTargetMonster.MonsterName, "Response target monster name equals expected")

							t.Logf("Checking target monster strenth >%d< >%d<", expectData.ActionTargetMonster.MonsterStrength, responseBody.Data[idx].ActionTargetMonster.MonsterStrength)
							require.NotEmpty(t, responseBody.Data[idx].ActionTargetMonster.MonsterStrength, "Response target monster strength is not empty")
							require.Equal(t, expectData.ActionTargetMonster.MonsterStrength, responseBody.Data[idx].ActionTargetMonster.MonsterStrength, "Response target monster strength equals expected")
							t.Logf("Checking target monster dexterity >%d< >%d<", expectData.ActionTargetMonster.MonsterDexterity, responseBody.Data[idx].ActionTargetMonster.MonsterDexterity)
							require.NotEmpty(t, responseBody.Data[idx].ActionTargetMonster.MonsterDexterity, "Response target monster dexterity is not empty")
							require.Equal(t, expectData.ActionTargetMonster.MonsterDexterity, responseBody.Data[idx].ActionTargetMonster.MonsterDexterity, "Response target monster dexterity equals expected")
							t.Logf("Checking target monster intelligence >%d< >%d<", expectData.ActionTargetMonster.MonsterIntelligence, responseBody.Data[idx].ActionTargetMonster.MonsterIntelligence)
							require.NotEmpty(t, responseBody.Data[idx].ActionTargetMonster.MonsterIntelligence, "Response target monster intelligence is not empty")
							require.Equal(t, expectData.ActionTargetMonster.MonsterIntelligence, responseBody.Data[idx].ActionTargetMonster.MonsterIntelligence, "Response target monster intelligence equals expected")

							t.Logf("Checking target monster current strenth >%d< >%d<", expectData.ActionTargetMonster.MonsterCurrentStrength, responseBody.Data[idx].ActionTargetMonster.MonsterCurrentStrength)
							require.NotEmpty(t, responseBody.Data[idx].ActionTargetMonster.MonsterCurrentStrength, "Response target monster current strength is not empty")
							require.Equal(t, expectData.ActionTargetMonster.MonsterCurrentStrength, responseBody.Data[idx].ActionTargetMonster.MonsterCurrentStrength, "Response target monster current strength equals expected")
							t.Logf("Checking target monster current dexterity >%d< >%d<", expectData.ActionTargetMonster.MonsterCurrentDexterity, responseBody.Data[idx].ActionTargetMonster.MonsterCurrentDexterity)
							require.NotEmpty(t, responseBody.Data[idx].ActionTargetMonster.MonsterCurrentDexterity, "Response target monster current dexterity is not empty")
							require.Equal(t, expectData.ActionTargetMonster.MonsterCurrentDexterity, responseBody.Data[idx].ActionTargetMonster.MonsterCurrentDexterity, "Response target monster current dexterity equals expected")
							t.Logf("Checking target monster current intelligence >%d< >%d<", expectData.ActionTargetMonster.MonsterCurrentIntelligence, responseBody.Data[idx].ActionTargetMonster.MonsterCurrentIntelligence)
							require.NotEmpty(t, responseBody.Data[idx].ActionTargetMonster.MonsterCurrentIntelligence, "Response target monster current intelligence is not empty")
							require.Equal(t, expectData.ActionTargetMonster.MonsterCurrentIntelligence, responseBody.Data[idx].ActionTargetMonster.MonsterCurrentIntelligence, "Response target monster current intelligence equals expected")

							t.Logf("Checking target monster health >%d< >%d<", expectData.ActionTargetMonster.MonsterHealth, responseBody.Data[idx].ActionTargetMonster.MonsterHealth)
							require.Equal(t, expectData.ActionTargetMonster.MonsterHealth, responseBody.Data[idx].ActionTargetMonster.MonsterHealth, "Response target monster health equals expected")
							t.Logf("Checking target monster fatigue >%d< >%d<", expectData.ActionTargetMonster.MonsterFatigue, responseBody.Data[idx].ActionTargetMonster.MonsterFatigue)
							require.Equal(t, expectData.ActionTargetMonster.MonsterFatigue, responseBody.Data[idx].ActionTargetMonster.MonsterFatigue, "Response target monster fatigue equals expected")

							t.Logf("Checking target monster current health >%d< >%d<", expectData.ActionTargetMonster.MonsterCurrentHealth, responseBody.Data[idx].ActionTargetMonster.MonsterCurrentHealth)
							require.Equal(t, expectData.ActionTargetMonster.MonsterCurrentHealth, responseBody.Data[idx].ActionTargetMonster.MonsterCurrentHealth, "Response target monster current health equals expected")
							t.Logf("Checking target monster current fatigue >%d< >%d<", expectData.ActionTargetMonster.MonsterCurrentFatigue, responseBody.Data[idx].ActionTargetMonster.MonsterCurrentFatigue)
							require.Equal(t, expectData.ActionTargetMonster.MonsterCurrentFatigue, responseBody.Data[idx].ActionTargetMonster.MonsterCurrentFatigue, "Response target monster current fatigue equals expected")
						}

						// Response target object
						t.Logf("Checking target object nil >%t< >%t<", isObjectNil(expectData.ActionTargetObject), isObjectNil(responseBody.Data[idx].ActionTargetObject))
						require.Equal(t, isObjectNil(expectData.ActionTargetObject), isObjectNil(responseBody.Data[idx].ActionTargetObject), "Response target object is nil or not nil as expected")
						if !isObjectNil(expectData.ActionTargetObject) {
							t.Logf("Checking target object name >%s< >%s<", expectData.ActionTargetObject.ObjectName, responseBody.Data[idx].ActionTargetObject.ObjectName)
							require.Equal(t, expectData.ActionTargetObject.ObjectName, responseBody.Data[idx].ActionTargetObject.ObjectName, "Response target object name equals expected")
							t.Logf("Checking target object description >%s< >%s<", expectData.ActionTargetObject.ObjectDescription, responseBody.Data[idx].ActionTargetObject.ObjectDescription)
							require.Equal(t, expectData.ActionTargetObject.ObjectDescription, responseBody.Data[idx].ActionTargetObject.ObjectDescription, "Response target object description equals expected")
						}

						// Response stashed object
						t.Logf("Checking stashed object nil >%t< >%t<", isObjectNil(expectData.ActionStashedObject), isObjectNil(responseBody.Data[idx].ActionStashedObject))
						require.Equal(t, isObjectNil(expectData.ActionStashedObject), isObjectNil(responseBody.Data[idx].ActionStashedObject), "Response stashed object is nil or not nil as expected")
						if !isObjectNil(expectData.ActionStashedObject) {
							t.Logf("Checking stashed object name >%s< >%s<", expectData.ActionStashedObject.ObjectName, responseBody.Data[idx].ActionStashedObject.ObjectName)
							require.Equal(t, expectData.ActionStashedObject.ObjectName, responseBody.Data[idx].ActionStashedObject.ObjectName, "Response stashed object name equals expected")
							t.Logf("Checking stashed object description >%s< >%s<", expectData.ActionStashedObject.ObjectDescription, responseBody.Data[idx].ActionStashedObject.ObjectDescription)
							require.Equal(t, expectData.ActionStashedObject.ObjectDescription, responseBody.Data[idx].ActionStashedObject.ObjectDescription, "Response stashed object description equals expected")

							t.Logf("Checking stashed object IsEquipped >%t< >%t<", expectData.ActionStashedObject.ObjectIsEquipped, responseBody.Data[idx].ActionStashedObject.ObjectIsEquipped)
							require.Equal(t, expectData.ActionStashedObject.ObjectIsEquipped, responseBody.Data[idx].ActionStashedObject.ObjectIsEquipped, "Response stashed object IsEquipped equals expected")
							t.Logf("Checking stashed object IsStashed >%t< >%t<", expectData.ActionStashedObject.ObjectIsStashed, responseBody.Data[idx].ActionStashedObject.ObjectIsStashed)
							require.Equal(t, expectData.ActionStashedObject.ObjectIsStashed, responseBody.Data[idx].ActionStashedObject.ObjectIsStashed, "Response stashed object IsStashed equals expected")
						}

						// Response equipped object
						t.Logf("Checking equipped object nil >%t< >%t<", isObjectNil(expectData.ActionEquippedObject), isObjectNil(responseBody.Data[idx].ActionEquippedObject))
						require.Equal(t, isObjectNil(expectData.ActionEquippedObject), isObjectNil(responseBody.Data[idx].ActionEquippedObject), "Response equipped object is nil or not nil as expected")
						if !isObjectNil(expectData.ActionEquippedObject) {
							t.Logf("Checking equipped object name >%s< >%s<", expectData.ActionEquippedObject.ObjectName, responseBody.Data[idx].ActionEquippedObject.ObjectName)
							require.Equal(t, expectData.ActionEquippedObject.ObjectName, responseBody.Data[idx].ActionEquippedObject.ObjectName, "Response equipped object name equals expected")
							t.Logf("Checking equipped object description >%s< >%s<", expectData.ActionEquippedObject.ObjectDescription, responseBody.Data[idx].ActionEquippedObject.ObjectDescription)
							require.Equal(t, expectData.ActionEquippedObject.ObjectDescription, responseBody.Data[idx].ActionEquippedObject.ObjectDescription, "Response equipped object description equals expected")

							t.Logf("Checking equipped object IsEquipped >%t< >%t<", expectData.ActionEquippedObject.ObjectIsEquipped, responseBody.Data[idx].ActionEquippedObject.ObjectIsEquipped)
							require.Equal(t, expectData.ActionEquippedObject.ObjectIsEquipped, responseBody.Data[idx].ActionEquippedObject.ObjectIsEquipped, "Response equipped object IsEquipped equals expected")
							t.Logf("Checking equipped object IsStashed >%t< >%t<", expectData.ActionEquippedObject.ObjectIsStashed, responseBody.Data[idx].ActionEquippedObject.ObjectIsStashed)
							require.Equal(t, expectData.ActionEquippedObject.ObjectIsStashed, responseBody.Data[idx].ActionEquippedObject.ObjectIsStashed, "Response equipped object IsStashed equals expected")
						}

						// Response dropped object
						t.Logf("Checking dropped object nil >%t< >%t<", isObjectNil(expectData.ActionDroppedObject), isObjectNil(responseBody.Data[idx].ActionDroppedObject))
						require.Equal(t, isObjectNil(expectData.ActionDroppedObject), isObjectNil(responseBody.Data[idx].ActionDroppedObject), "Response dropped object is nil or not nil as expected")
						if !isObjectNil(expectData.ActionDroppedObject) {
							t.Logf("Checking dropped object name >%s< >%s<", expectData.ActionDroppedObject.ObjectName, responseBody.Data[idx].ActionDroppedObject.ObjectName)
							require.Equal(t, expectData.ActionDroppedObject.ObjectName, responseBody.Data[idx].ActionDroppedObject.ObjectName, "Response dropped object name equals expected")
							t.Logf("Checking dropped object description >%s< >%s<", expectData.ActionDroppedObject.ObjectDescription, responseBody.Data[idx].ActionDroppedObject.ObjectDescription)
							require.Equal(t, expectData.ActionDroppedObject.ObjectDescription, responseBody.Data[idx].ActionDroppedObject.ObjectDescription, "Response dropped object description equals expected")

							t.Logf("Checking dropped object is_dropped >%t< >%t<", expectData.ActionDroppedObject.ObjectIsEquipped, responseBody.Data[idx].ActionDroppedObject.ObjectIsEquipped)
							require.Equal(t, expectData.ActionDroppedObject.ObjectIsEquipped, responseBody.Data[idx].ActionDroppedObject.ObjectIsEquipped, "Response dropped object is_dropped equals expected")
							t.Logf("Checking dropped object IsStashed >%t< >%t<", expectData.ActionDroppedObject.ObjectIsStashed, responseBody.Data[idx].ActionDroppedObject.ObjectIsStashed)
							require.Equal(t, expectData.ActionDroppedObject.ObjectIsStashed, responseBody.Data[idx].ActionDroppedObject.ObjectIsStashed, "Response dropped object IsStashed equals expected")
						}
					}
				}

				// Check dates and add teardown ID's
				for _, data := range responseBody.Data {
					require.False(t, data.ActionCreatedAt.IsZero(), "CreatedAt is not zero")
					if method == http.MethodPost {
						require.True(t, data.ActionUpdatedAt.IsZero(), "UpdatedAt is zero")
					}
				}
			}

			RunTestCase(t, th, &testCase, testFunc)
		})
	}
}
