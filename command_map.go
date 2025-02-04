package main

import (
	"errors"
	"fmt"
	"io"
)

func commandMapf(cfg *config, output io.Writer, args ...string) error {
	return commandMap(cfg, output, true)
}

func commandMapb(cfg *config, output io.Writer, args ...string) error {
	if cfg.prevLocationsURL == nil {
		return errors.New("you're on the first page")
	}

	return commandMap(cfg, output, false)
}

func commandMap(cfg *config, output io.Writer, goForward bool) error {
	var next_url *string
	if goForward {
		next_url = cfg.nextLocationsURL
	} else {
		next_url = cfg.prevLocationsURL
	}

	locationsResp, err := cfg.pokeapiClient.ListLocations(next_url)
	if err != nil {
		return err
	}

	cfg.nextLocationsURL = locationsResp.Next
	cfg.prevLocationsURL = locationsResp.Previous

	searchedPokemons, err := getSearchedPokemons(cfg, 1, 151)
	if err != nil {
		return err
	}
	cfg.knownEntities["locations"] = []string{}
	for _, loc := range locationsResp.Results {
		if locationContainsNew(cfg, loc.Name, searchedPokemons) {
			fmt.Fprintln(output, ToBold(loc.Name))
		} else {
			fmt.Fprintln(output, loc.Name)
		}
		cfg.knownEntities["locations"] = append(cfg.knownEntities["locations"], loc.Name)
	}
	return nil
}

func getSearchedPokemons(cfg *config, start, end int) (map[string]struct{}, error) {
	respPokemons, err := cfg.pokeapiClient.GetPokemonList(start, end)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve searched pokemons: %v", err)
	}
	searchedPokemons := map[string]struct{}{}
	for _, p := range respPokemons.Results {
		searchedPokemons[p.Name] = struct{}{}
	}
	return searchedPokemons, nil
}

func locationContainsNew(cfg *config, locationStr string, searchedPokemons map[string]struct{}) bool {
	loc, err := cfg.pokeapiClient.GetLocation(locationStr)
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
