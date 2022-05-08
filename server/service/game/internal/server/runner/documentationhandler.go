package runner

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"gitlab.com/alienspaces/go-mud/server/core/type/logger"
	"gitlab.com/alienspaces/go-mud/server/core/type/modeller"
)

// GetDocumentationHandler -
func (rnr *Runner) GetDocumentationHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp map[string]interface{}, l logger.Logger, m modeller.Modeller) error {

	l.Info("** get schema documentation handler ** p >%#v< m >%#v<", pp, m)

	docs, err := rnr.GenerateHandlerDocumentation(rnr.GetMessageConfigs(), rnr.GetHandlerConfigs())
	if err != nil {
		msg := fmt.Sprintf("unable to load schema documentation >%v<, cannot init runner", err)
		rnr.Log.Error(msg)
		return fmt.Errorf(msg)
	}

	// content type html
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// status
	w.WriteHeader(http.StatusOK)

	// content
	written, err := w.Write(docs)
	if err != nil {
		l.Warn("failed writing documentation >%v<", err)
		return err
	}

	l.Info("wrote >%d< bytes", written)

	return nil
}
