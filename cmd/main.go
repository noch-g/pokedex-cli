package main

import (
	"os"
	"time"

	"github.com/noch-g/pokedex-cli/internal/commands"
	"github.com/noch-g/pokedex-cli/internal/config"
	"github.com/noch-g/pokedex-cli/internal/pokeapi"
)

func main() {
	pokeClient := pokeapi.NewClient(5*time.Second, 5*time.Minute)
	cfg := config.NewConfig(pokeClient)
	cfg.Load("pokemons.json", commands.GetCommandNames())

	startRepl(cfg, os.Stdin, os.Stdout)
}
