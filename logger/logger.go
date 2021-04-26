package logger

import (
	"fmt"
	"runtime"

	"github.com/sirupsen/logrus"
)

type Event struct {
	id      int
	message string
}

type Logger struct {
	*logrus.Logger
}

func NewLogger() *Logger {
	var baseLogger = logrus.New()
	var standardLogger = &Logger{baseLogger}

	standardLogger.Formatter = &logrus.JSONFormatter{DisableHTMLEscape: true}

	return standardLogger
}

var (
	genericErrorMessage    = Event{-1, "Error: %s"}
	logOnEntryMessage      = Event{0, "Entry args: %s"}
	logOnExitMessage       = Event{1, "Exit args: %s"}
	invalidArgMessage      = Event{3, "Invalid arg: %s"}
	invalidArgValueMessage = Event{4, "Invalid value for argument: %s: %v"}
	missingArgMessage      = Event{5, "Missing arg: %s"}
	internalErrorMessage   = Event{500, "Internal error: %s"}
	badRequestErrorMessage = Event{400, "Bad request error: %s"}
)

func (l *Logger) SetLoggingLevel(level string) {
	switch level {
	case "info":
		l.SetLevel(logrus.InfoLevel)
	case "warn":
		l.SetLevel(logrus.WarnLevel)
	case "debug":
		l.SetLevel(logrus.DebugLevel)
	case "trace":
		l.SetLevel(logrus.TraceLevel)
	case "fatal":
		l.SetLevel(logrus.FatalLevel)
	case "panic":
		l.SetLevel(logrus.PanicLevel)
	default:
		l.SetLevel(logrus.ErrorLevel)
	}
}

func (l *Logger) GetContext() string {
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	return fmt.Sprintf("%s - %s:%d", frame.Function, frame.File, frame.Line)
}

func (l *Logger) LogOnEntryWithContext(context string, arguments ...interface{}) {
	l.WithField("function", "> "+context).Debugf(logOnEntryMessage.message, arguments)
}

func (l *Logger) LogOnExitWithContext(context string, arguments ...interface{}) {
	l.WithField("function", "< "+context).Debugf(logOnExitMessage.message, arguments)
}

func (l *Logger) LogOnInternalErrorWithContext(context string, arguments ...interface{}) {
	l.WithField("function", "< "+context).Errorf(internalErrorMessage.message, arguments)
}

func (l *Logger) LogOnBadRequestErrorWithContext(context string, arguments ...interface{}) {
	l.WithField("function", "< "+context).Errorf(badRequestErrorMessage.message, arguments)
}

func (l *Logger) LogOnErrorWithContext(context string, arguments ...interface{}) {
	l.WithField("function", "< "+context).Errorf(genericErrorMessage.message, arguments)
}

// func (l *Logger) InvalidArg(argumentName string) {
// 	l.Errorf(invalidArgMessage.message, argumentName)
// }

// func (l *Logger) InvalidArgValue(argumentName string, argumentValue string) {
// 	l.Errorf(invalidArgValueMessage.message, argumentName, argumentValue)
// }

// func (l *Logger) MissingArg(argumentName string) {
// 	l.Errorf(missingArgMessage.message, argumentName)
// }
