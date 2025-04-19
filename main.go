package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	initPokedex()

	scanner := bufio.NewScanner(os.Stdin)

	for true {
		fmt.Print("Pokedex > ")

		if scanner.Scan() {
			_input := scanner.Text()
			input := cleanInput(_input)

			cmd, arg := "", ""
			if len(input) > 0 {
				cmd = input[0]
			}
			if len(input) > 1 {
				arg = input[1]
			}

			if cmd == "" {
				continue
			} else if command, ok := commands[cmd]; ok {
				err := command.callback(arg)
				if err != nil {
					fmt.Printf("Error when calling %s: %v\n", command.name, err)
				}
			} else {
				fmt.Println("Unknown command")
			}
		}
	}
}
