package logger

import (
	"io"
	"log"
	"os"
	"path"

	"github.com/boxboxjason/jukebox/pkg/static"
	"github.com/boxboxjason/jukebox/pkg/utils/fileutils"
	"github.com/boxboxjason/jukebox/pkg/utils/timeutils"
)

const (
	MAX_LOG_SIZE = 1024 * 1024 * 10 // 10 MB
)

// Logger instances
var (
	LOG_DIR        = static.BuildJukeboxPath("logs")
	LOG_FILE       = path.Join(LOG_DIR, "server.log")
	LOG_ROTATE_ZIP = path.Join(LOG_DIR, "rotate.zip")
	debugLogger    *log.Logger
	infoLogger     *log.Logger
	errorLogger    *log.Logger
	criticalLogger *log.Logger
	fatalLogger    *log.Logger
)

func init() {
	// Check if the directory for the log directory exists
	if _, err := os.Stat(LOG_DIR); os.IsNotExist(err) {
		os.Mkdir(LOG_DIR, os.ModePerm)
	}

	file, err := os.OpenFile(LOG_FILE, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file:", err)
	}

	multi := io.MultiWriter(file, os.Stdout)

	// Initialize loggers
	debugLogger = log.New(multi, "DEBUG: ", log.Ldate|log.Ltime)
	infoLogger = log.New(multi, "INFO: ", log.Ldate|log.Ltime)
	errorLogger = log.New(multi, "ERROR: ", log.Ldate|log.Ltime)
	criticalLogger = log.New(multi, "CRITICAL: ", log.Ldate|log.Ltime)
	fatalLogger = log.New(multi, "FATAL: ", log.Ldate|log.Ltime)
}

// Debug logs a debug message.
func Debug(v ...interface{}) {
	debugLogger.Println(v...)
}

// Info logs an info message.
func Info(v ...interface{}) {
	infoLogger.Println(v...)
}

// Error logs an error message.
func Error(v ...interface{}) {
	errorLogger.Println(v...)
}

// Critical logs a critical message but does not exit the application.
func Critical(v ...interface{}) {
	criticalLogger.Println(v...)
}

// Fatal logs a fatal message and exits the application.
func Fatal(v ...interface{}) {
	fatalLogger.Fatalln(v...)
	os.Exit(1)
}

func isLogFileFull() bool {
	file, err := os.Open(LOG_FILE)
	if err != nil {
		log.Fatalln("Failed to open log file:", err)
	}

	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatalln("Failed to get log file info:", err)
	}

	return fileInfo.Size() >= MAX_LOG_SIZE
}

// RotateLogFile rotates the log file if it is full.
func RotateLogFile() {
	if isLogFileFull() {
		// Rename the current log file
		rotate_filename := timeutils.GetDatetimeString() + ".log"
		rotate_path := path.Join(LOG_DIR, rotate_filename)
		err := os.Rename(LOG_FILE, rotate_path)
		if err != nil {
			Error("Failed to rotate log file:", err)
			return
		}

		// Create a new log file
		file, err := os.OpenFile(LOG_FILE, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			Error("Failed to open the new log file:", err)
			return
		}

		// Update the loggers
		multi := io.MultiWriter(file, os.Stdout)
		debugLogger.SetOutput(multi)
		infoLogger.SetOutput(multi)
		errorLogger.SetOutput(multi)
		criticalLogger.SetOutput(multi)
		fatalLogger.SetOutput(multi)

		// Compress the rotated log file
		err = fileutils.CompressFiles(LOG_ROTATE_ZIP, []string{rotate_path})
		if err != nil {
			Error("Failed to compress the rotated log file:", err)
		}

		// Remove the rotated log file
		err = os.Remove(rotate_path)
		if err != nil {
			Error("Failed to remove the rotated log file:", err)
		}
	}
}
