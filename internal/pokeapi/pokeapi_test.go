package pokeapi

import (
	"fmt"
	"testing"
	"time"
)

func TestPokemonGet(t *testing.T) {
	pokeClient := NewClient(5*time.Second, 5*time.Minute)

	cases := []struct {
		input    string
		expected string
	}{
		{
			input:    "charmander",
			expected: "charmander fire 39",
		},
		{
			input:    "squirtle",
			expected: "squirtle water 44",
		},
	}
	for i, c := range cases {
		t.Run(fmt.Sprintf("Test case %v", i), func(t *testing.T) {
			pokemon, err := pokeClient.GetPokemon(c.input)
			if err != nil {
				t.Errorf("expected to get pokemon")
			}
			checkStr := fmt.Sprintf("%s %s %d", pokemon.Name, pokemon.Types[0].Type.Name, pokemon.Stats[0].BaseStat)
			if checkStr != c.expected {
				t.Errorf("expected different stats for pokemon\nExpected: %s\nGot: %s", c.expected, checkStr)
			}
		})
	}
}

func TestPokemonList(t *testing.T) {
	pokeClient := NewClient(5*time.Second, 5*time.Minute)

	cases := []struct {
		start int
		end   int
		first string
		last  string
	}{
		{
			start: 1,
			end:   151,
			first: "bulbasaur",
			last:  "mew",
		},
		{
			start: 10,
			end:   14,
			first: "caterpie",
			last:  "kakuna",
		},
	}
	for i, c := range cases {
		t.Run(fmt.Sprintf("Test case %v", i), func(t *testing.T) {
			pokemonList, err := pokeClient.GetPokemonList(c.start, c.end)
			if err != nil {
				t.Errorf("expected to get pokemon list")
			}
			if len(pokemonList.Results) != c.end-c.start+1 {
				t.Errorf("expected %d pokemons, got %d", c.end-c.start+1, pokemonList.Count)
			}
			firstPokemon := pokemonList.Results[0]
			lastPokemon := pokemonList.Results[len(pokemonList.Results)-1]
			if firstPokemon.Name != c.first {
				t.Errorf("expected first pokemon to be %s, got %s", c.first, firstPokemon.Name)
			}
			if lastPokemon.Name != c.last {
				t.Errorf("expected first pokemon to be %s, got %s", c.first, lastPokemon.Name)
			}
		})
	}
}

func TestLocationGet(t *testing.T) {
	pokeClient := NewClient(5*time.Second, 5*time.Minute)

	cases := []struct {
		locationName string
		expectedName string
		expectedId   int
	}{
		{
			locationName: "canalave-city-area",
			expectedName: "canalave-city-area",
			expectedId:   1,
		},
		{
			locationName: "mt-coronet-2f",
			expectedName: "mt-coronet-2f",
			expectedId:   12,
		},
	}
	for i, c := range cases {
		t.Run(fmt.Sprintf("Test case %v", i), func(t *testing.T) {
			location, err := pokeClient.GetLocation(c.locationName)
			if err != nil {
				t.Errorf("expected to get a location")
			}
			if location.Name != c.expectedName {
				t.Errorf("expected location name to be %s, got %s", c.expectedName, location.Name)
			}
			if location.ID != c.expectedId {
				t.Errorf("expected location ID to be %d, got %d", c.expectedId, location.ID)
			}
		})
	}
}

func TestLocationList(t *testing.T) {
	pokeClient := NewClient(5*time.Second, 5*time.Minute)

	cases := []struct {
		pageUrl           *string
		firstLocation     string
		lastLocation      string
		firstLocationNext string
		lastLocationNext  string
	}{
		{
			pageUrl:           nil,
			firstLocation:     "canalave-city-area",
			lastLocation:      "mt-coronet-1f-from-exterior",
			firstLocationNext: "mt-coronet-1f-route-216",
			lastLocationNext:  "solaceon-ruins-b3f-c",
		},
	}
	for i, c := range cases {
		t.Run(fmt.Sprintf("Test case %v", i), func(t *testing.T) {
			locationList, err := pokeClient.ListLocations(nil)
			if err != nil {
				t.Errorf("expected to get a locations list")
			}
			if len(locationList.Results) != 20 {
				t.Errorf("expected %d locations, got %d", 20, locationList.Count)
			}
			firstLocation := locationList.Results[0].Name
			lastLocation := locationList.Results[len(locationList.Results)-1].Name
			if firstLocation != c.firstLocation {
				t.Errorf("expected first location to be %s, got %s", c.firstLocation, firstLocation)
			}
			if lastLocation != c.lastLocation {
				t.Errorf("expected last location to be %s, got %s", c.lastLocation, lastLocation)
			}

			locationNextList, err := pokeClient.ListLocations(locationList.Next)
			if err != nil {
				t.Errorf("expected to get a locations list")
			}
			if len(locationNextList.Results) != 20 {
				t.Errorf("expected %d locations, got %d", 20, locationList.Count)
			}
			firstLocationNext := locationNextList.Results[0].Name
			lastLocationNext := locationNextList.Results[len(locationNextList.Results)-1].Name
			if firstLocationNext != c.firstLocationNext {
				t.Errorf("expected first next location to be %s, got %s", c.firstLocationNext, firstLocationNext)
			}
			if lastLocationNext != c.lastLocationNext {
				t.Errorf("expected last next location to be %s, got %s", c.lastLocationNext, lastLocationNext)
			}
		})
	}
}
