package entities

import "github.com/setanarut/tilecollider"

type SubmarineController struct {
	*Sprite
	MoveSpeed float64

	VelX, VelY float64
}

func (s *SubmarineController) UpdateMovement(inputX, inputY, acceleration, friction, maxSpeed float64, collider *tilecollider.Collider[uint8]) {
	// Accelerate based on input
	if inputX != 0 {
		s.VelX += inputX * acceleration
	} else {
		// apply friction when no input
		if s.VelX > 0 {
			s.VelX -= friction
			if s.VelX < 0 {
				s.VelX = 0
			}
		} else if s.VelX < 0 {
			s.VelX += friction
			if s.VelX > 0 {
				s.VelX = 0
			}
		}
	}

	if inputY != 0 {
		s.VelY += inputY * acceleration
	} else {
		if s.VelY > 0 {
			s.VelY -= friction
			if s.VelY < 0 {
				s.VelY = 0
			}
		} else if s.VelY < 0 {
			s.VelY += friction
			if s.VelY > 0 {
				s.VelY = 0
			}
		}
	}

	// Clamp max speed
	if s.VelX > maxSpeed {
		s.VelX = maxSpeed
	} else if s.VelX < -maxSpeed {
		s.VelX = -maxSpeed
	}
	if s.VelY > maxSpeed {
		s.VelY = maxSpeed
	} else if s.VelY < -maxSpeed {
		s.VelY = -maxSpeed
	}

	// Apply velocity through collider
	deltaX, deltaY := collider.Collide(
		s.X,
		s.Y,
		s.W,
		s.H,
		s.VelX,
		s.VelY,
		nil,
	)

	s.X += deltaX
	s.Y += deltaY
}
