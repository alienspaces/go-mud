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

func TestGetDungeonLocationHandler(t *testing.T) {

	// Test harness
	th, err := NewTestHarness()
	require.NoError(t, err, "New test data returns without error")

	type testCase struct {
		TestCase
		expectResponseBody func(data harness.Data) *schema.LocationResponse
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
		return rnr.HandlerConfig[getDungeonLocation]
	}

	testCaseRequestHeaders := func(data harness.Data) map[string]string {
		headers := map[string]string{
			"Authorization": "Bearer " + validAuthToken(),
		}
		return headers
	}

	testCaseResponseBody := func(body io.Reader) (interface{}, error) {
		var responseBody *schema.LocationResponse
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
						":dungeon_id":  data.DungeonRecs[0].ID,
						":location_id": data.LocationRecs[0].ID,
					}
					return params
				},
				RequestBody: func(data harness.Data) interface{} {
					return nil
				},
				ResponseBody: testCaseResponseBody,
				ResponseCode: http.StatusOK,
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
				Name:           "GET - Get with non-existant dungeon",
				HandlerConfig:  testCaseHandlerConfig,
				RequestHeaders: testCaseRequestHeaders,
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
				Name:           "GET - Get with non-existant location",
				HandlerConfig:  testCaseHandlerConfig,
				RequestHeaders: testCaseRequestHeaders,
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

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Logf("Running test >%s<", testCase.Name)

			testFunc := func(method string, body interface{}) {

				if testCase.TestResponseCode() != http.StatusOK {
					return
				}

				var responseBody *schema.LocationResponse
				if body != nil {
					responseBody = body.(*schema.LocationResponse)
				}

				// Validate response body
				if testCase.expectResponseBody != nil {
					require.NotNil(t, responseBody, "Response body is not nil")
					require.GreaterOrEqual(t, len(responseBody.Data), 0, "Response body data ")

					expectResponseBody := testCase.expectResponseBody(th.Data)

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

			RunTestCase(t, th, &testCase, testFunc)
		})
	}
}
