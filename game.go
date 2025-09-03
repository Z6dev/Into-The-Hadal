package main

import (
	"bytes"
	"fmt"
	"image/color"
	"math"

	"github.com/Z6dev/Boomtown/assets"
	"github.com/Z6dev/Boomtown/entities"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct {
	Initialized bool

	player *entities.SubmarineController
	cam    *entities.Camera

	tileImage *ebiten.Image

	// AUDIO AASAHDWAHDIWAHUDHWAIUDHUSAI
	audioCtx  *audio.Context
	thudSound *audio.Player
}

func (g *Game) Update() error {
	if !g.Initialized {
		stream, err := mp3.DecodeWithoutResampling(bytes.NewReader(assets.MetalThud_mp3))
		if err != nil {
			panic(err)
		}

		g.thudSound, err = g.audioCtx.NewPlayer(stream)
		if err != nil {
			panic(err)
		}
		g.Initialized = true
	}
	inputX, inputY := Axis(1)

	g.player.UpdateMovement(
		inputX,
		inputY,
		Acceleration,
		Friction,
		g.player.MoveSpeed,
		Collider,
	)

	handleCollision(g)

	g.cam.FollowTarget(g.player.X+32.0, g.player.Y+16.0, screenWidth, screenHeight)
	g.cam.Update(0.01666)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{4, 158, 209, 255})
	defer ebitenutil.DebugPrint(screen, fmt.Sprintf("%f\n%f\n%.2f", g.player.VelX, g.player.VelY, ebiten.ActualFPS()))

	opts := &ebiten.DrawImageOptions{}

	// Draw Player
	opts.GeoM.Translate(g.player.X, g.player.Y)
	opts.GeoM.Translate(g.cam.X, g.cam.Y)
	screen.DrawImage(g.player.Img, opts)

	opts.GeoM.Reset()

	tileW := Collider.TileSize[0]
	tileH := Collider.TileSize[1]

	// Compute visible bounds in tile coordinates
	startX := int((-g.cam.X) / float64(tileW))
	startY := int((-g.cam.Y) / float64(tileH))
	endX := int((-g.cam.X+screenWidth)/float64(tileW)) + 1
	endY := int((-g.cam.Y+screenHeight)/float64(tileH)) + 1

	// Clamp to map size
	if startX < 0 {
		startX = 0
	}
	if startY < 0 {
		startY = 0
	}
	if endX > len(TileMap[0]) {
		endX = len(TileMap[0])
	}
	if endY > len(TileMap) {
		endY = len(TileMap)
	}

	// Draw only visible tiles
	for y := startY; y < endY; y++ {
		for x := startX; x < endX; x++ {
			if TileMap[y][x] != 0 {
				opts.GeoM.Translate(float64(x*tileW), float64(y*tileH))
				opts.GeoM.Translate(g.cam.X, g.cam.Y)
				screen.DrawImage(g.tileImage, opts)
				opts.GeoM.Reset()
			}
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return int(screenWidth), int(screenHeight)
}

/*
@Z6dev is not happy with this code
What is this?

blegh
*/

/* ================= Helper functions non-game ================ */
func Axis(_ float64) (x, y float64) {
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		y = -1
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		y = 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		x = -1
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		x = 1
	}
	return
}

func handleCollision(g *Game) {
	const bounceFactor = 0.5 // elasticity factor (0 = no bounce, 1 = perfect bounce)

	for _, c := range Collider.Collisions {
		var impact float64

		// Vertical collisions
		if c.Normal[1] != 0 {
			impact = math.Abs(g.player.VelY)
			if impact > 0.5 { // threshold to avoid tiny bumps
				if g.thudSound != nil {
					vol := math.Min(1.0, impact/2.0) // scale speed to [0,1]
					g.thudSound.SetVolume(vol)
					g.thudSound.Rewind()
					g.thudSound.Play()
				}
			}
			g.player.VelY = -g.player.VelY * bounceFactor
			if math.Abs(g.player.VelY) < 0.05 {
				g.player.VelY = 0
			}
		}

		// Horizontal collision Right
		if c.Normal[0] == 1 {
			impact = math.Abs(g.player.VelX)
			if impact > 0.5 {
				if g.thudSound != nil {
					vol := math.Min(1.0, impact/2.0)
					g.thudSound.SetVolume(vol)
					g.thudSound.Rewind()
					g.thudSound.Play()
				}
				g.player.VelX = 0.4
				g.cam.Shake(0.2, 2)
			}
		}

		// Horizontal collision Left
		if c.Normal[0] == -1 {
			impact = math.Abs(g.player.VelX)
			if impact > 0.5 {
				if g.thudSound != nil {
					vol := math.Min(1.0, impact/2.0)
					g.thudSound.SetVolume(vol)
					g.thudSound.Rewind()
					g.thudSound.Play()
				}
				g.player.VelX = -0.4
				g.cam.Shake(0.2, 2)
			}
		}
	}
}
