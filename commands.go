// Stores the command struct and functions

package main

import (
	"fmt"
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
