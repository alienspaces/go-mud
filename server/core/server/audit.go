package server

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"gitlab.com/alienspaces/go-mud/server/core/type/logger"
	"gitlab.com/alienspaces/go-mud/server/core/type/modeller"
)

const RequesterTypeAPIKey = "api_key"

func (rnr *Runner) Audit(hc HandlerConfig, h Handle) (Handle, error) {
	authenTypes := ToAuthenticationSet(hc.MiddlewareConfig.AuthenTypes...)

	handle := func(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp map[string]interface{}, l logger.Logger, m modeller.Modeller) error {
		isMutatingHTTPMethod := hc.Method == http.MethodPost ||
			hc.Method == http.MethodPut ||
			hc.Method == http.MethodPatch ||
			hc.Method == http.MethodDelete

		if _, ok := authenTypes[AuthenTypePublic]; !isMutatingHTTPMethod || ok {
			return h(w, r, pp, qp, l, m)
		}

		ctx := r.Context()
		auth, ok := (ctx.Value(ctxKeyAuth)).(Authentication)
		if !ok {
			err := fmt.Errorf("failed to convert ctxKeyAuth value to Authentication")
			WriteSystemError(l, w, err)
			return err
		}

		correlationID, ok := ctx.Value(ctxKeyCorrelationID).(string)
		if correlationID == "" {
			err := fmt.Errorf("failed to convert ctxKeyCorrelationID value to string")
			WriteSystemError(l, w, err)
			return err
		} else if correlationID == "" {
			err := fmt.Errorf("failed to retrieve request ID using ctxKeyCorrelationID")
			WriteSystemError(l, w, err)
			return err
		}

		if err := rnr.SetAuditConfigFunc(m, l, RequesterTypeAPIKey, auth.HashedAPIKey, correlationID); err != nil {
			err := fmt.Errorf("failed to set audit config >%w<", err)
			WriteSystemError(l, w, err)
			return err
		}

		return h(w, r, pp, qp, l, m)
	}

	return handle, nil
}
