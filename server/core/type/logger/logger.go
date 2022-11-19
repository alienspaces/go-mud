package logger

// Logger -
type Logger interface {
	NewInstance() (Logger, error)
	Context(key, value string)
	WithFunctionContext(value string) Logger
	WithPackageContext(value string) Logger
	WithInstanceContext(value string) Logger
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
}
