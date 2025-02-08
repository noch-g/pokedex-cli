package commands

import (
	"errors"
	"fmt"
	"io"

	"github.com/noch-g/pokedex-cli/internal/config"
)

func commandExplore(cfg *config.Config, output io.Writer, args ...string) error {
	if len(args) != 1 {
		return errors.New("you must provide a location name")
	}
	locationStr := args[0]

	loc, err := cfg.PokeapiClient.GetLocation(locationStr)
	if err != nil {
		return err
	}
	fmt.Fprintln(output, "Exploring "+locationStr+"...")
	cfg.KnownEntities["wildPokemons"] = []string{}
	if len(loc.PokemonEncounters) == 0 {
		fmt.Fprintln(output, "No pokemon present in this area")
		return nil
	} else {
		fmt.Fprintln(output, "Found Pokemon:")
		for _, pokemon_encounter := range loc.PokemonEncounters {
			isNew := ""
			if _, ok := cfg.CaughtPokemon[pokemon_encounter.Pokemon.Name]; !ok {
				isNew = "(New!)"
			}
			pokemon, err := cfg.PokeapiClient.GetPokemon(pokemon_encounter.Pokemon.Name)
			if err != nil {
				fmt.Fprintf(output, " - %-10s\n", pokemon_encounter.Pokemon.Name)
			} else {
				fmt.Fprintf(output, " - %-10s #%03d %s\n", pokemon.Name, pokemon.ID, isNew)
			}
			cfg.KnownEntities["wildPokemons"] = append(cfg.KnownEntities["wildPokemons"], pokemon_encounter.Pokemon.Name)
		}
	}

	return nil
}
