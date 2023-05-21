package runner

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/go-mud/backend/core/server"
	"gitlab.com/alienspaces/go-mud/backend/schema/template"
	"gitlab.com/alienspaces/go-mud/backend/service/template/internal/harness"
)

func TestGetTemplate(t *testing.T) {

	// test dependencies
	c, l, s, err := newDefaultDependencies()
	require.NoError(t, err, "NewDefaultDependencies returns without error")

	// test harness
	th, err := newTestHarness(c, l, s, nil)
	require.NoError(t, err, "New test data returns without error")

	defer func() {
		db, err := s.GetDb()
		require.NoError(t, err, "getDb should return no error")

		err = db.Close()
		require.NoError(t, err, "close db should return no error")
	}()

	type testCase struct {
		TestCase
		expectResponse func(data *harness.Data) *template.Response
	}

	testCases := []testCase{
		{
			TestCase: TestCase{
				Name: "GET - Get existing",
				HandlerConfig: func(rnr *Runner) server.HandlerConfig {
					return rnr.HandlerConfig[HandlerGetOneTemplate]
				},
				RequestPathParams: func(data *harness.Data) map[string]string {
					params := map[string]string{
						":template_id": data.TemplateRecs[0].ID,
					}
					return params
				},
				ResponseCode: http.StatusOK,
				ResponseBody: testCaseResponseBodyGeneric[template.Response],
			},
			expectResponse: func(data *harness.Data) *template.Response {
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
				Name: "GET - Get non-existant",
				HandlerConfig: func(rnr *Runner) server.HandlerConfig {
					return rnr.HandlerConfig[HandlerGetOneTemplate]
				},
				RequestPathParams: func(data *harness.Data) map[string]string {
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
					er := tc.expectResponse(th.Data)
					if er != nil {
						ar := body.(template.Response)
						require.NotNil(t, ar.Data, "Response data is not nil")
					}
				}
			}

			RunTestCase(t, th, &tc, testFunc)
		})
	}
}

func TestCreateTemplate(t *testing.T) {

	// test dependencies
	c, l, s, err := newDefaultDependencies()
	require.NoError(t, err, "NewDefaultDependencies returns without error")

	// test harness
	th, err := newTestHarness(c, l, s, nil)
	require.NoError(t, err, "New test data returns without error")

	defer func() {
		db, err := s.GetDb()
		require.NoError(t, err, "getDb should return no error")

		err = db.Close()
		require.NoError(t, err, "close db should return no error")
	}()

	type testCase struct {
		TestCase
		expectResponse func(data *harness.Data) *template.Response
	}

	testCases := []testCase{
		{
			TestCase: TestCase{
				Name: "POST - Create without ID",
				HandlerConfig: func(rnr *Runner) server.HandlerConfig {
					return rnr.HandlerConfig[HandlerCreateOneTemplate]
				},
				RequestBody: func(data *harness.Data) interface{} {
					req := template.Request{
						Data: &template.Data{},
					}
					return &req
				},
				ResponseCode: http.StatusCreated,
				ResponseBody: testCaseResponseBodyGeneric[template.Response],
			},
		},
		{
			TestCase: TestCase{
				Name: "POST - Create with ID",
				HandlerConfig: func(rnr *Runner) server.HandlerConfig {
					return rnr.HandlerConfig[HandlerCreateOneTemplateWithID]
				},
				RequestPathParams: func(data *harness.Data) map[string]string {
					params := map[string]string{
						":template_id": "e3a9e0f8-ce9c-477b-8b93-cf4da03af4c9",
					}
					return params
				},
				RequestBody: func(data *harness.Data) interface{} {
					req := template.Request{
						Data: &template.Data{},
					}
					return &req
				},
				ResponseCode: http.StatusCreated,
				ResponseBody: testCaseResponseBodyGeneric[template.Response],
			},
			expectResponse: func(data *harness.Data) *template.Response {
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
					er := tc.expectResponse(th.Data)
					if er != nil {
						ar := body.(template.Response)
						require.NotNil(t, ar.Data, "Response data is not nil")
					}
				}
			}

			RunTestCase(t, th, &tc, testFunc)
		})
	}
}

func TestUpdateTemplate(t *testing.T) {

	// test dependencies
	c, l, s, err := newDefaultDependencies()
	require.NoError(t, err, "NewDefaultDependencies returns without error")

	// test harness
	th, err := newTestHarness(c, l, s, nil)
	require.NoError(t, err, "New test data returns without error")

	defer func() {
		db, err := s.GetDb()
		require.NoError(t, err, "getDb should return no error")

		err = db.Close()
		require.NoError(t, err, "close db should return no error")
	}()

	type testCase struct {
		TestCase
		expectResponse func(data *harness.Data) *template.Response
	}

	testCases := []testCase{
		{
			TestCase: TestCase{
				Name: "PUT - Update existing",
				HandlerConfig: func(rnr *Runner) server.HandlerConfig {
					return rnr.HandlerConfig[HandlerUpdateOneTemplate]
				},
				RequestPathParams: func(data *harness.Data) map[string]string {
					params := map[string]string{
						":template_id": data.TemplateRecs[0].ID,
					}
					return params
				},
				RequestBody: func(data *harness.Data) interface{} {
					req := template.Request{
						Data: &template.Data{
							ID: data.TemplateRecs[0].ID,
						},
					}
					return &req
				},
				ResponseCode: http.StatusOK,
				ResponseBody: testCaseResponseBodyGeneric[template.Response],
			},
			expectResponse: func(data *harness.Data) *template.Response {
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
				Name: "PUT - Update non-existing",
				HandlerConfig: func(rnr *Runner) server.HandlerConfig {
					return rnr.HandlerConfig[HandlerUpdateOneTemplate]
				},
				RequestPathParams: func(data *harness.Data) map[string]string {
					params := map[string]string{
						":template_id": "17c19414-2d15-4d20-8fc3-36fc10341dc8",
					}
					return params
				},
				RequestBody: func(data *harness.Data) interface{} {
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
				Name: "PUT - Update missing data",
				HandlerConfig: func(rnr *Runner) server.HandlerConfig {
					return rnr.HandlerConfig[HandlerUpdateOneTemplate]
				},
				RequestBody: func(data *harness.Data) interface{} {
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
					er := tc.expectResponse(th.Data)
					if er != nil {
						ar := body.(template.Response)
						require.NotNil(t, ar.Data, "Response data is not nil")
					}
				}
			}

			RunTestCase(t, th, &tc, testFunc)
		})
	}
}
