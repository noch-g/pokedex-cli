package settings

import "math"

const (
	Prompt        = "Pokedex > "
	SaveFilePath  = "pokemons.json"
	CacheFilePath = "cache.gob"
)

var CatchThreshold = map[string]int{
	"Pokeball":   40,
	"Superball":  60,
	"Hyperball":  85,
	"Masterball": math.MaxInt,
}
