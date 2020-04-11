package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

type Camera struct {
	X float64
	Y float64
}

type MessageFeedback struct {
	Message string
	Seconds float64
	X       float64
	Y       float64
	Time    int64
}

func (c Camera) InViewport(o Object) bool {
	return (c.X <= o.X() && c.X+float64(App.Width) >= o.X()) && (c.Y <= o.Y() && c.Y+float64(App.Height) >= o.Y())
}

func (c Camera) Draw(o Object, image int, screen *ebiten.Image) {
	relOpt := *o.Options
	relOpt.GeoM.Translate(-c.X, -c.Y)

	screen.DrawImage(o.Img[image], &relOpt)
}

func (c Camera) DrawFixed(o Object, image int, screen *ebiten.Image) {
	screen.DrawImage(o.Img[image], o.Options)
}

func (c Camera) DrawText(screen *ebiten.Image, msg string, x int, y int) {
	x -= int(c.X)
	y -= int(c.Y)
	ebitenutil.DebugPrintAt(screen, msg, x, y)
}

func (c Camera) DrawTextFixed(screen *ebiten.Image, msg string, x int, y int) {
	ebitenutil.DebugPrintAt(screen, msg, x, y)
}

func (c Camera) DrawRect(dst *ebiten.Image, x float64, y float64, width float64, height float64, clr color.Color) {
	x -= c.X
	y -= c.Y
	ebitenutil.DrawRect(dst, x, y, width, height, clr)
}

func (c Camera) DrawRectFixed(dst *ebiten.Image, x float64, y float64, width float64, height float64, clr color.Color) {
	ebitenutil.DrawRect(dst, x, y, width, height, clr)
}
