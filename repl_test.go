package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "hello world",
			expected: []string{"hello", "world"},
		},
		{
			input:    "YELLING at YOU",
			expected: []string{"yelling", "at", "you"},
		},
		{
			input:    "  lots    of white	 space     ",
			expected: []string{"lots", "of", "white", "space"},
		},
		{
			input:    "single",
			expected: []string{"single"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) == len(c.expected) {
			for i := range actual {
				word := actual[i]
				expectedWord := c.expected[i]
				if word != expectedWord {
					t.Errorf("test failed; %v != %v", actual, c.expected)
				}
			}
		} else {
			t.Errorf("test failed; %v != %v", actual, c.expected)
		}
	}
}
