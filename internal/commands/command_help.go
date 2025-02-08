package commands

import (
	"fmt"
	"io"

	"github.com/noch-g/pokedex-cli/internal/config"
)

func commandHelp(cfg *config.Config, output io.Writer, args ...string) error {
	fmt.Fprintln(output)
	fmt.Fprintln(output, "Welcome to the Pokedex!")
	fmt.Fprintln(output, "Usage:")
	fmt.Fprintln(output)
	for _, cmd := range GetCommands() {
		fmt.Fprintf(output, "%s: %s\n", cmd.name, cmd.description)
	}
	fmt.Fprintln(output)
	return nil
}