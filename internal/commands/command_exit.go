package commands

import (
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/noch-g/pokedex-cli/internal/config"
)

func commandExit(cfg *config.Config, output io.Writer, args ...string) error {
	fmt.Fprintf(output, "\rClosing the Pokedex... Goodbye!\n")
	if testing.Testing() {
		return fmt.Errorf("exit called during test")
	}
	cfg.SaveConf()
	cfg.PokeapiClient.SaveCache()
	os.Exit(0)
	return nil
}
