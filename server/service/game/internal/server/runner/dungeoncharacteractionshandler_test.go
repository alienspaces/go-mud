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
		// requestBody        func(data harness.Data) *schema.DungeonActionRequest
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
			// requestBody: func(data harness.Data) *schema.DungeonActionRequest {
			// 	res := schema.DungeonActionRequest{
			// 		Data: schema.DungeonActionRequestData{
			// 			Sentence: "look",
			// 		},
			// 	}
			// 	return &res
			// },
			expectResponseBody: func(data harness.Data) *schema.DungeonActionResponse {
				res := schema.DungeonActionResponse{
					Data: []schema.DungeonActionResponseData{
						{
							Action: schema.ActionData{
								Command:                   "look",
								TargetDungeonLocationName: "Cave Entrance",
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
			expectResponseBody: func(data harness.Data) *schema.DungeonActionResponse {
				res := schema.DungeonActionResponse{
					Data: []schema.DungeonActionResponseData{
						{
							Action: schema.ActionData{
								Command:                   "move",
								TargetDungeonLocationName: "Cave Tunnel",
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

			// test data
			if testCase.expectResponseBody != nil {
				expectResponseBody := testCase.expectResponseBody(th.Data)
				if expectResponseBody != nil {
					for idx, expectData := range expectResponseBody.Data {
						require.NotNil(t, responseBody.Data[idx], "Response body index is not empty")
						require.NotNil(t, responseBody.Data[idx].Action, "Response body action is not empty")
						require.Equal(t, responseBody.Data[idx].Action.Command, expectData.Action.Command)
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
