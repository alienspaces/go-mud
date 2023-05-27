package runner

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/go-mud/backend/core/server"
	schema "gitlab.com/alienspaces/go-mud/backend/schema/game"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/harness"
)

func TestPostCharacterHandler(t *testing.T) {

	// Test harness
	th, err := NewTestHarness()
	require.NoError(t, err, "New test data returns without error")

	type testCase struct {
		TestCase
		expectResponseBody func(data harness.Data) *schema.CharacterResponse
	}

	testCaseRequestHeaders := func(data harness.Data) map[string]string {
		headers := map[string]string{
			"X-Tx-Rollback": "true",
		}
		return headers
	}

	testCaseResponseDecoder := func(body io.Reader) (interface{}, error) {
		var responseBody *schema.CharacterResponse
		err = json.NewDecoder(body).Decode(&responseBody)
		return responseBody, err
	}

	testCases := []testCase{
		{
			TestCase: TestCase{
				Name: "create one with valid attributes",
				HandlerConfig: func(rnr *Runner) server.HandlerConfig {
					return rnr.HandlerConfig[postCharacter]
				},
				RequestHeaders: testCaseRequestHeaders,
				RequestPathParams: func(data harness.Data) map[string]string {
					return nil
				},
				RequestBody: func(data harness.Data) interface{} {
					res := schema.CharacterRequest{
						Data: schema.DungeonCharacterData{
							CharacterName:         gofakeit.Name() + gofakeit.Name(),
							CharacterStrength:     10,
							CharacterDexterity:    10,
							CharacterIntelligence: 10,
						},
					}
					return &res
				},
				ResponseDecoder: testCaseResponseDecoder,
				ResponseCode:    http.StatusOK,
			},
		},
	}

	for _, testCase := range testCases {

		t.Logf("Running test >%s<", testCase.Name)

		t.Run(testCase.Name, func(t *testing.T) {

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
					require.False(t, data.CharacterCreatedAt.IsZero(), "CreatedAt is not zero")
					if method == http.MethodPost {
						require.True(t, data.CharacterUpdatedAt.IsZero(), "UpdatedAt is zero")
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

	testCaseRequestHeaders := func(data harness.Data) map[string]string {
		headers := map[string]string{}
		return headers
	}

	testCaseResponseDecoder := func(body io.Reader) (interface{}, error) {
		var responseBody *schema.CharacterResponse
		err = json.NewDecoder(body).Decode(&responseBody)
		return responseBody, err
	}

	testCases := []testCase{
		{
			TestCase: TestCase{
				Name: "get many",
				HandlerConfig: func(rnr *Runner) server.HandlerConfig {
					return rnr.HandlerConfig[getCharacters]
				},
				RequestHeaders: testCaseRequestHeaders,
				RequestQueryParams: func(data harness.Data) map[string]string {
					params := map[string]string{
						"character_id": data.CharacterRecs[0].ID,
					}
					return params
				},
				RequestBody: func(data harness.Data) interface{} {
					return nil
				},
				ResponseDecoder: testCaseResponseDecoder,
				ResponseCode:    http.StatusOK,
			},
			expectResponseBody: func(data harness.Data) *schema.CharacterResponse {
				res := schema.CharacterResponse{
					Data: []schema.DungeonCharacterData{
						{
							CharacterID:               data.CharacterRecs[2].ID,
							CharacterName:             data.CharacterRecs[2].Name,
							CharacterStrength:         data.CharacterRecs[2].Strength,
							CharacterDexterity:        data.CharacterRecs[2].Dexterity,
							CharacterIntelligence:     data.CharacterRecs[2].Intelligence,
							CharacterHealth:           data.CharacterRecs[2].Health,
							CharacterFatigue:          data.CharacterRecs[2].Fatigue,
							CharacterCoins:            data.CharacterRecs[2].Coins,
							CharacterExperiencePoints: data.CharacterRecs[2].ExperiencePoints,
							CharacterAttributePoints:  data.CharacterRecs[2].AttributePoints,
						},
						{
							CharacterID:               data.CharacterRecs[1].ID,
							CharacterName:             data.CharacterRecs[1].Name,
							CharacterStrength:         data.CharacterRecs[1].Strength,
							CharacterDexterity:        data.CharacterRecs[1].Dexterity,
							CharacterIntelligence:     data.CharacterRecs[1].Intelligence,
							CharacterHealth:           data.CharacterRecs[1].Health,
							CharacterFatigue:          data.CharacterRecs[1].Fatigue,
							CharacterCoins:            data.CharacterRecs[1].Coins,
							CharacterExperiencePoints: data.CharacterRecs[1].ExperiencePoints,
							CharacterAttributePoints:  data.CharacterRecs[1].AttributePoints,
						},
						{
							CharacterID:               data.CharacterRecs[0].ID,
							CharacterName:             data.CharacterRecs[0].Name,
							CharacterStrength:         data.CharacterRecs[0].Strength,
							CharacterDexterity:        data.CharacterRecs[0].Dexterity,
							CharacterIntelligence:     data.CharacterRecs[0].Intelligence,
							CharacterHealth:           data.CharacterRecs[0].Health,
							CharacterFatigue:          data.CharacterRecs[0].Fatigue,
							CharacterCoins:            data.CharacterRecs[0].Coins,
							CharacterExperiencePoints: data.CharacterRecs[0].ExperiencePoints,
							CharacterAttributePoints:  data.CharacterRecs[0].AttributePoints,
						},
					},
				}
				return &res
			},
		},
		{
			TestCase: TestCase{
				Name: "get one",
				HandlerConfig: func(rnr *Runner) server.HandlerConfig {
					return rnr.HandlerConfig[getCharacter]
				},
				RequestHeaders: func(data harness.Data) map[string]string {
					headers := map[string]string{}
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
				ResponseDecoder: testCaseResponseDecoder,
				ResponseCode:    http.StatusOK,
			},
			expectResponseBody: func(data harness.Data) *schema.CharacterResponse {
				res := schema.CharacterResponse{
					Data: []schema.DungeonCharacterData{
						{
							CharacterID:               data.CharacterRecs[0].ID,
							CharacterName:             data.CharacterRecs[0].Name,
							CharacterStrength:         data.CharacterRecs[0].Strength,
							CharacterDexterity:        data.CharacterRecs[0].Dexterity,
							CharacterIntelligence:     data.CharacterRecs[0].Intelligence,
							CharacterHealth:           data.CharacterRecs[0].Health,
							CharacterFatigue:          data.CharacterRecs[0].Fatigue,
							CharacterCoins:            data.CharacterRecs[0].Coins,
							CharacterExperiencePoints: data.CharacterRecs[0].ExperiencePoints,
							CharacterAttributePoints:  data.CharacterRecs[0].AttributePoints,
						},
					},
				}
				return &res
			},
		},
		{
			TestCase: TestCase{
				Name: "Get one with unknown character id",
				HandlerConfig: func(rnr *Runner) server.HandlerConfig {
					return rnr.HandlerConfig[getCharacter]
				},
				RequestHeaders: func(data harness.Data) map[string]string {
					headers := map[string]string{}
					return headers
				},
				RequestPathParams: func(data harness.Data) map[string]string {
					params := map[string]string{
						":character_id": "a08eb991-759d-4671-8698-9f26056717e2",
					}
					return params
				},
				ResponseDecoder: testCaseResponseDecoder,
				ResponseCode:    http.StatusNotFound,
			},
		},
		{
			TestCase: TestCase{
				Name: "get one with invalid character id",
				HandlerConfig: func(rnr *Runner) server.HandlerConfig {
					return rnr.HandlerConfig[getCharacter]
				},
				RequestHeaders: func(data harness.Data) map[string]string {
					headers := map[string]string{}
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
				ResponseDecoder: testCaseResponseDecoder,
				ResponseCode:    http.StatusBadRequest,
			},
		},
	}

	for _, testCase := range testCases {

		t.Logf("Running test >%s<", testCase.Name)

		t.Run(testCase.Name, func(t *testing.T) {

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

					if expectResponseBody != nil {
						require.NotNil(t, responseBody.Data, "Response body is not nil")
						require.Equal(t, len(expectResponseBody.Data), len(responseBody.Data), "Response body length equals expected")
					}

					// Validate response body data
					for idx, expectData := range expectResponseBody.Data {
						require.NotNil(t, responseBody.Data[idx], "Response body index is not empty")

						// Validate character
						t.Logf("Checking character name >%s< >%s<", expectData.CharacterName, responseBody.Data[idx].CharacterName)
						require.Equal(t, expectData.CharacterName, responseBody.Data[idx].CharacterName, "Character name equals expected")
						t.Logf("Checking character strength >%d< >%d<", expectData.CharacterStrength, responseBody.Data[idx].CharacterStrength)
						require.Equal(t, expectData.CharacterStrength, responseBody.Data[idx].CharacterStrength, "Character strength equals expected")
						t.Logf("Checking character dexterity >%d< >%d<", expectData.CharacterDexterity, responseBody.Data[idx].CharacterDexterity)
						require.Equal(t, expectData.CharacterDexterity, responseBody.Data[idx].CharacterDexterity, "Character dexterity equals expected")
						t.Logf("Checking character intelligence >%d< >%d<", expectData.CharacterIntelligence, responseBody.Data[idx].CharacterIntelligence)
						require.Equal(t, expectData.CharacterIntelligence, responseBody.Data[idx].CharacterIntelligence, "Character intelligence equals expected")

						t.Logf("Checking character health >%d< >%d<", expectData.CharacterHealth, responseBody.Data[idx].CharacterHealth)
						require.Equal(t, expectData.CharacterHealth, responseBody.Data[idx].CharacterHealth, "Character health equals expected")
						t.Logf("Checking character fatigue >%d< >%d<", expectData.CharacterFatigue, responseBody.Data[idx].CharacterFatigue)
						require.Equal(t, expectData.CharacterFatigue, responseBody.Data[idx].CharacterFatigue, "Character fatigue equals expected")

						t.Logf("Checking character coins >%d< >%d<", expectData.CharacterCoins, responseBody.Data[idx].CharacterCoins)
						require.Equal(t, expectData.CharacterCoins, responseBody.Data[idx].CharacterCoins, "Character coins equals expected")
						t.Logf("Checking character experience points >%d< >%d<", expectData.CharacterExperiencePoints, responseBody.Data[idx].CharacterExperiencePoints)
						require.Equal(t, expectData.CharacterExperiencePoints, responseBody.Data[idx].CharacterExperiencePoints, "Character experience points equals expected")
						t.Logf("Checking character attribute points >%d< >%d<", expectData.CharacterAttributePoints, responseBody.Data[idx].CharacterAttributePoints)
						require.Equal(t, expectData.CharacterAttributePoints, responseBody.Data[idx].CharacterAttributePoints, "Character attribute points equals expected")
					}
				}

				// Check dates and add teardown ID's
				for _, data := range responseBody.Data {
					require.False(t, data.CharacterCreatedAt.IsZero(), "CreatedAt is not zero")
					if method == http.MethodPost {
						require.True(t, data.CharacterUpdatedAt.IsZero(), "UpdatedAt is zero")
					}
				}
			}

			RunTestCase(t, th, &testCase, testFunc)
		})
	}
}
