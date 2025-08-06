package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main()  {
  scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")
		if scanner.Scan() {
			input := scanner.Text()
			cleaned := cleanInput(input)
			fmt.Printf("Your command was: %s\n", cleaned[0])
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