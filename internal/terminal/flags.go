package terminal

// OutputFormat is the terminal output format
type OutputFormat string

// set of supported terminal output formats
const (
	OutputFormatJSON OutputFormat = "json"
	OutputFormatText OutputFormat = "text"
)

// set of supported terminal flag defaults
const (
	DefaultOutputFormat = string(OutputFormatText)
)
