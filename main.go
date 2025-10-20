package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/mrbrist/pokedex-go/internal/api"
	"github.com/mrbrist/pokedex-go/internal/pokecache"
)

type config struct {
	Next     string
	Previous string
	cache    *pokecache.Cache
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

var commands map[string]cliCommand

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
	}

	scanner := bufio.NewScanner(os.Stdin)

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading input:", err)
	}

	cfg := config{
		Next:  "https://pokeapi.co/api/v2/location-area",
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
				c.callback(&cfg)
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

func commandExit(cfg *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	for _, c := range commands {
		fmt.Printf("%s: %s\n", c.name, c.description)
	}
	return nil
}

func commandMap(cfg *config) error {
	la := api.GetLocationAreas(cfg.cache, cfg.Next)
	cfg.Next = la.Next
	cfg.Previous = la.Previous

	for _, loc := range la.Results {
		fmt.Println(loc.Name)
	}
	return nil
}

func commandMapb(cfg *config) error {
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
