package logger

type Logger interface {
	Infof(format string, args ...any)
	Warnf(format string, args ...any)
	Errorf(err error, format string, args ...any)
	Debugf(format string, args ...any)
	Track(operationName string, operation func() error) error
}
