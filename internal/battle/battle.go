package battle

import (
	"math/rand"

	"github.com/ahsanwtc/pokedexcli/internal/pokeapi"
)

func Attack(pokemon pokeapi.Pokemon) bool {
	exp := pokemon.BaseExperience
	pokemonAttack := exp + 100
	chance := rand.Intn(pokemonAttack)
	return chance > CAPTURE_THRESHOLD
}