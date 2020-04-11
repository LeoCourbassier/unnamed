package main

import (
	"log"
	"path/filepath"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

func CreateEnemy(wantedH, wantedW float64, path string, realH, realW float64, offsetX, offsetY float64, hasMass bool, collides bool, id int) PlayerObject {
	img, _, err := ebitenutil.NewImageFromFile(filepath.FromSlash(path), ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}

	options := &ebiten.DrawImageOptions{
		GeoM: ebiten.GeoM{},
	}
	w, h := img.Size()
	if wantedW == -1 {
		wantedW = float64(w)
	}

	if wantedH == -1 {
		wantedH = float64(h)
	}
	options.GeoM.Scale((wantedW / float64(w)), (wantedH / float64(h)))

	if realH == -1 {
		realH = float64(h)
	}

	if realW == -1 {
		realW = float64(w)
	}

	return PlayerObject{
		Object: Object{
			ID:            0,
			Img:           []*ebiten.Image{img},
			Options:       options,
			RealHeight:    realH,
			RealWidth:     realW,
			OffsetX:       offsetX,
			OffsetY:       offsetY,
			HasMass:       hasMass,
			isCollideable: collides,
			MaxHealth:     100,
			Health:        1,
			MeleeRange:    13.0,
			AttackDamage:  10,
			CritPercent:   100,
			Animation: Animations{
				CurrentAnimation: I0,
				FirstAnimation:   I0,
				LastAnimation:    I0,
				AnimationTicks:   -1,
				LoopAnimation:    false,
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
