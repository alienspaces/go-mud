package log

import (
	"os"
	"strings"

	"github.com/rs/zerolog"

	"gitlab.com/alienspaces/go-boilerplate/server/core/type/configurer"
	"gitlab.com/alienspaces/go-boilerplate/server/core/type/logger"
)

// Log -
type Log struct {
	log    zerolog.Logger
	fields map[string]interface{}
	Config configurer.Configurer
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

	// TODO: support different output writers primarily for testing purposes

	// logger
	l.log = zerolog.New(os.Stdout).With().Timestamp().Logger()

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

	l.log.Info().Msgf("Log level config >%s< actual >%s<", configLevel, l.log.GetLevel())

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
