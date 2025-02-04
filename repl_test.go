package main

import (
	"bytes"
	"fmt"
	"testing"
)

func TestStartRepl_ValidCommand(t *testing.T) {
	cfg := &config{}
	input := "help\nexit\n" // Simule l'entrée de l'utilisateur
	in := bytes.NewBufferString(input)
	out := new(bytes.Buffer)

	startRepl(cfg, in, out)

	result := out.String()
	fmt.Println(result)
}

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    " Charmander Bulbasaur PIKACHU ",
			expected: []string{"charmander", "bulbasaur", "pikachu"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("the number of elements did not match, actual: %s, expected: %s", actual, c.expected)
		}

		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("words did not match, actual: %s, expected: %s", word, expectedWord)
			}
		}
	}
}
