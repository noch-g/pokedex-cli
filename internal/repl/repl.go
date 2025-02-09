package repl

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"syscall"
	"testing"

	"github.com/noch-g/pokedex-cli/internal/commands"
	"github.com/noch-g/pokedex-cli/internal/config"
	"github.com/noch-g/pokedex-cli/internal/text"
	"golang.org/x/term"
)

func StartRepl(cfg *config.Config, input io.Reader, output io.Writer) {
	reader := bufio.NewReader(input)
	commands := commands.GetCommands()
	var history []string
	historyIndex := -1

	for {
		fmt.Fprint(output, config.GetPromptMessage())

		userInput, err := readInput(reader, &history, &historyIndex, cfg.KnownEntities, output)
		if err != nil {
			if err.Error() == "ctrl+C or ctrl+D called" {
				commands["exit"].Callback(cfg, output)
			}
			break
		}

		words := cleanInput(userInput)
		if len(words) == 0 {
			continue
		}

		history = append(history, userInput)
		historyIndex = len(history)

		args := []string{}
		if len(words) > 1 {
			args = words[1:]
		}
		commandName := words[0]

		command, exists := commands[commandName]
		text.StartFromClearLine(output)
		if exists {
			err := command.Callback(cfg, output, args...)
			if err != nil {
				fmt.Fprintln(output, err)
			}
		} else {
			fmt.Fprintf(output, "Unknown command: \"%s\". Type 'help' for a list of available commands.\n", commandName)
		}
	}
}

func readInput(reader *bufio.Reader, history *[]string, historyIndex *int, knownEntities map[string][]string, output io.Writer) (string, error) {
	if !testing.Testing() {
		oldState, err := term.MakeRaw(int(syscall.Stdin))
		if err != nil {
			return "", err
		}
		defer term.Restore(int(syscall.Stdin), oldState)
	}

	var inputSlice []rune
	cursorPos := 0
	for {
		char, _, err := reader.ReadRune()
		if err != nil {
			return "", err
		}

		// Handle Ctrl+C and Ctrl+D
		if char == 3 || char == 4 {
			fmt.Fprintf(output, "\n")
			return "", fmt.Errorf("ctrl+C or ctrl+D called")
		}

		// Handle Enter key
		if char == 10 || char == 13 {
			fmt.Fprintf(output, "\n")
			break
		}

		// Handle Arrow Keys
		if char == 27 {
			next, _ := reader.ReadByte()
			if next == 91 {
				key, _ := reader.ReadByte()
				if key == 65 && len(*history) > 0 { // Up Arrow (↑)
					if *historyIndex > 0 {
						*historyIndex--
					}
					inputSlice = []rune((*history)[*historyIndex])
					cursorPos = len(inputSlice)
					redrawLine(inputSlice, cursorPos, output)
					continue
				} else if key == 66 && len(*history) > 0 { // Down Arrow  (↓)
					if *historyIndex < len(*history)-1 {
						*historyIndex++
					} else {
						*historyIndex = len(*history) - 1
					}
					inputSlice = []rune((*history)[*historyIndex])
					cursorPos = len(inputSlice)
					redrawLine(inputSlice, cursorPos, output)
					continue
				} else if key == 67 { // Right Arrow (→)
					if cursorPos < len(inputSlice) {
						cursorPos++
						fmt.Fprintf(output, "\x1b[1C") // Move cursor to the right
					}
					continue
				} else if key == 68 { // Left Arrow (←)
					if cursorPos > 0 {
						cursorPos--
						fmt.Fprintf(output, "\x1b[1D") // Move cursor to the left
					}
					continue
				}
			}
		}

		// Handle Backspace (←)
		if char == 127 {
			if cursorPos > 0 {
				inputSlice = append(inputSlice[:cursorPos-1], inputSlice[cursorPos:]...)
				cursorPos--
				redrawLine(inputSlice, cursorPos, output)
			}
			continue
		}

		// Handle tab for auto-completion
		if char == 9 {
			currentInput := string(inputSlice)
			wordsInput := cleanInput(currentInput)

			if len(wordsInput) == 1 {
				autocomplete("", wordsInput[0], knownEntities["commands"], &inputSlice, &cursorPos, output)
				continue
			} else if len(wordsInput) == 2 {
				switch wordsInput[0] {
				case "explore":
					autocomplete("explore", wordsInput[1], knownEntities["locations"], &inputSlice, &cursorPos, output)
					continue
				case "inspect":
					autocomplete("inspect", wordsInput[1], knownEntities["pokemons"], &inputSlice, &cursorPos, output)
					continue
				case "catch":
					autocomplete("catch", wordsInput[1], knownEntities["wildPokemons"], &inputSlice, &cursorPos, output)
					continue
				default:
					continue
				}
			} else {
				continue
			}
		}
		inputSlice = append(inputSlice[:cursorPos], append([]rune{char}, inputSlice[cursorPos:]...)...)
		cursorPos++
		redrawLine(inputSlice, cursorPos, output)
	}

	return string(inputSlice), nil
}

func redrawLine(inputSlice []rune, cursorPos int, output io.Writer) {
	fmt.Fprint(output, "\r"+config.GetPromptMessage()+string(inputSlice)+" \x1b[K")
	placeCursor(cursorPos, output)
}

func placeCursor(cursorPos int, output io.Writer) {
	fmt.Fprintf(output, "\r\x1b[%dC", config.GetPromptLength()+cursorPos)
}

func cleanInput(text string) []string {
	output := strings.ToLower(text)
	words := strings.Fields(output)
	return words
}

func autocomplete(cmd string, strStart string, wordsDict []string, inputSlice *[]rune, cursorPos *int, output io.Writer) {
	suggestions := []string{}
	for _, entity := range wordsDict {
		if strings.HasPrefix(entity, strStart) {
			suggestions = append(suggestions, entity)
		}
	}
	if len(suggestions) == 0 {
		return
	}

	var newInput string

	if len(cmd) > 0 {
		newInput += cmd + " "
	}

	if len(suggestions) == 1 {
		newInput += suggestions[0] + " "
	} else if len(suggestions) > 1 {
		fmt.Println()
		for _, suggestion := range suggestions {
			text.StartFromClearLine(output)
			fmt.Println(suggestion)
		}
		newInput += text.LongestCommonPrefix(suggestions)
	}
	*inputSlice = []rune(newInput)
	*cursorPos = len(newInput)
	redrawLine(*inputSlice, *cursorPos, output)
}
