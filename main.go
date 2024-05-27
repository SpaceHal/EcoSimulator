package main

import (
	"ecosim/cats"
	"ecosim/foxes"
	"ecosim/grass"
	"ecosim/rabbits"
	"ecosim/world"
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var (
	katzen     []cats.Cat       //[]animal.Animal
	bunnies    []rabbits.Rabbit //[]animal.Animal
	fuechse    []foxes.Fox      //[]animal.Animal
	food       []grass.Grass    //[]animal.Animal
	welt       world.World
	tilesImage *ebiten.Image
	//waterImage *ebiten.Image
)

const (
	NumberOfCats    = 5
	NumberOfBunnies = 10
	NumberOfFoxes   = 5
	NumberOfGrass   = 20
	screenWidth     = 20 * 16 * 2
	screenHeight    = 20 * 16 * 2
)

type Game struct {
	counter int
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
	if g.counter%30 == 0 && g != nil {
		food = append(food, grass.New(&welt))
	}

	for _, b := range food {
		b.Update() // Position neu bestimmen
	}

	// Alle Katzen aktualisieren
	var livingKatzen []cats.Cat
	for _, katze := range katzen {
		if katze.IsAlive() {
			livingKatzen = append(livingKatzen, katze)
			neueKatze := katze.Update(&katzen, &bunnies, &food)
			if neueKatze != nil {
				livingKatzen = append(livingKatzen, neueKatze)
			}
		}
	}
	katzen = livingKatzen

	// Alle Hasen aktualisieren
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

	// Alle Füchse aktualisieren
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

	for _, f := range food {
		f.Draw(dst)
	}

	for _, b := range bunnies {
		b.Draw(dst)
	}

	for _, k := range katzen {
		k.Draw(dst)
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

	katzen = make([]cats.Cat, NumberOfCats)
	for i := 0; i < NumberOfCats; i++ {
		katzen[i] = cats.New(&welt)
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
