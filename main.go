package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Numpkens/pokedexcli/internal/pokecache"
)

type config struct {
	Client              http.Client
	Cache               pokecache.Cache
	NextLocationURL     *string
	PreviousLocationURL *string
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config, []string) error
}

type LocationAreaResponse struct {
	Count    int              `json:"count"`
	Next     *string          `json:"next"`
	Previous *string          `json:"previous"`
	Results  []LocationResult `json:"results"`
}

type LocationResult struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type LocationAreaDetailResponse struct {
	Name              string             `json:"name"`
	PokemonEncounters []PokemonEncounter `json:"pokemon_encounters"`
}

type PokemonEncounter struct {
	Pokemon struct {
		Name string `json:"name"`
	} `json:"pokemon"`
}

func fetchData(url string, cfg *config) ([]byte, error) {
	var data []byte

	if val, ok := cfg.Cache.Get(url); ok {
		data = val
		fmt.Fprintln(os.Stderr, "Cache HIT: Using cached data.")
		return data, nil
	}

	fmt.Fprintln(os.Stderr, "Cache MISS: Fetching data from PokeAPI...")

	res, err := cfg.Client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode > 399 {
		return nil, fmt.Errorf("bad status code: %d", res.StatusCode)
	}

	data, err = io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	cfg.Cache.Add(url, data)
	return data, nil
}

func commandExit(cfg *config, args []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config, args []string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()

	commands := getCommands()

	// Dynamic printing of commands
	for _, cmd := range commands {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}

	fmt.Println()
	return nil
}

func commandMap(cfg *config, args []string) error {
	url := "https://pokeapi.co/api/v2/location-area/"
	if cfg.NextLocationURL != nil {
		url = *cfg.NextLocationURL
	}

	data, err := fetchData(url, cfg)
	if err != nil {
		return err
	}

	var locationResp LocationAreaResponse
	err = json.Unmarshal(data, &locationResp)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	for _, result := range locationResp.Results {
		fmt.Println(result.Name)
	}

	cfg.NextLocationURL = locationResp.Next
	cfg.PreviousLocationURL = locationResp.Previous

	return nil
}

func commandMapb(cfg *config, args []string) error {
	if cfg.PreviousLocationURL == nil || *cfg.PreviousLocationURL == "" {
		fmt.Println("You're on the first page.")
		return nil
	}

	url := *cfg.PreviousLocationURL

	data, err := fetchData(url, cfg)
	if err != nil {
		return err
	}

	// Unmarshal logic from the old getLocations
	var locationResp LocationAreaResponse
	err = json.Unmarshal(data, &locationResp)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	for _, result := range locationResp.Results {
		fmt.Println(result.Name)
	}

	cfg.NextLocationURL = locationResp.Next
	cfg.PreviousLocationURL = locationResp.Previous

	return nil
}

func commandExplore(cfg *config, args []string) error {
	if len(args) == 0 {
		return errors.New("you must provide a location area name, e.g., explore canalave-city-area")
	}
	areaName := args[0]

	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s", areaName)
	fmt.Printf("Exploring %s...\n", areaName)

	data, err := fetchData(url, cfg)
	if err != nil {
		return err
	}

	var detailResp LocationAreaDetailResponse
	err = json.Unmarshal(data, &detailResp)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON for area detail: %w", err)
	}

	if len(detailResp.PokemonEncounters) == 0 {
		fmt.Println("No Pokémon found in this area.")
		return nil
	}

	fmt.Println("Found Pokemon:")
	for _, encounter := range detailResp.PokemonEncounters {
		fmt.Printf(" - %s\n", encounter.Pokemon.Name)
	}

	return nil
}

// --- Command Registry ---

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Displays the next 20 location areas",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the previous 20 location areas",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore <area_name>",
			description: "Displays the Pokémon found in a specific location area",
			callback:    commandExplore,
		},
	}
}

// --- REPL Core Logic ---

func cleanInput(text string) []string {
	lowered := strings.ToLower(text)
	words := strings.Fields(lowered)
	return words
}

func startREPL() {
	emptyString := ""

	cfg := &config{
		Client: http.Client{
			Timeout: time.Second * 10,
		},
		Cache: pokecache.NewCache(5 * time.Minute),

		PreviousLocationURL: &emptyString,
	}

	initialURL := "https://pokeapi.co/api/v2/location-area/"
	cfg.NextLocationURL = &initialURL

	scanner := bufio.NewScanner(os.Stdin)
	commands := getCommands()

	for {
		fmt.Fprintf(os.Stderr, "Pokedex > ")

		if !scanner.Scan() {
			break
		}

		text := scanner.Text()
		cleaned := cleanInput(text)

		if len(cleaned) == 0 {
			continue
		}

		commandName := cleaned[0]
		args := cleaned[1:]

		if command, ok := commands[commandName]; ok {
			err := command.callback(cfg, args)
			if err != nil {
				fmt.Printf("Error executing command %s: %v\n", commandName, err)
			}
		} else {
			fmt.Printf("Unknown command: %s\n", commandName)
		}
	}
}

func main() {
	startREPL()
}
