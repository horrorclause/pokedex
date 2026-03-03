package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {

	// Map of commands to be used
	commands := map[string]cliCommand{
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

	// Initiate a new scanner to capture user input
	scanner := bufio.NewScanner(os.Stdin)

	// Main REPL Loop
	for {
		fmt.Print("Pokedex > ")                 // CLI beginning
		scanner.Scan()                          // Waits for user input
		userInput := cleanInput(scanner.Text()) // Captures user input and cleans it

		if len(userInput) == 0 { // Prevents Panic for nothing submitted
			continue
		}

		commandInput := userInput[0] // Checks the first word for the command

		cmd, exists := commands[commandInput] // Checks if the command exists

		if exists {
			err := cmd.callback() // Captures any errors that may be thrown
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("Unknown command")
		}

	}
}
