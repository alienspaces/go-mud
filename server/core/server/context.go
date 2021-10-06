package server

import (
	"context"
)

// All context keys are of type string
type contextKey string

// Context keys
const (
	contextRolesKey    contextKey = "authRoles"
	contextIdentityKey contextKey = "authIdentity"
)

// addRolesContext -
func (rnr *Runner) addContextRoles(ctx context.Context, roles []string) (context.Context, error) {

	ctx = context.WithValue(ctx, contextRolesKey, roles)

	return ctx, nil
}

// getRolesContext -
func (rnr *Runner) getContextRoles(ctx context.Context) ([]string, error) {

	roles := ctx.Value(contextRolesKey)
	if roles != nil {
		return roles.([]string), nil
	}

	rnr.Log.Info("Context roles is nil")

	return nil, nil
}

// hasRoleContext -
func (rnr *Runner) hasContextRole(ctx context.Context, roleName string) (bool, error) {

	roles, err := rnr.getContextRoles(ctx)
	if err != nil {
		return false, err
	}
	if roles != nil {
		for _, contextRoleName := range roles {
			if contextRoleName == roleName {
				rnr.Log.Info("Context role >%s< name exists", contextRoleName)
				return true, nil
			}
		}
	}

	rnr.Log.Info("Context role name >%s< does not exist", roleName)

	return false, nil
}

// addIdentityContext -
func (rnr *Runner) addContextIdentity(ctx context.Context, identity map[string]interface{}) (context.Context, error) {

	ctx = context.WithValue(ctx, contextIdentityKey, identity)

	return ctx, nil
}

// getIdentityContext -
func (rnr *Runner) getContextIdentity(ctx context.Context) (map[string]interface{}, error) {

	identity := ctx.Value(contextIdentityKey)
	if identity != nil {
		return identity.(map[string]interface{}), nil
	}

	rnr.Log.Info("Context identity is nil")

	return nil, nil
}

// getIdentityContext -
func (rnr *Runner) getContextIdentityValue(ctx context.Context, identityKey string) (interface{}, error) {

	identity, err := rnr.getContextIdentity(ctx)
	if err != nil {
		return nil, err
	}
	if identity != nil {
		for contextIdentityKey, contextIdentityValue := range identity {
			if contextIdentityKey == identityKey {
				rnr.Log.Info("Context identity key >%s< value >%v<", contextIdentityKey, contextIdentityValue)
				return contextIdentityValue, nil
			}
		}
	}

	rnr.Log.Info("Context identity key >%s< does not exist", identityKey)

	return nil, nil
}
