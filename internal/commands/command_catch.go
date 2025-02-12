package commands

import (
	"errors"
	"fmt"
	"io"
	"math/rand"

	"github.com/noch-g/pokedex-cli/internal/config"
)

func commandCatch(cfg *config.Config, output io.Writer, args ...string) error {
	if len(args) != 1 {
		return errors.New("you must provide a pokemon name")
	}
	pokemonName := args[0]

	pokemon, err := cfg.PokeapiClient.GetPokemon(pokemonName)
	if err != nil {
		return err
	}

	res := rand.Intn(pokemon.BaseExperience)
	successChance := 100 * 40.0 / float64(pokemon.BaseExperience)

	fmt.Fprintf(output, "Throwing a Pokeball at %s... (success chance: %.2f%%)\n", pokemon.Name, successChance)
	if res > 40 {
		escapeMessages := []string{
			"Fail, %s broke free!",
			"Fail, %s slipped away!",
			"Fail, %s dodged the Pokéball!",
			"Fail, %s is too strong and ran away!",
		}
		fmt.Fprintf(output, escapeMessages[rand.Intn(len(escapeMessages))]+"\n", pokemon.Name)
		return nil
	}
	fmt.Fprintf(output, "%s was caught! (#%03d)\n", pokemon.Name, pokemon.ID)
	imgStr, err := cfg.PokeapiClient.RenderImage(&pokemon)
	if err != nil {
		fmt.Fprintf(output, "(Image could not be retrieved)\n")
	} else {
		fmt.Fprint(output, imgStr)
	}

	if _, ok := cfg.CaughtPokemon[pokemon.Name]; !ok {
		fmt.Fprintf(output, "The information was added to the pokedex (#%03d). You may now inspect it with the inspect command.\n", pokemon.ID)
		cfg.CaughtPokemon[pokemon.Name] = pokemon
		cfg.KnownEntities["pokemons"] = append(cfg.KnownEntities["pokemons"], pokemon.Name)
	} else {
		fmt.Fprintf(output, "You already had a %s, but it's always nice to make a new friend!\n", pokemon.Name)
	}

	return nil
}
