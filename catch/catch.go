package catch

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Pokemon struct {
	Name       string `json:"name"`
	Experience int    `json:"base_experience"`
	Tier       int    // You can use tiers based on BE ranges
	Height     int    `json:"height"`
	Weight     int    `json:"weight"`
	Stats      []Stat `json:"stats"`
	Types      []Type `json:"types"`
	// Add more fields as needed
}

type Stat struct {
	BaseStat int `json:"base_stat"`
	Stat     struct {
		Name string `json:"name"`
	} `json:"stat"`
}

type Type struct {
	Type struct {
		Name string `json:"name"`
	} `json:"type"`
}

func CatchPokemon(params []string) {
	// Implement the logic to catch the pokemon
	// For now, we'll just print a message and return true
	// Implement the logic to catch the pokemon
	// For now, we'll just print a message and return true
	if len(params) == 0 {
		fmt.Println("Please provide the name of the Pokemon to catch.")
		return
	}

	pokemonName := strings.Join(params, " ")
	fmt.Println("Attempting to catch", pokemonName)

	// Make API request to get Pokemon data
	res, err := http.Get("https://pokeapi.co/api/v2/pokemon/" + pokemonName)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	// Read response body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Define Pokemon struct

	// Unmarshal JSON into Pokemon struct
	var pokemon Pokemon
	err = json.Unmarshal(body, &pokemon)
	if err != nil {
		log.Fatal(err)
	}

	// Determine the Pokemon's tier based on base experience
	tier := determineTier(pokemon.Experience)

	// Calculate the catch rate multiplier based on tier
	multiplier := calculateMultiplier(tier)

	// Simulate the catch attempt
	success := simulateCatch(multiplier)

	// Print the result
	if success {
		fmt.Printf("Successfully caught %s!\n", pokemon.Name)
		AddToPokedex(pokemon.Name, pokemon.Experience, tier, pokemon.Height, pokemon.Weight)
		fmt.Println("stats : ", pokemon.Name, pokemon.Experience, tier)
	} else {
		fmt.Printf("Failed to catch %s.\n", pokemon.Name)
	}
}

// Function to determine the tier based on base experience
func determineTier(experience int) int {
	switch {
	case experience >= 50 && experience <= 100:
		return 1
	case experience >= 101 && experience <= 250:
		return 2
	case experience >= 251 && experience <= 500:
		return 3
	case experience >= 501:
		return 4
	default:
		return 1
	}
}

// Function to calculate the catch rate multiplier based on tier
func calculateMultiplier(tier int) float64 {
	switch tier {
	case 1:
		return 1.5 // Easy
	case 2:
		return 1.0 // Normal
	case 3:
		return 0.75 // Hard
	case 4:
		return 0.5 // Very Hard
	default:
		return 1.0
	}
}

// Function to simulate the catch attempt
func simulateCatch(multiplier float64) bool {
	r := rand.New(rand.NewSource(time.Now().UnixNano())) // Create a new random source
	successRate := r.Float64()                           // Generate a random number

	return successRate < multiplier
}

var CapturedPokemons = make(map[string]Pokemon)
var Mu sync.Mutex

func AddToPokedex(name string, experience int, tier int, height int, weight int) {
	Mu.Lock()
	defer Mu.Unlock()

	// Check if the Pokémon is already in the Pokedex
	if _, exists := CapturedPokemons[name]; exists {
		fmt.Println("You've already caught", name)
		return
	}

	// Add the Pokémon to the Pokedex
	pokemon := Pokemon{
		Name:       name,
		Experience: experience,
		Tier:       tier,
		Height:     height,
		Weight:     weight,
	}
	CapturedPokemons[name] = pokemon

	fmt.Println("Caught", name)
}
