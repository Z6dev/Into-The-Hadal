package main

import (
	"bytes"
	"image/color"
	"math/rand"

	"github.com/Z6dev/Into-The-Hadal/assets"
	"github.com/Z6dev/Into-The-Hadal/entities"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"github.com/aquilax/go-perlin"
	"github.com/setanarut/tilecollider"

	"github.com/hajimehoshi/ebiten/v2/audio"
)

func generateTerrain() [][]uint8 {
	grid := make([][]uint8, mapHeight)
	for i := range grid {
		grid[i] = make([]uint8, mapWidth)
	}

	// Perlin noise parameters
	alpha := 2.0    // persistence
	beta := 2.0     // frequency
	var n int32 = 7 // number of octaves
	var seed int64 = rand.Int63n(100)

	p := perlin.NewPerlin(alpha, beta, n, seed)

	for y := 0; y < mapHeight; y++ {
		for x := 0; x < mapWidth; x++ {
			noise := p.Noise2D(float64(x)/20.0, float64(y)/20.0)
			normalizedValue := (noise + 1) / 2
			if normalizedValue > 0.5 {
				grid[y][x] = 1
			} else {
				grid[y][x] = 0
			}
		}
	}
	return grid
}

func FindSpawnPoint(tilemap [][]uint8, tileSize float64) (float64, float64, bool) {
	rows := len(tilemap)
	if rows == 0 {
		return 0, 0, false
	}
	cols := len(tilemap[0])

	for y := 0; y < rows; y++ {
		for x := 0; x < cols-1; x++ { // -1 because player is 2 tiles wide
			// Check if this spot (2 wide, 1 high) is empty
			if tilemap[y][x] == 0 && tilemap[y][x+1] == 0 {
				// Convert to world coordinates (spawn at top-left of the area)
				return float64(x) * tileSize, float64(y) * tileSize, true
			}
		}
	}
	// No valid spawn found
	return 0, 0, false
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
	mapWidth  = 800 // Map height
	mapHeight = 800 // Map Width
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
	TileMap = generateTerrain()
	Collider = tilecollider.NewCollider(TileMap, 32, 32)

	tileW, tileH = Collider.TileSize[0], Collider.TileSize[1]

	// Init player
	spawnX, spawnY, ok := FindSpawnPoint(TileMap, 32)
	if !ok {
		panic("CANNOT FIND SPAWN POINT!")
	}
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
