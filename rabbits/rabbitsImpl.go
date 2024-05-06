package rabbits

import (
	"ecosim/animal"
	"ecosim/world"
	"math/rand"
)

type data struct {
	animal.Animal
}

func New(w *world.World, x, y float64) *data {
	var f *data
	f = new(data)
	(*f).Animal = animal.New(w, x, y)
	
	f.SetImageFromFile("rabbits/rabbit.png")
	f.SetEnergyLoss(((rand.Float64()*30 + 60) * 60))

	return f
}

