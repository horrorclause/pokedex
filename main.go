package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/horrorclause/pokedex/internal/pokecache"
)

func main() {

	cfg := &config{
		Cache:   pokecache.NewCache(5 * time.Second), // Setting 5 Second Cache interval
		Pokedex: make(map[string]Pokemon),
	}

	// Initiate a new scanner to capture user input
	scanner := bufio.NewScanner(os.Stdin)

	// Main REPL
	for {
		fmt.Print("Pokedex > ") // CLI beginning
		scanner.Scan()          // Waits for user input

		userInput := cleanInput(scanner.Text()) // Captures user input and cleans it
		if len(userInput) == 0 {                // Prevents Panic for nothing submitted
			continue
		}

		commandInput := userInput[0] // Checks the first word for the command

		cmd, exists := getCommands()[commandInput] // Checks if the command exists

		if exists {
			err := cmd.callback(cfg, userInput[1:]...) // Captures any errors that may be thrown
			if err != nil {
				fmt.Println(err)
			}
			continue
		} else {
			fmt.Println("Unknown command")
			continue
		}

	}
}
