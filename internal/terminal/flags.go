package terminal

// OutputFormat is the terminal output format
type OutputFormat string

// set of supported terminal output formats
const (
	OutputFormatText OutputFormat = "" // empty string allows for this to be flag's default value
	OutputFormatJSON OutputFormat = "json"
)

const (
	DefaultOutputFormat = string(OutputFormatText)
)
