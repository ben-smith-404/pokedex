package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func cleanInput(text string) []string {
	var strs []string
	split_strs := strings.Split(text, " ")
	for _, str := range split_strs {
		if len(str) > 0 {
			strs = append(strs, strings.ToLower((str)))
		}
	}
	return strs
}

func getDataFromAPI(url string) ([]byte, error) {
	data, ok := cache.Get(url)
	if ok {
		return data, nil
	}
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 299 {
		return nil, fmt.Errorf("Error: Status code %v does not indicate success. Your selection may not exist in the pokemon world.", res.StatusCode)
	}
	if err != nil {
		return nil, err
	}
	cache.Add(url, body)
	return body, nil
}

type LocationArea struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type LocationAreasResponse struct {
	Next     string         `json:"next"`
	Previous string         `json:"previous"`
	Results  []LocationArea `json:"results"`
}

func getLocationAreas(url string) (LocationAreasResponse, error) {
	var locationAreas LocationAreasResponse
	apiData, err := getDataFromAPI(url)
	if err != nil {
		return locationAreas, err
	}
	err = json.Unmarshal(apiData, &locationAreas)
	if err != nil {
		return locationAreas, err
	}
	return locationAreas, nil
}

// I did this one using JSON to go which is why it's much more comple than the above set of structs
type LocationAreaResponse struct {
	Location struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Name              string `json:"name"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

func getLocationArea(url string) (LocationAreaResponse, error) {
	var locationArea LocationAreaResponse
	apiData, err := getDataFromAPI(url)
	if err != nil {
		return locationArea, err
	}
	err = json.Unmarshal(apiData, &locationArea)
	if err != nil {
		return locationArea, err
	}
	return locationArea, nil
}

type Pokemon struct {
	Abilities []struct {
		Ability struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"ability"`
		IsHidden bool `json:"is_hidden"`
		Slot     int  `json:"slot"`
	} `json:"abilities"`
	BaseExperience int `json:"base_experience"`
	Forms          []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"forms"`
	Height                 int    `json:"height"`
	ID                     int    `json:"id"`
	IsDefault              bool   `json:"is_default"`
	LocationAreaEncounters string `json:"location_area_encounters"`
	Name                   string `json:"name"`
	Order                  int    `json:"order"`
	PastAbilities          []struct {
		Abilities []struct {
			Ability  interface{} `json:"ability"`
			IsHidden bool        `json:"is_hidden"`
			Slot     int         `json:"slot"`
		} `json:"abilities"`
		Generation struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"generation"`
	} `json:"past_abilities"`
	PastTypes []interface{} `json:"past_types"`
	Species   struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"species"`
	Stats []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
	Weight int `json:"weight"`
}

func getPokemon(url string) (Pokemon, error) {
	var pokemon Pokemon
	apiData, err := getDataFromAPI(url)
	if err != nil {
		return pokemon, err
	}
	err = json.Unmarshal(apiData, &pokemon)
	if err != nil {
		return pokemon, err
	}
	return pokemon, nil
}

var pokedex = map[string]Pokemon{}

func addToPokedex(pokemon string, pokemonData Pokemon) {
	pokedex[pokemon] = pokemonData
}
