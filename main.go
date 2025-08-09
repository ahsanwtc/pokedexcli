package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ahsanwtc/pokedexcli/internal/cache"
	"github.com/ahsanwtc/pokedexcli/internal/pokeapi"
)

type CliCommand struct {
	name        string
	description string
	callback    func(client *pokeapi.Client, parameters []string) error
}

const CACHE_TTL = 20 * time.Minute

func main()  {
	scanner := bufio.NewScanner(os.Stdin)
	commands := make(map[string]CliCommand)
	
	commands["exit"] = CliCommand{
		name:        "exit",
		description: "Exit the Pokedex",
		callback:    commandExit,
	}

	commands["help"] = CliCommand{
		name:        "help",
		description: "Displays a help message",
		callback: func(client *pokeapi.Client, parameters []string) error {
			fmt.Println("Welcome to the Pokedex!")
			fmt.Println("Usage:")
			fmt.Println("")
			fmt.Println("")
			for key, command := range commands {
				fmt.Printf("%s: %s\n", key, command.description)
			}

			return nil
		},
	}

	commands["map"] = CliCommand{
		name:        "map",
		description: "Fetches the next 20 locations of possible Pokémon encounters.",
		callback:    commandMap,
	}

	commands["mapb"] = CliCommand{
		name:        "mapb",
		description: "Fetches the previous 20 locations of possible Pokémon encounters.",
		callback:    commandMapB,
	}

	commands["explore"] = CliCommand{
		name:        "explore",
		description: "Fetches a list of all the Pokémon located at a location. explore <area_name>",
		callback:    commandExplore,
	}

	client := pokeapi.NewClient("https://pokeapi.co/api/v2/", cache.NewCache(CACHE_TTL))

	for {
		fmt.Print("Pokedex > ")
		if scanner.Scan() {
			input := scanner.Text()
			if input == "" {
				fmt.Println("Provide a command, or use `help` for information!")
				continue
			}

			cleaned := cleanInput(input)
			if command, ok := commands[cleaned[0]]; ok {
				command.callback(client, cleaned[1:])
			} else {
				fmt.Println("Unknown command")
			}
		}
	}
}

func cleanInput(text string) []string {
	var result []string
	if len(text) == 0 {
		return result
	}

	lowered := strings.ToLower(text)
	trimmed := strings.TrimSpace(lowered)
	return strings.Fields(trimmed)
}

func commandExit(client *pokeapi.Client, parameters []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandMap(client *pokeapi.Client, parameters []string) error {
	locationAreas, err := client.GetLocationAreas(pokeapi.Next)
	if err != nil {
		return  err
	}

	for _, result := range locationAreas.Results {
		fmt.Println(result.Name)
	}

	return nil
}

func commandMapB(client *pokeapi.Client, parameters []string) error {	
	locationAreas, err := client.GetLocationAreas(pokeapi.Previous)
	if err != nil {
		if err.Error() == "EMPTY_PREV" {
			fmt.Println("you're on the first page")
		}
		return  err
	}

	for _, result := range locationAreas.Results {
		fmt.Println(result.Name)
	}

	return nil
}

func commandExplore(client *pokeapi.Client, parameters []string) error {
	if len(parameters) == 0 || len(parameters) > 1 {
		return fmt.Errorf("invalid parameter length to search")
	}

	fmt.Printf("Exploring %s...\n", parameters[0])
	locationArea, err := client.GetLocationArea(parameters[0])
	if err != nil {
		if err.Error() == "EMPTY_PREV" {
			fmt.Println("you're on the first page")
		}
		return  err
	}

	fmt.Println("Found Pokemon:")
	for _, encounter := range locationArea.PokemonEncounters {
		fmt.Println(" - ", encounter.Pokemon.Name)
	}

	return nil
}
