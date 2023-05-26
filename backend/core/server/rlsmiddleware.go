package server

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"gitlab.com/alienspaces/go-mud/backend/core/queryparam"
	"gitlab.com/alienspaces/go-mud/backend/core/type/logger"
	"gitlab.com/alienspaces/go-mud/backend/core/type/modeller"
)

type RLS struct {
	Identifiers map[string][]string
}

func (rnr *Runner) RLSMiddleware(hc HandlerConfig, h Handle) (Handle, error) {
	handlerAuthenTypes := ToAuthenticationSet(hc.MiddlewareConfig.AuthenTypes...)

	handle := func(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m modeller.Modeller) error {
		l = Logger(l, "RLSMiddleware")

		if _, ok := handlerAuthenTypes[AuthenticationTypePublic]; ok {
			l.Debug("Handler name >%s< is public, not apply RLS", hc.Name)
			return h(w, r, pp, qp, l, m)
		}

		authenticatedRequest := AuthData(l, r)
		if authenticatedRequest == nil {
			err := fmt.Errorf("failed to read auth data")
			WriteSystemError(l, w, err)
			return err
		}

		if authenticatedRequest.Type == AuthenticatedTypeAPIKey && authenticatedRequest.RLSType == RLSTypeRestricted {
			if _, err := rnr.SetRLSFunc(l, m, *authenticatedRequest); err != nil {
				l.Warn("failed to set RLS >%v<", err)
				WriteSystemError(l, w, err)
				return err
			}
		}

		return h(w, r, pp, qp, l, m)
	}
	return handle, nil
}
