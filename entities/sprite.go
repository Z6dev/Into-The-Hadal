package entities

import "github.com/hajimehoshi/ebiten/v2"

type Sprite struct {
	X, Y, W, H float64
	Img        *ebiten.Image
}
