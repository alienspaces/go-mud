package server

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"

	coreerror "gitlab.com/alienspaces/go-mud/backend/core/error"
	"gitlab.com/alienspaces/go-mud/backend/core/type/logger"
	"gitlab.com/alienspaces/go-mud/backend/core/type/modeller"
)

// Authen -
func (rnr *Runner) Authen(hc HandlerConfig, h Handle) (Handle, error) {
	authenTypes := ToAuthenticationSet(hc.MiddlewareConfig.AuthenTypes...)

	handle := func(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp map[string]interface{}, l logger.Logger, m modeller.Modeller) error {

		if _, ok := authenTypes[AuthenTypePublic]; ok {
			return h(w, r, pp, qp, l, m)
		}

		xAuthorization := r.Header.Get("X-Authorization")

		apiKey, apiKeyErr := uuid.Parse(xAuthorization)
		if _, ok := authenTypes[AuthenTypeAPIKey]; ok && apiKeyErr == nil {
			return rnr.handleAuthByAPIKey(w, r, pp, qp, l, m, h, apiKey.String())
		}

		err := coreerror.NewClientUnauthenticatedError("X-Authorization header value is missing or invalid.")
		l.Warn(err.Error())
		WriteError(l, w, err)

		return err
	}

	return handle, nil
}

func (rnr *Runner) handleAuthByAPIKey(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp map[string]interface{}, l logger.Logger, m modeller.Modeller, h Handle, apiKey string) error {
	auth, err := rnr.AuthenticateByAPIKeyFunc(m, l, apiKey)
	if err != nil {
		l.Error("failed to validate hashed API key >%v<", err)
		WriteError(l, w, err)
		return err
	}

	if auth.IsValid {
		l.Context("API key", auth.HashedAPIKey)
		defer func() {
			l.Context("API key", "")
		}()

		ctx := r.Context()
		ctx = context.WithValue(ctx, ctxKeyAuth, auth)
		r = r.WithContext(ctx)
		return h(w, r, pp, qp, l, m)
	}

	m.Rollback()
	err = coreerror.NewClientUnauthenticatedError("X-Authorization header API key is missing or invalid.")
	WriteError(l, w, err)
	return err
}
