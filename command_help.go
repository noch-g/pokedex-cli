package main

import (
	"fmt"
	"io"
)

func commandHelp(cfg *config, output io.Writer, args ...string) error {
	fmt.Fprintln(output)
	fmt.Fprintln(output, "Welcome to the Pokedex!")
	fmt.Fprintln(output, "Usage:")
	fmt.Fprintln(output)
	for _, cmd := range getCommands() {
		fmt.Fprintf(output, "%s: %s\n", cmd.name, cmd.description)
	}
	fmt.Fprintln(output)
	return nil
}
