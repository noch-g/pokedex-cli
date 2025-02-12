package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/noch-g/pokedex-cli/internal/commands"
	"github.com/noch-g/pokedex-cli/internal/config"
	"github.com/noch-g/pokedex-cli/internal/terminal"
	"github.com/noch-g/pokedex-cli/internal/text"
)

func StartRepl(cfg *config.Config, input io.Reader, output io.Writer) {
	reader := bufio.NewReader(input)
	commands := commands.GetCommands()
	var history []string
	historyIndex := -1

	for {
		fmt.Fprint(output, config.GetPromptMessage())

		userInput, err := terminal.ReadInput(reader, &history, &historyIndex, cfg.KnownEntities, output)
		if err != nil {
			if err.Error() == "ctrl+C or ctrl+D called" {
				commands["exit"].Callback(cfg, output)
			}
			break
		}

		words := text.CleanInput(userInput)
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
