package commands

import (
	"errors"
	"fmt"
	"io"

	"github.com/noch-g/pokedex-cli/internal/config"
)

func commandInspect(cfg *config.Config, output io.Writer, args ...string) error {
	if len(args) != 1 {
		return errors.New("you must provide a pokemon name")
	}
	pokemonName := args[0]

	pokemon, ok := cfg.CaughtPokemon[pokemonName]
	if !ok {
		fmt.Fprintf(output, "You have not caught a %v yet\n", pokemonName)
		return nil
	}
	imgStr, err := cfg.PokeapiClient.RenderImage(&pokemon)
	if err != nil {
		fmt.Fprintln(output, "(Image could not be retrieved)")
	} else {
		fmt.Fprintln(output, imgStr)
	}
	fmt.Fprintf(output, "Name: %s\n", pokemon.Name)
	fmt.Fprintf(output, "Height: %v\n", pokemon.Height)
	fmt.Fprintf(output, "Weight: %v\n", pokemon.Weight)
	fmt.Fprintf(output, "Stats:")
	for _, stat := range pokemon.Stats {
		fmt.Fprintf(output, "  -%s: %v\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Fprintln(output, "Types:")
	for _, typeInfo := range pokemon.Types {
		fmt.Fprintf(output, "  - %s\n", typeInfo.Type.Name)
	}

	return nil
}
