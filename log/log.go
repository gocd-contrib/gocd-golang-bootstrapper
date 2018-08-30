package log

import (
	"os"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("gocd-golang-bootstrapper")
var format = logging.MustStringFormatter(
	`%{color}%{time:2006-01-02T15:04:05.999Z-07:00} [Bootstrap] [%{level:.4s}]%{color:reset} %{message}`,
)

// Init initializes a logger instance with STDERR backend
func Init() {
	backend := logging.NewLogBackend(os.Stderr, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)
	logging.SetBackend(backendFormatter)
}

// Critical logs a message using CRITICAL as log level.
func Critical(msg string) {
	log.Critical(msg)
}

// Criticalf logs a message using CRITICAL as log level.
func Criticalf(format string, args ...interface{}) {
	log.Criticalf(format, args...)
}

// Debugf logs a message using DEBUG as log level.
func Debugf(format string, args ...interface{}) {
	log.Debugf(format, args...)
}

// Info logs a message using INFO as log level.
func Info(msg string) {
	log.Info(msg)
}

// Infof logs a message using DEBUG as log level.
func Infof(format string, args ...interface{}) {
	log.Infof(format, args...)
}

// Warningf logs a message using EARNING as log level.
func Warningf(format string, args ...interface{}) {
	log.Warningf(format, args...)
}
