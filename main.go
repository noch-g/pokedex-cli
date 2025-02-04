package main

import (
	"os"
	"time"

	"github.com/noch-g/pokedex-cli/internal/pokeapi"
)

func main() {
	pokeClient := pokeapi.NewClient(5*time.Second, 5*time.Minute)
	cfg := &config{
		CaughtPokemon: map[string]pokeapi.Pokemon{},
		knownEntities: map[string][]string{
			"commands":     {},
			"pokemons":     {},
			"locations":    {},
			"wildPokemons": {},
		},
		pokeapiClient: pokeClient,
	}
	cfg.Load("pokemons.json")
	startRepl(cfg, os.Stdin, os.Stdout)
}
