package server

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"gitlab.com/alienspaces/go-mud/server/core/schema"
	"gitlab.com/alienspaces/go-mud/server/core/type/logger"
	"gitlab.com/alienspaces/go-mud/server/core/type/modeller"
)

// pathParamCache - path, method, []string
var pathParamCache map[string]map[string]map[string]ValidatePathParam

// queryParamCache - path, method, []string
var queryParamCache map[string]map[string][]string

var schemaValidator *schema.Validator

// Validate -
func (rnr *Runner) Validate(hc HandlerConfig, h HandlerFunc) (HandlerFunc, error) {

	rnr.Log.Info("** Validate ** cache query param lists")

	// cache query parameter lists
	err := rnr.validateCacheQueryParams(hc)
	if err != nil {
		rnr.Log.Warn("Failed caching query param list >%v<", err)
		return nil, err
	}

	rnr.Log.Info("** Validate ** cache schemas")

	// Cache path parameter validations
	err = rnr.validateCachePathParams(hc)
	if err != nil {
		rnr.Log.Warn("Failed caching path param list >%v<", err)
		return nil, err
	}

	// Validation schema configuration
	schemaConfig := schema.Config{
		Location:   hc.MiddlewareConfig.ValidateSchemaLocation,
		Main:       hc.MiddlewareConfig.ValidateSchemaRequestMain,
		References: hc.MiddlewareConfig.ValidateSchemaRequestReferences,
	}

	// Load configured schemas
	if hc.MiddlewareConfig.ValidateSchemaLocation != "" && hc.MiddlewareConfig.ValidateSchemaRequestMain != "" {
		err = rnr.validateCacheSchemas(schemaConfig)
		if err != nil {
			rnr.Log.Warn("Failed loading schemas >%v<", err)
			return nil, err
		}
	}

	handle := func(w http.ResponseWriter, r *http.Request, pp httprouter.Params, _ map[string]interface{}, l logger.Logger, m modeller.Modeller) {

		l.Debug("** Validate ** request URI >%s< method >%s<", r.RequestURI, r.Method)

		// validate path parameters
		err := rnr.validatePathParameters(l, hc.Path, r, pp)
		if err != nil {
			rnr.WriteResponse(l, w, rnr.ValidationError(err))
			return
		}

		// validate query parameters
		qp, err := rnr.validateQueryParameters(l, hc.Path, r)
		if err != nil {
			rnr.WriteResponse(l, w, rnr.ValidationError(err))
			return
		}

		// only validate requests where the method provides data
		if r.Method != http.MethodPost && r.Method != http.MethodPut {
			l.Debug("Skipping validation of URI >%s< method >%s<", r.RequestURI, r.Method)

			// delegate request
			h(w, r, pp, qp, l, m)
			return
		}

		if hc.MiddlewareConfig.ValidateSchemaLocation != "" && hc.MiddlewareConfig.ValidateSchemaRequestMain != "" {
			if schemaValidator == nil {
				rnr.WriteResponse(l, w, rnr.SystemError(fmt.Errorf("schema validator not available, cannot validate")))
				return
			}
		}

		if !schemaValidator.SchemaCached(schemaConfig) {
			l.Info("Not validating URI >%s< method >%s<", r.RequestURI, r.Method)

			// delegate request
			h(w, r, pp, qp, l, m)
			return
		}

		// data from context
		data := r.Context().Value(ContextKeyData)

		l.Debug("Data >%s<", data)
		switch data := data.(type) {
		case string:
			err = schemaValidator.Validate(
				schemaConfig,
				data,
			)
			if err != nil {
				rnr.WriteResponse(l, w, rnr.ValidationError(err))
				return
			}
		case nil:
			l.Warn("Data is nil")
			rnr.WriteResponse(l, w, rnr.SystemError(fmt.Errorf("Data is nil")))
			return
		default:
			l.Warn("Data type unsupported")
			rnr.WriteResponse(l, w, rnr.SystemError(fmt.Errorf("Data is nil")))
			return
		}

		// delegate request
		h(w, r, pp, qp, l, m)
	}

	return handle, nil
}

// validatePathParameters - validate path parameters
func (rnr *Runner) validatePathParameters(l logger.Logger, path string, r *http.Request, pp httprouter.Params) error {

	if len(pathParamCache) == 0 {
		l.Debug("Handler method >%s< path >%s< not configured with path param validations", r.Method, path)
		return nil
	}

	// Request context may be used to validate path parameter values
	ctx := r.Context()

	// Cached path parameter validation configuration
	params := pathParamCache[path][r.Method]

VALIDATE_PARAMS:
	for pathParamName, pathParamConfig := range params {
		// Provided path parameters
		for _, pathParam := range pp {
			if pathParam.Key == pathParamName {
				if pathParamConfig.MatchIdentity {
					identityValue, err := rnr.getContextIdentityValue(ctx, pathParam.Key)
					if err != nil {
						l.Warn("Failed getting context identity value >%s< >%v<", pathParam.Key, err)
						return err
					}
					if identityValue != pathParam.Value {
						msg := fmt.Sprintf("Mismatched path param value >%v< versus identity value >%v<", pathParam.Value, identityValue)
						l.Warn(msg)
						return fmt.Errorf(msg)
					}
					l.Info("Matched path param config name >%s< value >%v< with identity value >%v<", pathParamName, pathParam.Value, identityValue)
					continue VALIDATE_PARAMS
				}
			}
		}
		msg := fmt.Sprintf("Path param name >%s< not found in params", pathParamName)
		l.Warn(msg)
		return fmt.Errorf(msg)
	}

	return nil
}

// validateQueryParameters - validate query parameters
func (rnr *Runner) validateQueryParameters(l logger.Logger, path string, r *http.Request) (map[string]interface{}, error) {

	if len(queryParamCache) == 0 {
		l.Debug("Handler method >%s< path >%s< not configured with query params list", r.Method, path)
		return nil, nil
	}

	// Query parameters
	q := r.URL.Query()

	// Valid query parameters
	qp := map[string]interface{}{}

	// Cached allowed query parameter names
	params := queryParamCache[path][r.Method]

	for paramName, paramValue := range q {
		l.Info("Checking parameter >%s< >%s<", paramName, paramValue)

		found := false
		for _, param := range params {
			if paramName == param {
				found = true
			}
		}
		if !found {
			msg := fmt.Sprintf("Parameter >%s< not allowed", paramName)
			l.Warn(msg)
			return nil, fmt.Errorf(msg)
		}
		qp[paramName] = paramValue
	}

	l.Info("All parameters okay")

	return qp, nil
}

// validateCachePathParams -
func (rnr *Runner) validateCachePathParams(hc HandlerConfig) error {

	if len(hc.MiddlewareConfig.ValidatePathParams) == 0 {
		rnr.Log.Info("Handler method >%s< path >%s< not configured with path params list", hc.Method, hc.Path)
		return nil
	}

	rnr.Log.Info("Handler method >%s< path >%s< has path param validations >%#v<", hc.Method, hc.Path, hc.MiddlewareConfig.ValidatePathParams)

	if pathParamCache == nil {
		pathParamCache = map[string]map[string]map[string]ValidatePathParam{}
	}
	if pathParamCache[hc.Path] == nil {
		pathParamCache[hc.Path] = make(map[string]map[string]ValidatePathParam)
	}
	pathParamCache[hc.Path][hc.Method] = hc.MiddlewareConfig.ValidatePathParams

	return nil
}

// validateCacheQueryParams -
func (rnr *Runner) validateCacheQueryParams(hc HandlerConfig) error {

	if len(hc.MiddlewareConfig.ValidateQueryParams) == 0 {
		rnr.Log.Info("Handler method >%s< path >%s< not configured with query params list", hc.Method, hc.Path)
		return nil
	}

	rnr.Log.Info("Handler method >%s< path >%s< has query params list >%#v<", hc.Method, hc.Path, hc.MiddlewareConfig.ValidateQueryParams)

	if queryParamCache == nil {
		queryParamCache = map[string]map[string][]string{}
	}
	if queryParamCache[hc.Path] == nil {
		queryParamCache[hc.Path] = make(map[string][]string)
	}
	queryParamCache[hc.Path][hc.Method] = hc.MiddlewareConfig.ValidateQueryParams

	return nil
}

// validateCacheSchemas - load validation JSON schemas
func (rnr *Runner) validateCacheSchemas(schemaConfig schema.Config) error {

	if schemaValidator == nil {
		rnr.Log.Info("Schema validator is nil, creating new schema validator")
		var err error
		schemaValidator, err = schema.NewValidator(rnr.Config, rnr.Log)
		if err != nil {
			rnr.Log.Warn("Failed new schema validator >%v<", err)
			return err
		}
	}

	_, err := schemaValidator.LoadSchema(schemaConfig)
	if err != nil {
		rnr.Log.Warn("Failed loading schema >%v<", err)
		return err
	}

	return nil
}
