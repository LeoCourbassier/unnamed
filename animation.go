package main

import "fmt"

type AnimationsSprite int

type Animations struct {
	CurrentAnimation AnimationsSprite
	FirstAnimation   AnimationsSprite
	LastAnimation    AnimationsSprite
	AnimationTicks   int
	LoopAnimation    bool
	Ticks            int
}

const (
	I0  AnimationsSprite = iota
	I1  AnimationsSprite = iota
	I2  AnimationsSprite = iota
	I3  AnimationsSprite = iota
	W0  AnimationsSprite = iota
	W1  AnimationsSprite = iota
	W2  AnimationsSprite = iota
	W3  AnimationsSprite = iota
	W4  AnimationsSprite = iota
	W5  AnimationsSprite = iota
	J0  AnimationsSprite = iota
	J1  AnimationsSprite = iota
	J2  AnimationsSprite = iota
	J3  AnimationsSprite = iota
	A0  AnimationsSprite = iota
	A1  AnimationsSprite = iota
	A2  AnimationsSprite = iota
	A3  AnimationsSprite = iota
	A4  AnimationsSprite = iota
	A5  AnimationsSprite = iota
	AF0 AnimationsSprite = iota
	AF1 AnimationsSprite = iota
	AF2 AnimationsSprite = iota
	AF3 AnimationsSprite = iota
	AF4 AnimationsSprite = iota
	AF5 AnimationsSprite = iota
)

func (a *Animations) Update() {
	fmt.Println(a.Ticks)
	if a.Ticks > a.AnimationTicks {
		a.Ticks = 0
		if a.FirstAnimation != a.LastAnimation {
			a.CurrentAnimation++
		}
		if a.CurrentAnimation > a.LastAnimation {
			if Player.IsAttacking {
				Player.IsAttacking = false
				for i := range Enemies {
					for j := range Enemies[i].Damage {
						if Enemies[i].Damage[j].LastTick && Enemies[i].Damage[j].Giver == Player.ID {
							Enemies[i].Damage[j].LastTick = false
							break
						}
					}
				}
				a.CurrentAnimation = I0
				a.FirstAnimation = I0
				a.LastAnimation = I3
				a.AnimationTicks = 7
				a.LoopAnimation = true
			}
			if a.LoopAnimation {
				a.CurrentAnimation = a.FirstAnimation
			} else {
				a.CurrentAnimation = a.LastAnimation
			}
		}
	}
}
