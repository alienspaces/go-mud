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

func TestGetDungeonLocationHandler(t *testing.T) {

	// Test harness
	th, err := newTestHarness()
	require.NoError(t, err, "New test data returns without error")

	_, err = th.Setup()
	require.NoError(t, err, "Test data setup returns without error")
	defer func() {
		err = th.Teardown()
		require.NoError(t, err, "Test data teardown returns without error")
	}()

	type testCase struct {
		TestCase
		expectResponseBody func(data harness.Data) *schema.LocationResponse
	}

	testCaseHandlerConfig := func(rnr *Runner) server.HandlerConfig {
		return rnr.HandlerConfig[getDungeonLocation]
	}

	testCaseResponseDecoder := func(body io.Reader) (interface{}, error) {
		var responseBody *schema.LocationResponse
		err = json.NewDecoder(body).Decode(&responseBody)
		return responseBody, err
	}

	testCases := []testCase{
		{
			TestCase: TestCase{
				Name:          "GET - Get existing",
				HandlerConfig: testCaseHandlerConfig,
				RequestPathParams: func(data harness.Data) map[string]string {
					params := map[string]string{
						":dungeon_id":  data.DungeonRecs[0].ID,
						":location_id": data.LocationRecs[0].ID,
					}
					return params
				},
				RequestBody: func(data harness.Data) interface{} {
					return nil
				},
				ResponseDecoder: testCaseResponseDecoder,
				ResponseCode:    http.StatusOK,
			},
			expectResponseBody: func(data harness.Data) *schema.LocationResponse {
				res := schema.LocationResponse{
					Data: []schema.LocationData{
						{
							ID:          data.LocationRecs[0].ID,
							Name:        data.LocationRecs[0].Name,
							Description: data.LocationRecs[0].Description,
						},
					},
				}
				return &res
			},
		},
		{
			TestCase: TestCase{
				Name:          "GET - Get with non-existant dungeon",
				HandlerConfig: testCaseHandlerConfig,
				RequestPathParams: func(data harness.Data) map[string]string {
					params := map[string]string{
						":dungeon_id":  "dc73cc41-0c6d-4dc1-b79d-c41ea2761304",
						":location_id": data.LocationRecs[0].ID,
					}
					return params
				},
				RequestBody: func(data harness.Data) interface{} {
					return nil
				},
				ResponseCode: http.StatusNotFound,
			},
		},
		{
			TestCase: TestCase{
				Name:          "GET - Get with non-existant location",
				HandlerConfig: testCaseHandlerConfig,
				RequestPathParams: func(data harness.Data) map[string]string {
					params := map[string]string{
						":dungeon_id":  data.DungeonRecs[0].ID,
						":location_id": "f3677973-cef2-4c83-96db-4abd387f321d",
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

				var responseBody *schema.LocationResponse
				if body != nil {
					responseBody = body.(*schema.LocationResponse)
				}

				// Validate response body
				if tc.expectResponseBody != nil {
					require.NotNil(t, responseBody, "Response body is not nil")
					require.GreaterOrEqual(t, len(responseBody.Data), 0, "Response body data ")

					expectResponseBody := tc.expectResponseBody(th.Data)

					// Validate response body data
					for idx, expectData := range expectResponseBody.Data {
						require.NotNil(t, responseBody.Data[idx], "Response body index is not empty")

						// Validate location
						t.Logf("Checking location name >%s< >%s<", expectData.Name, responseBody.Data[idx].Name)
						require.Equal(t, expectData.Name, responseBody.Data[idx].Name)
						t.Logf("Checking location description >%s< >%s<", expectData.Description, responseBody.Data[idx].Description)
						require.Equal(t, expectData.Description, responseBody.Data[idx].Description)
					}
				}
			}

			RunTestCase(t, th, &tc, testFunc)
		})
	}
}
