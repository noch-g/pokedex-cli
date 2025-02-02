package main

import (
	"time"

	"github.com/noch-g/pokedex-cli/internal/pokeapi"
)

func main() {
	pokeClient := pokeapi.NewClient(5*time.Second, 5*time.Minute)
	cfg := &config{
		CaughtPokemon: map[string]pokeapi.Pokemon{},
		pokeapiClient: pokeClient,
	}
	cfg.Load("pokemons.json")
	startRepl(cfg)
}
