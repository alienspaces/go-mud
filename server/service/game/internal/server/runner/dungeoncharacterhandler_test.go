package runner

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/brianvoe/gofakeit"

	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/go-mud/server/core/auth"
	"gitlab.com/alienspaces/go-mud/server/core/server"
	"gitlab.com/alienspaces/go-mud/server/schema"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/harness"
)

func TestCreateDungeonCharacterHandler(t *testing.T) {

	// Test harness
	th, err := NewTestHarness()
	require.NoError(t, err, "New test data returns without error")

	type testCase struct {
		TestCase
		expectResponseBody func(data harness.Data) *schema.DungeonCharacterResponse
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

	testCaseResponseBody := func(body io.Reader) (interface{}, error) {
		var responseBody *schema.DungeonCharacterResponse
		err = json.NewDecoder(body).Decode(&responseBody)
		return responseBody, err
	}

	testCases := []testCase{
		{
			TestCase: TestCase{
				Name: "POST - Create one with valid attributes",
				HandlerConfig: func(rnr *Runner) server.HandlerConfig {
					return rnr.HandlerConfig[4]
				},
				RequestHeaders: func(data harness.Data) map[string]string {
					headers := map[string]string{
						"Authorization": "Bearer " + validAuthToken(),
					}
					return headers
				},
				RequestPathParams: func(data harness.Data) map[string]string {
					params := map[string]string{
						":dungeon_id": data.DungeonRecs[0].ID,
					}
					return params
				},
				RequestBody: func(data harness.Data) interface{} {
					res := schema.DungeonCharacterRequest{
						Data: schema.DungeonCharacterData{
							Name:         gofakeit.Name() + gofakeit.Name(),
							Strength:     10,
							Dexterity:    10,
							Intelligence: 10,
						},
					}
					return &res
				},
				ResponseBody: testCaseResponseBody,
				ResponseCode: http.StatusOK,
			},
		},
	}

	for _, testCase := range testCases {

		t.Logf("Running test >%s<", testCase.Name)

		testFunc := func(method string, body interface{}) {

			if testCase.TestResponseCode() != http.StatusOK {
				return
			}

			var responseBody *schema.DungeonCharacterResponse
			if body != nil {
				responseBody = body.(*schema.DungeonCharacterResponse)
			}

			// Validate response body
			if testCase.expectResponseBody != nil {
				require.NotNil(t, responseBody, "Response body is not nil")
				require.GreaterOrEqual(t, len(responseBody.Data), 0, "Response body data ")
			}

			for _, data := range responseBody.Data {
				require.False(t, data.CreatedAt.IsZero(), "CreatedAt is not zero")
				if method == http.MethodPost {
					require.True(t, data.UpdatedAt.IsZero(), "UpdatedAt is zero")
					t.Logf("Adding dungeon character ID >%s< for teardown", data.ID)
					th.AddDungeonCharacterTeardownID(data.ID)
				}
			}
		}

		RunTestCase(t, th, &testCase, testFunc)
	}
}

func TestGetDungeonCharacterHandler(t *testing.T) {

	// Test harness
	th, err := NewTestHarness()
	require.NoError(t, err, "New test data returns without error")

	type testCase struct {
		TestCase
		expectResponseBody func(data harness.Data) *schema.DungeonCharacterResponse
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

	testCaseResponseBody := func(body io.Reader) (interface{}, error) {
		var responseBody *schema.DungeonCharacterResponse
		err = json.NewDecoder(body).Decode(&responseBody)
		return responseBody, err
	}

	testCases := []testCase{
		{
			TestCase: TestCase{
				Name: "GET - Get many",
				HandlerConfig: func(rnr *Runner) server.HandlerConfig {
					return rnr.HandlerConfig[2]
				},
				RequestHeaders: func(data harness.Data) map[string]string {
					headers := map[string]string{
						"Authorization": "Bearer " + validAuthToken(),
					}
					return headers
				},
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
			expectResponseBody: func(data harness.Data) *schema.DungeonCharacterResponse {
				res := schema.DungeonCharacterResponse{
					Data: []schema.DungeonCharacterData{
						{
							ID:                  data.DungeonCharacterRecs[0].ID,
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
					},
				}
				return &res
			},
		},
		{
			TestCase: TestCase{
				Name: "GET - Get many with invalid dungeon ID",
				HandlerConfig: func(rnr *Runner) server.HandlerConfig {
					return rnr.HandlerConfig[2]
				},
				RequestHeaders: func(data harness.Data) map[string]string {
					headers := map[string]string{
						"Authorization": "Bearer " + validAuthToken(),
					}
					return headers
				},
				RequestPathParams: func(data harness.Data) map[string]string {
					params := map[string]string{
						":dungeon_id": "21954f35-76fb-4a4a-bc39-fba15432b28b",
					}
					return params
				},
				RequestBody: func(data harness.Data) interface{} {
					return nil
				},
				ResponseBody: testCaseResponseBody,
				ResponseCode: http.StatusNotFound,
			},
			expectResponseBody: func(data harness.Data) *schema.DungeonCharacterResponse {
				return nil
			},
		},
		{
			TestCase: TestCase{
				Name: "GET - Get one",
				HandlerConfig: func(rnr *Runner) server.HandlerConfig {
					return rnr.HandlerConfig[3]
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
					return nil
				},
				ResponseBody: testCaseResponseBody,
				ResponseCode: http.StatusOK,
			},
			expectResponseBody: func(data harness.Data) *schema.DungeonCharacterResponse {
				res := schema.DungeonCharacterResponse{
					Data: []schema.DungeonCharacterData{
						{
							ID:                  data.DungeonCharacterRecs[0].ID,
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
					},
				}
				return &res
			},
		},
		{
			TestCase: TestCase{
				Name: "GET - Get one with invalid dungeon ID",
				HandlerConfig: func(rnr *Runner) server.HandlerConfig {
					return rnr.HandlerConfig[3]
				},
				RequestHeaders: func(data harness.Data) map[string]string {
					headers := map[string]string{
						"Authorization": "Bearer " + validAuthToken(),
					}
					return headers
				},
				RequestPathParams: func(data harness.Data) map[string]string {
					params := map[string]string{
						":dungeon_id": "17c19414-2d15-4d20-8fc3-36fc10341dc8",
					}
					return params
				},
				RequestBody: func(data harness.Data) interface{} {
					return nil
				},
				ResponseBody: testCaseResponseBody,
				ResponseCode: http.StatusNotFound,
			},
		},
		{
			TestCase: TestCase{
				Name: "GET - Get one with invalid character ID",
				HandlerConfig: func(rnr *Runner) server.HandlerConfig {
					return rnr.HandlerConfig[3]
				},
				RequestHeaders: func(data harness.Data) map[string]string {
					headers := map[string]string{
						"Authorization": "Bearer " + validAuthToken(),
					}
					return headers
				},
				RequestPathParams: func(data harness.Data) map[string]string {
					params := map[string]string{
						":dungeon_id": "17c19414-2d15-4d20-8fc3-36fc10341dc8",
					}
					return params
				},
				RequestBody: func(data harness.Data) interface{} {
					return nil
				},
				ResponseBody: testCaseResponseBody,
				ResponseCode: http.StatusNotFound,
			},
		},
	}

	for _, testCase := range testCases {

		t.Logf("Running test >%s<", testCase.Name)

		testFunc := func(method string, body interface{}) {

			if testCase.TestResponseCode() != http.StatusOK {
				return
			}

			var responseBody *schema.DungeonCharacterResponse
			if body != nil {
				responseBody = body.(*schema.DungeonCharacterResponse)
			}

			// Validate response body
			if testCase.expectResponseBody != nil {
				require.NotNil(t, responseBody, "Response body is not nil")
				require.GreaterOrEqual(t, len(responseBody.Data), 0, "Response body data ")

				expectResponseBody := testCase.expectResponseBody(th.Data)

				// Validate response body data
				for idx, expectData := range expectResponseBody.Data {
					require.NotNil(t, responseBody.Data[idx], "Response body index is not empty")

					// Validate character
					t.Logf("Checking character name >%s< >%s<", expectData.Name, responseBody.Data[idx].Name)
					require.Equal(t, expectData.Name, responseBody.Data[idx].Name)
					t.Logf("Checking character strength >%d< >%d<", expectData.Strength, responseBody.Data[idx].Strength)
					require.Equal(t, expectData.Strength, responseBody.Data[idx].Strength)
					t.Logf("Checking character dexterity >%d< >%d<", expectData.Dexterity, responseBody.Data[idx].Dexterity)
					require.Equal(t, expectData.Dexterity, responseBody.Data[idx].Dexterity)
					t.Logf("Checking character intelligence >%d< >%d<", expectData.Intelligence, responseBody.Data[idx].Intelligence)
					require.Equal(t, expectData.Intelligence, responseBody.Data[idx].Intelligence)

					t.Logf("Checking character current strength >%d< >%d<", expectData.CurrentStrength, responseBody.Data[idx].CurrentStrength)
					require.Equal(t, expectData.CurrentStrength, responseBody.Data[idx].CurrentStrength)
					t.Logf("Checking character current dexterity >%d< >%d<", expectData.CurrentDexterity, responseBody.Data[idx].CurrentDexterity)
					require.Equal(t, expectData.CurrentDexterity, responseBody.Data[idx].CurrentDexterity)
					t.Logf("Checking character current intelligence >%d< >%d<", expectData.CurrentIntelligence, responseBody.Data[idx].CurrentIntelligence)
					require.Equal(t, expectData.CurrentIntelligence, responseBody.Data[idx].CurrentIntelligence)

					t.Logf("Checking character health >%d< >%d<", expectData.Health, responseBody.Data[idx].Health)
					require.Equal(t, expectData.Health, responseBody.Data[idx].Health)
					t.Logf("Checking character fatigue >%d< >%d<", expectData.Fatigue, responseBody.Data[idx].Fatigue)
					require.Equal(t, expectData.Fatigue, responseBody.Data[idx].Fatigue)

					t.Logf("Checking character current health >%d< >%d<", expectData.CurrentHealth, responseBody.Data[idx].CurrentHealth)
					require.Equal(t, expectData.CurrentHealth, responseBody.Data[idx].CurrentHealth)
					t.Logf("Checking character current fatigue >%d< >%d<", expectData.CurrentFatigue, responseBody.Data[idx].CurrentFatigue)
					require.Equal(t, expectData.CurrentFatigue, responseBody.Data[idx].CurrentFatigue)
				}
			}

			// Check dates and add teardown ID's
			for _, data := range responseBody.Data {
				require.False(t, data.CreatedAt.IsZero(), "CreatedAt is not zero")
				if method == http.MethodPost {
					require.True(t, data.UpdatedAt.IsZero(), "UpdatedAt is zero")
					t.Logf("Adding dungeon character ID >%s< for teardown", data.ID)
					th.AddDungeonCharacterTeardownID(data.ID)
				}
			}
		}

		RunTestCase(t, th, &testCase, testFunc)
	}
}
