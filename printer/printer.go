// printer/printer.go
package printer

import (
	"fmt"
)

// PrintPrompt prints the command prompt.
var cliName string = "Pokedex"
func PrintPrompt() {
	fmt.Print(cliName, "> ")
}

// PrintUnknown prints a message for an unknown command.
func PrintUnknown(text string) {
	fmt.Println(text, ": Unknown command")
}

// PrintHelp prints the help message.
func PrintHelp() {
	fmt.Printf(
		"Welcome to %v! These are the available commands:\n",
		cliName,
	)
	fmt.Println("help      - Show available commands")
	fmt.Println("clear     - Clear the terminal screen")
	fmt.Println("locations - Get Pokemon locations")
	fmt.Println("exit      - Closes your connection to pokedex")
	fmt.Println("explore   - Explore a location area")
	fmt.Println("catch     - Attempt to catch a Pokemon")
	fmt.Println("list      - List captured Pokemon in Pokedex")
	fmt.Println("inspect   - Inspect details of a captured Pokemon")
	// Add more commands as needed
}

// PrintMessage prints a general message.
func PrintMessage(message string) {
	fmt.Println(message)
}
