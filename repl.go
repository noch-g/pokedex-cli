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
		input, err := readInput(reader, &history, &historyIndex, commands)
		if err != nil {
			// fmt.Println("\nExiting REPL...")
			commands["exit"].callback(cfg)
			break
		}
		// fmt.Println("Input detected: " + input)

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

func readInput(reader *bufio.Reader, history *[]string, historyIndex *int, commands map[string]cliCommand) (string, error) {
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
			suggestions := []string{}

			// Check commands starting with current input
			for cmd := range commands {
				if strings.HasPrefix(cmd, currentInput) {
					suggestions = append(suggestions, cmd)
				}
			}

			if len(suggestions) == 1 {
				input.Reset()
				input.WriteString(suggestions[0] + " ")
				fmt.Print(GetPromptMessage() + suggestions[0] + " ")
				continue
			}

			if len(suggestions) > 1 {
				fmt.Println()
				StartFromClearLine()
				fmt.Println(strings.Join(suggestions, ", "))
				fmt.Print(GetPromptMessage() + currentInput)
				continue
			}
			continue
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
