package server

import (
	"context"
	"net/http"

	"github.com/julienschmidt/httprouter"

	coreerror "gitlab.com/alienspaces/go-mud/backend/core/error"
	"gitlab.com/alienspaces/go-mud/backend/core/queryparam"
	"gitlab.com/alienspaces/go-mud/backend/core/type/logger"
	"gitlab.com/alienspaces/go-mud/backend/core/type/modeller"
)

var unauthErr = coreerror.NewUnauthenticatedError("Authorization header value is missing or invalid.")

// AuthenMiddleware -
func (rnr *Runner) AuthenMiddleware(hc HandlerConfig, h Handle) (Handle, error) {
	handlerAuthenTypes := ToAuthenticationSet(hc.MiddlewareConfig.AuthenTypes...)

	handle := func(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m modeller.Modeller) error {
		l = Logger(l, "AuthenMiddleware")

		if _, ok := handlerAuthenTypes[AuthenticationTypePublic]; ok {
			l.Debug("Handler name >%s< is public, not authenticating", hc.Name)
			return h(w, r, pp, qp, l, m)
		}

		auth := r.Header.Get("Authorization")
		if auth == "" {
			l.Warn(unauthErr.Error())
			WriteError(l, w, unauthErr)
			return unauthErr
		}

		authenticatedRequest, err := rnr.AuthenticateRequestFunc(l, m, auth)
		if err != nil {
			l.Warn(unauthErr.Error())
			WriteError(l, w, unauthErr)
			return err
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, ctxKeyAuth, authenticatedRequest)
		r = r.WithContext(ctx)

		return h(w, r, pp, qp, l, m)
	}

	return handle, nil
}

func ToAuthenticationSet(authen ...AuthenticationType) map[AuthenticationType]struct{} {
	set := map[AuthenticationType]struct{}{}
	for _, p := range authen {
		set[p] = struct{}{}
	}
	return set
}
