package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/funniestname/pokemap"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

type config struct {
	Next     string
	Previous string
}

func main() {
	commands := getCommands()
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

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Display 20 location areas",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Display the previous 20 location areas",
			callback:    commandMapb,
		},
	}
}

func commandHelp() error {
	fmt.Println("Available commands:")
	for _, command := range getCommands() {
		fmt.Printf(" %s: %s\n", command.name, command.description)
	}
	return nil
}

func commandExit() error {
	fmt.Println("Exiting the pokedex...")
	os.Exit(0)
	return nil
}

func commandMap() error {
	url := "https://pokeapi.co/api/v2/location/"
	locations, err := pokemap.GetPokeMap(url)
	if err != nil {
		return err
	}
	for _, val := range locations {
		fmt.Println(val)
	}
	return nil
}

func commandMapb() error {
	return nil
}
