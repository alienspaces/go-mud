package runner

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/go-mud/server/core/auth"
	"gitlab.com/alienspaces/go-mud/server/core/server"
	"gitlab.com/alienspaces/go-mud/server/schema"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/harness"
)

func TestDungeonHandler(t *testing.T) {

	// Test harness
	th, err := NewTestHarness()
	require.NoError(t, err, "New test data returns without error")

	type testCase struct {
		TestCase
		responseBody func(data harness.Data) *schema.DungeonResponse
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
		return rnr.HandlerConfig[1]
	}

	testCaseRequestHeaders := func(data harness.Data) map[string]string {
		headers := map[string]string{
			"Authorization": "Bearer " + validAuthToken(),
		}
		return headers
	}

	testCases := []testCase{
		{
			TestCase: TestCase{
				Name:           "GET - Get existing",
				HandlerConfig:  testCaseHandlerConfig,
				RequestHeaders: testCaseRequestHeaders,
				RequestPathParams: func(data harness.Data) map[string]string {
					params := map[string]string{
						":dungeon_id": data.DungeonRecs[0].ID,
					}
					return params
				},
				RequestBody: func(data harness.Data) interface{} {
					return nil
				},
				ResponseCode: http.StatusOK,
			},
			responseBody: func(data harness.Data) *schema.DungeonResponse {
				res := schema.DungeonResponse{
					Data: []schema.DungeonData{
						{
							ID: data.DungeonRecs[0].ID,
						},
					},
				}
				return &res
			},
		},
		// {
		// 	name: "GET - Get non-existant",
		// 	config: func(rnr *Runner) server.HandlerConfig {
		// 		return rnr.HandlerConfig[1]
		// 	},
		// 	requestHeaders: func(data harness.Data) map[string]string {
		// 		headers := map[string]string{
		// 			"Authorization": "Bearer " + validAuthToken(),
		// 		}
		// 		return headers
		// 	},
		// 	requestParams: func(data harness.Data) map[string]string {
		// 		params := map[string]string{
		// 			":dungeon_id": "17c19414-2d15-4d20-8fc3-36fc10341dc8",
		// 		}
		// 		return params
		// 	},
		// 	requestData: func(data harness.Data) *schema.DungeonRequest {
		// 		return nil
		// 	},
		// 	responseCode: http.StatusNotFound,
		// },
	}

	// v, err := coreschema.NewValidator(th.Config, th.Log)
	// require.NoError(t, err, "Validator returns without error")

	for _, testCase := range testCases {

		t.Logf("Running test >%s<", testCase.Name)

		testFunc := func(method string, body io.Reader) {

			var responseBody *schema.DungeonResponse
			err = json.NewDecoder(body).Decode(&responseBody)
			require.NoError(t, err, "Decode returns without error")

			if testCase.TestResponseCode() != http.StatusOK {
				return
			}

			if testCase.responseBody != nil {
				expectResponseBody := testCase.responseBody(th.Data)
				if responseBody != nil {

					// Validate response body
					// jsonData, err := json.Marshal(responseBody)
					// require.NoError(t, err, "Marshal returns without error")

					// err = v.Validate(coreschema.Config{
					// 	Key:      "dungeonaction.create.response",
					// 	Location: "dungeonaction",
					// 	Main:     "create.response.schema.json",
					// 	References: []string{
					// 		"data.schema.json",
					// 	},
					// }, string(jsonData))
					// require.NoError(t, err, "Validates against schema without error")

					for idx := range expectResponseBody.Data {
						//
						// Response data
						require.NotNil(t, responseBody.Data[idx], "Response body index is not empty")
					}
				}
			}
		}

		RunTestCase(t, th, &testCase, testFunc)
	}
}
