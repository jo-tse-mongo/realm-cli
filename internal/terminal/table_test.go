package terminal

import (
	"testing"
)

func TestTable(t *testing.T) {

	t.Run("Table should work", func(t *testing.T) {
		testHeader := []string{"make this a really long header or something is this not more than 30 characters yet i think it is", "this", "header"}
		testValue := map[string]string {
			"make this a really long header or something is this not more than 30 characters yet i think it is": "hello",
			"this": "yeee",
			"header": "test this large string for something lolol whatever!",
		}
		testValue2 := map[string]string {
			"make this a really long header or something is this not more than 30 characters yet i think it is": "hello2",
			"this": "yee2",
			"header": "test this large string for something lolol whatever!1312312312",
		}
		testValues := []map[string]string{testValue, testValue2}
		table := NewTable(testHeader, testValues)

		print(table.String())
	})
}
