package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/mrbrist/pokedex-go/api"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
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
	}

	scanner := bufio.NewScanner(os.Stdin)

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading input:", err)
	}

	for {
		fmt.Print("Pokedex > ")
		if scanner.Scan() {
			input := cleanInput(scanner.Text())
			c, ok := commands[input[0]]
			if !ok {
				fmt.Println("Unknown command")
			} else {
				c.callback()
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

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp() error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	for _, c := range commands {
		fmt.Printf("%s: %s\n", c.name, c.description)
	}
	return nil
}

func commandMap() error {
	api.Map()
	return nil
}
