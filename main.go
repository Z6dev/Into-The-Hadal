package main

import (
	"bytes"
	"image/color"
	"math"
	"math/rand"

	"github.com/Z6dev/Into-The-Hadal/assets"
	"github.com/Z6dev/Into-The-Hadal/entities"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"github.com/aquilax/go-perlin"
	"github.com/setanarut/tilecollider"

	"github.com/hajimehoshi/ebiten/v2/audio"
)

// GenerateTerrain builds a 2D terrain map using layered Perlin noise.
// 0 = water, 1 = terrain
func GenerateTerrain(width, height int, seed int64) [][]uint8 {
	alpha := 2.0
	beta := 2.0
	octaves := 5

	// Base noise (big caves)
	p1 := perlin.NewPerlin(alpha, beta, int32(octaves), seed)
	// Detail noise (smaller features)
	p2 := perlin.NewPerlin(alpha, beta, int32(octaves), seed+42)

	terrain := make([][]uint8, height)
	for y := 0; y < height; y++ {
		terrain[y] = make([]uint8, width)

		depthFactor := float64(y) / float64(height)

		// Adjust threshold to control density
		// Start more solid near top, slightly emptier deep down
		threshold := 0.20 - 0.10*depthFactor

		for x := 0; x < width; x++ {
			// Stretched coordinates
			nx := float64(x) / 20.0 // horizontal scale
			ny := float64(y) / 50.0 // vertical scale (smaller than before)

			// Large-scale caves
			baseNoise := math.Abs(p1.Noise2D(nx, ny))

			// Smaller details (tunnels, cracks)
			detailNoise := p2.Noise2D(nx*2.0, ny*2.0) * 0.5

			// Combine them
			noiseVal := baseNoise*0.7 + detailNoise*0.3

			if noiseVal > threshold {
				terrain[y][x] = 1
			} else {
				terrain[y][x] = 0
			}
		}
	}
	return terrain
}

func PlayAudio(data []byte, ctx *audio.Context) {
	pl, err := ctx.NewPlayer(bytes.NewReader(data))
	if err != nil {
		panic(err)
	}
	pl.Play()
}

/* ================= Main and Init ================= */

const (
	screenWidth  float64 = 480
	screenHeight float64 = 320

	// BEWARE, MAP WIDTH AND MAP HEIGHT IS SWAPPED.
	// @Z6dev still doesn't know how to fix it 😭
	mapWidth  int = 800 // Map height
	mapHeight int = 800 // Map Width
)

var (
	TileMap [][]uint8

	Collider *tilecollider.Collider[uint8]

	tileW, tileH int
)

func main() {
	playerImg := ebiten.NewImage(64, 32)
	playerImg.Fill(color.RGBA{255, 0, 0, 255})

	tileImg, _, err := ebitenutil.NewImageFromReader(bytes.NewReader(assets.Sandtile_png))
	if err != nil {
		panic(err)
	}
	// Init terrain
	TileMap = GenerateTerrain(mapWidth, mapHeight, rand.Int63n(100))
	Collider = tilecollider.NewCollider(TileMap, 32, 32)

	tileW, tileH = Collider.TileSize[0], Collider.TileSize[1]

	// Init player
	var spawnX, spawnY float64 = float64(mapHeight/2) * 32.0, -32.0
	// GAME INIT!!!

	game := &Game{
		audioCtx: audio.NewContext(44100),

		player: &entities.SubmarineController{
			Sprite: &entities.Sprite{
				X:   spawnX,
				Y:   spawnY,
				W:   62,
				H:   31,
				Img: playerImg,
			},
			MoveSpeed:    2,
			Acceleration: 0.025,
			Friction:     0.01,
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
