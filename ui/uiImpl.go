package ui

import (
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
	nCats           int
	nBunnies        int
	nFoxes          int
	nGrass          int
	checkBoxGrass   checkboxes.Checkbox
	checkBoxBunnies checkboxes.Checkbox
	checkBoxCats    checkboxes.Checkbox
	checkBoxFoxes   checkboxes.Checkbox
	sliderGrass     sliders.Slider
	sliderBunnies   sliders.Slider
	sliderCats      sliders.Slider
	sliderFoxes     sliders.Slider
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
		nGrass:   NumberOfGrass,
		nBunnies: NumberOfBunnies,
		nCats:    NumberOfCats,
		nFoxes:   NumberOfFoxes,
	}

	u.checkBoxGrass = checkboxes.New(leftIndent,lineSpacingInPixels * 3,"Grass",true)
	u.checkBoxGrass.SetOnClicked(func() {
		if u.checkBoxGrass.IsChecked() {
			u.nGrass = NumberOfGrass
			u.sliderGrass.SetActive(true)
		} else {
			u.nGrass = 0
			u.sliderGrass.SetActive(false)
		}
	})

	u.sliderGrass = sliders.New(leftIndent,lineSpacingInPixels * 4,NumberOfGrass,100,"Anzahl Grassflächen: ",true)
	u.sliderGrass.SetOnMoved(func() {
		u.nGrass = u.sliderGrass.GetValue()
	})

	u.checkBoxBunnies = checkboxes.New(leftIndent,lineSpacingInPixels * 5,"Hasen",true)
	u.checkBoxBunnies.SetOnClicked(func() {
		if u.checkBoxBunnies.IsChecked() {
			u.nBunnies = NumberOfBunnies
			u.sliderBunnies.SetActive(true)
		} else {
			u.nBunnies = 0
			u.sliderBunnies.SetActive(false)
		}
	})

	u.sliderBunnies = sliders.New(leftIndent,lineSpacingInPixels * 6,NumberOfBunnies,50,"Anzahl Hasen: ",true)
	u.sliderBunnies.SetOnMoved(func() {
		u.nBunnies = u.sliderBunnies.GetValue()
	})

	u.checkBoxCats = checkboxes.New(leftIndent,lineSpacingInPixels * 7,"Katzen",true)
	u.checkBoxCats.SetOnClicked(func() {
		if u.checkBoxCats.IsChecked() {
			u.nCats = NumberOfCats
			u.sliderCats.SetActive(true)
		} else {
			u.nCats = 0
			u.sliderCats.SetActive(false)
		}
	})

	u.sliderCats = sliders.New(leftIndent,lineSpacingInPixels * 8,NumberOfGrass,50,"Anzahl Katzen: ",true)
	u.sliderCats.SetOnMoved(func() {
		u.nCats = u.sliderCats.GetValue()
	})

	u.checkBoxFoxes = checkboxes.New(leftIndent,lineSpacingInPixels * 9,"Fuechse",true)
	u.checkBoxFoxes.SetOnClicked(func() {
		if u.checkBoxFoxes.IsChecked() {
			u.nFoxes = NumberOfFoxes
			u.sliderFoxes.SetActive(true)
		} else {
			u.nFoxes = 0
			u.sliderFoxes.SetActive(false)
		}
	})

	u.sliderFoxes = sliders.New(leftIndent,lineSpacingInPixels * 10,NumberOfFoxes,30,"Anzahl Füchse: ",true)
	u.sliderFoxes.SetOnMoved(func() {
		u.nFoxes = u.sliderFoxes.GetValue()
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
