package main

import (
	"errors"
	"fmt"
	"math/rand"
)

func commandCatch(cfg *config, args ...string) error {
	if len(args) != 1 {
		return errors.New("you must provide a pokemon name")
	}
	pokemonName := args[0]

	pokemon, err := cfg.pokeapiClient.GetPokemon(pokemonName)
	if err != nil {
		return err
	}

	res := rand.Intn(pokemon.BaseExperience)
	successChance := 100 * 40.0 / float64(pokemon.BaseExperience)

	fmt.Printf("Throwing a Pokeball at %s... (success chance: %.2f%%)\n", pokemon.Name, successChance)
	if res > 40 {
		fmt.Printf("%s escaped!\n", pokemon.Name)
		return nil
	}
	fmt.Printf("%s was caught! (#%03d)\n", pokemon.Name, pokemon.ID)
	imgStr, err := cfg.pokeapiClient.RenderImage(&pokemon)
	if err != nil {
		fmt.Println("(Image could not be retrieved)")
	} else {
		fmt.Println(imgStr)
	}
	fmt.Printf("The information was added to the pokedex (#%03d). You may now inspect it with the inspect command.\n", pokemon.ID)

	cfg.CaughtPokemon[pokemon.Name] = pokemon

	return nil
}
