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
		}
		return headers
	}

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
				loRec, _ := data.GetObjectRecByName("rusted sword")
				lmRec, _ := data.GetMonsterRecByName("grumpy dwarf")

				res := schema.ActionResponse{
					Data: []schema.ActionResponseData{
						{
							Command:   "look",
							Narrative: fmt.Sprintf("%s looks", cRec.Name),
							Location: schema.ActionLocation{
								Name:        lRec.Name,
								Description: lRec.Description,
								Directions:  []string{"north"},
								Characters: []schema.ActionLocationCharacter{
									{
										Name: cRec.Name,
									},
								},
								Monsters: []schema.ActionLocationMonster{
									{
										Name: lmRec.Name,
									},
								},
								Objects: []schema.ActionLocationObject{
									{
										Name: loRec.Name,
									},
								},
							},
							Character: &schema.ActionCharacter{
								Name:                cRec.Name,
								Strength:            cRec.Strength,
								Dexterity:           cRec.Dexterity,
								Intelligence:        cRec.Intelligence,
								Health:              cRec.Health,
								Fatigue:             cRec.Fatigue,
								CurrentStrength:     ciRec.Strength,
								CurrentDexterity:    ciRec.Dexterity,
								CurrentIntelligence: ciRec.Intelligence,
								CurrentHealth:       ciRec.Health,
								CurrentFatigue:      ciRec.Fatigue,
								EquippedObjects: []schema.ActionObject{
									{
										Name: ceoRec.Name,
									},
								},
								StashedObjects: []schema.ActionObject{
									{
										Name: csoRec.Name,
									},
								},
							},
							Monster:         nil,
							EquippedObject:  nil,
							StashedObject:   nil,
							TargetObject:    nil,
							TargetCharacter: nil,
							TargetMonster:   nil,
							TargetLocation: &schema.ActionLocation{
								Name:        lRec.Name,
								Description: lRec.Description,
								Directions:  []string{"north"},
								Characters: []schema.ActionLocationCharacter{
									{
										Name: cRec.Name,
									},
								},
								Monsters: []schema.ActionLocationMonster{
									{
										Name: lmRec.Name,
									},
								},
								Objects: []schema.ActionLocationObject{
									{
										Name: loRec.Name,
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
							Command:   "move",
							Narrative: fmt.Sprintf("%s moves north", cRec.Name),
							Location: schema.ActionLocation{
								Name:        lRec.Name,
								Description: lRec.Description,
								Directions:  []string{"north", "south", "northwest"},
								Characters: []schema.ActionLocationCharacter{
									{
										Name: cRec.Name,
									},
								},
								Monsters: []schema.ActionLocationMonster{
									{
										Name: lmRec.Name,
									},
								},
								Objects: []schema.ActionLocationObject{
									{
										Name: loRec.Name,
									},
								},
							},
							Character: &schema.ActionCharacter{
								Name:                cRec.Name,
								Strength:            cRec.Strength,
								Dexterity:           cRec.Dexterity,
								Intelligence:        cRec.Intelligence,
								Health:              cRec.Health,
								Fatigue:             cRec.Fatigue,
								CurrentStrength:     ciRec.Strength,
								CurrentDexterity:    ciRec.Dexterity,
								CurrentIntelligence: ciRec.Intelligence,
								CurrentHealth:       ciRec.Health,
								CurrentFatigue:      ciRec.Fatigue,
								EquippedObjects: []schema.ActionObject{
									{
										Name: ceoRec.Name,
									},
								},
								StashedObjects: []schema.ActionObject{
									{
										Name: csoRec.Name,
									},
								},
							},
							Monster:         nil,
							EquippedObject:  nil,
							StashedObject:   nil,
							TargetObject:    nil,
							TargetCharacter: nil,
							TargetMonster:   nil,
							TargetLocation: &schema.ActionLocation{
								Name:        lRec.Name,
								Description: lRec.Description,
								Direction:   "north",
								Directions:  []string{"north", "south", "northwest"},
								Characters: []schema.ActionLocationCharacter{
									{
										Name: cRec.Name,
									},
								},
								Monsters: []schema.ActionLocationMonster{
									{
										Name: lmRec.Name,
									},
								},
								Objects: []schema.ActionLocationObject{
									{
										Name: loRec.Name,
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
				loRec, _ := data.GetObjectRecByName("rusted sword")
				lmRec, _ := data.GetMonsterRecByName("grumpy dwarf")
				tlRec, _ := data.GetLocationRecByName("cave tunnel")
				tloRec, _ := data.GetObjectRecByName("rusted helmet")
				tlmRec, _ := data.GetMonsterRecByName("angry goblin")

				res := schema.ActionResponse{
					Data: []schema.ActionResponseData{
						{
							Command:   "look",
							Narrative: fmt.Sprintf("%s looks north", cRec.Name),
							Location: schema.ActionLocation{
								Name:        lRec.Name,
								Description: lRec.Description,
								Directions:  []string{"north"},
								Characters: []schema.ActionLocationCharacter{
									{
										Name: cRec.Name,
									},
								},
								Monsters: []schema.ActionLocationMonster{
									{
										Name: lmRec.Name,
									},
								},
								Objects: []schema.ActionLocationObject{
									{
										Name: loRec.Name,
									},
								},
							},
							Character: &schema.ActionCharacter{
								Name:                cRec.Name,
								Strength:            cRec.Strength,
								Dexterity:           cRec.Dexterity,
								Intelligence:        cRec.Intelligence,
								Health:              cRec.Health,
								Fatigue:             cRec.Fatigue,
								CurrentStrength:     ciRec.Strength,
								CurrentDexterity:    ciRec.Dexterity,
								CurrentIntelligence: ciRec.Intelligence,
								CurrentHealth:       ciRec.Health,
								CurrentFatigue:      ciRec.Fatigue,
								EquippedObjects: []schema.ActionObject{
									{
										Name: ceoRec.Name,
									},
								},
								StashedObjects: []schema.ActionObject{
									{
										Name: csoRec.Name,
									},
								},
							},
							Monster:         nil,
							EquippedObject:  nil,
							StashedObject:   nil,
							TargetObject:    nil,
							TargetCharacter: nil,
							TargetMonster:   nil,
							TargetLocation: &schema.ActionLocation{
								Name:        tlRec.Name,
								Description: tlRec.Description,
								Direction:   "north",
								Directions:  []string{"north", "south", "northwest"},
								Monsters: []schema.ActionLocationMonster{
									{
										Name: tlmRec.Name,
									},
								},
								Objects: []schema.ActionLocationObject{
									{
										Name: tloRec.Name,
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
					res := schema.ActionRequest{
						Data: schema.ActionRequestData{
							Sentence: "look rusted sword",
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
				loRec, _ := data.GetObjectRecByName("rusted sword")
				lmRec, _ := data.GetMonsterRecByName("grumpy dwarf")
				toRec, _ := data.GetObjectRecByName("rusted sword")

				res := schema.ActionResponse{
					Data: []schema.ActionResponseData{
						{
							Command:   "look",
							Narrative: fmt.Sprintf("%s looks %s", cRec.Name, loRec.Name),
							Location: schema.ActionLocation{
								Name:        lRec.Name,
								Description: lRec.Description,
								Directions:  []string{"north"},
								Characters: []schema.ActionLocationCharacter{
									{
										Name: cRec.Name,
									},
								},
								Monsters: []schema.ActionLocationMonster{
									{
										Name: lmRec.Name,
									},
								},
								Objects: []schema.ActionLocationObject{
									{
										Name: loRec.Name,
									},
								},
							},
							Character: &schema.ActionCharacter{
								Name:                cRec.Name,
								Strength:            cRec.Strength,
								Dexterity:           cRec.Dexterity,
								Intelligence:        cRec.Intelligence,
								Health:              cRec.Health,
								Fatigue:             cRec.Fatigue,
								CurrentStrength:     ciRec.Strength,
								CurrentDexterity:    ciRec.Dexterity,
								CurrentIntelligence: ciRec.Intelligence,
								CurrentHealth:       ciRec.Health,
								CurrentFatigue:      ciRec.Fatigue,
								EquippedObjects: []schema.ActionObject{
									{
										Name: ceoRec.Name,
									},
								},
								StashedObjects: []schema.ActionObject{
									{
										Name: csoRec.Name,
									},
								},
							},
							Monster:        nil,
							EquippedObject: nil,
							StashedObject:  nil,
							TargetObject: &schema.ActionObject{
								Name:        toRec.Name,
								Description: toRec.Description,
							},
							TargetCharacter: nil,
							TargetMonster:   nil,
							TargetLocation:  nil,
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
					res := schema.ActionRequest{
						Data: schema.ActionRequestData{
							Sentence: "look Grumpy Dwarf",
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
				loRec, _ := data.GetObjectRecByName("rusted sword")
				lmRec, _ := data.GetMonsterRecByName("grumpy dwarf")
				tmRec, _ := data.GetMonsterRecByName("grumpy dwarf")
				tmiRec, _ := data.GetMonsterInstanceRecByName("grumpy dwarf")
				tmeoRec, _ := data.GetObjectRecByName("bone dagger")

				res := schema.ActionResponse{
					Data: []schema.ActionResponseData{
						{
							Command:   "look",
							Narrative: fmt.Sprintf("%s looks %s", cRec.Name, tmRec.Name),
							Location: schema.ActionLocation{
								Name:        lRec.Name,
								Description: lRec.Description,
								Directions:  []string{"north"},
								Characters: []schema.ActionLocationCharacter{
									{
										Name: cRec.Name,
									},
								},
								Monsters: []schema.ActionLocationMonster{
									{
										Name: lmRec.Name,
									},
								},
								Objects: []schema.ActionLocationObject{
									{
										Name: loRec.Name,
									},
								},
							},
							Character: &schema.ActionCharacter{
								Name:                cRec.Name,
								Strength:            cRec.Strength,
								Dexterity:           cRec.Dexterity,
								Intelligence:        cRec.Intelligence,
								Health:              cRec.Health,
								Fatigue:             cRec.Fatigue,
								CurrentStrength:     ciRec.Strength,
								CurrentDexterity:    ciRec.Dexterity,
								CurrentIntelligence: ciRec.Intelligence,
								CurrentHealth:       ciRec.Health,
								CurrentFatigue:      ciRec.Fatigue,
								EquippedObjects: []schema.ActionObject{
									{
										Name: ceoRec.Name,
									},
								},
								StashedObjects: []schema.ActionObject{
									{
										Name: csoRec.Name,
									},
								},
							},
							Monster:         nil,
							EquippedObject:  nil,
							StashedObject:   nil,
							TargetObject:    nil,
							TargetCharacter: nil,
							TargetMonster: &schema.ActionMonster{
								Name:                tmRec.Name,
								Strength:            tmRec.Strength,
								Dexterity:           tmRec.Dexterity,
								Intelligence:        tmRec.Intelligence,
								Health:              tmRec.Health,
								Fatigue:             tmRec.Fatigue,
								CurrentStrength:     tmiRec.Strength,
								CurrentDexterity:    tmiRec.Dexterity,
								CurrentIntelligence: tmiRec.Intelligence,
								CurrentHealth:       tmiRec.Health,
								CurrentFatigue:      tmiRec.Fatigue,
								EquippedObjects: []schema.ActionObject{
									{
										Name: tmeoRec.Name,
									},
								},
							},
							TargetLocation: nil,
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
					res := schema.ActionRequest{
						Data: schema.ActionRequestData{
							Sentence: "look barricade",
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
				loRec, _ := data.GetObjectRecByName("rusted sword")
				lmRec, _ := data.GetMonsterRecByName("grumpy dwarf")
				tcRec, _ := data.GetCharacterRecByName("barricade")
				tciRec, _ := data.GetCharacterInstanceRecByName("barricade")
				tceoRec, _ := data.GetObjectRecByName("dull bronze ring")

				res := schema.ActionResponse{
					Data: []schema.ActionResponseData{
						{
							Command:   "look",
							Narrative: fmt.Sprintf("%s looks %s", cRec.Name, tcRec.Name),
							Location: schema.ActionLocation{
								Name:        lRec.Name,
								Description: lRec.Description,
								Directions:  []string{"north"},
								Characters: []schema.ActionLocationCharacter{
									{
										Name: cRec.Name,
									},
								},
								Monsters: []schema.ActionLocationMonster{
									{
										Name: lmRec.Name,
									},
								},
								Objects: []schema.ActionLocationObject{
									{
										Name: loRec.Name,
									},
								},
							},
							Character: &schema.ActionCharacter{
								Name:                cRec.Name,
								Strength:            cRec.Strength,
								Dexterity:           cRec.Dexterity,
								Intelligence:        cRec.Intelligence,
								Health:              cRec.Health,
								Fatigue:             cRec.Fatigue,
								CurrentStrength:     ciRec.Strength,
								CurrentDexterity:    ciRec.Dexterity,
								CurrentIntelligence: ciRec.Intelligence,
								CurrentHealth:       ciRec.Health,
								CurrentFatigue:      ciRec.Fatigue,
								EquippedObjects: []schema.ActionObject{
									{
										Name: ceoRec.Name,
									},
								},
								StashedObjects: []schema.ActionObject{
									{
										Name: csoRec.Name,
									},
								},
							},
							Monster:        nil,
							EquippedObject: nil,
							StashedObject:  nil,
							TargetObject:   nil,
							TargetCharacter: &schema.ActionCharacter{
								Name:                tcRec.Name,
								Strength:            tcRec.Strength,
								Dexterity:           tcRec.Dexterity,
								Intelligence:        tcRec.Intelligence,
								Health:              tcRec.Health,
								Fatigue:             tcRec.Fatigue,
								CurrentStrength:     tciRec.Strength,
								CurrentDexterity:    tciRec.Dexterity,
								CurrentIntelligence: tciRec.Intelligence,
								CurrentHealth:       tciRec.Health,
								CurrentFatigue:      tciRec.Fatigue,
								EquippedObjects: []schema.ActionObject{
									{
										Name: tceoRec.Name,
									},
								},
							},
							TargetMonster:  nil,
							TargetLocation: nil,
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
					res := schema.ActionRequest{
						Data: schema.ActionRequestData{
							Sentence: "stash rusted sword",
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
				lmRec, _ := data.GetMonsterRecByName("grumpy dwarf")
				toRec, _ := data.GetObjectRecByName("rusted sword")

				res := schema.ActionResponse{
					Data: []schema.ActionResponseData{
						{
							Command:   "stash",
							Narrative: fmt.Sprintf("%s stashes %s", cRec.Name, toRec.Name),
							Location: schema.ActionLocation{
								Name:        lRec.Name,
								Description: lRec.Description,
								Directions:  []string{"north"},
								Characters: []schema.ActionLocationCharacter{
									{
										Name: cRec.Name,
									},
								},
								Monsters: []schema.ActionLocationMonster{
									{
										Name: lmRec.Name,
									},
								},
								Objects: []schema.ActionLocationObject{},
							},
							Character: &schema.ActionCharacter{
								Name:                cRec.Name,
								Strength:            cRec.Strength,
								Dexterity:           cRec.Dexterity,
								Intelligence:        cRec.Intelligence,
								Health:              cRec.Health,
								Fatigue:             cRec.Fatigue,
								CurrentStrength:     ciRec.Strength,
								CurrentDexterity:    ciRec.Dexterity,
								CurrentIntelligence: ciRec.Intelligence,
								CurrentHealth:       ciRec.Health,
								CurrentFatigue:      ciRec.Fatigue,
								EquippedObjects: []schema.ActionObject{
									{
										Name: ceoRec.Name,
									},
								},
								StashedObjects: []schema.ActionObject{
									{
										Name: toRec.Name,
									},
									{
										Name: csoRec.Name,
									},
								},
							},
							Monster:        nil,
							EquippedObject: nil,
							StashedObject: &schema.ActionObject{
								Name:        toRec.Name,
								Description: toRec.Description,
								IsEquipped:  false,
								IsStashed:   true,
							},
							DroppedObject: nil,
							TargetObject: &schema.ActionObject{
								Name:        toRec.Name,
								Description: toRec.Description,
								IsEquipped:  false,
								IsStashed:   true,
							},
							TargetCharacter: nil,
							TargetMonster:   nil,
							TargetLocation:  nil,
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
					res := schema.ActionRequest{
						Data: schema.ActionRequestData{
							Sentence: "equip rusted sword",
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
				lmRec, _ := data.GetMonsterRecByName("grumpy dwarf")
				toRec, _ := data.GetObjectRecByName("rusted sword")

				res := schema.ActionResponse{
					Data: []schema.ActionResponseData{
						{
							Command:   "equip",
							Narrative: fmt.Sprintf("%s equips %s", cRec.Name, toRec.Name),
							Location: schema.ActionLocation{
								Name:        lRec.Name,
								Description: lRec.Description,
								Directions:  []string{"north"},
								Characters: []schema.ActionLocationCharacter{
									{
										Name: cRec.Name,
									},
								},
								Monsters: []schema.ActionLocationMonster{
									{
										Name: lmRec.Name,
									},
								},
								Objects: []schema.ActionLocationObject{},
							},
							Character: &schema.ActionCharacter{
								Name:                cRec.Name,
								Strength:            cRec.Strength,
								Dexterity:           cRec.Dexterity,
								Intelligence:        cRec.Intelligence,
								Health:              cRec.Health,
								Fatigue:             cRec.Fatigue,
								CurrentStrength:     ciRec.Strength,
								CurrentDexterity:    ciRec.Dexterity,
								CurrentIntelligence: ciRec.Intelligence,
								CurrentHealth:       ciRec.Health,
								CurrentFatigue:      ciRec.Fatigue,
								EquippedObjects: []schema.ActionObject{
									{
										Name: toRec.Name,
									},
									{
										Name: ceoRec.Name,
									},
								},
								StashedObjects: []schema.ActionObject{
									{
										Name: csoRec.Name,
									},
								},
							},
							Monster: nil,
							EquippedObject: &schema.ActionObject{
								Name:        toRec.Name,
								Description: toRec.Description,
								IsEquipped:  true,
								IsStashed:   false,
							},
							StashedObject: nil,
							DroppedObject: nil,
							TargetObject: &schema.ActionObject{
								Name:        toRec.Name,
								Description: toRec.Description,
								IsEquipped:  true,
								IsStashed:   false,
							},
							TargetCharacter: nil,
							TargetMonster:   nil,
							TargetLocation:  nil,
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
					res := schema.ActionRequest{
						Data: schema.ActionRequestData{
							Sentence: "drop Dull Bronze Ring",
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
				loRec, _ := data.GetObjectRecByName("rusted sword")
				lmRec, _ := data.GetMonsterRecByName("grumpy dwarf")
				toRec, _ := data.GetObjectRecByName("dull bronze ring")

				res := schema.ActionResponse{
					Data: []schema.ActionResponseData{
						{
							Command:   "drop",
							Narrative: fmt.Sprintf("%s drops %s", cRec.Name, toRec.Name),
							Location: schema.ActionLocation{
								Name:        lRec.Name,
								Description: lRec.Description,
								Directions:  []string{"north"},
								Characters: []schema.ActionLocationCharacter{
									{
										Name: cRec.Name,
									},
								},
								Monsters: []schema.ActionLocationMonster{
									{
										Name: lmRec.Name,
									},
								},
								Objects: []schema.ActionLocationObject{
									{
										Name: loRec.Name,
									},
									{
										Name: toRec.Name,
									},
								},
							},
							Character: &schema.ActionCharacter{
								Name:                cRec.Name,
								Strength:            cRec.Strength,
								Dexterity:           cRec.Dexterity,
								Intelligence:        cRec.Intelligence,
								Health:              cRec.Health,
								Fatigue:             cRec.Fatigue,
								CurrentStrength:     ciRec.Strength,
								CurrentDexterity:    ciRec.Dexterity,
								CurrentIntelligence: ciRec.Intelligence,
								CurrentHealth:       ciRec.Health,
								CurrentFatigue:      ciRec.Fatigue,
								EquippedObjects:     []schema.ActionObject{},
								StashedObjects: []schema.ActionObject{
									{
										Name: csoRec.Name,
									},
								},
							},
							Monster:        nil,
							EquippedObject: nil,
							StashedObject:  nil,
							DroppedObject: &schema.ActionObject{
								Name:        toRec.Name,
								Description: toRec.Description,
								IsEquipped:  false,
								IsStashed:   false,
							},
							TargetObject: &schema.ActionObject{
								Name:        toRec.Name,
								Description: toRec.Description,
								IsEquipped:  false,
								IsStashed:   false,
							},
							TargetCharacter: nil,
							TargetMonster:   nil,
							TargetLocation:  nil,
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
						require.Equal(t, expectData.Command, responseBody.Data[idx].Command)

						// Narrative
						require.Equal(t, expectData.Narrative, responseBody.Data[idx].Narrative)

						// Current location
						t.Logf("Checking location name >%s< >%s<", expectData.Location.Name, responseBody.Data[idx].Location.Name)
						require.Equal(t, expectData.Location.Name, responseBody.Data[idx].Location.Name)
						t.Logf("Checking location description >%s< >%s<", expectData.Location.Description, responseBody.Data[idx].Location.Description)
						require.Equal(t, expectData.Location.Description, responseBody.Data[idx].Location.Description)
						t.Logf("Checking location directions >%s< >%s<", expectData.Location.Directions, responseBody.Data[idx].Location.Directions)
						require.Equal(t, expectData.Location.Directions, responseBody.Data[idx].Location.Directions)

						// Current location characters
						t.Logf("Checking character count >%d< >%d<", len(expectData.Location.Characters), len(responseBody.Data[idx].Location.Characters))
						require.Equal(t, len(expectData.Location.Characters), len(responseBody.Data[idx].Location.Characters), "Response characters count equals expected")
						if len(expectData.Location.Characters) > 0 {
							for cIdx, character := range expectData.Location.Characters {
								t.Logf("Checking character name >%s< >%s<", character.Name, responseBody.Data[idx].Location.Characters[cIdx].Name)
								require.Equal(t, character.Name, responseBody.Data[idx].Location.Characters[cIdx].Name, "Character name equals expected")
							}
						}
						if len(expectData.Location.Characters) == 0 {
							require.Equal(t, 0, len(responseBody.Data[idx].Location.Characters), "Location characters length is 0")
						}

						// Current location monsters
						t.Logf("Checking monster count >%d< >%d<", len(expectData.Location.Monsters), len(responseBody.Data[idx].Location.Monsters))
						require.Equal(t, len(expectData.Location.Monsters), len(responseBody.Data[idx].Location.Monsters), "Response monsters count equals expected")
						if len(expectData.Location.Monsters) > 0 {
							for mIdx, monster := range expectData.Location.Monsters {
								t.Logf("Checking monster name >%s< >%s<", monster.Name, responseBody.Data[idx].Location.Monsters[mIdx].Name)
								require.Equal(t, monster.Name, responseBody.Data[idx].Location.Monsters[mIdx].Name, "Monster name equals expected")
							}
						}
						if len(expectData.Location.Monsters) == 0 {
							require.Equal(t, 0, len(responseBody.Data[idx].Location.Monsters), "Location monsters length is 0")
						}

						// Current location objects (any order)
						t.Logf("Checking object count >%d< >%d<", len(expectData.Location.Objects), len(responseBody.Data[idx].Location.Objects))
						require.Equal(t, len(expectData.Location.Objects), len(responseBody.Data[idx].Location.Objects), "Response objects count equals expected")
						if len(expectData.Location.Objects) > 0 {
							for _, eObject := range expectData.Location.Objects {
								found := false
								for _, rObject := range responseBody.Data[idx].Location.Objects {
									t.Logf("Checking expected location object name >%s< response location object name >%s<", eObject.Name, rObject.Name)
									if eObject.Name == rObject.Name {
										found = true
									}
								}
								require.True(t, found, fmt.Sprintf("Location object name >%s< found", eObject.Name))
							}
						}

						// Target location
						if expectData.TargetLocation != nil {
							require.NotNil(t, responseBody.Data[idx].TargetLocation, "Response target location is not empty")
							t.Logf("Checking location name >%s< >%s<", expectData.TargetLocation.Name, responseBody.Data[idx].TargetLocation.Name)
							require.Equal(t, expectData.TargetLocation.Name, responseBody.Data[idx].TargetLocation.Name)
							t.Logf("Checking location description >%s< >%s<", expectData.TargetLocation.Description, responseBody.Data[idx].TargetLocation.Description)
							require.Equal(t, expectData.TargetLocation.Description, responseBody.Data[idx].TargetLocation.Description)
							t.Logf("Checking location direction >%s< >%s<", expectData.TargetLocation.Direction, responseBody.Data[idx].TargetLocation.Direction)
							require.Equal(t, expectData.TargetLocation.Direction, responseBody.Data[idx].TargetLocation.Direction)
							t.Logf("Checking location directions >%s< >%s<", expectData.TargetLocation.Directions, responseBody.Data[idx].TargetLocation.Directions)
							require.Equal(t, expectData.TargetLocation.Directions, responseBody.Data[idx].TargetLocation.Directions)

							// Target location characters
							t.Logf("Checking character count >%d< >%d<", len(expectData.TargetLocation.Characters), len(responseBody.Data[idx].TargetLocation.Characters))
							require.Equal(t, len(expectData.TargetLocation.Characters), len(responseBody.Data[idx].TargetLocation.Characters), "Response characters count equals expected")
							if len(expectData.TargetLocation.Characters) > 0 {
								for cIdx, character := range expectData.TargetLocation.Characters {
									t.Logf("Checking character name >%s< >%s<", character.Name, responseBody.Data[idx].TargetLocation.Characters[cIdx].Name)
									require.Equal(t, character.Name, responseBody.Data[idx].TargetLocation.Characters[cIdx].Name, "Character name equals expected")
								}
							}

							// Target location monsters
							t.Logf("Checking monster count >%d< >%d<", len(expectData.TargetLocation.Monsters), len(responseBody.Data[idx].TargetLocation.Monsters))
							require.Equal(t, len(expectData.TargetLocation.Monsters), len(responseBody.Data[idx].TargetLocation.Monsters), "Response monsters count equals expected")
							if len(expectData.TargetLocation.Monsters) > 0 {
								for mIdx, monster := range expectData.TargetLocation.Monsters {
									t.Logf("Checking monster name >%s< >%s<", monster.Name, responseBody.Data[idx].TargetLocation.Monsters[mIdx].Name)
									require.Equal(t, monster.Name, responseBody.Data[idx].TargetLocation.Monsters[mIdx].Name, "Monster name equals expected")
								}
							}

							// Target location objects
							t.Logf("Checking object count >%d< >%d<", len(expectData.TargetLocation.Objects), len(responseBody.Data[idx].TargetLocation.Objects))
							require.Equal(t, len(expectData.TargetLocation.Objects), len(responseBody.Data[idx].TargetLocation.Objects), "Response objects count equals expected")
							if len(expectData.TargetLocation.Objects) > 0 {
								for _, eObject := range expectData.TargetLocation.Objects {
									found := false
									for _, tObject := range responseBody.Data[idx].TargetLocation.Objects {
										t.Logf("Checking expected target location object name >%s< response target location object name >%s<", eObject.Name, tObject.Name)
										if eObject.Name == tObject.Name {
											found = true
										}
									}
									require.True(t, found, fmt.Sprintf("Target location object name >%s< found", eObject.Name))
								}
							}
						}

						// Response character
						t.Logf("Checking character nil >%t< >%t<", isCharacterNil(expectData.Character), isCharacterNil(responseBody.Data[idx].Character))
						require.Equal(t, isCharacterNil(expectData.Character), isCharacterNil(responseBody.Data[idx].Character), "Response character is nil or not nil as expected")
						if !isCharacterNil(expectData.Character) {
							t.Logf("Checking action character name >%s< >%s<", expectData.Character.Name, responseBody.Data[idx].Character.Name)
							require.Equal(t, expectData.Character.Name, responseBody.Data[idx].Character.Name, "Response character name equals expected")
							t.Logf("Checking action character strenth >%d< >%d<", expectData.Character.Strength, responseBody.Data[idx].Character.Strength)
							require.Equal(t, expectData.Character.Strength, responseBody.Data[idx].Character.Strength, "Response character strength equals expected")
							t.Logf("Checking action character dexterity >%d< >%d<", expectData.Character.Dexterity, responseBody.Data[idx].Character.Dexterity)
							require.Equal(t, expectData.Character.Dexterity, responseBody.Data[idx].Character.Dexterity, "Response character dexterity equals expected")
							t.Logf("Checking action character intelligence >%d< >%d<", expectData.Character.Intelligence, responseBody.Data[idx].Character.Intelligence)
							require.Equal(t, expectData.Character.Intelligence, responseBody.Data[idx].Character.Intelligence, "Response character intelligence equals expected")

							t.Logf("Checking action character current strenth >%d< >%d<", expectData.Character.CurrentStrength, responseBody.Data[idx].Character.CurrentStrength)
							require.Equal(t, expectData.Character.CurrentStrength, responseBody.Data[idx].Character.CurrentStrength, "Response character current strength equals expected")
							t.Logf("Checking action character current dexterity >%d< >%d<", expectData.Character.CurrentDexterity, responseBody.Data[idx].Character.CurrentDexterity)
							require.Equal(t, expectData.Character.CurrentDexterity, responseBody.Data[idx].Character.CurrentDexterity, "Response character current dexterity equals expected")
							t.Logf("Checking action character current intelligence >%d< >%d<", expectData.Character.CurrentIntelligence, responseBody.Data[idx].Character.CurrentIntelligence)
							require.Equal(t, expectData.Character.CurrentIntelligence, responseBody.Data[idx].Character.CurrentIntelligence, "Response character current intelligence equals expected")

							t.Logf("Checking action character health >%d< >%d<", expectData.Character.Health, responseBody.Data[idx].Character.Health)
							require.Equal(t, expectData.Character.Health, responseBody.Data[idx].Character.Health, "Response character health equals expected")
							t.Logf("Checking action character fatigue >%d< >%d<", expectData.Character.Fatigue, responseBody.Data[idx].Character.Fatigue)
							require.Equal(t, expectData.Character.Fatigue, responseBody.Data[idx].Character.Fatigue, "Response character fatigue equals expected")

							t.Logf("Checking action character current health >%d< >%d<", expectData.Character.CurrentHealth, responseBody.Data[idx].Character.CurrentHealth)
							require.Equal(t, expectData.Character.CurrentHealth, responseBody.Data[idx].Character.CurrentHealth, "Response character current health equals expected")
							t.Logf("Checking action character current fatigue >%d< >%d<", expectData.Character.CurrentFatigue, responseBody.Data[idx].Character.CurrentFatigue)
							require.Equal(t, expectData.Character.CurrentFatigue, responseBody.Data[idx].Character.CurrentFatigue, "Response character current fatigue equals expected")

							t.Logf("Checking action character equipped objects >%d< >%d<", len(expectData.Character.EquippedObjects), len(responseBody.Data[idx].Character.EquippedObjects))
							require.Equal(t, len(expectData.Character.EquippedObjects), len(responseBody.Data[idx].Character.EquippedObjects), "Response character equipped object count equals expected")

							t.Logf("Checking action character stashed objects >%d< >%d<", len(expectData.Character.StashedObjects), len(responseBody.Data[idx].Character.StashedObjects))
							require.Equal(t, len(expectData.Character.StashedObjects), len(responseBody.Data[idx].Character.StashedObjects), "Response character stashed object count equals expected")
						}

						// Response target character
						t.Logf("Checking target character nil >%t< >%t<", isCharacterNil(expectData.TargetCharacter), isCharacterNil(responseBody.Data[idx].TargetCharacter))
						require.Equal(t, isCharacterNil(expectData.TargetCharacter), isCharacterNil(responseBody.Data[idx].TargetCharacter), "Response target character is nil or not nil as expected")
						if !isCharacterNil(expectData.TargetCharacter) {
							t.Logf("Checking action target character name >%s< >%s<", expectData.TargetCharacter.Name, responseBody.Data[idx].TargetCharacter.Name)
							require.Equal(t, expectData.TargetCharacter.Name, responseBody.Data[idx].TargetCharacter.Name, "Response target character name equals expected")

							t.Logf("Checking action target character strenth >%d< >%d<", expectData.TargetCharacter.Strength, responseBody.Data[idx].TargetCharacter.Strength)
							require.Equal(t, expectData.TargetCharacter.Strength, responseBody.Data[idx].TargetCharacter.Strength, "Response target character strength equals expected")
							t.Logf("Checking action target character dexterity >%d< >%d<", expectData.TargetCharacter.Dexterity, responseBody.Data[idx].TargetCharacter.Dexterity)
							require.Equal(t, expectData.TargetCharacter.Dexterity, responseBody.Data[idx].TargetCharacter.Dexterity, "Response target character dexterity equals expected")
							t.Logf("Checking action target character intelligence >%d< >%d<", expectData.TargetCharacter.Intelligence, responseBody.Data[idx].TargetCharacter.Intelligence)
							require.Equal(t, expectData.TargetCharacter.Intelligence, responseBody.Data[idx].TargetCharacter.Intelligence, "Response target character intelligence equals expected")

							t.Logf("Checking target character current strenth >%d< >%d<", expectData.TargetCharacter.CurrentStrength, responseBody.Data[idx].TargetCharacter.CurrentStrength)
							require.Equal(t, expectData.TargetCharacter.CurrentStrength, responseBody.Data[idx].TargetCharacter.CurrentStrength, "Response target character current strength equals expected")
							t.Logf("Checking target character current dexterity >%d< >%d<", expectData.TargetCharacter.CurrentDexterity, responseBody.Data[idx].TargetCharacter.CurrentDexterity)
							require.Equal(t, expectData.TargetCharacter.CurrentDexterity, responseBody.Data[idx].TargetCharacter.CurrentDexterity, "Response target character current dexterity equals expected")
							t.Logf("Checking target character current intelligence >%d< >%d<", expectData.TargetCharacter.CurrentIntelligence, responseBody.Data[idx].TargetCharacter.CurrentIntelligence)
							require.Equal(t, expectData.TargetCharacter.CurrentIntelligence, responseBody.Data[idx].TargetCharacter.CurrentIntelligence, "Response target character current intelligence equals expected")

							t.Logf("Checking target character health >%d< >%d<", expectData.TargetCharacter.Health, responseBody.Data[idx].TargetCharacter.Health)
							require.Equal(t, expectData.TargetCharacter.Health, responseBody.Data[idx].TargetCharacter.Health, "Response target character health equals expected")
							t.Logf("Checking target character fatigue >%d< >%d<", expectData.TargetCharacter.Fatigue, responseBody.Data[idx].TargetCharacter.Fatigue)
							require.Equal(t, expectData.TargetCharacter.Fatigue, responseBody.Data[idx].TargetCharacter.Fatigue, "Response target character fatigue equals expected")

							t.Logf("Checking target character current health >%d< >%d<", expectData.TargetCharacter.CurrentHealth, responseBody.Data[idx].TargetCharacter.CurrentHealth)
							require.Equal(t, expectData.TargetCharacter.CurrentHealth, responseBody.Data[idx].TargetCharacter.CurrentHealth, "Response target character current health equals expected")
							t.Logf("Checking target character current fatigue >%d< >%d<", expectData.TargetCharacter.CurrentFatigue, responseBody.Data[idx].TargetCharacter.CurrentFatigue)
							require.Equal(t, expectData.TargetCharacter.CurrentFatigue, responseBody.Data[idx].TargetCharacter.CurrentFatigue, "Response target character current fatigue equals expected")

							t.Logf("Checking target character equipped objects >%d< >%d<", len(expectData.TargetCharacter.EquippedObjects), len(responseBody.Data[idx].TargetCharacter.EquippedObjects))
							require.Equal(t, len(expectData.TargetCharacter.EquippedObjects), len(responseBody.Data[idx].TargetCharacter.EquippedObjects), "Response target character equipped object count equals expected")

							t.Logf("Checking target character stashed objects >%d< >%d<", len(expectData.TargetCharacter.StashedObjects), len(responseBody.Data[idx].TargetCharacter.StashedObjects))
							require.Equal(t, len(expectData.TargetCharacter.StashedObjects), len(responseBody.Data[idx].TargetCharacter.StashedObjects), "Response target character stashed object count equals expected")
						}

						// Response monster
						t.Logf("Checking monster nil >%t< >%t<", isMonsterNil(expectData.Monster), isMonsterNil(responseBody.Data[idx].Monster))
						require.Equal(t, isMonsterNil(expectData.Monster), isMonsterNil(responseBody.Data[idx].Monster), "Response monster is nil or not nil as expected")
						if !isMonsterNil(expectData.Monster) {
							t.Logf("Checking monster name >%s< >%s<", expectData.Monster.Name, responseBody.Data[idx].Monster.Name)
							require.Equal(t, expectData.Monster.Name, responseBody.Data[idx].Monster.Name, "Response monster name equals expected")

							t.Logf("Checking monster strenth >%d< >%d<", expectData.Monster.Strength, responseBody.Data[idx].Monster.Strength)
							require.Equal(t, expectData.Monster.Strength, responseBody.Data[idx].Monster.Strength, "Response monster strength equals expected")
							t.Logf("Checking monster dexterity >%d< >%d<", expectData.Monster.Dexterity, responseBody.Data[idx].Monster.Dexterity)
							require.Equal(t, expectData.Monster.Dexterity, responseBody.Data[idx].Monster.Dexterity, "Response monster dexterity equals expected")
							t.Logf("Checking monster intelligence >%d< >%d<", expectData.Monster.Intelligence, responseBody.Data[idx].Monster.Intelligence)
							require.Equal(t, expectData.Monster.Intelligence, responseBody.Data[idx].Monster.Intelligence, "Response monster intelligence equals expected")

							t.Logf("Checking monster current strenth >%d< >%d<", expectData.Monster.CurrentStrength, responseBody.Data[idx].Monster.CurrentStrength)
							require.Equal(t, expectData.Monster.CurrentStrength, responseBody.Data[idx].Monster.CurrentStrength, "Response monster current strength equals expected")
							t.Logf("Checking monster current dexterity >%d< >%d<", expectData.Monster.CurrentDexterity, responseBody.Data[idx].Monster.CurrentDexterity)
							require.Equal(t, expectData.Monster.CurrentDexterity, responseBody.Data[idx].Monster.CurrentDexterity, "Response monster current dexterity equals expected")
							t.Logf("Checking monster current intelligence >%d< >%d<", expectData.Monster.CurrentIntelligence, responseBody.Data[idx].Monster.CurrentIntelligence)
							require.Equal(t, expectData.Monster.CurrentIntelligence, responseBody.Data[idx].Monster.CurrentIntelligence, "Response monster current intelligence equals expected")

							t.Logf("Checking monster health >%d< >%d<", expectData.Monster.Health, responseBody.Data[idx].Monster.Health)
							require.Equal(t, expectData.Monster.Health, responseBody.Data[idx].Monster.Health, "Response monster health equals expected")
							t.Logf("Checking monster fatigue >%d< >%d<", expectData.Monster.Fatigue, responseBody.Data[idx].Monster.Fatigue)
							require.Equal(t, expectData.Monster.Fatigue, responseBody.Data[idx].Monster.Fatigue, "Response monster fatigue equals expected")

							t.Logf("Checking monster current health >%d< >%d<", expectData.Monster.CurrentHealth, responseBody.Data[idx].Monster.CurrentHealth)
							require.Equal(t, expectData.Monster.CurrentHealth, responseBody.Data[idx].Monster.CurrentHealth, "Response monster current health equals expected")
							t.Logf("Checking monster current fatigue >%d< >%d<", expectData.Monster.CurrentFatigue, responseBody.Data[idx].Monster.CurrentFatigue)
							require.Equal(t, expectData.Monster.CurrentFatigue, responseBody.Data[idx].Monster.CurrentFatigue, "Response monster current fatigue equals expected")
						}

						// Response target monster
						t.Logf("Checking target monster nil >%t< >%t<", isMonsterNil(expectData.TargetMonster), isMonsterNil(responseBody.Data[idx].TargetMonster))
						require.Equal(t, isMonsterNil(expectData.TargetMonster), isMonsterNil(responseBody.Data[idx].TargetMonster), "Response target monster is nil or not nil as expected")
						if !isMonsterNil(expectData.TargetMonster) {
							t.Logf("Checking target monster name >%s< >%s<", expectData.TargetMonster.Name, responseBody.Data[idx].TargetMonster.Name)
							require.Equal(t, expectData.TargetMonster.Name, responseBody.Data[idx].TargetMonster.Name, "Response target monster name equals expected")

							t.Logf("Checking target monster strenth >%d< >%d<", expectData.TargetMonster.Strength, responseBody.Data[idx].TargetMonster.Strength)
							require.Equal(t, expectData.TargetMonster.Strength, responseBody.Data[idx].TargetMonster.Strength, "Response target monster strength equals expected")
							t.Logf("Checking target monster dexterity >%d< >%d<", expectData.TargetMonster.Dexterity, responseBody.Data[idx].TargetMonster.Dexterity)
							require.Equal(t, expectData.TargetMonster.Dexterity, responseBody.Data[idx].TargetMonster.Dexterity, "Response target monster dexterity equals expected")
							t.Logf("Checking target monster intelligence >%d< >%d<", expectData.TargetMonster.Intelligence, responseBody.Data[idx].TargetMonster.Intelligence)
							require.Equal(t, expectData.TargetMonster.Intelligence, responseBody.Data[idx].TargetMonster.Intelligence, "Response target monster intelligence equals expected")

							t.Logf("Checking target monster current strenth >%d< >%d<", expectData.TargetMonster.CurrentStrength, responseBody.Data[idx].TargetMonster.CurrentStrength)
							require.Equal(t, expectData.TargetMonster.CurrentStrength, responseBody.Data[idx].TargetMonster.CurrentStrength, "Response monster current strength equals expected")
							t.Logf("Checking target monster current dexterity >%d< >%d<", expectData.TargetMonster.CurrentDexterity, responseBody.Data[idx].TargetMonster.CurrentDexterity)
							require.Equal(t, expectData.TargetMonster.CurrentDexterity, responseBody.Data[idx].TargetMonster.CurrentDexterity, "Response monster current dexterity equals expected")
							t.Logf("Checking target monster current intelligence >%d< >%d<", expectData.TargetMonster.CurrentIntelligence, responseBody.Data[idx].TargetMonster.CurrentIntelligence)
							require.Equal(t, expectData.TargetMonster.CurrentIntelligence, responseBody.Data[idx].TargetMonster.CurrentIntelligence, "Response monster current intelligence equals expected")

							t.Logf("Checking target monster health >%d< >%d<", expectData.TargetMonster.Health, responseBody.Data[idx].TargetMonster.Health)
							require.Equal(t, expectData.TargetMonster.Health, responseBody.Data[idx].TargetMonster.Health, "Response target monster health equals expected")
							t.Logf("Checking target monster fatigue >%d< >%d<", expectData.TargetMonster.Fatigue, responseBody.Data[idx].TargetMonster.Fatigue)
							require.Equal(t, expectData.TargetMonster.Fatigue, responseBody.Data[idx].TargetMonster.Fatigue, "Response target monster fatigue equals expected")

							t.Logf("Checking target monster current health >%d< >%d<", expectData.TargetMonster.CurrentHealth, responseBody.Data[idx].TargetMonster.CurrentHealth)
							require.Equal(t, expectData.TargetMonster.CurrentHealth, responseBody.Data[idx].TargetMonster.CurrentHealth, "Response target monster current health equals expected")
							t.Logf("Checking target monster current fatigue >%d< >%d<", expectData.TargetMonster.CurrentFatigue, responseBody.Data[idx].TargetMonster.CurrentFatigue)
							require.Equal(t, expectData.TargetMonster.CurrentFatigue, responseBody.Data[idx].TargetMonster.CurrentFatigue, "Response target monster current fatigue equals expected")
						}

						// Response target object
						t.Logf("Checking target object nil >%t< >%t<", isObjectNil(expectData.TargetObject), isObjectNil(responseBody.Data[idx].TargetObject))
						require.Equal(t, isObjectNil(expectData.TargetObject), isObjectNil(responseBody.Data[idx].TargetObject), "Response target object is nil or not nil as expected")
						if !isObjectNil(expectData.TargetObject) {
							t.Logf("Checking target object name >%s< >%s<", expectData.TargetObject.Name, responseBody.Data[idx].TargetObject.Name)
							require.Equal(t, expectData.TargetObject.Name, responseBody.Data[idx].TargetObject.Name, "Response target object name equals expected")
							t.Logf("Checking target object description >%s< >%s<", expectData.TargetObject.Description, responseBody.Data[idx].TargetObject.Description)
							require.Equal(t, expectData.TargetObject.Description, responseBody.Data[idx].TargetObject.Description, "Response target object description equals expected")
						}

						// Response stashed object
						t.Logf("Checking stashed object nil >%t< >%t<", isObjectNil(expectData.StashedObject), isObjectNil(responseBody.Data[idx].StashedObject))
						require.Equal(t, isObjectNil(expectData.StashedObject), isObjectNil(responseBody.Data[idx].StashedObject), "Response stashed object is nil or not nil as expected")
						if !isObjectNil(expectData.StashedObject) {
							t.Logf("Checking stashed object name >%s< >%s<", expectData.StashedObject.Name, responseBody.Data[idx].StashedObject.Name)
							require.Equal(t, expectData.StashedObject.Name, responseBody.Data[idx].StashedObject.Name, "Response stashed object name equals expected")
							t.Logf("Checking stashed object description >%s< >%s<", expectData.StashedObject.Description, responseBody.Data[idx].StashedObject.Description)
							require.Equal(t, expectData.StashedObject.Description, responseBody.Data[idx].StashedObject.Description, "Response stashed object description equals expected")

							t.Logf("Checking stashed object is_equipped >%t< >%t<", expectData.StashedObject.IsEquipped, responseBody.Data[idx].StashedObject.IsEquipped)
							require.Equal(t, expectData.StashedObject.IsEquipped, responseBody.Data[idx].StashedObject.IsEquipped, "Response stashed object is_equipped equals expected")
							t.Logf("Checking stashed object is_stashed >%t< >%t<", expectData.StashedObject.IsStashed, responseBody.Data[idx].StashedObject.IsStashed)
							require.Equal(t, expectData.StashedObject.IsStashed, responseBody.Data[idx].StashedObject.IsStashed, "Response stashed object is_stashed equals expected")
						}

						// Response equipped object
						t.Logf("Checking equipped object nil >%t< >%t<", isObjectNil(expectData.EquippedObject), isObjectNil(responseBody.Data[idx].EquippedObject))
						require.Equal(t, isObjectNil(expectData.EquippedObject), isObjectNil(responseBody.Data[idx].EquippedObject), "Response equipped object is nil or not nil as expected")
						if !isObjectNil(expectData.EquippedObject) {
							t.Logf("Checking equipped object name >%s< >%s<", expectData.EquippedObject.Name, responseBody.Data[idx].EquippedObject.Name)
							require.Equal(t, expectData.EquippedObject.Name, responseBody.Data[idx].EquippedObject.Name, "Response equipped object name equals expected")
							t.Logf("Checking equipped object description >%s< >%s<", expectData.EquippedObject.Description, responseBody.Data[idx].EquippedObject.Description)
							require.Equal(t, expectData.EquippedObject.Description, responseBody.Data[idx].EquippedObject.Description, "Response equipped object description equals expected")

							t.Logf("Checking equipped object is_equipped >%t< >%t<", expectData.EquippedObject.IsEquipped, responseBody.Data[idx].EquippedObject.IsEquipped)
							require.Equal(t, expectData.EquippedObject.IsEquipped, responseBody.Data[idx].EquippedObject.IsEquipped, "Response equipped object is_equipped equals expected")
							t.Logf("Checking equipped object is_stashed >%t< >%t<", expectData.EquippedObject.IsStashed, responseBody.Data[idx].EquippedObject.IsStashed)
							require.Equal(t, expectData.EquippedObject.IsStashed, responseBody.Data[idx].EquippedObject.IsStashed, "Response equipped object is_stashed equals expected")
						}

						// Response dropped object
						t.Logf("Checking dropped object nil >%t< >%t<", isObjectNil(expectData.DroppedObject), isObjectNil(responseBody.Data[idx].DroppedObject))
						require.Equal(t, isObjectNil(expectData.DroppedObject), isObjectNil(responseBody.Data[idx].DroppedObject), "Response dropped object is nil or not nil as expected")
						if !isObjectNil(expectData.DroppedObject) {
							t.Logf("Checking dropped object name >%s< >%s<", expectData.DroppedObject.Name, responseBody.Data[idx].DroppedObject.Name)
							require.Equal(t, expectData.DroppedObject.Name, responseBody.Data[idx].DroppedObject.Name, "Response dropped object name equals expected")
							t.Logf("Checking dropped object description >%s< >%s<", expectData.DroppedObject.Description, responseBody.Data[idx].DroppedObject.Description)
							require.Equal(t, expectData.DroppedObject.Description, responseBody.Data[idx].DroppedObject.Description, "Response dropped object description equals expected")

							t.Logf("Checking dropped object is_dropped >%t< >%t<", expectData.DroppedObject.IsEquipped, responseBody.Data[idx].DroppedObject.IsEquipped)
							require.Equal(t, expectData.DroppedObject.IsEquipped, responseBody.Data[idx].DroppedObject.IsEquipped, "Response dropped object is_dropped equals expected")
							t.Logf("Checking dropped object is_stashed >%t< >%t<", expectData.DroppedObject.IsStashed, responseBody.Data[idx].DroppedObject.IsStashed)
							require.Equal(t, expectData.DroppedObject.IsStashed, responseBody.Data[idx].DroppedObject.IsStashed, "Response dropped object is_stashed equals expected")
						}
					}
				}

				// Check dates and add teardown ID's
				for _, data := range responseBody.Data {
					require.False(t, data.CreatedAt.IsZero(), "CreatedAt is not zero")
					if method == http.MethodPost {
						require.True(t, data.UpdatedAt.IsZero(), "UpdatedAt is zero")
						t.Logf("Adding dungeon character action ID >%s< for teardown", data.ID)
						th.AddActionTeardownID(data.ID)
					}
				}
			}

			RunTestCase(t, th, &testCase, testFunc)
		})
	}
}
