package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/go-boilerplate/server/core/config"
	"gitlab.com/alienspaces/go-boilerplate/server/core/log"
	"gitlab.com/alienspaces/go-boilerplate/server/core/type/configurer"
	"gitlab.com/alienspaces/go-boilerplate/server/core/type/logger"
)

// Common test handler function
var handlerFunc = func(rw http.ResponseWriter, req *http.Request) {

	// request data
	reqData := Request{}

	// close before returning
	defer req.Body.Close()

	err := json.NewDecoder(req.Body).Decode(&reqData)
	if err != nil && err.Error() != "EOF" {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	// response data
	respData, err := json.Marshal(&Response{})
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write(respData)
}

// NewDefaultDependencies -
func NewDefaultDependencies() (configurer.Configurer, logger.Logger, error) {

	// configurer
	c, err := config.NewConfig(nil, false)
	if err != nil {
		return nil, nil, err
	}

	configVars := []string{
		// general
		"APP_SERVER_HOST",
		// logger
		"APP_SERVER_LOG_LEVEL",
	}
	for _, key := range configVars {
		err = c.Add(key, false)
		if err != nil {
			return nil, nil, err
		}
	}

	// logger
	l, err := log.NewLogger(c)
	if err != nil {
		return nil, nil, err
	}

	return c, l, nil
}

func TestGetTemplate(t *testing.T) {

	c, l, err := NewDefaultDependencies()
	require.NoError(t, err, "NewDefaultDependencies returns without error")

	tests := []struct {
		name       string
		id         string
		serverFunc func(rw http.ResponseWriter, req *http.Request)
		expectErr  bool
	}{
		{
			name:       "Get resource success",
			id:         "bceac4ab-738a-4e62-a040-835e6fab331f",
			serverFunc: handlerFunc,
			expectErr:  false,
		},
		{
			name:       "Get resource not success",
			id:         "",
			serverFunc: handlerFunc,
			expectErr:  true,
		},
	}

	for _, tc := range tests {

		t.Logf("Running test >%s<", tc.name)

		func() {
			// Test HTTP server
			server := httptest.NewServer(http.HandlerFunc(tc.serverFunc))
			defer server.Close()

			// Client
			cl, err := NewClient(c, l)
			require.NoError(t, err, "NewClient returns without error")
			require.NotNil(t, cl, "NewClient returns a client")

			// Set host
			cl.Host = server.URL

			// Set max retries to speed up tests
			cl.MaxRetries = 2

			resp, err := cl.GetTemplate(tc.id)
			if tc.expectErr == true {
				require.Error(t, err, "GetTemplate returns with error")
				return
			}
			require.NoError(t, err, "GetTemplate returns without error")
			require.NotNil(t, resp, "GetTemplate returns a response")
		}()
	}
}

func TestGetTemplates(t *testing.T) {

	c, l, err := NewDefaultDependencies()
	require.NoError(t, err, "NewDefaultDependencies returns without error")

	tests := []struct {
		name       string
		params     map[string]string
		serverFunc func(rw http.ResponseWriter, req *http.Request)
		expectErr  bool
	}{
		{
			name: "Get resources with ID success",
			params: map[string]string{
				"id": "bceac4ab-738a-4e62-a040-835e6fab331f",
			},
			serverFunc: handlerFunc,
			expectErr:  false,
		},
		{
			name: "Get resources with params success",
			params: map[string]string{
				"blah": "bceac4ab-738a-4e62-a040-835e6fab331f",
			},
			serverFunc: handlerFunc,
			expectErr:  false,
		},
	}

	for _, tc := range tests {

		t.Logf("Running test >%s<", tc.name)

		func() {
			// Test HTTP server
			server := httptest.NewServer(http.HandlerFunc(tc.serverFunc))
			defer server.Close()

			// Client
			cl, err := NewClient(c, l)
			require.NoError(t, err, "NewClient returns without error")
			require.NotNil(t, cl, "NewClient returns a client")

			// Set host
			cl.Host = server.URL

			// Set max retries to speed up tests
			cl.MaxRetries = 2

			resp, err := cl.GetTemplates(tc.params)
			if tc.expectErr == true {
				require.Error(t, err, "GetTemplates returns with error")
				return
			}
			require.NoError(t, err, "GetTemplates returns without error")
			require.NotNil(t, resp, "GetTemplates returns a response")
		}()
	}
}

func TestCreateTemplate(t *testing.T) {

	c, l, err := NewDefaultDependencies()
	require.NoError(t, err, "NewDefaultDependencies returns without error")

	tests := []struct {
		name        string
		id          string
		requestData *Request
		serverFunc  func(rw http.ResponseWriter, req *http.Request)
		expectErr   bool
	}{
		{
			name: "Create resource with ID success",
			id:   gofakeit.UUID(),
			requestData: &Request{
				Data: Data{},
			},
			serverFunc: handlerFunc,
			expectErr:  false,
		},
		{
			name: "Create resource without ID success",
			requestData: &Request{
				Data: Data{},
			},
			serverFunc: handlerFunc,
			expectErr:  false,
		},
	}

	for _, tc := range tests {

		t.Logf("Running test >%s<", tc.name)

		func() {
			// Test HTTP server
			server := httptest.NewServer(http.HandlerFunc(tc.serverFunc))
			defer server.Close()

			// Client
			cl, err := NewClient(c, l)
			require.NoError(t, err, "NewClient returns without error")
			require.NotNil(t, cl, "NewClient returns a client")

			// Set host
			cl.Host = server.URL

			// Set max retries to speed up tests
			cl.MaxRetries = 2

			resp, err := cl.CreateTemplate(tc.id, tc.requestData)
			if tc.expectErr == true {
				require.Error(t, err, "CreateTemplate returns with error")
				return
			}
			require.NoError(t, err, "CreateTemplate returns without error")
			require.NotNil(t, resp, "CreateTemplate returns a response")
		}()
	}
}

func TestUpdateTemplate(t *testing.T) {

	c, l, err := NewDefaultDependencies()
	require.NoError(t, err, "NewDefaultDependencies returns without error")

	tests := []struct {
		name        string
		id          string
		requestData *Request
		serverFunc  func(rw http.ResponseWriter, req *http.Request)
		expectErr   bool
	}{
		{
			name: "Update resource with ID success",
			id:   gofakeit.UUID(),
			requestData: &Request{
				Data: Data{},
			},
			serverFunc: handlerFunc,
			expectErr:  false,
		},
		{
			name: "Update resource without ID fail",
			requestData: &Request{
				Data: Data{},
			},
			serverFunc: handlerFunc,
			expectErr:  true,
		},
	}

	for _, tc := range tests {

		t.Logf("Running test >%s<", tc.name)

		func() {
			// Test HTTP server
			server := httptest.NewServer(http.HandlerFunc(tc.serverFunc))
			defer server.Close()

			// Client
			cl, err := NewClient(c, l)
			require.NoError(t, err, "NewClient returns without error")
			require.NotNil(t, cl, "NewClient returns a client")

			// Set host
			cl.Host = server.URL

			// Set max retries to speed up tests
			cl.MaxRetries = 2

			resp, err := cl.UpdateTemplate(tc.id, tc.requestData)
			if tc.expectErr == true {
				require.Error(t, err, "UpdateTemplate returns with error")
				return
			}
			require.NoError(t, err, "UpdateTemplate returns without error")
			require.NotNil(t, resp, "UpdateTemplate returns a response")
		}()
	}
}
