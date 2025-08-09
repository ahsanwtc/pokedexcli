package battle

import (
	"github.com/ahsanwtc/pokedexcli/internal/pokeapi"
)

const CAPTURE_THRESHOLD int = 60

type Fight interface {
	Attack(pokemon pokeapi.Pokemon) bool
}