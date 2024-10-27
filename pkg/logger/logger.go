package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/boxboxjason/jukebox/pkg/utils/fileutils"
	"github.com/boxboxjason/jukebox/pkg/utils/timeutils"
)

const (
	MAX_LOG_SIZE = 1024 * 1024 * 10 // 10 MB
)

// Setup the logger to simple STDOUT logging for now
func init() {
	debugLogger = log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime)
	infoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
	errorLogger = log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime)
	criticalLogger = log.New(os.Stdout, "CRITICAL: ", log.Ldate|log.Ltime)
	fatalLogger = log.New(os.Stdout, "FATAL: ", log.Ldate|log.Ltime)
}

// Logger instances
var (
	// Logging directories and files
	LOG_DIR        string
	LOG_FILE       string
	LOG_ROTATE_ZIP string
	// Logging levels
	LOG_LEVEL  int = 1
	LOG_LEVELS     = map[string]int{
		"DEBUG":    0,
		"INFO":     1,
		"ERROR":    2,
		"CRITICAL": 3,
		"FATAL":    4,
	}
	// Mutexes
	mu_LOG_DIR        = &sync.RWMutex{}
	mu_LOG_FILE       = &sync.RWMutex{}
	mu_LOG_ROTATE_ZIP = &sync.RWMutex{}
	mu_LOG_LEVEL      = &sync.RWMutex{}
	// Loggers
	debugLogger    *log.Logger
	infoLogger     *log.Logger
	errorLogger    *log.Logger
	criticalLogger *log.Logger
	fatalLogger    *log.Logger
)

func SetupLogger(log_dir string, log_level string) {
	// Lock the mutexes
	mu_LOG_DIR.Lock()
	mu_LOG_FILE.Lock()
	mu_LOG_ROTATE_ZIP.Lock()
	defer mu_LOG_DIR.Unlock()
	defer mu_LOG_FILE.Unlock()
	defer mu_LOG_ROTATE_ZIP.Unlock()
	// Set the logging directories and files
	LOG_DIR = log_dir
	LOG_FILE = path.Join(LOG_DIR, "server.log")
	LOG_ROTATE_ZIP = path.Join(LOG_DIR, "rotate.zip")

	// Set the logging level
	mu_LOG_LEVEL.Lock()
	defer mu_LOG_LEVEL.Unlock()
	if _, ok := LOG_LEVELS[strings.ToUpper(log_level)]; !ok {
		fmt.Println("Invalid log level:", log_level, "Defaulting to INFO")
		LOG_LEVEL = LOG_LEVELS["INFO"]
	} else {
		LOG_LEVEL = LOG_LEVELS[strings.ToUpper(log_level)]
	}

	// Check if the directory for the log directory exists
	err := os.MkdirAll(LOG_DIR, os.ModePerm)
	if err != nil {
		fmt.Println("Failed to create log directory:", err)
		os.Exit(1)
	}

	file, err := os.OpenFile(LOG_FILE, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0640)
	if err != nil {
		fmt.Println("Fatal error: Failed to open log file:", err)
		os.Exit(1)
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
	if LOG_LEVEL <= 0 {
		debugLogger.Println(v...)
	}
}

// Info logs an info message.
func Info(v ...interface{}) {
	if LOG_LEVEL <= 1 {
		infoLogger.Println(v...)
	}
}

// Error logs an error message.
func Error(v ...interface{}) {
	if LOG_LEVEL <= 2 {
		errorLogger.Println(v...)
	}
}

// Critical logs a critical message but does not exit the application.
func Critical(v ...interface{}) {
	if LOG_LEVEL <= 3 {
		criticalLogger.Println(v...)
	}
}

// Fatal logs a fatal message and exits the application.
func Fatal(v ...interface{}) {
	if LOG_LEVEL <= 4 {
		fatalLogger.Println(v...)
	}
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
