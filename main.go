package main

import (
	"bufio"
	"fmt"
	"os"
)

type Config struct {
	previous string
	next     string
}

type cliCommand struct {
	name        string
	description string
	callback    func(config *Config) error
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
	registry["map"] = cliCommand{
		name:        "map",
		description: "gets a list of the next 20 location areas present in the Pokemons world",
		callback:    commandMap,
	}
	registry["mapb"] = cliCommand{
		name:        "mapb",
		description: "gets a list of the previous 20 location areas present in the Pokemon world",
		callback:    commandMapBack,
	}

	var config = Config{
		next:     "",
		previous: "",
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
			err := command.callback(&config)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			}
		}
	}
}

func commandExit(config *Config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	defer os.Exit(0)
	return nil
}

func commandHelp(config *Config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	for _, command := range registry {
		fmt.Printf("%v: %v\n", command.name, command.description)
	}
	return nil
}
