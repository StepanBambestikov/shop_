package log

import (
	"errors"
	"go.uber.org/zap"
	_log "log"
	"syscall"
)

type zapLogger struct {
	logger *zap.SugaredLogger
}

func (z *zapLogger) syncLogger() {
	// https://github.com/uber-go/zap/issues/991#issuecomment-962098428
	err := z.logger.Sync()
	if err != nil && !errors.Is(err, syscall.ENOTTY) {
		_log.Println("Sync error: ", err)
	}
}

func (z *zapLogger) Errorf(format string, args ...interface{}) {
	z.logger.Errorf(format, args...)
	z.syncLogger()
}

func (z *zapLogger) Error(args ...interface{}) {
	z.logger.Error(args...)
	z.syncLogger()
}

func (z *zapLogger) Fatalf(format string, args ...interface{}) {
	z.logger.Fatalf(format, args...)
	z.syncLogger()
}

func (z *zapLogger) Fatal(args ...interface{}) {
	z.logger.Fatal(args...)
	z.syncLogger()
}

func (z *zapLogger) Infof(format string, args ...interface{}) {
	z.logger.Infof(format, args...)
	z.syncLogger()
}

func (z *zapLogger) Info(args ...interface{}) {
	z.logger.Info(args...)
	z.syncLogger()
}

func (z *zapLogger) Warnf(format string, args ...interface{}) {
	z.logger.Warnf(format, args...)
	z.syncLogger()
}

func (z *zapLogger) Warn(args ...interface{}) {
	z.logger.Warn(args...)
	z.syncLogger()
}

func (z *zapLogger) Debugf(format string, args ...interface{}) {
	z.logger.Debugf(format, args...)
	z.syncLogger()
}

func (z *zapLogger) Debug(args ...interface{}) {
	z.logger.Debug(args...)
	z.syncLogger()
}

func NewZapLogger(debug bool) (ILogger, error) {
	var logger *zap.Logger
	var err error

	log := &zapLogger{}

	// AddCallerSkip is needed to print real caller, not this file (zap.go)
	if debug {
		logger, err = zap.NewDevelopment(zap.AddCallerSkip(2))
	} else {
		logger, err = zap.NewProduction(zap.AddCallerSkip(2))
	}
	if err != nil {
		return nil, err
	}

	log.logger = logger.Sugar()

	return log, nil
}
