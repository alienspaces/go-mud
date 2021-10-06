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

	"gitlab.com/alienspaces/go-boilerplate/server/core/auth"
	"gitlab.com/alienspaces/go-boilerplate/server/core/server"
	"gitlab.com/alienspaces/go-boilerplate/server/schema"
	"gitlab.com/alienspaces/go-boilerplate/server/service/template/internal/harness"
)

func TestTemplateHandler(t *testing.T) {

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
		requestData    func(data *harness.Data) *schema.TemplateRequest
		responseCode   int
		responseData   func(data *harness.Data) *schema.TemplateResponse
	}

	// validAuthToken - Generate a valid authentication token for this handler
	validAuthToken := func() string {
		authen, _ := auth.NewAuth(c, l)
		token, _ := authen.EncodeJWT(&auth.Claims{
			Roles:    []string{},
			Identity: map[string]interface{}{},
		})
		return token
	}

	tests := []TestCase{
		{
			name: "GET - Get existing",
			config: func(rnr *Runner) server.HandlerConfig {
				return rnr.HandlerConfig[1]
			},
			requestHeaders: func(data *harness.Data) map[string]string {
				headers := map[string]string{
					"Authorization": "Bearer " + validAuthToken(),
				}
				return headers
			},
			requestParams: func(data *harness.Data) map[string]string {
				params := map[string]string{
					":template_id": data.TemplateRecs[0].ID,
				}
				return params
			},
			requestData: func(data *harness.Data) *schema.TemplateRequest {
				return nil
			},
			responseCode: http.StatusOK,
			responseData: func(data *harness.Data) *schema.TemplateResponse {
				res := schema.TemplateResponse{
					Data: []schema.TemplateData{
						{
							ID: data.TemplateRecs[0].ID,
						},
					},
				}
				return &res
			},
		},
		{
			name: "GET - Get non-existant",
			config: func(rnr *Runner) server.HandlerConfig {
				return rnr.HandlerConfig[1]
			},
			requestHeaders: func(data *harness.Data) map[string]string {
				headers := map[string]string{
					"Authorization": "Bearer " + validAuthToken(),
				}
				return headers
			},
			requestParams: func(data *harness.Data) map[string]string {
				params := map[string]string{
					":template_id": "17c19414-2d15-4d20-8fc3-36fc10341dc8",
				}
				return params
			},
			requestData: func(data *harness.Data) *schema.TemplateRequest {
				return nil
			},
			responseCode: http.StatusNotFound,
		},
		{
			name: "POST - Create without ID",
			config: func(rnr *Runner) server.HandlerConfig {
				return rnr.HandlerConfig[2]
			},
			requestHeaders: func(data *harness.Data) map[string]string {
				headers := map[string]string{
					"Authorization": "Bearer " + validAuthToken(),
				}
				return headers
			},
			requestData: func(data *harness.Data) *schema.TemplateRequest {
				req := schema.TemplateRequest{
					Data: schema.TemplateData{},
				}
				return &req
			},
			responseCode: http.StatusOK,
		},
		{
			name: "POST - Create with ID",
			config: func(rnr *Runner) server.HandlerConfig {
				return rnr.HandlerConfig[3]
			},
			requestHeaders: func(data *harness.Data) map[string]string {
				headers := map[string]string{
					"Authorization": "Bearer " + validAuthToken(),
				}
				return headers
			},
			requestParams: func(data *harness.Data) map[string]string {
				params := map[string]string{
					":template_id": "e3a9e0f8-ce9c-477b-8b93-cf4da03af4c9",
				}
				return params
			},
			requestData: func(data *harness.Data) *schema.TemplateRequest {
				req := schema.TemplateRequest{
					Data: schema.TemplateData{},
				}
				return &req
			},
			responseCode: http.StatusOK,
			responseData: func(data *harness.Data) *schema.TemplateResponse {
				res := schema.TemplateResponse{
					Data: []schema.TemplateData{
						{
							ID: "e3a9e0f8-ce9c-477b-8b93-cf4da03af4c9",
						},
					},
				}
				return &res
			},
		},
		{
			name: "PUT - Update existing",
			config: func(rnr *Runner) server.HandlerConfig {
				return rnr.HandlerConfig[4]
			},
			requestHeaders: func(data *harness.Data) map[string]string {
				headers := map[string]string{
					"Authorization": "Bearer " + validAuthToken(),
				}
				return headers
			},
			requestParams: func(data *harness.Data) map[string]string {
				params := map[string]string{
					":template_id": data.TemplateRecs[0].ID,
				}
				return params
			},
			requestData: func(data *harness.Data) *schema.TemplateRequest {
				req := schema.TemplateRequest{
					Data: schema.TemplateData{
						ID: data.TemplateRecs[0].ID,
					},
				}
				return &req
			},
			responseCode: http.StatusOK,
			responseData: func(data *harness.Data) *schema.TemplateResponse {
				res := schema.TemplateResponse{
					Data: []schema.TemplateData{
						{
							ID: data.TemplateRecs[0].ID,
						},
					},
				}
				return &res
			},
		},
		{
			name: "PUT - Update non-existing",
			config: func(rnr *Runner) server.HandlerConfig {
				return rnr.HandlerConfig[4]
			},
			requestHeaders: func(data *harness.Data) map[string]string {
				headers := map[string]string{
					"Authorization": "Bearer " + validAuthToken(),
				}
				return headers
			},
			requestParams: func(data *harness.Data) map[string]string {
				params := map[string]string{
					":template_id": "17c19414-2d15-4d20-8fc3-36fc10341dc8",
				}
				return params
			},
			requestData: func(data *harness.Data) *schema.TemplateRequest {
				req := schema.TemplateRequest{
					Data: schema.TemplateData{
						ID: data.TemplateRecs[0].ID,
					},
				}
				return &req
			},
			responseCode: http.StatusNotFound,
		},
		{
			name: "PUT - Update missing data",
			requestHeaders: func(data *harness.Data) map[string]string {
				headers := map[string]string{
					"Authorization": "Bearer " + validAuthToken(),
				}
				return headers
			},
			config: func(rnr *Runner) server.HandlerConfig {
				return rnr.HandlerConfig[3]
			},
			requestData: func(data *harness.Data) *schema.TemplateRequest {
				return nil
			},
			responseCode: http.StatusBadRequest,
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
				t.Logf("Resulting requestPath >%s<", requestPath)
			}

			// request data
			data := tc.requestData(th.Data)

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
			require.Equalf(t, tc.responseCode, rec.Code, "%s - Response code equals expected", tc.name)

			res := schema.TemplateResponse{}
			err = json.NewDecoder(rec.Body).Decode(&res)
			require.NoError(t, err, "Decode returns without error")

			// response data
			var resData *schema.TemplateResponse
			if tc.responseData != nil {
				resData = tc.responseData(th.Data)
			}

			// test data
			if tc.responseCode == http.StatusOK {

				// response data
				if resData != nil {
					require.Equal(t, resData.Data[0].ID, res.Data[0].ID, "ID equals expected")
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
