package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/mrbrist/pokedex-go/internal/api"
	"github.com/mrbrist/pokedex-go/internal/pokecache"
)

type config struct {
	Next     string
	Previous string
	Base     string
	cache    *pokecache.Cache
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config, string) error
}

var commands map[string]cliCommand
var pokedex map[string]api.PokemonData

func main() {
	commands = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays help info",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Displays the next 20 locations of the Pokemon map",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the previous 20 locations of the Pokemon map",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "Displays the Pokemon in a selected area",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Attempts to catch the Pokemon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspects a Pokemon in your Pokedex",
			callback:    commandInspect,
		},
	}

	pokedex = map[string]api.PokemonData{}

	scanner := bufio.NewScanner(os.Stdin)

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading input:", err)
	}

	cfg := config{
		Next:  "https://pokeapi.co/api/v2/location-area",
		Base:  "https://pokeapi.co/api/v2/location-area/",
		cache: pokecache.NewCache(5 * time.Second),
	}

	for {
		fmt.Print("Pokedex > ")
		if scanner.Scan() {
			input := cleanInput(scanner.Text())
			c, ok := commands[input[0]]
			if !ok {
				fmt.Println("Unknown command")
			} else {
				if len(input) > 1 {
					p1 := input[1]
					c.callback(&cfg, p1)
				} else {
					c.callback(&cfg, "")
				}
			}
		}
	}
}

func cleanInput(text string) []string {
	input := strings.TrimSpace(text)
	if input == "" {
		return []string{}
	}

	words := strings.Fields(input)
	for i, w := range words {
		words[i] = strings.ToLower(w)
	}
	return words
}

func commandExit(cfg *config, p1 string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config, p1 string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	for _, c := range commands {
		fmt.Printf("%s: %s\n", c.name, c.description)
	}
	return nil
}

func commandMap(cfg *config, p1 string) error {
	la := api.GetLocationAreas(cfg.cache, cfg.Next)
	cfg.Next = la.Next
	cfg.Previous = la.Previous

	for _, loc := range la.Results {
		fmt.Println(loc.Name)
	}
	return nil
}

func commandMapb(cfg *config, p1 string) error {
	if cfg.Previous == "" {
		fmt.Println("You are on the first page")
	} else {
		la := api.GetLocationAreas(cfg.cache, cfg.Previous)
		cfg.Next = la.Next
		cfg.Previous = la.Previous

		for _, loc := range la.Results {
			fmt.Println(loc.Name)
		}
	}
	return nil
}

func commandExplore(cfg *config, p1 string) error {
	if p1 != "" {
		lad := api.GetExploreData(cfg.cache, cfg.Base, p1)

		for _, p := range lad.PokemonEncounters {
			fmt.Printf("- %s\n", p.Pokemon.Name)
		}
		return nil
	}
	fmt.Println("Please pick a location")
	return nil
}

func commandCatch(cfg *config, p1 string) error {
	if p1 != "" {
		fmt.Printf("Throwing a Pokeball at %s...\n", p1)
		pd := api.GetPokemonData(cfg.cache, p1)

		// This is 5050 not the best
		catchChance := rand.Intn(pd.BaseExperience)
		if catchChance > pd.BaseExperience/2 {
			pokedex[p1] = *pd
			fmt.Printf("%s was caught!\n", pd.Name)
		} else {
			fmt.Printf("%s escaped!\n", pd.Name)
		}

		return nil
	}
	fmt.Println("Please pick a Pokemon")
	return nil
}

func commandInspect(cfg *config, p1 string) error {
	if p1 != "" {
		p, ok := pokedex[p1]
		if !ok {
			fmt.Println("you have not caught that pokemon")
			return nil
		}
		// Display stats here
		fmt.Println(p)
		return nil
	}
	fmt.Println("Please pick a Pokemon")
	return nil
}
