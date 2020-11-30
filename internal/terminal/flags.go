package terminal

// OutputFormat is the terminal output format
type OutputFormat string

func (of OutputFormat) Set(val string) error {
	of = OutputFormat(val)
	return nil
}

func (of OutputFormat) String() string { return string(of) }

func (of OutputFormat) Type() string { return "terminal.OutputFormat" }

const (
	// set of supported terminal output formats
	OutputFormatJSON OutputFormat = "json"
	OutputFormatText OutputFormat = "" // empty string allows for this to be flag's default value
)
