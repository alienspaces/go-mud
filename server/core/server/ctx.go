package server

type contextKey int

const (
	ctxKeyAuth          contextKey = 1
	ctxKeyData          contextKey = 2
	ctxKeyCorrelationID contextKey = 3
)
