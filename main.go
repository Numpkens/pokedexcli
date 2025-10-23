package main

import (
	"fmt"
	"strings"
)

func cleanInput(text string) []string {
	lowered := strings.ToLower(text)

	words := strings.Fields(lowered)

	return words
}

func main() {
	// Example usage: Call cleanInput with a test string and assign the result.
	testString := "  Hello Pokémon World 123 "
	cleanedWords := cleanInput(testString)

	// Print the result since main() cannot return values.
	fmt.Println("Original:", testString)
	fmt.Println("Cleaned:", cleanedWords)
}
