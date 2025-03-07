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
	X, Y float64
}

type Player struct {
	Position Vec2
	Rotation float64
	Cooldown float64
}

func (p *Player) MoveForward(dt float64) { // TODO: how to sync dt between server and game?
	forwardDir := p.getForward()
	p.Position.X += PLAYER_SPEED * forwardDir.X * dt
	p.Position.Y += PLAYER_SPEED * forwardDir.Y * dt
}

func (p *Player) MoveBackward(dt float64) {
	forwardDir := p.getForward()
	p.Position.X -= PLAYER_SPEED * forwardDir.X * dt
	p.Position.Y -= PLAYER_SPEED * forwardDir.Y * dt
}

func (p *Player) RotateLeft(dt float64) {
	p.Rotation -= PLAYER_TURN_SPEED * dt
}

func (p *Player) RotateRight(dt float64) {
	p.Rotation += PLAYER_TURN_SPEED * dt
}

func (p *Player) Shoot() {
	panic("TODO")
}

func (p *Player) getForward() Vec2 {
	return Vec2{X: math.Cos(p.Rotation), Y: math.Sin(p.Rotation)}
}
