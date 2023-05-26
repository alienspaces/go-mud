package runner

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/go-mud/backend/core/server"
	schema "gitlab.com/alienspaces/go-mud/backend/schema/game"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/harness"
)

func TestGetDungeonHandler(t *testing.T) {

	// Test harness
	th, err := NewTestHarness()
	require.NoError(t, err, "New test data returns without error")

	type testCase struct {
		TestCase
		expectResponseBody func(data harness.Data) *schema.DungeonResponse
	}

	testCaseHandlerConfig := func(rnr *Runner) server.HandlerConfig {
		return rnr.HandlerConfig[getDungeon]
	}

	testCaseRequestHeaders := func(data harness.Data) map[string]string {
		headers := map[string]string{
					}
		return headers
	}

	testCaseResponseBody := func(body io.Reader) (interface{}, error) {
		var responseBody *schema.DungeonResponse
		err = json.NewDecoder(body).Decode(&responseBody)
		return responseBody, err
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
				ResponseBody: testCaseResponseBody,
				ResponseCode: http.StatusOK,
			},
			expectResponseBody: func(data harness.Data) *schema.DungeonResponse {
				res := schema.DungeonResponse{
					Data: []schema.DungeonData{
						{
							DungeonID:          data.DungeonRecs[0].ID,
							DungeonName:        data.DungeonRecs[0].Name,
							DungeonDescription: data.DungeonRecs[0].Description,
						},
					},
				}
				return &res
			},
		},
		{
			TestCase: TestCase{
				Name:           "GET - Get non-existant",
				HandlerConfig:  testCaseHandlerConfig,
				RequestHeaders: testCaseRequestHeaders,
				RequestPathParams: func(data harness.Data) map[string]string {
					params := map[string]string{
						":dungeon_id": "17c19414-2d15-4d20-8fc3-36fc10341dc8",
					}
					return params
				},
				RequestBody: func(data harness.Data) interface{} {
					return nil
				},
				ResponseCode: http.StatusNotFound,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			t.Logf("Running test >%s<", tc.Name)

			testFunc := func(method string, body interface{}) {

				if tc.TestResponseCode() != http.StatusOK {
					return
				}

				var responseBody *schema.DungeonResponse
				if body != nil {
					responseBody = body.(*schema.DungeonResponse)
				}

				// Validate response body
				if tc.expectResponseBody != nil {
					require.NotNil(t, responseBody, "Response body is not nil")
					require.GreaterOrEqual(t, len(responseBody.Data), 0, "Response body data ")

					expectResponseBody := tc.expectResponseBody(th.Data)

					// Validate response body data
					for idx, expectData := range expectResponseBody.Data {
						require.NotNil(t, responseBody.Data[idx], "Response body index is not empty")

						// Validate dungeon
						t.Logf("Checking dungeon name >%s< >%s<", expectData.DungeonName, responseBody.Data[idx].DungeonName)
						require.Equal(t, expectData.DungeonName, responseBody.Data[idx].DungeonName)
						t.Logf("Checking dungeon description >%s< >%s<", expectData.DungeonDescription, responseBody.Data[idx].DungeonDescription)
						require.Equal(t, expectData.DungeonDescription, responseBody.Data[idx].DungeonDescription)
					}
				}
			}

			RunTestCase(t, th, &tc, testFunc)
		})
	}
}
