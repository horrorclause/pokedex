// Stores the command struct and functions

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// Struct for saving the state of Prev and Next in map
type config struct {
	Next     *string
	Previous *string
}

// Struct for commands a user can use
type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

// Map of Commands to be used
func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exits the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays the help menu",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Look forward 20 results on the map",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "'Map Back' - Look back 20 results on the map",
			callback:    commandMapb,
		},
	}
}

// Exit function
func commandExit(cfg *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

// Help Function
func commandHelp(cfg *config) error {
	fmt.Println()
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("=======================")
	fmt.Println()

	for _, cmd := range getCommands() {
		fmt.Printf("%s: %s\n\n", cmd.name, cmd.description)
	}

	fmt.Println("-----------------------")
	return nil
}

// Map Function
func commandMap(cfg *config) error {

	// Create the struct for the Area JSON info
	type Area struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}

	// Create the struct for the entire response
	type Results struct {
		Count    int     `json:"count"`
		Next     *string `json:"next"`
		Previous *string `json:"previous"`
		Results  []Area  `json:"results"`
	}

	url := "https://pokeapi.co/api/v2/location-area/"

	// Checks to see if the Config file has a "Next" url listed
	// If it does, use it
	if cfg.Next != nil {
		url = *cfg.Next
	}

	res, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("An issue was encountered reaching the URL: %s", url)
	}

	defer res.Body.Close()

	var results Results

	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&results); err != nil {
		return fmt.Errorf("Error decoding response body")
	}

	for _, area := range results.Results {
		fmt.Println(area.Name)
	}

	// Stores the Next and Prev URL in the config file
	cfg.Next = results.Next
	cfg.Previous = results.Previous

	return nil
}

// Map Back Command, moves backwards
func commandMapb(cfg *config) error {

	// Create the struct for the Area JSON info
	type Area struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}

	// Create the struct for the entire response
	type Results struct {
		Count    int     `json:"count"`
		Next     *string `json:"next"`
		Previous *string `json:"previous"`
		Results  []Area  `json:"results"`
	}

	// Checks to see if the Config file has a Previous URL
	// If it doesnt, notify user
	if cfg.Previous == nil {
		return fmt.Errorf("You are on the first page.")
	}

	// Set the URL to the Previous set URL
	url := *cfg.Previous

	res, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("An issue was encountered reaching the URL: %s", url)
	}

	defer res.Body.Close()

	var results Results

	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&results); err != nil {
		return fmt.Errorf("Error decoding response body")
	}

	for _, area := range results.Results {
		fmt.Println(area.Name)
	}

	// Stores the Next and Prev URL in the config file
	cfg.Next = results.Next
	cfg.Previous = results.Previous

	return nil

}
