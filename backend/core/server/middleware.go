package server

import (
	"github.com/julienschmidt/httprouter"
)

type MiddlewareFunc func(hc HandlerConfig, h Handle) (Handle, error)

// defaultMiddleware provides a list of default middleware
func (rnr *Runner) defaultMiddlewareFuncs() []MiddlewareFunc {
	return []MiddlewareFunc{
		rnr.WaitMiddleware,
		rnr.ParamMiddleware,
		rnr.DataMiddleware,
		rnr.RLSMiddleware,
		rnr.AuthzMiddleware,
		rnr.AuthenMiddleware,
		rnr.TxMiddleware,
		rnr.CorrelationMiddleware,
	}
}

// ApplyMiddleware applies middleware by the assigned middleware function
func (rnr *Runner) ApplyMiddleware(hc HandlerConfig, h Handle) (httprouter.Handle, error) {
	l := Logger(rnr.Log, "ApplyMiddleware")

	middlewareFuncs := rnr.HandlerMiddlewareFuncs()

	var err error
	for idx := range middlewareFuncs {
		h, err = middlewareFuncs[idx](hc, h)
		if err != nil {
			l.Warn("failed adding middleware >%v<", err)
			return nil, err
		}
	}

	return rnr.HttpRouterHandlerWrapper(h), nil
}
