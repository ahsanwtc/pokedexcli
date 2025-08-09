package pokeapi

import (
	"net/http"
	"net/url"

	"github.com/ahsanwtc/pokedexcli/internal/cache"
)

type Config struct {
	Previous string
	Next     string
}

type Client struct {
	url *url.URL
	httpClient *http.Client
	Config *Config
	cache cache.Cache
}

type LocationAreas struct {
	Next     string `json:"next"`
	Previous string    `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
	} `json:"results"`
}

type Page string
const (
	Next Page = "next"
	Previous Page = "previous"
)