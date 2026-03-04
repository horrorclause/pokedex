// Stores the command struct and functions

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// Struct for commands a user can use
type cliCommand struct {
	name        string
	description string
	callback    func() error
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
			description: "Displays the map 20 locations at a time",
			callback:    commandMap,
		},
	}
}

// Exit function
func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

// Help Function
func commandHelp() error {
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
func commandMap() error {

	// Create the struct for the Area JSON info
	type Area struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}

	// Create the struct for the entire response
	type Results struct {
		Count    int    `json:"count"`
		Next     string `json:"next"`
		Previous string `json:"previous"`
		Results  []Area `json:"results"`
	}

	url := "https://pokeapi.co/api/v2/location-area/"

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

	return nil
}
