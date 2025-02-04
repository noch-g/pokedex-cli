package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"golang.org/x/term"
)

func startRepl(cfg *config) {
	reader := bufio.NewReader(os.Stdin)
	commands := getCommands()
	var history []string
	historyIndex := -1

	for {
		fmt.Print(GetPromptMessage())

		input, err := readInput(reader, &history, &historyIndex, cfg.knownEntities)
		if err != nil {
			commands["exit"].callback(cfg)
			break
		}

		words := cleanInput(input)
		if len(words) == 0 {
			continue
		}

		history = append(history, input)
		historyIndex = len(history)

		args := []string{}
		if len(words) > 1 {
			args = words[1:]
		}
		commandName := words[0]

		command, exists := commands[commandName]
		StartFromClearLine()
		if exists {
			err := command.callback(cfg, args...)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Printf("Unknown command: \"%s\"\n", commandName)
		}
	}
}

func readInput(reader *bufio.Reader, history *[]string, historyIndex *int, knownEntities map[string][]string) (string, error) {
	// Switch terminal to raw mode
	oldState, err := term.MakeRaw(int(syscall.Stdin))
	if err != nil {
		return "", err
	}
	defer term.Restore(int(syscall.Stdin), oldState)

	var input []rune
	cursorPos := 0
	for {
		char, _, err := reader.ReadRune()
		if err != nil {
			return "", err
		}

		// Handle Ctrl+C and Ctrl+D
		if char == 3 || char == 4 {
			fmt.Printf("\n")
			StartFromClearLine()
			return "", fmt.Errorf("Exiting")
		}

		// Handle Enter key
		if char == 10 || char == 13 {
			fmt.Printf("\n")
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
					input = []rune((*history)[*historyIndex])
					cursorPos = len(input)
					redrawLine(input, cursorPos)
					continue
				} else if key == 66 && len(*history) > 0 { // Down Arrow  (↓)
					if *historyIndex < len(*history)-1 {
						*historyIndex++
					} else {
						*historyIndex = len(*history) - 1
					}
					input = []rune((*history)[*historyIndex])
					cursorPos = len(input)
					redrawLine(input, cursorPos)
					continue
				} else if key == 67 { // Right Arrow (→)
					if cursorPos < len(input) {
						cursorPos++
						fmt.Print("\x1b[1C") // Move cursor to the right
					}
					continue
				} else if key == 68 { // Left Arrow (←)
					if cursorPos > 0 {
						cursorPos--
						fmt.Print("\x1b[1D") // Move cursor to the left
					}
					continue
				}
			}
		}

		// Handle Backspace (←)
		if char == 127 && cursorPos > 0 {
			input = append(input[:cursorPos-1], input[cursorPos:]...)
			cursorPos--
			redrawLine(input, cursorPos)
			continue
		}

		// Handle tab for auto-completion
		if char == 9 {
			currentInput := string(input)
			wordsInput := cleanInput(currentInput)

			if len(wordsInput) == 1 {
				autocomplete("", wordsInput[0], knownEntities["commands"], &input, &cursorPos)
				continue
			} else if len(wordsInput) == 2 {
				switch wordsInput[0] {
				case "explore":
					autocomplete("explore", wordsInput[1], knownEntities["locations"], &input, &cursorPos)
					continue
				case "inspect":
					autocomplete("inspect", wordsInput[1], knownEntities["pokemons"], &input, &cursorPos)
					continue
				case "catch":
					autocomplete("catch", wordsInput[1], knownEntities["wildPokemons"], &input, &cursorPos)
					continue
				default:
					continue
				}
			} else {
				continue
			}
		}
		input = append(input[:cursorPos], append([]rune{char}, input[cursorPos:]...)...)
		cursorPos++
		redrawLine(input, cursorPos)
	}

	return string(input), nil
}

func redrawLine(input []rune, cursorPos int) {
	fmt.Print("\r" + GetPromptMessage() + string(input) + " \x1b[K")
	placeCursor(cursorPos)
}

func placeCursor(cursorPos int) {
	fmt.Printf("\r\x1b[%dC", GetPromptLength()+cursorPos)
}

func cleanInput(text string) []string {
	output := strings.ToLower(text)
	words := strings.Fields(output)
	return words
}

func autocomplete(cmd string, strStart string, wordsDict []string, input *[]rune, cursorPos *int) {
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
			StartFromClearLine()
			fmt.Println(suggestion)
		}
		newInput += LongestCommonPrefix(suggestions)
	}
	*input = []rune(newInput)
	*cursorPos = len(newInput)
	redrawLine(*input, *cursorPos)
}
