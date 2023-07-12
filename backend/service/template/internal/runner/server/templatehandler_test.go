package runner

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/go-mud/backend/core/server"
	"gitlab.com/alienspaces/go-mud/backend/schema/template"
	"gitlab.com/alienspaces/go-mud/backend/service/template/internal/harness"
)

func Test_getTemplateHandler(t *testing.T) {

	// test harness
	th, err := newTestHarness()
	require.NoError(t, err, "New test data returns without error")

	_, err = th.Setup()
	require.NoError(t, err, "Test data setup returns without error")
	defer func() {
		err = th.Teardown()
		require.NoError(t, err, "Test data teardown returns without error")
	}()

	type testCase struct {
		TestCase
		expectResponse func(data harness.Data) *template.Response
	}

	testCaseResponseDecoder := func(body io.Reader) (interface{}, error) {
		var responseBody *template.Response
		err = json.NewDecoder(body).Decode(&responseBody)
		return responseBody, err
	}

	testCases := []testCase{
		{
			TestCase: TestCase{
				Name: "get existing",
				HandlerConfig: func(rnr *Runner) server.HandlerConfig {
					return rnr.HandlerConfig[getTemplate]
				},
				RequestPathParams: func(data harness.Data) map[string]string {
					params := map[string]string{
						":template_id": data.TemplateRecs[0].ID,
					}
					return params
				},
				ResponseCode:    http.StatusOK,
				ResponseDecoder: testCaseResponseDecoder,
			},
			expectResponse: func(data harness.Data) *template.Response {
				res := template.Response{
					Data: &template.Data{
						ID: data.TemplateRecs[0].ID,
					},
				}
				return &res
			},
		},
		{
			TestCase: TestCase{
				Name: "get non-existant",
				HandlerConfig: func(rnr *Runner) server.HandlerConfig {
					return rnr.HandlerConfig[getTemplate]
				},
				RequestPathParams: func(data harness.Data) map[string]string {
					params := map[string]string{
						":template_id": "17c19414-2d15-4d20-8fc3-36fc10341dc8",
					}
					return params
				},
				ResponseCode: http.StatusNotFound,
			},
		},
	}

	for _, tc := range testCases {
		t.Logf("Running test >%s<\n", tc.Name)

		t.Run(tc.Name, func(t *testing.T) {
			testFunc := func(method string, body interface{}) {
				if tc.ResponseCode != http.StatusOK && tc.ResponseCode != http.StatusCreated {
					return
				}
				if tc.expectResponse != nil {
					require.NotNil(t, body, "Response body is not nil")
					expectResponseBody := tc.expectResponse(th.Data)
					if expectResponseBody != nil {
						require.NotNil(t, expectResponseBody.Data, "Response data is not nil")
					}
				}
			}

			RunTestCase(t, th, &tc, testFunc)
		})
	}
}

func Test_postTemplatesHandler(t *testing.T) {

	// test harness
	th, err := newTestHarness()
	require.NoError(t, err, "New test data returns without error")

	_, err = th.Setup()
	require.NoError(t, err, "Test data setup returns without error")
	defer func() {
		err = th.Teardown()
		require.NoError(t, err, "Test data teardown returns without error")
	}()

	type testCase struct {
		TestCase
		expectResponse func(data harness.Data) *template.Response
	}

	testCaseResponseDecoder := func(body io.Reader) (interface{}, error) {
		var responseBody *template.Response
		err = json.NewDecoder(body).Decode(&responseBody)
		return responseBody, err
	}

	testCases := []testCase{
		{
			TestCase: TestCase{
				Name: "create without ID",
				HandlerConfig: func(rnr *Runner) server.HandlerConfig {
					return rnr.HandlerConfig[postTemplate]
				},
				RequestBody: func(data harness.Data) interface{} {
					req := template.Request{
						Data: &template.Data{},
					}
					return &req
				},
				ResponseCode:    http.StatusCreated,
				ResponseDecoder: testCaseResponseDecoder,
			},
		},
		{
			TestCase: TestCase{
				Name: "create with ID",
				HandlerConfig: func(rnr *Runner) server.HandlerConfig {
					return rnr.HandlerConfig[postTemplateWithID]
				},
				RequestPathParams: func(data harness.Data) map[string]string {
					params := map[string]string{
						":template_id": "e3a9e0f8-ce9c-477b-8b93-cf4da03af4c9",
					}
					return params
				},
				RequestBody: func(data harness.Data) interface{} {
					req := template.Request{
						Data: &template.Data{},
					}
					return &req
				},
				ResponseCode:    http.StatusCreated,
				ResponseDecoder: testCaseResponseDecoder,
			},
			expectResponse: func(data harness.Data) *template.Response {
				res := template.Response{
					Data: &template.Data{
						ID: "e3a9e0f8-ce9c-477b-8b93-cf4da03af4c9",
					},
				}
				return &res
			},
		},
	}

	for _, tc := range testCases {
		t.Logf("Running test >%s<\n", tc.Name)

		t.Run(tc.Name, func(t *testing.T) {
			testFunc := func(method string, body interface{}) {
				if tc.ResponseCode != http.StatusOK && tc.ResponseCode != http.StatusCreated {
					return
				}
				if tc.expectResponse != nil {
					require.NotNil(t, body, "Response body is not nil")
					expectResponseBody := tc.expectResponse(th.Data)
					if expectResponseBody != nil {
						require.NotNil(t, expectResponseBody.Data, "Response data is not nil")
					}
				}
			}

			RunTestCase(t, th, &tc, testFunc)
		})
	}
}

func Test_putTemplatesHandler(t *testing.T) {

	// test harness
	th, err := newTestHarness()
	require.NoError(t, err, "New test data returns without error")

	_, err = th.Setup()
	require.NoError(t, err, "Test data setup returns without error")
	defer func() {
		err = th.Teardown()
		require.NoError(t, err, "Test data teardown returns without error")
	}()

	type testCase struct {
		TestCase
		expectResponse func(data harness.Data) *template.Response
	}

	testCaseResponseDecoder := func(body io.Reader) (interface{}, error) {
		var responseBody *template.Response
		err = json.NewDecoder(body).Decode(&responseBody)
		return responseBody, err
	}

	testCases := []testCase{
		{
			TestCase: TestCase{
				Name: "update existing",
				HandlerConfig: func(rnr *Runner) server.HandlerConfig {
					return rnr.HandlerConfig[putTemplate]
				},
				RequestPathParams: func(data harness.Data) map[string]string {
					params := map[string]string{
						":template_id": data.TemplateRecs[0].ID,
					}
					return params
				},
				RequestBody: func(data harness.Data) interface{} {
					req := template.Request{
						Data: &template.Data{
							ID: data.TemplateRecs[0].ID,
						},
					}
					return &req
				},
				ResponseCode:    http.StatusOK,
				ResponseDecoder: testCaseResponseDecoder,
			},
			expectResponse: func(data harness.Data) *template.Response {
				res := template.Response{
					Data: &template.Data{
						ID: data.TemplateRecs[0].ID,
					},
				}
				return &res
			},
		},
		{
			TestCase: TestCase{
				Name: "update non-existing",
				HandlerConfig: func(rnr *Runner) server.HandlerConfig {
					return rnr.HandlerConfig[putTemplate]
				},
				RequestPathParams: func(data harness.Data) map[string]string {
					params := map[string]string{
						":template_id": "17c19414-2d15-4d20-8fc3-36fc10341dc8",
					}
					return params
				},
				RequestBody: func(data harness.Data) interface{} {
					req := template.Request{
						Data: &template.Data{
							ID: data.TemplateRecs[0].ID,
						},
					}
					return &req
				},
				ResponseCode: http.StatusNotFound,
			},
		},
		{
			TestCase: TestCase{
				Name: "update missing data",
				HandlerConfig: func(rnr *Runner) server.HandlerConfig {
					return rnr.HandlerConfig[putTemplate]
				},
				RequestBody: func(data harness.Data) interface{} {
					return nil
				},
				ResponseCode: http.StatusBadRequest,
			},
		},
	}

	for _, tc := range testCases {
		t.Logf("Running test >%s<\n", tc.Name)

		t.Run(tc.Name, func(t *testing.T) {
			testFunc := func(method string, body interface{}) {
				if tc.ResponseCode != http.StatusOK && tc.ResponseCode != http.StatusCreated {
					return
				}
				if tc.expectResponse != nil {
					require.NotNil(t, body, "Response body is not nil")
					expectResponseBody := tc.expectResponse(th.Data)
					if expectResponseBody != nil {
						require.NotNil(t, expectResponseBody.Data, "Response data is not nil")
					}
				}
			}

			RunTestCase(t, th, &tc, testFunc)
		})
	}
}
