package main

import (
	"bufio"
	"fmt"
	"os"
)

type config struct {
	previous string
	next     string
}

func (config config) setNext(next string) {
	config.next = next
}

func (config config) setPrevious(previous string) {
	config.previous = previous
}

type cliCommand struct {
	name        string
	description string
	callback    func() error
	config      config
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
		config: config{
			next:     "",
			previous: "",
		},
	}
	registry["mapb"] = cliCommand{
		name:        "mapb",
		description: "gets a list of the previous 20 location areas present in the Pokemon world",
		callback:    commandMapBack,
		config: config{
			next:     "",
			previous: "",
		},
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
			err := command.callback()
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			}
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
