package pokedex

import (
	"fmt"

	"github.com/ahsanwtc/pokedexcli/internal/pokeapi"
)

func NewDex() Dex {
	pokedex := &Pokedex{
		dex: make(map[string]pokeapi.Pokemon),
	}

	return pokedex
}

func (p* Pokedex) Add(pokemon pokeapi.Pokemon)  {
	p.dex[pokemon.Name] = pokemon
}

func (p* Pokedex) List()  {
	fmt.Println("List of caught Pokemons")
	for _, pokemon := range p.dex {
		fmt.Println(" - ", pokemon.Name)
	}
}

func (p* Pokedex) Inspect(pokemonName string) *pokeapi.Pokemon {
	pokemon, ok := p.dex[pokemonName]
	if !ok {
		return nil
	}

	return &pokemon
}