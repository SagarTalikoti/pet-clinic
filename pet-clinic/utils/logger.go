package utils

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Log = logrus.New()

func InitLogger() {

	fileLogger := &lumberjack.Logger{
		Filename:   "app.log", // File to store logs
		MaxSize:    5,         // Max size in MB before rotation
		MaxBackups: 3,         // Keep last 3 old logs
		MaxAge:     28,        // Delete older than 28 days
		Compress:   true,      // Compress rotated logs
	}

	// Combine outputs: write logs to both console and file
	multiWriter := io.MultiWriter(os.Stdout, fileLogger)

	// Set log output
	Log.SetOutput(multiWriter)

	// Use structured JSON format with timestamps
	Log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// Set minimum log level
	// This means all logs from DEBUG → ERROR will be shown
	Log.SetLevel(logrus.DebugLevel)

	// Test logs to confirm setup
	Log.Debug("Logger initialized in DEBUG mode (verbose output)")
	Log.Info("Logger initialized successfully")
	Log.Warn("Logger rotation is active")
	Log.Error("Test ERROR log (for verification only)")
}
