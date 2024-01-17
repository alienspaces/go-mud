package logger

type Level int

const (
	// DebugLevel -
	DebugLevel Level = 5
	// InfoLevel -
	InfoLevel Level = 4
	// WarnLevel -
	WarnLevel Level = 3
	// ErrorLevel -
	ErrorLevel Level = 2
)

const (
	ContextApplication string = "application"
	ContextPackage     string = "package"
	ContextFunction    string = "function"
	ContextCorrelation string = "correlation"
)

// Logger -
type Logger interface {
	NewInstance(fields map[string]interface{}) Logger
	WithContext(key, value string) Logger
	WithApplicationContext(value string) Logger
	WithPackageContext(value string) Logger
	WithFunctionContext(value string) Logger
	Write(level Level, msg string, args ...any)
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}
