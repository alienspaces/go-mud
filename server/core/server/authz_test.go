package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/go-boilerplate/server/core/type/logger"
	"gitlab.com/alienspaces/go-boilerplate/server/core/type/modeller"
)

const (
	testRoleDefault       string = "default"
	testRoleAdministrator string = "administrator"
	testIdentityPlayer   string = "account_id"
	testIdentityCompany   string = "company_id"
)

// Authz middleware tests
func TestAuthzMiddlewareConfig(t *testing.T) {

	// Test configuration
	type TestCase struct {
		name                            string
		runner                          func() TestRunner
		expectAllIdentitiesCachedConfig bool
		expectAnyIdentityCachedConfig   bool
		expectAllRolesCachedConfig      bool
		expectAnyRoleCachedConfig       bool
	}

	tests := []TestCase{
		{
			name: "Without authorization configured",
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
			expectAllIdentitiesCachedConfig: false,
			expectAnyIdentityCachedConfig:   false,
			expectAllRolesCachedConfig:      false,
			expectAnyRoleCachedConfig:       false,
		},
		{
			name: "With authorization configured",
			runner: func() TestRunner {
				r := TestRunner{}
				r.HandlerConfig = []HandlerConfig{
					{
						Method: http.MethodPost,
						Path:   "/test",
						MiddlewareConfig: MiddlewareConfig{
							AuthRequireAllIdentities: []string{
								testIdentityPlayer,
							},
							AuthRequireAnyIdentity: []string{
								testIdentityCompany,
							},
							AuthRequireAllRoles: []string{
								testRoleDefault,
							},
							AuthRequireAnyRole: []string{
								testRoleAdministrator,
							},
						},
					},
					{
						Method: http.MethodGet,
						Path:   "/test/:id",
						MiddlewareConfig: MiddlewareConfig{
							AuthRequireAllIdentities: []string{
								testIdentityPlayer,
							},
							AuthRequireAnyIdentity: []string{
								testIdentityCompany,
							},
							AuthRequireAllRoles: []string{
								testRoleDefault,
							},
							AuthRequireAnyRole: []string{
								testRoleAdministrator,
							},
						},
					},
					{
						Method: http.MethodPut,
						Path:   "/test/:id",
						MiddlewareConfig: MiddlewareConfig{
							AuthRequireAllIdentities: []string{
								testIdentityPlayer,
							},
							AuthRequireAnyIdentity: []string{
								testIdentityCompany,
							},
							AuthRequireAllRoles: []string{
								testRoleDefault,
							},
							AuthRequireAnyRole: []string{
								testRoleAdministrator,
							},
						},
					},
				}

				return r
			},
			expectAllIdentitiesCachedConfig: true,
			expectAnyIdentityCachedConfig:   true,
			expectAllRolesCachedConfig:      true,
			expectAnyRoleCachedConfig:       true,
		},
	}

	for _, tc := range tests {

		t.Logf("Running test >%s<", tc.name)

		c, l, s, err := NewDefaultDependencies()
		require.NoError(t, err, "NewDefaultDependencies returns without error")

		tr := tc.runner()

		err = tr.Init(c, l, s)
		require.NoError(t, err, "Runner Init returns without error")

		// Clear auths caches
		authzAllIdentitiesCache = nil
		authzAnyIdentityCache = nil
		authzAllRolesCache = nil
		authzAnyRoleCache = nil

		// Register authorization
		for _, hc := range tr.HandlerConfig {
			nextHandlerFunc, err := tr.Authz(hc, hc.HandlerFunc)
			require.NoError(t, err, "Authz returns without error")
			require.NotNil(t, nextHandlerFunc, "Authz return the next handler function")
		}

		for _, hc := range tr.HandlerConfig {

			if tc.expectAllIdentitiesCachedConfig == true {
				// Path cached
				cachedPath, ok := authzAllIdentitiesCache[hc.Path]
				require.True(t, ok, "Request path found in authz all identities cache")
				require.NotEmpty(t, cachedPath, "Cached path data is not empty")
				// Method cached
				cachedMethod, ok := cachedPath[hc.Method]
				require.True(t, ok, "Request method found in authz all identities cache")
				require.NotEmpty(t, cachedMethod, "Cached method data is not empty")
			} else {
				// Path not cached
				cachedPath, ok := authzAllIdentitiesCache[hc.Path]
				require.False(t, ok, "Request path found in authz all identities cache")
				require.Empty(t, cachedPath, "Cached path data is not empty")
			}

			if tc.expectAnyIdentityCachedConfig == true {
				// Path cached
				cachedPath, ok := authzAnyIdentityCache[hc.Path]
				require.True(t, ok, "Request path found in authz any identities cache")
				require.NotEmpty(t, cachedPath, "Cached path data is not empty")
				// Method cached
				cachedMethod, ok := cachedPath[hc.Method]
				require.True(t, ok, "Request method found in authz any identities cache")
				require.NotEmpty(t, cachedMethod, "Cached method data is not empty")
			} else {
				// Path not cached
				cachedPath, ok := authzAnyIdentityCache[hc.Path]
				require.False(t, ok, "Request path found in authz any identities cache")
				require.Empty(t, cachedPath, "Cached path data is not empty")
			}

			if tc.expectAllRolesCachedConfig == true {
				// Path cached
				cachedPath, ok := authzAllRolesCache[hc.Path]
				require.True(t, ok, "Request path found in authz all roles cache")
				require.NotEmpty(t, cachedPath, "Cached path data is not empty")
				// Method cached
				cachedMethod, ok := cachedPath[hc.Method]
				require.True(t, ok, "Request method found in authz all roles cache")
				require.NotEmpty(t, cachedMethod, "Cached method data is not empty")
			} else {
				// Path not cached
				cachedPath, ok := authzAllRolesCache[hc.Path]
				require.False(t, ok, "Request path found in authz all roles cache")
				require.Empty(t, cachedPath, "Cached path data is not empty")
			}

			if tc.expectAnyRoleCachedConfig == true {
				// Path cached
				cachedPath, ok := authzAnyRoleCache[hc.Path]
				require.True(t, ok, "Request path found in authz any role cache")
				require.NotEmpty(t, cachedPath, "Cached path data is not empty")
				// Method cached
				cachedMethod, ok := cachedPath[hc.Method]
				require.True(t, ok, "Request method found in authz any role cache")
				require.NotEmpty(t, cachedMethod, "Cached method data is not empty")
			} else {
				// Path not cached
				cachedPath, ok := authzAnyRoleCache[hc.Path]
				require.False(t, ok, "Request path found in authz any role cache")
				require.Empty(t, cachedPath, "Cached path data is not empty")
			}
		}
	}
}

func TestAuthzHandler(t *testing.T) {

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
		requestContext func(tr TestRunner, ctx context.Context) context.Context
		responseCode   int
	}

	// handlerFunc - Default handler function
	handlerFunc := func(w http.ResponseWriter, r *http.Request, pathParams httprouter.Params, queryParams map[string]interface{}, l logger.Logger, m modeller.Modeller) {
		return
	}

	tests := []TestCase{
		{
			name: "Without authorization configured and no context",
			runner: func() TestRunner {
				rnr := TestRunner{}
				rnr.Init(c, l, s)
				rnr.HandlerConfig = []HandlerConfig{
					{
						Method:           http.MethodPost,
						Path:             "/test",
						MiddlewareConfig: MiddlewareConfig{},
						HandlerFunc:      handlerFunc,
					},
				}
				return rnr
			},
			responseCode: http.StatusOK,
		},
		{
			name: "With authorization any identity configured and valid context",
			runner: func() TestRunner {
				rnr := TestRunner{}
				rnr.Init(c, l, s)
				rnr.HandlerConfig = []HandlerConfig{
					{
						Method: http.MethodPost,
						Path:   "/test",
						MiddlewareConfig: MiddlewareConfig{
							AuthRequireAnyIdentity: []string{
								testIdentityPlayer,
							},
						},
						HandlerFunc: handlerFunc,
					},
				}
				return rnr
			},
			requestContext: func(tr TestRunner, ctx context.Context) context.Context {
				ctx, _ = tr.addContextIdentity(ctx, map[string]interface{}{
					testIdentityPlayer: gofakeit.UUID(),
				})
				return ctx
			},
			responseCode: http.StatusOK,
		},
		{
			name: "With authorization all identities configured and valid context",
			runner: func() TestRunner {
				rnr := TestRunner{}
				rnr.Init(c, l, s)
				rnr.HandlerConfig = []HandlerConfig{
					{
						Method: http.MethodPost,
						Path:   "/test",
						MiddlewareConfig: MiddlewareConfig{
							AuthRequireAllIdentities: []string{
								testIdentityPlayer,
								testIdentityCompany,
							},
						},
						HandlerFunc: handlerFunc,
					},
				}
				return rnr
			},
			requestContext: func(tr TestRunner, ctx context.Context) context.Context {
				ctx, _ = tr.addContextIdentity(ctx, map[string]interface{}{
					testIdentityPlayer: gofakeit.UUID(),
					testIdentityCompany: gofakeit.UUID(),
				})
				return ctx
			},
			responseCode: http.StatusOK,
		},
		{
			name: "With authorization any identity configured and invalid context",
			runner: func() TestRunner {
				rnr := TestRunner{}
				rnr.Init(c, l, s)
				rnr.HandlerConfig = []HandlerConfig{
					{
						Method: http.MethodPost,
						Path:   "/test",
						MiddlewareConfig: MiddlewareConfig{
							AuthRequireAnyIdentity: []string{
								testIdentityPlayer,
							},
						},
						HandlerFunc: handlerFunc,
					},
				}
				return rnr
			},
			requestContext: func(tr TestRunner, ctx context.Context) context.Context {
				return ctx
			},
			responseCode: http.StatusUnauthorized,
		},
		{
			name: "With authorization all identities configured and invalid context",
			runner: func() TestRunner {
				rnr := TestRunner{}
				rnr.Init(c, l, s)
				rnr.HandlerConfig = []HandlerConfig{
					{
						Method: http.MethodPost,
						Path:   "/test",
						MiddlewareConfig: MiddlewareConfig{
							AuthRequireAllIdentities: []string{
								testIdentityPlayer,
								testIdentityCompany,
							},
						},
						HandlerFunc: handlerFunc,
					},
				}
				return rnr
			},
			requestContext: func(tr TestRunner, ctx context.Context) context.Context {
				ctx, _ = tr.addContextIdentity(ctx, map[string]interface{}{
					testIdentityPlayer: gofakeit.UUID(),
				})
				return ctx
			},
			responseCode: http.StatusUnauthorized,
		},
		{
			name: "With authorization any role configured and valid context",
			runner: func() TestRunner {
				rnr := TestRunner{}
				rnr.Init(c, l, s)
				rnr.HandlerConfig = []HandlerConfig{
					{
						Method: http.MethodPost,
						Path:   "/test",
						MiddlewareConfig: MiddlewareConfig{
							AuthRequireAnyRole: []string{
								testRoleDefault,
							},
						},
						HandlerFunc: handlerFunc,
					},
				}
				return rnr
			},
			requestContext: func(tr TestRunner, ctx context.Context) context.Context {
				ctx, _ = tr.addContextRoles(ctx, []string{
					testRoleDefault,
				})
				return ctx
			},
			responseCode: http.StatusOK,
		},
		{
			name: "With authorization all roles configured and valid context",
			runner: func() TestRunner {
				rnr := TestRunner{}
				rnr.Init(c, l, s)
				rnr.HandlerConfig = []HandlerConfig{
					{
						Method: http.MethodPost,
						Path:   "/test",
						MiddlewareConfig: MiddlewareConfig{
							AuthRequireAllRoles: []string{
								testRoleDefault,
								testRoleAdministrator,
							},
						},
						HandlerFunc: handlerFunc,
					},
				}
				return rnr
			},
			requestContext: func(tr TestRunner, ctx context.Context) context.Context {
				ctx, _ = tr.addContextRoles(ctx, []string{
					testRoleDefault,
					testRoleAdministrator,
				})
				return ctx
			},
			responseCode: http.StatusOK,
		},
		{
			name: "With authorization any role configured and invalid context",
			runner: func() TestRunner {
				rnr := TestRunner{}
				rnr.Init(c, l, s)
				rnr.HandlerConfig = []HandlerConfig{
					{
						Method: http.MethodPost,
						Path:   "/test",
						MiddlewareConfig: MiddlewareConfig{
							AuthRequireAnyRole: []string{
								testRoleDefault,
							},
						},
						HandlerFunc: handlerFunc,
					},
				}
				return rnr
			},
			requestContext: func(tr TestRunner, ctx context.Context) context.Context {
				return ctx
			},
			responseCode: http.StatusUnauthorized,
		},
		{
			name: "With authorization all roles configured and invalid context",
			runner: func() TestRunner {
				rnr := TestRunner{}
				rnr.Init(c, l, s)
				rnr.HandlerConfig = []HandlerConfig{
					{
						Method: http.MethodPost,
						Path:   "/test",
						MiddlewareConfig: MiddlewareConfig{
							AuthRequireAllRoles: []string{
								testRoleDefault,
								testRoleAdministrator,
							},
						},
						HandlerFunc: handlerFunc,
					},
				}
				return rnr
			},
			requestContext: func(tr TestRunner, ctx context.Context) context.Context {
				ctx, _ = tr.addContextRoles(ctx, []string{
					testRoleDefault,
				})
				return ctx
			},
			responseCode: http.StatusUnauthorized,
		},
	}

	for _, tc := range tests {

		t.Logf("Running test >%s<", tc.name)

		tr := tc.runner()

		// Clear auths caches
		authzAllIdentitiesCache = nil
		authzAnyIdentityCache = nil
		authzAllRolesCache = nil
		authzAnyRoleCache = nil

		// Register authentication
		hc := tr.HandlerConfig[0]
		hf, err := tr.Authz(hc, hc.HandlerFunc)
		require.NoError(t, err, "Authz returns without error")
		require.NotNil(t, hf, "Authz return the next handler function")

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

		// request context
		if tc.requestContext != nil {
			ctx := tc.requestContext(tr, req.Context())
			req = req.WithContext(ctx)
		}

		// recorder
		rec := httptest.NewRecorder()

		// serve
		rtr.ServeHTTP(rec, req)

		// test status
		require.Equalf(t, tc.responseCode, rec.Code, "%s - Response code equals expected", tc.name)
	}
}
