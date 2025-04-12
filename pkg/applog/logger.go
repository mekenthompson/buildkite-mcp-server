package applog

import (
	"context"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

// contextKey is a private type for context keys to avoid collisions
type contextKey int

// loggerKey is the key for logger values in contexts.
const loggerKey contextKey = 0

// AppLogger wraps a slog.Logger with platform-specific file handling
type AppLogger struct {
	*slog.Logger
	appName string
	logFile *os.File
}

// Close closes the logger's file
func (l *AppLogger) Close() error {
	return l.logFile.Close()
}

// GetLogFilePath returns the full path to the log file
func (l *AppLogger) GetLogFilePath() (string, error) {
	return GetAppDataDir(l.appName)
}

// createLogFile creates a log file with a processed template name
func createLogFile(appName, logFileNameTemplate string) (*os.File, string, error) {
	logDir, err := GetAppDataDir(appName)
	if err != nil {
		return nil, "", err
	}

	// Process the template if it contains any template placeholders
	logFileName := logFileNameTemplate
	if strings.Contains(logFileNameTemplate, "{{") {
		logFileName = processLogFileNameTemplate(logFileNameTemplate)
	}

	logPath := filepath.Join(logDir, logFileName)
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, "", err
	}

	return file, logPath, nil
}

// NewTextLogger creates a new slog logger with TextHandler that writes to the platform-specific location
func NewTextLogger(appName, logFileNameTemplate string, opts *slog.HandlerOptions) (*AppLogger, error) {
	file, _, err := createLogFile(appName, logFileNameTemplate)
	if err != nil {
		return nil, err
	}

	handler := slog.NewTextHandler(file, opts)
	logger := slog.New(handler)

	return &AppLogger{
		Logger:  logger,
		appName: appName,
		logFile: file,
	}, nil
}

// NewJSONLogger creates a new slog logger with JSONHandler that writes to the platform-specific location
func NewJSONLogger(appName, logFileNameTemplate string, opts *slog.HandlerOptions) (*AppLogger, error) {
	file, _, err := createLogFile(appName, logFileNameTemplate)
	if err != nil {
		return nil, err
	}

	handler := slog.NewJSONHandler(file, opts)
	logger := slog.New(handler)

	return &AppLogger{
		Logger:  logger,
		appName: appName,
		logFile: file,
	}, nil
}

// SetDefaultLogger sets the slog default logger to use the platform-specific location
func SetDefaultLogger(appName, logFileNameTemplate string, jsonFormat bool, opts *slog.HandlerOptions) (*AppLogger, error) {
	var logger *AppLogger
	var err error

	if jsonFormat {
		logger, err = NewJSONLogger(appName, logFileNameTemplate, opts)
	} else {
		logger, err = NewTextLogger(appName, logFileNameTemplate, opts)
	}

	if err != nil {
		return nil, err
	}

	// Set the default logger
	slog.SetDefault(logger.Logger)

	return logger, nil
}

// MultiLogger allows writing to multiple outputs (both file and stderr for example)
type MultiLogger struct {
	*slog.Logger
	appName string
	logFile *os.File
}

// NewMultiLogger creates a logger that writes to both the platform-specific log file
// and another writer (like os.Stderr)
func NewMultiLogger(appName, logFileNameTemplate string, additionalWriter io.Writer, jsonFormat bool, opts *slog.HandlerOptions) (*MultiLogger, error) {
	file, _, err := createLogFile(appName, logFileNameTemplate)
	if err != nil {
		return nil, err
	}

	// Create a multi-writer that sends output to both the file and the additional writer
	multiWriter := io.MultiWriter(file, additionalWriter)

	var handler slog.Handler
	if jsonFormat {
		handler = slog.NewJSONHandler(multiWriter, opts)
	} else {
		handler = slog.NewTextHandler(multiWriter, opts)
	}

	logger := slog.New(handler)

	return &MultiLogger{
		Logger:  logger,
		appName: appName,
		logFile: file,
	}, nil
}

// Close closes the logger's file
func (l *MultiLogger) Close() error {
	return l.logFile.Close()
}

// --- Context Integration ---

// WithLogger returns a new Context that carries the provided logger
func WithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

// Ctx returns the Logger stored in ctx, if any.
// If no logger is found, it returns the default logger.
func Ctx(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(loggerKey).(*slog.Logger); ok {
		return logger
	}
	return slog.Default()
}

// WithContextAttrs enriches a logger with values from the context
// This is a helper function but users are free to use the logger directly
func WithContextAttrs(ctx context.Context, logger *slog.Logger, attrs ...any) *slog.Logger {
	// Extract common values from context if they exist
	// Users can customize this based on their context keys
	if reqID, ok := ctx.Value("request_id").(string); ok {
		logger = logger.With("request_id", reqID)
	}
	if userID, ok := ctx.Value("user_id").(string); ok {
		logger = logger.With("user_id", userID)
	}
	// Add any additional attributes
	if len(attrs) > 0 {
		logger = logger.With(attrs...)
	}
	return logger
}
