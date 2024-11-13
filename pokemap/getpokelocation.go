package pokemap

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type PokeLocation struct {
	Count    int     `json:"count"`
	Next     string  `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	} `json:"results"`
}

func GetPokeMap(url string) (PokeLocation, error) {
	resp, err := http.Get(url)
	if err != nil {
		return PokeLocation{}, err
	}

	var locations PokeLocation

	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&locations); err != nil {
		return PokeLocation{}, fmt.Errorf("error decoding response body: %v", err)
	}

	return locations, nil

}
