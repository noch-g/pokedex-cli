package commands

import (
	"fmt"
	"io"
	"strings"

	"github.com/noch-g/pokedex-cli/internal/config"
)

func commandPokedex(cfg *config.Config, output io.Writer, args ...string) error {
	if len(cfg.CaughtPokemon) == 0 {
		fmt.Fprintln(output, "Your pokedex is empty for now. Try to use the command catch <pokemon>")
		return nil
	}
	fmt.Fprintf(output, "Your pokedex:\n")
	fmt.Fprintln(output)

	const columns = 5
	const rows = int((151 + columns) / columns)

	table := make([][]string, rows)
	for i := range table {
		table[i] = make([]string, columns)
	}

	var i int
	for _, pokemon := range cfg.CaughtPokemon {
		i = pokemon.ID
		if i > 151 {
			// Skip newer generations for now
			continue
		}

		row := (i - 1) % rows
		col := (i - 1) / rows
		table[row][col] = fmt.Sprintf("#%03d %-15s", i, strings.ToUpper(pokemon.Name[:1])+pokemon.Name[1:])
	}

	for i, row := range table {
		for j, entry := range row {
			if entry != "" {
				fmt.Fprint(output, entry)
			} else if (i+1)+j*rows <= 151 {
				fmt.Fprintf(output, "#%03d %-15s", (i+1)+j*rows, "   ???")
			}
		}
		fmt.Fprintln(output)
	}
	fmt.Fprintln(output)
	return nil
}
