package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

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

	res, err := c.httpClient.Get(url)
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

func toLocationAreas(data []byte, locationAreas *LocationAreas) error {
	err := json.Unmarshal(data, locationAreas)
	if err != nil {
		return err
	}
	return nil
}
