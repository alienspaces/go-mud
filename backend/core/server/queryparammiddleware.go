package server

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"

	coreerror "gitlab.com/alienspaces/go-mud/backend/core/error"
	"gitlab.com/alienspaces/go-mud/backend/core/jsonschema"
	"gitlab.com/alienspaces/go-mud/backend/core/queryparam"
	"gitlab.com/alienspaces/go-mud/backend/core/type/logger"
	"gitlab.com/alienspaces/go-mud/backend/core/type/modeller"
)

// QueryParamMiddleware -
func (rnr *Runner) QueryParamMiddleware(hc HandlerConfig, h Handle) (Handle, error) {

	handle := func(w http.ResponseWriter, r *http.Request, pp httprouter.Params, _ *queryparam.QueryParams, l logger.Logger, m modeller.Modeller) error {
		l = Logger(l, "QueryParamMiddleware")

		l.Debug("Request method >%s< path >%s<", r.Method, r.RequestURI)

		err := validateQueryParams(l, r.URL.Query(), hc.MiddlewareConfig.ValidateQueryParams)
		if err != nil {
			WriteError(l, w, coreerror.ProcessQueryParamError(err))
			return err
		}

		qp, err := buildQueryParams(l, r.URL.Query())
		if err != nil {
			l.Warn("failed to build query params >%#v< >%v<", r.URL.Query(), err)
			WriteError(l, w, err)
			return err
		}

		return h(w, r, pp, qp, l, m)
	}

	return handle, nil
}

// validateQueryParams validates any provided parameters
func validateQueryParams(l logger.Logger, q url.Values, paramSchema jsonschema.SchemaWithReferences) error {
	if len(q) == 0 {
		return nil
	}

	if paramSchema.IsEmpty() {
		for k := range q {
			return coreerror.NewQueryParamError(fmt.Sprintf("Query parameter >%s< not allowed.", k))
		}
	}

	qJSON := queryParamsToJSON(q)
	result, err := jsonschema.Validate(paramSchema, qJSON)
	if err != nil {
		l.Warn("failed validate query params due to schema validation logic >%v<", err)
		return err
	}

	if !result.Valid() {
		err := coreerror.NewSchemaValidationError(result.Errors())
		l.Warn("failed validate query params >%#v<", err)
		return err
	}

	l.Info("all parameters okay")

	return nil
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

		if strings.HasSuffix(k, "[]") {
			arr := "["

			for _, v := range v {
				arr += parseValue(v)
				arr += ","
			}

			if len(arr) > 1 {
				arr = arr[:len(arr)-1] // remove extra comma
			}
			arr += "]"
			jsonBuilder.WriteString(arr)
		} else if len(v) == 0 {
			jsonBuilder.WriteString(`""`)
		} else {
			jsonBuilder.WriteString(parseValue(v[0]))
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
	// Attempt to parse an integer prior to
	// a boolean as parse bool will also accept
	// 1 and 0 as valid booleans
	i, err := strconv.Atoi(v)
	if err != nil {
		b, err := strconv.ParseBool(v)
		if err != nil {
			return fmt.Sprintf("%#v", v)
		}
		return fmt.Sprintf("%#v", b)
	}
	return fmt.Sprintf("%#v", i)
}
