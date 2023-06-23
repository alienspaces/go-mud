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
				RequestPathParams: func(data harness.Data) map[string]string {
					return nil
				},
				RequestBody: func(data harness.Data) interface{} {
					res := schema.CharacterRequest{
						Data: schema.DungeonCharacterData{
							Name:         gofakeit.Name() + gofakeit.Name(),
							Strength:     10,
							Dexterity:    10,
							Intelligence: 10,
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
					require.False(t, data.CreatedAt.IsZero(), "CreatedAt is not zero")
					if method == http.MethodPost {
						require.True(t, data.UpdatedAt.IsZero(), "UpdatedAt is zero")
					}
				}
			}

			RunTestCase(t, th, &testCase, testFunc)
		})
	}
}

func TestGetCharacterHandler(t *testing.T) {

	th, err := NewTestHarness()
	require.NoError(t, err, "New test data returns without error")

	_, err = th.Setup()
	require.NoError(t, err, "Test data setup returns without error")
	defer func() {
		err = th.Teardown()
		require.NoError(t, err, "Test data teardown returns without error")
	}()

	type testCase struct {
		TestCase
		expectResponseBody func(data harness.Data) *schema.CharacterResponse
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
				RequestQueryParams: func(data harness.Data) map[string]interface{} {
					params := map[string]interface{}{
						"name": data.CharacterRecs[0].Name,
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
						t.Logf("Checking character name >%s< >%s<", expectData.Name, responseBody.Data[idx].Name)
						require.Equal(t, expectData.Name, responseBody.Data[idx].Name, "Character name equals expected")
						t.Logf("Checking character strength >%d< >%d<", expectData.Strength, responseBody.Data[idx].Strength)
						require.Equal(t, expectData.Strength, responseBody.Data[idx].Strength, "Character strength equals expected")
						t.Logf("Checking character dexterity >%d< >%d<", expectData.Dexterity, responseBody.Data[idx].Dexterity)
						require.Equal(t, expectData.Dexterity, responseBody.Data[idx].Dexterity, "Character dexterity equals expected")
						t.Logf("Checking character intelligence >%d< >%d<", expectData.Intelligence, responseBody.Data[idx].Intelligence)
						require.Equal(t, expectData.Intelligence, responseBody.Data[idx].Intelligence, "Character intelligence equals expected")

						t.Logf("Checking character health >%d< >%d<", expectData.Health, responseBody.Data[idx].Health)
						require.Equal(t, expectData.Health, responseBody.Data[idx].Health, "Character health equals expected")
						t.Logf("Checking character fatigue >%d< >%d<", expectData.Fatigue, responseBody.Data[idx].Fatigue)
						require.Equal(t, expectData.Fatigue, responseBody.Data[idx].Fatigue, "Character fatigue equals expected")

						t.Logf("Checking character coins >%d< >%d<", expectData.Coins, responseBody.Data[idx].Coins)
						require.Equal(t, expectData.Coins, responseBody.Data[idx].Coins, "Character coins equals expected")
						t.Logf("Checking character experience points >%d< >%d<", expectData.ExperiencePoints, responseBody.Data[idx].ExperiencePoints)
						require.Equal(t, expectData.ExperiencePoints, responseBody.Data[idx].ExperiencePoints, "Character experience points equals expected")
						t.Logf("Checking character attribute points >%d< >%d<", expectData.AttributePoints, responseBody.Data[idx].AttributePoints)
						require.Equal(t, expectData.AttributePoints, responseBody.Data[idx].AttributePoints, "Character attribute points equals expected")
					}
				}

				// Check dates and add teardown ID's
				for _, data := range responseBody.Data {
					require.False(t, data.CreatedAt.IsZero(), "CreatedAt is not zero")
					if method == http.MethodPost {
						require.True(t, data.UpdatedAt.IsZero(), "UpdatedAt is zero")
					}
				}
			}

			RunTestCase(t, th, &testCase, testFunc)
		})
	}
}
