package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/MohamedAboMousallam/pokedexcli/catch"
	"github.com/MohamedAboMousallam/pokedexcli/commands"
	"github.com/MohamedAboMousallam/pokedexcli/locations"
	"github.com/MohamedAboMousallam/pokedexcli/pokedex"
	"github.com/MohamedAboMousallam/pokedexcli/printer"
)

func main() {
	Commands := map[string]interface{}{
		"help":      printer.PrintHelp,
		"clear":     commands.ClearScreen,
		"locations": locations.GetLocationsCommand,
		"explore":   locations.ExploreLocation,
		"catch":     catch.CatchPokemon,
		"add":       catch.AddToPokedex,
		"list":      pokedex.ListCapturedPokemons,
		"inspect":   pokedex.InspectPokemon,
	}

	reader := bufio.NewScanner(os.Stdin)
	printer.PrintPrompt()
	for reader.Scan() {
		text := commands.CleanInput(reader.Text())

		// Split the user input into command and parameters
		splitResult := strings.Fields(text)
		if len(splitResult) > 0 {
			command := splitResult[0]
			var params []string

			// Extract parameters if available
			if len(splitResult) > 1 {
				params = splitResult[1:]
			}

			// Check if the command exists
			if cmdFunc, exists := Commands[command]; exists {
				switch f := cmdFunc.(type) {
				case func():
					// Call functions without arguments
					f()
				case func(string):
					// Call functions that expect a string argument
					if len(splitResult) > 1 {
						f(strings.Join(splitResult[1:], " "))
					} else {
						fmt.Println("Please provide a valid argument for the", command, "command.")
					}
				case func([]string):
					f(params)
					if len(splitResult) > 1 {
						f(splitResult[1:])
					} else {
						fmt.Println("Please provide a valid argument for the", command, "command.")
					}
				default:
					fmt.Println("Invalid command type")
				}
			} else if strings.EqualFold("exit", command) {
				return
			} else {
				commands.HandleCommand(command, params)
			}
		} else {
			fmt.Println("Invalid input. Please enter a command.")
		}

		printer.PrintPrompt()
	}

	// Print an additional line if we encountered an EOF character
	fmt.Println()
}
