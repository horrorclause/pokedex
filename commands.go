// Stores the command struct and functions

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"

	"github.com/horrorclause/pokedex/internal/pokecache"
)

// Pokemon Stats Struct commandInspect
type PokemonStat struct {
	BaseStat int `json:"base_stat"`
	Stat     struct {
		Name string `json:"name"`
	} `json:"stat"`
}

// Pokem Type Struct commandInspect
type PokemonType struct {
	Type struct {
		Name string `json:"name"`
	} `json:"type"`
}

// Struct for pokemon
type Pokemon struct {
	Name           string        `json:"name"`
	BaseExperience int           `json:"base_experience"`
	Height         int           `json:"height"`
	Weight         int           `json:"weight"`
	Stats          []PokemonStat `json:"stats"`
	Types          []PokemonType `json:"types"`
}

// Config Struct for saving the state of Prev and Next in map
type config struct {
	Next     *string
	Previous *string
	Cache    *pokecache.Cache
	Pokedex  map[string]Pokemon
}

// CLI Commands Struct for commands a user can use
type cliCommand struct {
	name        string
	description string
	callback    func(*config, ...string) error
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
		"explore": {
			name:        "explore",
			description: "Dive deeper into an area and discover what pokemon await!",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Try and catch a Pokemon you encounter",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspect Pokemon you have already caught",
			callback:    commandInspect,
		},
	}
}

// Exit Command function
func commandExit(cfg *config, args ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

// Help Command Function
func commandHelp(cfg *config, args ...string) error {
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

// Map Command Function
func commandMap(cfg *config, args ...string) error {

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

	// Check cache FIRST
	var body []byte
	cachedData, found := cfg.Cache.Get(url)

	if found {

		fmt.Println("[+] Using cached data...")
		fmt.Println()

		// Use cached data
		body = cachedData
	} else {

		// Make HTTP request
		res, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("An issue was encountered reaching the URL: %s", url)
		}

		defer res.Body.Close()

		// Read response body into bytes
		body, err = io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("error reading body")
		}

		// Save to cache
		cfg.Cache.Add(url, body)

	}

	var result Results
	err := json.Unmarshal(body, &result)
	if err != nil {
		return fmt.Errorf("Error decoding")
	}

	//Print to Terminal
	for _, area := range result.Results {
		fmt.Println(area.Name)
	}

	// Stores the Next and Prev URL in the config file
	cfg.Next = result.Next
	cfg.Previous = result.Previous

	return nil
}

// Map Back Command, moves backwards
func commandMapb(cfg *config, args ...string) error {

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

	// Check cache FIRST
	var body []byte
	cachedData, found := cfg.Cache.Get(url)

	if found {

		fmt.Println("[+] Using cached data...")
		fmt.Println()

		// Use cached data
		body = cachedData
	} else {

		// Make HTTP request
		res, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("An issue was encountered reaching the URL: %s", url)
		}

		defer res.Body.Close()

		// Read response body into bytes
		body, err = io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("error reading body")
		}

		// Save to cache
		cfg.Cache.Add(url, body)

	}

	var result Results
	err := json.Unmarshal(body, &result)
	if err != nil {
		return fmt.Errorf("Error decoding")
	}

	// Being output to terminal
	for _, area := range result.Results {
		fmt.Println(area.Name)
	}

	// Stores the Next and Prev URL in the config file
	cfg.Next = result.Next
	cfg.Previous = result.Previous

	return nil

}

// Explore Command, dive deeper into discovered locations
func commandExplore(cfg *config, args ...string) error {

	// ExploreResult contains a slice of PokemonEncounter
	// Each PokemonEncounter contains one Pokemon
	// Each Pokemon contains a Name string

	// Innermost level of pokemon encounter, the pokemon's name
	type Pokemon struct {
		Name string `json:"name"`
	}

	// Middle level, this is the Encounter itself
	type PokemonEncounter struct {
		Pokemon Pokemon `json:"pokemon"`
	}

	// Top level, this is the whole response for the area to explore
	type ExploreResult struct {
		PokemonEncounter []PokemonEncounter `json:"pokemon_encounters"`
	}

	if len(args) == 0 {
		return fmt.Errorf("No location supplied to explore")
	}

	url := "https://pokeapi.co/api/v2/location-area/"

	// First index of the arguments would be the specific location
	locationName := args[0]

	url = url + locationName

	fmt.Printf("Exploring %s...\n", locationName)

	// Checking if the data is cached first
	var body []byte

	cachedData, found := cfg.Cache.Get(url)

	if found {

		fmt.Println("[+] Using cached data...")
		fmt.Println()

		// Use cached data
		body = cachedData
	} else {

		// Make HTTP request
		res, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("An issue was encountered reaching the URL: %s", url)
		}

		defer res.Body.Close()

		// Read response body into bytes
		body, err = io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("error reading body")
		}

		// Save to cache
		cfg.Cache.Add(url, body)

	}

	var result ExploreResult

	err := json.Unmarshal(body, &result)
	if err != nil {
		return fmt.Errorf("Error decoding the result")
	}

	for _, pokemon := range result.PokemonEncounter {
		fmt.Println(pokemon.Pokemon.Name)
	}

	return nil
}

// Catch Command, try and catch some Pokemon
func commandCatch(cfg *config, args ...string) error {

	url := "https://pokeapi.co/api/v2/pokemon/"

	if len(args) == 0 {
		return fmt.Errorf("No pokemon listed...")
	}

	pokemonName := args[0]

	url = url + pokemonName

	fmt.Printf("Throwing a Pokeball at %s...", pokemonName)
	fmt.Println()

	var body []byte

	cachedData, found := cfg.Cache.Get(url)

	if found {

		fmt.Println("[+] Using cached data...")
		fmt.Println()

		// Use cached data
		body = cachedData
	} else {

		// Make HTTP request
		res, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("An issue was encountered reaching the URL: %s", url)
		}

		defer res.Body.Close()

		// Read response body into bytes
		body, err = io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("error reading body")
		}

		// Save to cache
		cfg.Cache.Add(url, body)

	}

	var result Pokemon

	err := json.Unmarshal(body, &result)
	if err != nil {
		return fmt.Errorf("Error decoding data")
	}

	randomNum := rand.Intn(200)

	if randomNum >= result.BaseExperience {
		cfg.Pokedex[result.Name] = result
		fmt.Printf("%s was caught!", result.Name)
		fmt.Println()

	} else {
		fmt.Printf("%s got away!", result.Name)
		fmt.Println()

	}

	return nil
}

// Inspect Command, inspect caught Pokemon
func commandInspect(cfg *config, args ...string) error {

	if len(args) == 0 {
		return fmt.Errorf("No Pokemon name supplied.")
	}

	pokemonName := args[0]

	pokemon, found := cfg.Pokedex[pokemonName]
	if !found {
		fmt.Printf("%s is not found, go and catch one!", pokemonName)
		fmt.Println()
		return nil
	}

	// Printing Pokemon details
	fmt.Println("Pokemon Details:")
	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)

	fmt.Println("Stats:")
	for _, stat := range pokemon.Stats {
		fmt.Printf(" -%s: %d", stat.Stat.Name, stat.BaseStat)
		fmt.Println()
	}

	fmt.Println("Type:")
	for _, t := range pokemon.Types {
		fmt.Printf(" - %s", t.Type.Name)
		fmt.Println()
	}

	return nil
}
