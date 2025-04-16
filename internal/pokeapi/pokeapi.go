package pokeapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/logan-waite/bootdev_pokedex/internal/pokecache"
	"io"
	"net/http"
	"time"
)

var BASE_URL = "https://pokeapi.co/api/v2"

// Location Paginator
type paginator struct {
	next     string
	previous string
}

var locationPaginator = paginator{}

// Named Api Resource
type NamedApiResource struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type NamedApiResourceList struct {
	Count    int                `json:"count"`
	Next     string             `json:"next"`
	Previous string             `json:"previous"`
	Results  []NamedApiResource `json:"results"`
}

// Caching
var apiCache = pokecache.NewCache(10 * time.Second)

func cachedGet(url string) ([]byte, error) {
	var data []byte
	if val, exists := apiCache.Get(url); exists {
		fmt.Println("hitting cache")
		data = val
	} else {
		fmt.Println("hitting api")
		res, err := http.Get(url)
		if err != nil {
			return nil, errors.New("error getting location-areas from PokeAPI")
		}
		defer res.Body.Close()

		val, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, errors.New("error reading data from response body")
		}

		data = val
		apiCache.Add(url, data)
	}

	return data, nil
}

func GetLocationAreas(paginate string) ([]NamedApiResource, error) {
	url := BASE_URL + "/location-area/"
	if paginate == "next" && locationPaginator != (paginator{}) {
		url = locationPaginator.next
	} else if paginate == "prev" {
		if locationPaginator != (paginator{}) && locationPaginator.previous == "" {
			return nil, fmt.Errorf("No previous map to return to; use `map` instead")
		} else {
			url = locationPaginator.previous
		}
	}

	data, err := cachedGet(url)
	if err != nil {
		return nil, err
	}

	var locations NamedApiResourceList
	err = json.Unmarshal(data, &locations)
	if err != nil {
		return nil, errors.New("unable to parse location-area JSON")
	}

	locationPaginator.next = locations.Next
	locationPaginator.previous = locations.Previous

	return locations.Results, nil

}
