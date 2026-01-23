package main

import (
	"bufio"
	"fmt"
	pokecache "github/ben-smith-404/pokedexcli/internal/pokecache"
	"os"
	"time"
)

var cache = pokecache.NewCache(5 * time.Second)

var pokedex = map[string]Pokemon{}

func main() {
	registry := regsterCommands()

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
				fmt.Printf("%v\n", err)
			}
		}
	}
}
