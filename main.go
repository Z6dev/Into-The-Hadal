package main

import (
	"bytes"
	"image/color"

	"github.com/Z6dev/Boomtown/assets"
	"github.com/Z6dev/Boomtown/entities"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"github.com/setanarut/tilecollider"

	"github.com/hajimehoshi/ebiten/v2/audio"
)

/* ================= Main and Init ================= */
const screenWidth, screenHeight float64 = 480, 320
const (
	Acceleration = 0.012 // how fast player ramp up
	Friction     = 0.002 // how fast player slow down

	mapWidth  = 80
	mapHeight = 40
)

var (
	TileMap [][]uint8

	Collider *tilecollider.Collider[uint8]
)

func main() {
	playerImg := ebiten.NewImage(64, 32)
	playerImg.Fill(color.RGBA{255, 0, 0, 255})

	tileImg, _, err := ebitenutil.NewImageFromReader(bytes.NewReader(assets.Sandtile_png))
	if err != nil {
		panic(err)
	}
	// Init terrain
	TileMap = generateTerrain()
	Collider = tilecollider.NewCollider(TileMap, 32, 32)

	// GAME INIT!!!

	game := &Game{
		audioCtx: audio.NewContext(44100),

		player: &entities.SubmarineController{
			Sprite: &entities.Sprite{
				X:   100,
				Y:   100,
				W:   64,
				H:   32,
				Img: playerImg,
			},
			MoveSpeed: 2,
		},

		cam:       entities.NewCamera(0, 0),
		tileImage: tileImg,
	}

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Into The Hadal")

	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
