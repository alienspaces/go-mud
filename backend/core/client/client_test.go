package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/go-mud/backend/core/config"
	"gitlab.com/alienspaces/go-mud/backend/core/log"
	"gitlab.com/alienspaces/go-mud/backend/core/type/configurer"
	"gitlab.com/alienspaces/go-mud/backend/core/type/logger"
)

type requestData struct {
	Request
	Name         string
	Strength     int
	Intelligence int
	Dexterity    int
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
func NewDefaultDependencies() (configurer.Configurer, logger.Logger, error) {

	c, err := config.NewConfigWithDefaults(nil, false)
	if err != nil {
		return nil, nil, err
	}

	l := log.NewLogger(c)

	return c, l, nil
}

func TestRetryRequest(t *testing.T) {

	c, l, err := NewDefaultDependencies()
	require.NoError(t, err, "NewDefaultDependencies returns without error")

	tests := []clientTestCase{
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
			requestData: &requestData{
				Name:     "Brain Basher",
				Strength: 10,
			},
			serverFunc: func(rw http.ResponseWriter, req *http.Request) {
				requestData := requestData{}
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
			requestData: &requestData{
				Name:     "Stone Fist",
				Strength: 20,
			},
			serverFunc: func(rw http.ResponseWriter, req *http.Request) {
				rw.WriteHeader(http.StatusBadRequest)
			},
			expectErr: true,
		},
		{
			name:   "Post resource error - missing request data",
			method: http.MethodPost,
			path:   "/api/trolls/:troll_id",
			params: map[string]string{
				"troll_id": "52fdfc07-2182-454f-963f-5f0f9a621d72",
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
			path:   "/api/orcs",
			params: map[string]string{
				"id": "52fdfc07-2182-454f-963f-5f0f9a621d72",
			},
			requestData: &requestData{
				Name:     "Brain Basher",
				Strength: 10,
			},
			serverFunc: func(rw http.ResponseWriter, req *http.Request) {
				requestData := requestData{}
				err := json.NewDecoder(req.Body).Decode(&requestData)
				if err != nil {
					rw.WriteHeader(http.StatusBadRequest)
					return
				}

				id := strings.TrimPrefix(req.URL.Path, "/api/orcs/")
				if id != "52fdfc07-2182-454f-963f-5f0f9a621d72" {
					rw.WriteHeader(http.StatusNotFound)
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
			name:   "Put resource NotFound",
			method: http.MethodPut,
			path:   "/api/orcs",
			params: map[string]string{
				"id": "52fdfc07-2182-454f-963f-5f0f9a621d72",
			},
			requestData: &requestData{
				Name:     "Brain Basher",
				Strength: 10,
			},
			serverFunc: func(rw http.ResponseWriter, req *http.Request) {
				rw.WriteHeader(http.StatusNotFound)
			},
			expectErr: true,
		},
		{
			name:   "Put resource error - missing request data",
			method: http.MethodPut,
			path:   "/api/orcs",
			params: map[string]string{
				"id": "52fdfc07-2182-454f-963f-5f0f9a621d72",
			},
			serverFunc: func(rw http.ResponseWriter, req *http.Request) {
				rw.WriteHeader(http.StatusBadRequest)
			},
			expectErr: true,
		},
		// TODO delete method not yet supported in core/client/client.go
		//{
		//	name:   "Delete resource OK",
		//	method: http.MethodDelete,
		//	path:   "/api/orcs",
		//	params: map[string]string{
		//		"id": "52fdfc07-2182-454f-963f-5f0f9a621d72",
		//	},
		//	requestData: &requestData{
		//		Name:     "Brain Basher",
		//		Strength: 10,
		//	},
		//	serverFunc: func(rw http.ResponseWriter, req *http.Request) {
		//		id := strings.TrimPrefix(req.URL.Path, "/api/orcs/")
		//		if id != "52fdfc07-2182-454f-963f-5f0f9a621d72" {
		//			rw.WriteHeader(http.StatusNotFound)
		//			return
		//		}
		//		rw.WriteHeader(http.StatusOK)
		//	},
		//	expectErr: false,
		//},
		{
			name:   "Delete resource NotFound",
			method: http.MethodDelete,
			path:   "/api/orcs",
			params: map[string]string{
				"id": "52fdfc07-abcd-abcd-abcde-5f0f9a621d72",
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
			cl, err := NewClient(c, l)
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
			name:   "Build GET URL with additional ID",
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
			name:   "Build GET URL with additional ID and path param without colon",
			method: http.MethodGet,
			path:   "/games/game_id/mages/test/test_id",
			params: map[string]string{
				"id":      "52fdfc07-2182-454f-963f-5f0f9a621d72",
				"game_id": "3fa1b1b7-cca9-435e-b2d6-a8c03be21bf1",
				"test_id": "52fdfc07-2182-454f-963f-5f0f9a621abc",
			},
			expectErr: false,
			expectURL: "http://example.com/api/games/3fa1b1b7-cca9-435e-b2d6-a8c03be21bf1/mages/test/52fdfc07-2182-454f-963f-5f0f9a621abc/52fdfc07-2182-454f-963f-5f0f9a621d72",
		},
		{
			name:   "Build GET URL without additional ID",
			method: http.MethodGet,
			path:   "/games/:game_id/mages",
			params: map[string]string{
				"game_id": "3fa1b1b7-cca9-435e-b2d6-a8c03be21bf1",
			},
			expectErr: false,
			expectURL: "http://example.com/api/games/3fa1b1b7-cca9-435e-b2d6-a8c03be21bf1/mages",
		},
		{
			name:   "Build GET URL with query params",
			method: http.MethodGet,
			path:   "/games/:game_id/mages",
			params: map[string]string{
				"game_id": "3fa1b1b7-cca9-435e-b2d6-a8c03be21bf1",
				"qp1":     "1",
			},
			expectErr: false,
			expectURL: "http://example.com/api/games/3fa1b1b7-cca9-435e-b2d6-a8c03be21bf1/mages?qp1=1",
		},
		{
			name:   "Build POST URL with additional ID",
			method: http.MethodPost,
			path:   "/games/:game_id/mages",
			params: map[string]string{
				"id":      "52fdfc07-2182-454f-963f-5f0f9a621d72",
				"game_id": "3fa1b1b7-cca9-435e-b2d6-a8c03be21bf1",
			},
			expectErr: false,
			expectURL: "http://example.com/api/games/3fa1b1b7-cca9-435e-b2d6-a8c03be21bf1/mages/52fdfc07-2182-454f-963f-5f0f9a621d72",
		},
		{
			name:   "Build POST URL without additional ID",
			method: http.MethodPost,
			path:   "/games/:game_id/mages",
			params: map[string]string{
				"game_id": "3fa1b1b7-cca9-435e-b2d6-a8c03be21bf1",
			},
			expectErr: false,
			expectURL: "http://example.com/api/games/3fa1b1b7-cca9-435e-b2d6-a8c03be21bf1/mages",
		},
		{
			name:   "Build PUT URL with ID",
			method: http.MethodPut,
			path:   "/games/:game_id/mages",
			params: map[string]string{
				"id":      "52fdfc07-2182-454f-963f-5f0f9a621d72",
				"game_id": "3fa1b1b7-cca9-435e-b2d6-a8c03be21bf1",
			},
			expectErr: false,
			expectURL: "http://example.com/api/games/3fa1b1b7-cca9-435e-b2d6-a8c03be21bf1/mages/52fdfc07-2182-454f-963f-5f0f9a621d72",
		},
		{
			name:   "Build PUT URL without ID",
			method: http.MethodPut,
			path:   "/games/:game_id/mages",
			params: map[string]string{
				"game_id": "3fa1b1b7-cca9-435e-b2d6-a8c03be21bf1",
			},
			expectErr: true,
			expectURL: "http://example.com/api/games/3fa1b1b7-cca9-435e-b2d6-a8c03be21bf1/mages",
		},
		{
			name:   "Build DELETE URL with ID",
			method: http.MethodDelete,
			path:   "/games/:game_id/mages",
			params: map[string]string{
				"id":      "52fdfc07-2182-454f-963f-5f0f9a621d72",
				"game_id": "3fa1b1b7-cca9-435e-b2d6-a8c03be21bf1",
			},
			expectErr: false,
			expectURL: "http://example.com/api/games/3fa1b1b7-cca9-435e-b2d6-a8c03be21bf1/mages/52fdfc07-2182-454f-963f-5f0f9a621d72",
		},
		{
			name:   "Build DELETE URL without ID",
			method: http.MethodDelete,
			path:   "/games/:game_id/mages",
			params: map[string]string{
				"game_id": "3fa1b1b7-cca9-435e-b2d6-a8c03be21bf1",
			},
			expectErr: true,
			expectURL: "http://example.com/api/games/3fa1b1b7-cca9-435e-b2d6-a8c03be21bf1/mages",
		},
		{
			name:   "Build GET URL with empty params",
			method: http.MethodGet,
			path:   "/games/:game_id/mages",
			params: map[string]string{
				"id":      "",
				"game_id": "",
			},
			expectErr: true,
		},
		// prefixed with colon
		{
			name:   "Build GET URL with additional ID - prefixed with colon",
			method: http.MethodGet,
			path:   "/games/:game_id/mages",
			params: map[string]string{
				":id":      "52fdfc07-2182-454f-963f-5f0f9a621d72",
				":game_id": "3fa1b1b7-cca9-435e-b2d6-a8c03be21bf1",
			},
			expectErr: false,
			expectURL: "http://example.com/api/games/3fa1b1b7-cca9-435e-b2d6-a8c03be21bf1/mages/52fdfc07-2182-454f-963f-5f0f9a621d72",
		},
		{
			name:   "Build GET URL without additional ID - prefixed with colon",
			method: http.MethodGet,
			path:   "/games/:game_id/mages",
			params: map[string]string{
				":game_id": "3fa1b1b7-cca9-435e-b2d6-a8c03be21bf1",
			},
			expectErr: false,
			expectURL: "http://example.com/api/games/3fa1b1b7-cca9-435e-b2d6-a8c03be21bf1/mages",
		},
		{
			name:   "Build GET URL with query params - prefixed with colon",
			method: http.MethodGet,
			path:   "/games/:game_id/mages",
			params: map[string]string{
				":game_id": "3fa1b1b7-cca9-435e-b2d6-a8c03be21bf1",
				":qp1":     "1",
			},
			expectErr: false,
			expectURL: "http://example.com/api/games/3fa1b1b7-cca9-435e-b2d6-a8c03be21bf1/mages?qp1=1",
		},
		{
			name:   "Build POST URL with additional ID - prefixed with colon",
			method: http.MethodPost,
			path:   "/games/:game_id/mages",
			params: map[string]string{
				":id":      "52fdfc07-2182-454f-963f-5f0f9a621d72",
				":game_id": "3fa1b1b7-cca9-435e-b2d6-a8c03be21bf1",
			},
			expectErr: false,
			expectURL: "http://example.com/api/games/3fa1b1b7-cca9-435e-b2d6-a8c03be21bf1/mages/52fdfc07-2182-454f-963f-5f0f9a621d72",
		},
		{
			name:   "Build POST URL without additional ID - prefixed with colon",
			method: http.MethodPost,
			path:   "/games/:game_id/mages",
			params: map[string]string{
				":game_id": "3fa1b1b7-cca9-435e-b2d6-a8c03be21bf1",
			},
			expectErr: false,
			expectURL: "http://example.com/api/games/3fa1b1b7-cca9-435e-b2d6-a8c03be21bf1/mages",
		},
		{
			name:   "Build PUT URL with ID - prefixed with colon",
			method: http.MethodPut,
			path:   "/games/:game_id/mages",
			params: map[string]string{
				":id":      "52fdfc07-2182-454f-963f-5f0f9a621d72",
				":game_id": "3fa1b1b7-cca9-435e-b2d6-a8c03be21bf1",
			},
			expectErr: false,
			expectURL: "http://example.com/api/games/3fa1b1b7-cca9-435e-b2d6-a8c03be21bf1/mages/52fdfc07-2182-454f-963f-5f0f9a621d72",
		},
		{
			name:   "Build PUT URL without ID - prefixed with colon",
			method: http.MethodPut,
			path:   "/games/:game_id/mages",
			params: map[string]string{
				":game_id": "3fa1b1b7-cca9-435e-b2d6-a8c03be21bf1",
			},
			expectErr: true,
			expectURL: "http://example.com/api/games/3fa1b1b7-cca9-435e-b2d6-a8c03be21bf1/mages",
		},
		{
			name:   "Build DELETE URL with ID - prefixed with colon",
			method: http.MethodDelete,
			path:   "/games/:game_id/mages",
			params: map[string]string{
				":id":      "52fdfc07-2182-454f-963f-5f0f9a621d72",
				":game_id": "3fa1b1b7-cca9-435e-b2d6-a8c03be21bf1",
			},
			expectErr: false,
			expectURL: "http://example.com/api/games/3fa1b1b7-cca9-435e-b2d6-a8c03be21bf1/mages/52fdfc07-2182-454f-963f-5f0f9a621d72",
		},
		{
			name:   "Build DELETE URL without ID - prefixed with colon",
			method: http.MethodDelete,
			path:   "/games/:game_id/mages",
			params: map[string]string{
				":game_id": "3fa1b1b7-cca9-435e-b2d6-a8c03be21bf1",
			},
			expectErr: true,
			expectURL: "http://example.com/api/games/3fa1b1b7-cca9-435e-b2d6-a8c03be21bf1/mages",
		},
	}

	for _, tc := range tests {

		t.Logf("Running test >%s<", tc.name)

		t.Run(tc.name, func(t *testing.T) {

			url, err := cl.BuildURL(tc.method, tc.path, tc.params)
			if tc.expectErr == true {
				require.Error(t, err, "BuildURL returns with error")
				return
			}
			t.Logf("Resulting URL >%s<", url)
			require.NoError(t, err, "BuildURL returns without error")
			require.Equal(t, tc.expectURL, url, "BuildURL returns expected URL")
		})
	}
}
