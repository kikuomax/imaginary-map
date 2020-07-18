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
// The second return value is
// `false` if one or more of `event.Zoom`, `event.X`, `event.Y` are negative,
// `true` otherwise.
// The first return value is undefined if the second one is `false`.
func (event GetTileEvent) ToTile () (maptile.Tile, bool) {
	ok := (event.X >= 0) && (event.Y >= 0) && (event.Zoom >= 0)
	tile := maptile.New(
		uint32(event.X),
		uint32(event.Y),
		maptile.Zoom(event.Zoom),
	)
	return tile, ok
}
