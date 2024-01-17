package runner

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	coreerror "gitlab.com/alienspaces/go-mud/backend/core/error"
	"gitlab.com/alienspaces/go-mud/backend/core/jsonschema"
	"gitlab.com/alienspaces/go-mud/backend/core/queryparam"
	"gitlab.com/alienspaces/go-mud/backend/core/server"
	"gitlab.com/alienspaces/go-mud/backend/core/type/logger"
	"gitlab.com/alienspaces/go-mud/backend/core/type/modeller"
	schema "gitlab.com/alienspaces/go-mud/backend/schema/template"
	"gitlab.com/alienspaces/go-mud/backend/service/template/internal/model"
	"gitlab.com/alienspaces/go-mud/backend/service/template/internal/record"
)

const (
	getTemplates       string = "get-templates"
	getTemplate        string = "get-template"
	postTemplate       string = "post-template"
	postTemplateWithID string = "post-template-with-id"
	putTemplate        string = "put-template"
)

func (rnr *Runner) TemplateHandlerConfig(hc map[string]server.HandlerConfig) map[string]server.HandlerConfig {

	requestSchema := &jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "schema/template",
			Name:     "main.request.schema.json",
		},
		References: []jsonschema.Schema{
			{
				Location: "schema/template",
				Name:     "data.request.schema.json",
			},
		},
	}

	responseSchema := &jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "schema/template",
			Name:     "main.response.schema.json",
		},
		References: []jsonschema.Schema{
			{
				Location: "schema/template",
				Name:     "data.response.schema.json",
			},
		},
	}

	return mergeHandlerConfigs(hc, map[string]server.HandlerConfig{
		getTemplates: {
			Method:      http.MethodGet,
			Path:        "/api/templates",
			HandlerFunc: rnr.getTemplatesHandler,
			MiddlewareConfig: server.MiddlewareConfig{
				AuthenTypes: []server.AuthenticationType{
					server.AuthenticationTypePublic,
				},
				ValidateResponseSchema: responseSchema,
			},
			DocumentationConfig: server.DocumentationConfig{
				Document:    true,
				Description: "Query templates.",
			},
		},
		getTemplate: {
			Method:      http.MethodGet,
			Path:        "/api/templates/:template_id",
			HandlerFunc: rnr.getTemplateHandler,
			MiddlewareConfig: server.MiddlewareConfig{
				AuthenTypes: []server.AuthenticationType{
					server.AuthenticationTypePublic,
				},
				ValidateResponseSchema: responseSchema,
			},
			DocumentationConfig: server.DocumentationConfig{
				Document:    true,
				Description: "Get a template.",
			},
		},
		postTemplate: {
			Method:      http.MethodPost,
			Path:        "/api/templates",
			HandlerFunc: rnr.postTemplatesHandler,
			MiddlewareConfig: server.MiddlewareConfig{
				AuthenTypes: []server.AuthenticationType{
					server.AuthenticationTypePublic,
				},
				ValidateRequestSchema:  requestSchema,
				ValidateResponseSchema: responseSchema,
			},
			DocumentationConfig: server.DocumentationConfig{
				Document:    true,
				Description: "Create a template.",
			},
		},
		postTemplateWithID: {
			Method:      http.MethodPost,
			Path:        "/api/templates/:template_id",
			HandlerFunc: rnr.postTemplatesHandler,
			MiddlewareConfig: server.MiddlewareConfig{
				AuthenTypes: []server.AuthenticationType{
					server.AuthenticationTypePublic,
				},
				ValidateRequestSchema:  requestSchema,
				ValidateResponseSchema: responseSchema,
			},
			DocumentationConfig: server.DocumentationConfig{
				Document:    true,
				Description: "Create a template.",
			},
		},
		putTemplate: {
			Method:      http.MethodPut,
			Path:        "/api/templates/:template_id",
			HandlerFunc: rnr.putTemplatesHandler,
			MiddlewareConfig: server.MiddlewareConfig{
				AuthenTypes: []server.AuthenticationType{
					server.AuthenticationTypePublic,
				},
				ValidateRequestSchema:  requestSchema,
				ValidateResponseSchema: responseSchema,
			},
			DocumentationConfig: server.DocumentationConfig{
				Document:    true,
				Description: "Update a template.",
			},
		},
	})
}

// getTemplateHandler -
func (rnr *Runner) getTemplateHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m modeller.Modeller) error {
	l = loggerWithFunctionContext(l, "getTemplateHandler")
	l.Info("** Get template handler **")

	id := pp.ByName("template_id")

	if id == "" {
		err := coreerror.NewNotFoundError("template", id)
		server.WriteError(l, w, err)
		return err
	}

	l.Info("Getting template ID >%s<", id)

	rec, err := m.(*model.Model).GetTemplateRec(id, nil)
	if err != nil {
		l.Warn("failed getting template record >%v<", err)
		server.WriteError(l, w, err)
		return err
	}

	// Resource not found
	if rec == nil {
		err := coreerror.NewNotFoundError("template", id)
		server.WriteError(l, w, err)
		return err
	}

	// response data
	data, err := rnr.RecordToTemplateResponseData(rec)
	if err != nil {
		l.Warn("failed mapping template response >%v<", err)
		server.WriteError(l, w, err)
		return err
	}

	res := schema.Response{
		Data: data,
	}

	err = server.WriteResponse(l, w, http.StatusOK, res)
	if err != nil {
		l.Warn("Failed writing response >%v<", err)
		return err
	}

	return nil
}

// getTemplatesHandler -
func (rnr *Runner) getTemplatesHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m modeller.Modeller) error {
	l = loggerWithFunctionContext(l, "getTemplatesHandler")

	opts := queryparam.ToSQLOptions(qp)

	l.Info("Querying dungeon records with opts >%#v<", opts)

	recs, err := m.(*model.Model).GetTemplateRecs(opts)
	if err != nil {
		server.WriteError(l, w, err)
		return err
	}

	// assign response properties
	data := []*schema.Data{}
	for _, rec := range recs {

		// response data
		responseData, err := rnr.RecordToTemplateResponseData(rec)
		if err != nil {
			server.WriteSystemError(l, w, err)
			return err
		}

		data = append(data, responseData)
	}

	res := schema.CollectionResponse{
		Data: data,
	}

	err = server.WriteResponse(l, w, http.StatusOK, res)
	if err != nil {
		l.Warn("Failed writing response >%v<", err)
		return err
	}

	return nil
}

// postTemplatesHandler -
func (rnr *Runner) postTemplatesHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m modeller.Modeller) error {
	l.Info("Post templates handler pp >%#v< qp >%#v<", pp, qp)

	// parameters
	id := pp.ByName("template_id")

	req, err := server.ReadRequest(l, r, &schema.Request{})
	if err != nil {
		server.WriteSystemError(l, w, err)
		return err
	}

	rec := &record.Template{}

	// assign request properties
	rec.ID = id

	// record data
	rec, err = rnr.TemplateRequestDataToRecord(req.Data)
	if err != nil {
		server.WriteSystemError(l, w, err)
		return err
	}

	rec, err = m.(*model.Model).CreateTemplateRec(rec)
	if err != nil {
		server.WriteError(l, w, err)
		return err
	}

	// response data
	responseData, err := rnr.RecordToTemplateResponseData(rec)
	if err != nil {
		server.WriteSystemError(l, w, err)
		return err
	}

	// assign response properties
	res := schema.Response{
		Data: responseData,
	}

	err = server.WriteResponse(l, w, http.StatusCreated, res)
	if err != nil {
		l.Warn("Failed writing response >%v<", err)
		return err
	}

	return nil
}

// putTemplatesHandler -
func (rnr *Runner) putTemplatesHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m modeller.Modeller) error {
	l.Info("Put templates handler pp >%#v< qp >%#v<", pp, qp)

	// parameters
	id := pp.ByName("template_id")

	l.Info("Updating resource ID >%s<", id)

	_, err := m.(*model.Model).GetTemplateRec(id, nil)
	if err != nil {
		server.WriteError(l, w, err)
		return err
	}

	req, err := server.ReadRequest(l, r, &schema.Request{})
	if err != nil {
		server.WriteSystemError(l, w, err)
		return err
	}

	// record data
	rec, err := rnr.TemplateRequestDataToRecord(req.Data)
	if err != nil {
		server.WriteSystemError(l, w, err)
		return err
	}

	rec.ID = id

	rec, err = m.(*model.Model).UpdateTemplateRec(rec)
	if err != nil {
		server.WriteError(l, w, err)
		return err
	}

	// response data
	responseData, err := rnr.RecordToTemplateResponseData(rec)
	if err != nil {
		server.WriteSystemError(l, w, err)
		return err
	}

	// assign response properties
	res := schema.Response{
		Data: responseData,
	}

	err = server.WriteResponse(l, w, http.StatusOK, res)
	if err != nil {
		l.Warn("Failed writing response >%v<", err)
		return err
	}

	return nil
}

// TemplateRequestDataToRecord -
func (rnr *Runner) TemplateRequestDataToRecord(data *schema.Data) (*record.Template, error) {
	return &record.Template{}, nil
}

// RecordToTemplateResponseData -
func (rnr *Runner) RecordToTemplateResponseData(templateRec *record.Template) (*schema.Data, error) {
	data := schema.Data{
		ID:        templateRec.ID,
		CreatedAt: templateRec.CreatedAt,
		UpdatedAt: templateRec.UpdatedAt.Time,
	}
	return &data, nil
}
