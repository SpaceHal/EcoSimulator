package graphics

import (
	"log"
	"image/color"
	"strconv"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	_ "image/png"
	"bytes"
	"golang.org/x/image/font/gofont/goregular"
)
	
type data struct {
	x			float64
	y 			float64
}

const (
	padding	= 16
	graphHeight = 128+64
	graphWidth = 224
	nPoints = 1000
	yMax = 50
)

var (
	colorGrass = color.RGBA{10,150,10,255}
	colorRabbit = color.RGBA{120,120,120,255}
	colorCat = color.RGBA{0,0,0,255}
	colorFox = color.RGBA{180,50,50,255}
)

var (
	historyRabbits []int = make([]int,nPoints)
	historyCats	   []int = make([]int,nPoints)
	historyFoxes   []int = make([]int,nPoints)
	historyGrass   []int = make([]int,nPoints)
	faceSource 	   *text.GoTextFaceSource
)

func init() {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	if err != nil {
		log.Fatal(err)
	}
	faceSource = s
}


func New(x,y float64) *data {
	u := &data {
		x: x,
		y:y,
	}
	return u
}

func (u *data) drawLine(dst *ebiten.Image, history []int, c color.RGBA) {
	var x,y,x1,y1 float32
	x = float32(u.x+padding)
	y = float32(u.y+padding+graphHeight)-float32(history[0]*graphHeight/yMax)
	for i:=0;i<nPoints;i++ {
		x1=float32(u.x+padding)+float32(i*graphWidth/nPoints)
		y1=float32(u.y+padding+graphHeight)-float32(history[i]*graphHeight/yMax)
		if float64(y1)<u.y+padding {
			y1 = float32(u.y+padding)
		}
		vector.StrokeLine(dst,x,y,x1,y1,2,c,false)
		x=x1
		y=y1
	}
}

func drawText(dst *ebiten.Image, x,y float64, c color.RGBA, str string) {
	op := &text.DrawOptions{}
	op.GeoM.Translate(x, y)
	op.ColorScale.ScaleWithColor(c)
	op.PrimaryAlign = text.AlignEnd
	op.SecondaryAlign = text.AlignCenter
	text.Draw(dst, str, &text.GoTextFace{
		Source: faceSource,
		Size:   14,
	}, op)
}

func (u *data) Draw(dst *ebiten.Image) {
	ebitenutil.DrawRect(dst, u.x+padding,u.y+padding,2,graphHeight,color.RGBA{0x40, 0x40, 0x40, 0xff})
	ebitenutil.DrawRect(dst, u.x+padding,u.y+padding+graphHeight,graphWidth,2,color.RGBA{0x40, 0x40, 0x40, 0xff})

	u.drawLine(dst,historyGrass,colorGrass)
	u.drawLine(dst,historyRabbits,colorRabbit)
	u.drawLine(dst,historyCats,colorCat)
	u.drawLine(dst,historyFoxes,colorFox)
	
	var x,y float64 
	x = u.x+padding+graphWidth
	y = u.y+graphHeight+padding*3
	drawText(dst,x,y,colorGrass,"Grassflächen: "+strconv.Itoa(historyGrass[nPoints-1]))
	drawText(dst,x,y+padding*2,colorRabbit,"Hasen: "+strconv.Itoa(historyRabbits[nPoints-1]))
	drawText(dst,x,y+padding*4,colorCat,"Katzen: "+strconv.Itoa(historyCats[nPoints-1]))
	drawText(dst,x,y+padding*6,colorFox,"Füchse: "+strconv.Itoa(historyFoxes[nPoints-1]))
}

func (u *data) Update(nG,nR,nC,nF int) {
	historyGrass = append(historyGrass[1:],nG)
	historyRabbits = append(historyRabbits[1:],nR)
	historyCats = append(historyCats[1:],nC)
	historyFoxes = append(historyFoxes[1:],nF)
}
