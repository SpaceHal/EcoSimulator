package checkboxes

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/gofont/goregular"
	"log"
	"bytes"
	"image/color"
)

type data struct {
	x              float64
	y              float64
	text           string
	checked        bool
	mouseDown      bool
	onClicked      func()
}

const (
	checkboxWidth       = 16
	checkboxHeight      = 16
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

func New(x,y float64, text string, checked bool) *data {
	return &data{
		x:       x,
		y:       y,
		text:    text,
		checked: checked,
	}
}

func (c *data) Update() {
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
			if c.onClicked != nil {
				c.onClicked()
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

func (c *data) Draw(dst *ebiten.Image) {
	offset := float64(checkboxWidth) / 2
	ebitenutil.DrawCircle(dst, c.x+offset, c.y+offset, offset, color.RGBA{255, 255, 255, 255})
	if c.checked {
		ebitenutil.DrawCircle(dst, c.x+offset, c.y+offset, offset-2, color.RGBA{0x80, 0x80, 0x80, 0xff})
	}

	drawText(dst, c.x+checkboxWidth+padding, c.y+checkboxHeight/2, uiFontSize, c.text)
}

func (c *data) IsChecked() bool {
	return c.checked
}

func (c *data) SetOnClicked(f func ()) {
	c.onClicked = f
}
	
