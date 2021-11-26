package runner

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"gitlab.com/alienspaces/go-mud/server/schema"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/require"
	"gitlab.com/alienspaces/go-mud/server/core/auth"
	"gitlab.com/alienspaces/go-mud/server/core/server"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/harness"
)

func TestDungeonCharacterActionHandler(t *testing.T) {

	// Test harness
	th, err := NewTestHarness()
	require.NoError(t, err, "New test data returns without error")

	type TestCase struct {
		name               string
		config             func(rnr *Runner) server.HandlerConfig
		requestHeaders     func(data harness.Data) map[string]string
		requestParams      func(data harness.Data) map[string]string
		queryParams        func(data harness.Data) map[string]string
		requestBody        func(data harness.Data) *schema.DungeonActionRequest
		expectResponseBody func(data harness.Data) *schema.DungeonActionResponse
		expectResponseCode int
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

	tests := []TestCase{
		{
			name: "POST - look",
			config: func(rnr *Runner) server.HandlerConfig {
				return rnr.HandlerConfig[6]
			},
			requestHeaders: func(data harness.Data) map[string]string {
				headers := map[string]string{
					"Authorization": "Bearer " + validAuthToken(),
				}
				return headers
			},
			requestParams: func(data harness.Data) map[string]string {
				params := map[string]string{
					":dungeon_id":   data.DungeonRecs[0].ID,
					":character_id": data.DungeonCharacterRecs[0].ID,
				}
				return params
			},
			requestBody: func(data harness.Data) *schema.DungeonActionRequest {
				res := schema.DungeonActionRequest{
					Data: schema.DungeonActionRequestData{
						Sentence: "look",
					},
				}
				return &res
			},
			expectResponseBody: func(data harness.Data) *schema.DungeonActionResponse {
				res := schema.DungeonActionResponse{
					Data: []schema.DungeonActionResponseData{
						{
							Action: schema.ActionData{
								Command:                   "look",
								TargetDungeonLocationName: "Cave Entrance",
							},
						},
					},
				}
				return &res
			},
			expectResponseCode: http.StatusOK,
		},
		{
			name: "POST - move north",
			config: func(rnr *Runner) server.HandlerConfig {
				return rnr.HandlerConfig[6]
			},
			requestHeaders: func(data harness.Data) map[string]string {
				headers := map[string]string{
					"Authorization": "Bearer " + validAuthToken(),
				}
				return headers
			},
			requestParams: func(data harness.Data) map[string]string {
				params := map[string]string{
					":dungeon_id":   data.DungeonRecs[0].ID,
					":character_id": data.DungeonCharacterRecs[0].ID,
				}
				return params
			},
			requestBody: func(data harness.Data) *schema.DungeonActionRequest {
				res := schema.DungeonActionRequest{
					Data: schema.DungeonActionRequestData{
						Sentence: "move north",
					},
				}
				return &res
			},
			expectResponseBody: func(data harness.Data) *schema.DungeonActionResponse {
				res := schema.DungeonActionResponse{
					Data: []schema.DungeonActionResponseData{
						{
							Action: schema.ActionData{
								Command:                   "move",
								TargetDungeonLocationName: "Cave Tunnel",
							},
						},
					},
				}
				return &res
			},
			expectResponseCode: http.StatusOK,
		},
		// {
		// 	name: "POST - empty",
		// 	config: func(rnr *Runner) server.HandlerConfig {
		// 		return rnr.HandlerConfig[6]
		// 	},
		// 	requestHeaders: func(data harness.Data) map[string]string {
		// 		headers := map[string]string{
		// 			"Authorization": "Bearer " + validAuthToken(),
		// 		}
		// 		return headers
		// 	},
		// 	requestParams: func(data harness.Data) map[string]string {
		// 		params := map[string]string{
		// 			":dungeon_id":   data.DungeonRecs[0].ID,
		// 			":character_id": data.DungeonCharacterRecs[0].ID,
		// 		}
		// 		return params
		// 	},
		// 	requestBody: func(data harness.Data) *schema.DungeonActionRequest {
		// 		res := schema.DungeonActionRequest{
		// 			Data: schema.DungeonActionRequestData{
		// 				Sentence: "",
		// 			},
		// 		}
		// 		return &res
		// 	},
		// 	responseCode: http.StatusBadRequest,
		// },
	}

	for _, tc := range tests {

		t.Logf("Running test >%s<", tc.name)

		func() {
			rnr := NewRunner()

			err = rnr.Init(th.Config, th.Log, th.Store)
			require.NoError(t, err, "Runner init returns without error")

			err = th.Setup()
			require.NoError(t, err, "Test data setup returns without error")
			defer func() {
				err = th.Teardown()
				require.NoError(t, err, "Test data teardown returns without error")
			}()

			// config
			cfg := tc.config(rnr)

			// handler
			h, _ := rnr.DefaultMiddleware(cfg, cfg.HandlerFunc)

			// router
			rtr := httprouter.New()

			switch cfg.Method {
			case http.MethodGet:
				rtr.GET(cfg.Path, h)
			case http.MethodPost:
				rtr.POST(cfg.Path, h)
			case http.MethodPut:
				rtr.PUT(cfg.Path, h)
			case http.MethodDelete:
				rtr.DELETE(cfg.Path, h)
			default:
				//
			}

			// request params
			requestParams := map[string]string{}
			if tc.requestParams != nil {
				requestParams = tc.requestParams(th.Data)
			}

			requestPath := cfg.Path
			for paramKey, paramValue := range requestParams {
				requestPath = strings.Replace(requestPath, paramKey, paramValue, 1)
			}

			// query params
			queryParams := map[string]string{}
			if tc.queryParams != nil {
				queryParams = tc.queryParams(th.Data)
			}

			if len(queryParams) > 0 {
				count := 0
				for paramKey, paramValue := range queryParams {
					if count == 0 {
						requestPath = requestPath + `?`
					} else {
						requestPath = requestPath + `&`
					}
					t.Logf("Adding parameter key >%s< param >%s<", paramKey, paramValue)
					requestPath = fmt.Sprintf("%s%s=%s", requestPath, paramKey, url.QueryEscape(paramValue))
				}
				t.Logf("Resulting requestPath >%s<", requestPath)
			}

			// request data
			data := tc.requestBody(th.Data)

			var req *http.Request

			if data != nil {
				jsonData, err := json.Marshal(data)
				require.NoError(t, err, "Marshal returns without error")

				req, err = http.NewRequest(cfg.Method, requestPath, bytes.NewBuffer(jsonData))
				require.NoError(t, err, "NewRequest returns without error")
			} else {
				req, err = http.NewRequest(cfg.Method, requestPath, nil)
				require.NoError(t, err, "NewRequest returns without error")
			}

			// request headers
			requestHeaders := map[string]string{}
			if tc.requestHeaders != nil {
				requestHeaders = tc.requestHeaders(th.Data)
			}

			for headerKey, headerVal := range requestHeaders {
				req.Header.Add(headerKey, headerVal)
			}

			// recorder
			rec := httptest.NewRecorder()

			// serve
			rtr.ServeHTTP(rec, req)

			// test status
			require.Equalf(t, tc.expectResponseCode, rec.Code, "%s - Response code equals expected", tc.name)

			responseBody := schema.DungeonActionResponse{}
			err = json.NewDecoder(rec.Body).Decode(&responseBody)
			require.NoError(t, err, "Decode returns without error")

			if tc.expectResponseCode != http.StatusOK {
				return
			}

			// test data
			if tc.expectResponseBody != nil {
				expectResponseBody := tc.expectResponseBody(th.Data)
				if expectResponseBody != nil {
					for idx, expectData := range expectResponseBody.Data {
						require.NotNil(t, responseBody.Data[idx], "Response body index is not empty")
						require.NotNil(t, responseBody.Data[idx].Action, "Response body action is not empty")
						require.Equal(t, responseBody.Data[idx].Action.Command, expectData.Action.Command)
					}
				}
			}

			require.NotNil(t, responseBody, "Response body is not nil")
			require.GreaterOrEqual(t, len(responseBody.Data), 0, "Response body data ")

			for _, data := range responseBody.Data {

				// record timestamps
				require.False(t, data.Action.CreatedAt.IsZero(), "CreatedAt is not zero")
				if cfg.Method == http.MethodPost {
					require.True(t, data.Action.UpdatedAt.IsZero(), "UpdatedAt is zero")
				}
				if cfg.Method == http.MethodPut {
					require.False(t, data.Action.UpdatedAt.IsZero(), "UpdatedAt is not zero")
				}

				if cfg.Method == http.MethodPost {
					if len(responseBody.Data) != 0 {
						t.Logf("Adding dungeon character action ID >%s< for teardown", data.Action.ID)
						th.AddDungeonCharacterActionTeardownID(data.Action.ID)
					}
				}
			}
		}()
	}
}
