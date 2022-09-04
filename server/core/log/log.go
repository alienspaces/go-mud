package log

import (
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog"

	"gitlab.com/alienspaces/go-mud/server/core/type/configurer"
	"gitlab.com/alienspaces/go-mud/server/core/type/logger"
)

// Log -
type Log struct {
	log      zerolog.Logger
	fields   map[string]interface{}
	Config   configurer.Configurer
	IsPretty bool
}

var _ logger.Logger = &Log{}

// Level -
type Level uint32

const (
	// DebugLevel -
	DebugLevel = 5
	// InfoLevel -
	InfoLevel = 4
	// WarnLevel -
	WarnLevel = 3
	// ErrorLevel -
	ErrorLevel = 2
)

var levelMap = map[Level]zerolog.Level{
	// DebugLevel -
	DebugLevel: zerolog.DebugLevel,
	// InfoLevel -
	InfoLevel: zerolog.InfoLevel,
	// WarnLevel -
	WarnLevel: zerolog.WarnLevel,
	// ErrorLevel -
	ErrorLevel: zerolog.ErrorLevel,
}

// NewLogger returns a logger
func NewLogger(c configurer.Configurer) (*Log, error) {

	l := Log{
		fields: make(map[string]interface{}),
		Config: c,
	}

	err := l.Init()
	if err != nil {
		return nil, err
	}

	return &l, nil
}

// Init initializes logger
func (l *Log) Init() error {

	l.log = zerolog.New(os.Stdout).With().Timestamp().Logger()

	if l.Config.Get("APP_SERVER_LOG_PRETTY") == "true" {
		l.IsPretty = true
	}

	if l.IsPretty {
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

	// logger level
	configLevel := l.Config.Get("APP_SERVER_LOG_LEVEL")

	level := strings.ToLower(configLevel)

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

	return nil
}

// NewInstance - Create a new log instance based off configuration of this instance
func (l *Log) NewInstance() (logger.Logger, error) {

	i := Log{
		fields: make(map[string]interface{}),
		Config: l.Config,
	}

	err := i.Init()
	if err != nil {
		return nil, err
	}

	return &i, nil
}

// Printf -
func (l *Log) Printf(format string, args ...interface{}) {
	l.log.Printf(format, args...)
}

// Level -
func (l *Log) Level(level Level) {
	if lvl, ok := levelMap[level]; ok {
		l.log = l.log.Level(lvl)
	}
}

// Context - set logging
func (l *Log) Context(key, value string) {
	if value == "" {
		delete(l.fields, key)
		return
	}
	l.fields[key] = value
}

// WithFunctionContext - Shallow copied logger instance with new function context and existing package context
func (l *Log) WithFunctionContext(value string) logger.Logger {
	ctxLog := *l
	fields := map[string]interface{}{
		"function": value,
	}
	if value, ok := ctxLog.fields["instance"]; ok {
		fields["instance"] = value
	}
	if value, ok := ctxLog.fields["package"]; ok {
		fields["package"] = value
	}
	if value, ok := ctxLog.fields["correlation-id"]; ok {
		fields["correlation-id"] = value
	}

	ctxLog.fields = fields
	return &ctxLog
}

// WithInstanceContext - Shallow copied logger instance with new function context and existing package context
func (l *Log) WithInstanceContext(value string) logger.Logger {
	ctxLog := *l
	fields := map[string]interface{}{
		"instance": value,
	}
	if value, ok := ctxLog.fields["function"]; ok {
		fields["function"] = value
	}
	if value, ok := ctxLog.fields["package"]; ok {
		fields["package"] = value
	}
	if value, ok := ctxLog.fields["correlation-id"]; ok {
		fields["correlation-id"] = value
	}

	ctxLog.fields = fields
	return &ctxLog
}

// WithPackageContext - Shallow copied logger instance with new package context and existing function context
func (l *Log) WithPackageContext(value string) logger.Logger {
	ctxLog := *l
	fields := map[string]interface{}{
		"package": value,
	}
	if value, ok := ctxLog.fields["instance"]; ok {
		fields["instance"] = value
	}
	if value, ok := ctxLog.fields["function"]; ok {
		fields["function"] = value
	}
	if value, ok := ctxLog.fields["correlation-id"]; ok {
		fields["correlation-id"] = value
	}
	ctxLog.fields = fields
	return &ctxLog
}

// Debug -
func (l *Log) Debug(msg string, args ...interface{}) {
	ctxLog := l.log.With().Fields(l.fields).Logger()
	ctxLog.Debug().Msgf(msg, args...)
}

// Info -
func (l *Log) Info(msg string, args ...interface{}) {
	ctxLog := l.log.With().Fields(l.fields).Logger()
	ctxLog.Info().Msgf(msg, args...)
}

// Warn -
func (l *Log) Warn(msg string, args ...interface{}) {
	ctxLog := l.log.With().Fields(l.fields).Logger()
	ctxLog.Warn().Msgf(msg, args...)
}

// Error -
func (l *Log) Error(msg string, args ...interface{}) {
	ctxLog := l.log.With().Fields(l.fields).Logger()
	ctxLog.Error().Msgf(msg, args...)
}
