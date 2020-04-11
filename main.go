package main

import (
	"fmt"
	"image/color"
	_ "image/png"
	"log"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

var Player PlayerObject
var Coin Object
var Tiles []Object
var Enemies []Object
var Background []Object
var keys []Control
var Debug bool
var JumpDebounce func(f func())
var InputDebounce func(f func())
var GravityDebounce func(f func())
var Gravity Control
var TimeDelta float64
var LastJumpTime time.Time
var App *Window

var Message MessageFeedback
var MainCamera Camera

type Window struct {
	Height int
	Width  int
}

type Control struct {
	Key ebiten.Key
	Tx  float64
	Ty  float64
}

func CreateCoin(wantedH, wantedW float64, gravity bool) Object {
	img, _, err := ebitenutil.NewImageFromFile(filepath.FromSlash("assets/coin.png"), ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}

	options := &ebiten.DrawImageOptions{
		GeoM: ebiten.GeoM{},
	}
	options.GeoM.Translate(150, 150)
	w, h := img.Size()
	options.GeoM.Scale((wantedW / float64(w)), (wantedH / float64(h)))

	return Object{
		Img:           []*ebiten.Image{img},
		Options:       options,
		RealHeight:    303.0,
		RealWidth:     303.0,
		OffsetX:       107.0,
		OffsetY:       107.0,
		HasMass:       gravity,
		isCollideable: true,
	}
}

func init() {
	App = &Window{
		Height: 600,
		Width:  800,
	}
	MainCamera = Camera{
		X: 0,
		Y: 0,
	}
	Debug = false
	JumpDebounce = NewDebouncer(50 * time.Millisecond)
	InputDebounce = NewDebouncer(100 * time.Microsecond)

	rand.Seed(time.Now().UnixNano())
	Player = CreatePlayer(100, 150)
	MainCamera.X = -(float64(App.Width/2) - 75) + Player.RawX()
	MainCamera.Y = -(float64(App.Height/2) - 50) + Player.RawY()

	Coin = CreateCoin(64, 64, true)

	floor := CreateObject(-1, float64(App.Width), "assets/Background.png", -1, -1, 0, 0, false, false, -1)
	floor.Options.GeoM.Translate(0, -floor.Height()+float64(App.Height))
	Background = append(Background, floor)

	floor = CreateObject(32, float64(App.Width), "assets/grass.png", -1, -1, 0, 0, false, true, -1)
	floor.Options.GeoM.Translate(500, 400)
	Tiles = append(Tiles, floor)

	floor = CreateObject(100, float64(App.Width), "assets/grass.png", -1, -1, 0, 0, false, true, -1)
	floor.Options.GeoM.Translate(0, 533)
	Tiles = append(Tiles, floor)

	enemy := CreateObject(80, 100, "assets/bat/bat_walk0.png", -1, -1, 0, 0, true, true, 1)
	for i := 1; i < 5; i++ {
		img, _, err := ebitenutil.NewImageFromFile(filepath.FromSlash(fmt.Sprintf("assets/bat/bat_walk%d.png", i)), ebiten.FilterDefault)
		if err != nil {
			log.Fatal(err)
		}
		enemy.Img = append(enemy.Img, img)
	}
	enemy.Options.GeoM.Translate(250, 150)
	enemy.MaxHealth = 100
	enemy.Health = 100
	Enemies = append(Enemies, enemy)

	enemy = CreateObject(80, 100, "assets/bat/bat_walk0.png", -1, -1, 0, 0, true, true, 2)
	for i := 1; i < 5; i++ {
		img, _, err := ebitenutil.NewImageFromFile(filepath.FromSlash(fmt.Sprintf("assets/bat/bat_walk%d.png", i)), ebiten.FilterDefault)
		if err != nil {
			log.Fatal(err)
		}
		enemy.Img = append(enemy.Img, img)
	}
	enemy.Options.GeoM.Translate(550, 150)
	enemy.MaxHealth = 100
	enemy.Health = 100
	Enemies = append(Enemies, enemy)

	Gravity = Control{
		Key: ebiten.KeyDown,
		Tx:  0,
		Ty:  5,
	}
	keys = []Control{
		{Key: ebiten.KeyUp, Tx: 0, Ty: -10},
		{Key: ebiten.KeyDown, Tx: 0, Ty: 0},
		{Key: ebiten.KeyLeft, Tx: Player.Speed * -3, Ty: 0},
		{Key: ebiten.KeyRight, Tx: Player.Speed * 3, Ty: 0},
	}
	TimeDelta = 5
}

func update(screen *ebiten.Image) error {
	if ebiten.IsDrawingSkipped() {
		fmt.Println("skipped")
		return nil
	}

	TimeDelta = float64(time.Now().UnixNano()-LastJumpTime.UnixNano()) * (math.Pow(10, -9))

	MainCamera.DrawFixed(Background[0], 0, screen)
	for _, tile := range Tiles {
		MainCamera.Draw(tile, 0, screen)
	}

	Player.Update(screen)

	if ebiten.IsKeyPressed(ebiten.KeyD) {
		JumpDebounce(func() {
			Debug = !Debug
		})
	}

	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		JumpDebounce(func() {
			os.Exit(0)
		})
	}

	if Player.Intersects(Coin) {
		Player.Score++
		Player.Health += 10
		Coin.ResetXY()
		newX := math.Max(rand.Float64()*float64(App.Width)-Coin.RealWidth+1, 0)
		newY := math.Max(rand.Float64()*float64(App.Height)-Coin.RealHeight+1, 0)
		Coin.Options.GeoM.Translate(newX, newY)
	}

	applyGravity()

	Player.Combat(Enemies)

	for _, e := range Enemies {
		e.Update(screen)
	}

	MainCamera.Draw(Coin, 0, screen)
	if Debug {
		MainCamera.DrawRect(screen, Coin.X(), Coin.Y(), Coin.Width(), Coin.Height(), color.White)

		for _, o := range Tiles {
			if o.isCollideable {
				MainCamera.DrawRect(screen, o.X(), o.Y(), o.Width(), o.Height(), color.White)
			}
		}
		for _, o := range Enemies {
			if o.isCollideable {
				MainCamera.DrawRect(screen, o.X(), o.Y(), o.Width(), o.Height(), color.White)
			}
		}
		MainCamera.DrawRect(screen, Player.X(), Player.Y(), Player.Width(), Player.Height(), color.RGBA{
			A: 0xFF,
			R: 0xFF,
			G: 0x00,
			B: 0x00,
		})
	}

	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %.2f", ebiten.CurrentFPS()))
	return nil
}

func main() {
	if err := ebiten.Run(update, App.Width, App.Height, 1, "Unnamed"); err != nil {
		log.Fatal(err)
	}
}

func applyGravity() {
	isColliding := Player.IntersectsArray(Tiles)
	haveSidewayException := Player.SidewayExceptionArray(Tiles)
	if (!isColliding || haveSidewayException) && Player.HasMass {
		Player.Move(Gravity.Tx, Gravity.Ty)
		Player.IsGrounded = false
	} else if isColliding {
		Player.IsGrounded = true
	}

	if !Coin.IntersectsArray(Tiles) && Coin.HasMass {
		Coin.Options.GeoM.Translate(Gravity.Tx, Gravity.Ty)
	}

	for i, e := range Enemies {
		if !e.IntersectsArray(Tiles) && e.HasMass {
			Enemies[i].Options.GeoM.Translate(Gravity.Tx, Gravity.Ty)
		}
	}

	if Player.IsJumping {
		if jumpFn(TimeDelta) > 0 {
			Player.IsJumping = false
		} else {
			Player.Move(Gravity.Tx, jumpFn(TimeDelta)*Gravity.Ty)
		}
	}
}
