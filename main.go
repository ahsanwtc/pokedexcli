package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ahsanwtc/pokedexcli/internal/battle"
	"github.com/ahsanwtc/pokedexcli/internal/cache"
	"github.com/ahsanwtc/pokedexcli/internal/pokeapi"
	"github.com/ahsanwtc/pokedexcli/internal/pokedex"
)

type CommandEnv struct {
	client *pokeapi.Client
	pokedex pokedex.Dex
}

type CliCommand struct {
	name        string
	description string
	callback    func(env *CommandEnv, parameters []string) error
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
		callback: func(env *CommandEnv, parameters []string) error {
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

	commands["catch"] = CliCommand{
		name:        "catch",
		description: "Try to catch a pokemon. catch <pokemon_name>",
		callback:    commandCatch,
	}

	commands["list"] = CliCommand{
		name:        "list",
		description: "Show a list of all the caught Pokemons",
		callback:    commandList,
	}

	commands["inspect"] = CliCommand{
		name:        "inspect",
		description: "Inspect a captured Pokemon. inspect <pokemon_name>",
		callback:    commandInspect,
	}

	commandEnv := &CommandEnv{
		client: pokeapi.NewClient("https://pokeapi.co/api/v2/", cache.NewCache(CACHE_TTL)),
		pokedex: pokedex.NewDex(),
	}

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
				err := command.callback(commandEnv, cleaned[1:])
				if err != nil {
					fmt.Println(err)
				}
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

func commandExit(env *CommandEnv, parameters []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandMap(env *CommandEnv, parameters []string) error {
	if env.client == nil {
		return fmt.Errorf("pokeapi client is not configured")
	}

	locationAreas, err := env.client.GetLocationAreas(pokeapi.Next)
	if err != nil {
		return  err
	}

	for _, result := range locationAreas.Results {
		fmt.Println(result.Name)
	}

	return nil
}

func commandMapB(env *CommandEnv, parameters []string) error {
	if env.client == nil {
		return fmt.Errorf("pokeapi client is not configured")
	}

	locationAreas, err := env.client.GetLocationAreas(pokeapi.Previous)
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

func commandExplore(env *CommandEnv, parameters []string) error {
	if env.client == nil {
		return fmt.Errorf("pokeapi client is not configured")
	}

	if len(parameters) == 0 || len(parameters) > 1 {
		return fmt.Errorf("invalid parameter length to search")
	}

	fmt.Printf("Exploring %s...\n", parameters[0])
	locationArea, err := env.client.GetLocationArea(parameters[0])
	if err != nil {
		return  err
	}

	fmt.Println("Found Pokemon:")
	for _, encounter := range locationArea.PokemonEncounters {
		fmt.Println(" - ", encounter.Pokemon.Name)
	}

	return nil
}

func commandCatch(env *CommandEnv, parameters []string) error {
	if env.client == nil {
		return fmt.Errorf("pokeapi client is not configured")
	}

	if env.pokedex == nil {
		return fmt.Errorf("pokedex is not configured")
	}

	if len(parameters) == 0 || len(parameters) > 1 {
		return fmt.Errorf("invalid parameter length to catch a pokemon")
	}

	pokemonName := parameters[0]
	pokemon, err := env.client.GetPokemon(pokemonName)
	if err != nil {
		return  err
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)
	captured := battle.Attack(*pokemon)
	if captured {
		env.pokedex.Add(*pokemon)
		fmt.Printf("%s was caught!\n", pokemon.Name)
	} else {
		fmt.Printf("%s escaped\n", pokemon.Name)
	}

	return nil
}

func commandList(env *CommandEnv, parameters []string) error {
	if env.pokedex == nil {
		return fmt.Errorf("pokedex is not configured")
	}

	env.pokedex.List()

	return nil
}

func commandInspect(env *CommandEnv, parameters []string) error {
	if env.pokedex == nil {
		return fmt.Errorf("pokedex is not configured")
	}

	if len(parameters) == 0 || len(parameters) > 1 {
		return fmt.Errorf("invalid parameter length to inspect a pokemon")
	}

	pokemonName := parameters[0]
	pokemon := env.pokedex.Inspect(pokemonName)
	if pokemon == nil {
		return fmt.Errorf("%s has not been captured yet", pokemonName)
	}

	fmt.Println("Name: ", pokemon.Name)
	fmt.Println("Height: ", pokemon.Height)
	fmt.Println("Weight: ", pokemon.Weight)
	fmt.Println("Stats:")
	for _, stat := range pokemon.Stats {
		fmt.Printf(" -%s: %v\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Println("Types:")
	for _, _type := range pokemon.Types {
		fmt.Println(" -", _type.Type.Name)
	}

	return nil
}
