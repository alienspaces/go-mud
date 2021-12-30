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
			TestCase: TestCase{
				Name:              "POST - look",
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
								Name:         data.DungeonCharacterRecs[0].Name,
								Strength:     data.DungeonCharacterRecs[0].Strength,
								Dexterity:    data.DungeonCharacterRecs[0].Dexterity,
								Intelligence: data.DungeonCharacterRecs[0].Intelligence,
								Health:       data.DungeonCharacterRecs[0].Health,
								Fatigue:      data.DungeonCharacterRecs[0].Fatigue,
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
			TestCase: TestCase{
				Name:              "POST - move north",
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
								Directions:  []string{"north", "south"},
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
								Name:         data.DungeonCharacterRecs[0].Name,
								Strength:     data.DungeonCharacterRecs[0].Strength,
								Dexterity:    data.DungeonCharacterRecs[0].Dexterity,
								Intelligence: data.DungeonCharacterRecs[0].Intelligence,
								Health:       data.DungeonCharacterRecs[0].Health,
								Fatigue:      data.DungeonCharacterRecs[0].Fatigue,
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
								Directions:  []string{"north", "south"},
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
			TestCase: TestCase{
				Name:              "POST - look north",
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
								Name:         data.DungeonCharacterRecs[0].Name,
								Strength:     data.DungeonCharacterRecs[0].Strength,
								Dexterity:    data.DungeonCharacterRecs[0].Dexterity,
								Intelligence: data.DungeonCharacterRecs[0].Intelligence,
								Health:       data.DungeonCharacterRecs[0].Health,
								Fatigue:      data.DungeonCharacterRecs[0].Fatigue,
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
								Directions:  []string{"north", "south"},
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
			TestCase: TestCase{
				Name:              "POST - look rusted sword",
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
								Name:         data.DungeonCharacterRecs[0].Name,
								Strength:     data.DungeonCharacterRecs[0].Strength,
								Dexterity:    data.DungeonCharacterRecs[0].Dexterity,
								Intelligence: data.DungeonCharacterRecs[0].Intelligence,
								Health:       data.DungeonCharacterRecs[0].Health,
								Fatigue:      data.DungeonCharacterRecs[0].Fatigue,
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
			TestCase: TestCase{
				Name:              "POST - look white cat",
				HandlerConfig:     testCaseHandlerConfig,
				RequestHeaders:    testCaseRequestHeaders,
				RequestPathParams: testCaseRequestPathParams,
				RequestBody: func(data harness.Data) interface{} {
					res := schema.DungeonActionRequest{
						Data: schema.DungeonActionRequestData{
							Sentence: "look white cat",
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
								Name:         data.DungeonCharacterRecs[0].Name,
								Strength:     data.DungeonCharacterRecs[0].Strength,
								Dexterity:    data.DungeonCharacterRecs[0].Dexterity,
								Intelligence: data.DungeonCharacterRecs[0].Intelligence,
								Health:       data.DungeonCharacterRecs[0].Health,
								Fatigue:      data.DungeonCharacterRecs[0].Fatigue,
							},
							Monster:         nil,
							EquippedObject:  nil,
							StashedObject:   nil,
							TargetObject:    nil,
							TargetCharacter: nil,
							TargetMonster: &schema.MonsterDetailedData{
								Name:         data.DungeonMonsterRecs[0].Name,
								Strength:     data.DungeonMonsterRecs[0].Strength,
								Dexterity:    data.DungeonMonsterRecs[0].Dexterity,
								Intelligence: data.DungeonMonsterRecs[0].Intelligence,
								Health:       data.DungeonMonsterRecs[0].Health,
								Fatigue:      data.DungeonMonsterRecs[0].Fatigue,
							},
							TargetLocation: nil,
						},
					},
				}
				return &res
			},
		},
		{
			TestCase: TestCase{
				Name:              "POST - empty",
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

					// Current location monsters
					t.Logf("Checking monster count >%d< >%d<", len(expectData.Location.Monsters), len(responseBody.Data[idx].Location.Monsters))
					require.Equal(t, len(expectData.Location.Monsters), len(responseBody.Data[idx].Location.Monsters), "Response monsters count equals expected")
					if len(expectData.Location.Monsters) > 0 {
						for mIdx, monster := range expectData.Location.Monsters {
							t.Logf("Checking monster name >%s< >%s<", monster.Name, responseBody.Data[idx].Location.Monsters[mIdx].Name)
							require.Equal(t, monster.Name, responseBody.Data[idx].Location.Monsters[mIdx].Name, "Monster name equals expected")
						}
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
						t.Logf("Checking action target character health >%d< >%d<", expectData.TargetCharacter.Health, responseBody.Data[idx].TargetCharacter.Health)
						require.Equal(t, expectData.TargetCharacter.Health, responseBody.Data[idx].TargetCharacter.Health, "Response target character health equals expected")
						t.Logf("Checking action target character fatigue >%d< >%d<", expectData.TargetCharacter.Fatigue, responseBody.Data[idx].TargetCharacter.Fatigue)
						require.Equal(t, expectData.TargetCharacter.Fatigue, responseBody.Data[idx].TargetCharacter.Fatigue, "Response target character fatigue equals expected")
					}

					// Response monster
					t.Logf("Checking monster nil >%t< >%t<", isMonsterNil(expectData.Monster), isMonsterNil(responseBody.Data[idx].Monster))
					require.Equal(t, isMonsterNil(expectData.Monster), isMonsterNil(responseBody.Data[idx].Monster), "Response monster is nil or not nil as expected")
					if !isMonsterNil(expectData.Monster) {
						t.Logf("Checking action monster name >%s< >%s<", expectData.Monster.Name, responseBody.Data[idx].Monster.Name)
						require.Equal(t, expectData.Monster.Name, responseBody.Data[idx].Monster.Name, "Response monster name equals expected")
						t.Logf("Checking action monster strenth >%d< >%d<", expectData.Monster.Strength, responseBody.Data[idx].Monster.Strength)
						require.Equal(t, expectData.Monster.Strength, responseBody.Data[idx].Monster.Strength, "Response monster strength equals expected")
						t.Logf("Checking action monster dexterity >%d< >%d<", expectData.Monster.Dexterity, responseBody.Data[idx].Monster.Dexterity)
						require.Equal(t, expectData.Monster.Dexterity, responseBody.Data[idx].Monster.Dexterity, "Response monster dexterity equals expected")
						t.Logf("Checking action monster intelligence >%d< >%d<", expectData.Monster.Intelligence, responseBody.Data[idx].Monster.Intelligence)
						require.Equal(t, expectData.Monster.Intelligence, responseBody.Data[idx].Monster.Intelligence, "Response monster intelligence equals expected")
						t.Logf("Checking action monster health >%d< >%d<", expectData.Monster.Health, responseBody.Data[idx].Monster.Health)
						require.Equal(t, expectData.Monster.Health, responseBody.Data[idx].Monster.Health, "Response monster health equals expected")
						t.Logf("Checking action monster fatigue >%d< >%d<", expectData.Monster.Fatigue, responseBody.Data[idx].Monster.Fatigue)
						require.Equal(t, expectData.Monster.Fatigue, responseBody.Data[idx].Monster.Fatigue, "Response monster fatigue equals expected")
					}

					// Response target monster
					t.Logf("Checking target monster nil >%t< >%t<", isMonsterNil(expectData.TargetMonster), isMonsterNil(responseBody.Data[idx].TargetMonster))
					require.Equal(t, isMonsterNil(expectData.TargetMonster), isMonsterNil(responseBody.Data[idx].TargetMonster), "Response target monster is nil or not nil as expected")
					if !isMonsterNil(expectData.TargetMonster) {
						t.Logf("Checking action target monster name >%s< >%s<", expectData.TargetMonster.Name, responseBody.Data[idx].TargetMonster.Name)
						require.Equal(t, expectData.TargetMonster.Name, responseBody.Data[idx].TargetMonster.Name, "Response target monster name equals expected")
						t.Logf("Checking action target monster strenth >%d< >%d<", expectData.TargetMonster.Strength, responseBody.Data[idx].TargetMonster.Strength)
						require.Equal(t, expectData.TargetMonster.Strength, responseBody.Data[idx].TargetMonster.Strength, "Response target monster strength equals expected")
						t.Logf("Checking action target monster dexterity >%d< >%d<", expectData.TargetMonster.Dexterity, responseBody.Data[idx].TargetMonster.Dexterity)
						require.Equal(t, expectData.TargetMonster.Dexterity, responseBody.Data[idx].TargetMonster.Dexterity, "Response target monster dexterity equals expected")
						t.Logf("Checking action target monster intelligence >%d< >%d<", expectData.TargetMonster.Intelligence, responseBody.Data[idx].TargetMonster.Intelligence)
						require.Equal(t, expectData.TargetMonster.Intelligence, responseBody.Data[idx].TargetMonster.Intelligence, "Response target monster intelligence equals expected")
						t.Logf("Checking action target monster health >%d< >%d<", expectData.TargetMonster.Health, responseBody.Data[idx].TargetMonster.Health)
						require.Equal(t, expectData.TargetMonster.Health, responseBody.Data[idx].TargetMonster.Health, "Response target monster health equals expected")
						t.Logf("Checking action target monster fatigue >%d< >%d<", expectData.TargetMonster.Fatigue, responseBody.Data[idx].TargetMonster.Fatigue)
						require.Equal(t, expectData.TargetMonster.Fatigue, responseBody.Data[idx].TargetMonster.Fatigue, "Response target monster fatigue equals expected")
					}

					// Response target object
					t.Logf("Checking target object nil >%t< >%t<", isObjectNil(expectData.TargetObject), isObjectNil(responseBody.Data[idx].TargetObject))
					require.Equal(t, isObjectNil(expectData.TargetObject), isObjectNil(responseBody.Data[idx].TargetObject), "Response target object is nil or not nil as expected")
					if !isObjectNil(expectData.TargetObject) {
						t.Logf("Checking action target object name >%s< >%s<", expectData.TargetObject.Name, responseBody.Data[idx].TargetObject.Name)
						require.Equal(t, expectData.TargetObject.Name, responseBody.Data[idx].TargetObject.Name, "Response target object name equals expected")
						t.Logf("Checking action target object description >%s< >%s<", expectData.TargetObject.Description, responseBody.Data[idx].TargetObject.Description)
						require.Equal(t, expectData.TargetObject.Description, responseBody.Data[idx].TargetObject.Description, "Response target object description equals expected")
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
