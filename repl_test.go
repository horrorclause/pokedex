package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {

	// Struct creation for various test cases
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    " hello world ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "HEllo world HOW ArE you ",
			expected: []string{"hello", "world", "how", "are", "you"},
		},
		{
			input:    "thats a BULBasaur  ",
			expected: []string{"thats", "a", "bulbasaur"},
		},
		{
			input:    "  whats that  BEHIND the BuShes  ",
			expected: []string{"whats", "that", "behind", "the", "bushes"},
		},
		{
			input:    "pikachu",
			expected: []string{"pikachu"},
		},
		{
			input:    "hello    world",
			expected: []string{"hello", "world"},
		},
		{
			input:    "",
			expected: []string{},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)

		// Check the length of the Actual case vs. the Expected case
		// Fail the test if they don't match
		if len(actual) != len(c.expected) {
			t.Errorf("Length mismatch for %s. Received length: %d, Expected length: %d",
				c.input, len(actual), len(c.expected))
			continue
		}

		// Work through the cleaned "actual" words in the slice
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]

			// If the Actual cleaned word does not match the Expected
			// word, fail the test
			if word != expectedWord {
				t.Errorf("Word mismatch for input '%s' at index '%d'. Received: %s, Expected: %s",
					c.input, i, word, expectedWord)
			}
		}

	}
}
