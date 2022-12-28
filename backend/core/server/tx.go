package server

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	"gitlab.com/alienspaces/go-mud/backend/core/type/logger"
	"gitlab.com/alienspaces/go-mud/backend/core/type/modeller"
)

// Tx -
func (rnr *Runner) Tx(h Handle) (Handle, error) {

	handle := func(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp map[string]interface{}, l logger.Logger, _ modeller.Modeller) error {
		l = HTTPLogger(l, "Tx")

		// NOTE: The modeller is created and initialised with every request instead of
		// creating and assigning to a runner struct "Model" property at start up.
		// This prevents directly accessing a shared property from with the handler
		// function which is running in a goroutine. Otherwise accessing the "Model"
		// property would require locking and block simultaneous requests.

		l.Info("Initialising database transaction")

		// Force a transaction rollback
		xTxRollback := r.Header.Get("X-Tx-Rollback")

		m, err := rnr.InitModeller(l)
		if err != nil {
			l.Error("failed initialising database transaction, cannot authen >%v<", err)
			WriteSystemError(l, w, err)
			return err
		}

		err = h(w, r, pp, qp, l, m)
		if err != nil || xTxRollback == "true" {
			l.Info("Rolling back database transaction")

			if err := m.Rollback(); err != nil {
				l.Warn("failed Tx rollback >%v<", err)
				return err
			}
			return err
		}

		l.Info("Committing database transaction")

		err = m.Commit()
		if err != nil {
			l.Warn("failed Tx commit >%v<", err)
			WriteSystemError(l, w, err)
			return err
		}

		return nil
	}

	return handle, nil
}
