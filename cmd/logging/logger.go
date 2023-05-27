package logging

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

func NewLogger() (*logrus.Logger, error) {
	processName, err := getProcessName()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return &logrus.Logger{}, err
	}

	logFname := fmt.Sprintf("%s.log", processName)

	logFile := &lumberjack.Logger{
		Filename:   logFname,
		MaxSize:    10, // In megabytes
		MaxBackups: 0,
		MaxAge:     365, // In days
		Compress:   true,
	}

	logger := logrus.New()
	logger.SetOutput(io.MultiWriter(os.Stdout, logFile))
	logger.SetFormatter(&logrus.JSONFormatter{})

	// Add this line for logging filename and line number!
	logger.SetReportCaller(true)

	return logger, nil
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

func getProcessName() (string, error) {
	fullexecpath, err := os.Executable()
	if err != nil {
		return "", err
	}

	_, execname := filepath.Split(fullexecpath)
	ext := filepath.Ext(execname)
	name := execname[:len(execname)-len(ext)]

	return name, nil
}

func init() {
	Logger, _ = NewLogger()
}

var Logger *logrus.Logger
