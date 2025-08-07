package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type CliCommand struct {
	name        string
	description string
	callback    func() error
}


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
	callback: func() error {
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

	for {
		fmt.Print("Pokedex > ")
		if scanner.Scan() {
			input := scanner.Text()
			cleaned := cleanInput(input)
			if command, ok := commands[cleaned[0]]; ok {
				command.callback()
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

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}
