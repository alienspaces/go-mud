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
		}
		return headers
	}

	testCaseResponseBody := func(body io.Reader) (interface{}, error) {
		var responseBody *schema.DungeonCharacterResponse
		err = json.NewDecoder(body).Decode(&responseBody)
		return responseBody, err
	}

	testCases := []testCase{
		// TODO: Implement headers that control application behvaiour. In this case we want the
		// database transaction rolled back at the end of the request so we don't have to deal
		// with managing adding teardown ID's to the test harness.
		{
			TestCase: TestCase{
				Name:           "POST - Enter existing",
				HandlerConfig:  testCaseHandlerConfig,
				RequestHeaders: testCaseRequestHeaders,
				RequestPathParams: func(data harness.Data) map[string]string {
					dRec, _ := data.GetDungeonRecByName("Cave")
					cRec, _ := data.GetCharacterRecByName("Legislate")
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

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Logf("Running test >%s<", testCase.Name)

			// TODO: Re-implement

			// testFunc := func(method string, body interface{}) {
			// 	if testCase.TestResponseCode() != http.StatusOK {
			// 		return
			// 	}
			// 	var responseBody *schema.DungeonCharacterResponse
			// 	if body != nil {
			// 		responseBody = body.(*schema.DungeonCharacterResponse)
			// 	}
			// 	// Check dates and add teardown ID's
			// 	for _, data := range responseBody.Data {
			// 		require.False(t, data.CreatedAt.IsZero(), "CreatedAt is not zero")
			// 	}
			// }
			// RunTestCase(t, th, &testCase, testFunc)
		})
	}
}
