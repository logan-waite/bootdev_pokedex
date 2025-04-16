package pokeapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

var BASE_URL = "https://pokeapi.co/api/v2"

// Location Paginator
type paginator struct {
	next     string
	previous string
}

var locationPaginator = paginator{}

func initLocationPaginator() {
	locationPaginator = paginator{
		next:     "",
		previous: "",
	}
}

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

	res, err := http.Get(url)
	if err != nil {
		return nil, errors.New("error getting location-areas from PokeAPI")
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.New("error reading data from response body")
	}

	var locations NamedApiResourceList
	err = json.Unmarshal(data, &locations)
	if err != nil {
		return nil, errors.New("unable to parse location-areas JSON")
	}

	locationPaginator.next = locations.Next
	locationPaginator.previous = locations.Previous

	return locations.Results, nil

}
