package main

import (
	"fmt"
	"strings"
)

func commandPokedex(cfg *config, args ...string) error {
	if len(cfg.CaughtPokemon) == 0 {
		fmt.Println("Your pokedex is empty for now. Try to use the command catch <pokemon>")
		return nil
	}
	fmt.Printf("Your pokedex:\n")
	fmt.Println()

	const columns = 5
	const rows = int((151 + columns) / columns)

	// Create a table with correct dimensions
	table := make([][]string, rows)
	for i := range table {
		table[i] = make([]string, columns)
	}

	var i int
	for _, pokemon := range cfg.CaughtPokemon {
		i = pokemon.ID
		if i > 151 {
			// Skip for now
			continue
		}

		row := (i - 1) % rows
		col := (i - 1) / rows
		table[row][col] = fmt.Sprintf("#%03d %-15s", i, strings.ToUpper(pokemon.Name[:1])+pokemon.Name[1:])
	}

	// Print row by row
	for i, row := range table {
		for j, entry := range row {
			if entry != "" {
				fmt.Print(entry)
			} else if (i+1)+j*rows <= 151 {
				fmt.Printf("#%03d %-15s", (i+1)+j*rows, "   ???")
			}
		}
		fmt.Println()
	}
	fmt.Println()
	return nil
}
