package server

import (
	"encoding/json"
	"encoding/xml"
	"net/http"

	coreerror "gitlab.com/alienspaces/go-mud/backend/core/error"
	"gitlab.com/alienspaces/go-mud/backend/core/type/logger"
)

func WriteError(l logger.Logger, w http.ResponseWriter, errs ...error) {
	if len(errs) == 0 {
		// This is a logic error!
		l.Error("no errors passed to WriteError")
		writeSystemError(l, w)
		return
	}

	results, err := coreerror.ToErrors(errs...)
	if err != nil {
		l.Error("system error >%#v< >%v<", errs, err)
		writeSystemError(l, w)
		return
	}

	if len(results) > 1 {
		for _, e := range results {
			if e.HttpStatusCode != http.StatusBadRequest {
				l.Error("the slice of errors must contain 400 errors only >%#v<", e)
				writeSystemError(l, w)
				return
			}
		}
	}

	status := results[0].HttpStatusCode
	if status == http.StatusInternalServerError {
		l.Error(results[0].Error())
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	w.WriteHeader(status)
	l.Info("write response status >%d<", status)

	if err := json.NewEncoder(w).Encode(results); err != nil {
		l.Error("failed writing response >%v<", err)
		writeSystemError(l, w)
		return
	}
}

func WriteMalformedError(l logger.Logger, w http.ResponseWriter, err error) {
	e := coreerror.NewMalformedError()
	l.Warn("malformed data >%v< >%v", err, e)
	WriteError(l, w, e)
}

func WriteNotFoundError(l logger.Logger, w http.ResponseWriter, entity string, id string) {
	e := coreerror.NewNotFoundError(entity, id)
	l.Warn("not found error >%v<", e)

	WriteError(l, w, e)
}

func WriteUnauthorizedError(l logger.Logger, w http.ResponseWriter, err error) {
	e := coreerror.NewUnauthorizedError()
	l.Error("unauthorized error >%v< >%v<", err, e)

	WriteError(l, w, e)
}

func WriteUnavailableError(l logger.Logger, w http.ResponseWriter, err error) {
	e := coreerror.NewUnavailableError()
	l.Error("unavailable error >%v< >%v<", err, e)

	WriteError(l, w, e)
}

func WriteSystemError(l logger.Logger, w http.ResponseWriter, err error) {
	l.Error("system error >%v<", err)

	writeSystemError(l, w)
}

func writeSystemError(l logger.Logger, w http.ResponseWriter) {
	err := coreerror.NewInternalError("system")
	e, _ := coreerror.ToError(err)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	status := e.HttpStatusCode
	w.WriteHeader(e.HttpStatusCode)

	l.Info("write response status >%d<", status)

	if err := json.NewEncoder(w).Encode([]coreerror.Error{e}); err != nil {
		l.Error("failed writing response >%v<", e)
	}
}

// WriteXMLErrorResponse responds with an 200 HTTP Status Code. For Service Cloud to retry message delivery, a nack (false) should be sent instead.
func WriteXMLErrorResponse(l logger.Logger, w http.ResponseWriter, s interface{}, e error) {
	l.Debug("writing error response >%+v<", s)

	if e != nil && !coreerror.IsError(e) {
		l.Error("system error >%v<", e)
	}

	w.Header().Set("Content-Type", "text/xml; charset=utf-8")

	w.WriteHeader(http.StatusOK)

	if _, err := w.Write([]byte(xml.Header)); err != nil {
		l.Error("failed writing response >%v<", err)
		return
	}

	if err := xml.NewEncoder(w).Encode(s); err != nil {
		l.Error("failed encoding response >%v<", err)
	}
}
