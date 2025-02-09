package config

import (
	"encoding/json"
	"os"

	"github.com/noch-g/pokedex-cli/internal/pokeapi"
	"github.com/noch-g/pokedex-cli/internal/text"
)

const Prompt = "Pokedex > "

type Config struct {
	CaughtPokemon    map[string]pokeapi.Pokemon `json:"pokemons"`
	KnownEntities    map[string][]string        `json:"-"`
	PokeapiClient    pokeapi.Client             `json:"-"`
	NextLocationsURL *string
	PrevLocationsURL *string
}

func NewConfig(pokeapiClient pokeapi.Client) *Config {
	return &Config{
		CaughtPokemon: make(map[string]pokeapi.Pokemon),
		KnownEntities: make(map[string][]string),
		PokeapiClient: pokeapiClient,
	}
}

func (cfg *Config) Save(filename string) error {
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

func (cfg *Config) Load(filename string, commandNames []string) error {
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
	if len(cfg.CaughtPokemon) > 0 {
		for pokemonName := range cfg.CaughtPokemon {
			cfg.KnownEntities["pokemons"] = append(cfg.KnownEntities["pokemons"], pokemonName)
		}
	}
	cfg.KnownEntities["commands"] = append(cfg.KnownEntities["commands"], commandNames...)
	return nil
}

func GetPromptLength() int {
	return len(Prompt)
}

func GetPromptMessage() string {
	return text.ToBold(Prompt)
}
