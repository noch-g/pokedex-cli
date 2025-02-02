package main

import (
	"errors"
	"fmt"
)

func commandInspect(cfg *config, args ...string) error {
	if len(args) != 1 {
		return errors.New("you must provide a pokemon name")
	}
	pokemonName := args[0]

	pokemon, ok := cfg.CaughtPokemon[pokemonName]
	if !ok {
		fmt.Printf("You have not caught a %v yet\n", pokemonName)
		return nil
	}
	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %v\n", pokemon.Height)
	fmt.Printf("Weight: %v\n", pokemon.Weight)
	fmt.Println("Stats:")
	for _, stat := range pokemon.Stats {
		fmt.Printf("  -%s: %v\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Println("Types:")
	for _, typeInfo := range pokemon.Types {
		fmt.Printf("  - %s\n", typeInfo.Type.Name)
	}

	return nil
}
