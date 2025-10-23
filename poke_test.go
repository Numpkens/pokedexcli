package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		// 1. User's primary example (trimming and splitting)
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		// 2. User's secondary example (casing)
		{
			input:    "Charmander Bulbasaur PIKACHU",
			expected: []string{"charmander", "bulbasaur", "pikachu"},
		},
		// 3. Single word test
		{
			input:    " SingleWord ",
			expected: []string{"singleword"},
		},
		// 4. Input with tabs and newlines (should still split/clean)
		{
			input:    "First\t\nsecond\t\nTHIRD",
			expected: []string{"first", "second", "third"},
		},
		// 5. Empty string
		{
			input:    "",
			expected: []string{},
		},
		// 6. Only whitespace
		{
			input:    " \t \n ",
			expected: []string{},
		},
	}

	// Loop over all defined test cases
	for _, c := range cases {
		actual := cleanInput(c.input)

		// 1. Check the length of the actual slice against the expected slice
		if len(actual) != len(c.expected) {
			t.Errorf("Input: %q. Expected length %d (%v), but got %d (%v)",
				c.input, len(c.expected), c.expected, len(actual), actual)
			continue // Stop and move to the next case if lengths don't match
		}

		// 2. Check each word in the slice for correctness
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]

			if word != expectedWord {
				t.Errorf("Input: %q. At index %d, expected word %q but got %q",
					c.input, i, expectedWord, word)
			}
		}
	}
}
