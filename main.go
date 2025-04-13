package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for true {
		fmt.Print("Pokedex > ")

		if scanner.Scan() {
			input := scanner.Text()
			cleaned := cleanInput(input)
			fmt.Println("Your command was:", cleaned[0])
		}
	}
}
