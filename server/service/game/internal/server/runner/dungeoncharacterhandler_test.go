package runner

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/brianvoe/gofakeit"

	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/go-mud/server/core/auth"
	coreschema "gitlab.com/alienspaces/go-mud/server/core/schema"
	"gitlab.com/alienspaces/go-mud/server/core/server"
	"gitlab.com/alienspaces/go-mud/server/schema"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/harness"
)

func TestDungeonCharacterHandler(t *testing.T) {

	// Test harness
	th, err := NewTestHarness()
	require.NoError(t, err, "New test data returns without error")

	type testCase struct {
		TestCase
		responseBody func(data harness.Data) *schema.DungeonCharacterResponse
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
				ResponseCode: http.StatusOK,
			},
			responseBody: func(data harness.Data) *schema.DungeonCharacterResponse {
				res := schema.DungeonCharacterResponse{
					Data: []schema.DungeonCharacterData{
						{
							ID: data.DungeonCharacterRecs[0].ID,
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
				ResponseCode: http.StatusNotFound,
			},
			responseBody: func(data harness.Data) *schema.DungeonCharacterResponse {
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
				ResponseCode: http.StatusOK,
			},
			responseBody: func(data harness.Data) *schema.DungeonCharacterResponse {
				res := schema.DungeonCharacterResponse{
					Data: []schema.DungeonCharacterData{
						{
							ID: data.DungeonCharacterRecs[0].ID,
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
				ResponseCode: http.StatusNotFound,
			},
		},
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
				ResponseCode: http.StatusOK,
			},
		},
	}

	v, err := coreschema.NewValidator(th.Config, th.Log)
	require.NoError(t, err, "Validator returns without error")

	for _, testCase := range testCases {

		t.Logf("Running test >%s<", testCase.Name)

		testFunc := func(method string, body io.Reader) {

			if testCase.TestResponseCode() != http.StatusOK {
				return
			}

			var responseBody *schema.DungeonCharacterResponse
			err = json.NewDecoder(body).Decode(&responseBody)
			require.NoError(t, err, "Decode returns without error")

			if testCase.responseBody != nil {

				// Validate response body
				require.NotNil(t, responseBody, "Response body is not nil")

				jsonData, err := json.Marshal(responseBody)
				require.NoError(t, err, "Marshal returns without error")

				err = v.Validate(coreschema.Config{
					Key:      "dungeoncharacter.response",
					Location: "dungeoncharacter",
					Main:     "main.schema.json",
					References: []string{
						"data.schema.json",
					},
				}, string(jsonData))
				require.NoError(t, err, "Validates against schema without error")
			}

			require.NotNil(t, responseBody, "Response body is not nil")
			require.GreaterOrEqual(t, len(responseBody.Data), 0, "Response body data ")

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
