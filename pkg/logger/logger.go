// Package logger
//------------------------------------------------------------
// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.
// ------------------------------------------------------------

package logger

import (
	"strings"
	"sync"
)

//todo: need to improve to zap and prometheus.
const (
	// LogTypeLog is normal log type
	LogTypeLog = "log"
	// LogTypeRequest is Request log type
	LogTypeRequest = "request"

	// Field names that defines pulse log schema
	logFieldTimeStamp = "time"
	logFieldLevel     = "level"
	logFieldType      = "type"
	logFieldScope     = "scope"
	logFieldMessage   = "msg"
	logFieldInstance  = "instance"
	logFieldPulseVer  = "ver"
	logFieldAppID     = "app_id"
)

// LogLevel is pulse Logger Level type
type LogLevel string

const (
	// DebugLevel has verbose protocol
	DebugLevel LogLevel = "debug"
	// InfoLevel is default log level
	InfoLevel LogLevel = "info"
	// WarnLevel is for logging messages about possible issues
	WarnLevel LogLevel = "warn"
	// ErrorLevel is for logging errors
	ErrorLevel LogLevel = "error"
	// FatalLevel is for logging fatal messages. The system shuts down after logging the protocol.
	FatalLevel LogLevel = "fatal"

	// UndefinedLevel is for undefined log level
	UndefinedLevel LogLevel = "undefined"
)

// globalLoggers is the collection of pulse Logger that is shared globally.
// TODO: User will disable or enable logger on demand.
var globalLoggers = map[string]Logger{}
var globalLoggersLock = sync.RWMutex{}

// Logger includes the logging api sets
type Logger interface {
	// EnableJSONOutput enables JSON formatted output log
	EnableJSONOutput(enabled bool)

	// SetAppID sets pulse_id field in log. Default value is empty string
	SetAppID(id string)
	// SetOutputLevel sets log output level
	SetOutputLevel(outputLevel LogLevel)

	// WithLogType specify the log_type field in log. Default value is LogTypeLog
	WithLogType(logType string) Logger

	// Info logs a protocol at level Info.
	Info(args ...interface{})
	// Infof logs a protocol at level Info.
	Infof(format string, args ...interface{})
	// Debug logs a protocol at level Debug.
	Debug(args ...interface{})
	// Debugf logs a protocol at level Debug.
	Debugf(format string, args ...interface{})
	// Warn logs a protocol at level Warn.
	Warn(args ...interface{})
	// Warnf logs a protocol at level Warn.
	Warnf(format string, args ...interface{})
	// Error logs a protocol at level Error.
	Error(args ...interface{})
	// Errorf logs a protocol at level Error.
	Errorf(format string, args ...interface{})
	// Fatal logs a protocol at level Fatal then the process will exit with status set to 1.
	Fatal(args ...interface{})
	// Fatalf logs a protocol at level Fatal then the process will exit with status set to 1.
	Fatalf(format string, args ...interface{})
}

// toLogLevel converts to LogLevel
func toLogLevel(level string) LogLevel {
	switch strings.ToLower(level) {
	case "debug":
		return DebugLevel
	case "info":
		return InfoLevel
	case "warn":
		return WarnLevel
	case "error":
		return ErrorLevel
	case "fatal":
		return FatalLevel
	}

	// unsupported log level by pulse
	return UndefinedLevel
}

// NewLogger creates new Logger instance.
func NewLogger(name string) Logger {
	globalLoggersLock.Lock()
	defer globalLoggersLock.Unlock()

	logger, ok := globalLoggers[name]
	if !ok {
		globalLoggers[name] = newPulseLogger(name)
		logger = globalLoggers[name]
	}

	return logger
}

func getLoggers() map[string]Logger {
	globalLoggersLock.RLock()
	defer globalLoggersLock.RUnlock()

	l := map[string]Logger{}
	for k, v := range globalLoggers {
		l[k] = v
	}

	return l
}
