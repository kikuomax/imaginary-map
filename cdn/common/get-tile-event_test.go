package common

import (
	"testing"
)

func TestGetTileEventToTile (t *testing.T) {
	t.Run(
		"GetTileEvent{X:0,Y:0,Zoom:0}.ToTile() should be maptile.Tile{X:0,Y:0,Z:0}",
		func (t *testing.T) {
			event := GetTileEvent {
				X: 0,
				Y: 0,
				Zoom: 0,
			}
			tile, ok := event.ToTile()
			if !ok {
				t.Error("expected GetTileEvent{X:0,Y:0,Zoom:0}.ToTile() is ok but not")
			}
			if (tile.X != 0) || (tile.Y != 0) || (tile.Z != 0) {
				t.Errorf(
					"expected GetTileEvent{X:0,Y:0,Zoom:0}.ToTile() is maptile.Tile{X:0,Y:0:Z:0} but got maptile.Tile{X:%v,Y:%v,Z:%v}",
					tile.X,
					tile.Y,
					tile.Z)
			}
		})
	t.Run(
		"GetTileEvent{X:4,Y:2,Zoom:3}.ToTile() should be maptile.Tile{X:4,Y:2,Zoom:3}",
		func (t *testing.T) {
			event := GetTileEvent {
				X: 4,
				Y: 2,
				Zoom: 3,
			}
			tile, ok := event.ToTile()
			if !ok {
				t.Error("expected GetTileEvent{X:4,Y:2,Zoom:3}.ToTile() is ok but not")
			}
			if (tile.X != 4) || (tile.Y != 2) || (tile.Z != 3) {
				t.Errorf(
					"expected GetTileEvent{X:4,Y:2,Zoom:3}.ToTile() is maptile.Tile{X:4,Y:2,Z:3} but got maptile.Tile{X:%v,Y:%v,Z:%v}",
					tile.X,
					tile.Y,
					tile.Z)
			}
		})
	t.Run(
		"GetTileEvent{X:-1,Y:0,Zoom:0}.ToTile() should not be ok",
		func (t *testing.T) {
			event := GetTileEvent {
				X: -1,
				Y: 0,
				Zoom: 0,
			}
			_, ok := event.ToTile()
			if ok {
				t.Error("expected GetTileEvent{X:-1,Y:0,Zoom:0}.ToTile() is not ok but is")
			}
		})
	t.Run(
		"GetTileEvent{X:0,Y:-1,Zoom:0}.ToTile() should not be ok",
		func (t *testing.T) {
			event := GetTileEvent {
				X: 0,
				Y: -1,
				Zoom: 0,
			}
			_, ok := event.ToTile()
			if ok {
				t.Error("expected GetTileEvent{X:0,Y:-1,Zoom:0}.ToTile() is not ok but is")
			}
		})
	t.Run(
		"GetTileEvent{X:0,Y:0,Zoom:-1}.ToTile() should not be ok",
		func (t *testing.T) {
			event := GetTileEvent {
				X: 0,
				Y: 0,
				Zoom: -1,
			}
			_, ok := event.ToTile()
			if ok {
				t.Error("expected GetTileEvent{X:0,Y:0,Zoom:-1}.ToTile() is not ok but is")
			}
		})
}
