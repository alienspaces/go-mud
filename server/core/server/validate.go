package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"

	coreerror "gitlab.com/alienspaces/go-mud/server/core/error"
	"gitlab.com/alienspaces/go-mud/server/core/jsonschema"
	"gitlab.com/alienspaces/go-mud/server/core/type/logger"
	"gitlab.com/alienspaces/go-mud/server/core/type/modeller"
)

// Validate -
func (rnr *Runner) Validate(hc HandlerConfig, h Handle) (Handle, error) {

	handle := func(w http.ResponseWriter, r *http.Request, pp httprouter.Params, _ map[string]interface{}, l logger.Logger, m modeller.Modeller) error {
		l = Logger(l, "Validate")

		l.Info("Request method >%s< path >%s<", r.Method, r.RequestURI)

		qp, err := validateQueryParameters(l, r.URL.Query(), hc.MiddlewareConfig.ValidateQueryParams)
		if err != nil {
			WriteError(l, w, coreerror.ProcessQueryParamError(err))
			return err
		}

		if r.Method != http.MethodPost && r.Method != http.MethodPut && r.Method != http.MethodPatch {
			l.Debug("Skipping validation of URI >%s< method >%s<", r.RequestURI, r.Method)
			return h(w, r, pp, qp, l, m)
		}

		requestSchema := hc.MiddlewareConfig.ValidateRequestSchema
		schemaMain := requestSchema.Main
		if schemaMain.Name == "" || schemaMain.Location == "" {
			l.Debug("Not validating URI >%s< method >%s<", r.RequestURI, r.Method)
			return h(w, r, pp, qp, l, m)
		}

		data := r.Context().Value(ctxKeyData)

		l.Info("Validation schemas >%#v<", requestSchema)
		l.Info("Validation data >%s<", data)

		result, err := jsonschema.Validate(requestSchema, data)
		if err != nil {
			l.Warn("failed validation >%v<", err)

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
			l.Warn("failed validation >%#v<", err)
			WriteError(l, w, err)
			return err
		}

		return h(w, r, pp, qp, l, m)
	}

	return handle, nil
}

// validateQueryParameters validates any provided parameters
func validateQueryParameters(l logger.Logger, q url.Values, paramSchema jsonschema.SchemaWithReferences) (map[string]interface{}, error) {
	if len(q) == 0 {
		return nil, nil
	}

	if paramSchema.IsEmpty() {
		for k := range q {
			return nil, coreerror.NewQueryParamError(fmt.Sprintf("Query parameter >%s< not allowed.", k))
		}
	}

	qJSON := queryParamsToJSON(q)
	result, err := jsonschema.Validate(paramSchema, qJSON)
	if err != nil {
		l.Warn("failed validate query params due to schema validation logic >%v<", err)
		return nil, err
	}

	if !result.Valid() {
		err := coreerror.NewSchemaValidationError(result.Errors())
		l.Warn("failed validate query params >%#v<", err)
		return nil, err
	}

	l.Info("all parameters okay")

	return buildQueryParams(q), nil
}

func queryParamsToJSON(q url.Values) string {
	if len(q) == 0 {
		return ""
	}

	jsonBuilder := strings.Builder{}
	jsonBuilder.WriteString("{")

	for k, v := range q {
		jsonBuilder.WriteString(parseKey(k))
		jsonBuilder.WriteString(":")

		switch len(v) {
		case 0:
			jsonBuilder.WriteString(`""`)
		case 1:
			jsonBuilder.WriteString(parseValue(v[0]))
		default:
			arr := "["

			for _, v := range v {
				arr += parseValue(v)
				arr += ","
			}

			arr = arr[0 : len(arr)-1] // remove extra comma
			arr += "]"
			jsonBuilder.WriteString(arr)
		}

		jsonBuilder.WriteString(",")
	}

	qpJSON := jsonBuilder.String()
	qpJSON = qpJSON[0 : len(qpJSON)-1] // remove extra comma
	qpJSON += "}"
	return qpJSON
}

func parseKey(k string) string {
	return `"` + k + `"`
}

func parseValue(v string) string {
	i, err := strconv.Atoi(v)
	if err != nil {
		return fmt.Sprintf("%#v", v)
	}

	return fmt.Sprintf("%#v", i)
}

func buildQueryParams(q url.Values) map[string]interface{} {
	if len(q) == 0 {
		return nil
	}

	qp := map[string]interface{}{}

	for key, value := range q {
		if len(value) == 0 {
			continue
		}

		qp[key] = value
	}

	return qp
}
