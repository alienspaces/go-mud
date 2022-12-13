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
func (rnr *Runner) Handler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp map[string]interface{}, l logger.Logger, m modeller.Modeller) error {
	l = Logger(l, "Handler")

	l.Info("Using Go M.U.D game handler")

	fmt.Fprint(w, "Hello from Go M.U.D!\n")

	return nil
}

// Router -
func (rnr *Runner) Router(r *httprouter.Router) error {
	l := Logger(rnr.Log, "Router")

	l.Info("Using Go M.U.D game router")

	return nil
}

// Middleware -
func (rnr *Runner) Middleware(h server.Handle) (server.Handle, error) {
	l := Logger(rnr.Log, "Middleware")

	l.Info("Using Go M.U.D game middleware")

	return h, nil
}
