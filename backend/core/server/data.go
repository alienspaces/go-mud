package server

import (
	"bytes"
	"context"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"gitlab.com/alienspaces/go-mud/backend/core/type/logger"
	"gitlab.com/alienspaces/go-mud/backend/core/type/modeller"
)

// Data -
func (rnr *Runner) Data(h Handle) (Handle, error) {

	handle := func(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp map[string]interface{}, l logger.Logger, m modeller.Modeller) error {

		// read body into a string
		buf := new(bytes.Buffer)
		if r.Body != nil {
			bytes, err := buf.ReadFrom(r.Body)
			if err != nil {
				l.Warn("failed reading data buffer >%v<", err)
				WriteSystemError(l, w, err)
				return err
			}
			l.Debug("** Data read >%d< bytes ** ", bytes)
		}
		data := buf.String()

		// Add data to context
		ctx := context.WithValue(r.Context(), ctxKeyData, data)

		// delegate request
		return h(w, r.WithContext(ctx), pp, qp, l, m)
	}

	return handle, nil
}
