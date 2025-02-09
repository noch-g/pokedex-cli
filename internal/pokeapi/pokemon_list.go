package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/noch-g/pokedex-cli/internal/logger"
)

func (c *Client) GetPokemonList(start, end int) (RespPokemons, error) {
	url := baseURL + fmt.Sprintf("/pokemon?offset=%d&limit=%d", start-1, end-start+1)

	// Check cache before request
	if val, ok := c.cache.Get(url); ok {
		logger.Debug("FROM CACHE", "url", url)
		pokemonResp := RespPokemons{}
		err := json.Unmarshal(val, &pokemonResp)
		if err != nil {
			return RespPokemons{}, err
		}
		return pokemonResp, nil
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return RespPokemons{}, err
	}
	logger.Debug("GET", "url", url)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return RespPokemons{}, err
	}
	defer resp.Body.Close()

	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		return RespPokemons{}, err
	}

	pokemonResp := RespPokemons{}
	err = json.Unmarshal(dat, &pokemonResp)
	if err != nil {
		return RespPokemons{}, err
	}

	c.cache.Add(url, dat)
	return pokemonResp, nil
}
