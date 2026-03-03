package main

import (
	"strings"
)

func cleanInput(text string) []string {

	userInput := strings.Fields(strings.ToLower(text))

	return userInput
}
