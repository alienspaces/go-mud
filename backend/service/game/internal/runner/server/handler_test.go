package runner

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/go-mud/backend/core/jsonschema"
	"gitlab.com/alienspaces/go-mud/backend/core/server"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/harness"
)

// TestCaser is the interface required by the RunTestCase function
type TestCaser interface {
	TestName() string
	TestHandlerConfig(rnr *Runner) server.HandlerConfig
	TestDataConfig() *harness.DataConfig
	TestRequestHeaders(data harness.Data) map[string]string
	TestRequestPathParams(data harness.Data) map[string]string
	TestRequestQueryParams(data harness.Data) map[string]interface{}
	TestRequestBody(data harness.Data) interface{}
	TestResponseDecoder(body io.Reader) (interface{}, error)
	TestResponseCode() int
	TestShouldSetupTeardown() bool
	TestShouldTxCommit() bool
}

// TestCase is the base test case class for all tests cases to extend
type TestCase struct {
	Skip                bool
	Name                string
	HandlerConfig       func(rnr *Runner) server.HandlerConfig
	DataConfig          func() *harness.DataConfig
	RequestHeaders      func(data harness.Data) map[string]string
	RequestPathParams   func(data harness.Data) map[string]string
	RequestQueryParams  func(data harness.Data) map[string]interface{}
	RequestBody         func(data harness.Data) interface{}
	ResponseDecoder     func(body io.Reader) (interface{}, error)
	ResponseCode        int
	ShouldSetupTeardown bool // Should running the test automatically run the harness data setup and teardown
	ShouldTxCommit      bool // Should running the test automatically include the rollback header
}

//lint:ignore U1000 - testing struct implements interface
var _tc TestCaser = &TestCase{}

func (t *TestCase) TestName() string {
	return t.Name
}

func (t *TestCase) TestHandlerConfig(rnr *Runner) server.HandlerConfig {
	if t.HandlerConfig != nil {
		return t.HandlerConfig(rnr)
	}
	return server.HandlerConfig{}
}

func (t *TestCase) TestDataConfig() *harness.DataConfig {
	if t.DataConfig != nil {
		return t.DataConfig()
	}
	return nil
}

func (t *TestCase) TestRequestHeaders(data harness.Data) map[string]string {
	headers := map[string]string{}
	if t.RequestHeaders != nil {
		headers = t.RequestHeaders(data)
	}

	if !t.ShouldTxCommit {
		headers[server.HeaderXTxRollback] = "true"
	}

	return headers
}

func (t *TestCase) TestRequestPathParams(data harness.Data) map[string]string {
	if t.RequestPathParams != nil {
		return t.RequestPathParams(data)
	}
	return nil
}

func (t *TestCase) TestRequestQueryParams(data harness.Data) map[string]interface{} {
	if t.RequestQueryParams != nil {
		return t.RequestQueryParams(data)
	}
	return nil
}

func (t *TestCase) TestRequestBody(data harness.Data) interface{} {
	var b interface{}
	if t.RequestBody != nil {
		b = t.RequestBody(data)
	}
	return b
}

func (t *TestCase) TestResponseDecoder(body io.Reader) (interface{}, error) {
	var b interface{}
	var err error
	if t.ResponseDecoder != nil {
		b, err = t.ResponseDecoder(body)
	}
	return b, err
}

func (t *TestCase) TestResponseCode() int {
	return t.ResponseCode
}

func (t *TestCase) TestShouldSetupTeardown() bool {
	return t.ShouldSetupTeardown
}

func (t *TestCase) TestShouldTxCommit() bool {
	return t.ShouldTxCommit
}

func RunTestCase(t *testing.T, th *harness.Testing, tc TestCaser, tf func(method string, body interface{})) {

	rnr, err := NewRunner(th.Config, th.Log)
	require.NoError(t, err, "Runner init returns without error")

	err = rnr.Init(th.Store)
	require.NoError(t, err, "Runner init returns without error")

	if tc.TestShouldSetupTeardown() {
		dataConfig := tc.TestDataConfig()
		if dataConfig != nil {
			th.DataConfig = *dataConfig
		}
		_, err = th.Setup()
		require.NoError(t, err, "Test data setup returns without error")
		defer func() {
			err = th.Teardown()
			require.NoError(t, err, "Test data teardown returns without error")
		}()
	}

	// Handler config
	cfg := tc.TestHandlerConfig(rnr)

	// Handler
	h, _ := rnr.ApplyMiddleware(cfg, cfg.HandlerFunc)

	// Router
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

	// Request params
	requestParams := tc.TestRequestPathParams(th.Data)

	requestPath := cfg.Path
	for paramKey, paramValue := range requestParams {
		requestPath = strings.Replace(requestPath, paramKey, paramValue, 1)
	}

	// Query params
	queryParams := tc.TestRequestQueryParams(th.Data)

	if len(queryParams) > 0 {
		requestPath += `?`
		for paramKey, paramValue := range queryParams {
			t.Logf("Adding parameter key >%s< param >%s<", paramKey, paramValue)
			switch v := paramValue.(type) {
			case int:
				requestPath = fmt.Sprintf("%s%s=%d&", requestPath, paramKey, v)
			case string:
				requestPath = fmt.Sprintf("%s%s=%s&", requestPath, paramKey, url.QueryEscape(v))
			case bool:
				requestPath = fmt.Sprintf("%s%s=%t&", requestPath, paramKey, v)
			default:
				t.Errorf("Unsupported query parameter type for value >%v<", v)
			}
		}
		t.Logf("Request path with query params >%s<", requestPath)
	}

	t.Logf(">>> Request path >%s<", requestPath)

	// Request data
	data := tc.TestRequestBody(th.Data)

	var req *http.Request

	if data != nil {
		jsonData, err := json.Marshal(data)
		require.NoError(t, err, "Marshal returns without error")

		t.Logf("++++ Posting JSON data >%s<", jsonData)

		req, err = http.NewRequest(cfg.Method, requestPath, bytes.NewBuffer(jsonData))
		require.NoError(t, err, "NewRequest returns without error")
	} else {
		req, err = http.NewRequest(cfg.Method, requestPath, nil)
		require.NoError(t, err, "NewRequest returns without error")
	}

	// Request headers
	requestHeaders := tc.TestRequestHeaders(th.Data)

	for headerKey, headerVal := range requestHeaders {
		req.Header.Add(headerKey, headerVal)
	}

	// Recorder
	recorder := httptest.NewRecorder()

	// Serve
	rtr.ServeHTTP(recorder, req)

	// Test status
	require.Equalf(t, tc.TestResponseCode(), recorder.Code, "%s - Response code equals expected", tc.TestName())

	var responseBody interface{}

	// Validate response body
	if recorder.Code == http.StatusOK || recorder.Code == http.StatusCreated {

		// Response body
		responseBody, err = tc.TestResponseDecoder(recorder.Body)
		require.NoError(t, err, "Response body decodes without error")

		if responseBody != nil {
			jsonData, err := json.Marshal(responseBody)
			require.NoError(t, err, "Marshal returns without error")

			testResponseSchema(t, cfg, jsonData)

			result, err := jsonschema.Validate(cfg.MiddlewareConfig.ValidateResponseSchema, string(jsonData))
			require.NoError(t, err, "Validates against JSON response schema without error")
			require.NotNil(t, result, "JSON response schema validation result is not nil")
			require.True(t, result.Valid(), fmt.Sprintf("JSON response schema validation result is valid >%+v<", result.Errors()))
		}
	}

	if tf != nil {
		tf(cfg.Method, responseBody)
	}
}

func testResponseSchema(t *testing.T, hc server.HandlerConfig, actualRes interface{}) {

	schema := hc.MiddlewareConfig.ValidateResponseSchema
	schemaMain := schema.Main
	require.NotEmpty(t, schemaMain.Location, "handler >%s %s< ValidateResponseSchema main location path should not be empty", hc.Method, hc.Path)
	require.NotEmpty(t, schemaMain.Name, "handler >%s %s< ValidateResponseSchema main filename should not be empty", hc.Method, hc.Path)

	for _, r := range schema.References {
		require.NotEmpty(t, r.Location, "handler >%s %s< ValidateResponseSchema reference location path should not be empty", hc.Method, hc.Path)
		require.NotEmpty(t, r.Name, "handler >%s %s< ValidateResponseSchema reference filename should not be empty", hc.Method, hc.Path)
	}

	testSchemaHelper(t, schema, actualRes)
}

func testSchemaHelper(t *testing.T, s *jsonschema.SchemaWithReferences, actualRes interface{}) {
	result, err := jsonschema.Validate(s, actualRes)
	require.NoError(t, err, "schema validation should not error")
	err = jsonschema.MapError(result)
	require.NoError(t, err, "schema validation results should be empty")
}
