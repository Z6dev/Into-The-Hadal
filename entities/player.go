package entities

import "github.com/setanarut/tilecollider"

type SubmarineController struct {
	*Sprite
	MoveSpeed float64

	VelX, VelY float64
}

func (sub *SubmarineController) UpdateMovement(inputX, inputY, acceleration, friction, maxSpeed float64, collider *tilecollider.Collider[uint8]) {
	// Accelerate based on input
	if inputX != 0 {
		sub.VelX += inputX * acceleration
	} else {
		// apply friction when no input
		if sub.VelX > 0 {
			sub.VelX -= friction
			if sub.VelX < 0 {
				sub.VelX = 0
			}
		} else if sub.VelX < 0 {
			sub.VelX += friction
			if sub.VelX > 0 {
				sub.VelX = 0
			}
		}
	}

	if inputY != 0 {
		sub.VelY += inputY * acceleration
	} else {
		if sub.VelY > 0 {
			sub.VelY -= friction
			if sub.VelY < 0 {
				sub.VelY = 0
			}
		} else if sub.VelY < 0 {
			sub.VelY += friction
			if sub.VelY > 0 {
				sub.VelY = 0
			}
		}
	}

	// Clamp max speed
	if sub.VelX > maxSpeed {
		sub.VelX = maxSpeed
	} else if sub.VelX < -maxSpeed {
		sub.VelX = -maxSpeed
	}
	if sub.VelY > maxSpeed {
		sub.VelY = maxSpeed
	} else if sub.VelY < -maxSpeed {
		sub.VelY = -maxSpeed
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
