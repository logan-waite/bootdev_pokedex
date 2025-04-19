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

// Caching
var apiCache = pokecache.NewCache(10 * time.Second)

func cachedGet(url string) ([]byte, error) {
	var data []byte
	if val, exists := apiCache.Get(url); exists {
		data = val
	} else {
		res, err := http.Get(url)
		if err != nil {
			return nil, errors.New("error getting response from PokeAPI")
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

/*** GetLocationAreas ***/
// Types
type paginator struct {
	next     string
	previous string
}

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

var locationPaginator = paginator{}

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

/*** GetLocationAreaData ***/
// Types
type Encounter struct {
	MinLevel        int              `json:"min_level"`
	MaxLevel        int              `json:"max_level"`
	ConditionValues NamedApiResource `json:"conditionValues"`
	Change          int              `json:"change"`
	Method          NamedApiResource `json:"method"`
}

type VersionEncounterDetail struct {
	Version          NamedApiResource `json:"version"`
	MaxChance        int              `json:"max_chance"`
	EncounterDetails []Encounter      `json:"encounter_details"`
}

type EncounterMethodRate struct {
	EncounterMethod NamedApiResource         `json:"encounter_method"`
	VersionDetails  []VersionEncounterDetail `json:"version_details"`
}

type Name struct {
	Name     string           `json:"Name"`
	Language NamedApiResource `json:"Language"`
}

type PokemonEncounter struct {
	Pokemon        NamedApiResource         `json:"pokemon"`
	VersionDetails []VersionEncounterDetail `json:"version_details"`
}

type LocationArea struct {
	Id                   int                   `json:"id"`
	Name                 string                `json:"name"`
	GameIndex            int                   `json:"game_index"`
	EncounterMethodRates []EncounterMethodRate `json:"encounter_method_rates"`
	Location             NamedApiResource      `json:"location"`
	Names                []Name                `json:"names"`
	PokemonEncounters    []PokemonEncounter    `json:"pokemon_encounters"`
}

func GetLocationAreaData(location string) (LocationArea, error) {
	url := BASE_URL + "/location-area/" + location

	data, err := cachedGet(url)
	if err != nil {
		return LocationArea{}, err
	}

	var result LocationArea
	err = json.Unmarshal(data, &result)
	if err != nil {
		return LocationArea{}, errors.New("unable to parse location JSON")
	}

	return result, nil

}

/*** GetPokemon ***/
// Types
type PokemonAbility struct {
	IsHidden bool             `json:"is_hidden"`
	Slot     int              `json:"slot"`
	Ability  NamedApiResource `json:"ability"`
}
type VersionGameIndex struct {
	GameIndex int              `json:"game_index"`
	Version   NamedApiResource `json:"version"`
}
type PokemonHeldItem struct {
	Item           NamedApiResource `json:"item"`
	VersionDetails []struct {
		Rarity  int              `json:"rarity"`
		Version NamedApiResource `json:"version"`
	} `json:"version_details"`
}
type PokemonMove struct {
	Move                NamedApiResource `json:"move"`
	VersionGroupDetails []struct {
		LevelLearnedAt int `json:"level_learned_at"`
		VersionGroup   struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"version_group"`
		MoveLearnMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"move_learn_method"`
		Order int `json:"order"`
	} `json:"version_group_details"`
}
type PokemonSprites struct {
	BackDefault      string `json:"back_default"`
	BackFemale       any    `json:"back_female"`
	BackShiny        string `json:"back_shiny"`
	BackShinyFemale  any    `json:"back_shiny_female"`
	FrontDefault     string `json:"front_default"`
	FrontFemale      any    `json:"front_female"`
	FrontShiny       string `json:"front_shiny"`
	FrontShinyFemale any    `json:"front_shiny_female"`
	Other            struct {
		DreamWorld struct {
			FrontDefault string `json:"front_default"`
			FrontFemale  any    `json:"front_female"`
		} `json:"dream_world"`
		Home struct {
			FrontDefault     string `json:"front_default"`
			FrontFemale      any    `json:"front_female"`
			FrontShiny       string `json:"front_shiny"`
			FrontShinyFemale any    `json:"front_shiny_female"`
		} `json:"home"`
		OfficialArtwork struct {
			FrontDefault string `json:"front_default"`
			FrontShiny   string `json:"front_shiny"`
		} `json:"official-artwork"`
		Showdown struct {
			BackDefault      string `json:"back_default"`
			BackFemale       any    `json:"back_female"`
			BackShiny        string `json:"back_shiny"`
			BackShinyFemale  any    `json:"back_shiny_female"`
			FrontDefault     string `json:"front_default"`
			FrontFemale      any    `json:"front_female"`
			FrontShiny       string `json:"front_shiny"`
			FrontShinyFemale any    `json:"front_shiny_female"`
		} `json:"showdown"`
	} `json:"other"`
	Versions struct {
		GenerationI struct {
			RedBlue struct {
				BackDefault  string `json:"back_default"`
				BackGray     string `json:"back_gray"`
				FrontDefault string `json:"front_default"`
				FrontGray    string `json:"front_gray"`
			} `json:"red-blue"`
			Yellow struct {
				BackDefault  string `json:"back_default"`
				BackGray     string `json:"back_gray"`
				FrontDefault string `json:"front_default"`
				FrontGray    string `json:"front_gray"`
			} `json:"yellow"`
		} `json:"generation-i"`
		GenerationIi struct {
			Crystal struct {
				BackDefault  string `json:"back_default"`
				BackShiny    string `json:"back_shiny"`
				FrontDefault string `json:"front_default"`
				FrontShiny   string `json:"front_shiny"`
			} `json:"crystal"`
			Gold struct {
				BackDefault  string `json:"back_default"`
				BackShiny    string `json:"back_shiny"`
				FrontDefault string `json:"front_default"`
				FrontShiny   string `json:"front_shiny"`
			} `json:"gold"`
			Silver struct {
				BackDefault  string `json:"back_default"`
				BackShiny    string `json:"back_shiny"`
				FrontDefault string `json:"front_default"`
				FrontShiny   string `json:"front_shiny"`
			} `json:"silver"`
		} `json:"generation-ii"`
		GenerationIii struct {
			Emerald struct {
				FrontDefault string `json:"front_default"`
				FrontShiny   string `json:"front_shiny"`
			} `json:"emerald"`
			FireredLeafgreen struct {
				BackDefault  string `json:"back_default"`
				BackShiny    string `json:"back_shiny"`
				FrontDefault string `json:"front_default"`
				FrontShiny   string `json:"front_shiny"`
			} `json:"firered-leafgreen"`
			RubySapphire struct {
				BackDefault  string `json:"back_default"`
				BackShiny    string `json:"back_shiny"`
				FrontDefault string `json:"front_default"`
				FrontShiny   string `json:"front_shiny"`
			} `json:"ruby-sapphire"`
		} `json:"generation-iii"`
		GenerationIv struct {
			DiamondPearl struct {
				BackDefault      string `json:"back_default"`
				BackFemale       any    `json:"back_female"`
				BackShiny        string `json:"back_shiny"`
				BackShinyFemale  any    `json:"back_shiny_female"`
				FrontDefault     string `json:"front_default"`
				FrontFemale      any    `json:"front_female"`
				FrontShiny       string `json:"front_shiny"`
				FrontShinyFemale any    `json:"front_shiny_female"`
			} `json:"diamond-pearl"`
			HeartgoldSoulsilver struct {
				BackDefault      string `json:"back_default"`
				BackFemale       any    `json:"back_female"`
				BackShiny        string `json:"back_shiny"`
				BackShinyFemale  any    `json:"back_shiny_female"`
				FrontDefault     string `json:"front_default"`
				FrontFemale      any    `json:"front_female"`
				FrontShiny       string `json:"front_shiny"`
				FrontShinyFemale any    `json:"front_shiny_female"`
			} `json:"heartgold-soulsilver"`
			Platinum struct {
				BackDefault      string `json:"back_default"`
				BackFemale       any    `json:"back_female"`
				BackShiny        string `json:"back_shiny"`
				BackShinyFemale  any    `json:"back_shiny_female"`
				FrontDefault     string `json:"front_default"`
				FrontFemale      any    `json:"front_female"`
				FrontShiny       string `json:"front_shiny"`
				FrontShinyFemale any    `json:"front_shiny_female"`
			} `json:"platinum"`
		} `json:"generation-iv"`
		GenerationV struct {
			BlackWhite struct {
				Animated struct {
					BackDefault      string `json:"back_default"`
					BackFemale       any    `json:"back_female"`
					BackShiny        string `json:"back_shiny"`
					BackShinyFemale  any    `json:"back_shiny_female"`
					FrontDefault     string `json:"front_default"`
					FrontFemale      any    `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale any    `json:"front_shiny_female"`
				} `json:"animated"`
				BackDefault      string `json:"back_default"`
				BackFemale       any    `json:"back_female"`
				BackShiny        string `json:"back_shiny"`
				BackShinyFemale  any    `json:"back_shiny_female"`
				FrontDefault     string `json:"front_default"`
				FrontFemale      any    `json:"front_female"`
				FrontShiny       string `json:"front_shiny"`
				FrontShinyFemale any    `json:"front_shiny_female"`
			} `json:"black-white"`
		} `json:"generation-v"`
		GenerationVi struct {
			OmegarubyAlphasapphire struct {
				FrontDefault     string `json:"front_default"`
				FrontFemale      any    `json:"front_female"`
				FrontShiny       string `json:"front_shiny"`
				FrontShinyFemale any    `json:"front_shiny_female"`
			} `json:"omegaruby-alphasapphire"`
			XY struct {
				FrontDefault     string `json:"front_default"`
				FrontFemale      any    `json:"front_female"`
				FrontShiny       string `json:"front_shiny"`
				FrontShinyFemale any    `json:"front_shiny_female"`
			} `json:"x-y"`
		} `json:"generation-vi"`
		GenerationVii struct {
			Icons struct {
				FrontDefault string `json:"front_default"`
				FrontFemale  any    `json:"front_female"`
			} `json:"icons"`
			UltraSunUltraMoon struct {
				FrontDefault     string `json:"front_default"`
				FrontFemale      any    `json:"front_female"`
				FrontShiny       string `json:"front_shiny"`
				FrontShinyFemale any    `json:"front_shiny_female"`
			} `json:"ultra-sun-ultra-moon"`
		} `json:"generation-vii"`
		GenerationViii struct {
			Icons struct {
				FrontDefault string `json:"front_default"`
				FrontFemale  any    `json:"front_female"`
			} `json:"icons"`
		} `json:"generation-viii"`
	} `json:"versions"`
}
type PokemonCries struct {
	Latest string `json:"latest"`
	Legacy string `json:"legacy"`
}
type PokemonStat struct {
	BaseStat int              `json:"base_stat"`
	Effort   int              `json:"effort"`
	Stat     NamedApiResource `json:"stat"`
}
type PokemonType struct {
	Slot int              `json:"slot"`
	Type NamedApiResource `json:"type"`
}
type PokemonPastType struct {
	Generation NamedApiResource `json:"generation"`
	Types      []PokemonType    `json:"types"`
}
type PokemonAbilityPast struct {
	Generation NamedApiResource `json:"generation"`
	Abilities  []PokemonAbility `json:"abilities"`
}
type Pokemon struct {
	ID                     int                  `json:"id"`
	Name                   string               `json:"name"`
	BaseExperience         int                  `json:"base_experience"`
	Height                 int                  `json:"height"`
	IsDefault              bool                 `json:"is_default"`
	Order                  int                  `json:"order"`
	Weight                 int                  `json:"weight"`
	Abilities              []PokemonAbility     `json:"abilities"`
	Forms                  []NamedApiResource   `json:"forms"`
	GameIndices            []VersionGameIndex   `json:"game_indices"`
	HeldItems              []PokemonHeldItem    `json:"held_items"`
	LocationAreaEncounters string               `json:"location_area_encounters"`
	Moves                  []PokemonMove        `json:"moves"`
	Species                NamedApiResource     `json:"species"`
	Sprites                PokemonSprites       `json:"sprites"`
	Cries                  PokemonCries         `json:"cries"`
	Stats                  []PokemonStat        `json:"stats"`
	Types                  []PokemonType        `json:"types"`
	PastTypes              []PokemonPastType    `json:"past_types"`
	PastAbilities          []PokemonAbilityPast `json:"past_abilities"`
}

func GetPokemon(pokemon string) (Pokemon, error) {
	url := BASE_URL + "/pokemon/" + pokemon

	data, err := cachedGet(url)
	if err != nil {
		return Pokemon{}, err
	}

	var result Pokemon
	err = json.Unmarshal(data, &result)
	if err != nil {
		return Pokemon{}, errors.New("unable to parse Pokemon JSON")
	}

	return result, nil
}
