package runner

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"gitlab.com/alienspaces/go-mud/server/schema"

	"github.com/stretchr/testify/require"
	"gitlab.com/alienspaces/go-mud/server/core/auth"
	"gitlab.com/alienspaces/go-mud/server/core/server"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/harness"
)

func TestCreateDungeonCharacterActionHandler(t *testing.T) {

	// Test harness
	th, err := NewTestHarness()
	require.NoError(t, err, "New test data returns without error")

	type testCase struct {
		TestCase
		expectResponseBody func(data harness.Data) *schema.DungeonActionResponse
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
		return rnr.HandlerConfig[6]
	}

	testCaseRequestHeaders := func(data harness.Data) map[string]string {
		headers := map[string]string{
			"Authorization": "Bearer " + validAuthToken(),
		}
		return headers
	}

	testCaseRequestPathParams := func(data harness.Data) map[string]string {
		params := map[string]string{
			":dungeon_id":   data.DungeonRecs[0].ID,
			":character_id": data.DungeonCharacterRecs[0].ID,
		}
		return params
	}

	testCaseResponseBody := func(body io.Reader) (interface{}, error) {
		var responseBody *schema.DungeonActionResponse
		err = json.NewDecoder(body).Decode(&responseBody)
		return responseBody, err
	}

	testCases := []testCase{
		{
			// Look at the current room
			TestCase: TestCase{
				Skip:              false,
				Name:              "Look at the current room",
				HandlerConfig:     testCaseHandlerConfig,
				RequestHeaders:    testCaseRequestHeaders,
				RequestPathParams: testCaseRequestPathParams,
				RequestBody: func(data harness.Data) interface{} {
					res := schema.DungeonActionRequest{
						Data: schema.DungeonActionRequestData{
							Sentence: "look",
						},
					}
					return &res
				},
				ResponseBody: testCaseResponseBody,
				ResponseCode: http.StatusOK,
			},
			expectResponseBody: func(data harness.Data) *schema.DungeonActionResponse {
				res := schema.DungeonActionResponse{
					Data: []schema.DungeonActionResponseData{
						{
							Command: "look",
							Location: schema.LocationData{
								Name:        data.DungeonLocationRecs[0].Name,
								Description: data.DungeonLocationRecs[0].Description,
								Directions:  []string{"north"},
								Characters: []schema.CharacterData{
									{
										Name: data.DungeonCharacterRecs[0].Name,
									},
								},
								Monsters: []schema.MonsterData{
									{
										Name: data.DungeonMonsterRecs[0].Name,
									},
								},
								Objects: []schema.ObjectData{
									{
										Name: data.DungeonObjectRecs[0].Name,
									},
								},
							},
							Character: &schema.CharacterDetailedData{
								Name:                data.DungeonCharacterRecs[0].Name,
								Strength:            data.DungeonCharacterRecs[0].Strength,
								Dexterity:           data.DungeonCharacterRecs[0].Dexterity,
								Intelligence:        data.DungeonCharacterRecs[0].Intelligence,
								CurrentStrength:     data.DungeonCharacterRecs[0].CurrentStrength,
								CurrentDexterity:    data.DungeonCharacterRecs[0].CurrentDexterity,
								CurrentIntelligence: data.DungeonCharacterRecs[0].CurrentIntelligence,
								Health:              data.DungeonCharacterRecs[0].Health,
								Fatigue:             data.DungeonCharacterRecs[0].Fatigue,
								CurrentHealth:       data.DungeonCharacterRecs[0].CurrentHealth,
								CurrentFatigue:      data.DungeonCharacterRecs[0].CurrentFatigue,
							},
							Monster:         nil,
							EquippedObject:  nil,
							StashedObject:   nil,
							TargetObject:    nil,
							TargetCharacter: nil,
							TargetMonster:   nil,
							TargetLocation: &schema.LocationData{
								Name:        data.DungeonLocationRecs[0].Name,
								Description: data.DungeonLocationRecs[0].Description,
								Directions:  []string{"north"},
								Characters: []schema.CharacterData{
									{
										Name: data.DungeonCharacterRecs[0].Name,
									},
								},
								Monsters: []schema.MonsterData{
									{
										Name: data.DungeonMonsterRecs[0].Name,
									},
								},
								Objects: []schema.ObjectData{
									{
										Name: data.DungeonObjectRecs[0].Name,
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
				Skip:              false,
				Name:              "Move north from the current room",
				HandlerConfig:     testCaseHandlerConfig,
				RequestHeaders:    testCaseRequestHeaders,
				RequestPathParams: testCaseRequestPathParams,
				RequestBody: func(data harness.Data) interface{} {
					res := schema.DungeonActionRequest{
						Data: schema.DungeonActionRequestData{
							Sentence: "move north",
						},
					}
					return &res
				},
				ResponseBody: testCaseResponseBody,
				ResponseCode: http.StatusOK,
			},
			expectResponseBody: func(data harness.Data) *schema.DungeonActionResponse {
				res := schema.DungeonActionResponse{
					Data: []schema.DungeonActionResponseData{
						{
							Command: "move",
							Location: schema.LocationData{
								Name:        data.DungeonLocationRecs[1].Name,
								Description: data.DungeonLocationRecs[1].Description,
								Directions:  []string{"north", "south", "northwest"},
								Characters: []schema.CharacterData{
									{
										Name: data.DungeonCharacterRecs[0].Name,
									},
								},
								Monsters: []schema.MonsterData{
									{
										Name: data.DungeonMonsterRecs[1].Name,
									},
								},
							},
							Character: &schema.CharacterDetailedData{
								Name:                data.DungeonCharacterRecs[0].Name,
								Strength:            data.DungeonCharacterRecs[0].Strength,
								Dexterity:           data.DungeonCharacterRecs[0].Dexterity,
								Intelligence:        data.DungeonCharacterRecs[0].Intelligence,
								CurrentStrength:     data.DungeonCharacterRecs[0].CurrentStrength,
								CurrentDexterity:    data.DungeonCharacterRecs[0].CurrentDexterity,
								CurrentIntelligence: data.DungeonCharacterRecs[0].CurrentIntelligence,
								Health:              data.DungeonCharacterRecs[0].Health,
								Fatigue:             data.DungeonCharacterRecs[0].Fatigue,
								CurrentHealth:       data.DungeonCharacterRecs[0].CurrentHealth,
								CurrentFatigue:      data.DungeonCharacterRecs[0].CurrentFatigue,
							},
							Monster:         nil,
							EquippedObject:  nil,
							StashedObject:   nil,
							TargetObject:    nil,
							TargetCharacter: nil,
							TargetMonster:   nil,
							TargetLocation: &schema.LocationData{
								Name:        data.DungeonLocationRecs[1].Name,
								Description: data.DungeonLocationRecs[1].Description,
								Direction:   "north",
								Directions:  []string{"north", "south", "northwest"},
								Characters: []schema.CharacterData{
									{
										Name: data.DungeonCharacterRecs[0].Name,
									},
								},
								Monsters: []schema.MonsterData{
									{
										Name: data.DungeonMonsterRecs[1].Name,
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
				Skip:              false,
				Name:              "Look north from the current room",
				HandlerConfig:     testCaseHandlerConfig,
				RequestHeaders:    testCaseRequestHeaders,
				RequestPathParams: testCaseRequestPathParams,
				RequestBody: func(data harness.Data) interface{} {
					res := schema.DungeonActionRequest{
						Data: schema.DungeonActionRequestData{
							Sentence: "look north",
						},
					}
					return &res
				},
				ResponseBody: testCaseResponseBody,
				ResponseCode: http.StatusOK,
			},
			expectResponseBody: func(data harness.Data) *schema.DungeonActionResponse {
				res := schema.DungeonActionResponse{
					Data: []schema.DungeonActionResponseData{
						{
							Command: "look",
							Location: schema.LocationData{
								Name:        data.DungeonLocationRecs[0].Name,
								Description: data.DungeonLocationRecs[0].Description,
								Directions:  []string{"north"},
								Characters: []schema.CharacterData{
									{
										Name: data.DungeonCharacterRecs[0].Name,
									},
								},
								Monsters: []schema.MonsterData{
									{
										Name: data.DungeonMonsterRecs[0].Name,
									},
								},
								Objects: []schema.ObjectData{
									{
										Name: data.DungeonObjectRecs[0].Name,
									},
								},
							},
							Character: &schema.CharacterDetailedData{
								Name:                data.DungeonCharacterRecs[0].Name,
								Strength:            data.DungeonCharacterRecs[0].Strength,
								Dexterity:           data.DungeonCharacterRecs[0].Dexterity,
								Intelligence:        data.DungeonCharacterRecs[0].Intelligence,
								CurrentStrength:     data.DungeonCharacterRecs[0].CurrentStrength,
								CurrentDexterity:    data.DungeonCharacterRecs[0].CurrentDexterity,
								CurrentIntelligence: data.DungeonCharacterRecs[0].CurrentIntelligence,
								Health:              data.DungeonCharacterRecs[0].Health,
								Fatigue:             data.DungeonCharacterRecs[0].Fatigue,
								CurrentHealth:       data.DungeonCharacterRecs[0].CurrentHealth,
								CurrentFatigue:      data.DungeonCharacterRecs[0].CurrentFatigue,
							},
							Monster:         nil,
							EquippedObject:  nil,
							StashedObject:   nil,
							TargetObject:    nil,
							TargetCharacter: nil,
							TargetMonster:   nil,
							TargetLocation: &schema.LocationData{
								Name:        data.DungeonLocationRecs[1].Name,
								Description: data.DungeonLocationRecs[1].Description,
								Direction:   "north",
								Directions:  []string{"north", "south", "northwest"},
								Monsters: []schema.MonsterData{
									{
										Name: data.DungeonMonsterRecs[1].Name,
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
				Skip:              false,
				Name:              "Look at an item in the current room",
				HandlerConfig:     testCaseHandlerConfig,
				RequestHeaders:    testCaseRequestHeaders,
				RequestPathParams: testCaseRequestPathParams,
				RequestBody: func(data harness.Data) interface{} {
					res := schema.DungeonActionRequest{
						Data: schema.DungeonActionRequestData{
							Sentence: "look rusted sword",
						},
					}
					return &res
				},
				ResponseBody: testCaseResponseBody,
				ResponseCode: http.StatusOK,
			},
			expectResponseBody: func(data harness.Data) *schema.DungeonActionResponse {
				res := schema.DungeonActionResponse{
					Data: []schema.DungeonActionResponseData{
						{
							Command: "look",
							Location: schema.LocationData{
								Name:        data.DungeonLocationRecs[0].Name,
								Description: data.DungeonLocationRecs[0].Description,
								Directions:  []string{"north"},
								Characters: []schema.CharacterData{
									{
										Name: data.DungeonCharacterRecs[0].Name,
									},
								},
								Monsters: []schema.MonsterData{
									{
										Name: data.DungeonMonsterRecs[0].Name,
									},
								},
								Objects: []schema.ObjectData{
									{
										Name: data.DungeonObjectRecs[0].Name,
									},
								},
							},
							Character: &schema.CharacterDetailedData{
								Name:                data.DungeonCharacterRecs[0].Name,
								Strength:            data.DungeonCharacterRecs[0].Strength,
								Dexterity:           data.DungeonCharacterRecs[0].Dexterity,
								Intelligence:        data.DungeonCharacterRecs[0].Intelligence,
								CurrentStrength:     data.DungeonCharacterRecs[0].CurrentStrength,
								CurrentDexterity:    data.DungeonCharacterRecs[0].CurrentDexterity,
								CurrentIntelligence: data.DungeonCharacterRecs[0].CurrentIntelligence,
								Health:              data.DungeonCharacterRecs[0].Health,
								Fatigue:             data.DungeonCharacterRecs[0].Fatigue,
								CurrentHealth:       data.DungeonCharacterRecs[0].CurrentHealth,
								CurrentFatigue:      data.DungeonCharacterRecs[0].CurrentFatigue,
							},
							Monster:        nil,
							EquippedObject: nil,
							StashedObject:  nil,
							TargetObject: &schema.ObjectDetailedData{
								Name:        data.DungeonObjectRecs[0].Name,
								Description: data.DungeonObjectRecs[0].Description,
								IsEquipped:  data.DungeonObjectRecs[0].IsEquipped,
								IsStashed:   data.DungeonObjectRecs[0].IsStashed,
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
				Skip:              false,
				Name:              "Look at a monster in the current room",
				HandlerConfig:     testCaseHandlerConfig,
				RequestHeaders:    testCaseRequestHeaders,
				RequestPathParams: testCaseRequestPathParams,
				RequestBody: func(data harness.Data) interface{} {
					res := schema.DungeonActionRequest{
						Data: schema.DungeonActionRequestData{
							Sentence: "look Grumpy Dwarf",
						},
					}
					return &res
				},
				ResponseBody: testCaseResponseBody,
				ResponseCode: http.StatusOK,
			},
			expectResponseBody: func(data harness.Data) *schema.DungeonActionResponse {
				res := schema.DungeonActionResponse{
					Data: []schema.DungeonActionResponseData{
						{
							Command: "look",
							Location: schema.LocationData{
								Name:        data.DungeonLocationRecs[0].Name,
								Description: data.DungeonLocationRecs[0].Description,
								Directions:  []string{"north"},
								Characters: []schema.CharacterData{
									{
										Name: data.DungeonCharacterRecs[0].Name,
									},
								},
								Monsters: []schema.MonsterData{
									{
										Name: data.DungeonMonsterRecs[0].Name,
									},
								},
								Objects: []schema.ObjectData{
									{
										Name: data.DungeonObjectRecs[0].Name,
									},
								},
							},
							Character: &schema.CharacterDetailedData{
								Name:                data.DungeonCharacterRecs[0].Name,
								Strength:            data.DungeonCharacterRecs[0].Strength,
								Dexterity:           data.DungeonCharacterRecs[0].Dexterity,
								Intelligence:        data.DungeonCharacterRecs[0].Intelligence,
								CurrentStrength:     data.DungeonCharacterRecs[0].CurrentStrength,
								CurrentDexterity:    data.DungeonCharacterRecs[0].CurrentDexterity,
								CurrentIntelligence: data.DungeonCharacterRecs[0].CurrentIntelligence,
								Health:              data.DungeonCharacterRecs[0].Health,
								Fatigue:             data.DungeonCharacterRecs[0].Fatigue,
								CurrentHealth:       data.DungeonCharacterRecs[0].CurrentHealth,
								CurrentFatigue:      data.DungeonCharacterRecs[0].CurrentFatigue,
							},
							Monster:         nil,
							EquippedObject:  nil,
							StashedObject:   nil,
							TargetObject:    nil,
							TargetCharacter: nil,
							TargetMonster: &schema.MonsterDetailedData{
								Name:                data.DungeonMonsterRecs[0].Name,
								Strength:            data.DungeonMonsterRecs[0].Strength,
								Dexterity:           data.DungeonMonsterRecs[0].Dexterity,
								Intelligence:        data.DungeonMonsterRecs[0].Intelligence,
								CurrentStrength:     data.DungeonMonsterRecs[0].CurrentStrength,
								CurrentDexterity:    data.DungeonMonsterRecs[0].CurrentDexterity,
								CurrentIntelligence: data.DungeonMonsterRecs[0].CurrentIntelligence,
								Health:              data.DungeonMonsterRecs[0].Health,
								Fatigue:             data.DungeonMonsterRecs[0].Fatigue,
								CurrentHealth:       data.DungeonMonsterRecs[0].CurrentHealth,
								CurrentFatigue:      data.DungeonMonsterRecs[0].CurrentFatigue,
							},
							TargetLocation: nil,
						},
					},
				}
				return &res
			},
		},
		{
			// Look at a character in the current room
			TestCase: TestCase{
				Skip:              false,
				Name:              "Look at a character in the current room",
				HandlerConfig:     testCaseHandlerConfig,
				RequestHeaders:    testCaseRequestHeaders,
				RequestPathParams: testCaseRequestPathParams,
				RequestBody: func(data harness.Data) interface{} {
					res := schema.DungeonActionRequest{
						Data: schema.DungeonActionRequestData{
							Sentence: "look barricade",
						},
					}
					return &res
				},
				ResponseBody: testCaseResponseBody,
				ResponseCode: http.StatusOK,
			},
			expectResponseBody: func(data harness.Data) *schema.DungeonActionResponse {
				res := schema.DungeonActionResponse{
					Data: []schema.DungeonActionResponseData{
						{
							Command: "look",
							Location: schema.LocationData{
								Name:        data.DungeonLocationRecs[0].Name,
								Description: data.DungeonLocationRecs[0].Description,
								Directions:  []string{"north"},
								Characters: []schema.CharacterData{
									{
										Name: data.DungeonCharacterRecs[0].Name,
									},
								},
								Monsters: []schema.MonsterData{
									{
										Name: data.DungeonMonsterRecs[0].Name,
									},
								},
								Objects: []schema.ObjectData{
									{
										Name: data.DungeonObjectRecs[0].Name,
									},
								},
							},
							Character: &schema.CharacterDetailedData{
								Name:                data.DungeonCharacterRecs[0].Name,
								Strength:            data.DungeonCharacterRecs[0].Strength,
								Dexterity:           data.DungeonCharacterRecs[0].Dexterity,
								Intelligence:        data.DungeonCharacterRecs[0].Intelligence,
								CurrentStrength:     data.DungeonCharacterRecs[0].CurrentStrength,
								CurrentDexterity:    data.DungeonCharacterRecs[0].CurrentDexterity,
								CurrentIntelligence: data.DungeonCharacterRecs[0].CurrentIntelligence,
								Health:              data.DungeonCharacterRecs[0].Health,
								Fatigue:             data.DungeonCharacterRecs[0].Fatigue,
								CurrentHealth:       data.DungeonCharacterRecs[0].CurrentHealth,
								CurrentFatigue:      data.DungeonCharacterRecs[0].CurrentFatigue,
							},
							Monster:        nil,
							EquippedObject: nil,
							StashedObject:  nil,
							TargetObject:   nil,
							TargetCharacter: &schema.CharacterDetailedData{
								Name:                data.DungeonCharacterRecs[0].Name,
								Strength:            data.DungeonCharacterRecs[0].Strength,
								Dexterity:           data.DungeonCharacterRecs[0].Dexterity,
								Intelligence:        data.DungeonCharacterRecs[0].Intelligence,
								CurrentStrength:     data.DungeonCharacterRecs[0].CurrentStrength,
								CurrentDexterity:    data.DungeonCharacterRecs[0].CurrentDexterity,
								CurrentIntelligence: data.DungeonCharacterRecs[0].CurrentIntelligence,
								Health:              data.DungeonCharacterRecs[0].Health,
								Fatigue:             data.DungeonCharacterRecs[0].Fatigue,
								CurrentHealth:       data.DungeonCharacterRecs[0].CurrentHealth,
								CurrentFatigue:      data.DungeonCharacterRecs[0].CurrentFatigue,
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
				Skip:              false,
				Name:              "Stash object that is in the current room",
				HandlerConfig:     testCaseHandlerConfig,
				RequestHeaders:    testCaseRequestHeaders,
				RequestPathParams: testCaseRequestPathParams,
				RequestBody: func(data harness.Data) interface{} {
					res := schema.DungeonActionRequest{
						Data: schema.DungeonActionRequestData{
							Sentence: "stash rusted sword",
						},
					}
					return &res
				},
				ResponseBody: testCaseResponseBody,
				ResponseCode: http.StatusOK,
			},
			expectResponseBody: func(data harness.Data) *schema.DungeonActionResponse {
				res := schema.DungeonActionResponse{
					Data: []schema.DungeonActionResponseData{
						{
							Command: "stash",
							Location: schema.LocationData{
								Name:        data.DungeonLocationRecs[0].Name,
								Description: data.DungeonLocationRecs[0].Description,
								Directions:  []string{"north"},
								Characters: []schema.CharacterData{
									{
										Name: data.DungeonCharacterRecs[0].Name,
									},
								},
								Monsters: []schema.MonsterData{
									{
										Name: data.DungeonMonsterRecs[0].Name,
									},
								},
								Objects: []schema.ObjectData{},
							},
							Character: &schema.CharacterDetailedData{
								Name:                data.DungeonCharacterRecs[0].Name,
								Strength:            data.DungeonCharacterRecs[0].Strength,
								Dexterity:           data.DungeonCharacterRecs[0].Dexterity,
								Intelligence:        data.DungeonCharacterRecs[0].Intelligence,
								CurrentStrength:     data.DungeonCharacterRecs[0].CurrentStrength,
								CurrentDexterity:    data.DungeonCharacterRecs[0].CurrentDexterity,
								CurrentIntelligence: data.DungeonCharacterRecs[0].CurrentIntelligence,
								Health:              data.DungeonCharacterRecs[0].Health,
								Fatigue:             data.DungeonCharacterRecs[0].Fatigue,
								CurrentHealth:       data.DungeonCharacterRecs[0].CurrentHealth,
								CurrentFatigue:      data.DungeonCharacterRecs[0].CurrentFatigue,
							},
							Monster:        nil,
							EquippedObject: nil,
							StashedObject: &schema.ObjectDetailedData{
								Name:        data.DungeonObjectRecs[0].Name,
								Description: data.DungeonObjectRecs[0].Description,
								IsEquipped:  false,
								IsStashed:   true,
							},
							DroppedObject: nil,
							TargetObject: &schema.ObjectDetailedData{
								Name:        data.DungeonObjectRecs[0].Name,
								Description: data.DungeonObjectRecs[0].Description,
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
				Skip:              false,
				Name:              "Equip object that is in the current room",
				HandlerConfig:     testCaseHandlerConfig,
				RequestHeaders:    testCaseRequestHeaders,
				RequestPathParams: testCaseRequestPathParams,
				RequestBody: func(data harness.Data) interface{} {
					res := schema.DungeonActionRequest{
						Data: schema.DungeonActionRequestData{
							Sentence: "equip rusted sword",
						},
					}
					return &res
				},
				ResponseBody: testCaseResponseBody,
				ResponseCode: http.StatusOK,
			},
			expectResponseBody: func(data harness.Data) *schema.DungeonActionResponse {
				res := schema.DungeonActionResponse{
					Data: []schema.DungeonActionResponseData{
						{
							Command: "equip",
							Location: schema.LocationData{
								Name:        data.DungeonLocationRecs[0].Name,
								Description: data.DungeonLocationRecs[0].Description,
								Directions:  []string{"north"},
								Characters: []schema.CharacterData{
									{
										Name: data.DungeonCharacterRecs[0].Name,
									},
								},
								Monsters: []schema.MonsterData{
									{
										Name: data.DungeonMonsterRecs[0].Name,
									},
								},
								Objects: []schema.ObjectData{},
							},
							Character: &schema.CharacterDetailedData{
								Name:                data.DungeonCharacterRecs[0].Name,
								Strength:            data.DungeonCharacterRecs[0].Strength,
								Dexterity:           data.DungeonCharacterRecs[0].Dexterity,
								Intelligence:        data.DungeonCharacterRecs[0].Intelligence,
								CurrentStrength:     data.DungeonCharacterRecs[0].CurrentStrength,
								CurrentDexterity:    data.DungeonCharacterRecs[0].CurrentDexterity,
								CurrentIntelligence: data.DungeonCharacterRecs[0].CurrentIntelligence,
								Health:              data.DungeonCharacterRecs[0].Health,
								Fatigue:             data.DungeonCharacterRecs[0].Fatigue,
								CurrentHealth:       data.DungeonCharacterRecs[0].CurrentHealth,
								CurrentFatigue:      data.DungeonCharacterRecs[0].CurrentFatigue,
							},
							Monster: nil,
							EquippedObject: &schema.ObjectDetailedData{
								Name:        data.DungeonObjectRecs[0].Name,
								Description: data.DungeonObjectRecs[0].Description,
								IsEquipped:  true,
								IsStashed:   false,
							},
							StashedObject: nil,
							DroppedObject: nil,
							TargetObject: &schema.ObjectDetailedData{
								Name:        data.DungeonObjectRecs[0].Name,
								Description: data.DungeonObjectRecs[0].Description,
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
				Skip:              false,
				Name:              "Drop object that is equipped",
				HandlerConfig:     testCaseHandlerConfig,
				RequestHeaders:    testCaseRequestHeaders,
				RequestPathParams: testCaseRequestPathParams,
				RequestBody: func(data harness.Data) interface{} {
					res := schema.DungeonActionRequest{
						Data: schema.DungeonActionRequestData{
							Sentence: "drop Dull Bronze Ring",
						},
					}
					return &res
				},
				ResponseBody: testCaseResponseBody,
				ResponseCode: http.StatusOK,
			},
			expectResponseBody: func(data harness.Data) *schema.DungeonActionResponse {
				res := schema.DungeonActionResponse{
					Data: []schema.DungeonActionResponseData{
						{
							Command: "drop",
							Location: schema.LocationData{
								Name:        data.DungeonLocationRecs[0].Name,
								Description: data.DungeonLocationRecs[0].Description,
								Directions:  []string{"north"},
								Characters: []schema.CharacterData{
									{
										Name: data.DungeonCharacterRecs[0].Name,
									},
								},
								Monsters: []schema.MonsterData{
									{
										Name: data.DungeonMonsterRecs[0].Name,
									},
								},
								Objects: []schema.ObjectData{
									{
										Name: data.DungeonObjectRecs[0].Name,
									},
									{
										Name: data.DungeonObjectRecs[2].Name,
									},
								},
							},
							Character: &schema.CharacterDetailedData{
								Name:                data.DungeonCharacterRecs[0].Name,
								Strength:            data.DungeonCharacterRecs[0].Strength,
								Dexterity:           data.DungeonCharacterRecs[0].Dexterity,
								Intelligence:        data.DungeonCharacterRecs[0].Intelligence,
								CurrentStrength:     data.DungeonCharacterRecs[0].CurrentStrength,
								CurrentDexterity:    data.DungeonCharacterRecs[0].CurrentDexterity,
								CurrentIntelligence: data.DungeonCharacterRecs[0].CurrentIntelligence,
								Health:              data.DungeonCharacterRecs[0].Health,
								Fatigue:             data.DungeonCharacterRecs[0].Fatigue,
								CurrentHealth:       data.DungeonCharacterRecs[0].CurrentHealth,
								CurrentFatigue:      data.DungeonCharacterRecs[0].CurrentFatigue,
							},
							Monster:        nil,
							EquippedObject: nil,
							StashedObject:  nil,
							DroppedObject: &schema.ObjectDetailedData{
								Name:        data.DungeonObjectRecs[2].Name,
								Description: data.DungeonObjectRecs[2].Description,
								IsEquipped:  false,
								IsStashed:   false,
							},
							TargetObject: &schema.ObjectDetailedData{
								Name:        data.DungeonObjectRecs[2].Name,
								Description: data.DungeonObjectRecs[2].Description,
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
				Skip:              false,
				Name:              "Submit nothing",
				HandlerConfig:     testCaseHandlerConfig,
				RequestHeaders:    testCaseRequestHeaders,
				RequestPathParams: testCaseRequestPathParams,
				RequestBody: func(data harness.Data) interface{} {
					res := schema.DungeonActionRequest{
						Data: schema.DungeonActionRequestData{
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

	isCharacterNil := func(c *schema.CharacterDetailedData) bool {
		return c == nil
	}
	isMonsterNil := func(m *schema.MonsterDetailedData) bool {
		return m == nil
	}
	isObjectNil := func(o *schema.ObjectDetailedData) bool {
		return o == nil
	}

	for _, testCase := range testCases {

		if testCase.Skip {
			t.Logf("Skipping test >%s<", testCase.Name)
			continue
		}

		t.Logf("Running test >%s<", testCase.Name)

		testFunc := func(method string, body interface{}) {

			if testCase.TestResponseCode() != http.StatusOK {
				return
			}

			var responseBody *schema.DungeonActionResponse
			if body != nil {
				responseBody = body.(*schema.DungeonActionResponse)
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

					// Current location objects
					t.Logf("Checking object count >%d< >%d<", len(expectData.Location.Objects), len(responseBody.Data[idx].Location.Objects))
					require.Equal(t, len(expectData.Location.Objects), len(responseBody.Data[idx].Location.Objects), "Response objects count equals expected")
					if len(expectData.Location.Objects) > 0 {
						for oIdx, object := range expectData.Location.Objects {
							t.Logf("Checking object name >%s< >%s<", object.Name, responseBody.Data[idx].Location.Objects[oIdx].Name)
							require.Equal(t, object.Name, responseBody.Data[idx].Location.Objects[oIdx].Name, "Object name equals expected")
						}
					}
					if len(expectData.Location.Objects) == 0 {
						require.Equal(t, 0, len(responseBody.Data[idx].Location.Objects), "Location objects length is 0")
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
							for oIdx, object := range expectData.TargetLocation.Objects {
								t.Logf("Checking object name >%s< >%s<", object.Name, responseBody.Data[idx].TargetLocation.Objects[oIdx].Name)
								require.Equal(t, object.Name, responseBody.Data[idx].TargetLocation.Objects[oIdx].Name, "Object name equals expected")
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

						t.Logf("Checking action character health >%d< >%d<", expectData.Character.Health, responseBody.Data[idx].Character.Health)
						require.Equal(t, expectData.Character.Health, responseBody.Data[idx].Character.Health, "Response character health equals expected")
						t.Logf("Checking action character fatigue >%d< >%d<", expectData.Character.Fatigue, responseBody.Data[idx].Character.Fatigue)
						require.Equal(t, expectData.Character.Fatigue, responseBody.Data[idx].Character.Fatigue, "Response character fatigue equals expected")
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
						require.Equal(t, expectData.TargetCharacter.CurrentStrength, responseBody.Data[idx].TargetCharacter.CurrentStrength, "Response character current strength equals expected")
						t.Logf("Checking target character current dexterity >%d< >%d<", expectData.TargetCharacter.CurrentDexterity, responseBody.Data[idx].TargetCharacter.CurrentDexterity)
						require.Equal(t, expectData.TargetCharacter.CurrentDexterity, responseBody.Data[idx].TargetCharacter.CurrentDexterity, "Response character current dexterity equals expected")
						t.Logf("Checking target character current intelligence >%d< >%d<", expectData.TargetCharacter.CurrentIntelligence, responseBody.Data[idx].TargetCharacter.CurrentIntelligence)
						require.Equal(t, expectData.TargetCharacter.CurrentIntelligence, responseBody.Data[idx].TargetCharacter.CurrentIntelligence, "Response character current intelligence equals expected")

						t.Logf("Checking target character health >%d< >%d<", expectData.TargetCharacter.Health, responseBody.Data[idx].TargetCharacter.Health)
						require.Equal(t, expectData.TargetCharacter.Health, responseBody.Data[idx].TargetCharacter.Health, "Response target character health equals expected")
						t.Logf("Checking target character fatigue >%d< >%d<", expectData.TargetCharacter.Fatigue, responseBody.Data[idx].TargetCharacter.Fatigue)
						require.Equal(t, expectData.TargetCharacter.Fatigue, responseBody.Data[idx].TargetCharacter.Fatigue, "Response target character fatigue equals expected")

						t.Logf("Checking target character current health >%d< >%d<", expectData.TargetCharacter.CurrentHealth, responseBody.Data[idx].TargetCharacter.CurrentHealth)
						require.Equal(t, expectData.TargetCharacter.CurrentHealth, responseBody.Data[idx].TargetCharacter.CurrentHealth, "Response target character current health equals expected")
						t.Logf("Checking target character current fatigue >%d< >%d<", expectData.TargetCharacter.CurrentFatigue, responseBody.Data[idx].TargetCharacter.CurrentFatigue)
						require.Equal(t, expectData.TargetCharacter.CurrentFatigue, responseBody.Data[idx].TargetCharacter.CurrentFatigue, "Response target character current fatigue equals expected")
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
					th.AddDungeonCharacterActionTeardownID(data.ID)
				}
			}
		}

		RunTestCase(t, th, &testCase, testFunc)
	}
}
