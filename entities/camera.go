package entities

import (
	"math"
	"math/rand"
)

type Camera struct {
	X, Y float64

	// shake parameters
	shakeDuration float64
	shakeStrength float64
	shakeTime     float64
}

// NewCamera creates a new Camera instance
func NewCamera(x, y float64) *Camera {
	return &Camera{
		X: x,
		Y: y,
	}
}

// FollowTarget centers the camera on the given target
func (c *Camera) FollowTarget(targetX, targetY, screenWidth, screenHeight float64) {
	c.X = -targetX + screenWidth/2.0
	c.Y = -targetY + screenHeight/2.0
}

// Constrain keeps the camera within the tilemap boundaries
func (c *Camera) Constrain(tilemapWidthPixels, tilemapHeightPixels, screenWidth, screenHeight float64) {
	c.X = math.Min(c.X, 0.0)
	c.Y = math.Min(c.Y, 0.0)

	c.X = math.Max(c.X, screenWidth-tilemapWidthPixels)
	c.Y = math.Max(c.Y, screenHeight-tilemapHeightPixels)
}

// Shake starts a camera shake with the given duration (seconds) and strength (pixels)
func (c *Camera) Shake(duration, strength float64) {
	c.shakeDuration = duration
	c.shakeStrength = strength
	c.shakeTime = duration
}

// Update should be called every frame with deltaTime (seconds)
// Applies shake effect if active
func (c *Camera) Update(deltaTime float64) {
	if c.shakeTime > 0 {
		c.shakeTime -= deltaTime
		strength := c.shakeStrength * (c.shakeTime / c.shakeDuration) // decay

		offsetX := (rand.Float64()*2 - 1) * strength
		offsetY := (rand.Float64()*2 - 1) * strength

		c.X += offsetX
		c.Y += offsetY
	}
}
