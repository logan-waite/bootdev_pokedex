package main

import (
	"fmt"
	"github.com/logan-waite/bootdev_pokedex/internal/pokeapi"
	"os"
)

// Command Registry
type cliCommand struct {
	name        string
	description string
	callback    func() error
}

var commands = map[string]cliCommand{}

func initCommands() {
	commands = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "List the next 20 locations in the database",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "list the previous 20 locations in the database",
			callback:    commandMapb,
		},
	}
}

// Command Callbacks
func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp() error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Print("Usage:\n\n")

	for _, command := range commands {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	fmt.Print("\n")
	return nil
}

func commandMap() error {
	result, err := pokeapi.GetLocationAreas("next")
	if err != nil {
		return err
	}
	for _, location := range result {
		fmt.Println(location.Name)
	}
	return nil
}

func commandMapb() error {
	result, err := pokeapi.GetLocationAreas("prev")
	if err != nil {
		return err
	}
	for _, location := range result {
		fmt.Println(location.Name)
	}
	return nil
}
