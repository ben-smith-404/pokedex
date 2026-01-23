package main

import (
	"fmt"
	"math/rand/v2"
	"os"
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

func regsterCommands() map[string]cliCommand {
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
	registry["inspect"] = cliCommand{
		name:        "inspect",
		description: "inspect a pokemon in your pokedex, you can only inspect pokemon you have previousy caught",
		callback:    commandInspect,
	}
	registry["pokedex"] = cliCommand{
		name:        "pokedex",
		description: "list the pokemn you have caught in your pokedex",
		callback:    commandPokedex,
	}
	return registry
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
	url := baseURL + "/location-area/" + urlExtender
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
		return fmt.Errorf("please include the name of the pokemon you are trying to catch")
	}
	url := baseURL + "/pokemon/" + urlExtender
	pokemon, err := getPokemon(url)
	if err != nil {
		return err
	}
	fmt.Printf("Throwing a Pokeball at %v...\n", urlExtender)
	// the strongest pokemon appear to have a base exp of approx 400 and the weakest appear to be around 38
	// generating an arbitrary catch chance by dividing 3500 by the base experience creates a low % chance for
	// higher tier pokemon and very high for lower tier.
	catchChance := 3500 / pokemon.BaseExperience
	// use a random number between 1-100 to determine if the player rolled a % low enough to catch the pokemon
	// e.g. if the catch chance is 10, the player must roll 10 or less (10% chance)
	if rand.IntN(100) <= catchChance {
		fmt.Printf("%v was caught!\n", urlExtender)
		addToPokedex(urlExtender, pokemon)
	} else {
		fmt.Printf("%v escaped!\n", urlExtender)
	}
	return nil
}

func commandInspect(config *Config, urlExtender string) error {
	pokemon, err := getFromPokedex(urlExtender)
	if err != nil {
		return err
	}
	fmt.Println("Name: " + pokemon.Name)
	fmt.Printf("Height: %v\n", pokemon.Height)
	fmt.Printf("Weight: %v\n", pokemon.Weight)
	fmt.Println("Stats:")
	fmt.Printf("  -hp: %v\n", pokemon.Stats[0].BaseStat)
	fmt.Printf("  -attack: %v\n", pokemon.Stats[1].BaseStat)
	fmt.Printf("  -defence: %v\n", pokemon.Stats[2].BaseStat)
	fmt.Printf("  -special-attack: %v\n", pokemon.Stats[3].BaseStat)
	fmt.Printf("  -special-defence: %v\n", pokemon.Stats[4].BaseStat)
	fmt.Printf("  -speed: %v\n", pokemon.Stats[5].BaseStat)
	fmt.Println("Types:")
	fmt.Printf("  - %v\n", pokemon.Types[0].Type.Name)
	fmt.Printf("  - %v\n", pokemon.Types[1].Type.Name)
	return nil
}

func commandPokedex(config *Config, urlExtender string) error {
	if len(pokedex) == 0 {
		fmt.Println("you haven't caught any pokemon yet")
		return nil
	}
	fmt.Println("Your Pokedex:")
	for name := range pokedex {
		fmt.Printf(" - %v\n", name)
	}
	return nil
}
