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
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 299 {
		return nil, fmt.Errorf("Error: Status code %v does not indicate success", res.StatusCode)
	}
	if err != nil {
		return nil, err
	}
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

const baseURL = "https://pokeapi.co/api/v2/"

func commandMap(config *Config) error {
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

func commandMapBack(config *Config) error {
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
