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

	testCases := []testCase{
		{
			TestCase: TestCase{
				Name: "POST - look",
				HandlerConfig: func(rnr *Runner) server.HandlerConfig {
					return rnr.HandlerConfig[6]
				},
				RequestHeaders: func(data harness.Data) map[string]string {
					headers := map[string]string{
						"Authorization": "Bearer " + validAuthToken(),
					}
					return headers
				},
				RequestPathParams: func(data harness.Data) map[string]string {
					params := map[string]string{
						":dungeon_id":   data.DungeonRecs[0].ID,
						":character_id": data.DungeonCharacterRecs[0].ID,
					}
					return params
				},
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
							Action: schema.ActionData{
								Command:                   "look",
								TargetDungeonLocationName: data.DungeonLocationRecs[0].Name,
							},
							Location: schema.LocationData{
								Name:        data.DungeonLocationRecs[0].Name,
								Description: data.DungeonLocationRecs[0].Description,
								Directions:  []string{"north"},
							},
							Character: &schema.CharacterData{
								Name: data.DungeonCharacterRecs[0].Name,
							},
						},
					},
				}
				return &res
			},
		},
		{
			TestCase: TestCase{
				Name: "POST - move north",
				HandlerConfig: func(rnr *Runner) server.HandlerConfig {
					return rnr.HandlerConfig[6]
				},
				RequestHeaders: func(data harness.Data) map[string]string {
					headers := map[string]string{
						"Authorization": "Bearer " + validAuthToken(),
					}
					return headers
				},
				RequestPathParams: func(data harness.Data) map[string]string {
					params := map[string]string{
						":dungeon_id":   data.DungeonRecs[0].ID,
						":character_id": data.DungeonCharacterRecs[0].ID,
					}
					return params
				},
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
							Action: schema.ActionData{
								Command:                   "move",
								TargetDungeonLocationName: data.DungeonLocationRecs[1].Name,
							},
							Location: schema.LocationData{
								Name:        data.DungeonLocationRecs[1].Name,
								Description: data.DungeonLocationRecs[1].Description,
								Directions:  []string{"north", "south"},
							},
							Character: &schema.CharacterData{
								Name: data.DungeonCharacterRecs[0].Name,
							},
						},
					},
				}
				return &res
			},
		},
		{
			TestCase: TestCase{
				Name: "POST - empty",
				HandlerConfig: func(rnr *Runner) server.HandlerConfig {
					return rnr.HandlerConfig[6]
				},
				RequestHeaders: func(data harness.Data) map[string]string {
					headers := map[string]string{
						"Authorization": "Bearer " + validAuthToken(),
					}
					return headers
				},
				RequestPathParams: func(data harness.Data) map[string]string {
					params := map[string]string{
						":dungeon_id":   data.DungeonRecs[0].ID,
						":character_id": data.DungeonCharacterRecs[0].ID,
					}
					return params
				},
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

	for _, testCase := range testCases {

		t.Logf("Running test >%s<", testCase.Name)

		testFunc := func(method string, body io.Reader) {

			var responseBody *schema.DungeonActionResponse
			err = json.NewDecoder(body).Decode(&responseBody)
			require.NoError(t, err, "Decode returns without error")

			if testCase.TestResponseCode() != http.StatusOK {
				return
			}

			isCharacterNil := func(c *schema.CharacterData) bool {
				return c == nil
			}
			isMonsterNil := func(c *schema.MonsterData) bool {
				return c == nil
			}

			// test data
			if testCase.responseBody != nil {
				expectResponseBody := testCase.responseBody(th.Data)
				if responseBody != nil {
					for idx, expectData := range expectResponseBody.Data {
						// Response data
						require.NotNil(t, responseBody.Data[idx], "Response body index is not empty")
						// Response action
						require.NotNil(t, responseBody.Data[idx].Action, "Response body action is not empty")
						require.Equal(t, responseBody.Data[idx].Action.Command, expectData.Action.Command)
						t.Logf("Checking action target dungeon object name >%s< >%s<", responseBody.Data[idx].Action.TargetDungeonObjectName, expectData.Action.TargetDungeonObjectName)
						require.Equal(t, responseBody.Data[idx].Action.TargetDungeonObjectName, expectData.Action.TargetDungeonObjectName)
						t.Logf("Checking action target dungeon character name >%s< >%s<", responseBody.Data[idx].Action.TargetDungeonCharacterName, expectData.Action.TargetDungeonCharacterName)
						require.Equal(t, responseBody.Data[idx].Action.TargetDungeonCharacterName, expectData.Action.TargetDungeonCharacterName)
						t.Logf("Checking action target dungeon monster name >%s< >%s<", responseBody.Data[idx].Action.TargetDungeonMonsterName, expectData.Action.TargetDungeonMonsterName)
						require.Equal(t, responseBody.Data[idx].Action.TargetDungeonMonsterName, expectData.Action.TargetDungeonMonsterName)
						t.Logf("Checking action target dungeon location name >%s< >%s<", responseBody.Data[idx].Action.TargetDungeonLocationName, expectData.Action.TargetDungeonLocationName)
						require.Equal(t, responseBody.Data[idx].Action.TargetDungeonLocationName, expectData.Action.TargetDungeonLocationName)
						// Response location
						require.NotNil(t, responseBody.Data[idx].Location, "Response body location is not empty")
						t.Logf("Checking location name >%s< >%s<", responseBody.Data[idx].Location.Name, expectData.Location.Name)
						require.Equal(t, responseBody.Data[idx].Location.Name, expectData.Location.Name)
						t.Logf("Checking location description >%s< >%s<", responseBody.Data[idx].Location.Description, expectData.Location.Description)
						require.Equal(t, responseBody.Data[idx].Location.Description, expectData.Location.Description)
						t.Logf("Checking location directions >%s< >%s<", responseBody.Data[idx].Location.Directions, expectData.Location.Directions)
						require.Equal(t, responseBody.Data[idx].Location.Directions, expectData.Location.Directions)
						// Response character
						require.Equal(t, isCharacterNil(responseBody.Data[idx].Character), isCharacterNil(expectData.Character), "Response body character is nil or not nil as expected")
						if !isCharacterNil(expectData.Character) {
							t.Logf("Checking action character name >%s< >%s<", responseBody.Data[idx].Character.Name, expectData.Character.Name)
						}

						// Response monster
						require.Equal(t, isMonsterNil(responseBody.Data[idx].Monster), isMonsterNil(expectData.Monster), "Response body character is nil or not nil as expected")
						if !isMonsterNil(expectData.Monster) {
							t.Logf("Checking action monster name >%s< >%s<", responseBody.Data[idx].Monster.Name, expectData.Monster.Name)
						}
					}
				}
			}

			require.NotNil(t, responseBody, "Response body is not nil")
			require.GreaterOrEqual(t, len(responseBody.Data), 0, "Response body data ")

			for _, data := range responseBody.Data {

				// record timestamps
				require.False(t, data.Action.CreatedAt.IsZero(), "CreatedAt is not zero")
				if method == http.MethodPost {
					require.True(t, data.Action.UpdatedAt.IsZero(), "UpdatedAt is zero")
				}
				if method == http.MethodPut {
					require.False(t, data.Action.UpdatedAt.IsZero(), "UpdatedAt is not zero")
				}

				if method == http.MethodPost {
					if len(responseBody.Data) != 0 {
						t.Logf("Adding dungeon character action ID >%s< for teardown", data.Action.ID)
						th.AddDungeonCharacterActionTeardownID(data.Action.ID)
					}
				}
			}
		}

		RunTestCase(t, th, &testCase, testFunc)
	}
}
