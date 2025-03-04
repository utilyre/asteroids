package main

import "math"

const (
	PLAYER_RADIUS         = 20.0
	PLAYER_TURN_SPEED     = 300.0
	PLAYER_SPEED          = 200.0
	PLAYER_SHOOT_SPEED    = 500.0
	PLAYER_SHOOT_COOLDOWN = 0.3
)

type Vec2 struct {
	x, y float64
}

type Player struct {
	Position Vec2
	Rotation float64
	Cooldown float64
}

func (p *Player) Move(dt float64, forward bool) {
	forwardDir := p.getForward()
	if forward {
		p.Position.x += PLAYER_SPEED * forwardDir.x * dt
		p.Position.y += PLAYER_SPEED * forwardDir.y * dt
	} else {
		p.Position.x -= PLAYER_SPEED * forwardDir.x * dt
		p.Position.y -= PLAYER_SPEED * forwardDir.y * dt
	}
}

func (p *Player) Rotate(dt float64, right bool) {
	if right {
		p.Rotation += PLAYER_TURN_SPEED * dt
	} else {
		p.Rotation -= PLAYER_TURN_SPEED * dt
	}
}

func (p *Player) Shoot() {
	panic("TODO")
}

func (p *Player) getForward() Vec2 {
	return Vec2{x: math.Cos(p.Rotation), y: math.Sin(p.Rotation)}
}
