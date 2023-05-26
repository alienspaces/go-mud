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

	"gitlab.com/alienspaces/go-mud/backend/core/collection/set"
	coreerror "gitlab.com/alienspaces/go-mud/backend/core/error"
	"gitlab.com/alienspaces/go-mud/backend/core/jsonschema"
	"gitlab.com/alienspaces/go-mud/backend/core/server"
	"gitlab.com/alienspaces/go-mud/backend/service/template/internal/harness"
)

type ExpectedErrorResponse struct {
	err  coreerror.Error // For convenience, this allows the test to specify the error, without wrapping it in a slice
	errs set.Set[coreerror.ErrorCode]
}

type TestCaser interface {
	TestName() string
	TestHandlerConfig(rnr *Runner) server.HandlerConfig
	TestRequestHeaders(data *harness.Data) map[string]string
	TestRequestPathParams(data *harness.Data) map[string]string
	TestRequestQueryParams(data *harness.Data) map[string]interface{}
	TestRequestBody(data *harness.Data) interface{}
	TestResponseBody(body io.Reader) (interface{}, error)
	TestResponseCode() int
	TestResponseError(data *harness.Data) *ExpectedErrorResponse
}

type TestCase struct {
	Skip               bool
	Name               string
	HandlerConfig      func(rnr *Runner) server.HandlerConfig
	RequestHeaders     func(data *harness.Data) map[string]string
	RequestPathParams  func(data *harness.Data) map[string]string
	RequestQueryParams func(data *harness.Data) map[string]interface{}
	RequestBody        func(data *harness.Data) interface{}
	ResponseBody       func(body io.Reader) (interface{}, error)
	ResponseCode       int
	ResponseError      func(*harness.Data) ExpectedErrorResponse
	ShouldTxCommit     bool
}

//lint:ignore U1000 - testing struct implements interface
var _tc TestCaser = &TestCase{}

func (t *TestCase) TestName() string {
	return t.Name
}

func (t *TestCase) TestHandlerConfig(rnr *Runner) server.HandlerConfig {
	return t.HandlerConfig(rnr)
}

func (t *TestCase) TestRequestHeaders(data *harness.Data) map[string]string {
	headers := map[string]string{}
	if t.RequestHeaders != nil {
		headers = t.RequestHeaders(data)
	}

	if !t.ShouldTxCommit {
		headers[server.HeaderXTxRollback] = "true"
	}

	return headers
}

func (t *TestCase) TestRequestPathParams(data *harness.Data) map[string]string {
	pp := map[string]string{}
	if t.RequestPathParams != nil {
		pp = t.RequestPathParams(data)
	}
	return pp
}

func (t *TestCase) TestRequestQueryParams(data *harness.Data) map[string]interface{} {
	qp := map[string]interface{}{}
	if t.RequestQueryParams != nil {
		qp = t.RequestQueryParams(data)
	}
	return qp
}

func testCaseResponseBodyGeneric[T any](body io.Reader) (any, error) {
	var responseBody T
	err := json.NewDecoder(body).Decode(&responseBody)
	return responseBody, err
}

func (t *TestCase) TestRequestBody(data *harness.Data) interface{} {
	var b interface{}
	if t.RequestBody != nil {
		b = t.RequestBody(data)
	}
	return b
}

func (t *TestCase) TestResponseBody(body io.Reader) (interface{}, error) {
	var b interface{}
	var err error
	if t.ResponseBody != nil {
		b, err = t.ResponseBody(body)
	}
	return b, err
}

func (t *TestCase) TestResponseCode() int {
	return t.ResponseCode
}

func (t *TestCase) TestResponseError(data *harness.Data) *ExpectedErrorResponse {
	if t.ResponseError != nil {
		respErr := t.ResponseError(data)

		if respErr.errs == nil {
			respErr.errs = set.New(respErr.err.ErrorCode)
		}

		return &respErr
	}
	return nil
}

func testResponseSchema(t *testing.T, hc server.HandlerConfig, actualRes interface{}) {
	schema := hc.MiddlewareConfig.ValidateResponseSchema
	schemaMain := schema.Main
	require.NotEmpty(t, schemaMain.Location, "handler >%s %s< schema main location path should not be empty", hc.Method, hc.Path)
	require.NotEmpty(t, schemaMain.Name, "handler >%s %s< schema main filename should not be empty", hc.Method, hc.Path)

	for _, r := range schema.References {
		require.NotEmpty(t, r.Location, "handler >%s %s< schema ref location path should not be empty", hc.Method, hc.Path)
		require.NotEmpty(t, r.Name, "handler >%s %s< schema ref filename should not be empty", hc.Method, hc.Path)
	}

	testSchemaHelper(t, schema, actualRes)
}

func testSchemaHelper(t *testing.T, s jsonschema.SchemaWithReferences, actualRes interface{}) {
	result, err := jsonschema.Validate(s, actualRes)
	require.NoError(t, err, "schema validation should not error")

	err = jsonschema.MapError(result)
	require.NoError(t, err, "schema validation results should be empty")
}

func RunTestCase(t *testing.T, th *harness.Testing, tc TestCaser, tf func(method string, body interface{})) {

	rnr, err := NewRunner(th.Config, th.Log)
	require.NoError(t, err, "Runner init returns without error")

	err = rnr.Init(th.Store)
	require.NoError(t, err, "Runner init returns without error")

	_, err = th.Setup()
	require.NoError(t, err, "Test data setup returns without error")
	defer func() {
		err = th.Teardown()
		require.NoError(t, err, "Test data teardown returns without error")
	}()

	// config
	cfg := tc.TestHandlerConfig(rnr)

	// handler
	h, _ := rnr.ApplyMiddleware(cfg, cfg.HandlerFunc)

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
	requestParams := tc.TestRequestPathParams(th.Data)

	requestPath := cfg.Path
	for paramKey, paramValue := range requestParams {
		requestPath = strings.Replace(requestPath, paramKey, paramValue, 1)
	}

	// query params
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
		t.Logf("Resulting requestPath >%s<", requestPath)
	}

	// request data
	data := tc.TestRequestBody(th.Data)

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
	requestHeaders := tc.TestRequestHeaders(th.Data)

	for headerKey, headerVal := range requestHeaders {
		req.Header.Add(headerKey, headerVal)
	}

	// recorder
	rec := httptest.NewRecorder()

	// serve
	rtr.ServeHTTP(rec, req)

	// test status
	if tc.TestResponseCode() != rec.Code {
		t.Logf("%s", rec.Body.String())
	}
	require.Equalf(t, tc.TestResponseCode(), rec.Code, "%s - Response code equals expected", tc.TestName())

	// Test expected error response
	expectedErr := tc.TestResponseError(th.Data)
	if expectedErr != nil {
		var actualErrs []coreerror.Error

		err = json.NewDecoder(rec.Body).Decode(&actualErrs)
		require.NoError(t, err, "Decode returns without error")

		for _, actual := range actualErrs {
			ok := expectedErr.errs.Contains(actual.ErrorCode)
			require.True(t, ok, "expected >%#v< actual >%v<", expectedErr, actual.ErrorCode)
		}
	}

	responseBody, err := tc.TestResponseBody(rec.Body)
	require.NoError(t, err, "Response body decodes without error")

	// Test response body
	if rec.Code == http.StatusOK || rec.Code == http.StatusCreated {
		require.NotNil(t, responseBody, "Response body is not nil")

		jsonData, err := json.Marshal(responseBody)
		require.NoError(t, err, "Marshal returns without error")

		testResponseSchema(t, cfg, jsonData)

		result, err := jsonschema.Validate(cfg.MiddlewareConfig.ValidateResponseSchema, string(jsonData))
		require.NoError(t, err, "Validates against JSON response schema without error")
		require.NotNil(t, result, "JSON response schema validation result is not nil")
		require.True(t, result.Valid(), fmt.Sprintf("JSON response schema validation result is valid >%+v<", result.Errors()))
	}

	if tf != nil {
		tf(cfg.Method, responseBody)
	}
}
