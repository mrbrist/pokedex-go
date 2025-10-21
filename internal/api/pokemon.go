package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/mrbrist/pokedex-go/internal/pokecache"
)

func GetPokemonData(c *pokecache.Cache, pokemon string) *PokemonData {
	url := "https://pokeapi.co/api/v2/pokemon/" + pokemon

	cachedBody, ok := c.Get(url)
	var body []byte

	if ok {
		body = cachedBody
	} else {
		res, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}
		dBody, err := io.ReadAll(res.Body)
		defer res.Body.Close()

		if res.StatusCode > 299 {
			log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, dBody)
		}
		if err != nil {
			log.Fatal(err)
		}
		c.Add(url, dBody)
		body = dBody
	}

	pd := PokemonData{}
	err := json.Unmarshal(body, &pd)
	if err != nil {
		fmt.Println(err)
	}

	return &pd
}
