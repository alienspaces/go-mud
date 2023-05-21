package runner

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/go-mud/backend/core/auth"
	"gitlab.com/alienspaces/go-mud/backend/core/server"
	schema "gitlab.com/alienspaces/go-mud/backend/schema/game"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/harness"
)

func TestPostDungeonCharacterEnterHandler(t *testing.T) {

	// Test harness
	th, err := NewTestHarness()
	require.NoError(t, err, "New test data returns without error")

	type testCase struct {
		TestCase
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
		return rnr.HandlerConfig[postDungeonCharacterEnter]
	}

	testCaseRequestHeaders := func(data harness.Data) map[string]string {
		headers := map[string]string{
			"Authorization": "Bearer " + validAuthToken(),
			"X-Tx-Rollback": "true",
		}
		return headers
	}

	testCaseResponseBody := func(body io.Reader) (interface{}, error) {
		var responseBody *schema.DungeonCharacterResponse
		err = json.NewDecoder(body).Decode(&responseBody)
		return responseBody, err
	}

	testCases := []testCase{
		{
			TestCase: TestCase{
				Name:           "POST - Enter existing",
				HandlerConfig:  testCaseHandlerConfig,
				RequestHeaders: testCaseRequestHeaders,
				RequestPathParams: func(data harness.Data) map[string]string {
					dRec, _ := data.GetDungeonRecByName(harness.DungeonNameCave)
					cRec, _ := data.GetCharacterRecByName(harness.CharacterNameBolster)
					params := map[string]string{
						":dungeon_id":   dRec.ID,
						":character_id": cRec.ID,
					}
					return params
				},
				ResponseBody: testCaseResponseBody,
				ResponseCode: http.StatusOK,
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

				var responseBody *schema.DungeonCharacterResponse
				if body != nil {
					responseBody = body.(*schema.DungeonCharacterResponse)
				}

				for _, data := range responseBody.Data {
					require.NotEmpty(t, data.DungeonID, "Data DungeonID is not empty")
					require.NotEmpty(t, data.DungeonName, "Data DungeonName is not empty")
					require.NotEmpty(t, data.LocationID, "Data LocationID is not empty")
					require.NotEmpty(t, data.LocationName, "Data LocationName is not empty")

					require.NotEmpty(t, data.CharacterID, "Data CharacterID is not empty")
					require.NotEmpty(t, data.CharacterName, "Data CharacterName is not empty")
					require.NotEmpty(t, data.CharacterStrength, "Data CharacterStrength is not empty")
					require.NotEmpty(t, data.CharacterDexterity, "Data CharacterDexterity is not empty")
					require.NotEmpty(t, data.CharacterIntelligence, "Data CharacterIntelligence is not empty")
					require.NotEmpty(t, data.CharacterCurrentStrength, "Data CharacterCurrentStrength is not empty")
					require.NotEmpty(t, data.CharacterCurrentDexterity, "Data CharacterCurrentDexterity is not empty")
					require.NotEmpty(t, data.CharacterCurrentIntelligence, "Data CharacterCurrentIntelligence is not empty")
					require.NotEmpty(t, data.CharacterHealth, "Data CharacterHealth is not empty")
					require.NotEmpty(t, data.CharacterFatigue, "Data CharacterFatigue is not empty")
					require.NotEmpty(t, data.CharacterCurrentHealth, "Data CharacterCurrentHealth is not empty")
					require.NotEmpty(t, data.CharacterCurrentFatigue, "Data CharacterCurrentFatigue is not empty")

					require.False(t, data.CharacterCreatedAt.IsZero(), "Data CreatedAt is not zero")
				}
			}
			RunTestCase(t, th, &tc, testFunc)
		})
	}
}

func TestPostDungeonCharacterExitHandler(t *testing.T) {

	// Test harness
	th, err := NewTestHarness()
	require.NoError(t, err, "New test data returns without error")

	type testCase struct {
		TestCase
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
		return rnr.HandlerConfig[postDungeonCharacterExit]
	}

	testCaseRequestHeaders := func(data harness.Data) map[string]string {
		headers := map[string]string{
			"Authorization": "Bearer " + validAuthToken(),
			"X-Tx-Rollback": "true",
		}
		return headers
	}

	testCaseResponseBody := func(body io.Reader) (interface{}, error) {
		var responseBody *schema.DungeonCharacterResponse
		err = json.NewDecoder(body).Decode(&responseBody)
		return responseBody, err
	}

	testCases := []testCase{
		{
			TestCase: TestCase{
				Name:           "POST - Exit existing",
				HandlerConfig:  testCaseHandlerConfig,
				RequestHeaders: testCaseRequestHeaders,
				RequestPathParams: func(data harness.Data) map[string]string {
					dRec, _ := data.GetDungeonRecByName(harness.DungeonNameCave)
					cRec, _ := data.GetCharacterRecByName(harness.CharacterNameBarricade)
					params := map[string]string{
						":dungeon_id":   dRec.ID,
						":character_id": cRec.ID,
					}
					return params
				},
				ResponseBody: testCaseResponseBody,
				ResponseCode: http.StatusOK,
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

				var responseBody *schema.DungeonCharacterResponse
				if body != nil {
					responseBody = body.(*schema.DungeonCharacterResponse)
				}

				for _, data := range responseBody.Data {
					require.Empty(t, data.DungeonID, "Data DungeonID is empty")
					require.Empty(t, data.LocationID, "Data LocationID is empty")

					require.NotEmpty(t, data.CharacterID, "Data CharacterID is not empty")
					require.NotEmpty(t, data.CharacterName, "Data CharacterName is not empty")
					require.NotEmpty(t, data.CharacterStrength, "Data CharacterStrength is not empty")
					require.NotEmpty(t, data.CharacterDexterity, "Data CharacterDexterity is not empty")
					require.NotEmpty(t, data.CharacterIntelligence, "Data CharacterIntelligence is not empty")
					require.NotEmpty(t, data.CharacterHealth, "Data CharacterHealth is not empty")
					require.NotEmpty(t, data.CharacterFatigue, "Data CharacterFatigue is not empty")

					require.False(t, data.CharacterCreatedAt.IsZero(), "Data CreatedAt is not zero")
					require.False(t, data.CharacterUpdatedAt.IsZero(), "Data UpdatedAt is not zero")
				}
			}
			RunTestCase(t, th, &tc, testFunc)
		})
	}
}
