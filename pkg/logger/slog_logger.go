package logger

import (
	"fmt"
	"log/slog"
	"os"
	"time"
)

type SlogLogger struct {
	logger *slog.Logger
}

func NewSlogLogger() *SlogLogger {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	return &SlogLogger{
		logger: slog.New(handler),
	}
}

func (l *SlogLogger) Infof(format string, args ...any) {
	l.logger.Info(fmt.Sprintf(format, args...))
}

func (l *SlogLogger) Warnf(format string, args ...any) {
	l.logger.Warn(fmt.Sprintf(format, args...))
}

func (l *SlogLogger) Errorf(err error, format string, args ...any) {
	l.logger.Error(fmt.Sprintf(format, args...), slog.Any("err", err))
}

func (l *SlogLogger) Debugf(format string, args ...any) {
	l.logger.Debug(fmt.Sprintf(format, args...))
}

func (l *SlogLogger) Track(operationName string, operation func() error) error {
	start := time.Now()

	err := operation()

	l.Infof(
		"Operation finished: operation=%s, duration=%s, error=%v",
		operationName,
		time.Since(start),
		err,
	)

	return err
}
