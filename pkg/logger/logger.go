package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"sync"
	"time"
)

const (
	MAX_LOG_SIZE           = 1024 * 1024 * 10 // 10 MB
	DEFAULT_CHANNEL_BUFFER = 1000
)

var (
	logChannel  chan string
	stopLogging chan struct{}
	logWg       sync.WaitGroup
	logFile     *os.File
	logMutex    sync.Mutex
	logLevel    int = 1 // Default log level: INFO
	logFilePath string
	loggerReady sync.WaitGroup
)

// Log levels
const (
	DEBUG = iota
	INFO
	ERROR
	CRITICAL
	FATAL
)

func SetupLogger(logDir string, level string) {
	var err error

	// Map of log levels
	logLevels := map[string]int{
		"DEBUG":    DEBUG,
		"INFO":     INFO,
		"ERROR":    ERROR,
		"CRITICAL": CRITICAL,
		"FATAL":    FATAL,
	}

	if l, ok := logLevels[level]; ok {
		logLevel = l
	} else {
		fmt.Println("Invalid log level. Defaulting to INFO.")
		logLevel = INFO
	}

	logFilePath = path.Join(logDir, "jukebox.log")

	// Create log directory if it doesn't exist
	if err = os.MkdirAll(logDir, os.ModePerm); err != nil {
		log.Fatalf("Failed to create log directory: %v", err)
	}

	logFile, err = os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	// Initialize the logging channel and control signals
	logChannel = make(chan string, DEFAULT_CHANNEL_BUFFER)
	stopLogging = make(chan struct{})

	// Prepare loggerReady wait group
	loggerReady.Add(1)

	// Start the asynchronous logger goroutine
	go asyncLogger()

	// Wait for the logger to be ready
	loggerReady.Wait()
}

func asyncLogger() {
	defer logWg.Done()
	logWg.Add(1)

	logger := log.New(io.MultiWriter(logFile, os.Stdout), "", log.Ldate|log.Ltime)

	// Signal that the logger is ready
	loggerReady.Done()

	for {
		select {
		case logMessage := <-logChannel:
			logger.Println(logMessage)
			checkLogRotation() // Rotate log if needed
		case <-stopLogging:
			// Flush remaining logs
			for len(logChannel) > 0 {
				logger.Println(<-logChannel)
			}
			return
		}
	}
}

func checkLogRotation() {
	logMutex.Lock()
	defer logMutex.Unlock()

	fileInfo, err := logFile.Stat()
	if err != nil {
		fmt.Println("Failed to get log file info:", err)
		return
	}

	if fileInfo.Size() >= MAX_LOG_SIZE {
		// Rotate log file
		newFileName := fmt.Sprintf("%s.%d", logFilePath, time.Now().Unix())
		if err := os.Rename(logFilePath, newFileName); err != nil {
			fmt.Println("Failed to rotate log file:", err)
			return
		}

		logFile.Close()
		logFile, err = os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			fmt.Println("Failed to open new log file:", err)
			return
		}
	}
}

func Debug(v ...interface{}) {
	if logLevel <= DEBUG {
		logMessage("DEBUG", v...)
	}
}

func Info(v ...interface{}) {
	if logLevel <= INFO {
		logMessage("INFO", v...)
	}
}

func Error(v ...interface{}) {
	if logLevel <= ERROR {
		logMessage("ERROR", v...)
	}
}

func Critical(v ...interface{}) {
	if logLevel <= CRITICAL {
		logMessage("CRITICAL", v...)
	}
}

func Fatal(v ...interface{}) {
	if logLevel <= FATAL {
		logMessage("FATAL", v...)
		ShutdownLogger()
		os.Exit(1)
	}
}

func logMessage(level string, v ...interface{}) {
	msg := fmt.Sprintf("%s: %s", level, fmt.Sprintln(v...))
	if len(msg) > 0 && msg[len(msg)-1] == '\n' {
		msg = msg[:len(msg)-1]
	}
	select {
	case logChannel <- msg:
	default:
		fmt.Println("Log channel full, dropping log:", msg)
	}
}

func ShutdownLogger() {
	close(stopLogging)
	logWg.Wait()
	if logFile != nil {
		logFile.Close()
	}
}
