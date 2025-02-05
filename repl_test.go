package main

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/noch-g/pokedex-cli/internal/pokeapi"
)

type commandTestCase struct {
	name     string
	input    string
	expected []string
}

func TestStartRepl_ValidCommands(t *testing.T) {
	pokeClient := pokeapi.NewClient(5*time.Second, 5*time.Minute)

	cfg := &config{
		CaughtPokemon: map[string]pokeapi.Pokemon{},
		knownEntities: map[string][]string{
			"commands":     {},
			"pokemons":     {},
			"locations":    {},
			"wildPokemons": {},
		},
		pokeapiClient: pokeClient,
	}
	commands := []commandTestCase{
		{
			name:     "Check mapb before any location loaded",
			input:    "mapb",
			expected: []string{"you're on the first page"},
		},
		{
			name:     "Load locations with map",
			input:    "map",
			expected: []string{"canalave-city-area"},
		},
		{
			name:     "Load locations with map",
			input:    "map",
			expected: []string{"great-marsh-area-1"},
		},
		{
			name:     "Check mapb after locations loaded",
			input:    "mapb",
			expected: []string{"canalave-city-area", "eterna-city-area"},
		},
	}

	for _, cmd := range commands {
		t.Run(cmd.name, func(t *testing.T) {
			in := bytes.NewBufferString(cmd.input + "\nexit\n")
			out := new(bytes.Buffer)

			startRepl(cfg, in, out)

			result := out.String()
			for _, expected := range cmd.expected {
				if !strings.Contains(result, expected) {
					t.Errorf("\nCommand: %q\nExpected: %s\nGot: %s",
						cmd.input,
						expected,
						result)
				}
			}
		})
	}
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
