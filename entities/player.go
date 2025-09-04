package entities

import "github.com/setanarut/tilecollider"

type SubmarineController struct {
	*Sprite
	MoveSpeed    float64
	Acceleration float64
	Friction     float64

	VelX, VelY float64
}

func (sub *SubmarineController) UpdateMovement(inputX, inputY float64, collider *tilecollider.Collider[uint8]) {
	// Accelerate based on input
	if inputX != 0 {
		sub.VelX += inputX * sub.Acceleration
	} else {
		// apply friction when no input
		if sub.VelX > 0 {
			sub.VelX -= sub.Friction
			if sub.VelX < 0 {
				sub.VelX = 0
			}
		} else if sub.VelX < 0 {
			sub.VelX += sub.Friction
			if sub.VelX > 0 {
				sub.VelX = 0
			}
		}
	}

	if inputY != 0 {
		sub.VelY += inputY * sub.Acceleration
	} else {
		if sub.VelY > 0 {
			sub.VelY -= sub.Friction
			if sub.VelY < 0 {
				sub.VelY = 0
			}
		} else if sub.VelY < 0 {
			sub.VelY += sub.Friction
			if sub.VelY > 0 {
				sub.VelY = 0
			}
		}
	}

	// Clamp max speed
	if sub.VelX > sub.MoveSpeed {
		sub.VelX = sub.MoveSpeed
	} else if sub.VelX < -sub.MoveSpeed {
		sub.VelX = -sub.MoveSpeed
	}
	if sub.VelY > sub.MoveSpeed {
		sub.VelY = sub.MoveSpeed
	} else if sub.VelY < -sub.MoveSpeed {
		sub.VelY = -sub.MoveSpeed
	}

	// Apply velocity through collider
	deltaX, deltaY := collider.Collide(
		sub.X,
		sub.Y,
		sub.W,
		sub.H,
		sub.VelX,
		sub.VelY,
		nil,
	)

	sub.X += deltaX
	sub.Y += deltaY
}
