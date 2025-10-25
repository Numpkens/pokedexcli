package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func cleanInput(text string) []string {
	lowered := strings.ToLower(text)

	words := strings.Fields(lowered)

	return words
}

func startREPL() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Fprintf(os.Stderr, "Pokedex > ")

		if !scanner.Scan() {
			break // Exit loop on EOF or error.
		}

		text := scanner.Text()
		cleaned := cleanInput(text)

		if len(cleaned) == 0 {
			continue
		}

		commandName := cleaned[0]

		fmt.Printf("Your command was: %s\n", commandName)
	}
}

func main() {
	startREPL()
}

