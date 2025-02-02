package main

import (
	"fmt"
)

func commandPokedex(cfg *config, args ...string) error {
	if len(cfg.CaughtPokemon) == 0 {
		fmt.Println("Your pokedex is empty for now. Try to use the command catch <pokemon>")
		return nil
	}
	fmt.Printf("Your pokedex:\n")
	for _, pokemon := range cfg.CaughtPokemon {
		fmt.Printf(" - %s\n", pokemon.Name)
	}
	return nil
}
