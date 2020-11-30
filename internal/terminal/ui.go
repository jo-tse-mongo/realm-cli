package terminal

import (
	"fmt"
	"io"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
)

// UI is a terminal UI
type UI interface {
	AskOne(prompt survey.Prompt, answer interface{}) error
	Print(logs ...Log) error
}

// NewUI creates a new terminal UI
func NewUI(config UIConfig, in io.Reader, out, err io.Writer) UI {
	noColor := config.DisableColors
	if config.OutputFormat == OutputFormatJSON {
		noColor = true
	}
	color.NoColor = noColor

	return &ui{
		config,
		fdReader{in},
		fdWriter{out},
		err,
	}
}

type ui struct {
	config UIConfig
	in     fdReader
	out    fdWriter
	err    io.Writer
}

func (ui *ui) AskOne(prompt survey.Prompt, answer interface{}) error {
	return survey.AskOne(
		prompt,
		answer,
		survey.WithStdio(ui.in, ui.out, ui.err),
	)
}

func (ui *ui) Print(logs ...Log) error {
	for _, log := range logs {
		output, outputErr := log.Print(ui.config.OutputFormat)
		if outputErr != nil {
			return outputErr
		}

		var writer io.Writer
		switch log.Level {
		case LogLevelError:
			writer = ui.err
		default:
			writer = ui.out
		}

		if _, err := fmt.Fprintln(writer, output); err != nil {
			return err
		}
	}
	return nil
}

// UIConfig holds the global config for the CLI ui
type UIConfig struct {
	DisableColors bool
	OutputFormat  OutputFormat
	OutputTarget  string
}

// FileDescriptor is a file descriptor
type FileDescriptor interface {
	Fd() uintptr
}

// fdReader wraps an io.Reader and exposes the FileDesriptor interface on it
// the underlying io.Reader's Fd() implementation will be used if it exists
type fdReader struct {
	io.Reader
}

func (r fdReader) Fd() uintptr {
	if fd, ok := r.Reader.(FileDescriptor); ok {
		return fd.Fd()
	}
	return 0
}

// fdWriter wraps an io.Writer and exposes the FileDesriptor interface on it
// the underlying io.Writer's Fd() implementation will be used if it exists
type fdWriter struct {
	io.Writer
}

func (w fdWriter) Fd() uintptr {
	if fd, ok := w.Writer.(FileDescriptor); ok {
		return fd.Fd()
	}
	return 0
}
