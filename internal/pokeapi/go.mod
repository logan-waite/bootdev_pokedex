module github.com/logan-waite/bootdev_pokedex/internal/pokeapi

go 1.23.3

replace github.com/logan-waite/bootdev_pokedex/internal/pokecache v0.0.0 => ../pokecache

require (
	github.com/logan-waite/bootdev_pokedex/internal/pokecache v0.0.0
)
