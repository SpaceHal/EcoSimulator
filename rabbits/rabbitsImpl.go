package rabbits

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

func New(w *world.World, x, y float64) *data {
	var f *data
	f = new(data)
	(*f).Animal = animal.New(w, x, y)
	f.makeAnimal()
	f.SetEnergyLoss(((rand.Float64()*30 + 60) * 60))

	return f
}

func (f *data) makeAnimal() *ebiten.Image {
	var img *ebiten.Image
	var err error
	img, _, err = ebitenutil.NewImageFromFile("rabbits/rabbit.png")
	if err == nil {
		//f.SetImage(img)
	}
	return img
}
