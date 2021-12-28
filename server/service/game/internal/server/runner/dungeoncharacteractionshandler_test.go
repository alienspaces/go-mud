package runner

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"gitlab.com/alienspaces/go-mud/server/schema"

	"github.com/stretchr/testify/require"
	"gitlab.com/alienspaces/go-mud/server/core/auth"
	coreschema "gitlab.com/alienspaces/go-mud/server/core/schema"
	"gitlab.com/alienspaces/go-mud/server/core/server"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/harness"
)

func TestDungeonCharacterActionHandler(t *testing.T) {

	// Test harness
	th, err := NewTestHarness()
	require.NoError(t, err, "New test data returns without error")

	type testCase struct {
		TestCase
		responseBody func(data harness.Data) *schema.DungeonActionResponse
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
				ResponseCode: http.StatusOK,
			},
			responseBody: func(data harness.Data) *schema.DungeonActionResponse {
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
							Character: &schema.CharacterData{
								Name: data.DungeonCharacterRecs[0].Name,
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
				ResponseCode: http.StatusOK,
			},
			responseBody: func(data harness.Data) *schema.DungeonActionResponse {
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
							Character: &schema.CharacterData{
								Name: data.DungeonCharacterRecs[0].Name,
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
				ResponseCode: http.StatusOK,
			},
			responseBody: func(data harness.Data) *schema.DungeonActionResponse {
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
							Character: &schema.CharacterData{
								Name: data.DungeonCharacterRecs[0].Name,
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
				ResponseCode: http.StatusOK,
			},
			responseBody: func(data harness.Data) *schema.DungeonActionResponse {
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
							Character: &schema.CharacterData{
								Name: data.DungeonCharacterRecs[0].Name,
							},
							Monster:        nil,
							EquippedObject: nil,
							StashedObject:  nil,
							TargetObject: &schema.ObjectData{
								Name:        data.DungeonObjectRecs[0].Name,
								Description: data.DungeonObjectRecs[0].Description,
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
				ResponseCode: http.StatusOK,
			},
			responseBody: func(data harness.Data) *schema.DungeonActionResponse {
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
							Character: &schema.CharacterData{
								Name: data.DungeonCharacterRecs[0].Name,
							},
							Monster:         nil,
							EquippedObject:  nil,
							StashedObject:   nil,
							TargetObject:    nil,
							TargetCharacter: nil,
							TargetMonster: &schema.MonsterData{
								Name: data.DungeonMonsterRecs[0].Name,
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
				ResponseCode: http.StatusBadRequest,
			},
		},
	}

	v, err := coreschema.NewValidator(th.Config, th.Log)
	require.NoError(t, err, "Validator returns without error")

	for _, testCase := range testCases {

		t.Logf("Running test >%s<", testCase.Name)

		testFunc := func(method string, body io.Reader) {

			if testCase.TestResponseCode() != http.StatusOK {
				return
			}

			var responseBody *schema.DungeonActionResponse
			err = json.NewDecoder(body).Decode(&responseBody)
			require.NoError(t, err, "Decode returns without error")

			isCharacterNil := func(c *schema.CharacterData) bool {
				return c == nil
			}
			isMonsterNil := func(c *schema.MonsterData) bool {
				return c == nil
			}

			if testCase.responseBody != nil {

				// Validate response body
				require.NotNil(t, responseBody, "Response body is not nil")

				jsonData, err := json.Marshal(responseBody)
				require.NoError(t, err, "Marshal returns without error")

				err = v.Validate(coreschema.Config{
					Key:      "dungeonaction.create.response",
					Location: "dungeonaction",
					Main:     "create.response.schema.json",
					References: []string{
						"data.schema.json",
					},
				}, string(jsonData))
				require.NoError(t, err, "Validates against schema without error")

				expectResponseBody := testCase.responseBody(th.Data)
				for idx, expectData := range expectResponseBody.Data {

					// Response data
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
						require.Equal(t, expectData.Character.Name, responseBody.Data[idx].Character.Name, "Response monster name equals expected")
					}

					// Response monster
					t.Logf("Checking monster nil >%t< >%t<", isMonsterNil(expectData.Monster), isMonsterNil(responseBody.Data[idx].Monster))
					require.Equal(t, isMonsterNil(expectData.Monster), isMonsterNil(responseBody.Data[idx].Monster), "Response monster is nil or not nil as expected")
					if !isMonsterNil(expectData.Monster) {
						t.Logf("Checking action monster name >%s< >%s<", expectData.Monster.Name, responseBody.Data[idx].Monster.Name)
						require.Equal(t, expectData.Monster.Name, responseBody.Data[idx].Monster.Name, "Response character name equals expected")
					}
				}
			}

			require.NotNil(t, responseBody, "Response body is not nil")
			require.GreaterOrEqual(t, len(responseBody.Data), 0, "Response body data ")

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
