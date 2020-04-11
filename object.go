package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"math/rand"
	"path/filepath"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

type Object struct {
	ID            int
	Img           []*ebiten.Image
	Options       *ebiten.DrawImageOptions
	RealHeight    float64
	RealWidth     float64
	OffsetX       float64
	OffsetY       float64
	HasMass       bool
	isCollideable bool
	MaxHealth     float64
	Health        float64
	MeleeRange    float64
	AttackDamage  float64
	CritPercent   float64
	Damage        []CombatRegistry
	Animation     Animations
}

type CombatRegistry struct {
	Giver    int
	Quantity int
	LastTick bool
}

func CreateObject(wantedH, wantedW float64, path string, realH, realW float64, offsetX, offsetY float64, hasMass bool, collides bool, id int) Object {
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

	return Object{
		ID:            id,
		Img:           []*ebiten.Image{img},
		Options:       options,
		RealHeight:    realH,
		RealWidth:     realW,
		OffsetX:       offsetX,
		OffsetY:       offsetY,
		HasMass:       hasMass,
		isCollideable: collides,
	}
}

func (o Object) Intersects(other Object) bool {
	leftX := math.Max(o.X(), other.X())
	rightX := math.Min(o.X()+o.Width(), other.X()+other.Width())
	topY := math.Max(o.Y(), other.Y())
	bottomY := math.Min(o.Y()+o.Height(), other.Y()+other.Height())

	return leftX < rightX && topY < bottomY
}

func (o Object) IntersectsSideways(other Object) bool {
	if !o.Intersects(other) {
		return false
	}

	topO := o.Y()
	bottomO := topO + o.Height()

	topOther := other.Y()
	//bottomOther := topOther + other.Height()

	return !(bottomO-10 <= topOther) //|| bottomOther-10 <= topO)
}

func (o Object) SidewayException(other Object) bool {

	if !o.IntersectsSideways(other) {
		return false
	}

	leftO := o.X()
	rightO := leftO + o.Width()

	leftOther := other.X()
	rightOther := leftOther + other.Width()

	return (rightO-10 <= leftOther || rightOther-10 <= leftO)
}

func (o Object) IntersectsArray(other []Object) bool {
	for _, obj := range other {
		if o.Intersects(obj) && obj.isCollideable {
			return true
		}
	}
	return false
}

func (o Object) IntersectsArraySideways(other []Object) bool {
	for _, obj := range other {
		if o.IntersectsSideways(obj) && obj.isCollideable {
			return true
		}
	}
	return false
}

func (o Object) SidewayExceptionArray(other []Object) bool {
	for _, obj := range other {
		if o.SidewayException(obj) && obj.isCollideable {
			return true
		}
	}
	return false
}

func (o Object) RawX() float64 {
	return o.Options.GeoM.Element(0, 2)
}

func (o Object) RawY() float64 {
	return o.Options.GeoM.Element(1, 2)
}

func (o Object) X() float64 {
	if o.ScaleX() < 0 {
		return o.Options.GeoM.Element(0, 2) - o.Width() - o.OffsetX*math.Abs(o.ScaleX())
	}
	return o.Options.GeoM.Element(0, 2) + o.OffsetX*math.Abs(o.ScaleX())
}

func (o Object) Y() float64 {
	return o.Options.GeoM.Element(1, 2) + o.OffsetY*math.Abs(o.ScaleY())
}

func (o Object) DebugXY() string {
	return fmt.Sprintf("(%.2f, %.2f)", o.X(), o.Y())
}

func (o Object) ResetXY() {
	o.Options.GeoM.Translate(-o.X(), -o.Y())
}

func (o Object) ScaleX() float64 {
	return o.Options.GeoM.Element(0, 0)
}

func (o Object) ScaleY() float64 {
	return o.Options.GeoM.Element(1, 1)
}

func (o Object) Width() float64 {
	return o.RealWidth * math.Abs(o.ScaleX())
}

func (o Object) Height() float64 {
	return o.RealHeight * math.Abs(o.ScaleY())
}

func (o Object) Reflect() {
	x := o.X()
	y := o.Y()
	o.ResetXY()
	o.Options.GeoM.Scale(-1, 1)
	o.Options.GeoM.Translate(x+o.Width(), y)
}

func (o Object) Range() float64 {
	return o.MeleeRange * math.Abs(o.ScaleX())
}

func (o Object) FacingEnemy(other Object) bool {
	if o.ScaleX() < 0 {
		return o.X() >= other.X()
	}

	return o.X() <= other.X()
}

func (o Object) WillCritAttack() bool {
	return rand.Float64()*100 <= o.CritPercent
}

func (o Object) InAttackRange(other Object) bool {
	x := 0.0

	topO := o.Y()
	bottomO := topO + o.Height()

	topOther := other.Y()
	bottomOther := topOther + other.Height()

	sameHeight := !(bottomO-10 <= topOther || bottomOther-10 <= topO)

	if o.ScaleX() < 0 {
		if other.ScaleX() < 0 {
			x = other.X()
		} else {
			x = other.X() + other.Width()
		}

		return o.X()-o.Range() <= x && o.FacingEnemy(other) && sameHeight
	}

	if other.ScaleX() < 0 {
		x = other.X() - other.Width()
	} else {
		x = other.X()
	}
	return o.X()+o.Width()+o.Range() >= x && o.FacingEnemy(other) && sameHeight
}

func (o *Object) Update(screen *ebiten.Image) {
	o.Animation.Update()
	o.Draw(screen)
}

func (o *Object) Draw(screen *ebiten.Image) {
	MainCamera.Draw(*o, int(o.Animation.CurrentAnimation), screen)

	MainCamera.DrawRect(screen, o.X(), o.Y()-20, o.Width(), 16, color.Gray16{0xCCCF})
	barWidth := (o.Health / o.MaxHealth) * o.Width()
	MainCamera.DrawRect(screen, o.X(), o.Y()-20, barWidth, 16, color.RGBA{
		A: 0xFF,
		R: 0xFF,
		G: 0x00,
		B: 0x00,
	})
	MainCamera.DrawText(screen, fmt.Sprintf("%.0f/%.0f", o.Health, o.MaxHealth), int(o.X()+o.Width()/2-20), int(o.Y()-20))
}
