package server

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	"gitlab.com/alienspaces/go-mud/server/core/type/logger"
	"gitlab.com/alienspaces/go-mud/server/core/type/modeller"
)

// Tx -
func (rnr *Runner) Tx(h Handle) (Handle, error) {

	handle := func(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp map[string]interface{}, l logger.Logger, m modeller.Modeller) error {

		// NOTE: The modeller is created and initialised with every request instead of
		// creating and assigning to a runner struct "Model" property at start up.
		// This prevents directly accessing a shared property from with the handler
		// function which is running in a goroutine. Otherwise accessing the "Model"
		// property would require locking and block simultaneous requests.

		l.Info("** Tx ** initialising database transaction")

		m, err := rnr.InitTx(l)
		if err != nil {
			l.Error("failed initialising database transaction, cannot authen >%v<", err)
			WriteSystemError(l, w, err)
			return err
		}

		err = h(w, r, pp, qp, l, m)
		if err != nil {
			l.Warn("** Tx ** rolling back database transaction")

			if err := m.Rollback(); err != nil {
				l.Warn("failed Tx rollback >%v<", err)
				return err
			}
			return err
		}

		l.Info("** Tx ** committing database transaction")

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
