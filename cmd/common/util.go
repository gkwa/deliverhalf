package cmd

import (
	"io"
	"log"
	"os"

	"gopkg.in/natefinch/lumberjack.v2"
)

func SetupLogger() *log.Logger {
	logFile := &lumberjack.Logger{
		Filename:   "deliverhalf.log",
		MaxSize:    10, // In megabytes
		MaxBackups: 0,
		MaxAge:     365, // In days
		Compress:   true,
	}
	defer logFile.Close()
	logWriter := io.MultiWriter(logFile, os.Stderr)
	return log.New(logWriter, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.LUTC|log.Lshortfile)
}

func FileExists(logger *log.Logger, filePath string) bool {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}
