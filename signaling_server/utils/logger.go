package utils

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/fatih/color"
)

type LogLevel string

const (
	SUCCESS LogLevel = "SUCCESS"
	ERROR   LogLevel = "ERROR"
	INFO    LogLevel = "INFO"
	DEBUG   LogLevel = "DEBUG"
)

var levelPrefix = map[LogLevel]string{
	SUCCESS: "[+]",
	ERROR:   "[-]",
	INFO:    "[!]",
	DEBUG:   "[~]",
}

type Logger struct {
	enableDebug bool
	output      io.Writer
}

// New creates a new Logger that writes to stdout
func New(enableDebug bool) *Logger {
	return &Logger{
		enableDebug: enableDebug,
		output:      os.Stdout,
	}
}

// formatLogLine builds the full log message
func formatLogLine(level LogLevel, msg string) string {
	timestamp := time.Now().Format("15:04:05.000")
	prefix := levelPrefix[level]
	return fmt.Sprintf("[%s] %s [%s] %s", timestamp, prefix, level, msg)
}

func (l *Logger) log(level LogLevel, msg string) {
	if level == DEBUG && !l.enableDebug {
		return
	}
	line := formatLogLine(level, msg)

	switch level {
	case SUCCESS:
		color.New(color.FgGreen).Fprintln(l.output, line)
	case ERROR:
		color.New(color.FgRed).Fprintln(l.output, line)
	case INFO:
		color.New(color.FgYellow).Fprintln(l.output, line)
	case DEBUG:
		color.New(color.FgCyan).Fprintln(l.output, line)
	default:
		fmt.Fprintln(l.output, line)
	}
}

// Public methods
func (l *Logger) Success(msg string) { l.log(SUCCESS, msg) }
func (l *Logger) Error(msg string)   { l.log(ERROR, msg) }
func (l *Logger) Info(msg string)    { l.log(INFO, msg) }
func (l *Logger) Debug(msg string)   { l.log(DEBUG, msg) }

// Public methods (formatted)
func (l *Logger) Successf(msg string, args ...any) { l.log(SUCCESS, fmt.Sprintf(msg, args...)) }
func (l *Logger) Errorf(msg string, args ...any)   { l.log(ERROR, fmt.Sprintf(msg, args...)) }
func (l *Logger) Infof(msg string, args ...any)    { l.log(INFO, fmt.Sprintf(msg, args...)) }
func (l *Logger) Debugf(msg string, args ...any)   { l.log(DEBUG, fmt.Sprintf(msg, args...)) }

// LogWithContext adds structured metadata to logs
func (l *Logger) LogWithContext(level LogLevel, component, action, sessionID, extra string) {
	if level == DEBUG && !l.enableDebug {
		return
	}
	msg := formatContext(component, action, sessionID, extra)
	l.log(level, msg)
}

func formatContext(component, action, sessionID, extra string) string {
	base := fmt.Sprintf("%s | %s | Session: %s", component, action, sessionID)
	if extra != "" {
		base += " | " + extra
	}
	return base
}
