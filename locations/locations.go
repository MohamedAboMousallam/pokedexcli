package locations

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/MohamedAboMousallam/pokedexcli/pokecache"
	"github.com/manifoldco/promptui"
)

var cache *pokecache.Cache

func init() {
	cache = pokecache.GetDefaultCache()
}

var StartURL = "https://pokeapi.co/api/v2/location/"

type APIResponse struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
}

func MakeRequest(url string) (*APIResponse, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var apiResponse APIResponse
	err = json.NewDecoder(response.Body).Decode(&apiResponse)
	if err != nil {
		return nil, err
	}

	return &apiResponse, nil
}

func GetLocationsFromAPI(url string) ([]string, *APIResponse, error) {

	if cachedData, ok := cache.Get(url); ok {
		var locations []string
		err := json.Unmarshal(cachedData, &locations)
		if err != nil {
			return nil, nil, err
		}
		fmt.Println("Using cached data", cachedData, locations)
		return locations, nil, nil // No need to call MakeRequest here
	}

	res, err := http.Get(url)
	if err != nil {
		return nil, nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, nil, err
	}

	if res.StatusCode > 299 {
		return nil, nil, fmt.Errorf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
	}

	var locations []string

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, nil, err
	}

	if results, ok := result["results"].([]interface{}); ok {
		for _, location := range results {
			if locMap, ok := location.(map[string]interface{}); ok {
				if name, ok := locMap["name"].(string); ok {
					locations = append(locations, name)
				}
			}
		}
	}

	dataToCache, err := json.Marshal(locations)
	if err != nil {
		return nil, nil, err
	}
	cache.Add(url, dataToCache)
	fmt.Println("Added cached data as ")

	apiResponse, err := MakeRequest(url)
	if err != nil {
		return nil, nil, err
	}

	return locations, apiResponse, nil
}

func DisplayLocations(locations []string) {
	fmt.Println("Pokemon Locations:")
	for _, location := range locations {
		fmt.Println(location)
	}
}

func HandleLocationInput(apiResponse *APIResponse) error {
	for {
		prompt := promptui.Prompt{
			Label: "Welcome to Pokedex Locations Type 'next' or 'previous' or 'exit' to return to the main menu",
		}

		result, err := prompt.Run()
		if err != nil {
			return err
		}
		splitResult := strings.Fields(result)
		if len(splitResult) > 0 {
			switch strings.ToLower(splitResult[0]) {
			case "next":
				if apiResponse.Next == "" {
					fmt.Println("No next URL available.")
				} else {
					StartURL = apiResponse.Next
					GetLocations(StartURL)
				}
			case "previous":
				if apiResponse.Previous == "" {
					fmt.Println("No previous URL available.")
				} else {
					StartURL = apiResponse.Previous
					GetLocations(StartURL)
				}
			case "explore":
				if len(splitResult) < 2 {
					fmt.Println("Please provide a location area name to explore.")
				} else {
					areaName := strings.Join(splitResult[1:], "-")
					ExploreLocation(areaName)
				}
			case "exit":
				fmt.Println("Exiting the locator map. Godspeed champ!")
				os.Exit(0)
			default:
				fmt.Println("Invalid input. Type 'next' or 'previous'.")
			}
		}
	}
}

// The main GetLocations function
func GetLocations(url string) {
	if url == "" {
		url = StartURL
	}
	locations, apiResponse, err := GetLocationsFromAPI(StartURL)
	if err != nil {
		log.Fatal(err)
	}

	DisplayLocations(locations)
	HandleLocationInput(apiResponse)
}

func ExploreLocation(areaName string) {
	// Construct the URL for the location area
	//areaName = "canalave-city-area"

	fmt.Print("areaname you're looking for is: ", areaName)
	locationURL := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s/", areaName)

	// Make the request to the PokeAPI
	response, err := http.Get(locationURL)
	if err != nil {
		fmt.Println("Error fetching data:", err)
		return
	}
	defer response.Body.Close()
	// Check the response status code
	if response.StatusCode != http.StatusOK {
		fmt.Println("API returned an error:", response.Status)
		return
	}

	// Print the raw response for debugging
	// body, _ := io.ReadAll(response.Body)
	// fmt.Println("Raw Response Body:", string(body))

	// Parse the response
	var locationDetails map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&locationDetails)
	if err != nil {
		fmt.Println("Error decoding response:", err)
		return
	}

	// Extract and display the Pokemon names
	if pokemonList, ok := locationDetails["pokemon_encounters"].([]interface{}); ok {
		fmt.Println("Pokemon in", areaName, ":")
		for _, pokemon := range pokemonList {
			if pokemonDetails, ok := pokemon.(map[string]interface{}); ok {
				if pokemonName, ok := pokemonDetails["pokemon"].(map[string]interface{})["name"].(string); ok {
					fmt.Println("-", pokemonName)
				}
			}
		}
	} else {
		fmt.Println("No Pokemon encounters found in", areaName)
	}
}

func GetLocationsCommand(params []string) {
	// Ignore any command-line arguments
	GetLocations(StartURL)
}
