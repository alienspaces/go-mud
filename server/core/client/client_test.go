package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/go-boilerplate/server/core/config"
	"gitlab.com/alienspaces/go-boilerplate/server/core/log"
	"gitlab.com/alienspaces/go-boilerplate/server/core/type/configurer"
	"gitlab.com/alienspaces/go-boilerplate/server/core/type/logger"
)

// NewDefaultDependencies -
func NewDefaultDependencies() (configurer.Configurer, logger.Logger, error) {

	// configurer
	c, err := config.NewConfig(nil, false)
	if err != nil {
		return nil, nil, err
	}

	configVars := []string{
		// general
		"APP_HOST",
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

func TestRetryRequest(t *testing.T) {

	c, l, err := NewDefaultDependencies()
	require.NoError(t, err, "NewDefaultDependencies returns without error")

	type RequestData struct {
		Request
		Name         string
		Strength     int
		Intelligence int
		Dexterity    int
	}

	tests := []struct {
		name        string
		method      string
		path        string
		params      map[string]string
		requestData *RequestData
		expectURL   string
		serverFunc  func(rw http.ResponseWriter, req *http.Request)
		expectErr   bool
	}{
		{
			name:   "Get resource OK",
			method: http.MethodGet,
			path:   "/api/games/:game_id/mages",
			params: map[string]string{
				"id":      "52fdfc07-2182-454f-963f-5f0f9a621d72",
				"game_id": "3fa1b1b7-cca9-435e-b2d6-a8c03be21bf1",
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
			path:   "/api/kobolds/:kobold_id",
			serverFunc: func(rw http.ResponseWriter, req *http.Request) {
				rw.WriteHeader(http.StatusBadRequest)
				return
			},
			expectErr: true,
		},
		{
			name:   "Post resource OK",
			method: http.MethodPost,
			path:   "/api/orcs/:orc_id",
			params: map[string]string{
				"orc_id": "52fdfc07-2182-454f-963f-5f0f9a621d72",
			},
			requestData: &RequestData{
				Name:     "Brain Basher",
				Strength: 10,
			},
			serverFunc: func(rw http.ResponseWriter, req *http.Request) {
				requestData := RequestData{}
				err := json.NewDecoder(req.Body).Decode(&requestData)
				if err != nil {
					rw.WriteHeader(http.StatusBadRequest)
					return
				}

				if requestData.Name != "Brain Basher" {
					rw.WriteHeader(http.StatusBadRequest)
					return
				}
				if requestData.Strength != 10 {
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
			path:   "/api/trolls/:troll_id",
			params: map[string]string{
				"troll_id": "52fdfc07-2182-454f-963f-5f0f9a621d72",
			},
			requestData: &RequestData{
				Name:     "Stone Fist",
				Strength: 20,
			},
			serverFunc: func(rw http.ResponseWriter, req *http.Request) {
				rw.WriteHeader(http.StatusBadRequest)
			},
			expectErr: true,
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

			// Host
			cl.Host = server.URL

			// set max retries to speed tests up
			cl.MaxRetries = 2

			resp := &Response{}
			err = cl.RetryRequest(tc.method, tc.path, tc.params, tc.requestData, resp)
			if tc.expectErr == true {
				require.Error(t, err, "RetryRequest returns with error")
				return
			}
			require.NoError(t, err, "RetryRequest returns without error")
			require.NotNil(t, resp, "RetryRequest returns a response")
		}()
	}
}

func TestBuildURL(t *testing.T) {

	c, l, err := NewDefaultDependencies()
	require.NoError(t, err, "NewDefaultDependencies returns without error")

	// Client
	cl, err := NewClient(c, l)

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
			name:   "Build URL with additional ID",
			method: http.MethodGet,
			path:   "/games/:game_id/mages",
			params: map[string]string{
				"id":      "52fdfc07-2182-454f-963f-5f0f9a621d72",
				"game_id": "3fa1b1b7-cca9-435e-b2d6-a8c03be21bf1",
			},
			expectErr: false,
			expectURL: "http://example.com/api/games/3fa1b1b7-cca9-435e-b2d6-a8c03be21bf1/mages/52fdfc07-2182-454f-963f-5f0f9a621d72",
		},
		{
			name:   "Build URL without additional ID",
			method: http.MethodGet,
			path:   "/games/:game_id/mages",
			params: map[string]string{
				"game_id": "3fa1b1b7-cca9-435e-b2d6-a8c03be21bf1",
			},
			expectErr: false,
			expectURL: "http://example.com/api/games/3fa1b1b7-cca9-435e-b2d6-a8c03be21bf1/mages",
		},
	}

	for _, tc := range tests {

		t.Logf("Running test >%s<", tc.name)

		func() {

			url, err := cl.BuildURL(tc.method, tc.path, tc.params)
			if tc.expectErr == true {
				require.Error(t, err, "BuildURL returns with error")
				return
			}
			t.Logf("Resulting URL >%s<", url)
			require.NoError(t, err, "BuildURL returns without error")
			require.Equal(t, tc.expectURL, url, "BuildURL returns expected URL")
		}()
	}
}
