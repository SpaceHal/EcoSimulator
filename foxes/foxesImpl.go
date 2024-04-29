package foxes

import (
	"ecosim/animal"
	"ecosim/world"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type data struct {
	animal.Animal
}

func New(w *world.World, x, y float64, preys *[]animal.Animal) *data {
	var f *data
	f = new(data)
	(*f).Animal = animal.New(w, x, y)
	f.makeAnimal()
	f.SetEnergyLoss(((rand.Float64()*30 + 60) * 60))
	f.SetPreys(preys)
	f.SetMaxVel(0.7)
	f.SetViewMag(200)

	return f
}

func (f *data) makeAnimal() *ebiten.Image {
	var img *ebiten.Image
	var err error
	img, _, err = ebitenutil.NewImageFromFile("foxes/fox.png")
	if err == nil {
		//f.SetImage(img)
		f.SetColorRGB(50, 150, 250)
	}
	return img
}
