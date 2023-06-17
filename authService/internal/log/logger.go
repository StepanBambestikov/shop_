package log

import (
	_log "log"
)

type ILogger interface {
	Errorf(format string, args ...interface{})
	Error(args ...interface{})
	Fatalf(format string, args ...interface{})
	Fatal(args ...interface{})
	Infof(format string, args ...interface{})
	Info(args ...interface{})
	Warnf(format string, args ...interface{})
	Warn(args ...interface{})
	Debugf(format string, args ...interface{})
	Debug(args ...interface{})
}

var logger ILogger

func InitLogger() {
	if logger == nil {
		logger_, err := NewZapLogger(true)
		if err != nil {
			_log.Fatal("could not init logger: %w", err)
			return
		}
		logger = logger_
	}
}

func Errorf(format string, args ...interface{}) {
	logger.Errorf(format, args...)
}

func Error(args ...interface{}) {
	logger.Error(args...)
}

func Fatalf(format string, args ...interface{}) {
	logger.Fatalf(format, args...)
}

func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}

func Infof(format string, args ...interface{}) {
	logger.Infof(format, args...)
}

func Info(args ...interface{}) {
	logger.Info(args...)
}

func Warnf(format string, args ...interface{}) {
	logger.Warnf(format, args...)
}

func Warn(args ...interface{}) {
	logger.Warn(args...)
}

func Debugf(format string, args ...interface{}) {
	logger.Debugf(format, args...)
}

func Debug(args ...interface{}) {
	logger.Debug(args...)
}
