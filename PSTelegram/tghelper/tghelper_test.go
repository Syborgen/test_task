package tghelper

import "testing"

func TestIsDate(t *testing.T) {
	testCases := []struct {
		input    string
		expected bool
	}{
		{
			input:    "123",
			expected: false,
		},
		{
			input:    "aa",
			expected: false,
		},
		{
			input:    "2022-12-05",
			expected: true,
		},
	}

	for _, testCase := range testCases {
		actual := isDate(testCase.input)
		if actual != testCase.expected {
			t.Errorf("Expected: %v, actual: %v", testCase.expected, actual)
		}
	}
}

func TestParseArguments(t *testing.T) {
	testCases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "123 456 789",
			expected: []string{"123", "456", "789"},
		},
		{
			input:    "one two three",
			expected: []string{"one", "two", "three"},
		},
		{
			input:    "2022-12-05 11:43:66 two",
			expected: []string{"2022-12-05 11:43:66", "two"},
		},
		{
			input:    "one two 2022-12-05 11:43:66 2022-12-05 11:43:66",
			expected: []string{"one", "two", "2022-12-05 11:43:66", "2022-12-05 11:43:66"},
		},
	}

	for _, testCase := range testCases {
		actual := ParseArguments(testCase.input)

		if len(actual) != len(testCase.expected) {
			t.Errorf("Expected: %v, actual: %v", testCase.expected, actual)
			continue
		}

		for i := range testCase.expected {
			if testCase.expected[i] != actual[i] {
				t.Errorf("Expected: %v, actual: %v", testCase.expected, actual)
				break
			}
		}
	}
}
