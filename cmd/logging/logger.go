package logging

import (
	"io"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

func NewLogger() (*logrus.Logger, error) {
	logFile := &lumberjack.Logger{
		Filename:   "deliverhalf.log",
		MaxSize:    10, // In megabytes
		MaxBackups: 0,
		MaxAge:     365, // In days
		Compress:   true,
	}

	logger := logrus.New()
	logger.SetOutput(io.MultiWriter(os.Stdout, logFile))
	logger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: time.RFC3339,
		DisableColors:   true,
		FullTimestamp:   true,
	})

	// Add this line forlog.filename and line number!
	logger.SetReportCaller(true)

	logger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		DisableColors:   true,
		FullTimestamp:   true,
	})

	logger.SetReportCaller(true)

	return logger, nil
}

var Logger *logrus.Logger

func init() {
	Logger, _ = NewLogger()
}

func ParseLogLevel(level string) logrus.Level {
	switch level {
	case "trace":
		return logrus.TraceLevel
	case "debug":
		return logrus.DebugLevel
	case "info":
		return logrus.InfoLevel
	case "warn", "warning":
		return logrus.WarnLevel
	case "error":
		return logrus.ErrorLevel
	case "fatal":
		return logrus.FatalLevel
	case "panic":
		return logrus.PanicLevel
	default:
		return logrus.TraceLevel
	}
}
