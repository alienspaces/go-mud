package server

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"

	"gitlab.com/alienspaces/go-mud/backend/core/type/logger"
	"gitlab.com/alienspaces/go-mud/backend/core/type/modeller"
)

// Correlation -
func (rnr *Runner) Correlation(h Handle) (Handle, error) {

	handle := func(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp map[string]interface{}, l logger.Logger, _ modeller.Modeller) error {
		l = loggerWithContext(l, "Correlation")

		path := r.RequestURI
		if path == "" && r.URL != nil {
			path = r.URL.Path
		}

		l.Info("Request method >%s< path >%s<", r.Method, path)

		correlationID := r.Header.Get("X-Correlation-ID")
		if correlationID == "" {
			correlationID = uuid.New().String()
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, ctxKeyCorrelationID, correlationID)
		r = r.WithContext(ctx)

		l.Context("correlation-id", correlationID)

		return h(w, r, pp, nil, l, nil)
	}

	return handle, nil
}
