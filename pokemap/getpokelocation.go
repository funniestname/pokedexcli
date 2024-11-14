package pokemap

import (
	"encoding/json"
	"fmt"
	"math/rand"
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

func GetCatch(pokemon string) (PokemonData, bool, error) {
	resp, err := http.Get("https://pokeapi.co/api/v2/pokemon/" + pokemon)
	if err != nil {
		return PokemonData{}, false, err
	}

	var pokemonData PokemonData

	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&pokemonData); err != nil {
		return PokemonData{}, false, err
	}

	fmt.Println("Throwing a Pokeball at " + pokemon + "...")

	catchRate := 1.0 - float64(pokemonData.BaseExperience)/1000.0
	if rand.Float64() < catchRate {
		return pokemonData, true, nil
	}
	return PokemonData{}, false, nil
}
