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

// ParamMiddleware -
func (rnr *Runner) ParamMiddleware(hc HandlerConfig, h Handle) (Handle, error) {

	// Used to verify path parameters are resolved
	pathParams := extractPathParams(hc.Path)

	handle := func(w http.ResponseWriter, r *http.Request, pp httprouter.Params, _ *queryparam.QueryParams, l logger.Logger, m modeller.Modeller) error {
		l = Logger(l, "ParamMiddleware")

		l.Debug("Request method >%s< path >%s<", r.Method, r.RequestURI)

		queryParamValues := r.URL.Query()
		pathParamValues := map[string][]string{}

		if hc.MiddlewareConfig.ValidateParamsConfig != nil {

			cfg := hc.MiddlewareConfig.ValidateParamsConfig

			// Validate query parameters
			if cfg.QueryParamSchema != nil {
				l.Debug("Validating query parameters >%#v<", queryParamValues)
				err := validateParams(l, queryParamValues, cfg.QueryParamSchema)
				if err != nil {
					WriteError(l, w, coreerror.ProcessParamError(err))
					return err
				}
			}

			// Collect path parameter values
			for _, pathParam := range pathParams {
				pathParamValue := pp.ByName(pathParam)
				pathParamValues[pathParam] = []string{
					pathParamValue,
				}
				if !cfg.ExcludePathParamsFromQueryParams {
					l.Info("Adding path param >%s< Value >%s<", pathParam, pathParamValue)
					queryParamValues[pathParam] = []string{
						pathParamValue,
					}
				}
			}

			// Validate path parameters
			if cfg.PathParamSchema != nil {
				l.Debug("Validating path parameters >%#v<", pathParamValues)
				err := validateParams(l, pathParamValues, cfg.PathParamSchema)
				if err != nil {
					WriteError(l, w, coreerror.ProcessParamError(err))
					return err
				}
			}
		}

		// Build resulting query parameters
		qp, err := queryparam.BuildQueryParams(l, queryParamValues)
		if err != nil {
			l.Warn("failed to build query params >%#v< >%v<", r.URL.Query(), err)
			WriteError(l, w, err)
			return err
		}

		return h(w, r, pp, qp, l, m)
	}

	return handle, nil
}

func extractPathParams(p string) []string {
	parts := strings.Split(p, "/")
	params := []string{}
	for _, part := range parts {
		if strings.HasPrefix(part, ":") {
			params = append(params, strings.TrimPrefix(part, ":"))
		}
	}
	return params
}

// validateParams validates any provided parameters
func validateParams(l logger.Logger, q url.Values, paramSchema *jsonschema.SchemaWithReferences) error {
	if len(q) == 0 {
		return nil
	}

	if paramSchema == nil {
		return nil
	}

	qJSON := paramsToJSON(q)
	result, err := jsonschema.Validate(paramSchema, qJSON)
	if err != nil {
		l.Warn("failed validate params due to schema validation logic >%v<", err)
		return err
	}

	if !result.Valid() {
		err := coreerror.NewInvalidJSONError(result.Errors())
		l.Warn("failed validate params >%#v<", err)
		return err
	}

	l.Info("all parameters okay")

	return nil
}

func paramsToJSON(q url.Values) string {
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
	return `"` + strings.Split(k, ":")[0] + `"`
}

func parseValue(v string) string {
	// Attempt to parse an integer prior to a boolean as parse
	// bool will also accept 1 and 0 as valid booleans.
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
