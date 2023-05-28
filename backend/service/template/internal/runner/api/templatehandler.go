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
	"gitlab.com/alienspaces/go-mud/backend/schema/template"
	"gitlab.com/alienspaces/go-mud/backend/service/template/internal/model"
	"gitlab.com/alienspaces/go-mud/backend/service/template/internal/record"
)

const (
	HandlerGetManyTemplates        = "get-many-templates"
	HandlerGetOneTemplate          = "get-one-template"
	HandlerCreateOneTemplate       = "create-one-template"
	HandlerCreateOneTemplateWithID = "create-one-template-with-id"
	HandlerUpdateOneTemplate       = "update-one-template"
)

func (rnr *Runner) templateHandlerConfig(l logger.Logger, config map[string]server.HandlerConfig) (map[string]server.HandlerConfig, error) {

	config[HandlerGetOneTemplate] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/templates",
		HandlerFunc: rnr.getTemplatesHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypePublic,
			},
			ValidateResponseSchema: &jsonschema.SchemaWithReferences{
				Main: jsonschema.Schema{
					Name: "main.response.schema.json",
				},
				References: []jsonschema.Schema{
					{
						Name: "data.response.schema.json",
					},
				},
			},
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:    true,
			Description: "Query templates.",
		},
	}

	config[HandlerGetOneTemplate] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/templates/:template_id",
		HandlerFunc: rnr.getTemplateHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypePublic,
			},
			ValidateResponseSchema: &jsonschema.SchemaWithReferences{
				Main: jsonschema.Schema{
					Name: "main.response.schema.json",
				},
				References: []jsonschema.Schema{
					{
						Name: "data.response.schema.json",
					},
				},
			},
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:    true,
			Description: "Get a template.",
		},
	}

	config[HandlerCreateOneTemplate] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/api/templates",
		HandlerFunc: rnr.PostTemplatesHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypePublic,
			},
			ValidateRequestSchema: &jsonschema.SchemaWithReferences{
				Main: jsonschema.Schema{
					Name: "main.request.schema.json",
				},
				References: []jsonschema.Schema{
					{
						Name: "data.request.schema.json",
					},
				},
			},
			ValidateResponseSchema: &jsonschema.SchemaWithReferences{
				Main: jsonschema.Schema{
					Name: "main.response.schema.json",
				},
				References: []jsonschema.Schema{
					{
						Name: "data.response.schema.json",
					},
				},
			},
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:    true,
			Description: "Create a template.",
		},
	}

	config[HandlerCreateOneTemplateWithID] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/api/templates/:template_id",
		HandlerFunc: rnr.PostTemplatesHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypePublic,
			},
			ValidateRequestSchema: &jsonschema.SchemaWithReferences{
				Main: jsonschema.Schema{
					Name: "main.request.schema.json",
				},
				References: []jsonschema.Schema{
					{
						Name: "data.request.schema.json",
					},
				},
			},
			ValidateResponseSchema: &jsonschema.SchemaWithReferences{
				Main: jsonschema.Schema{
					Name: "main.response.schema.json",
				},
				References: []jsonschema.Schema{
					{
						Name: "data.response.schema.json",
					},
				},
			},
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:    true,
			Description: "Create a template.",
		},
	}

	config[HandlerUpdateOneTemplate] = server.HandlerConfig{
		Method:      http.MethodPut,
		Path:        "/api/templates/:template_id",
		HandlerFunc: rnr.PutTemplatesHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypePublic,
			},
			ValidateRequestSchema: &jsonschema.SchemaWithReferences{
				Main: jsonschema.Schema{
					Name: "main.request.schema.json",
				},
				References: []jsonschema.Schema{
					{
						Name: "data.request.schema.json",
					},
				},
			},
			ValidateResponseSchema: &jsonschema.SchemaWithReferences{
				Main: jsonschema.Schema{
					Name: "main.response.schema.json",
				},
				References: []jsonschema.Schema{
					{
						Name: "data.response.schema.json",
					},
				},
			},
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:    true,
			Description: "Update a template.",
		},
	}

	return config, nil
}

// getTemplateHandler -
func (rnr *Runner) getTemplateHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m modeller.Modeller) error {

	l.Info("** get template handler ** p >%#v< m >%#v<", pp, m)

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

	// response data
	responseData, err := rnr.RecordToTemplateResponseData(rec)
	if err != nil {
		l.Warn("failed mapping template response >%v<", err)
		server.WriteError(l, w, err)
		return err
	}

	res := template.Response{
		Data: responseData,
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
	l.Info("Get templates handler pp >%#v< qp >%#v<", pp, qp)

	var err error

	l.Info("Querying template records")

	// query parameters
	opts := queryparam.ToSQLOptions(qp)
	if err != nil {
		l.Warn("failed to map collection query params to repository params >%v<", err)
		server.WriteError(l, w, err)
		return err
	}

	recs, err := m.(*model.Model).GetTemplateRecs(opts)
	if err != nil {
		server.WriteError(l, w, err)
		return err
	}

	// assign response properties
	data := []*template.Data{}
	for _, rec := range recs {

		// response data
		responseData, err := rnr.RecordToTemplateResponseData(rec)
		if err != nil {
			server.WriteSystemError(l, w, err)
			return err
		}

		data = append(data, responseData)
	}

	res := template.CollectionResponse{
		Data: data,
	}

	err = server.WriteResponse(l, w, http.StatusOK, res)
	if err != nil {
		l.Warn("Failed writing response >%v<", err)
		return err
	}

	return nil
}

// PostTemplatesHandler -
func (rnr *Runner) PostTemplatesHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m modeller.Modeller) error {
	l.Info("Post templates handler pp >%#v< qp >%#v<", pp, qp)

	// parameters
	id := pp.ByName("template_id")

	req, err := server.ReadRequest(l, r, &template.Request{})
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
	res := template.Response{
		Data: responseData,
	}

	err = server.WriteResponse(l, w, http.StatusCreated, res)
	if err != nil {
		l.Warn("Failed writing response >%v<", err)
		return err
	}

	return nil
}

// PutTemplatesHandler -
func (rnr *Runner) PutTemplatesHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m modeller.Modeller) error {
	l.Info("Put templates handler pp >%#v< qp >%#v<", pp, qp)

	// parameters
	id := pp.ByName("template_id")

	l.Info("Updating resource ID >%s<", id)

	_, err := m.(*model.Model).GetTemplateRec(id, nil)
	if err != nil {
		server.WriteError(l, w, err)
		return err
	}

	req, err := server.ReadRequest(l, r, &template.Request{})
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
	res := template.Response{
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
func (rnr *Runner) TemplateRequestDataToRecord(data *template.Data) (*record.Template, error) {
	return &record.Template{}, nil
}

// RecordToTemplateResponseData -
func (rnr *Runner) RecordToTemplateResponseData(templateRec *record.Template) (*template.Data, error) {
	data := template.Data{
		ID:        templateRec.ID,
		CreatedAt: templateRec.CreatedAt,
		UpdatedAt: templateRec.UpdatedAt.Time,
	}
	return &data, nil
}
