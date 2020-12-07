package terminal

import (
	"errors"
	"fmt"
	"strings"

	"github.com/fatih/color"
)

const (
	logFieldHeaders = "headers"
	logFieldData    = "data"
)

var (
	// gutter is the gap between table columns
	gutter = strings.Repeat(" ", 2)

	tableFields = []string{logFieldData, logFieldHeaders}
)

type table struct {
	headers      []string
	data         []map[string]string
	columnWidths map[string]int
}

func newTable(headers []string, data []map[string]interface{}) table {
	var t table

	if len(headers) == 0 {
		return t
	}

	t.headers = headers
	t.data = make([]map[string]string, 0, len(data))
	t.columnWidths = make(map[string]int, len(headers))

	for _, header := range headers {
		t.columnWidths[header] = len(header)
	}

	for _, row := range data {
		if len(row) == 0 {
			continue
		}
		r := make(map[string]string)
		for _, header := range t.headers {
			// convert the interface data into strings for printing
			value := parseValue(row[header])
			if width := len(value); width > t.columnWidths[header] {
				t.columnWidths[header] = width
			}
			r[header] = value
		}
		t.data = append(t.data, r)
	}
	return t
}

func (t table) Message() (string, error) {
	if err := t.validate(); err != nil {
		return "", err
	}
	return fmt.Sprintf(`
%s
%s
%s`, t.headerString(), t.dividerString(), t.dataString()), nil
}

func (t table) Payload() ([]string, map[string]interface{}, error) {
	if err := t.validate(); err != nil {
		return nil, nil, err
	}
	return tableFields, map[string]interface{}{
		logFieldHeaders: t.headers,
		logFieldData:    t.data,
	}, nil
}

func (t table) headerString() string {
	rows := make([]string, len(t.headers))
	for i, header := range t.headers {
		rows[i] = fmt.Sprintf("%s%s",
			color.New(color.Bold).SprintFunc()(header),
			strings.Repeat(" ", t.columnWidths[header]-len(header)),
		)
	}
	return strings.Join(rows, gutter)
}

func (t table) dataString() string {
	rows := make([]string, len(t.data))
	for i, dataMap := range t.data {
		cells := make([]string, len(t.headers))
		for j, header := range t.headers {
			cells[j] = fmt.Sprintf(
				"%s%s",
				dataMap[header],
				strings.Repeat(" ", t.columnWidths[header]-len(dataMap[header])),
			)
		}
		rows[i] = strings.Join(cells, gutter)
	}
	return strings.Join(rows, "\n")
}

func (t table) dividerString() string {
	dashes := make([]string, len(t.headers))
	for i, header := range t.headers {
		dashes[i] = strings.Repeat("-", t.columnWidths[header])
	}
	return strings.Join(dashes, gutter)
}

func (t table) validate() error {
	if len(t.headers) == 0 {
		return errors.New("cannot create a table without headers")
	}
	return nil
}

func parseValue(value interface{}) string {
	parsed := ""
	switch v := value.(type) {
	case nil: // leave zero-value
	case string:
		parsed = v
	case fmt.Stringer:
		parsed = v.String()
	default:
		parsed = fmt.Sprintf("%v", v)
	}
	return parsed
}
