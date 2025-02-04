package main

import (
	"fmt"
	"io"
	"os"
	"testing"
)

func commandExit(cfg *config, output io.Writer, args ...string) error {
	fmt.Fprintf(output, "Closing the Pokedex... Goodbye!\n")
	if testing.Testing() {
		return fmt.Errorf("exit called during test")
	}
	cfg.Save("pokemons.json")
	os.Exit(0)
	return nil
}
