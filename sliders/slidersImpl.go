package sliders

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/gofont/goregular"
	"log"
	"bytes"
	"image/color"
	"strconv"
)

type data struct {
	x              	float64
	y              	float64
	text           	string
	mouseDown       bool
	onMoved func()
	currentValue    int
	maxValue        int
	active          bool
}

const (
	sliderWidth         = 256
	sliderHeight        = 8
	padding             = 8
	uiFontSize          = 16
)
var (
	uiFaceSource *text.GoTextFaceSource
)

func init() {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	if err != nil {
		log.Fatal(err)
	}
	uiFaceSource = s
}

func New(x,y float64, current,max int, text string, active bool) *data {
	return &data{
		x:       		x,
		y:       		y,
		currentValue: 	current,
		maxValue: 		max,
		text:    		text,
		active: 		active,
	}
}

func (s *data) Update() {
	if s.active && ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		var padding float64 = sliderHeight * 2
		if s.x <= float64(x) && float64(x) < s.x+sliderWidth && s.y-padding <= float64(y) && float64(y) < s.y+sliderHeight+padding {
			s.mouseDown = true
			s.currentValue = int(1 + (float64(x)-s.x)*float64(s.maxValue)/sliderWidth)
			s.onMoved()
		} else {
			s.mouseDown = false
		}
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

func (s *data) Draw(dst *ebiten.Image) {
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
	drawText(dst, s.x+sliderWidth+16, s.y+sliderHeight/2, uiFontSize, s.text+strconv.Itoa(n))
}

func (s *data) GetValue() int {
	return s.currentValue
}

func (s *data) SetValue(v int) {
	s.currentValue = v
}

func (s *data) SetActive(a bool) {
	s.active = a
}

func (s *data) SetOnMoved(f func ()) {
	s.onMoved = f
}
	
