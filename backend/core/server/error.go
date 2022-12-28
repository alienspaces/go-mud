package server

import (
	"encoding/json"
	"encoding/xml"
	"net/http"

	coreerror "gitlab.com/alienspaces/go-mud/backend/core/error"
	"gitlab.com/alienspaces/go-mud/backend/core/type/logger"
)

func WriteError(l logger.Logger, w http.ResponseWriter, e error) {
	l = HTTPLogger(l, "WriteError")

	eres, err := coreerror.ToError(e)
	if err != nil {
		l.Error("System error >%v<", err)
		err = coreerror.GetRegistryError(coreerror.Internal)
		WriteSystemError(l, w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	status := eres.HttpStatusCode
	w.WriteHeader(eres.HttpStatusCode)

	l.Warn("Write error response status >%d<", status)

	if err := json.NewEncoder(w).Encode(e); err != nil {
		l.Error("failed writing response >%v<", err)
		err = coreerror.GetRegistryError(coreerror.Internal)
		WriteSystemError(l, w, err)
		return
	}
}

func WriteNotFoundError(l logger.Logger, w http.ResponseWriter, entity string, id string) {
	l = HTTPLogger(l, "WriteNotFoundError")

	e := coreerror.NewNotFoundError(entity, id)
	l.Warn("Resource not found >%v<", e)

	WriteError(l, w, e)
}

func WriteUnavailableError(l logger.Logger, w http.ResponseWriter, err error) {
	l = HTTPLogger(l, "WriteUnavailableError")

	e := coreerror.NewUnavailableError()
	l.Error("Service unavailable >%v< >%v<", err, e)

	WriteError(l, w, e)
}

func WriteSystemError(l logger.Logger, w http.ResponseWriter, err error) {
	l = HTTPLogger(l, "WriteSystemError")

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	e, err := coreerror.ToError(err)
	if err != nil {
		l.Error("failed writing response >%v<", err)
	}

	status := e.HttpStatusCode
	w.WriteHeader(e.HttpStatusCode)

	l.Warn("Write error response status >%d<", status)

	err = json.NewEncoder(w).Encode(e)
	if err != nil {
		l.Error("failed writing response >%v<", err)
	}
}

// WriteXMLErrorResponse responds with an 200 HTTP Status Code. For Service Cloud to retry message delivery, a nack (false) should be sent instead.
func WriteXMLErrorResponse(l logger.Logger, w http.ResponseWriter, s interface{}, err error) {
	l = HTTPLogger(l, "WriteXMLErrorResponse")

	w.Header().Set("Content-Type", "text/xml; charset=utf-8")

	e, err := coreerror.ToError(err)
	if err != nil {
		l.Error("failed writing response >%v<", err)
	}

	status := e.HttpStatusCode
	w.WriteHeader(e.HttpStatusCode)

	l.Info("Write error response status >%d<", status)

	if _, err := w.Write([]byte(xml.Header)); err != nil {
		l.Error("failed writing response >%v<", err)
		return
	}

	if err := xml.NewEncoder(w).Encode(s); err != nil {
		l.Error("failed encoding response >%v<", err)
	}
}
