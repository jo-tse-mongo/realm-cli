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
	// TODO: make the header row bold - after everything's done since it's from the colors package
	// 	table design like in the doc
	Headers []string
	Values  []map[string]string
	// TODO: need to require the header array into the correct order - OrderedMaps in titled json to be shown in .Payload
	maxHeaderWidth map[string]int
}

func newTable(headers []string, values []map[string]string) *Table {
	if headers == nil || len(headers) == 0 ||
		values == nil || len(values) == 0 {
		return nil
	}

	maxHeaderWidth := make(map[string]int)
	for _, h := range headers {
		temp := len(h)
		if temp > maxHeaderWidth[h] {
			maxHeaderWidth[h] = temp
		}
	}

	for _, v := range values {
		for k1, v1 := range v {
			// possible new lines
			temp := len(v1)
			if temp > maxHeaderWidth[k1] {
				maxHeaderWidth[k1] = temp
			}
		}
	}

	table := &Table{
		headers,
		values,
		maxHeaderWidth,
	}
	return table
}

func (t Table) writeHeader() string {
	var tempBuilder strings.Builder
	var result strings.Builder

	for _, h := range t.Headers {
		spaces := strings.Repeat(" ", t.maxHeaderWidth[h]-len(h)+2)
		tempBuilder.WriteString("%s")
		tempBuilder.WriteString(spaces)

		result.WriteString(fmt.Sprintf(tempBuilder.String(), color.New(color.Bold).SprintFunc()(h)))
		tempBuilder.Reset()
	}
	return result.String()
}

func (t Table) writeValues() string {
	var result strings.Builder
	var tempBuilder strings.Builder

	for _, valueMap := range t.Values {
		for _, header := range t.Headers {
			spaces := strings.Repeat(" ", t.maxHeaderWidth[header]-len(valueMap[header])+2)

			tempBuilder.WriteString("%s")
			tempBuilder.WriteString(spaces)

			result.WriteString(fmt.Sprintf(tempBuilder.String(), valueMap[header]))
			tempBuilder.Reset()
		}
		result.WriteString("\n")
	}
	return result.String()
}

func (t Table) writeDividers() string {
	var divider strings.Builder
	var tempBuilder strings.Builder
	for _, header := range t.Headers {
		colLen := t.maxHeaderWidth[header]
		tempBuilder.WriteString(strings.Repeat("-", colLen))
		tempBuilder.WriteString("  ")
		divider.WriteString(tempBuilder.String())
		tempBuilder.Reset()
	}
	return divider.String()
}

func (t Table) Message() string {
	if len(t.Headers) == 0 || len(t.Values) == 0 {
		return ""
	}
	var tempBuilder strings.Builder
	tempBuilder.WriteString("\n")
	tempBuilder.WriteString(strings.Join([]string{t.writeHeader(), t.writeDividers(), t.writeValues()}, "\n"))
	return tempBuilder.String()
}

func (t Table) Payload() ([]string, map[string]interface{}, error) {
	return t.Headers, map[string]interface{} {
		logFieldDoc: t.Values,
	}, nil
}