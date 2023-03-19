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
	TestRequestQueryParams(data harness.Data) map[string]string
	TestRequestBody(data harness.Data) interface{}
	TestResponseBody(body io.Reader) (interface{}, error)
	TestResponseCode() int
}

// TestCase is the base test case class for all tests cases to extend
type TestCase struct {
	Skip               bool
	Name               string
	HandlerConfig      func(rnr *Runner) server.HandlerConfig
	DataConfig         func() *harness.DataConfig
	RequestHeaders     func(data harness.Data) map[string]string
	RequestPathParams  func(data harness.Data) map[string]string
	RequestQueryParams func(data harness.Data) map[string]string
	RequestBody        func(data harness.Data) interface{}
	ResponseBody       func(body io.Reader) (interface{}, error)
	ResponseCode       int
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
	if t.RequestHeaders != nil {
		return t.RequestHeaders(data)
	}
	return nil
}

func (t *TestCase) TestRequestPathParams(data harness.Data) map[string]string {
	if t.RequestPathParams != nil {
		return t.RequestPathParams(data)
	}
	return nil
}

func (t *TestCase) TestRequestQueryParams(data harness.Data) map[string]string {
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

func RunTestCase(t *testing.T, th *harness.Testing, tc TestCaser, tf func(method string, body interface{})) {

	rnr, err := NewRunner(th.Config, th.Log)
	require.NoError(t, err, "Runner init returns without error")

	err = rnr.Init(th.Store)
	require.NoError(t, err, "Runner init returns without error")

	// Data config
	dataConfig := tc.TestDataConfig()
	if dataConfig != nil {
		th.DataConfig = *dataConfig
	}

	err = th.Setup()
	require.NoError(t, err, "Test data setup returns without error")
	defer func() {
		err = th.Teardown()
		require.NoError(t, err, "Test data teardown returns without error")
	}()

	// Handler config
	handlerConfig := tc.TestHandlerConfig(rnr)

	// Handler
	h, _ := rnr.DefaultMiddleware(handlerConfig, handlerConfig.HandlerFunc)

	// Router
	rtr := httprouter.New()

	switch handlerConfig.Method {
	case http.MethodGet:
		rtr.GET(handlerConfig.Path, h)
	case http.MethodPost:
		rtr.POST(handlerConfig.Path, h)
	case http.MethodPut:
		rtr.PUT(handlerConfig.Path, h)
	case http.MethodDelete:
		rtr.DELETE(handlerConfig.Path, h)
	default:
		//
	}

	// Request params
	requestParams := tc.TestRequestPathParams(th.Data)

	requestPath := handlerConfig.Path
	for paramKey, paramValue := range requestParams {
		requestPath = strings.Replace(requestPath, paramKey, paramValue, 1)
	}

	// Query params
	queryParams := tc.TestRequestQueryParams(th.Data)

	if len(queryParams) > 0 {
		count := 0
		for paramKey, paramValue := range queryParams {
			if count == 0 {
				requestPath = requestPath + `?`
			} else {
				requestPath = requestPath + `&`
			}
			t.Logf("> Adding parameter key >%s< param >%s<", paramKey, paramValue)
			requestPath = fmt.Sprintf("%s%s=%s", requestPath, paramKey, url.QueryEscape(paramValue))
		}
	}

	// Request data
	data := tc.TestRequestBody(th.Data)

	var req *http.Request

	if data != nil {
		jsonData, err := json.Marshal(data)
		require.NoError(t, err, "Marshal returns without error")

		req, err = http.NewRequest(handlerConfig.Method, requestPath, bytes.NewBuffer(jsonData))
		require.NoError(t, err, "NewRequest returns without error")
	} else {
		req, err = http.NewRequest(handlerConfig.Method, requestPath, nil)
		require.NoError(t, err, "NewRequest returns without error")
	}

	// Request headers
	requestHeaders := tc.TestRequestHeaders(th.Data)

	for headerKey, headerVal := range requestHeaders {
		req.Header.Add(headerKey, headerVal)
	}

	// Recorder
	rec := httptest.NewRecorder()

	// Serve
	rtr.ServeHTTP(rec, req)

	// Test status
	require.Equalf(t, tc.TestResponseCode(), rec.Code, "%s - Response code equals expected", tc.TestName())

	// Response body
	responseBody, err := tc.TestResponseBody(rec.Body)
	require.NoError(t, err, "Response body decodes without error")

	// Validate response body
	if rec.Code == 200 || rec.Code == 201 {
		require.NotNil(t, responseBody, "Response body is not nil")

		jsonData, err := json.Marshal(responseBody)
		require.NoError(t, err, "Marshal returns without error")

		result, err := jsonschema.Validate(handlerConfig.MiddlewareConfig.ValidateResponseSchema, string(jsonData))
		require.NoError(t, err, "Validates against schema without error")
		t.Logf("Validation result errors >%+v< valid >%t<", result.Errors(), result.Valid())

		require.True(t, result.Valid(), "Validates against schema")
	}

	tf(handlerConfig.Method, responseBody)
}
