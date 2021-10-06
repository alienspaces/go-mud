package server

import (
	"bytes"
	"context"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"gitlab.com/alienspaces/go-boilerplate/server/core/type/logger"
	"gitlab.com/alienspaces/go-boilerplate/server/core/type/modeller"
)

// Data -
func (rnr *Runner) Data(hc HandlerConfig, h HandlerFunc) (HandlerFunc, error) {

	handle := func(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp map[string]interface{}, l logger.Logger, m modeller.Modeller) {

		// read body into a string
		buf := new(bytes.Buffer)
		if r.Body != nil {
			bytes, err := buf.ReadFrom(r.Body)
			if err != nil {
				l.Warn("Failed reading data buffer >%v<", err)
				http.Error(w, "Server Error", http.StatusInternalServerError)
				return
			}
			l.Debug("** Data read >%d< bytes ** ", bytes)
		}
		data := buf.String()

		// Add data to context
		ctx := context.WithValue(r.Context(), ContextKeyData, data)

		// delegate request
		h(w, r.WithContext(ctx), pp, qp, l, m)
	}

	return handle, nil
}
