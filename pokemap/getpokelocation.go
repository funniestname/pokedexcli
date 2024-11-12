package pokemap

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func GetPokeMap(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	var locations []string

	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&locations); err != nil {
		return nil, fmt.Errorf("error decoding response body: %v", err)
	}

	return locations, nil

}
