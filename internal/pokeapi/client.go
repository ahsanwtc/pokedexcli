package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

func NewClient(baseURL string) *Client {
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		log.Fatal(err)
	}

	return &Client{
		url: parsedURL,
		httpClient: &http.Client{},
		Config: &Config{},
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

	var locationAreas LocationAreas
	err = json.Unmarshal(body, &locationAreas)
	if err != nil {
		return  nil, err
	}

	c.Config.Next = locationAreas.Next
	c.Config.Previous = locationAreas.Previous

	return  &locationAreas, nil
}
