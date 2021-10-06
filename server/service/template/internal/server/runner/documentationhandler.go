package runner

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	"gitlab.com/alienspaces/go-boilerplate/server/core/type/logger"
	"gitlab.com/alienspaces/go-boilerplate/server/core/type/modeller"
)

// GetDocumentationHandler -
func (rnr *Runner) GetDocumentationHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp map[string]interface{}, l logger.Logger, m modeller.Modeller) {

	l.Info("** Get documentation handler ** p >%#v< m >%#v<", pp, m)

	// content type html
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// status
	w.WriteHeader(http.StatusOK)

	documentation, err := rnr.GenerateHandlerDocumentation()
	if err != nil {
		l.Warn("Failed getting handler documentation >%v<", err)

		// system error
		res := rnr.SystemError(err)

		err = rnr.WriteResponse(l, w, res)
		if err != nil {
			l.Warn("Failed writing response >%v<", err)
			return
		}
		return
	}

	// content
	written, err := w.Write(documentation)
	if err != nil {
		l.Warn("Failed writing documentation >%v<", err)
		return
	}

	l.Info("Wrote >%s< bytes", written)
}
