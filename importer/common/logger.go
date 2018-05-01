package common

import (
	"fmt"
)

// LogProperty contains log message metadata
type LogProperty struct {
	Name  string
	Value string
}

// Logger logs messages
type Logger interface {
	Log(level string, message string)
	Info(message string)
	Debug(message string)
	Warning(message string)
	Error(message string)
}

// ConsoleLogger logs messsages to the console
type ConsoleLogger struct{}

// NewLogger creates a new logger
func NewLogger() Logger {
	return &ConsoleLogger{}
}

// Log logs a message and additional args for a given level
func (c ConsoleLogger) Log(level string, message string) {
	m := fmt.Sprintf("%s: %s", level, message)
	fmt.Println(m)
}

// Info logs a message for the INFO level
func (c ConsoleLogger) Info(message string) {
	c.Log("INFO", message)
}

// Debug logs a message for the DEBUG level
func (c ConsoleLogger) Debug(message string) {
	c.Log("DEBUG", message)
}

// Warning logs a message for the WARNING level
func (c ConsoleLogger) Warning(message string) {
	c.Log("WARNING", message)
}

// Error logs a message and error for the ERROR level
func (c ConsoleLogger) Error(message string) {
	c.Log("ERROR", message)
}
