package multi

import (
	"os"

	"github.com/sirupsen/logrus"
)

func NewLogger() *StandardLogger {
	baseLogger := logrus.New()
	baseLogger.SetOutput(os.Stdout)
	baseLogger.SetLevel(logrus.InfoLevel)
	var sl = &StandardLogger{baseLogger}
	sl.Formatter = &logrus.JSONFormatter{}
	return sl
}

var (
	internalError = LogEvent{InternalError, "Internal Error: %v"}
)

func (l *StandardLogger) LogInternalError(message string) {
	l.Errorf(internalError.message, message)
}
