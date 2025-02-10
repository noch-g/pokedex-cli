package main

import (
	"os"
	"time"

	"github.com/noch-g/pokedex-cli/internal/commands"
	"github.com/noch-g/pokedex-cli/internal/config"
	"github.com/noch-g/pokedex-cli/internal/logger"
	"github.com/noch-g/pokedex-cli/internal/pokeapi"
	"github.com/noch-g/pokedex-cli/internal/repl"
	"github.com/noch-g/pokedex-cli/internal/settings"
)

func main() {
	logger.InitLogger()
	pokeClient := pokeapi.NewClient(5*time.Second, 24*time.Hour)
	cfg := config.NewConfig(pokeClient)
	cfg.Load(settings.SaveFilePath, commands.GetCommandNames())

	repl.StartRepl(cfg, os.Stdin, os.Stdout)
}
