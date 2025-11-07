// Package logger provides a centralized logging facility for gojekyll.
package logger

import (
	"fmt"
	"io"
	"os"
	"sync"
)

// Level represents the severity of a log message.
type Level int

const (
	// DebugLevel for detailed debugging information
	DebugLevel Level = iota
	// InfoLevel for general informational messages
	InfoLevel
	// WarnLevel for warning messages
	WarnLevel
	// ErrorLevel for error messages
	ErrorLevel
)

// Logger is the interface for logging operations.
type Logger struct {
	level      Level
	quiet      bool
	out        io.Writer
	err        io.Writer
	labelWidth int
	mu         sync.Mutex
}

var defaultLogger = &Logger{
	level: InfoLevel,
	quiet: false,
	out:   os.Stdout,
	err:   os.Stderr,
}

// New creates a new Logger instance.
func New() *Logger {
	return &Logger{
		level: InfoLevel,
		quiet: false,
		out:   os.Stdout,
		err:   os.Stderr,
	}
}

// Default returns the default logger instance.
func Default() *Logger {
	return defaultLogger
}

// SetLevel sets the minimum log level.
func (l *Logger) SetLevel(level Level) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

// SetQuiet enables or disables quiet mode.
// In quiet mode, Info and Debug messages are suppressed.
func (l *Logger) SetQuiet(quiet bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.quiet = quiet
}

// SetOutput sets the output destination for Info and Debug messages.
func (l *Logger) SetOutput(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.out = w
}

// SetErrorOutput sets the output destination for Warning and Error messages.
func (l *Logger) SetErrorOutput(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.err = w
}

// Debug logs a debug message.
func (l *Logger) Debug(msg string, args ...interface{}) {
	l.log(DebugLevel, msg, args...)
}

// Info logs an informational message.
func (l *Logger) Info(msg string, args ...interface{}) {
	l.log(InfoLevel, msg, args...)
}

// Warn logs a warning message.
func (l *Logger) Warn(msg string, args ...interface{}) {
	l.log(WarnLevel, msg, args...)
}

// Error logs an error message.
func (l *Logger) Error(msg string, args ...interface{}) {
	l.log(ErrorLevel, msg, args...)
}

// Label logs a message with a label prefix (used for formatted output).
// The label is left-padded to maintain alignment.
func (l *Logger) Label(label string, msg string, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.quiet {
		return
	}

	if len(label) > l.labelWidth {
		l.labelWidth = len(label)
	}

	formatted := fmt.Sprintf(msg, args...)
	fmt.Fprintf(l.out, "%*s %s\n", l.labelWidth, label, formatted)
}

// Println logs arguments separated by spaces with a newline.
func (l *Logger) Println(args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.quiet {
		return
	}

	fmt.Fprintln(l.out, args...)
}

// Printf logs a formatted message.
func (l *Logger) Printf(format string, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.quiet {
		return
	}

	fmt.Fprintf(l.out, format, args...)
}

func (l *Logger) log(level Level, msg string, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	// Check if we should suppress this message
	if level < l.level {
		return
	}

	if l.quiet && (level == InfoLevel || level == DebugLevel) {
		return
	}

	// Choose output destination
	out := l.out
	if level >= WarnLevel {
		out = l.err
	}

	// Format and write the message
	formatted := fmt.Sprintf(msg, args...)
	if len(args) > 0 {
		fmt.Fprintln(out, formatted)
	} else {
		fmt.Fprintln(out, msg)
	}
}

// Package-level convenience functions that use the default logger

// SetLevel sets the minimum log level on the default logger.
func SetLevel(level Level) {
	defaultLogger.SetLevel(level)
}

// SetQuiet enables or disables quiet mode on the default logger.
func SetQuiet(quiet bool) {
	defaultLogger.SetQuiet(quiet)
}

// Debug logs a debug message using the default logger.
func Debug(msg string, args ...interface{}) {
	defaultLogger.Debug(msg, args...)
}

// Info logs an informational message using the default logger.
func Info(msg string, args ...interface{}) {
	defaultLogger.Info(msg, args...)
}

// Warn logs a warning message using the default logger.
func Warn(msg string, args ...interface{}) {
	defaultLogger.Warn(msg, args...)
}

// Error logs an error message using the default logger.
func Error(msg string, args ...interface{}) {
	defaultLogger.Error(msg, args...)
}

// Label logs a message with a label prefix using the default logger.
func Label(label string, msg string, args ...interface{}) {
	defaultLogger.Label(label, msg, args...)
}

// Println logs arguments using the default logger.
func Println(args ...interface{}) {
	defaultLogger.Println(args...)
}

// Printf logs a formatted message using the default logger.
func Printf(format string, args ...interface{}) {
	defaultLogger.Printf(format, args...)
}
