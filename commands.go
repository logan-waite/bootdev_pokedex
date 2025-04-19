package main

import (
	"fmt"
	"math/rand/v2"
	"os"

	"github.com/logan-waite/bootdev_pokedex/internal/pokeapi"
)

func initPokedex() {
	initCommands()
	initPokemonList()
}

// List of Caught Pokemon
var pokemonList map[string]pokeapi.Pokemon

func initPokemonList() {
	pokemonList = map[string]pokeapi.Pokemon{}
}

// Command Registry
type cliCommand struct {
	name        string
	description string
	callback    func(arg string) error
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
		"explore": {
			name:        "explore",
			description: "explore a location (explore <location name>)",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "attempt to catch a pokemon (catch <pokemon name>)",
			callback:    commandCatch,
		},
	}
}

// Command Callbacks
func commandExit(_ string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(_ string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Print("Usage:\n\n")

	for _, command := range commands {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	fmt.Print("\n")
	return nil
}

func commandMap(_ string) error {
	result, err := pokeapi.GetLocationAreas("next")
	if err != nil {
		return err
	}
	for _, location := range result {
		fmt.Println(location.Name)
	}
	return nil
}

func commandMapb(_ string) error {
	result, err := pokeapi.GetLocationAreas("prev")
	if err != nil {
		return err
	}
	for _, location := range result {
		fmt.Println(location.Name)
	}
	return nil
}

func commandExplore(locationArg string) error {
	if locationArg == "" {
		return fmt.Errorf("need a location to explore")
	}
	location, err := pokeapi.GetLocationAreaData(locationArg)
	if err != nil {
		return err
	}
	for _, pokemon := range location.PokemonEncounters {
		fmt.Println(pokemon.Pokemon.Name)
	}
	return nil
}

func commandCatch(pokemonArg string) error {
	fmt.Printf("Throwing a Pokeball at %s\n", pokemonArg)
	if pokemonArg == "" {
		return fmt.Errorf("need a pokemon to catch!")
	}
	pokemon, err := pokeapi.GetPokemon(pokemonArg)
	if err != nil {
		return err
	}
	target := pokemon.BaseExperience - (25 + (pokemon.BaseExperience / 10))
	attempt := rand.IntN(pokemon.BaseExperience)
	fmt.Printf("%v >= %v (out of %v)\n", attempt, target, pokemon.BaseExperience)
	if attempt >= target {
		fmt.Printf("%s was caught!\n", pokemon.Name)
		pokemonList[pokemon.Name] = pokemon
	} else {
		fmt.Printf("%s escaped!\n", pokemon.Name)
	}
	return nil
}
