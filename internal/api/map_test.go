package api

import (
	"testing"
	"time"

	"github.com/mrbrist/pokedex-go/internal/pokecache"
)

func TestGetLocationAreas(t *testing.T) {
	cache := pokecache.NewCache(5 * time.Second)
	la := GetLocationAreas(cache, "https://pokeapi.co/api/v2/location-area")
	if la == nil {
		t.Errorf("Could not get data")
		return
	}
}
