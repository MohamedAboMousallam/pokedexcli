package commands

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/MohamedAboMousallam/pokedexcli/catch"
	"github.com/MohamedAboMousallam/pokedexcli/locations"
	"github.com/MohamedAboMousallam/pokedexcli/pokedex"
	"github.com/MohamedAboMousallam/pokedexcli/printer"
)

var mu sync.Mutex

// HandleCommand handles the user's input command.
func HandleCommand(command string, params []string) {
	switch strings.ToLower(command) {
	case "help":
		printer.PrintHelp()
	case "clear":
		ClearScreen()
	case "locations":
		locations.GetLocations(locations.StartURL)
	case "explore":
		if len(params) > 0 {
			locations.ExploreLocation(strings.Join(params, "-"))
		} else {
			fmt.Println("Please provide a location area name to explore.")
		}
	case "catch":
		catch.CatchPokemon(params)
	case "list":
		pokedex.ListCapturedPokemons()
	case "inspect":
		pokedex.InspectPokemon(params)
	default:
		fmt.Println("Unknown command. Type 'help' for available commands.")
	}
}

func ClearScreen() {
	cmd := exec.Command("cmd", "/c", "cls", "clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func CleanInput(text string) string {
	output := strings.TrimSpace(text)
	output = strings.ToLower(output)
	return output
}

// AddToPokedex adds a Pokemon to the user's Pokedex.
// func AddToPokedex(name string, experience int, tier int, height int, weight int, stats []pokedex.Stat, types []pokedex.Type) {
// 	mu.Lock()
// 	defer mu.Unlock()

// 	pokemon := pokedex.Pokemon{
// 		Name:       name,
// 		Experience: experience,
// 		Tier:       tier,
// 		Height:     height,
// 		Weight:     weight,
// 		Stats:      stats,
// 		Types:      types,
// 	}
// 	pokedex.AddToPokedex(name, experience, tier, height, weight, stats, types)

// 	fmt.Printf("Caught %s! Added to Pokedex.\n", name)
// }

// // ListCapturedPokemons lists all Pokemon in the user's Pokedex.
// func ListCapturedPokemons() {
// 	pokedex.ListCapturedPokemons()
// }

// // InspectPokemon inspects details about a specific Pokemon.
// func InspectPokemon(params []string) {
// 	pokedex.InspectPokemon(params)
// }
