package main

import (
	"bytes"
	"image/color"

	"github.com/Z6dev/Into-The-Hadal/assets"
	"github.com/Z6dev/Into-The-Hadal/entities"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"github.com/aquilax/go-perlin"
	"github.com/setanarut/tilecollider"

	"github.com/hajimehoshi/ebiten/v2/audio"
)

func generateTerrain() [][]uint8 {
	grid := make([][]uint8, mapWidth)
	for i := range grid {
		grid[i] = make([]uint8, mapHeight)
	}

	// Perlin noise parameters
	alpha := 2.0    // persistence
	beta := 2.0     // frequency
	var n int32 = 3 // number of octaves
	seed := int64(42)

	p := perlin.NewPerlin(alpha, beta, n, seed)

	for x := 0; x < mapWidth; x++ {
		for y := 0; y < mapHeight; y++ {
			// Scale coordinates so the noise doesn’t look too chaotic
			noise := p.Noise2D(float64(x)/20.0, float64(y)/20.0)

			// Normalize noise [-1,1] -> [0,1]
			normalized := (noise + 1) / 2

			// Threshold to decide solid vs empty
			if normalized > 0.5 {
				grid[x][y] = 1 // solid seabed
			} else {
				grid[x][y] = 0 // water/empty
			}
		}
	}
	return grid
}

func PlayAudio(data []byte, ctx *audio.Context) {
	pl, err := ctx.NewPlayer(bytes.NewReader(data))
	if err != nil {
		panic(err)
	}
	pl.Play()
}

/* ================= Main and Init ================= */
const screenWidth, screenHeight float64 = 480, 320
const (
	Acceleration = 0.012 // how fast player ramp up
	Friction     = 0.002 // how fast player slow down

	mapWidth  = 400
	mapHeight = 60
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
				H:   31,
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
