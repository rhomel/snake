package images

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	TileWidth  = 10
	TileHeight = 10
)

func NewTile() *ebiten.Image {
	return NewColoredTile(color.White)
}

func NewColoredTile(c color.Color) *ebiten.Image {
	i := ebiten.NewImage(TileWidth, TileHeight)
	i.Fill(c)
	return i
}
