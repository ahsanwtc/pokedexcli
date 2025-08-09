package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"

	"github.com/ahsanwtc/pokedexcli/internal/cache"
)

func NewClient(baseURL string, cache cache.Cache) *Client {
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		log.Fatal(err)
	}

	return &Client{
		url: parsedURL,
		httpClient: &http.Client{},
		Config: &Config{},
		cache: cache,
	}
}

func (c *Client) GetLocationAreas(page Page) (*LocationAreas, error)  {
	if c.Config.Next == "" {
		url := c.url.ResolveReference(&url.URL{Path: "location-area"})
		c.Config.Next = url.String()
	}

	url := c.Config.Next
	if page == Previous {
		if c.Config.Previous == "" {
			return nil, fmt.Errorf("EMPTY_PREV")	
		}
		url = c.Config.Previous
	}

	data, ok := c.cache.Get(url)
	var locationAreas LocationAreas
	
	if ok {
		err := toLocationAreas(data, &locationAreas)
		if err != nil {
			return nil, err
		}
		return &locationAreas, nil
	}

	fmt.Println("cache miss", url)

	res, err := c.doRequest("GET", url)
	if err != nil {
		return  nil, err
	}
	defer res.Body.Close()
	
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed with status code: %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return  nil, err
	}

	c.cache.Set(url, body)

	err = json.Unmarshal(body, &locationAreas)
	if err != nil {
		return  nil, err
	}

	c.Config.Next = locationAreas.Next
	c.Config.Previous = locationAreas.Previous

	return  &locationAreas, nil
}

func (c *Client) GetLocationArea(locationAreaName string) (*LocationArea, error)  {
	resourcePath := path.Join("location-area", locationAreaName)
	url := c.url.ResolveReference(&url.URL{Path: resourcePath}).String()

	data, ok := c.cache.Get(url)
	var locationArea LocationArea
	
	if ok {
		err := toLocationArea(data, &locationArea)
		if err != nil {
			return nil, err
		}
		return &locationArea, nil
	}

	res, err := c.doRequest("GET", url)
	if err != nil {
		return  nil, err
	}
	defer res.Body.Close()
	
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed with status code: %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return  nil, err
	}

	c.cache.Set(url, body)

	err = json.Unmarshal(body, &locationArea)
	if err != nil {
		return  nil, err
	}

	return  &locationArea, nil
}

func (c *Client) GetPokemon(pokemonName string) (*Pokemon, error)  {
	resourcePath := path.Join("pokemon", pokemonName)
	url := c.url.ResolveReference(&url.URL{Path: resourcePath}).String()

	data, ok := c.cache.Get(url)
	var pokemon Pokemon
	
	if ok {
		err := toPokemon(data, &pokemon)
		if err != nil {
			return nil, err
		}
		return &pokemon, nil
	}

	res, err := c.doRequest("GET", url)
	if err != nil {
		return  nil, err
	}
	defer res.Body.Close()
	
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed with status code: %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return  nil, err
	}

	c.cache.Set(url, body)

	err = json.Unmarshal(body, &pokemon)
	if err != nil {
		return  nil, err
	}

	return  &pokemon, nil
}

func toLocationAreas(data []byte, locationAreas *LocationAreas) error {
	err := json.Unmarshal(data, locationAreas)
	if err != nil {
		return err
	}
	return nil
}

func toLocationArea(data []byte, locationArea *LocationArea) error {
	err := json.Unmarshal(data, locationArea)
	if err != nil {
		return err
	}
	return nil
}

func toPokemon(data []byte, pokemon *Pokemon) error {
	err := json.Unmarshal(data, pokemon)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) doRequest(method string, url string) (*http.Response, error) {
	var response *http.Response
	var err error
	switch method {
		case http.MethodGet:
			response, err = c.httpClient.Get(url)
		default:
			err = fmt.Errorf("unsupported method")
	}
	return response, err
}
