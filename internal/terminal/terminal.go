package terminal

import (
	"bufio"
	"fmt"
	"io"
	"slices"
	"strings"
	"syscall"
	"testing"

	"github.com/noch-g/pokedex-cli/internal/config"
	"github.com/noch-g/pokedex-cli/internal/text"
	"golang.org/x/term"
)

func ReadInput(reader *bufio.Reader, history *[]string, historyIndex *int, knownEntities map[string][]string, output io.Writer) (string, error) {
	if !testing.Testing() {
		oldState, err := term.MakeRaw(int(syscall.Stdin))
		if err != nil {
			return "", err
		}
		defer term.Restore(int(syscall.Stdin), oldState)
	}

	var inputSlice []rune
	var isSecondTab = false
	var doubleTapDetected = false
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
			text.StartFromClearLine(output)
			break
		}

		// Handle tab for auto-completion
		doubleTapDetected = false
		if char == 9 {
			currentInput := string(inputSlice)
			wordsInput := text.CleanInput(currentInput)

			if isSecondTab {
				isSecondTab = false
				doubleTapDetected = true
			} else {
				isSecondTab = true
			}

			if len(wordsInput) == 0 && doubleTapDetected {
				autocomplete("", "", knownEntities["commands"], &inputSlice, &cursorPos, output)
				continue
			} else if len(wordsInput) == 1 {
				switch wordsInput[0] {
				case "inspect":
					if doubleTapDetected {
						autocomplete("inspect", "", knownEntities["pokemons"], &inputSlice, &cursorPos, output)
					}
					continue
				case "explore":
					if doubleTapDetected {
						autocomplete("explore", "", knownEntities["locations"], &inputSlice, &cursorPos, output)
					}
					continue
				case "catch":
					if doubleTapDetected {
						autocomplete("catch", "", knownEntities["wildPokemons"], &inputSlice, &cursorPos, output)
					}
					continue
				case "compare":
					if doubleTapDetected {
						autocomplete("compare", "", knownEntities["pokemons"], &inputSlice, &cursorPos, output)
					}
					continue
				default:
					autocomplete("", wordsInput[0], knownEntities["commands"], &inputSlice, &cursorPos, output)
					continue
				}
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
				case "compare":
					if doubleTapDetected && slices.Contains(knownEntities["pokemons"], wordsInput[1]) {
						autocomplete("compare "+wordsInput[1], "", knownEntities["pokemons"], &inputSlice, &cursorPos, output)
					} else {
						autocomplete("compare", wordsInput[1], knownEntities["pokemons"], &inputSlice, &cursorPos, output)
					}
					continue
				default:
					continue
				}
			} else if len(wordsInput) == 3 {
				switch wordsInput[0] {
				case "compare":
					autocomplete("compare "+wordsInput[1], wordsInput[2], knownEntities["pokemons"], &inputSlice, &cursorPos, output)
					continue
				default:
					continue
				}
			} else {
				continue
			}
		} else {
			isSecondTab = false
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
		columnWidth := max(15, text.LongestWordLength(suggestions)+5)
		columns := 4
		var returnNeeded bool = false

		fmt.Fprintln(output)
		text.StartFromClearLine(output)
		for i, suggestion := range suggestions {
			returnNeeded = true
			fmt.Fprintf(output, "%-*s", columnWidth, suggestion) // Left-align with fixed width
			if (i+1)%columns == 0 {                              // Move to new line after fixed number of columns
				returnNeeded = false
				fmt.Fprint(output, "\n\r")
			}
		}
		if returnNeeded {
			fmt.Fprintf(output, "\n")
		}
		newInput += text.LongestCommonPrefix(suggestions)
	}
	*inputSlice = []rune(newInput)
	*cursorPos = len(newInput)
	redrawLine(*inputSlice, *cursorPos, output)
}
