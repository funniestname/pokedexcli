package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/funniestname/pokedexcli/pokecache"
	"github.com/funniestname/pokedexcli/pokemap"
)

type cliCommand struct {
	name        string
	description string
	params      []string
	callback    func([]string) error
}

type Config struct {
	cache    *pokecache.Cache
	Next     string
	Previous *string
}

func main() {
	myCache := pokecache.NewCache(1 * time.Minute)
	cfg := &Config{
		cache:    myCache,
		Next:     "https://pokeapi.co/api/v2/location/",
		Previous: nil,
	}
	commands := cfg.getCommands()
	repl(commands)
}

func repl(commands map[string]cliCommand) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("pokedex > ")
		scanner.Scan()
		input := scanner.Text()

		if len(input) == 0 {
			continue
		}

		words := strings.Fields(input)
		commandName := words[0]
		args := words[1:]

		command, ok := commands[commandName]
		if !ok {
			fmt.Printf("Unknown command: %s\n", commandName)
			continue
		}

		err := command.callback(args)
		if err != nil {
			fmt.Println("Error executing command:", err)
		}
	}
}

func (cfg *Config) getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			params:      []string{},
			callback:    cfg.commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the pokedex",
			params:      []string{},
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Display 20 location areas",
			params:      []string{},
			callback:    cfg.commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Display the previous 20 location areas",
			params:      []string{},
			callback:    cfg.commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "Display Pokemon in a specified location",
			params:      []string{"location_name"},
			callback:    cfg.commandExplore,
		},
	}
}

func (cfg *Config) commandHelp([]string) error {
	fmt.Println("Available commands:")
	for _, command := range cfg.getCommands() {
		fmt.Printf(" %s: %s\n", command.name, command.description)
	}
	return nil
}

func commandExit([]string) error {
	fmt.Println("Exiting the pokedex...")
	os.Exit(0)
	return nil
}

func (cfg *Config) commandMap([]string) error {
	values, ok := cfg.cache.Get(cfg.Next)
	if ok {
		var locations pokemap.PokeLocation
		err := json.Unmarshal(values, &locations)
		if err != nil {
			return err
		}
		cfg.Next = locations.Next
		cfg.Previous = locations.Previous

		for _, val := range locations.Results {
			fmt.Println(val.Name)
		}
		return nil
	}
	locations, err := pokemap.GetPokeMap(cfg.Next)
	if err != nil {
		return err
	}

	locationBytes, err := json.Marshal(locations)
	if err != nil {
		return err
	}
	cfg.cache.Add(cfg.Next, locationBytes)

	cfg.Next = locations.Next
	cfg.Previous = locations.Previous

	for _, val := range locations.Results {
		fmt.Println(val.Name)
	}
	return nil
}

func (cfg *Config) commandMapb([]string) error {
	if cfg.Previous == nil {
		fmt.Println("No previous page found")
		return nil
	}

	values, ok := cfg.cache.Get(*cfg.Previous)
	if ok {
		var locations pokemap.PokeLocation
		err := json.Unmarshal(values, &locations)
		if err != nil {
			return err
		}
		cfg.Next = locations.Next
		cfg.Previous = locations.Previous

		for _, val := range locations.Results {
			fmt.Println(val.Name)
		}
		return nil
	}

	locations, err := pokemap.GetPokeMap(*cfg.Previous)
	if err != nil {
		return err
	}

	locationBytes, err := json.Marshal(locations)
	if err != nil {
		return err
	}
	cfg.cache.Add(*cfg.Previous, locationBytes)

	cfg.Next = locations.Next
	cfg.Previous = locations.Previous

	for _, val := range locations.Results {
		fmt.Println(val.Name)
	}
	return nil
}

func (cfg *Config) commandExplore(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: explore 'area'")
	}
	fullUrl := "https://pokeapi.co/api/v2/location-area/" + args[0]
	fmt.Println("Exploring " + args[0])

	values, ok := cfg.cache.Get(fullUrl)
	if ok {
		var pokemonEncounters pokemap.PokemonEncounters
		err := json.Unmarshal(values, &pokemonEncounters)
		if err != nil {
			return err
		}
		fmt.Println("Found Pokemon:")
		for _, val := range pokemonEncounters {
			fmt.Println(" - " + val.Pokemon.Name)
		}
		return nil
	}

	pokemonEncounters, err := pokemap.GetPokemon(fullUrl)
	if err != nil {
		return err
	}

	pokemonBytes, err := json.Marshal(pokemonEncounters)
	if err != nil {
		return err
	}
	cfg.cache.Add(fullUrl, pokemonBytes)

	fmt.Println("Found Pokemon:")
	for _, val := range pokemonEncounters {
		fmt.Println(" - " + val.Pokemon.Name)
	}
	return nil
}
