package main

import (
	"strings"
	"testing"
)

func TestExport(t *testing.T) {
	maxRecords = 1000
	type TestCase struct {
		input    string
		expected string
		name     string
	}
	tests := []TestCase{
		{"0xf503017d7baf7fbc0fff7492b751025c6a78179b", "0xf503017d7baf7fbc0fff7492b751025c6a78179b has 1000 appearances", "basic address"},
		{"0xf503017d7baf7fbc0fff7492b751025c6A78179B", "0xf503017d7baf7fbc0fff7492b751025c6a78179b has 1000 appearances", "address w/ capitals"},
		{"0xf503017d7baf7fbc0fff7492b751025c6a78179", "Invalid address or ENS name: 0xf503017d7baf7fbc0fff7492b751025c6a78179", "bogus address"},
		{"trueblocks.eth", "trueblocks.eth (0xf503017d7baf7fbc0fff7492b751025c6a78179b) has 1000 appearances", "ENS name"},
		{"trueblock\n.eth", "trueblock .eth (0x0) has 0 appearances", "invalid ENS name"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var a App
			actual := strings.Replace(a.Export(test.input), "\n", " ", -1)
			if actual != test.expected {
				t.Errorf("TestExport(%s): expected %s, actual %s", test.input, test.expected, actual)
			}
		})
	}
}
