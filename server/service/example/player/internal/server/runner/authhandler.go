package runner

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	"gitlab.com/alienspaces/go-boilerplate/server/constant"
	"gitlab.com/alienspaces/go-boilerplate/server/core/auth"
	"gitlab.com/alienspaces/go-boilerplate/server/core/type/logger"
	"gitlab.com/alienspaces/go-boilerplate/server/core/type/modeller"
	"gitlab.com/alienspaces/go-boilerplate/server/schema"
	"gitlab.com/alienspaces/go-boilerplate/server/service/player/internal/model"
)

// PostAuthHandler -
func (rnr *Runner) PostAuthHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp map[string]interface{}, l logger.Logger, m modeller.Modeller) {

	l.Info("** Post auth handler ** p >%#v< m >#%v<", pp, m)

	req := schema.AuthRequest{}

	err := rnr.ReadRequest(l, r, &req)
	if err != nil {
		rnr.WriteSystemError(l, w, err)
		return
	}

	accountRec, err := m.(*model.Model).AuthVerify(model.AuthData{
		Provider:          req.Data.Provider,
		ProviderAccountID: req.Data.ProviderAccountID,
		ProviderToken:     req.Data.ProviderToken,
		PlayerEmail:       req.Data.PlayerEmail,
		PlayerName:        req.Data.PlayerName,
	})
	if err != nil {
		rnr.WriteUnauthorizedError(l, w, err)
		return
	}

	a, err := auth.NewAuth(rnr.Config, rnr.Log)
	if err != nil {
		rnr.WriteUnauthorizedError(l, w, err)
		return
	}

	// TODO: Expand on account roles
	roles := []string{
		constant.AuthRoleDefault,
	}

	identity := map[string]interface{}{
		constant.AuthIdentityPlayerID: accountRec.ID,
	}

	claims := auth.Claims{
		Roles:    roles,
		Identity: identity,
	}

	tokenString, err := a.EncodeJWT(&claims)
	if err != nil {
		rnr.WriteUnauthorizedError(l, w, err)
		return
	}

	// assign response properties
	res := schema.AuthResponse{
		Data: []schema.AuthData{
			{
				Provider:          req.Data.Provider,
				ProviderAccountID: req.Data.ProviderAccountID,
				ProviderToken:     req.Data.ProviderToken,
				PlayerID:          accountRec.ID,
				PlayerName:        accountRec.Name,
				PlayerEmail:       accountRec.Email,
				Token:             tokenString,
				CreatedAt:         accountRec.CreatedAt,
				UpdatedAt:         accountRec.UpdatedAt.Time,
			},
		},
	}

	err = rnr.WriteResponse(l, w, res)
	if err != nil {
		l.Warn("Failed writing response >%v<", err)
		return
	}
}
