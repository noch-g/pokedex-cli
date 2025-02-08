package repl

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/noch-g/pokedex-cli/internal/config"
	"github.com/noch-g/pokedex-cli/internal/pokeapi"
)

type commandTestCase struct {
	name     string
	input    string
	expected []string
}

func TestStartRepl_ValidCommands(t *testing.T) {
	cfg := setupTestConfig()
	commands := []commandTestCase{
		{
			name:     "Check mapb before any location loaded",
			input:    "mapb",
			expected: []string{config.GetPromptMessage() + "mapb", "you're on the first page"},
		},
		{
			name:     "Load locations with map",
			input:    "map",
			expected: []string{config.GetPromptMessage() + "map", "canalave-city-area"},
		},
		{
			name:     "Load locations with map x2",
			input:    "map",
			expected: []string{config.GetPromptMessage() + "map", "mt-coronet-1f-route-216"},
		},
		{
			name:     "Check mapb after locations loaded",
			input:    "mapb",
			expected: []string{config.GetPromptMessage() + "mapb", "canalave-city-area", "eterna-city-area"},
		},
	}

	inR, inW := io.Pipe()
	outR, outW := io.Pipe()

	go StartRepl(cfg, inR, outW)

	// Make this higher if cmd is expected to output more lines
	outputChan := make(chan string, 100)
	go func() {
		reader := bufio.NewReader(outR)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				close(outputChan)
				return
			}
			outputChan <- line
		}
	}()

	for _, cmd := range commands {
		t.Run(cmd.name, func(t *testing.T) {
			fmt.Fprintf(inW, "%s\n", cmd.input)
			var output []string
			timeout := time.After(5 * time.Second)
		collectNbExpectedLines:
			for len(output) < len(cmd.expected) {
				select {
				case line, ok := <-outputChan:
					if !ok {
						break
					}
					output = append(output, line)
				case <-timeout:
					t.Errorf("Timeout waiting for command %s output. Got %d lines, expected %d",
						cmd.input, len(output), len(cmd.expected))
					break collectNbExpectedLines
				}
			}

			if !areAllExpectedFound(output, cmd) {
				var errorStr strings.Builder
				errorStr.WriteString(fmt.Sprintf("Cmd %s failed,\n", cmd.input))
				errorStr.WriteString(fmt.Sprintf("Expected:\n%s\n", strings.Join(cmd.expected, "\n")))
				errorStr.WriteString(fmt.Sprintf("Received:\n%s\n", strings.Join(output, "")))
				t.Error(errorStr.String())
			}

			flushOutputChan(outputChan)
		})
	}

	// Clean up
	fmt.Fprintf(inW, "exit\n")
	inW.Close()
	outW.Close()
}

func flushOutputChan(outputChan <-chan string) {
	for {
		select {
		case <-outputChan:
		default:
			// If there is no data in the channel for 500ms, consider it flushed
			time.Sleep(500 * time.Millisecond)
			select {
			case <-outputChan:
			default:
				return
			}
		}
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

func setupTestConfig() *config.Config {
	pokeClient := pokeapi.NewClient(5*time.Second, 5*time.Minute)
	return config.NewConfig(pokeClient)
}

func areAllExpectedFound(output []string, cmd commandTestCase) bool {
	for _, expected := range cmd.expected {
		found := false
		for _, out := range output {
			if strings.Contains(out, expected) {
				found = true
				break
			}
		}
		if !found {
			fmt.Println("Could not find expected output:", expected)
			return false
		}
	}
	return true
}
