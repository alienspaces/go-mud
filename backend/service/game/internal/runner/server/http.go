package runner

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"gitlab.com/alienspaces/go-mud/backend/core/server"
	"gitlab.com/alienspaces/go-mud/backend/core/type/logger"
	"gitlab.com/alienspaces/go-mud/backend/core/type/modeller"
)

// Handler - default handler
func (rnr *Runner) Handler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp map[string]interface{}, l logger.Logger, m modeller.Modeller) error {
	l = loggerWithContext(l, "Handler")

	l.Info("Using Go M.U.D game handler")

	fmt.Fprint(w, "Hello from Go M.U.D!\n")

	return nil
}

// Router -
func (rnr *Runner) Router(r *httprouter.Router) error {
	l := loggerWithContext(rnr.Log, "Router")

	l.Info("Using Go M.U.D game router")

	return nil
}

// Middleware -
func (rnr *Runner) Middleware(h server.Handle) (server.Handle, error) {
	l := loggerWithContext(rnr.Log, "Middleware")

	l.Info("Using Go M.U.D game middleware")

	return h, nil
}