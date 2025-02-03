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

		// Read user input with terminal settings to detect arrow keys
		input, err := readInput(reader, &history, &historyIndex, cfg.knownEntities)
		if err != nil {
			commands["exit"].callback(cfg)
			break
		}

		words := cleanInput(input)
		if len(words) == 0 {
			continue
		}

		// Store command in history
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

	var input strings.Builder

	for {
		char, err := reader.ReadByte()
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

		// Handle Up and Down Arrow Keys
		if char == 27 {
			next, _ := reader.ReadByte()
			if next == 91 {
				key, _ := reader.ReadByte()
				if key == 65 && len(*history) > 0 { // Up Arrow Pressed
					if *historyIndex > 0 {
						*historyIndex--
					}
				} else if key == 66 && len(*history) > 0 { // Down Arrow  (↓)
					if *historyIndex < len(*history)-1 {
						*historyIndex++
					} else {
						*historyIndex = len(*history) - 1 // Effacer l’entrée courante
					}
				} else {
					continue
				}

				fmt.Print(GetPromptMessage() + (*history)[*historyIndex])
				input.Reset()
				input.WriteString((*history)[*historyIndex])
				continue
			}
		}

		// Handle Backspace (←)
		if char == 127 {
			if input.Len() > 0 {
				str := input.String()
				input.Reset()
				input.WriteString(str[:len(str)-1]) // Supprime le dernier caractère
				fmt.Print("\b \b")                  // Efface visuellement
			}
			continue
		}

		// Handle tab for auto-completion
		if char == 9 {
			currentInput := input.String()
			wordsInput := cleanInput(currentInput)

			if len(wordsInput) == 1 {
				autocomplete("", wordsInput[0], knownEntities["commands"], &input)
				continue
			} else if len(wordsInput) == 2 {
				switch wordsInput[0] {
				case "explore":
					autocomplete("explore", wordsInput[1], knownEntities["locations"], &input)
					continue
				case "inspect":
					autocomplete("inspect", wordsInput[1], knownEntities["pokemons"], &input)
					continue
				case "catch":
					autocomplete("catch", wordsInput[1], knownEntities["wildPokemons"], &input)
					continue
				default:
					continue
				}
			} else {
				continue
			}
		}

		// Append character to input
		input.WriteByte(char)
		fmt.Print(string(char))
	}

	return input.String(), nil
}

func cleanInput(text string) []string {
	output := strings.ToLower(text)
	words := strings.Fields(output)
	return words
}

func autocomplete(cmd string, strStart string, wordsDict []string, input *strings.Builder) {
	suggestions := []string{}
	for _, entity := range wordsDict {
		if strings.HasPrefix(entity, strStart) {
			suggestions = append(suggestions, entity)
		}
	}

	var newInput string
	if len(suggestions) == 1 {
		input.Reset()
		if len(cmd) > 0 {
			newInput += cmd + " "
		}
		newInput += suggestions[0] + " "
		input.WriteString(newInput)
		fmt.Print(GetPromptMessage() + newInput)
	} else if len(suggestions) > 1 {
		fmt.Println()
		StartFromClearLine()
		for _, suggestion := range suggestions {
			StartFromClearLine()
			fmt.Println(suggestion)
		}
		prefix := LongestCommonPrefix(suggestions)

		input.Reset()
		if len(cmd) > 0 {
			newInput += cmd + " "
		}
		newInput += prefix
		input.WriteString(newInput)
		fmt.Print(GetPromptMessage() + newInput)
	}
}
