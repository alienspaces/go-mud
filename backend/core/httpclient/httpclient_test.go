package httpclient

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/go-mud/backend/core/config"
	"gitlab.com/alienspaces/go-mud/backend/core/log"
	"gitlab.com/alienspaces/go-mud/backend/core/type/logger"
)

type requestData struct {
	Request
	Name string
	Age  int
}

type clientTestCase struct {
	name        string
	method      string
	path        string
	params      map[string]string
	requestData *requestData
	serverFunc  func(rw http.ResponseWriter, req *http.Request)
	expectErr   bool
}

// NewDefaultDependencies -
func NewDefaultDependencies() (logger.Logger, error) {

	c, err := config.NewConfig(nil, false)
	if err != nil {
		return nil, err
	}

	l, err := log.NewLogger(c)
	if err != nil {
		return nil, err
	}

	return l, nil
}

func TestRetryRequest(t *testing.T) {

	l, err := NewDefaultDependencies()
	require.NoError(t, err, "NewDefaultDependencies returns without error")

	tests := []clientTestCase{
		{
			name:   "Get resource OK",
			method: http.MethodGet,
			path:   "/api/collections/:collection_id/members",
			params: map[string]string{
				"id":            "52fdfc07-2182-454f-963f-5f0f9a621d72",
				"collection_id": "3fa1b1b7-cca9-435e-b2d6-a8c03be21bf1",
			},
			serverFunc: func(rw http.ResponseWriter, req *http.Request) {
				respData, err := json.Marshal(&Response{})
				if err != nil {
					l.Warn("Failed encoding data >%v<", err)
					rw.WriteHeader(http.StatusInternalServerError)
					return
				}
				rw.WriteHeader(http.StatusOK)
				rw.Write(respData)
			},
			expectErr: false,
		},
		{
			name:   "Get resource BadRequest",
			method: http.MethodGet,
			path:   "/api/collections/:collection_id",
			serverFunc: func(rw http.ResponseWriter, req *http.Request) {
				rw.WriteHeader(http.StatusBadRequest)
			},
			expectErr: true,
		},
		{
			name:   "Post resource OK",
			method: http.MethodPost,
			path:   "/api/collections/:collection_id",
			params: map[string]string{
				"collection_id": "52fdfc07-2182-454f-963f-5f0f9a621d72",
			},
			requestData: &requestData{
				Name: "John",
				Age:  10,
			},
			serverFunc: func(rw http.ResponseWriter, req *http.Request) {
				requestData := requestData{}
				err := json.NewDecoder(req.Body).Decode(&requestData)
				if err != nil {
					rw.WriteHeader(http.StatusBadRequest)
					return
				}

				if requestData.Name != "John" {
					rw.WriteHeader(http.StatusBadRequest)
					return
				}
				if requestData.Age != 10 {
					rw.WriteHeader(http.StatusBadRequest)
					return
				}

				respData, err := json.Marshal(&Response{})
				if err != nil {
					l.Warn("Failed encoding data >%v<", err)
					rw.WriteHeader(http.StatusInternalServerError)
					return
				}
				rw.WriteHeader(http.StatusOK)
				rw.Write(respData)
			},
			expectErr: false,
		},
		{
			name:   "Post resource BadRequest",
			method: http.MethodPost,
			path:   "/api/collections/:collection_id",
			params: map[string]string{
				"collection_id": "52fdfc07-2182-454f-963f-5f0f9a621d72",
			},
			requestData: &requestData{
				Name: "Mary",
				Age:  20,
			},
			serverFunc: func(rw http.ResponseWriter, req *http.Request) {
				rw.WriteHeader(http.StatusBadRequest)
			},
			expectErr: true,
		},
		{
			name:   "Post resource error - missing request data",
			method: http.MethodPost,
			path:   "/api/collections/:collection_id",
			params: map[string]string{
				"collection_id": "52fdfc07-2182-454f-963f-5f0f9a621d72",
			},
			requestData: nil,
			serverFunc: func(rw http.ResponseWriter, req *http.Request) {
				rw.WriteHeader(http.StatusBadRequest)
			},
			expectErr: true,
		},
		{
			name:   "Put resource OK",
			method: http.MethodPut,
			path:   "/api/collections/:collection_id",
			params: map[string]string{
				"collection_id": "52fdfc07-2182-454f-963f-5f0f9a621d72",
			},
			requestData: &requestData{
				Name: "John",
				Age:  10,
			},
			serverFunc: func(rw http.ResponseWriter, req *http.Request) {
				requestData := requestData{}
				err := json.NewDecoder(req.Body).Decode(&requestData)
				if err != nil {
					rw.WriteHeader(http.StatusBadRequest)
					return
				}

				id := strings.TrimPrefix(req.URL.Path, "/api/collections/")
				if id != "52fdfc07-2182-454f-963f-5f0f9a621d72" {
					rw.WriteHeader(http.StatusNotFound)
					return
				}

				if requestData.Name != "John" {
					rw.WriteHeader(http.StatusBadRequest)
					return
				}
				if requestData.Age != 10 {
					rw.WriteHeader(http.StatusBadRequest)
					return
				}

				respData, err := json.Marshal(&Response{})
				if err != nil {
					l.Warn("Failed encoding data >%v<", err)
					rw.WriteHeader(http.StatusInternalServerError)
					return
				}
				rw.WriteHeader(http.StatusOK)
				rw.Write(respData)
			},
			expectErr: false,
		},
		{
			name:   "Put resource NotFound",
			method: http.MethodPut,
			path:   "/api/collections",
			params: map[string]string{
				"id": "52fdfc07-2182-454f-963f-5f0f9a621d72",
			},
			requestData: &requestData{
				Name: "John",
				Age:  10,
			},
			serverFunc: func(rw http.ResponseWriter, req *http.Request) {
				rw.WriteHeader(http.StatusNotFound)
			},
			expectErr: true,
		},
		{
			name:   "Put resource error - missing request data",
			method: http.MethodPut,
			path:   "/api/collections",
			params: map[string]string{
				"id": "52fdfc07-2182-454f-963f-5f0f9a621d72",
			},
			serverFunc: func(rw http.ResponseWriter, req *http.Request) {
				rw.WriteHeader(http.StatusBadRequest)
			},
			expectErr: true,
		},
		// TODO delete method not yet supported in core/client/client.go
		{
			name:   "Delete resource OK",
			method: http.MethodDelete,
			path:   "/api/collections/:collection_id",
			params: map[string]string{
				"collection_id": "52fdfc07-2182-454f-963f-5f0f9a621d72",
			},
			requestData: &requestData{
				Name: "John",
				Age:  10,
			},
			serverFunc: func(rw http.ResponseWriter, req *http.Request) {
				id := strings.TrimPrefix(req.URL.Path, "/api/collections/")
				if id != "52fdfc07-2182-454f-963f-5f0f9a621d72" {
					rw.WriteHeader(http.StatusNotFound)
					return
				}
				rw.WriteHeader(http.StatusOK)
			},
			expectErr: false,
		},
		{
			name:   "Delete resource NotFound",
			method: http.MethodDelete,
			path:   "/api/collections/:collection_id",
			params: map[string]string{
				"collection_id": "52fdfc07-abcd-abcd-abcde-5f0f9a621d72",
			},
			serverFunc: func(rw http.ResponseWriter, req *http.Request) {
				rw.WriteHeader(http.StatusNotFound)
			},
			expectErr: true,
		},
	}

	testRequest := func(t *testing.T, tc clientTestCase, methodName string, resp *Response, err error) {
		if tc.expectErr == true {
			require.Errorf(t, err, "%s returns with error", methodName)
			return
		}
		require.NoErrorf(t, err, "%s returns without error", methodName)
		require.NotNilf(t, resp, "%s returns a response", methodName)
	}

	for _, tc := range tests {

		t.Logf("Running test >%s<", tc.name)

		t.Run(tc.name, func(t *testing.T) {
			// Test HTTP server
			server := httptest.NewServer(http.HandlerFunc(tc.serverFunc))
			defer server.Close()

			// Client
			cl, err := NewClient(l)
			require.NoError(t, err, "NewClient returns without error")
			require.NotNil(t, cl, "NewClient returns a client")

			// Host
			cl.Host = server.URL

			// set max retries to speed tests up
			cl.MaxRetries = 2

			resp := &Response{}
			err = cl.RetryRequest(tc.method, tc.path, tc.params, tc.requestData, resp)
			testRequest(t, tc, "RetryRequest", resp, err)

			if tc.method == http.MethodGet {
				err = cl.Get(tc.path, tc.params, resp)
				testRequest(t, tc, "Get", resp, err)
			} else if tc.method == http.MethodPost {
				if tc.requestData != nil {
					err = cl.Create(tc.path, tc.params, tc.requestData, resp)
				} else {
					// for reqData to be in nil in Create, the interface type must be null
					err = cl.Create(tc.path, tc.params, nil, resp)
				}
				testRequest(t, tc, "Post", resp, err)
			} else if tc.method == http.MethodPut {
				if tc.requestData != nil {
					_ = cl.Update(tc.path, tc.params, tc.requestData, resp)
				} else {
					// for reqData to be in nil in Update, the interface type must be null
					_ = cl.Update(tc.path, tc.params, nil, resp)
				}
				err = cl.Update(tc.path, tc.params, tc.requestData, resp)
				testRequest(t, tc, "Update", resp, err)
			} else if tc.method == http.MethodDelete {
				err = cl.Delete(tc.path, tc.params, resp)
				testRequest(t, tc, "Delete", resp, err)
			}
		})
	}
}

func Test_buildURL(t *testing.T) {

	l, err := NewDefaultDependencies()
	require.NoError(t, err, "NewDefaultDependencies returns without error")

	// Client
	cl, err := NewClient(l)

	// Host
	cl.Host = "http://example.com"

	require.NoError(t, err, "NewClient returns without error")
	require.NotNil(t, cl, "NewClient returns a client")

	// Set base path
	cl.Path = "/api"

	tests := []struct {
		name      string
		method    string
		path      string
		params    map[string]string
		expectErr bool
		expectURL string
	}{
		{
			name:   "Build GET URL",
			method: http.MethodGet,
			path:   "/collections/:collection_id/members/:member_id",
			params: map[string]string{
				"collection_id": "3fa1b1b7-cca9-435e-b2d6-a8c03be21bf1",
				"member_id":     "52fdfc07-2182-454f-963f-5f0f9a621d72",
			},
			expectErr: false,
			expectURL: "http://example.com/api/collections/3fa1b1b7-cca9-435e-b2d6-a8c03be21bf1/members/52fdfc07-2182-454f-963f-5f0f9a621d72",
		},
		{
			name:   "Build GET URL with path params without colon and path params without colon",
			method: http.MethodGet,
			path:   "/collections/collection_id/members/member_id",
			params: map[string]string{
				"collection_id": "3fa1b1b7-cca9-435e-b2d6-a8c03be21bf1",
				"member_id":     "52fdfc07-2182-454f-963f-5f0f9a621d72",
			},
			expectErr: false,
			expectURL: "http://example.com/api/collections/3fa1b1b7-cca9-435e-b2d6-a8c03be21bf1/members/52fdfc07-2182-454f-963f-5f0f9a621d72",
		},
		{
			name:   "Build GET URL with query params",
			method: http.MethodGet,
			path:   "/collections/:collection_id/members",
			params: map[string]string{
				"collection_id": "3fa1b1b7-cca9-435e-b2d6-a8c03be21bf1",
				"qp1":           "1",
			},
			expectErr: false,
			expectURL: "http://example.com/api/collections/3fa1b1b7-cca9-435e-b2d6-a8c03be21bf1/members?qp1=1",
		},
		{
			name:   "Build POST URL",
			method: http.MethodPost,
			path:   "/collections/:collection_id/members",
			params: map[string]string{
				"collection_id": "3fa1b1b7-cca9-435e-b2d6-a8c03be21bf1",
			},
			expectErr: false,
			expectURL: "http://example.com/api/collections/3fa1b1b7-cca9-435e-b2d6-a8c03be21bf1/members",
		},
		{
			name:   "Build PUT URL",
			method: http.MethodPut,
			path:   "/collections/:collection_id/members",
			params: map[string]string{
				"collection_id": "3fa1b1b7-cca9-435e-b2d6-a8c03be21bf1",
			},
			expectErr: false,
			expectURL: "http://example.com/api/collections/3fa1b1b7-cca9-435e-b2d6-a8c03be21bf1/members",
		},
		{
			name:   "Build DELETE URL",
			method: http.MethodDelete,
			path:   "/collections/:collection_id/members",
			params: map[string]string{
				"collection_id": "3fa1b1b7-cca9-435e-b2d6-a8c03be21bf1",
			},
			expectErr: false,
			expectURL: "http://example.com/api/collections/3fa1b1b7-cca9-435e-b2d6-a8c03be21bf1/members",
		},
		{
			name:   "Build GET URL with empty params",
			method: http.MethodGet,
			path:   "/collections/:collection_id/members",
			params: map[string]string{
				"collection_id": "",
			},
			expectErr: true,
		},
		// prefixed with colon
		{
			name:   "Build GET URL prefixed with colon",
			method: http.MethodGet,
			path:   "/collections/:collection_id/members/:member_id",
			params: map[string]string{
				":collection_id": "3fa1b1b7-cca9-435e-b2d6-a8c03be21bf1",
				":member_id":     "52fdfc07-2182-454f-963f-5f0f9a621d72",
			},
			expectErr: false,
			expectURL: "http://example.com/api/collections/3fa1b1b7-cca9-435e-b2d6-a8c03be21bf1/members/52fdfc07-2182-454f-963f-5f0f9a621d72",
		},
		{
			name:   "Build GET URL with path params without colon and params prefixed with colon",
			method: http.MethodGet,
			path:   "/collections/collection_id/members/member_id",
			params: map[string]string{
				":collection_id": "3fa1b1b7-cca9-435e-b2d6-a8c03be21bf1",
				":member_id":     "52fdfc07-2182-454f-963f-5f0f9a621d72",
			},
			expectErr: false,
			expectURL: "http://example.com/api/collections/3fa1b1b7-cca9-435e-b2d6-a8c03be21bf1/members/52fdfc07-2182-454f-963f-5f0f9a621d72",
		},
		{
			name:   "Build GET URL with query params prefixed with colon",
			method: http.MethodGet,
			path:   "/collections/:collection_id/members",
			params: map[string]string{
				":collection_id": "3fa1b1b7-cca9-435e-b2d6-a8c03be21bf1",
				":qp1":           "1",
			},
			expectErr: false,
			expectURL: "http://example.com/api/collections/3fa1b1b7-cca9-435e-b2d6-a8c03be21bf1/members?qp1=1",
		},
		{
			name:   "Build POST URL prefixed with colon",
			method: http.MethodPost,
			path:   "/collections/:collection_id/members",
			params: map[string]string{
				":collection_id": "3fa1b1b7-cca9-435e-b2d6-a8c03be21bf1",
			},
			expectErr: false,
			expectURL: "http://example.com/api/collections/3fa1b1b7-cca9-435e-b2d6-a8c03be21bf1/members",
		},
		{
			name:   "Build PUT URL prefixed with colon",
			method: http.MethodPut,
			path:   "/collections/:collection_id/members",
			params: map[string]string{
				":collection_id": "3fa1b1b7-cca9-435e-b2d6-a8c03be21bf1",
			},
			expectErr: false,
			expectURL: "http://example.com/api/collections/3fa1b1b7-cca9-435e-b2d6-a8c03be21bf1/members",
		},
		{
			name:   "Build DELETE URL prefixed with colon",
			method: http.MethodDelete,
			path:   "/collections/:collection_id/members",
			params: map[string]string{
				":collection_id": "3fa1b1b7-cca9-435e-b2d6-a8c03be21bf1",
			},
			expectErr: false,
			expectURL: "http://example.com/api/collections/3fa1b1b7-cca9-435e-b2d6-a8c03be21bf1/members",
		},
		{
			name:   "Build GET URL with empty params prefixed with colon",
			method: http.MethodGet,
			path:   "/collections/:collection_id/members",
			params: map[string]string{
				":collection_id": "",
			},
			expectErr: true,
		},
	}

	for _, tc := range tests {

		t.Logf("Running test >%s<", tc.name)

		t.Run(tc.name, func(t *testing.T) {

			url, err := cl.buildURL(tc.method, tc.path, tc.params)
			if tc.expectErr == true {
				require.Error(t, err, "buildURL returns with error")
				return
			}
			t.Logf("Resulting URL >%s<", url)
			require.NoError(t, err, "buildURL returns without error")
			require.Equal(t, tc.expectURL, url, "buildURL returns expected URL")
		})
	}
}
