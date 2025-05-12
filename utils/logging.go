package utils

import (
	"log/slog"
	"os"
	"path/filepath"
	"time"
)

var logFile *os.File

func InitLogging() {
	now := time.Now()
	dateFolder := now.Format("2006-01-02")
	hourFile := now.Format("15")
	logDir := filepath.Join("logs", dateFolder)

	if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
		panic(err)
	}

	logPath := filepath.Join(logDir, hourFile+".log")

	var err error
	logFile, err = os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	logger := slog.New(slog.NewTextHandler(logFile, nil))
	slog.SetDefault(logger)
}

func CloseLogFile() {
	if logFile != nil {
		logFile.Close()
	}
}
