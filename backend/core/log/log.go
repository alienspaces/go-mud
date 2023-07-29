package log

import (
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog"

	"gitlab.com/alienspaces/go-mud/backend/core/type/configurer"
	"gitlab.com/alienspaces/go-mud/backend/core/type/logger"
)

// Environment configuration keys
const (
	EnvKeyAppServerLogLevel  = "APP_SERVER_LOG_LEVEL"
	EnvKeyAppServerLogPretty = "APP_SERVER_LOG_PRETTY"
)

// Log -
type Log struct {
	log    zerolog.Logger
	fields map[string]interface{}
	Config Config
}

type Config struct {
	Level  string
	Pretty bool
}

var _ logger.Logger = &Log{}

var levelMap = map[logger.Level]zerolog.Level{
	// DebugLevel -
	logger.DebugLevel: zerolog.DebugLevel,
	// InfoLevel -
	logger.InfoLevel: zerolog.InfoLevel,
	// WarnLevel -
	logger.WarnLevel: zerolog.WarnLevel,
	// ErrorLevel -
	logger.ErrorLevel: zerolog.ErrorLevel,
}

func NewDefaultLogger() *Log {
	l := Log{
		fields: make(map[string]interface{}),
		Config: Config{
			Level:  "info",
			Pretty: false,
		},
	}

	l.Init()
	return &l
}

// NewLogger returns a logger
func NewLogger(c configurer.Configurer) (*Log, error) {

	l := Log{
		fields: make(map[string]interface{}),
		Config: Config{
			Level:  c.Get(EnvKeyAppServerLogLevel),
			Pretty: c.Get(EnvKeyAppServerLogPretty) == "true",
		},
	}

	l.Init()
	return &l, nil
}

// NewLoggerWithConfig returns a logger with the provided configuration
func NewLoggerWithConfig(c Config) (*Log, error) {

	l := Log{
		fields: make(map[string]interface{}),
		Config: c,
	}

	l.Init()
	return &l, nil
}

// Init initializes logger
func (l *Log) Init() {

	l.log = zerolog.New(os.Stdout).With().Timestamp().Logger()

	// Pretty
	if l.Config.Pretty {
		output := zerolog.ConsoleWriter{
			Out: os.Stdout,
			// The following adds colour to the value of additional log fields, a nice shade of purple actually..
			FormatFieldValue: func(i interface{}) string {
				if i != nil {
					return fmt.Sprintf("\x1b[%dm%v\x1b[0m", 35, i)
				}
				return ""
			},
		}
		l.log = l.log.Output(output)
	}

	// Level
	level := strings.ToLower(l.Config.Level)

	switch level {
	case "debug":
		l.log = l.log.Level(zerolog.DebugLevel)
	case "info":
		l.log = l.log.Level(zerolog.InfoLevel)
	case "warn":
		l.log = l.log.Level(zerolog.WarnLevel)
	case "error":
		l.log = l.log.Level(zerolog.ErrorLevel)
	default:
		l.log = l.log.Level(zerolog.DebugLevel)
	}
}

// NewInstance - Create a new log instance based off configuration of this instance
func (l *Log) NewInstance(fields map[string]interface{}) logger.Logger {
	return &Log{
		fields: fields,
		Config: l.Config,
		log:    l.log.With().Fields(fields).Logger(),
	}
}

// Printf -
func (l *Log) Printf(format string, args ...interface{}) {
	l.log.Printf(format, args...)
}

// Level -
func (l *Log) Level(level logger.Level) {
	if lvl, ok := levelMap[level]; ok {
		l.log = l.log.Level(lvl)
	}
}

func (l *Log) WithContext(key, value string) logger.Logger {
	fields := map[string]interface{}{
		key: value,
	}
	for field, value := range l.fields {
		if field == key {
			continue
		}
		fields[field] = value
	}
	i := l.NewInstance(fields)
	return i
}

// WithApplicationContext - New logger instance with new application context and existing package and function context
func (l *Log) WithApplicationContext(value string) logger.Logger {
	return l.WithContext(logger.ContextApplication, value)
}

// WithPackageContext - New logger instance with new package context and existing application and function context
func (l *Log) WithPackageContext(value string) logger.Logger {
	return l.WithContext(logger.ContextPackage, value)
}

// WithFunctionContext - New logger instance with new function context and existing application and package context
func (l *Log) WithFunctionContext(value string) logger.Logger {
	return l.WithContext(logger.ContextFunction, value)
}

func (l *Log) Write(lvl logger.Level, msg string, args ...any) {
	switch lvl {
	case logger.DebugLevel:
		l.Debug(msg, args...)
	case logger.InfoLevel:
		l.Info(msg, args...)
	case logger.WarnLevel:
		l.Warn(msg, args...)
	case logger.ErrorLevel:
		l.Error(msg, args...)
	}
}

// Debug -
func (l *Log) Debug(msg string, args ...any) {
	l.log.Debug().Msgf(msg, args...)
}

// Info -
func (l *Log) Info(msg string, args ...any) {
	l.log.Info().Msgf(msg, args...)
}

// Warn -
func (l *Log) Warn(msg string, args ...any) {
	l.log.Warn().Msgf(msg, args...)
}

// Error -
func (l *Log) Error(msg string, args ...any) {
	l.log.Error().Msgf(msg, args...)
}
