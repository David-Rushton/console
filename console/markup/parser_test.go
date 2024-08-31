package markup_test

import (
	"testing"

	"example.com/console/console/markup"
)

const passed = "✅"
const failed = "❌"

func TestParse(t *testing.T) {
	testCases := []struct {
		given    string
		expected string
	}{
		{
			given:    "text with <red>red</> section.",
			expected: "text with \x1b[31mred\x1b[0m section.",
		},
		{
			given:    "text with <green>green</> section.",
			expected: "text with \x1b[32mgreen\x1b[0m section.",
		},
		{
			given:    "text with <yellow>yellow</> section.",
			expected: "text with \x1b[33myellow\x1b[0m section.",
		},
	}

	t.Log("Should convert markup to an ANSI escaped string.")

	for _, testCase := range testCases {
		t.Logf("\tWhen input: %v", testCase.given)

		actual := markup.Parse(testCase.given)
		if actual != testCase.expected {
			t.Errorf("\t\tFailed %v\n\t\tExpected %v\n\t\tActual %v", failed, testCase.expected, actual)
		} else {
			t.Logf("\t\tPassed %v\n\n\tExpected matches actual for: %v", passed, testCase.expected)
		}
	}
}
