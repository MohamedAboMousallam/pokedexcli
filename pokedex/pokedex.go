package pokedex

import (
	"fmt"
	"strings"
	"github.com/MohamedAboMousallam/pokedexcli/catch"
)


// ListCapturedPokemons lists all Pokémon in the Pokedex.
func ListCapturedPokemons() {
	catch.Mu.Lock()
	defer catch.Mu.Unlock()

	if len(catch.CapturedPokemons) == 0 {
		fmt.Println("Oh no your pokedex is empty, go and catch some pokemons!.")
		return
	}

	fmt.Println("Pokémons in Your Pokedex:")
	for _, pokemon := range catch.CapturedPokemons {
		fmt.Printf("%s (Experience: %d, Tier: %d)\n", pokemon.Name, pokemon.Experience, pokemon.Tier)
	}
}

func InspectPokemon(params []string) {
	if len(params) == 0 {
		fmt.Println("Please provide the name of the Pokemon to inspect.")
		return
	}

	pokemonName := strings.Join(params, " ")

	// Check if the Pokémon is in the Pokedex
	catch.Mu.Lock()
	pokemon, exists := catch.CapturedPokemons[pokemonName]
	catch.Mu.Unlock()

	if !exists {
		fmt.Printf("You haven't caught %s yet.\n", pokemonName)
		return
	}

	// Print Pokemon details
	fmt.Printf("Details for %s:\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)
	fmt.Println("Stats:")
	for _, stat := range pokemon.Stats {
		fmt.Printf("  %s: %d\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Print("Types: ")
	for i, pokemonType := range pokemon.Types {
		fmt.Print(pokemonType.Type.Name)
		if i < len(pokemon.Types)-1 {
			fmt.Print(", ")
		}
	}
	fmt.Println()
}
