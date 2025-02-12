package commands

import (
	"errors"
	"fmt"
	"io"
	"math/rand"

	"github.com/noch-g/pokedex-cli/internal/config"
	"github.com/noch-g/pokedex-cli/internal/settings"
)

func commandCatch(cfg *config.Config, output io.Writer, args ...string) error {
	if len(args) < 1 {
		return errors.New("you must provide a pokemon name")
	}
	pokemonName := args[0]
	ballType := "Pokeball"

	for _, arg := range args[1:] {
		if arg == "--super-ball" {
			ballType = "Superball"
		} else if arg == "--hyper-ball" {
			ballType = "Hyperball"
		} else if arg == "--master-ball" {
			ballType = "Masterball"
		}
	}
	catchThreshold := settings.CatchThreshold[ballType]

	pokemon, err := cfg.PokeapiClient.GetPokemon(pokemonName)
	if err != nil {
		return err
	}

	res := rand.Intn(pokemon.BaseExperience)
	successChance := min(100, 100*float64(catchThreshold)/float64(pokemon.BaseExperience))

	fmt.Fprintf(output, "Throwing a %s at %s... (success chance: %.2f%%)\n", ballType, pokemon.Name, successChance)
	if res > catchThreshold {
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
