package cmd

import (
	"io"
	"log"
	"os"

	"gopkg.in/natefinch/lumberjack.v2"
)


func SetupLogger() *log.Logger {
	logFile := &lumberjack.Logger{
		Filename:   "fetchmeta.log",
		MaxSize:    1, // In megabytes
		MaxBackups: 0,
		MaxAge:     365, // In days
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
