package server

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"gitlab.com/alienspaces/go-boilerplate/server/core/type/logger"
	"gitlab.com/alienspaces/go-boilerplate/server/core/type/modeller"
)

// authzAllIdentitiesCache - path, method
var authzAllIdentitiesCache map[string]map[string][]string

// authzAnyIdentityCache - path, method
var authzAnyIdentityCache map[string]map[string][]string

// authzAllRolesCache - path, method
var authzAllRolesCache map[string]map[string][]string

// authzAnyRoleCache - path, method
var authzAnyRoleCache map[string]map[string][]string

// Authz -
func (rnr *Runner) Authz(hc HandlerConfig, h HandlerFunc) (HandlerFunc, error) {

	// Cache authz configuration
	err := rnr.authzCacheConfig(hc)
	if err != nil {
		rnr.Log.Warn("Failed caching authz config >%v<", err)
		return nil, err
	}

	handle := func(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp map[string]interface{}, l logger.Logger, m modeller.Modeller) {

		l.Debug("** Authz **")

		err := rnr.handleAuthz(r, l, m, hc)
		if err != nil {
			l.Warn("Failed authz >%v<", err)
			rnr.WriteUnauthorizedError(l, w, err)
			return
		}

		h(w, r, pp, qp, l, m)
	}

	return handle, nil
}

// handleAuthz -
func (rnr *Runner) handleAuthz(r *http.Request, l logger.Logger, m modeller.Modeller, hc HandlerConfig) error {

	// Authorization may need to check identities and roles from context
	ctx := r.Context()

	if authzAllIdentitiesCache == nil && authzAnyIdentityCache == nil && authzAllRolesCache == nil && authzAnyRoleCache == nil {
		l.Info("Authz not configured")
		return nil
	}

	// Identities
	var authzIdentityMethods map[string][]string
	var authzIdentities []string

	// Check all identities
	authzIdentityMethods = authzAllIdentitiesCache[hc.Path]
	if authzIdentityMethods != nil {
		authzIdentities = authzIdentityMethods[hc.Method]
		if authzIdentities != nil {
			for _, identityKey := range authzIdentities {
				identityValue, err := rnr.getContextIdentityValue(ctx, identityKey)
				if err != nil {
					l.Warn("Failed checking context identity key >%s< >%v<", identityKey, err)
					return err
				}
				if identityValue == nil {
					msg := fmt.Sprintf("Context missing identity key >%s< value", identityKey)
					l.Warn(msg)
					return fmt.Errorf(msg)
				}
				l.Info("Have identity key >%s< value >%s<", identityKey, identityValue)
			}
		}
	}

	// Check any identities
	authzIdentityMethods = authzAnyIdentityCache[hc.Path]
	if authzIdentityMethods != nil {
		authzIdentities = authzIdentityMethods[hc.Method]
		if authzIdentities != nil {
			found := false
			for _, identityKey := range authzIdentities {
				identityValue, err := rnr.getContextIdentityValue(ctx, identityKey)
				if err != nil {
					l.Warn("Failed checking context identity key >%s< >%v<", identityKey, err)
					return err
				}
				if identityValue != nil {
					l.Info("Have identity key >%s< value >%s<", identityKey, identityValue)
					found = true
					break
				}
				l.Info("Context missing identity key >%s< value", identityKey)
			}
			if !found {
				msg := fmt.Sprintf("Missing any identities")
				l.Warn(msg)
				return fmt.Errorf(msg)
			}
		}
	}

	// Roles
	var authzRoleMethods map[string][]string
	var authzRoles []string

	// Check all roles
	authzRoleMethods = authzAllRolesCache[hc.Path]
	if authzRoleMethods != nil {
		l.Info("Have roles >%#v<", authzRoleMethods[hc.Method])
		authzRoles = authzRoleMethods[hc.Method]
		if authzRoles != nil {
			for _, roleName := range authzRoles {
				hasRole, err := rnr.hasContextRole(ctx, roleName)
				if err != nil {
					l.Warn("Failed checking context role >%s< >%v<", roleName, err)
					return err
				}
				if hasRole != true {
					msg := fmt.Sprintf("Context missing role >%s< value", roleName)
					l.Warn(msg)
					return fmt.Errorf(msg)
				}
				l.Info("Have role name >%s<", roleName)
			}
		}
	}

	// Check any roles
	authzRoleMethods = authzAnyRoleCache[hc.Path]
	if authzRoleMethods != nil {
		authzRoles = authzRoleMethods[hc.Method]
		if authzRoles != nil {
			found := false
			for _, roleName := range authzRoles {
				hasRole, err := rnr.hasContextRole(ctx, roleName)
				if err != nil {
					l.Warn("Failed checking context role >%s< >%v<", roleName, err)
					return err
				}
				if hasRole == true {
					l.Info("Have role")
					found = true
					break
				}
				l.Info("Missing role >%s<", roleName)
			}
			if !found {
				msg := fmt.Sprintf("Missing any role")
				l.Warn(msg)
				return fmt.Errorf(msg)
			}
		}
	}

	return nil
}

// authzCacheConfig - cache authz configuration
func (rnr *Runner) authzCacheConfig(hc HandlerConfig) error {

	// Cache required identities
	if hc.MiddlewareConfig.AuthRequireAllIdentities != nil {
		if authzAllIdentitiesCache == nil {
			authzAllIdentitiesCache = make(map[string]map[string][]string)
		}
		if authzAllIdentitiesCache[hc.Path] == nil {
			authzAllIdentitiesCache[hc.Path] = make(map[string][]string)
		}
		for _, authRequiredIdentity := range hc.MiddlewareConfig.AuthRequireAllIdentities {
			rnr.Log.Info("Adding required all path >%s< method >%s< identity >%s<", hc.Path, hc.Method, authRequiredIdentity)
			authzAllIdentitiesCache[hc.Path][hc.Method] = append(authzAllIdentitiesCache[hc.Path][hc.Method], authRequiredIdentity)
		}
	}

	// Cache any identities
	if hc.MiddlewareConfig.AuthRequireAnyIdentity != nil {
		if authzAnyIdentityCache == nil {
			authzAnyIdentityCache = make(map[string]map[string][]string)
		}
		if authzAnyIdentityCache[hc.Path] == nil {
			authzAnyIdentityCache[hc.Path] = make(map[string][]string)
		}
		for _, authRequiredIdentity := range hc.MiddlewareConfig.AuthRequireAnyIdentity {
			rnr.Log.Info("Adding required any path >%s< method >%s< identity >%s<", hc.Path, hc.Method, authRequiredIdentity)
			authzAnyIdentityCache[hc.Path][hc.Method] = append(authzAnyIdentityCache[hc.Path][hc.Method], authRequiredIdentity)
		}
	}

	// Cache required roles
	if hc.MiddlewareConfig.AuthRequireAllRoles != nil {
		if authzAllRolesCache == nil {
			authzAllRolesCache = make(map[string]map[string][]string)
		}
		if authzAllRolesCache[hc.Path] == nil {
			authzAllRolesCache[hc.Path] = make(map[string][]string)
		}
		for _, authRequiredRole := range hc.MiddlewareConfig.AuthRequireAllRoles {
			rnr.Log.Info("Adding required all path >%s< method >%s< role >%s<", hc.Path, hc.Method, authRequiredRole)
			authzAllRolesCache[hc.Path][hc.Method] = append(authzAllRolesCache[hc.Path][hc.Method], authRequiredRole)
		}
	}

	// Cache any roles
	if hc.MiddlewareConfig.AuthRequireAnyRole != nil {
		if authzAnyRoleCache == nil {
			authzAnyRoleCache = make(map[string]map[string][]string)
		}
		if authzAnyRoleCache[hc.Path] == nil {
			authzAnyRoleCache[hc.Path] = make(map[string][]string)
		}
		for _, authRequiredRole := range hc.MiddlewareConfig.AuthRequireAnyRole {
			rnr.Log.Info("Adding required any path >%s< method >%s< role >%s<", hc.Path, hc.Method, authRequiredRole)
			authzAnyRoleCache[hc.Path][hc.Method] = append(authzAnyRoleCache[hc.Path][hc.Method], authRequiredRole)
		}
	}

	return nil
}
