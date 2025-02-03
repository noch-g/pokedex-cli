package main

import (
	"encoding/json"
	"os"

	"github.com/noch-g/pokedex-cli/internal/pokeapi"
)

type config struct {
	CaughtPokemon    map[string]pokeapi.Pokemon `json:"pokemons"`
	pokeapiClient    pokeapi.Client
	nextLocationsURL *string
	prevLocationsURL *string
}

func (cfg *config) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(cfg)
	if err != nil {
		return err
	}

	return nil
}

func (cfg *config) Load(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(cfg)
	if err != nil {
		return err
	}

	return nil
}
