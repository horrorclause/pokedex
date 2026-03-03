package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// TODO: Create support for REPL
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")                 // CLI beginning
		scanner.Scan()                          // Waits for user input
		userInput := cleanInput(scanner.Text()) // Captures user input and cleans it

		if len(userInput) > 0 {
			fmt.Println("Your command was:", userInput[0]) // Prints first word user submitted
		}
	}

}
