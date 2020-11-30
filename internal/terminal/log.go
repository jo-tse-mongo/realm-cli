package terminal

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// LogLevel is the level of a terminal log
type LogLevel string

// set of supported log levels
const (
	LogLevelInfo  LogLevel = "info"
	LogLevelError LogLevel = "error"
)

var (
	allLogLevels = []LogLevel{LogLevelInfo, LogLevelError}

	longestLogLevel = func() LogLevel {
		var longest LogLevel
		var longestLength int
		for _, level := range allLogLevels {
			if l := len(level); l > longestLength {
				longest = level
				longestLength = l
			}
		}
		return longest
	}()
)

// Messenger produces a message to display in the UI
type Messenger interface {
	Message(outputFormat OutputFormat) (string, error)
}

// Log is a terminal log
type Log struct {
	Level LogLevel
	Time  time.Time
	Messenger
}

// NewTextLog creates a new log with a text message
func NewTextLog(message string) Log {
	return Log{LogLevelInfo, time.Now(), TextMessage{message}}
}

// NewJSONLog creates a new log with a JSON document
func NewJSONLog(data map[string]interface{}) Log {
	return Log{LogLevelInfo, time.Now(), JSONDocument{data}}
}

// NewTitledJSONLog creates a new log with a titled JSON document
func NewTitledJSONLog(title string, data map[string]interface{}) Log {
	return Log{LogLevelInfo, time.Now(), TitledJSONDocument{JSONDocument{data}, title}}
}

// NewErrorLog creates a new error log
func NewErrorLog(err error) Log {
	return Log{LogLevelError, time.Now(), errorMessage{err}}
}

func (l Log) Message(outputFormat OutputFormat) (string, error) {
	msg, msgErr := l.Messenger.Message(outputFormat)
	if msgErr != nil {
		return "", msgErr
	}

	switch outputFormat {
	case OutputFormatJSON:
		return jsonLog(l, msg)
	case OutputFormatText:
		return textLog(l, msg), nil
	default:
		return "", fmt.Errorf("unsupported output format type: %s", outputFormat)
	}
}

type errorMessage struct {
	error
}

func (e errorMessage) Message(outputFormat OutputFormat) (string, error) {
	return e.Error(), nil
}

type logPayload struct {
	Level   LogLevel  `json:"level"`
	Time    time.Time `json:"time"`
	Message string    `json:"message"`
}

func jsonLog(l Log, msg string) (string, error) {
	data, err := json.Marshal(logPayload{l.Level, l.Time, msg})
	return string(data), err
}

func textLog(l Log, msg string) string {
	return fmt.Sprintf(fmt.Sprintf("%%-%ds", len(longestLogLevel)+1)+"%s: %s",
		strings.ToUpper(string(l.Level)),
		l.Time.In(time.Local).Format("15:04:05.000"),
		msg,
	)
}
