package runner

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"gitlab.com/alienspaces/go-mud/server/core/server"
	"gitlab.com/alienspaces/go-mud/server/core/type/logger"
	"gitlab.com/alienspaces/go-mud/server/core/type/modeller"
)

// Handler - default handler
func (rnr *Runner) Handler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp map[string]interface{}, l logger.Logger, m modeller.Modeller) {

	l.Info("** Game handler **")

	fmt.Fprint(w, "Hello from game!\n")
}

// Router -
func (rnr *Runner) Router(r *httprouter.Router) error {

	rnr.Log.Info("** Game Router **")

	return nil
}

// Middleware -
func (rnr *Runner) Middleware(h server.HandlerFunc) (server.HandlerFunc, error) {

	rnr.Log.Info("** Game Middleware **")

	return h, nil
}
