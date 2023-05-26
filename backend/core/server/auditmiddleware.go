package server

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"gitlab.com/alienspaces/go-mud/backend/core/queryparam"
	"gitlab.com/alienspaces/go-mud/backend/core/type/logger"
	"gitlab.com/alienspaces/go-mud/backend/core/type/modeller"
)

const (
	RequesterTypeAPIKey          = "api_key"
	RequesterTypeServiceCloudJWT = "service_cloud_jwt"
)

func (rnr *Runner) AuditMiddleware(hc HandlerConfig, h Handle) (Handle, error) {
	authenTypes := ToAuthenticationSet(hc.MiddlewareConfig.AuthenTypes...)

	handle := func(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m modeller.Modeller) error {
		l = Logger(l, "AuditMiddleware")

		isMutatingHTTPMethod := hc.Method == http.MethodPost ||
			hc.Method == http.MethodPut ||
			hc.Method == http.MethodPatch ||
			hc.Method == http.MethodDelete

		if _, ok := authenTypes[AuthenticationTypePublic]; !isMutatingHTTPMethod || ok {
			return h(w, r, pp, qp, l, m)
		}

		ctx := r.Context()
		correlationID, _ := ctx.Value(ctxKeyCorrelationID).(string)
		if correlationID == "" {
			err := fmt.Errorf("failed to retrieve request ID using ctxKeyCorrelationID")
			WriteSystemError(l, w, err)
			return err
		}

		auth := AuthData(l, r)
		if auth == nil {
			err := fmt.Errorf("failed to read auth data")
			WriteSystemError(l, w, err)
			return err
		}

		var requesterID string
		switch id := auth.User.ID.(type) {
		case string:
			requesterID = id
		case []byte:
			requesterID = base64.URLEncoding.EncodeToString(id)
		default:
			err := fmt.Errorf("unknown auth user ID type ")
			l.Warn(err.Error())
			return err
		}

		auditReq := AuditRequest{
			RequestID:      correlationID,
			RequesterID:    requesterID,
			RequesterName:  auth.User.Name,
			RequesterEmail: auth.User.Email,
		}

		switch auth.Type {
		case AuthenticatedTypeAPIKey:
			l.Debug("Setting API key audit data")

			auditReq.RequesterType = RequesterTypeAPIKey

			if err := rnr.SetAuditConfigFunc(l, m, auditReq); err != nil {
				err := fmt.Errorf("failed to set audit config for API key >%w<", err)
				WriteSystemError(l, w, err)
				return err
			}
		case AuthenticatedTypeUser:
			l.Debug("Setting JWT audit data")

			auditReq.RequesterType = RequesterTypeServiceCloudJWT

			if err := rnr.SetAuditConfigFunc(l, m, auditReq); err != nil {
				err := fmt.Errorf("failed to set audit config for API key >%w<", err)
				WriteSystemError(l, w, err)
				return err
			}
		default:
			err := fmt.Errorf("unknown authentication type, failed to set audit config >%s<", auth.Type)
			WriteSystemError(l, w, err)
			return err
		}

		return h(w, r, pp, qp, l, m)
	}

	return handle, nil
}
