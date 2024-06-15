package main

import (
	"ecosim/cats"
	"ecosim/foxes"
	"ecosim/graphics"
	"ecosim/grass"
	"ecosim/rabbits"
	"ecosim/ui"
	"ecosim/world"

	"bytes"
	"image"
	"image/color"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2/examples/resources/images"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/gofont/goregular"

	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var (
	katzen       []cats.Cat       //[]entity.Entity
	bunnies      []rabbits.Rabbit //[]entity.Entity
	fuechse      []foxes.Fox      //[]entity.Entity
	food         []grass.Grass    //[]entity.Entity
	welt         world.World
	userInt      ui.UI
	grafik       graphics.Graphics
	tilesImage   *ebiten.Image
	waterImage   *ebiten.Image
	gameRunning  bool
	uiImage      *ebiten.Image
	uiFaceSource *text.GoTextFaceSource
)

func init() {
	// Decode an image from the image file's byte slice.
	img, _, err := image.Decode(bytes.NewReader(images.UI_png))
	if err != nil {
		log.Fatal(err)
	}
	uiImage = ebiten.NewImageFromImage(img)

	s, err := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	if err != nil {
		log.Fatal(err)
	}
	uiFaceSource = s
}

const (
	uiFontSize          = 16
	lineSpacingInPixels = 16
)

const (
	graphicsWidth = 256
	buttonHeight  = 32
	buttonPadding = 8
	buttonWidth   = graphicsWidth - 2*buttonPadding
	scale         = 2.5
	screenWidth   = 20*16*scale + graphicsWidth
	screenHeight  = 20 * 16 * scale
)

type Game struct {
	counter int
	button  *Button
}

type Input struct {
	mouseButtonState int
}

type Button struct {
	x         float64
	y         float64
	text      string
	mouseDown bool
	onPressed func(b *Button)
}

func (b *Button) Update() {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if b.x <= float64(x) && float64(x) < b.x+buttonWidth && b.y <= float64(y) && float64(y) < b.y+buttonHeight {
			b.mouseDown = true
		} else {
			b.mouseDown = false
		}
	} else {
		if b.mouseDown {
			if b.onPressed != nil {
				b.onPressed(b)
			}
		}
		b.mouseDown = false
	}
}

func (b *Button) Draw(dst *ebiten.Image) {
	ebitenutil.DrawRect(dst, b.x, b.y, buttonWidth, buttonHeight, color.RGBA{0x80, 0x80, 0x80, 0xff})
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(b.x+(buttonWidth/2)), float64(b.y+buttonHeight/2))
	op.ColorScale.ScaleWithColor(color.Black)
	op.LineSpacing = 16
	op.PrimaryAlign = text.AlignCenter
	op.SecondaryAlign = text.AlignCenter
	text.Draw(dst, b.text, &text.GoTextFace{
		Source: uiFaceSource,
		Size:   uiFontSize,
	}, op)
}

func (b *Button) SetOnPressed(f func(b *Button)) {
	b.onPressed = f
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

	if gameRunning && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		ok := welt.IsLand(int(mx), int(my))
		if ok {
			food = append(food, grass.New(&welt))
		}
	}
	if gameRunning && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		welt.ToggleGround(mx, my)
	}

	g.button.Update()

	if !gameRunning {
		userInt.Update()
		return nil
	}

	// Grass löschen ...
	eoa := len(food)
	if eoa > 0 {
		for i := 0; i < eoa; i++ {
			if !food[i].IsAlive() {
				food = append(food[:i], food[i+1:]...)
				eoa--
			}
		}
	}
	// neues Grass
	if g.counter%30 == 0 && g != nil {
		food = append(food, grass.New(&welt))
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
			newRabbit := bunny.Update(&bunnies, &food)
			if newRabbit != nil {
				livingRabbits = append(livingRabbits, newRabbit)
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

	grafik.Update(len(food), len(bunnies), len(katzen), len(fuechse))
	return nil
}

func drawGame(g *Game, dst *ebiten.Image) {
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
}

func (g *Game) Draw(dst *ebiten.Image) {
	ebitenutil.DrawRect(dst, screenWidth-graphicsWidth, 0, graphicsWidth, screenHeight, color.RGBA{0xff, 0xff, 0xff, 0xff})
	grafik.Draw(dst)
	g.button.Draw(dst)

	if !gameRunning {
		userInt.Draw(dst)
	} else {
		drawGame(g, dst)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func resetGrass() {
	n := userInt.GetNumberOfGrass()
	food = make([]grass.Grass, n)
	for i := 0; i < n; i++ {
		food[i] = grass.New(&welt)
	}
}

func resetBunnies() {
	n := userInt.GetNumberOfBunnies()
	bunnies = make([]rabbits.Rabbit, n)
	for i := 0; i < n; i++ {
		bunnies[i] = rabbits.New(&welt)
	}
}

func resetCats() {
	n := userInt.GetNumberOfCats()
	katzen = make([]cats.Cat, n)
	for i := 0; i < n; i++ {
		katzen[i] = cats.New(&welt)
	}
}

func resetFoxes() {
	n := userInt.GetNumberOfFoxes()
	fuechse = make([]foxes.Fox, n)
	for i := 0; i < n; i++ {
		fuechse[i] = foxes.New(&welt)
	}
}

func main() {
	g := &Game{
		counter: 0,
	}

	welt = world.New(screenWidth-graphicsWidth, screenHeight, scale, tilesImage)
	userInt = ui.New()
	grafik = graphics.New(screenWidth-graphicsWidth, 0)

	resetGrass()
	resetBunnies()
	resetCats()
	resetFoxes()

	g.button = &Button{
		x:    screenWidth - buttonWidth - buttonPadding,
		y:    screenHeight - buttonPadding - buttonHeight,
		text: "Start",
	}

	g.button.SetOnPressed(func(b *Button) {
		resetGrass()
		resetBunnies()
		resetCats()
		resetFoxes()
		gameRunning = !gameRunning
		if gameRunning {
			b.text = "Stop"
		} else {
			b.text = "Start"
		}
	})

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("EcoSim")

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
