package server

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"

	coreerror "gitlab.com/alienspaces/go-mud/backend/core/error"
	"gitlab.com/alienspaces/go-mud/backend/core/type/logger"
	"gitlab.com/alienspaces/go-mud/backend/core/type/modeller"
)

// Authz -
func (rnr *Runner) Authz(hc HandlerConfig, h Handle) (Handle, error) {
	authenTypes := ToAuthenticationSet(hc.MiddlewareConfig.AuthenTypes...)
	authzPermissions := ToAuthorizationPermissionsSet(hc.MiddlewareConfig.AuthzPermissions...)

	handle := func(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp map[string]interface{}, l logger.Logger, m modeller.Modeller) error {

		if _, ok := authenTypes[AuthenTypePublic]; ok {
			return h(w, r, pp, qp, l, m)
		}

		ctx := r.Context()
		auth, ok := ctx.Value(ctxKeyAuth).(Authentication)
		if !ok {
			err := fmt.Errorf("failed to convert ctxKeyAuth value >%#v< type to coreserver.Authentication", auth)
			l.Error(err.Error())
			WriteSystemError(l, w, err)
			return err
		}

		auth, err := rnr.GetAuthorizationsByHashedAPIKeyFunc(m, l, auth)
		if err != nil {
			l.Error("failed to validate hashed API key >%v<", err)
			WriteError(l, w, err)
			return err
		}

		for p := range authzPermissions {
			if _, ok := auth.Permissions[string(p)]; !ok {
				err := coreerror.NewUnauthorizedError()
				l.Error(err.Error())
				WriteError(l, w, err)
				return err
			}
		}

		return h(w, r, pp, qp, l, m)
	}

	return handle, nil
}
