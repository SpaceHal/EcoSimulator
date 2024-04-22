package main

import (
	"ecosim/animal"
	"ecosim/world"
	"ecosim/foxes"
	"fmt"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var (
	bunnies    []animal.Animal
	welt       world.World
	tilesImage *ebiten.Image
	waterImage *ebiten.Image
)

const (
	NumberOfBunnies = 20
	screenWidth     = 20 * 16 * 3
	screenHeight    = 20 * 16 * 3
)

type Game struct {
	counter int
	//layers  [][]int // Tile map

	//grid bool
	//aa   bool
	//line bool
}

func (g *Game) Update() error {
	g.counter++

	// Kachel-Gitter anzeigen
	if inpututil.IsKeyJustPressed(ebiten.KeyG) {
		welt.ToggleGrid()
	}

	// Kachel-Gitter anzeigen
	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		welt.ToggleDebug()
	}

	// Mausposition einlesen
	mx, my := ebiten.CursorPosition()
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		welt.ToggleGround(mx, my)
	}

	for _, b := range bunnies {
		//b.SeeOthers(bunnies[:])

		b.Update(bunnies) // Position neu bestimmen
	}
	return nil
}

func (g *Game) Draw(dst *ebiten.Image) {

	welt.Draw(dst, g.counter)

	for _, b := range bunnies {
		//b.Separate(bunnies[:])
		b.Draw(dst, b.GetImage()) // Ein Jäger
	}
	// Text im Fenster
	msg := fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f", ebiten.ActualTPS(), ebiten.ActualFPS())
	ebitenutil.DebugPrint(dst, msg)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	g := &Game{
		counter: 0,
	}

	welt = world.New(screenWidth, screenHeight, tilesImage)

	bunnies = make([]animal.Animal, NumberOfBunnies)
	for i := 0; i < NumberOfBunnies; i++ {
		bunnies[i] = foxes.New(&welt, (rand.Float64()/2+0.5)*screenWidth/2, (rand.Float64()/2+0.5)*screenHeight/2)
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("EcoSim")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
