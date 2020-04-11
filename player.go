package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"path/filepath"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

type PlayerObject struct {
	Object
	Score          int
	IsJumping      bool
	IsGrounded     bool
	Speed          float64
	FacingRight    bool
	AirSeconds     float64
	IsAttacking    bool
	IsStrongAttack bool
	Crited         bool
}

func CreatePlayer(wantedH, wantedW float64) PlayerObject {
	var img *ebiten.Image
	var err error
	imgs := []*ebiten.Image{}
	for i := 0; i < 4; i++ {
		img, _, err = ebitenutil.NewImageFromFile(filepath.FromSlash(fmt.Sprintf("assets/player/individual/adventurer-idle-2-0%d.png", i)), ebiten.FilterDefault)
		if err != nil {
			log.Fatal(err)
		}
		imgs = append(imgs, img)
	}

	for i := 0; i < 6; i++ {
		img, _, err = ebitenutil.NewImageFromFile(filepath.FromSlash(fmt.Sprintf("assets/player/individual/adventurer-run-0%d.png", i)), ebiten.FilterDefault)
		if err != nil {
			log.Fatal(err)
		}
		imgs = append(imgs, img)
	}

	for i := 0; i < 4; i++ {
		img, _, err = ebitenutil.NewImageFromFile(filepath.FromSlash(fmt.Sprintf("assets/player/individual/adventurer-jump-0%d.png", i)), ebiten.FilterDefault)
		if err != nil {
			log.Fatal(err)
		}
		imgs = append(imgs, img)
	}

	for i := 0; i < 6; i++ {
		img, _, err = ebitenutil.NewImageFromFile(filepath.FromSlash(fmt.Sprintf("assets/player/individual/adventurer-attack2-0%d.png", i)), ebiten.FilterDefault)
		if err != nil {
			log.Fatal(err)
		}
		imgs = append(imgs, img)
	}

	for i := 0; i < 6; i++ {
		img, _, err = ebitenutil.NewImageFromFile(filepath.FromSlash(fmt.Sprintf("assets/player/individual/adventurer-attack3-0%d.png", i)), ebiten.FilterDefault)
		if err != nil {
			log.Fatal(err)
		}
		imgs = append(imgs, img)
	}

	options := &ebiten.DrawImageOptions{
		GeoM: ebiten.GeoM{},
	}
	options.GeoM.Translate(50, 50)
	w, h := img.Size()
	options.GeoM.Scale((wantedW / float64(w)), (wantedH / float64(h)))

	return PlayerObject{
		Object: Object{
			ID:         0,
			Img:        imgs,
			Options:    options,
			RealHeight: 32.0,
			RealWidth:  19.0,
			OffsetX:    15.0,
			OffsetY:    6.0,
			/*RealHeight: 200.0,
			RealWidth:  142.0,
			OffsetX:    46.0,
			OffsetY:    22.0,*/
			HasMass:       true,
			isCollideable: true,
			MaxHealth:     100,
			Health:        1,
			MeleeRange:    13.0,
			AttackDamage:  10,
			CritPercent:   100,
			Animation: Animations{
				CurrentAnimation: I0,
				FirstAnimation:   I0,
				LastAnimation:    I3,
				AnimationTicks:   7,
				LoopAnimation:    true,
				Ticks:            0,
			},
		},
		Score:       0,
		IsJumping:   false,
		Speed:       1.2,
		FacingRight: true,
		IsGrounded:  false,
		AirSeconds:  0.50,
		IsAttacking: false,
		Crited:      false,
	}
}

func (o *PlayerObject) Move(x, y float64) {
	o.Options.GeoM.Translate(x, y)

	MainCamera.X += x
	MainCamera.Y += y
}

func jumpFn(t float64) float64 {
	return math.Exp(t*2) - 4
}

func (p *PlayerObject) Update(screen *ebiten.Image) {
	p.Animation.Update()
	p.CheckInputs()
	p.Draw(screen)
}

func (p *PlayerObject) Draw(screen *ebiten.Image) {
	MainCamera.Draw(Player.Object, int(p.Animation.CurrentAnimation), screen)
	if Player.Crited {
		if float64(time.Now().UnixNano()-Message.Time)*math.Pow(10, -9) < Message.Seconds {
			MainCamera.DrawText(screen, Message.Message, int(Message.X), int(Message.Y))
			Message.Y--
		} else {
			Player.Crited = false
		}
	}

	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Score:%d", Player.Score), 700, 500)
	MainCamera.DrawRectFixed(screen, 20, 20, 300, 32, color.Gray16{0xCCCF})
	barWidth := (Player.Health / Player.MaxHealth) * 300
	MainCamera.DrawRectFixed(screen, 20, 20, barWidth, 32, color.RGBA{
		A: 0xFF,
		R: 0xFF,
		G: 0x00,
		B: 0x00,
	})
	MainCamera.DrawTextFixed(screen, fmt.Sprintf("%.0f/%.0f", Player.Health, Player.MaxHealth), 150, 26)
}

func (p *PlayerObject) CheckInputs() {
	hasWalked := false
	InputDebounce(func() {
		for _, k := range keys {
			if ebiten.IsKeyPressed(k.Key) && !Player.IsAttacking {
				if k.Key == ebiten.KeyUp && !Player.IsJumping && Player.IsGrounded {
					LastJumpTime = time.Now()
					Player.IsJumping = true
					p.Animation.CurrentAnimation = J0
					p.Animation.FirstAnimation = J0
					p.Animation.LastAnimation = J3
					p.Animation.AnimationTicks = 2
					p.Animation.LoopAnimation = false
					Player.IsGrounded = false
					Player.IsAttacking = false
				} else if k.Key == ebiten.KeyLeft && Player.FacingRight {
					Player.Reflect()
					Player.FacingRight = false
				} else if k.Key == ebiten.KeyRight && !Player.FacingRight {
					Player.Reflect()
					Player.FacingRight = true
				} else if k.Key != ebiten.KeyUp {
					if !Player.IntersectsArraySideways(Tiles) && !Player.IsAttacking {
						hasWalked = true
						if (p.Animation.CurrentAnimation < W0 || p.Animation.CurrentAnimation > W5) && !Player.IsJumping && Player.IsGrounded {
							p.Animation.CurrentAnimation = W0
							p.Animation.FirstAnimation = W0
							p.Animation.LastAnimation = W5
							p.Animation.LoopAnimation = true
							p.Animation.AnimationTicks = 4
						}
						Player.Move(k.Tx, k.Ty)
					}
				}
			}
		}
		if Player.IsGrounded && !hasWalked && !Player.IsAttacking && !Player.IsJumping && (p.Animation.CurrentAnimation < I0 || p.Animation.CurrentAnimation > I3) {
			p.Animation.CurrentAnimation = I0
			p.Animation.FirstAnimation = I0
			p.Animation.LastAnimation = I3
			p.Animation.AnimationTicks = 7
			p.Animation.LoopAnimation = true
		}
	})

	zPressed := ebiten.IsKeyPressed(ebiten.KeyZ)
	if zPressed || ebiten.IsKeyPressed(ebiten.KeyX) {
		if !Player.IsAttacking && Player.IsGrounded {
			Player.IsAttacking = true
			if zPressed {
				p.Animation.CurrentAnimation = A0
				p.Animation.FirstAnimation = A0
				p.Animation.LastAnimation = A5
				p.Animation.AnimationTicks = 5
				Player.IsStrongAttack = false
			} else {
				p.Animation.CurrentAnimation = AF0
				p.Animation.FirstAnimation = AF0
				p.Animation.LastAnimation = AF5
				p.Animation.AnimationTicks = 5
				Player.IsStrongAttack = true
			}
			p.Animation.LoopAnimation = true
		}
	}
}
