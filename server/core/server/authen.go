package server

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"

	"gitlab.com/alienspaces/go-boilerplate/server/core/auth"
	"gitlab.com/alienspaces/go-boilerplate/server/core/type/logger"
	"gitlab.com/alienspaces/go-boilerplate/server/core/type/modeller"
)

// auth
var authen *auth.Auth

// authenCache - path, method
var authenCache map[string]map[string][]string

// Authen -
func (rnr *Runner) Authen(hc HandlerConfig, h HandlerFunc) (HandlerFunc, error) {

	var err error
	if authen == nil {
		authen, err = auth.NewAuth(rnr.Config, rnr.Log)
		if err != nil {
			rnr.Log.Warn("Failed new auth >%v<", err)
			return nil, err
		}
	}

	// Cache authen configuration
	err = rnr.authenCacheConfig(hc)
	if err != nil {
		rnr.Log.Warn("Failed caching authen config >%v<", err)
		return nil, err
	}

	handle := func(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp map[string]interface{}, l logger.Logger, m modeller.Modeller) {

		l.Info("** Authen **")

		ctx, err := rnr.handleAuthen(r, l, m, hc)
		if err != nil {
			l.Warn("Failed authen >%v<", err)
			rnr.WriteUnauthorizedError(l, w, err)
			return
		}

		h(w, r.WithContext(ctx), pp, qp, l, m)
	}

	return handle, nil
}

// handleAuthen -
func (rnr *Runner) handleAuthen(r *http.Request, l logger.Logger, m modeller.Modeller, hc HandlerConfig) (context.Context, error) {

	// Authentication may add roles and identities to request context
	ctx := r.Context()

	if authenCache == nil {
		l.Info("Authen not configured")
		return ctx, nil
	}

	authenMethods := authenCache[hc.Path]
	if authenMethods == nil {
		l.Info("Authentication not configured for path >%s<", hc.Path)
		return ctx, nil
	}

	authenTypes := authenMethods[hc.Method]
	if authenTypes == nil {
		l.Info("Authentication not configured for path >%s< method >%s<", hc.Path, hc.Method)
		return ctx, nil
	}

	for _, authenType := range authenTypes {
		switch authenType {
		case auth.AuthTypeJWT:
			l.Info("** Authen ** JWT")
			// Get authentication token
			authString := r.Header.Get("Authorization")
			if authString == "" {
				msg := "Authorization header is empty"
				l.Warn(msg)
				return ctx, fmt.Errorf(msg)
			}
			if strings.Contains(authString, "Bearer ") {
				authString = strings.Split(authString, "Bearer ")[1]
			}

			// Decode authentication token
			claims, err := authen.DecodeJWT(authString)
			if err != nil {
				l.Warn("Failed authenticating token >%v<", err)
				return ctx, err
			}

			l.Info("Have claims >%#v<", claims)

			// Add roles to request context
			ctx, err = rnr.addContextRoles(ctx, claims.Roles)
			if err != nil {
				l.Warn("Failed adding roles context >%v<", err)
				return ctx, err
			}

			// Add identity to request context
			ctx, err = rnr.addContextIdentity(ctx, claims.Identity)
			if err != nil {
				return ctx, err
			}

		default:
			// Unsupported authentication configuration
			msg := "Unsupported authentication configuration"
			return ctx, fmt.Errorf(msg)
		}
	}

	return ctx, nil
}

// authenCacheConfig - cache authen configuration
func (rnr *Runner) authenCacheConfig(hc HandlerConfig) error {

	if hc.MiddlewareConfig.AuthTypes != nil {
		if authenCache == nil {
			authenCache = make(map[string]map[string][]string)
		}
		if authenCache[hc.Path] == nil {
			authenCache[hc.Path] = make(map[string][]string)
		}
		for _, authType := range hc.MiddlewareConfig.AuthTypes {
			authenCache[hc.Path][hc.Method] = append(authenCache[hc.Path][hc.Method], authType)
		}
	}

	return nil
}
