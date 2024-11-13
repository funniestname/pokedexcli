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
	callback    func() error
}

type Config struct {
	cache    *pokecache.Cache
	Next     string
	Previous *string
}

func main() {
	myCache := pokecache.NewCache(1 * time.Minute)
	cm := &Config{
		cache:    myCache,
		Next:     "https://pokeapi.co/api/v2/location/",
		Previous: nil,
	}
	commands := cm.getCommands()
	repl(commands)
}

func repl(commands map[string]cliCommand) {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("pokedex > ")
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input: ", err)
			continue
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Split(line, " ")
		commandName := parts[0]
		command, ok := commands[commandName]
		if !ok {
			fmt.Printf("Unknown command: %s\n", commandName)
			continue
		}

		err = command.callback()
		if err != nil {
			fmt.Println("Error executing command:", err)
		}
	}
}

func (cm *Config) getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    cm.commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Display 20 location areas",
			callback:    cm.commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Display the previous 20 location areas",
			callback:    cm.commandMapb,
		},
	}
}

func (cm *Config) commandHelp() error {
	fmt.Println("Available commands:")
	for _, command := range cm.getCommands() {
		fmt.Printf(" %s: %s\n", command.name, command.description)
	}
	return nil
}

func commandExit() error {
	fmt.Println("Exiting the pokedex...")
	os.Exit(0)
	return nil
}

func (cm *Config) commandMap() error {
	values, ok := cm.cache.Get(cm.Next)
	if ok {
		var locations pokemap.PokeLocation
		err := json.Unmarshal(values, &locations)
		if err != nil {
			return err
		}
		cm.Next = locations.Next
		cm.Previous = locations.Previous

		for _, val := range locations.Results {
			fmt.Println(val.Name)
		}
		return nil
	}
	locations, err := pokemap.GetPokeMap(cm.Next)
	if err != nil {
		return err
	}

	locationBytes, err := json.Marshal(locations)
	if err != nil {
		return err
	}
	cm.cache.Add(cm.Next, locationBytes)

	cm.Next = locations.Next
	cm.Previous = locations.Previous

	for _, val := range locations.Results {
		fmt.Println(val.Name)
	}
	return nil
}

func (cm *Config) commandMapb() error {
	if cm.Previous == nil {
		fmt.Println("No previous page found")
		return nil
	}

	values, ok := cm.cache.Get(*cm.Previous)
	if ok {
		var locations pokemap.PokeLocation
		err := json.Unmarshal(values, &locations)
		if err != nil {
			return err
		}
		cm.Next = locations.Next
		cm.Previous = locations.Previous

		for _, val := range locations.Results {
			fmt.Println(val.Name)
		}
		return nil
	}

	locations, err := pokemap.GetPokeMap(*cm.Previous)
	if err != nil {
		return err
	}

	locationBytes, err := json.Marshal(locations)
	if err != nil {
		return err
	}
	cm.cache.Add(*cm.Previous, locationBytes)

	cm.Next = locations.Next
	cm.Previous = locations.Previous

	for _, val := range locations.Results {
		fmt.Println(val.Name)
	}
	return nil
}
