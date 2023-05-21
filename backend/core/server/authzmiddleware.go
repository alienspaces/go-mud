package server

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"

	coreerror "gitlab.com/alienspaces/go-mud/backend/core/error"
	"gitlab.com/alienspaces/go-mud/backend/core/queryparam"
	"gitlab.com/alienspaces/go-mud/backend/core/type/logger"
	"gitlab.com/alienspaces/go-mud/backend/core/type/modeller"
)

// AuthzMiddleware -
func (rnr *Runner) AuthzMiddleware(hc HandlerConfig, h Handle) (Handle, error) {
	authenTypes := ToAuthenticationSet(hc.MiddlewareConfig.AuthenTypes...)
	authzPermissions := ToAuthorizationPermissionsSet(hc.MiddlewareConfig.AuthzPermissions...)

	handle := func(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m modeller.Modeller) error {
		l = Logger(l, "AuthzMiddleware")

		if _, ok := authenTypes[AuthenticationTypePublic]; ok {
			l.Debug("handler name >%s< is public, not checking permissions", hc.Name)
			return h(w, r, pp, qp, l, m)
		}

		authenticatedRequest := AuthData(l, r)
		if authenticatedRequest == nil {
			err := fmt.Errorf("failed to read auth data")
			WriteSystemError(l, w, err)
			return err
		}

		for _, permission := range authenticatedRequest.Permissions {
			if _, ok := authzPermissions[permission]; ok {
				return h(w, r, pp, qp, l, m)
			}
		}

		err := fmt.Errorf("authenticated request >%v< does not contain any required permissions >%v<", authenticatedRequest, hc.MiddlewareConfig.AuthzPermissions)
		l.Warn(err.Error())
		WriteError(l, w, coreerror.NewUnauthorizedError())
		return err
	}
	return handle, nil
}

func ToAuthorizationPermissionsSet(permissions ...AuthorizedPermission) map[AuthorizedPermission]struct{} {
	set := map[AuthorizedPermission]struct{}{}
	for _, p := range permissions {
		set[p] = struct{}{}
	}
	return set
}
