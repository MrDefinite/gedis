package clientsdk

import (
	"github.com/sirupsen/logrus"
	"os"
)

const (
	// PanicLevel level, highest level of severity. Logs and then calls panic with the
	// message passed to Debug, Info, ...
	Panic = 0
	// FatalLevel level. Logs and then calls `os.Exit(1)`. It will exit even if the
	// logging level is set to Panic.
	Fatal = 1
	// ErrorLevel level. Logs. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	Error = 2
	// WarnLevel level. Non-critical entries that deserve eyes.
	Warn = 3
	// InfoLevel level. General operational entries about what's going on inside the
	// application.
	Info = 4
	// DebugLevel level. Usually only enabled when debugging. Very verbose logging.
	Debug = 5
)

func (gc *Gclient) SetLogLevel(level uint8) {
	switch level {
	case Panic:
		gc.logLevel = logrus.PanicLevel
		break
	case Fatal:
		gc.logLevel = logrus.FatalLevel
		break
	case Error:
		gc.logLevel = logrus.ErrorLevel
		break
	case Warn:
		gc.logLevel = logrus.WarnLevel
		break
	case Info:
		gc.logLevel = logrus.InfoLevel
		break
	case Debug:
		gc.logLevel = logrus.DebugLevel
		break
	default:
		gc.logLevel = defaultLogLevel
		break
	}

	log.Level = gc.logLevel
}

func (gc *Gclient) SetLogOutput(output string) {
	// If output is empty, output to console, otherwise to file
	if output == "" {
		log.Out = os.Stdout
		return
	}

	// TODO: to file
}
