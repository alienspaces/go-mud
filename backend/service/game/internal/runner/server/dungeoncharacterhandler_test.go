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

func TestPostDungeonCharacterEnterHandler(t *testing.T) {

	// Test harness
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
	}

	testCaseHandlerConfig := func(rnr *Runner) server.HandlerConfig {
		return rnr.HandlerConfig[postDungeonCharacterEnter]
	}

	testCaseResponseDecoder := func(body io.Reader) (interface{}, error) {
		var responseBody *schema.DungeonCharacterResponse
		err = json.NewDecoder(body).Decode(&responseBody)
		return responseBody, err
	}

	testCases := []testCase{
		{
			TestCase: TestCase{
				Name:          "POST - Enter existing",
				HandlerConfig: testCaseHandlerConfig,
				RequestPathParams: func(data harness.Data) map[string]string {
					dRec, _ := data.GetDungeonRecByName(harness.DungeonNameCave)
					cRec, _ := data.GetCharacterRecByName(harness.CharacterNameBolster)
					params := map[string]string{
						":dungeon_id":   dRec.ID,
						":character_id": cRec.ID,
					}
					return params
				},
				ResponseDecoder: testCaseResponseDecoder,
				ResponseCode:    http.StatusOK,
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
					require.NotEmpty(t, data.ID, "Data ID is not empty")
					require.NotEmpty(t, data.Name, "Data Name is not empty")
					require.NotEmpty(t, data.Strength, "Data Strength is not empty")
					require.NotEmpty(t, data.Dexterity, "Data Dexterity is not empty")
					require.NotEmpty(t, data.Intelligence, "Data Intelligence is not empty")
					require.NotEmpty(t, data.CurrentStrength, "Data CurrentStrength is not empty")
					require.NotEmpty(t, data.CurrentDexterity, "Data CurrentDexterity is not empty")
					require.NotEmpty(t, data.CurrentIntelligence, "Data CurrentIntelligence is not empty")
					require.NotEmpty(t, data.Health, "Data Health is not empty")
					require.NotEmpty(t, data.Fatigue, "Data Fatigue is not empty")
					require.NotEmpty(t, data.CurrentHealth, "Data CurrentHealth is not empty")
					require.NotEmpty(t, data.CurrentFatigue, "Data CurrentFatigue is not empty")

					require.NotEmpty(t, data.Dungeon, "Data Dungeon is not empty")
					require.NotEmpty(t, data.Dungeon.ID, "Data Dungeon ID is not empty")
					require.NotEmpty(t, data.Dungeon.Name, "Data Dungeon Name is not empty")

					require.NotEmpty(t, data.Location, "Data Location is not empty")
					require.NotEmpty(t, data.Location.ID, "Data Location ID is not empty")
					require.NotEmpty(t, data.Location.Name, "Data Location Name is not empty")

					require.False(t, data.CreatedAt.IsZero(), "Data CreatedAt is not zero")
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

	testCaseHandlerConfig := func(rnr *Runner) server.HandlerConfig {
		return rnr.HandlerConfig[postDungeonCharacterExit]
	}

	testCaseResponseDecoder := func(body io.Reader) (interface{}, error) {
		var responseBody *schema.DungeonCharacterResponse
		err = json.NewDecoder(body).Decode(&responseBody)
		return responseBody, err
	}

	testCases := []testCase{
		{
			TestCase: TestCase{
				Name:          "POST - Exit existing",
				HandlerConfig: testCaseHandlerConfig,
				RequestPathParams: func(data harness.Data) map[string]string {
					dRec, _ := data.GetDungeonRecByName(harness.DungeonNameCave)
					cRec, _ := data.GetCharacterRecByName(harness.CharacterNameBarricade)
					params := map[string]string{
						":dungeon_id":   dRec.ID,
						":character_id": cRec.ID,
					}
					return params
				},
				ResponseDecoder: testCaseResponseDecoder,
				ResponseCode:    http.StatusOK,
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
					require.NotEmpty(t, data.ID, "Data ID is not empty")
					require.NotEmpty(t, data.Name, "Data Name is not empty")
					require.NotEmpty(t, data.Strength, "Data Strength is not empty")
					require.NotEmpty(t, data.Dexterity, "Data Dexterity is not empty")
					require.NotEmpty(t, data.Intelligence, "Data Intelligence is not empty")
					require.NotEmpty(t, data.Health, "Data Health is not empty")
					require.NotEmpty(t, data.Fatigue, "Data Fatigue is not empty")

					require.Empty(t, data.Dungeon, "Data Dungeon is empty")
					require.Empty(t, data.Location, "Data Location is empty")

					require.False(t, data.CreatedAt.IsZero(), "Data CreatedAt is not zero")
					require.False(t, data.UpdatedAt.IsZero(), "Data UpdatedAt is not zero")
				}
			}
			RunTestCase(t, th, &tc, testFunc)
		})
	}
}
