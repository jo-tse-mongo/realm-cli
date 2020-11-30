package terminal

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/10gen/realm-cli/internal/utils/test/assert"
)

func TestLogConstructor(t *testing.T) {
	for _, tc := range []struct {
		ctor              string
		log               Log
		expectedLevel     LogLevel
		exepctedMessenger Messenger
	}{
		{
			ctor:              "NewTextLog",
			log:               NewTextLog("oh yeah"),
			expectedLevel:     LogLevelInfo,
			exepctedMessenger: TextMessage{"oh yeah"},
		},
		{
			ctor:              "NewJSONLog",
			log:               NewJSONLog(map[string]interface{}{"a": "ayyy"}),
			expectedLevel:     LogLevelInfo,
			exepctedMessenger: JSONDocument{map[string]interface{}{"a": "ayyy"}},
		},
		{
			ctor:              "NewTitledJSONLog",
			log:               NewTitledJSONLog("Test Title", map[string]interface{}{"a": "ayyy"}),
			expectedLevel:     LogLevelInfo,
			exepctedMessenger: TitledJSONDocument{JSONDocument{map[string]interface{}{"a": "ayyy"}}, "Test Title"},
		},
		{
			ctor:              "NewErrorLog",
			log:               NewErrorLog(errors.New("oh noz")),
			expectedLevel:     LogLevelError,
			exepctedMessenger: errorMessage{errors.New("oh noz")},
		},
	} {
		t.Run(fmt.Sprintf("%s should create the expected Log", tc.ctor), func(t *testing.T) {
			time.Sleep(1 * time.Millisecond)
			assert.True(t, time.Now().After(tc.log.Time), "now should be later than the log's timestamp")
			assert.Equal(t, tc.expectedLevel, tc.log.Level)
			assert.Equal(t, tc.exepctedMessenger, tc.log.Messenger)
		})
	}
}

func TestLogMessage(t *testing.T) {
	for _, tc := range []struct {
		description  string
		outputFormat OutputFormat
		level        LogLevel
		messenger    Messenger
		expectedMsg  string
	}{
		{
			description:  "JSON",
			outputFormat: OutputFormatJSON,
			level:        LogLevelInfo,
			messenger:    TextMessage{"this is a test log"},
			expectedMsg:  `{"level":"info","time":"1989-06-22T11:54:00Z","message":"this is a test log"}`,
		},
		{
			description:  "text",
			outputFormat: OutputFormatText,
			level:        LogLevelInfo,
			messenger:    TextMessage{"this is a test log"},
			expectedMsg:  "INFO  07:54:00.000: this is a test log",
		},
		{
			description:  "JSON",
			outputFormat: OutputFormatJSON,
			level:        LogLevelInfo,
			messenger:    JSONDocument{map[string]interface{}{"a": true, "b": 1, "c": "sea"}},
			expectedMsg:  `{"level":"info","time":"1989-06-22T11:54:00Z","message":"{\"a\":true,\"b\":1,\"c\":\"sea\"}"}`,
		},
		{
			description:  "text",
			outputFormat: OutputFormatText,
			level:        LogLevelInfo,
			messenger:    JSONDocument{map[string]interface{}{"a": true, "b": 1, "c": "sea"}},
			expectedMsg: `INFO  07:54:00.000: {
  "a": true,
  "b": 1,
  "c": "sea"
}`,
		},
		{
			outputFormat: OutputFormatJSON,
			level:        LogLevelError,
			messenger:    errorMessage{errors.New("something bad happened")},
			expectedMsg:  `{"level":"error","time":"1989-06-22T11:54:00Z","message":"something bad happened"}`,
		},
		{
			outputFormat: OutputFormatText,
			level:        LogLevelError,
			messenger:    errorMessage{errors.New("something bad happened")},
			expectedMsg:  `ERROR 07:54:00.000: something bad happened`,
		},
	} {
		t.Run(fmt.Sprintf("With %s output format, %T should return the expected message", tc.description, tc.messenger), func(t *testing.T) {
			log := Log{
				tc.level,
				time.Date(1989, 6, 22, 11, 54, 0, 0, time.UTC),
				tc.messenger,
			}

			msg, err := log.Message(tc.outputFormat)
			assert.Nil(t, err)
			assert.Equal(t, tc.expectedMsg, msg)
		})
	}

	t.Run("Should return an error with an unknown output format", func(t *testing.T) {
		log := Log{
			LogLevelInfo,
			time.Date(1989, 6, 22, 11, 54, 0, 0, time.UTC),
			TextMessage{"this is a test log"},
		}

		_, err := log.Message(OutputFormat("eggcorn"))
		assert.Equal(t, errors.New("unsupported output format type: eggcorn"), err)
	})

	t.Run("Should propagate an error that occurs while producing the messenger message", func(t *testing.T) {
		failLog := Log{LogLevelInfo, time.Now(), failMessage{}}
		_, err := failLog.Message(OutputFormatText)
		assert.Equal(t, errFailMessage, err)
	})
}

var errFailMessage = errors.New("something bad happened")

type failMessage struct{}

func (f failMessage) Message(outputFormat OutputFormat) (string, error) {
	return "", errFailMessage
}
