package server

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"

	"gitlab.com/alienspaces/go-mud/backend/core/queryparam"
	"gitlab.com/alienspaces/go-mud/backend/core/type/logger"
	"gitlab.com/alienspaces/go-mud/backend/core/type/modeller"
)

// CorrelationMiddleware -
func (rnr *Runner) CorrelationMiddleware(hc HandlerConfig, h Handle) (Handle, error) {

	handle := func(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, _ modeller.Modeller) error {
		l = Logger(l, "CorrelationMiddleware")

		correlationID := r.Header.Get("X-Correlation-ID")
		if correlationID == "" {
			correlationID = uuid.New().String()
			l.Debug("Generated correlation ID >%s<", correlationID)
		}
		w.Header().Set("X-Correlation-ID", correlationID)

		ctx := r.Context()
		ctx = context.WithValue(ctx, ctxKeyCorrelationID, correlationID)
		r = r.WithContext(ctx)

		l.Context("correlation-id", correlationID)

		// NOTE: Log every request here at info log level
		l.Info("Request method >%s< path >%s<", r.Method, r.RequestURI)

		return h(w, r, pp, nil, l, nil)
	}

	return handle, nil
}
