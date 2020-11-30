package terminal

// OutputFormat is the terminal output format
type OutputFormat string

func (of *OutputFormat) String() string { return string(*of) }

func (of *OutputFormat) Type() string { return "terminal.OutputFormat" }

func (of *OutputFormat) Set(val string) error {
	*of = OutputFormat(val)
	return nil
}

// set of supported terminal output formats
const (
	OutputFormatText OutputFormat = "" // zero-valued to be flag's default
	OutputFormatJSON OutputFormat = "json"
)
