package foxes

import (
	"ecosim/animal";
	"ecosim/world";
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil")

type data struct {
	animal.Animal
}

func New(w *world.World, x, y float64) *data {
	var f *data
	f = new(data)
	(*f).Animal = animal.New(w,x,y)
	return f
}

func (f *data) GetImage () *ebiten.Image {
	var img *ebiten.Image
	var err error
	img, _, err = ebitenutil.NewImageFromFile("foxes/fox.png")
	if err != nil {
		return (*f).Animal.GetImage()
	}
	return img
}
