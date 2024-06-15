package ui

import (
	. "ecosim/config"
	"ecosim/checkboxes"
	"ecosim/sliders"
	"bytes"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/images"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/gofont/goregular"
	"image"
	"image/color"
	_ "image/png"
	"log"
)

type data struct {
	nCats 		int
	nBunnies  	int
	nFoxes		int
	nGrass  	int
	cbGrass   	checkboxes.Checkbox
	cbBunnies 	checkboxes.Checkbox
	cbCats    	checkboxes.Checkbox
	cbFoxes   	checkboxes.Checkbox
	sGrass    	sliders.Slider
	sBunnies  	sliders.Slider
	sCats     	sliders.Slider
	sFoxes    	sliders.Slider
}

var (
	uiImage         *ebiten.Image
	uiFaceSource    *text.GoTextFaceSource
)

const (
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

func New() *data {
	u := &data{
		nGrass:   GrassStartNumber,
		nBunnies: BunnyStartNumber,
		nCats:    CatStartNumber,
		nFoxes:   FoxStartNumber,
	}

	u.cbGrass = checkboxes.New(leftIndent,lineSpacingInPixels * 3,"Grass",true)
	u.cbGrass.SetOnClicked(func() {
		if u.cbGrass.IsChecked() {
			u.nGrass = GrassStartNumber
			u.sGrass.SetActive(true)
		} else {
			u.nGrass = 0
			u.sGrass.SetActive(false)
		}
	})

	u.sGrass = sliders.New(leftIndent,lineSpacingInPixels * 4,GrassStartNumber,100,"Anzahl Grassflächen: ",true)
	u.sGrass.SetOnMoved(func() {
		u.nGrass = u.sGrass.GetValue()
	})

	u.cbBunnies = checkboxes.New(leftIndent,lineSpacingInPixels * 5,"Hasen",true)
	u.cbBunnies.SetOnClicked(func() {
		if u.cbBunnies.IsChecked() {
			u.nBunnies = BunnyStartNumber
			u.sBunnies.SetActive(true)
		} else {
			u.nBunnies = 0
			u.sBunnies.SetActive(false)
		}
	})

	u.sBunnies = sliders.New(leftIndent,lineSpacingInPixels * 6,BunnyStartNumber,50,"Anzahl Hasen: ",true)
	u.sBunnies.SetOnMoved(func() {
		u.nBunnies = u.sBunnies.GetValue()
	})

	u.cbCats = checkboxes.New(leftIndent,lineSpacingInPixels * 7,"Katzen",true)
	u.cbCats.SetOnClicked(func() {
		if u.cbCats.IsChecked() {
			u.nCats = CatStartNumber
			u.sCats.SetActive(true)
		} else {
			u.nCats = 0
			u.sCats.SetActive(false)
		}
	})

	u.sCats = sliders.New(leftIndent,lineSpacingInPixels * 8,CatStartNumber,50,"Anzahl Katzen: ",true)
	u.sCats.SetOnMoved(func() {
		u.nCats = u.sCats.GetValue()
	})

	u.cbFoxes = checkboxes.New(leftIndent,lineSpacingInPixels * 9,"Fuechse",true)
	u.cbFoxes.SetOnClicked(func() {
		if u.cbFoxes.IsChecked() {
			u.nFoxes = FoxStartNumber
			u.sFoxes.SetActive(true)
		} else {
			u.nFoxes = 0
			u.sFoxes.SetActive(false)
		}
	})

	u.sFoxes = sliders.New(leftIndent,lineSpacingInPixels * 10,FoxStartNumber,30,"Anzahl Füchse: ",true)
	u.sFoxes.SetOnMoved(func() {
		u.nFoxes = u.sFoxes.GetValue()
	})

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
	u.cbGrass.Draw(dst)
	u.sGrass.Draw(dst)
	u.cbBunnies.Draw(dst)
	u.sBunnies.Draw(dst)
	u.cbCats.Draw(dst)
	u.sCats.Draw(dst)
	u.cbFoxes.Draw(dst)
	u.sFoxes.Draw(dst)
}

func (u *data) Update() {
	u.cbGrass.Update()
	u.sGrass.Update()
	u.cbBunnies.Update()
	u.sBunnies.Update()
	u.cbCats.Update()
	u.sCats.Update()
	u.cbFoxes.Update()
	u.sFoxes.Update()
}
