package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/go-boilerplate/server/core/auth"
	"gitlab.com/alienspaces/go-boilerplate/server/core/type/logger"
	"gitlab.com/alienspaces/go-boilerplate/server/core/type/modeller"
)

// Authen middleware tests
func TestAuthenMiddlewareConfig(t *testing.T) {

	// Test configuration
	type TestCase struct {
		name               string
		runner             func() TestRunner
		expectCachedConfig bool
	}

	tests := []TestCase{
		{
			name: "Without authentication configured",
			runner: func() TestRunner {
				r := TestRunner{}
				r.HandlerConfig = []HandlerConfig{
					{
						Method:           http.MethodPost,
						Path:             "/test",
						MiddlewareConfig: MiddlewareConfig{},
					},
				}

				return r
			},
			expectCachedConfig: false,
		},
		{
			name: "With authentication configured",
			runner: func() TestRunner {
				r := TestRunner{}
				r.HandlerConfig = []HandlerConfig{
					{
						Method: http.MethodPost,
						Path:   "/test",
						MiddlewareConfig: MiddlewareConfig{
							AuthTypes: []string{
								AuthTypeJWT,
							},
						},
					},
					{
						Method: http.MethodGet,
						Path:   "/test/:id",
						MiddlewareConfig: MiddlewareConfig{
							AuthTypes: []string{
								AuthTypeJWT,
							},
						},
					},
					{
						Method: http.MethodPut,
						Path:   "/test/:id",
						MiddlewareConfig: MiddlewareConfig{
							AuthTypes: []string{
								AuthTypeJWT,
							},
						},
					},
				}

				return r
			},
			expectCachedConfig: true,
		},
	}

	for _, tc := range tests {

		t.Logf("Running test >%s<", tc.name)

		c, l, s, err := NewDefaultDependencies()
		require.NoError(t, err, "NewDefaultDependencies returns without error")

		tr := tc.runner()

		err = tr.Init(c, l, s)
		require.NoError(t, err, "Runner Init returns without error")

		// Clear authentication cache
		authenCache = nil

		// Register authentication
		for _, hc := range tr.HandlerConfig {
			nextHandlerFunc, err := tr.Authen(hc, hc.HandlerFunc)
			require.NoError(t, err, "Authen returns without error")
			require.NotNil(t, nextHandlerFunc, "Authen return the next handler function")
		}

		for _, hc := range tr.HandlerConfig {

			if tc.expectCachedConfig == true {
				// Path cached
				cachedPath, ok := authenCache[hc.Path]
				require.True(t, ok, "Request path found in authen cache")
				require.NotEmpty(t, cachedPath, "Cached path data is not empty")
				// Method cached
				cachedMethod, ok := cachedPath[hc.Method]
				require.True(t, ok, "Request method found in authen cache")
				require.NotEmpty(t, cachedMethod, "Cached method data is not empty")
			} else {
				// Path not cached
				cachedPath, ok := authenCache[hc.Path]
				require.False(t, ok, "Request path found in authen cache")
				require.Empty(t, cachedPath, "Cached path data is not empty")
			}
		}
	}
}

func TestAuthenHandler(t *testing.T) {

	c, l, s, err := NewDefaultDependencies()
	require.NoError(t, err, "NewDefaultDependencies returns without error")

	const (
		testRoleName    string = "administrator"
		testIdentityKey string = "account"
	)

	// Test configuration
	type TestCase struct {
		name           string
		runner         func() TestRunner
		requestHeaders func() map[string]string
		responseCode   int
	}

	// validAuthToken - Generate a valid authentication token for this handler
	validAuthToken := func(roles []string, identity map[string]interface{}) string {
		authen, _ := auth.NewAuth(c, l)
		token, _ := authen.EncodeJWT(&auth.Claims{
			Roles:    roles,
			Identity: identity,
		})
		return token
	}

	tests := []TestCase{
		{
			name: "Without authentication configured and no token",
			runner: func() TestRunner {
				r := TestRunner{}
				r.HandlerConfig = []HandlerConfig{
					{
						Method:           http.MethodPost,
						Path:             "/test",
						MiddlewareConfig: MiddlewareConfig{},
						HandlerFunc: func(w http.ResponseWriter, r *http.Request, pathParams httprouter.Params, queryParams map[string]interface{}, l logger.Logger, m modeller.Modeller) {
							return
						},
					},
				}

				return r
			},
			responseCode: http.StatusOK,
		},
		{
			name: "With authentication configured and valid token",
			runner: func() TestRunner {
				rnr := TestRunner{}
				rnr.Init(c, l, s)
				rnr.HandlerConfig = []HandlerConfig{
					{
						Method: http.MethodPost,
						Path:   "/test",
						MiddlewareConfig: MiddlewareConfig{
							AuthTypes: []string{
								AuthTypeJWT,
							},
						},
						HandlerFunc: func(w http.ResponseWriter, r *http.Request, pathParams httprouter.Params, queryParams map[string]interface{}, l logger.Logger, m modeller.Modeller) {
							ctx := r.Context()
							hasRole, err := rnr.hasContextRole(ctx, testRoleName)
							if err != nil {
								w.WriteHeader(http.StatusUnauthorized)
								return
							}
							identityValue, err := rnr.getContextIdentityValue(ctx, testIdentityKey)
							if err != nil {
								w.WriteHeader(http.StatusUnauthorized)
								return
							}
							if hasRole && identityValue != nil {
								w.WriteHeader(http.StatusOK)
								return
							}
							w.WriteHeader(http.StatusUnauthorized)
							return
						},
					},
				}
				return rnr
			},
			requestHeaders: func() map[string]string {
				roles := []string{
					testRoleName,
				}
				identity := map[string]interface{}{
					testIdentityKey: gofakeit.UUID(),
				}
				headers := map[string]string{
					"Authorization": "Bearer " + validAuthToken(roles, identity),
				}
				return headers
			},
			responseCode: http.StatusOK,
		},
		{
			name: "With authentication configured and no token",
			runner: func() TestRunner {
				rnr := TestRunner{}
				rnr.Init(c, l, s)
				rnr.HandlerConfig = []HandlerConfig{
					{
						Method: http.MethodPost,
						Path:   "/test",
						MiddlewareConfig: MiddlewareConfig{
							AuthTypes: []string{
								AuthTypeJWT,
							},
						},
						HandlerFunc: func(w http.ResponseWriter, r *http.Request, pathParams httprouter.Params, queryParams map[string]interface{}, l logger.Logger, m modeller.Modeller) {
							return
						},
					},
				}
				return rnr
			},
			responseCode: http.StatusUnauthorized,
		},
	}

	for _, tc := range tests {

		t.Logf("Running test >%s<", tc.name)

		tr := tc.runner()

		// Clear authentication cache
		authenCache = nil

		// Register authentication
		hc := tr.HandlerConfig[0]
		hf, err := tr.Authen(hc, hc.HandlerFunc)
		require.NoError(t, err, "Authen returns without error")
		require.NotNil(t, hf, "Authen return the next handler function")

		// wrap everything in a httprouter Handler
		h := func(w http.ResponseWriter, r *http.Request, pp httprouter.Params) {
			// delegate
			hf(w, r, pp, nil, l, nil)
		}

		// router
		rtr := httprouter.New()

		switch hc.Method {
		case http.MethodGet:
			rtr.GET(hc.Path, h)
		case http.MethodPost:
			rtr.POST(hc.Path, h)
		case http.MethodPut:
			rtr.PUT(hc.Path, h)
		case http.MethodDelete:
			rtr.DELETE(hc.Path, h)
		default:
			//
		}

		req, err := http.NewRequest(hc.Method, hc.Path, nil)
		require.NoError(t, err, "NewRequest returns without error")

		// request headers
		requestHeaders := map[string]string{}
		if tc.requestHeaders != nil {
			requestHeaders = tc.requestHeaders()
		}

		for headerKey, headerVal := range requestHeaders {
			req.Header.Add(headerKey, headerVal)
		}

		// recorder
		rec := httptest.NewRecorder()

		// serve
		rtr.ServeHTTP(rec, req)

		// test status
		require.Equalf(t, tc.responseCode, rec.Code, "%s - Response code equals expected", tc.name)
	}
}
