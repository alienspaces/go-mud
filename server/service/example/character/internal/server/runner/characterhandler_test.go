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

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/go-boilerplate/server/constant"
	"gitlab.com/alienspaces/go-boilerplate/server/core/auth"
	"gitlab.com/alienspaces/go-boilerplate/server/core/server"
	"gitlab.com/alienspaces/go-boilerplate/server/schema"
	"gitlab.com/alienspaces/go-boilerplate/server/service/character/internal/harness"
)

func TestCharacterHandler(t *testing.T) {

	// Test harness
	th, err := NewTestHarness()
	require.NoError(t, err, "New test data returns without error")

	//  Test dependencies
	c, l, s, err := th.NewDefaultDependencies()
	require.NoError(t, err, "NewDefaultDependencies returns without error")

	type TestCase struct {
		name           string
		config         func(rnr *Runner) server.HandlerConfig
		requestHeaders func(data *harness.Data) map[string]string
		requestParams  func(data *harness.Data) map[string]string
		queryParams    func(data *harness.Data) map[string]string
		requestData    func(data *harness.Data) *schema.CharacterRequest
		responseCode   int
		responseData   func(data *harness.Data) *schema.CharacterResponse
	}

	// validAuthToken - Generate a valid authentication token for this handler
	validAuthToken := func(roles []string, identity map[string]interface{}) string {
		authen, _ := auth.NewAuth(c, l)
		token, _ := authen.EncodeJWT(&auth.Claims{
			Roles:    roles,
			Identity: identity,
		})
		return token
	}

	tests := []TestCase{
		{
			name: "GET - Get existing resource, admin role, no identity",
			config: func(rnr *Runner) server.HandlerConfig {
				return rnr.HandlerConfig[1]
			},
			requestHeaders: func(data *harness.Data) map[string]string {
				roles := []string{
					constant.AuthRoleAdministrator,
				}
				identity := map[string]interface{}{}
				headers := map[string]string{
					"Authorization": "Bearer " + validAuthToken(roles, identity),
				}
				return headers
			},
			requestParams: func(data *harness.Data) map[string]string {
				params := map[string]string{
					":character_id": data.CharacterRecs[0].ID,
				}
				return params
			},
			requestData: func(data *harness.Data) *schema.CharacterRequest {
				return nil
			},
			responseCode: http.StatusOK,
			responseData: func(data *harness.Data) *schema.CharacterResponse {
				res := schema.CharacterResponse{
					Data: []schema.CharacterData{
						{
							ID:               data.CharacterRecs[0].ID,
							PlayerID:         data.CharacterRecs[0].PlayerID,
							Name:             data.CharacterRecs[0].Name,
							Avatar:           data.CharacterRecs[0].Avatar,
							Strength:         data.CharacterRecs[0].Strength,
							Dexterity:        data.CharacterRecs[0].Dexterity,
							Intelligence:     data.CharacterRecs[0].Intelligence,
							AttributePoints:  data.CharacterRecs[0].AttributePoints,
							ExperiencePoints: data.CharacterRecs[0].ExperiencePoints,
							Coins:            data.CharacterRecs[0].Coins,
						},
					},
				}
				return &res
			},
		},
		{
			name: "GET - Get missing resource, admin role, no identity",
			config: func(rnr *Runner) server.HandlerConfig {
				return rnr.HandlerConfig[1]
			},
			requestHeaders: func(data *harness.Data) map[string]string {
				roles := []string{
					constant.AuthRoleAdministrator,
				}
				identity := map[string]interface{}{}
				headers := map[string]string{
					"Authorization": "Bearer " + validAuthToken(roles, identity),
				}
				return headers
			},
			requestParams: func(data *harness.Data) map[string]string {
				params := map[string]string{
					":character_id": "17c19414-2d15-4d20-8fc3-36fc10341dc8",
				}
				return params
			},
			requestData: func(data *harness.Data) *schema.CharacterRequest {
				return nil
			},
			responseCode: http.StatusNotFound,
			responseData: func(data *harness.Data) *schema.CharacterResponse {
				return nil
			},
		},
		{
			name: "GET - Get existing resource, default role, account identity matches",
			config: func(rnr *Runner) server.HandlerConfig {
				return rnr.HandlerConfig[3]
			},
			requestHeaders: func(data *harness.Data) map[string]string {
				roles := []string{
					constant.AuthRoleDefault,
				}
				identity := map[string]interface{}{
					"player_id": data.CharacterRecs[0].PlayerID,
				}
				headers := map[string]string{
					"Authorization": "Bearer " + validAuthToken(roles, identity),
				}
				return headers
			},
			requestParams: func(data *harness.Data) map[string]string {
				params := map[string]string{
					":player_id":    data.CharacterRecs[0].PlayerID,
					":character_id": data.CharacterRecs[0].ID,
				}
				return params
			},
			requestData: func(data *harness.Data) *schema.CharacterRequest {
				return nil
			},
			responseCode: http.StatusOK,
			responseData: func(data *harness.Data) *schema.CharacterResponse {
				res := schema.CharacterResponse{
					Data: []schema.CharacterData{
						{
							ID:               data.CharacterRecs[0].ID,
							PlayerID:         data.CharacterRecs[0].PlayerID,
							Name:             data.CharacterRecs[0].Name,
							Avatar:           data.CharacterRecs[0].Avatar,
							Strength:         data.CharacterRecs[0].Strength,
							Dexterity:        data.CharacterRecs[0].Dexterity,
							Intelligence:     data.CharacterRecs[0].Intelligence,
							AttributePoints:  data.CharacterRecs[0].AttributePoints,
							ExperiencePoints: data.CharacterRecs[0].ExperiencePoints,
							Coins:            data.CharacterRecs[0].Coins,
						},
					},
				}
				return &res
			},
		},
		{
			name: "POST - Create without character ID, default role, account identity matches",
			config: func(rnr *Runner) server.HandlerConfig {
				return rnr.HandlerConfig[4]
			},
			requestHeaders: func(data *harness.Data) map[string]string {
				roles := []string{
					constant.AuthRoleDefault,
				}
				identity := map[string]interface{}{
					"player_id": data.CharacterRecs[0].ID,
				}
				headers := map[string]string{
					"Authorization": "Bearer " + validAuthToken(roles, identity),
				}
				return headers
			},
			requestParams: func(data *harness.Data) map[string]string {
				params := map[string]string{
					":player_id": data.CharacterRecs[0].ID,
				}
				return params
			},
			requestData: func(data *harness.Data) *schema.CharacterRequest {
				req := schema.CharacterRequest{
					Data: schema.CharacterData{
						Name: "Veronica The Incredible",
					},
				}
				return &req
			},
			responseCode: http.StatusOK,
		},
		{
			name: "POST - Create with ID, default role, account identity matches",
			config: func(rnr *Runner) server.HandlerConfig {
				return rnr.HandlerConfig[5]
			},
			requestHeaders: func(data *harness.Data) map[string]string {
				roles := []string{
					constant.AuthRoleDefault,
				}
				identity := map[string]interface{}{
					"player_id": data.CharacterRecs[0].PlayerID,
				}
				headers := map[string]string{
					"Authorization": "Bearer " + validAuthToken(roles, identity),
				}
				return headers
			},
			requestParams: func(data *harness.Data) map[string]string {
				params := map[string]string{
					":player_id":    data.CharacterRecs[0].PlayerID,
					":character_id": "e3a9e0f8-ce9c-477b-8b93-cf4da03af4c9",
				}
				return params
			},
			requestData: func(data *harness.Data) *schema.CharacterRequest {
				req := schema.CharacterRequest{
					Data: schema.CharacterData{
						Name: "Audrey The Amazing",
					},
				}
				return &req
			},
			responseCode: http.StatusOK,
			responseData: func(data *harness.Data) *schema.CharacterResponse {
				res := schema.CharacterResponse{
					Data: []schema.CharacterData{
						{
							ID:              "e3a9e0f8-ce9c-477b-8b93-cf4da03af4c9",
							PlayerID:        data.CharacterRecs[0].PlayerID,
							Name:            "Audrey The Amazing",
							AttributePoints: 32,
						},
					},
				}
				return &res
			},
		},
		{
			name: "PUT - Update with ID, default role, account identity matches",
			config: func(rnr *Runner) server.HandlerConfig {
				return rnr.HandlerConfig[6]
			},
			requestHeaders: func(data *harness.Data) map[string]string {
				roles := []string{
					constant.AuthRoleDefault,
				}
				identity := map[string]interface{}{
					"player_id": data.CharacterRecs[0].PlayerID,
				}
				headers := map[string]string{
					"Authorization": "Bearer " + validAuthToken(roles, identity),
				}
				return headers
			},
			requestParams: func(data *harness.Data) map[string]string {
				params := map[string]string{
					":player_id":    data.CharacterRecs[0].PlayerID,
					":character_id": data.CharacterRecs[0].ID,
				}
				return params
			},
			requestData: func(data *harness.Data) *schema.CharacterRequest {
				req := schema.CharacterRequest{
					Data: schema.CharacterData{
						ID:               data.CharacterRecs[0].ID,
						PlayerID:         data.CharacterRecs[0].PlayerID,
						Name:             "Barricade Block",
						Avatar:           data.CharacterRecs[0].Avatar,
						Strength:         data.CharacterRecs[0].Strength,
						Dexterity:        data.CharacterRecs[0].Dexterity,
						Intelligence:     data.CharacterRecs[0].Intelligence,
						AttributePoints:  data.CharacterRecs[0].AttributePoints,
						ExperiencePoints: data.CharacterRecs[0].ExperiencePoints,
						Coins:            data.CharacterRecs[0].Coins,
					},
				}
				return &req
			},
			responseCode: http.StatusOK,
			responseData: func(data *harness.Data) *schema.CharacterResponse {
				res := schema.CharacterResponse{
					Data: []schema.CharacterData{
						{
							ID:               data.CharacterRecs[0].ID,
							PlayerID:         data.CharacterRecs[0].PlayerID,
							Name:             "Barricade Block",
							Avatar:           data.CharacterRecs[0].Avatar,
							Strength:         data.CharacterRecs[0].Strength,
							Dexterity:        data.CharacterRecs[0].Dexterity,
							Intelligence:     data.CharacterRecs[0].Intelligence,
							AttributePoints:  data.CharacterRecs[0].AttributePoints,
							ExperiencePoints: data.CharacterRecs[0].ExperiencePoints,
							Coins:            data.CharacterRecs[0].Coins,
						},
					},
				}
				return &res
			},
		},
		{
			name: "PUT - Update with ID, default role, account identity matches, missing data",
			config: func(rnr *Runner) server.HandlerConfig {
				return rnr.HandlerConfig[6]
			},
			requestHeaders: func(data *harness.Data) map[string]string {
				roles := []string{
					constant.AuthRoleDefault,
				}
				identity := map[string]interface{}{
					"player_id": data.CharacterRecs[0].ID,
				}
				headers := map[string]string{
					"Authorization": "Bearer " + validAuthToken(roles, identity),
				}
				return headers
			},
			requestParams: func(data *harness.Data) map[string]string {
				params := map[string]string{
					":player_id":    data.CharacterRecs[0].ID,
					":character_id": data.CharacterRecs[0].ID,
				}
				return params
			},
			requestData: func(data *harness.Data) *schema.CharacterRequest {
				return nil
			},
			responseCode: http.StatusBadRequest,
			responseData: func(data *harness.Data) *schema.CharacterResponse {
				return nil
			},
		},
	}

	for _, tc := range tests {

		t.Logf("Running test >%s<", tc.name)

		func() {
			rnr := NewRunner()

			err = rnr.Init(c, l, s)
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
			}

			t.Logf("Resulting requestPath >%s<", requestPath)

			// request data
			reqData := tc.requestData(th.Data)

			var req *http.Request

			if reqData != nil {
				jsonData, err := json.Marshal(reqData)
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
			require.Equalf(t, tc.responseCode, rec.Code, "%s - Response code equals expected", tc.name)

			res := schema.CharacterResponse{}
			err = json.NewDecoder(rec.Body).Decode(&res)
			require.NoError(t, err, "Decode returns without error")

			// response data
			var resData *schema.CharacterResponse
			if tc.responseData != nil {
				resData = tc.responseData(th.Data)
			}

			// test data
			if tc.responseCode == http.StatusOK {
				require.NotEmpty(t, res.Data, "Data is not empty")

				// response data
				if resData != nil {
					require.Equal(t, resData.Data[0].ID, res.Data[0].ID, "ID equals expected")
					require.Equal(t, resData.Data[0].PlayerID, res.Data[0].PlayerID, "PlayerID equals expected")
					require.Equal(t, resData.Data[0].Name, res.Data[0].Name, "Name equals expected")
					require.Equal(t, resData.Data[0].Strength, res.Data[0].Strength, "Strength equals expected")
					require.Equal(t, resData.Data[0].Dexterity, res.Data[0].Dexterity, "Dexterity equals expected")
					require.Equal(t, resData.Data[0].Intelligence, res.Data[0].Intelligence, "Intelligence equals expected")
					require.Equal(t, resData.Data[0].AttributePoints, res.Data[0].AttributePoints, "Attribute points equals expected")
					require.Equal(t, resData.Data[0].ExperiencePoints, res.Data[0].ExperiencePoints, "Experience points equals expected")
					require.Equal(t, resData.Data[0].Coins, res.Data[0].Coins, "Coins equals expected")
				}

				// record timestamps
				require.False(t, res.Data[0].CreatedAt.IsZero(), "CreatedAt is not zero")
				if cfg.Method == http.MethodPost {
					require.True(t, res.Data[0].UpdatedAt.IsZero(), "UpdatedAt is zero")
				}
				if cfg.Method == http.MethodPut {
					require.False(t, res.Data[0].UpdatedAt.IsZero(), "UpdatedAt is not zero")
				}
			}
		}()
	}
}
