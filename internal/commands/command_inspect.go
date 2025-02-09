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
	fmt.Fprintln(output, pokemon.GetStatsStr())
	return nil
}
