package terminal

import (
	"fmt"
	"github.com/fatih/color"
	"strings"
)

// Table
// Headers is an ordered list of strings to specify the order of printing
// Values is a map of a string to another string to print, where the key should be a value in Headers; if not, we will not print
type Table struct {
	// TODO: start the result of Message with a newline
	// TODO: make the header row bold - after everything's done since it's from the colors package
	// 	table design like in the doc
	Headers   []string
	Values    []map[string]string
	// TODO: need to require the header array into the correct order - OrderedMaps in titled json to be shown in .Payload
	// TODO: what can live in struct or what can be done as a local variable - constructor keep the max widths of the header

	// TODO: remove
	MaxLength int
	HeaderNewlines int
	ValueNewlines  int
	// TODO: change to maxHeader
	headerToWidth  map[string]int

	Divider        string
}

// TODO: change the max width of each column to be determined by the max of the values ANd the headers
// TODO print to file option??
func NewTable(headers []string, values []map[string]string) *Table {
	if headers == nil || len(headers) == 0 ||
		values == nil || len(values) == 0 {
		return nil
	}
	// TODO: move max to const ?
	maxLength := 30
	headerNewlines := 1
	valueNewlines := 1
	headerToWidth := make(map[string]int)
	for _, h := range headers {
		headerToWidth[h] = len(h)
		if len(h) > maxLength {
			temp := len(h) / maxLength
			if temp > headerNewlines {
				headerNewlines = temp
			}
			headerToWidth[h] = maxLength
		}
	}
	for _, v := range values {
		for k1, v1 := range v {
			// possible new lines
			if len(v1) > len(k1) {
				var temp int
				if len(k1) > maxLength {
					temp = len(v1) / maxLength
				} else {
					temp = len(v1) / len(k1)
				}
				if temp > valueNewlines {
					valueNewlines = temp
				}
			}
		}
	}
	table := &Table{
		headers,
		values,
		maxLength,
		headerToWidth,
		WriteDividers(headers, maxLength),
		headerNewlines,
		valueNewlines,
	}
	return table
}

// TODO; there's some similarity with writeHeader and writeValue so consolidate that
func (t Table) WriteHeader() string {
	var tempBuilder strings.Builder
	var result strings.Builder
	for i := 0; i < t.HeaderNewlines+1; i++ {
		for _, header := range t.Headers {
			toWrite := header
			if len(header) > t.MaxLength {
				start := i * t.MaxLength
				end := start + t.MaxLength

				if end > len(header) {
					end = len(header)
				}
				if start > len(header) {
					toWrite = strings.Repeat(" ", len(header))
				} else {
					toWrite = header[start:end]
				}
				if len(toWrite) < t.MaxLength {
					toWrite += strings.Repeat(" ", t.MaxLength-len(toWrite))
				}
			} else {
				if i > 0 {
					toWrite = strings.Repeat(" ", len(header))
				}
			}

			tempBuilder.WriteString("| ")
			tempBuilder.WriteString("%s")
			tempBuilder.WriteString(" |")

			result.WriteString(fmt.Sprintf(tempBuilder.String(), color.New(color.Bold).SprintFunc()(toWrite)))
			tempBuilder.Reset()
		}
		if i+1 < t.HeaderNewlines+1 {
			result.WriteString("\n")
		}
	}
	return result.String()
}

func (t Table) WriteValues() string {
	var result strings.Builder
	var tempBuilder strings.Builder

	for _, valueMap := range t.Values {
		result.WriteString("\n")
		for i := 0; i < t.ValueNewlines+1; i++ {
			for _, header := range t.Headers {
				var toWrite string
				maxWidth := t.HeaderToWidth[header]
				if value, isPresent := valueMap[header]; isPresent {
					toWrite = value
					if len(value) > maxWidth {
						start := i * maxWidth
						end := start + maxWidth

						if end > len(value) {
							end = len(value)
						}
						if start > len(value) {
							toWrite = strings.Repeat(" ", maxWidth)
						} else {
							toWrite = value[start:end]
						}
					} else {
						if i > 0 {
							toWrite = strings.Repeat(" ", len(value))
						}
					}
				}

				if len(toWrite) < maxWidth {
					toWrite += strings.Repeat(" ", maxWidth-len(toWrite))
				}

				tempBuilder.WriteString("| ")
				tempBuilder.WriteString("%s")
				tempBuilder.WriteString(" |")

				result.WriteString(fmt.Sprintf(tempBuilder.String(), color.New(color.Bold).SprintFunc()(toWrite)))
				tempBuilder.Reset()
			}

			if i+1 < t.ValueNewlines+1 {
				result.WriteString("\n")
			}
		}
		result.WriteString("\n")
		result.WriteString(t.Divider)
	}
	return result.String()
}

func WriteDividers(header []string, maxLength int) string {
	var divider strings.Builder
	var tempBuilder strings.Builder
	for _, value := range header {
		colLen := len(value)
		if colLen > maxLength {
			colLen = maxLength
		}
		tempBuilder.WriteString("+")
		tempBuilder.WriteString(strings.Repeat("-", colLen+2))
		tempBuilder.WriteString("+")
		divider.WriteString(tempBuilder.String())
		tempBuilder.Reset()
	}
	return divider.String()
}

// TODO: if value longer than longest (50) then split with newline at len(value) - 2?
func (t Table) Header() string {
	if len(t.Headers) == 0 {
		return ""
	}
	return strings.Join([]string{t.Divider, t.WriteHeader(), t.Divider}, "\n")
}

func (t Table) Value() string {
	if len(t.Values) == 0 {
		return ""
	}
	return t.WriteValues()

}

func (t Table) Message() string {
	return t.Header() + t.Value()
}
