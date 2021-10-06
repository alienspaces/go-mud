package server

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"

	"gitlab.com/alienspaces/go-boilerplate/server/core/type/logger"
	"gitlab.com/alienspaces/go-boilerplate/server/core/type/modeller"
)

// Correlation -
func (rnr *Runner) Correlation(h HandlerFunc) (HandlerFunc, error) {

	handle := func(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp map[string]interface{}, l logger.Logger, _ modeller.Modeller) {

		lc, err := l.NewInstance()
		if err != nil {
			rnr.Log.Warn("Failed new log instance >%v<", err)
			http.Error(w, "Server Error", http.StatusInternalServerError)
			return
		}

		correlationID := r.Header.Get("X-Correlation-ID")
		if correlationID == "" {
			correlationID = uuid.New().String()
			lc.Debug("Generated correlation ID >%s<", correlationID)
		}
		lc.Context("correlation-id", correlationID)

		// NOTE: Log every request here at info log level
		lc.Info("Request Method >%s< Path >%s<", r.Method, r.RequestURI)

		h(w, r, pp, nil, lc, nil)
	}

	return handle, nil
}
