package terminal

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/iancoleman/orderedmap"
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
		for _, level := range allLogLevels {
			if l := len(level); l > len(longest) {
				longest = level
			}
		}
		return longest
	}()
)

// LogData produces the log data
type LogData interface {
	Message() (string, error)
	Payload() ([]string, map[string]interface{}, error)
}

// Log is a terminal log
type Log struct {
	Level LogLevel
	Time  time.Time
	TZ    *time.Location
	Data  LogData
}

// NewTextLog creates a new log with a text message
func NewTextLog(message string) Log {
	return Log{LogLevelInfo, time.Now(), time.Local, textMessage(message)}
}

// NewJSONLog creates a new log with a JSON document
func NewJSONLog(data map[string]interface{}) Log {
	return Log{LogLevelInfo, time.Now(), time.Local, jsonDocument{data}}
}

// NewTitledJSONLog creates a new log with a titled JSON document
func NewTitledJSONLog(title string, data map[string]interface{}) Log {
	return Log{LogLevelInfo, time.Now(), time.Local, titledJSONDocument{title, jsonDocument{data}}}
}

// NewErrorLog creates a new error log
func NewErrorLog(err error) Log {
	return Log{LogLevelError, time.Now(), time.Local, errorMessage{err}}
}

// Print produces the log output based on the specified format
func (l Log) Print(outputFormat OutputFormat) (string, error) {
	switch outputFormat {
	case OutputFormatText:
		return l.textLog()
	case OutputFormatJSON:
		return l.jsonOutput()
	default:
		return "", fmt.Errorf("unsupported output format type: %s", outputFormat)
	}
}

func (l Log) textLog() (string, error) {
	message, err := l.Data.Message()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(
		fmt.Sprintf("%%-%ds", len(longestLogLevel)+1)+"%s: %s",
		strings.ToUpper(string(l.Level)),
		l.Time.In(l.TZ).Format("15:04:05.000"),
		message,
	), nil
}

const (
	logFieldLevel = "level"
	logFieldTime  = "time"
)

func (l Log) jsonOutput() (string, error) {
	out := orderedmap.New()
	out.Set(logFieldLevel, l.Level)
	out.Set(logFieldTime, l.Time)

	keys, payload, err := l.Data.Payload()
	if err != nil {
		return "", err
	}
	for _, key := range keys {
		out.Set(key, payload[key])
	}

	output, outputErr := json.Marshal(out)
	return string(output), outputErr
}
