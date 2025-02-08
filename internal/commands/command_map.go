package commands

import (
	"errors"
	"fmt"
	"io"

	"github.com/noch-g/pokedex-cli/internal/config"
	"github.com/noch-g/pokedex-cli/internal/text"
)

func commandMapf(cfg *config.Config, output io.Writer, args ...string) error {
	return commandMap(cfg, output, true)
}

func commandMapb(cfg *config.Config, output io.Writer, args ...string) error {
	if cfg.PrevLocationsURL == nil {
		return errors.New("you're on the first page")
	}

	return commandMap(cfg, output, false)
}

func commandMap(cfg *config.Config, output io.Writer, goForward bool) error {
	var next_url *string
	if goForward {
		next_url = cfg.NextLocationsURL
	} else {
		next_url = cfg.PrevLocationsURL
	}

	locationsResp, err := cfg.PokeapiClient.ListLocations(next_url)
	if err != nil {
		return err
	}

	cfg.NextLocationsURL = locationsResp.Next
	cfg.PrevLocationsURL = locationsResp.Previous

	searchedPokemons, err := getSearchedPokemons(cfg, 1, 151)
	if err != nil {
		return err
	}
	cfg.KnownEntities["locations"] = []string{}
	for _, loc := range locationsResp.Results {
		if locationContainsNew(cfg, loc.Name, searchedPokemons) {
			fmt.Fprintln(output, text.ToBold(loc.Name))
		} else {
			fmt.Fprintln(output, loc.Name)
		}
		cfg.KnownEntities["locations"] = append(cfg.KnownEntities["locations"], loc.Name)
	}
	return nil
}

func getSearchedPokemons(cfg *config.Config, start, end int) (map[string]struct{}, error) {
	respPokemons, err := cfg.PokeapiClient.GetPokemonList(start, end)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve searched pokemons: %v", err)
	}
	searchedPokemons := map[string]struct{}{}
	for _, p := range respPokemons.Results {
		searchedPokemons[p.Name] = struct{}{}
	}
	return searchedPokemons, nil
}

func locationContainsNew(cfg *config.Config, locationStr string, searchedPokemons map[string]struct{}) bool {
	loc, err := cfg.PokeapiClient.GetLocation(locationStr)
	if err != nil {
		return false
	}

	for _, pokemon_encounter := range loc.PokemonEncounters {
		if _, ok := cfg.CaughtPokemon[pokemon_encounter.Pokemon.Name]; !ok {
			if _, ok := searchedPokemons[pokemon_encounter.Pokemon.Name]; ok {
				return true
			}
		}
	}
	return false
}
