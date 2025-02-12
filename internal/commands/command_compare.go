package commands

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/noch-g/pokedex-cli/internal/config"
	"golang.org/x/term"
)

func commandCompare(cfg *config.Config, output io.Writer, args ...string) error {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return fmt.Errorf("error getting terminal size: %v", err)
	}

	if len(args) != 2 {
		return errors.New("you must provide two pokemon names")
	}
	pokemonName1 := args[0]
	pokemonName2 := args[1]

	pokemon1, ok := cfg.CaughtPokemon[pokemonName1]
	if !ok {
		return fmt.Errorf("you have not caught a %v yet. compare only works with caught pokemons", pokemonName1)
	}
	pokemon2, ok := cfg.CaughtPokemon[pokemonName2]
	if !ok {
		return fmt.Errorf("you have not caught a %v yet. compare only works with caught pokemons", pokemonName2)
	}

	img1, err := cfg.PokeapiClient.RenderImage(&pokemon1)
	if err != nil {
		return fmt.Errorf("error rendering image for %v: %v", pokemonName1, err)
	}
	img2, err := cfg.PokeapiClient.RenderImage(&pokemon2)
	if err != nil {
		return fmt.Errorf("error rendering image for %v: %v", pokemonName2, err)
	}

	comparedImageStr, err := makeImagesSideBySide(img1, img2, width)
	if err != nil {
		return fmt.Errorf("error making images side by side: %v", err)
	}

	fmt.Fprintf(output, "%v\n", comparedImageStr)

	comparedStatsStr := makeStatsSideBySide(pokemon1.GetStatsStr(), pokemon2.GetStatsStr(), width)
	fmt.Fprintf(output, "%v\n", comparedStatsStr)

	return nil
}

func makeStatsSideBySide(stats1, stats2 string, tWidth int) string {
	stats1Lines := strings.Split(stats1, "\n")
	stats2Lines := strings.Split(stats2, "\n")

	stats2LongestLine := 0
	for _, line := range stats2Lines {
		if len(line) > stats2LongestLine {
			stats2LongestLine = len(line)
		}
	}

	combinedLines := make([]string, max(len(stats1Lines), len(stats2Lines)))
	for i := range combinedLines {
		if i < len(stats1Lines) && i < len(stats2Lines) {
			combinedLines[i] = stats1Lines[i] + strings.Repeat(" ", tWidth-stats2LongestLine-len(stats1Lines[i])) + stats2Lines[i]
		} else if i < len(stats1Lines) {
			combinedLines[i] = stats1Lines[i]
		} else {
			combinedLines[i] = strings.Repeat(" ", tWidth-stats2LongestLine) + stats2Lines[i]
		}
	}
	return strings.Join(combinedLines, "\n")
}

func makeImagesSideBySide(img1, img2 string, tWidth int) (string, error) {
	img1Lines := strings.Split(img1, "\n")
	img2Lines := strings.Split(img2, "\n")
	img1Lines = img1Lines[:len(img1Lines)-1]
	img2Lines = img2Lines[:len(img2Lines)-1]

	img1Width := len(img1Lines[0])
	img2Width := len(img2Lines[0])

	if img1Width > tWidth/2 {
		img1Lines = trimSides(img1Lines, img1Width-tWidth/2)
		img1Width = len(img1Lines[0])
	}
	if img2Width > tWidth/2 {
		img2Lines = trimSides(img2Lines, img2Width-tWidth/2)
		img2Width = len(img2Lines[0])
	}

	vPadding := len(img1Lines) - len(img2Lines)
	if vPadding > 0 {
		img2Lines = applyPadding(img2Lines, vPadding)
	} else if vPadding < 0 {
		img1Lines = applyPadding(img1Lines, -vPadding)
	}

	combinedLines := make([]string, max(len(img1Lines), len(img2Lines)))
	gap := tWidth - img1Width - img2Width
	for i := range combinedLines {
		if i < len(img1Lines) {
			combinedLines[i] = img1Lines[i]
		} else {
			combinedLines[i] = strings.Repeat(" ", img1Width)
		}
		combinedLines[i] += strings.Repeat(" ", gap)
		if i < len(img2Lines) {
			combinedLines[i] += img2Lines[i]
		} else {
			combinedLines[i] += strings.Repeat(" ", img2Width)
		}
	}
	return strings.Join(combinedLines, "\n"), nil
}

func applyPadding(lines []string, padding int) []string {
	for i := 0; i < padding/2; i++ {
		lines = append([]string{strings.Repeat(" ", len(lines[0]))}, lines...)
	}
	return lines
}

func trimSides(lines []string, pixelsToRemove int) []string {
	for line := range lines {
		lines[line] = lines[line][2*(pixelsToRemove/2) : len(lines[line])-2*(pixelsToRemove-pixelsToRemove/2)]
	}
	return lines
}
