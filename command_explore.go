package main

import (
	"errors"
	"fmt"
)

func commandExplore(cfg *config, args ...string) error {
	if len(args) != 1 {
		return errors.New("you must provide a location name")
	}
	locationStr := args[0]

	loc, err := cfg.pokeapiClient.GetLocation(locationStr)
	if err != nil {
		return err
	}
	fmt.Println("Exploring " + locationStr + "...")
	if len(loc.PokemonEncounters) == 0 {
		fmt.Println("No pokemon present in this area")
		return nil
	} else {
		fmt.Println("Found Pokemon:")
		for _, pokemon_encounter := range loc.PokemonEncounters {
			isNew := ""
			if _, ok := cfg.CaughtPokemon[pokemon_encounter.Pokemon.Name]; !ok {
				isNew = "(New!)"
			}
			pokemon, err := cfg.pokeapiClient.GetPokemon(pokemon_encounter.Pokemon.Name)
			if err != nil {
				fmt.Printf(" - %-10s", pokemon_encounter.Pokemon.Name)
			} else {
				fmt.Printf(" - %-10s #%03d %s\n", pokemon.Name, pokemon.ID, isNew)
			}
		}
	}

	return nil
}
