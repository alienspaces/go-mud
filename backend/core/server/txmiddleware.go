package server

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	"gitlab.com/alienspaces/go-mud/backend/core/queryparam"
	"gitlab.com/alienspaces/go-mud/backend/core/type/logger"
	"gitlab.com/alienspaces/go-mud/backend/core/type/modeller"
)

// HeaderXTxRollback is used to rollback API transactions during testing.
//
// The only handler tests where the API Server must commit the DB transaction are those that contain further DB queries,
// which must be able to see changes made by the API server tx. This is due to the current default transaction isolation
// level (read committed).
const HeaderXTxRollback = "X-Tx-Rollback"

// TxMiddleware -
func (rnr *Runner) TxMiddleware(hc HandlerConfig, h Handle) (Handle, error) {

	handle := func(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, _ modeller.Modeller) error {
		l = Logger(l, "TxMiddleware")

		// NOTE: The modeller is created and initialised with every request instead of
		// creating and assigning to a runner struct "Model" property at start up.
		// This prevents directly accessing a shared property from with the handler
		// function which is running in a goroutine. Otherwise accessing the "Model"
		// property would require locking and block simultaneous requests.

		l.Debug("initialising database transaction")

		m, err := rnr.InitTx(l)
		if err != nil {
			l.Error("failed initialising database transaction >%v<", err)
			WriteSystemError(l, w, err)
			return err
		}

		err = h(w, r, pp, qp, l, m)
		if err != nil {
			l.Warn("rolling back database transaction")

			if err := m.Rollback(); err != nil {
				l.Warn("failed Tx rollback >%v<", err)
				return err
			}
			return err
		}

		l.Debug("committing database transaction")

		if r.Header.Get(HeaderXTxRollback) != "" {
			if err = m.Rollback(); err != nil {
				l.Warn("failed Tx commit >%v<", err)
				WriteSystemError(l, w, err)
				return err
			}
		} else if err = m.Commit(); err != nil {
			l.Warn("failed Tx commit >%v<", err)
			WriteSystemError(l, w, err)
			return err
		}

		return nil
	}

	return handle, nil
}
