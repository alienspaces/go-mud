package runner

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	coreerror "gitlab.com/alienspaces/go-mud/server/core/error"
	"gitlab.com/alienspaces/go-mud/server/core/server"
	"gitlab.com/alienspaces/go-mud/server/core/type/logger"
	"gitlab.com/alienspaces/go-mud/server/core/type/modeller"
	"gitlab.com/alienspaces/go-mud/server/schema"
	"gitlab.com/alienspaces/go-mud/server/service/template/internal/model"
	"gitlab.com/alienspaces/go-mud/server/service/template/internal/record"
)

// GetTemplatesHandler -
func (rnr *Runner) GetTemplatesHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp map[string]interface{}, l logger.Logger, m modeller.Modeller) error {

	l.Info("** Get templates handler ** p >%#v< m >%#v<", pp, m)

	var recs []*record.Template
	var rec *record.Template
	var err error

	// Path parameters
	id := pp.ByName("template_id")

	// Single resource
	if id != "" {

		l.Info("Getting template record ID >%s<", id)

		rec, err = m.(*model.Model).GetTemplateRec(id, false)
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

		recs = append(recs, rec)

	} else {

		l.Info("Querying template records")

		params := make(map[string]interface{})
		for paramName, paramValue := range qp {
			l.Info("Querying template records with param name >%s< value >%v<", paramName, paramValue)
			params[paramName] = paramValue
		}

		recs, err = m.(*model.Model).GetTemplateRecs(params, nil, false)
		if err != nil {
			l.Warn("failed getting template record >%v<", err)
			server.WriteError(l, w, err)
			return err
		}
	}

	// Assign response properties
	data := []schema.TemplateData{}
	for _, rec := range recs {

		// Response data
		responseData, err := rnr.RecordToTemplateResponseData(rec)
		if err != nil {
			l.Warn("failed getting template record >%v<", err)
			server.WriteError(l, w, err)
			return err
		}

		data = append(data, responseData)
	}

	res := schema.TemplateResponse{
		Data: data,
	}

	err = server.WriteResponse(l, w, res)
	if err != nil {
		l.Warn("Failed writing response >%v<", err)
		return err
	}

	return nil
}

// PostTemplatesHandler -
func (rnr *Runner) PostTemplatesHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp map[string]interface{}, l logger.Logger, m modeller.Modeller) error {

	l.Info("** Post templates handler ** p >%#v< m >#%v<", pp, m)

	// Path parameters
	id := pp.ByName("template_id")

	req := schema.TemplateRequest{}

	err := server.ReadRequest(l, r, &req)
	if err != nil {
		server.WriteError(l, w, err)
		return err
	}

	rec := record.Template{}

	// Assign request properties
	rec.ID = id

	// Record data
	err = rnr.TemplateRequestDataToRecord(req.Data, &rec)
	if err != nil {
		server.WriteSystemError(l, w, err)
		return err
	}

	err = m.(*model.Model).CreateTemplateRec(&rec)
	if err != nil {
		server.WriteSystemError(l, w, err)
		return err
	}

	// Response data
	responseData, err := rnr.RecordToTemplateResponseData(&rec)
	if err != nil {
		server.WriteSystemError(l, w, err)
		return err
	}

	// Assign response properties
	res := schema.TemplateResponse{
		Data: []schema.TemplateData{
			responseData,
		},
	}

	err = server.WriteResponse(l, w, res)
	if err != nil {
		l.Warn("Failed writing response >%v<", err)
		return err
	}

	return nil
}

// PutTemplatesHandler -
func (rnr *Runner) PutTemplatesHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp map[string]interface{}, l logger.Logger, m modeller.Modeller) error {

	l.Info("** Put templates handler ** p >%#v< m >#%v<", pp, m)

	// Path parameters
	id := pp.ByName("template_id")

	l.Info("Updating resource ID >%s<", id)

	rec, err := m.(*model.Model).GetTemplateRec(id, false)
	if err != nil {
		server.WriteSystemError(l, w, err)
		return err
	}

	// Resource not found
	if rec == nil {
		err := coreerror.NewNotFoundError("template", id)
		server.WriteError(l, w, err)
		return err
	}

	req := schema.TemplateRequest{}

	err = server.ReadRequest(l, r, &req)
	if err != nil {
		server.WriteSystemError(l, w, err)
		return err
	}

	// Record data
	err = rnr.TemplateRequestDataToRecord(req.Data, rec)
	if err != nil {
		server.WriteSystemError(l, w, err)
		return err
	}

	err = m.(*model.Model).UpdateTemplateRec(rec)
	if err != nil {
		server.WriteSystemError(l, w, err)
		return err
	}

	// Response data
	responseData, err := rnr.RecordToTemplateResponseData(rec)
	if err != nil {
		server.WriteSystemError(l, w, err)
		return err
	}

	// Assign response properties
	res := schema.TemplateResponse{
		Data: []schema.TemplateData{
			responseData,
		},
	}

	err = server.WriteResponse(l, w, res)
	if err != nil {
		l.Warn("Failed writing response >%v<", err)
		return err
	}

	return nil
}

// TemplateRequestDataToRecord -
func (rnr *Runner) TemplateRequestDataToRecord(data schema.TemplateData, rec *record.Template) error {

	return nil
}

// RecordToTemplateResponseData -
func (rnr *Runner) RecordToTemplateResponseData(templateRec *record.Template) (schema.TemplateData, error) {

	data := schema.TemplateData{
		ID:        templateRec.ID,
		CreatedAt: templateRec.CreatedAt,
		UpdatedAt: templateRec.UpdatedAt.Time,
	}

	return data, nil
}
