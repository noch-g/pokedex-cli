package commands

import (
	"io"

	"github.com/noch-g/pokedex-cli/internal/config"
)

type cliCommand struct {
	name        string
	description string
	Callback    func(*config.Config, io.Writer, ...string) error
}

func GetCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			Callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Get the next page of locations",
			Callback:    commandMapf,
		},
		"mapb": {
			name:        "mapb",
			description: "Get the previous page of locations",
			Callback:    commandMapb,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			Callback:    commandExit,
		},
		"explore": {
			name:        "explore <location_name>",
			description: "Explore a location",
			Callback:    commandExplore,
		},
		"catch": {
			name:        "catch <pokemon_name>",
			description: "Attempt to catch a pokemon",
			Callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect <pokemon_name>",
			description: "Inspect a pokemon that you have caught",
			Callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "show the pokemons you have caught",
			Callback:    commandPokedex,
		},
	}
}

func GetCommandNames() []string {
	commands := GetCommands()
	commandNames := make([]string, 0, len(commands))
	for name := range commands {
		commandNames = append(commandNames, name)
	}
	return commandNames
}
