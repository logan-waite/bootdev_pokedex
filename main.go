package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	initCommands()

	scanner := bufio.NewScanner(os.Stdin)

	for true {
		fmt.Print("Pokedex > ")

		if scanner.Scan() {
			_input := scanner.Text()
			input := cleanInput(_input)

			if command, ok := commands[input[0]]; ok {
				err := command.callback()
				if err != nil {
					fmt.Printf("Error when calling %s: %v\n", command.name, err)
				}
			} else {
				fmt.Println("Unknown command")
			}
		}
	}
}
