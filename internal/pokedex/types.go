package pokedex

import (
	"github.com/ahsanwtc/pokedexcli/internal/pokeapi"
)

type Pokedex struct {
	dex map[string]pokeapi.Pokemon
}

type Dex interface {
	Add(pokemon pokeapi.Pokemon)
	List()
	Inspect(pokemonName string) *pokeapi.Pokemon
}