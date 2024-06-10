package ui

import (
	"bytes"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/images"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/gofont/goregular"
	"image"
	"image/color"
	_ "image/png"
	"log"
	"strconv"
)

type data struct {
	nCats           int
	nBunnies        int
	nFoxes          int
	nGrass          int
	checkBoxGrass   *CheckBox
	checkBoxBunnies *CheckBox
	checkBoxCats    *CheckBox
	checkBoxFoxes   *CheckBox
	sliderGrass     *Slider
	sliderBunnies   *Slider
	sliderCats      *Slider
	sliderFoxes     *Slider
}

var (
	NumberOfCats    = 5
	NumberOfBunnies = 10
	NumberOfFoxes   = 5
	NumberOfGrass   = 20
	uiImage         *ebiten.Image
	uiFaceSource    *text.GoTextFaceSource
)

const (
	checkboxWidth       = 16
	checkboxHeight      = 16
	sliderWidth         = 256
	sliderHeight        = 8
	padding             = 8
	uiFontSize          = 16
	lineSpacingInPixels = 60
	leftIndent          = 32
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

type Slider struct {
	x               float64
	y               float64
	textBase        string
	mouseDown       bool
	onSliderChanged func(s *Slider)
	currentValue    int
	maxValue        int
	active          bool
}

type CheckBox struct {
	x              float64
	y              float64
	text           string
	checked        bool
	mouseDown      bool
	onCheckChanged func(c *CheckBox)
}

func (c *CheckBox) Update() {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if c.x <= float64(x) && float64(x) < c.x+checkboxWidth+padding && c.y <= float64(y) && float64(y) < c.y+checkboxHeight {
			c.mouseDown = true
		} else {
			c.mouseDown = false
		}
	} else {
		if c.mouseDown {
			c.checked = !c.checked
			if c.onCheckChanged != nil {
				c.onCheckChanged(c)
			}
		}
		c.mouseDown = false
	}
}

func drawText(dst *ebiten.Image, x, y float64, size float64, str string) {
	op := &text.DrawOptions{}
	op.GeoM.Translate(x, y)
	op.ColorScale.ScaleWithColor(color.White)
	op.PrimaryAlign = text.AlignStart
	op.SecondaryAlign = text.AlignCenter
	text.Draw(dst, str, &text.GoTextFace{
		Source: uiFaceSource,
		Size:   size,
	}, op)
}

func (c *CheckBox) Draw(dst *ebiten.Image) {
	offset := float64(checkboxWidth) / 2
	ebitenutil.DrawCircle(dst, c.x+offset, c.y+offset, offset, color.RGBA{255, 255, 255, 255})
	if c.checked {
		ebitenutil.DrawCircle(dst, c.x+offset, c.y+offset, offset-2, color.RGBA{0x80, 0x80, 0x80, 0xff})
	}

	drawText(dst, c.x+checkboxWidth+padding, c.y+checkboxHeight/2, uiFontSize, c.text)
}

func (s *Slider) Update() {
	if s.active && ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		var padding float64 = sliderHeight * 2
		if s.x <= float64(x) && float64(x) < s.x+sliderWidth && s.y-padding <= float64(y) && float64(y) < s.y+sliderHeight+padding {
			s.mouseDown = true
			s.currentValue = int(1 + (float64(x)-s.x)*float64(s.maxValue)/sliderWidth)
			s.onSliderChanged(s)
		} else {
			s.mouseDown = false
		}
	}
}

func (s *Slider) Draw(dst *ebiten.Image) {
	var n int
	var colorSlider, colorHandle color.Color
	if s.active {
		n = s.currentValue
		colorSlider = color.RGBA{0x80, 0x80, 0x80, 0xff}
		colorHandle = color.RGBA{255, 255, 255, 255}
	} else {
		n = 0
		colorSlider = color.RGBA{0x40, 0x40, 0x40, 0xff}
		colorHandle = color.RGBA{0x40, 0x40, 0x40, 0xff}
	}
	ebitenutil.DrawRect(dst, s.x, s.y, sliderWidth, sliderHeight, colorSlider)
	xPos := s.x + float64(n*sliderWidth/s.maxValue)
	ebitenutil.DrawCircle(dst, xPos, s.y+sliderHeight/2, sliderHeight, colorHandle)
	drawText(dst, s.x+sliderWidth+16, s.y+sliderHeight/2, uiFontSize, s.textBase+strconv.Itoa(n))
}

func New() *data {
	u := &data{
		nGrass:   NumberOfGrass,
		nBunnies: NumberOfBunnies,
		nCats:    NumberOfCats,
		nFoxes:   NumberOfFoxes,
	}

	u.checkBoxGrass = &CheckBox{
		x:       leftIndent,
		y:       lineSpacingInPixels * 3,
		text:    "Grass",
		checked: true,
	}

	u.checkBoxGrass.onCheckChanged = func(c *CheckBox) {
		if c.checked {
			u.nGrass = NumberOfGrass
			u.sliderGrass.active = true
		} else {
			u.nGrass = 0
			u.sliderGrass.active = false
		}
	}

	u.sliderGrass = &Slider{
		x:            leftIndent,
		y:            lineSpacingInPixels * 4,
		maxValue:     100,
		textBase:     "Anzahl Grassflächen: ",
		currentValue: NumberOfGrass,
		active:       true,
	}

	u.sliderGrass.onSliderChanged = func(s *Slider) {
		u.nGrass = s.currentValue
	}

	u.checkBoxBunnies = &CheckBox{
		x:       leftIndent,
		y:       lineSpacingInPixels * 5,
		text:    "Hasen",
		checked: true,
	}

	u.checkBoxBunnies.onCheckChanged = func(c *CheckBox) {
		if c.checked {
			u.nBunnies = NumberOfBunnies
			u.sliderBunnies.active = true
		} else {
			u.nBunnies = 0
			u.sliderBunnies.active = false
		}
	}

	u.sliderBunnies = &Slider{
		x:            leftIndent,
		y:            lineSpacingInPixels * 6,
		maxValue:     50,
		textBase:     "Anzahl Hasen: ",
		currentValue: NumberOfBunnies,
		active:       true,
	}

	u.sliderBunnies.onSliderChanged = func(s *Slider) {
		u.nBunnies = s.currentValue
	}

	u.checkBoxCats = &CheckBox{
		x:       leftIndent,
		y:       lineSpacingInPixels * 7,
		text:    "Katzen",
		checked: true,
	}

	u.checkBoxCats.onCheckChanged = func(c *CheckBox) {
		if c.checked {
			u.nCats = NumberOfCats
			u.sliderCats.active = true
		} else {
			u.nCats = 0
			u.sliderCats.active = false
		}
	}

	u.sliderCats = &Slider{
		x:            leftIndent,
		y:            lineSpacingInPixels * 8,
		maxValue:     50,
		textBase:     "Anzahl Katzen: ",
		currentValue: NumberOfCats,
		active:       true,
	}

	u.sliderCats.onSliderChanged = func(s *Slider) {
		u.nCats = s.currentValue
	}

	u.checkBoxFoxes = &CheckBox{
		x:       leftIndent,
		y:       lineSpacingInPixels * 9,
		text:    "Fuechse",
		checked: true,
	}

	u.checkBoxFoxes.onCheckChanged = func(c *CheckBox) {
		if c.checked {
			u.nFoxes = NumberOfFoxes
			u.sliderFoxes.active = true
		} else {
			u.nFoxes = 0
			u.sliderFoxes.active = false
		}
	}

	u.sliderFoxes = &Slider{
		x:            leftIndent,
		y:            lineSpacingInPixels * 10,
		maxValue:     30,
		textBase:     "Anzahl Füchse: ",
		currentValue: NumberOfFoxes,
		active:       true,
	}

	u.sliderFoxes.onSliderChanged = func(s *Slider) {
		u.nFoxes = s.currentValue
	}

	return u
}

func (u *data) GetNumberOfGrass() int {
	return u.nGrass
}

func (u *data) GetNumberOfBunnies() int {
	return u.nBunnies
}

func (u *data) GetNumberOfCats() int {
	return u.nCats
}

func (u *data) GetNumberOfFoxes() int {
	return u.nFoxes
}

func (u *data) Draw(dst *ebiten.Image) {
	drawText(dst, leftIndent, uiFontSize*1.5, uiFontSize*1.5, "Einstellungen")
	var beschreibung [3]string = [3]string{
		"Dieses Programm simuliert die Jäger-Beute-Beziehung zwischen verschiedenen",
		"Tieren und Pflanzen. Die einzelnen Komponenten können aktiviert oder",
		"desaktiviert werden. Die Anfangspopulationen können festgelegt werden."}
	for i := 0; i < len(beschreibung); i++ {
		drawText(dst, leftIndent, lineSpacingInPixels+float64(i*uiFontSize), uiFontSize, beschreibung[i])
	}
	u.checkBoxGrass.Draw(dst)
	u.sliderGrass.Draw(dst)
	u.checkBoxBunnies.Draw(dst)
	u.sliderBunnies.Draw(dst)
	u.checkBoxCats.Draw(dst)
	u.sliderCats.Draw(dst)
	u.checkBoxFoxes.Draw(dst)
	u.sliderFoxes.Draw(dst)
}

func (u *data) Update() {
	u.checkBoxGrass.Update()
	u.sliderGrass.Update()
	u.checkBoxBunnies.Update()
	u.sliderBunnies.Update()
	u.checkBoxCats.Update()
	u.sliderCats.Update()
	u.checkBoxFoxes.Update()
	u.sliderFoxes.Update()
}
