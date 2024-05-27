package main

import (
	"ecosim/foxes"
	"ecosim/grass"
	"ecosim/rabbits"
	"ecosim/world"
	"fmt"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var (
	bunnies    []rabbits.Rabbit //[]animal.Animal
	fuechse    []foxes.Fox      //[]animal.Animal
	food       []grass.Grass    //[]animal.Animal
	welt       world.World
	tilesImage *ebiten.Image
	//waterImage *ebiten.Image
)

const (
	NumberOfBunnies = 10
	NumberOfFoxes   = 5
	NumberOfGrass   = 20
	screenWidth     = 20 * 16 * 3
	screenHeight    = 20 * 16 * 3
)

type Game struct {
	counter int
}

func randBetween(a, b float64) float64 {
	return a + rand.Float64()*b
}

func (g *Game) Update() error {
	g.counter++

	if inpututil.IsKeyJustPressed(ebiten.KeyG) {
		welt.ToggleGrid()
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		welt.ToggleDebug()
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		welt.ToggleStatistics()
	}

	mx, my := ebiten.CursorPosition()
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		ok := welt.IsLand(int(mx), int(my))
		if ok {
			food = append(food, grass.New(&welt))
		}
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		welt.ToggleGround(mx, my)
	}

	// Grass löschen ...
	eoa := len(food)
	for i := 0; i < eoa; i++ {
		if !food[i].IsAlive() {
			food = append(food[:i], food[i+1:]...)
			eoa--
			//fmt.Println("Ein Karotte weniger.", eoa, "leben noch.")
		}
	}

	// neues Grass
	if g.counter%30 == 0 {
		food = append(food, grass.New(&welt))
	}

	for _, b := range food {
		b.Update() // Position neu bestimmen
	}

	// Alle toten Hasen löschen ...
	var livingRabbits []rabbits.Rabbit
	for _, bunny := range bunnies {
		if bunny.IsAlive() {
			livingRabbits = append(livingRabbits, bunny)
			newFuchs := bunny.Update(&bunnies, &food)
			if newFuchs != nil {
				livingRabbits = append(livingRabbits, newFuchs)
			}
		}
	}
	bunnies = livingRabbits

	// Alle toten Füchse löschen ...
	var livingFuechse []foxes.Fox
	for _, fuchs := range fuechse {
		if fuchs.IsAlive() {
			livingFuechse = append(livingFuechse, fuchs)
			newFuchs := fuchs.Update(&fuechse, &bunnies)
			if newFuchs != nil {
				livingFuechse = append(livingFuechse, newFuchs)
			}
		}
	}
	fuechse = livingFuechse

	return nil
}

func (g *Game) Draw(dst *ebiten.Image) {

	welt.Draw(dst, g.counter)

	for _, gr := range food {
		gr.Draw(dst)
	}

	for _, b := range bunnies {
		b.Draw(dst)
	}

	for _, f := range fuechse {
		f.Draw(dst)
	}

	// Text im Fenster
	msg := fmt.Sprintf("FPS: %0.2f\n Essen:\t %d \n Hasen:\t%d \n Füchse:\t%d ",
		ebiten.ActualFPS(), len(food), len(bunnies), len(fuechse))
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

	food = make([]grass.Grass, NumberOfGrass)
	for i := 0; i < NumberOfGrass; i++ {
		food[i] = grass.New(&welt)
	}

	bunnies = make([]rabbits.Rabbit, NumberOfBunnies)
	for i := 0; i < NumberOfBunnies; i++ {
		bunnies[i] = rabbits.New(&welt)
	}

	fuechse = make([]foxes.Fox, NumberOfFoxes)
	for i := 0; i < NumberOfFoxes; i++ {
		fuechse[i] = foxes.New(&welt)
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("EcoSim")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
