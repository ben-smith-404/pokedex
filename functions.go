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
		return nil, fmt.Errorf("Error: Status code %v does not indicate success", res.StatusCode)
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

type LocationAreaResponse struct {
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	GameIndex int `json:"game_index"`
	ID        int `json:"id"`
	Location  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Name  string `json:"name"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance          int           `json:"chance"`
				ConditionValues []interface{} `json:"condition_values"`
				MaxLevel        int           `json:"max_level"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
				MinLevel int `json:"min_level"`
			} `json:"encounter_details"`
			MaxChance int `json:"max_chance"`
			Version   struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
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
