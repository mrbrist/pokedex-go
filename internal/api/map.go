package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/mrbrist/pokedex-go/internal/pokecache"
)

type LocationAreas struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func GetLocationAreas(c *pokecache.Cache, url string) *LocationAreas {
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

	la := LocationAreas{}
	err := json.Unmarshal(body, &la)
	if err != nil {
		fmt.Println(err)
	}

	return &la
}
