package pokemap

import (
	"encoding/json"
	"fmt"
	"net/http"
)

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

func GetPokemon(url string) (PokemonEncounters, error) {
	resp, err := http.Get(url)
	if err != nil {
		return PokemonEncounters{}, err
	}

	var locationData LocationData

	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&locationData); err != nil {
		return PokemonEncounters{}, err
	}
	return locationData.PokemonEncounters, nil
}
