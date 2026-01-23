package main

import (
	"bufio"
	"fmt"
	pokecache "github/ben-smith-404/pokedexcli/internal"
	"os"
	"time"
)

type Config struct {
	previous string
	next     string
}

var config = Config{
	next:     "",
	previous: "",
}

type cliCommand struct {
	name        string
	description string
	callback    func(config *Config, urlExtender string) error
}

var registry = map[string]cliCommand{}

var cache = pokecache.NewCache(5 * time.Second)

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
	registry["explore"] = cliCommand{
		name:        "explore",
		description: "returns the list of pokemon in the specified location area",
		callback:    commandExplore,
	}
	registry["catch"] = cliCommand{
		name:        "catch",
		description: "attempt to catch a pokemon by typing 'catch' and the name of the pokemon you want to catch",
		callback:    commandCatch,
	}

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := scanner.Text()
		words := cleanInput(input)
		urlExtender := ""
		if len(words) > 1 {
			urlExtender = words[1]
		}
		command, ok := registry[words[0]]
		if !ok {
			fmt.Println("Unknown command")
		} else {
			err := command.callback(&config, urlExtender)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			}
		}
	}
}

func commandExit(config *Config, urlExtender string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	defer os.Exit(0)
	return nil
}

func commandHelp(config *Config, urlExtender string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	for _, command := range registry {
		fmt.Printf("%v: %v\n", command.name, command.description)
	}
	return nil
}

const baseURL = "https://pokeapi.co/api/v2/"

func commandMap(config *Config, urlExtender string) error {
	var url string
	if config.next == "" {
		url = baseURL + "location-area/"
	} else {
		url = config.next
	}
	locationAreas, err := getLocationAreas(url)
	if err != nil {
		return err
	}
	config.next = locationAreas.Next
	config.previous = locationAreas.Previous
	for _, locationArea := range locationAreas.Results {
		fmt.Println(locationArea.Name)
	}
	return nil
}

func commandMapBack(config *Config, urlExtender string) error {
	var url string
	if config.previous == "" {
		return fmt.Errorf("you're on the first page")
	} else {
		url = config.previous
	}
	locationAreas, err := getLocationAreas(url)
	if err != nil {
		return err
	}
	config.next = locationAreas.Next
	config.previous = locationAreas.Previous
	for _, locationArea := range locationAreas.Results {
		fmt.Println(locationArea.Name)
	}
	return nil
}

func commandExplore(config *Config, urlExtender string) error {
	if urlExtender == "" {
		return fmt.Errorf("you need to specify the location you want to explore")
	}
	url := baseURL + "/location-area/" + urlExtender + "/"
	locationArea, err := getLocationArea(url)
	if err != nil {
		return err
	}
	for _, pokemon := range locationArea.PokemonEncounters {
		fmt.Println(pokemon.Pokemon.Name)
	}
	return nil
}

func commandCatch(config *Config, urlExtender string) error {
	if urlExtender == "" {
		return fmt.Errorf("please specify the name of the pokemon you are trying to catch")
	}
	return nil
}
