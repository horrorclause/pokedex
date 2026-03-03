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

// Exit function
func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

// Help Functions
func commandHelp() error {
	fmt.Println(`Welcome to the Pokedex!
	
Usage:

help: Displays a help message
exit: Exit the Pokedex
	`)
	return nil
}
