package server

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/julienschmidt/httprouter"

	coreerror "gitlab.com/alienspaces/go-mud/backend/core/error"
	"gitlab.com/alienspaces/go-mud/backend/core/jsonschema"
	"gitlab.com/alienspaces/go-mud/backend/core/queryparam"
	"gitlab.com/alienspaces/go-mud/backend/core/type/logger"
	"gitlab.com/alienspaces/go-mud/backend/core/type/modeller"
)

// DataMiddleware -
func (rnr *Runner) DataMiddleware(hc HandlerConfig, h Handle) (Handle, error) {

	handle := func(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m modeller.Modeller) error {
		l = Logger(l, "DataMiddleware")

		if r.Method == http.MethodGet {
			return h(w, r, pp, qp, l, m)
		}

		// read body into a string
		buf := new(bytes.Buffer)
		if r.Body != nil {
			bytes, err := buf.ReadFrom(r.Body)
			if err != nil {
				l.Warn("failed reading data buffer >%v<", err)
				WriteSystemError(l, w, err)
				return err
			}
			l.Debug("Read >%d< bytes ** ", bytes)
		}
		data := buf.String()

		if r.Method != http.MethodPost && r.Method != http.MethodPut && r.Method != http.MethodPatch {
			l.Debug("skipping validation of URI >%s< method >%s<", r.RequestURI, r.Method)
			return h(w, r, pp, qp, l, m)
		}

		requestSchema := hc.MiddlewareConfig.ValidateRequestSchema
		schemaMain := requestSchema.Main
		if schemaMain.Name == "" || schemaMain.Location == "" {
			l.Debug("not validating URI >%s< method >%s<", r.RequestURI, r.Method)
			return h(w, r, pp, qp, l, m)
		}

		l.Info("Schemas >%#v<", requestSchema)
		l.Info("Data >%s<", data)

		result, err := jsonschema.Validate(requestSchema, data)
		if err != nil {
			l.Warn("failed validate >%v<", err)

			var jsonSyntaxError *json.SyntaxError
			if errors.As(err, &jsonSyntaxError) || errors.Is(err, io.ErrUnexpectedEOF) {
				WriteError(l, w, coreerror.NewInvalidBodyError(""))
			} else if errors.Is(err, io.EOF) {
				WriteError(l, w, coreerror.NewInvalidBodyError("Request body is empty."))
			} else {
				WriteSystemError(l, w, err)
			}

			return err
		}

		if !result.Valid() {
			err := coreerror.NewSchemaValidationError(result.Errors())
			l.Warn("failed validate >%#v<", err)
			WriteError(l, w, err)
			return err
		}

		// Add data to context
		ctx := context.WithValue(r.Context(), ctxKeyData, data)

		// delegate request
		return h(w, r.WithContext(ctx), pp, qp, l, m)
	}

	return handle, nil
}
