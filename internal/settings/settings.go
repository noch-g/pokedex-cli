package settings

import "math"

const (
	Prompt        = "Pokedex > "
	SaveFilePath  = "pokemons.json"
	CacheFilePath = "cache.gob"
)

var CatchThreshold = map[string]int{
	"Pokeball":   40,
	"Greatball":  60,
	"Ultraball":  80,
	"Masterball": math.MaxInt,
}
