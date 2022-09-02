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

func TestPostCharacterHandler(t *testing.T) {

	// Test harness
	th, err := NewTestHarness()
	require.NoError(t, err, "New test data returns without error")

	type testCase struct {
		TestCase
		expectResponseBody func(data harness.Data) *schema.CharacterResponse
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

	testCaseRequestHeaders := func(data harness.Data) map[string]string {
		headers := map[string]string{
			"Authorization": "Bearer " + validAuthToken(),
		}
		return headers
	}

	testCaseResponseBody := func(body io.Reader) (interface{}, error) {
		var responseBody *schema.CharacterResponse
		err = json.NewDecoder(body).Decode(&responseBody)
		return responseBody, err
	}

	testCases := []testCase{
		{
			TestCase: TestCase{
				Name: "POST - Create one with valid attributes",
				HandlerConfig: func(rnr *Runner) server.HandlerConfig {
					return rnr.HandlerConfig[postCharacter]
				},
				RequestHeaders: testCaseRequestHeaders,
				RequestPathParams: func(data harness.Data) map[string]string {
					return nil
				},
				RequestBody: func(data harness.Data) interface{} {
					res := schema.CharacterRequest{
						Data: schema.CharacterData{
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
		t.Run(testCase.Name, func(t *testing.T) {
			t.Logf("Running test >%s<", testCase.Name)

			testFunc := func(method string, body interface{}) {

				if testCase.TestResponseCode() != http.StatusOK {
					return
				}

				var responseBody *schema.CharacterResponse
				if body != nil {
					responseBody = body.(*schema.CharacterResponse)
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
						th.AddCharacterTeardownID(data.ID)
					}
				}
			}

			RunTestCase(t, th, &testCase, testFunc)
		})
	}
}

func TestGetCharacterHandler(t *testing.T) {

	// Test harness
	th, err := NewTestHarness()
	require.NoError(t, err, "New test data returns without error")

	type testCase struct {
		TestCase
		expectResponseBody func(data harness.Data) *schema.CharacterResponse
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

	testCaseRequestHeaders := func(data harness.Data) map[string]string {
		headers := map[string]string{
			"Authorization": "Bearer " + validAuthToken(),
		}
		return headers
	}

	testCaseResponseBody := func(body io.Reader) (interface{}, error) {
		var responseBody *schema.CharacterResponse
		err = json.NewDecoder(body).Decode(&responseBody)
		return responseBody, err
	}

	testCases := []testCase{
		{
			TestCase: TestCase{
				Name: "GET - Get many",
				HandlerConfig: func(rnr *Runner) server.HandlerConfig {
					return rnr.HandlerConfig[getCharacters]
				},
				RequestHeaders: testCaseRequestHeaders,
				RequestPathParams: func(data harness.Data) map[string]string {
					params := map[string]string{}
					return params
				},
				RequestBody: func(data harness.Data) interface{} {
					return nil
				},
				ResponseBody: testCaseResponseBody,
				ResponseCode: http.StatusOK,
			},
			expectResponseBody: func(data harness.Data) *schema.CharacterResponse {
				res := schema.CharacterResponse{
					Data: []schema.CharacterData{
						{
							ID:               data.CharacterRecs[0].ID,
							Name:             data.CharacterRecs[0].Name,
							Strength:         data.CharacterRecs[0].Strength,
							Dexterity:        data.CharacterRecs[0].Dexterity,
							Intelligence:     data.CharacterRecs[0].Intelligence,
							Health:           data.CharacterRecs[0].Health,
							Fatigue:          data.CharacterRecs[0].Fatigue,
							Coins:            data.CharacterRecs[0].Coins,
							ExperiencePoints: data.CharacterRecs[0].ExperiencePoints,
							AttributePoints:  data.CharacterRecs[0].AttributePoints,
						},
					},
				}
				return &res
			},
		},
		{
			TestCase: TestCase{
				Name: "GET - Get one",
				HandlerConfig: func(rnr *Runner) server.HandlerConfig {
					return rnr.HandlerConfig[getCharacter]
				},
				RequestHeaders: func(data harness.Data) map[string]string {
					headers := map[string]string{
						"Authorization": "Bearer " + validAuthToken(),
					}
					return headers
				},
				RequestPathParams: func(data harness.Data) map[string]string {
					params := map[string]string{
						":character_id": data.CharacterRecs[0].ID,
					}
					return params
				},
				RequestBody: func(data harness.Data) interface{} {
					return nil
				},
				ResponseBody: testCaseResponseBody,
				ResponseCode: http.StatusOK,
			},
			expectResponseBody: func(data harness.Data) *schema.CharacterResponse {
				res := schema.CharacterResponse{
					Data: []schema.CharacterData{
						{
							ID:               data.CharacterRecs[0].ID,
							Name:             data.CharacterRecs[0].Name,
							Strength:         data.CharacterRecs[0].Strength,
							Dexterity:        data.CharacterRecs[0].Dexterity,
							Intelligence:     data.CharacterRecs[0].Intelligence,
							Health:           data.CharacterRecs[0].Health,
							Fatigue:          data.CharacterRecs[0].Fatigue,
							Coins:            data.CharacterRecs[0].Coins,
							ExperiencePoints: data.CharacterRecs[0].ExperiencePoints,
							AttributePoints:  data.CharacterRecs[0].AttributePoints,
						},
					},
				}
				return &res
			},
		},
		{
			TestCase: TestCase{
				Name: "GET - Get one with incorrect character ID",
				HandlerConfig: func(rnr *Runner) server.HandlerConfig {
					return rnr.HandlerConfig[getCharacter]
				},
				RequestHeaders: func(data harness.Data) map[string]string {
					headers := map[string]string{
						"Authorization": "Bearer " + validAuthToken(),
					}
					return headers
				},
				RequestPathParams: func(data harness.Data) map[string]string {
					params := map[string]string{
						":character_id": "a08eb991-759d-4671-8698-9f26056717e2",
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
					return rnr.HandlerConfig[getCharacter]
				},
				RequestHeaders: func(data harness.Data) map[string]string {
					headers := map[string]string{
						"Authorization": "Bearer " + validAuthToken(),
					}
					return headers
				},
				RequestPathParams: func(data harness.Data) map[string]string {
					params := map[string]string{
						":character_id": "notauuid-759d-4671-8698-notauuid",
					}
					return params
				},
				RequestBody: func(data harness.Data) interface{} {
					return nil
				},
				ResponseBody: testCaseResponseBody,
				ResponseCode: http.StatusBadRequest,
			},
		},
	}

	for _, testCase := range testCases {

		t.Logf("Running test >%s<", testCase.Name)

		testFunc := func(method string, body interface{}) {

			if testCase.TestResponseCode() != http.StatusOK {
				return
			}

			var responseBody *schema.CharacterResponse
			if body != nil {
				responseBody = body.(*schema.CharacterResponse)
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

					t.Logf("Checking character health >%d< >%d<", expectData.Health, responseBody.Data[idx].Health)
					require.Equal(t, expectData.Health, responseBody.Data[idx].Health)
					t.Logf("Checking character fatigue >%d< >%d<", expectData.Fatigue, responseBody.Data[idx].Fatigue)
					require.Equal(t, expectData.Fatigue, responseBody.Data[idx].Fatigue)

					t.Logf("Checking character coins >%d< >%d<", expectData.Coins, responseBody.Data[idx].Coins)
					require.Equal(t, expectData.Coins, responseBody.Data[idx].Coins)
					t.Logf("Checking character experience points >%d< >%d<", expectData.ExperiencePoints, responseBody.Data[idx].ExperiencePoints)
					require.Equal(t, expectData.ExperiencePoints, responseBody.Data[idx].ExperiencePoints)
					t.Logf("Checking character attribute points >%d< >%d<", expectData.AttributePoints, responseBody.Data[idx].AttributePoints)
					require.Equal(t, expectData.AttributePoints, responseBody.Data[idx].AttributePoints)
				}
			}

			// Check dates and add teardown ID's
			for _, data := range responseBody.Data {
				require.False(t, data.CreatedAt.IsZero(), "CreatedAt is not zero")
				if method == http.MethodPost {
					require.True(t, data.UpdatedAt.IsZero(), "UpdatedAt is zero")
					t.Logf("Adding dungeon character ID >%s< for teardown", data.ID)
					th.AddCharacterTeardownID(data.ID)
				}
			}
		}

		RunTestCase(t, th, &testCase, testFunc)
	}
}
