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

func commandMap() error {
	url := baseURL + "location-area/"
	if registry["map"].config.next != "" {
		url = registry["map"].config.next
	}
	fmt.Println(registry)
	locationAreas, err := getLocationAreas(url)
	if err != nil {
		return err
	}
	registry["map"].config.setNext(locationAreas.Next)
	registry["map"].config.setPrevious(locationAreas.Previous)
	for _, locationArea := range locationAreas.Results {
		fmt.Println(locationArea.Name)
	}
	return nil
}

func commandMapBack() error {
	url := registry["map"].config.previous
	if url == "" {
		return fmt.Errorf("you're on the first page")
	}
	locationAreas, err := getLocationAreas(url)
	if err != nil {
		return err
	}
	registry["map"].config.setNext(locationAreas.Next)
	registry["map"].config.setPrevious(locationAreas.Previous)
	for _, locationArea := range locationAreas.Results {
		fmt.Println(locationArea.Name)
	}
	return nil
}
