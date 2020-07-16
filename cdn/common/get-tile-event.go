package common

import (
	"github.com/paulmach/orb/maptile"
)

// Input type for a get-tile event.
type GetTileEvent struct {
	Zoom int `json:"zoom"`
	X int `json:"x"`
	Y int `json:"y"`
}

// Converts a `GetTileEvent` into a `maptile.Tile`
func (event GetTileEvent) ToTile () maptile.Tile {
	return maptile.New(
		uint32(event.X),
		uint32(event.Y),
		maptile.Zoom(event.Zoom),
	)
}
