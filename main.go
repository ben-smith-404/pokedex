package main

import (
	"bufio"
	"fmt"
	"os"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

var registry = map[string]cliCommand{}

func main() {
	// register the commands we need
	registry["exit"] = cliCommand{
		name:        "exit",
		description: "exit the application",
		callback:    commandExit,
	}
	registry["help"] = cliCommand{
		name:        "help",
		description: "explain the commands you can use to interact with the Pokedex",
		callback:    commandHelp,
	}

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := scanner.Text()
		words := cleanInput(input)
		command, ok := registry[words[0]]
		if !ok {
			fmt.Println("Unknown command")
		} else {
			command.callback()
		}
	}
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	defer os.Exit(0)
	return nil
}

func commandHelp() error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	for _, command := range registry {
		fmt.Printf("%v: %v\n", command.name, command.description)
	}
	return nil
}
