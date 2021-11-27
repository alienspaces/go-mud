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
	"gitlab.com/alienspaces/go-mud/server/core/server"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/harness"
)

type TestCaser interface {
	TestName() string
	TestHandlerConfig(rnr *Runner) server.HandlerConfig
	TestRequestHeaders(data harness.Data) map[string]string
	TestRequestPathParams(data harness.Data) map[string]string
	TestRequestQueryParams(data harness.Data) map[string]string
	TestRequestBody(data harness.Data) interface{}
	TestResponseCode() int
}

type TestCase struct {
	Name               string
	HandlerConfig      func(rnr *Runner) server.HandlerConfig
	RequestHeaders     func(data harness.Data) map[string]string
	RequestPathParams  func(data harness.Data) map[string]string
	RequestQueryParams func(data harness.Data) map[string]string
	RequestBody        func(data harness.Data) interface{}
	ResponseCode       int
}

//lint:ignore U1000 - testing struct implements interface
var _testCase TestCaser = &TestCase{}

func (t *TestCase) TestName() string {
	return t.Name
}

func (t *TestCase) TestHandlerConfig(rnr *Runner) server.HandlerConfig {
	return t.HandlerConfig(rnr)
}

func (t *TestCase) TestRequestHeaders(data harness.Data) map[string]string {
	return t.RequestHeaders(data)
}

func (t *TestCase) TestRequestPathParams(data harness.Data) map[string]string {
	pp := map[string]string{}
	if t.RequestPathParams != nil {
		pp = t.RequestPathParams(data)
	}
	return pp
}

func (t *TestCase) TestRequestQueryParams(data harness.Data) map[string]string {
	qp := map[string]string{}
	if t.RequestQueryParams != nil {
		qp = t.RequestQueryParams(data)
	}
	return qp
}

func (t *TestCase) TestRequestBody(data harness.Data) interface{} {
	var b interface{}
	if t.RequestBody != nil {
		b = t.RequestBody(data)
	}
	return b
}

func (t *TestCase) TestResponseCode() int {
	return t.ResponseCode
}

func RunTestCase(t *testing.T, th *harness.Testing, testCase TestCaser, testFunc func(method string, body io.Reader)) {
	rnr := NewRunner()

	err := rnr.Init(th.Config, th.Log, th.Store)
	require.NoError(t, err, "Runner init returns without error")

	err = th.Setup()
	require.NoError(t, err, "Test data setup returns without error")
	defer func() {
		err = th.Teardown()
		require.NoError(t, err, "Test data teardown returns without error")
	}()

	// config
	cfg := testCase.TestHandlerConfig(rnr)

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
	requestParams := testCase.TestRequestPathParams(th.Data)

	requestPath := cfg.Path
	for paramKey, paramValue := range requestParams {
		requestPath = strings.Replace(requestPath, paramKey, paramValue, 1)
	}

	// query params
	queryParams := testCase.TestRequestQueryParams(th.Data)

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
	data := testCase.TestRequestBody(th.Data)

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
	requestHeaders := testCase.TestRequestHeaders(th.Data)

	for headerKey, headerVal := range requestHeaders {
		req.Header.Add(headerKey, headerVal)
	}

	// recorder
	rec := httptest.NewRecorder()

	// serve
	rtr.ServeHTTP(rec, req)

	// test status
	require.Equalf(t, testCase.TestResponseCode(), rec.Code, "%s - Response code equals expected", testCase.TestName())

	testFunc(cfg.Method, rec.Body)
}
